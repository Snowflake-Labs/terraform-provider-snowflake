package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// AlertBuilder abstracts the creation of sql queries for a snowflake alert.
type AlertBuilder struct {
	name      string
	db        string
	schema    string
	warehouse string
	schedule  string
	comment   string
	condition string
	action    string
	disabled  bool
}

// GetFullName prepends db and schema to in parameter.
func (tb *AlertBuilder) GetFullName(name string) string {
	var n strings.Builder
	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, tb.db, tb.schema, name))
	return n.String()
}

// QualifiedName prepends the db and schema and escapes everything nicely.
func (tb *AlertBuilder) QualifiedName() string {
	return tb.GetFullName(tb.name)
}

// Name returns the name of the alert.
func (tb *AlertBuilder) Name() string {
	return tb.name
}

// WithWarehouse adds a warehouse to the AlertBuilder.
func (tb *AlertBuilder) WithWarehouse(s string) *AlertBuilder {
	tb.warehouse = s
	return tb
}

// WithSchedule adds a schedule to the AlertBuilder.
func (tb *AlertBuilder) WithSchedule(s string) *AlertBuilder {
	tb.schedule = s
	return tb
}

// WithComment adds a comment to the AlertBuilder.
func (tb *AlertBuilder) WithComment(c string) *AlertBuilder {
	tb.comment = c
	return tb
}

// WithCondition adds a condition to the AlertBuilder.
func (tb *AlertBuilder) WithCondition(condition string) *AlertBuilder {
	tb.condition = condition
	return tb
}

// WithAction adds a sql statement to the AlertBuilder.
func (tb *AlertBuilder) WithAction(action string) *AlertBuilder {
	tb.action = action
	return tb
}

// Alert returns a pointer to a Builder that abstracts the DDL operations for a alert.
//
// Supported DDL operations are:
//   - CREATE ALERT
//   - ALTER ALERT
//   - DROP ALERT
//   - DESCRIBE ALERT
//   - SHOW ALERTS
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/alerts)

func NewAlertBuilder(name, db, schema string) *AlertBuilder {
	return &AlertBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL that will create a new alert.
func (tb *AlertBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` ALERT %v`, tb.QualifiedName()))
	q.WriteString(fmt.Sprintf(` WAREHOUSE = "%v"`, EscapeString(tb.warehouse)))
	q.WriteString(fmt.Sprintf(` SCHEDULE = '%v'`, EscapeString(tb.schedule)))

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	q.WriteString(fmt.Sprintf(` IF (EXISTS ( %v ))`, EscapeString(tb.condition)))
	q.WriteString(fmt.Sprintf(` THEN %v`, EscapeString(tb.action)))

	return q.String()
}

// ChangeWarehouse returns the sql that will change the warehouse for the alert.
func (tb *AlertBuilder) ChangeWarehouse(newWh string) string {
	return fmt.Sprintf(`ALTER alert %v SET WAREHOUSE = "%v"`, tb.QualifiedName(), EscapeString(newWh))
}

// ChangeSchedule returns the sql that will change the schedule for the alert.
func (tb *AlertBuilder) ChangeSchedule(newSchedule string) string {
	return fmt.Sprintf(`ALTER alert %v SET SCHEDULE = '%v'`, tb.QualifiedName(), EscapeString(newSchedule))
}

// RemoveSchedule returns the sql that will remove the schedule for the alert.
func (tb *AlertBuilder) RemoveSchedule() string {
	return fmt.Sprintf(`ALTER alert %v UNSET SCHEDULE`, tb.QualifiedName())
}

// ChangeComment returns the sql that will change the comment for the alert.
func (tb *AlertBuilder) ChangeComment(newComment string) string {
	return fmt.Sprintf(`ALTER alert %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(newComment))
}

// RemoveComment returns the sql that will remove the comment for the alert.
func (tb *AlertBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER alert %v UNSET COMMENT`, tb.QualifiedName())
}

// ChangeCondition returns the sql that will update the WHEN condition for the alert.
func (tb *AlertBuilder) ChangeCondition(newCondition string) string {
	return fmt.Sprintf(`ALTER ALERT %v MODIFY CONDITION EXISTS ( %v )`, tb.QualifiedName(), newCondition)
}

// ChangeAction returns the sql that will update the sql the alert executes.
func (tb *AlertBuilder) ChangeAction(newAction string) string {
	return fmt.Sprintf(`ALTER ALERT %v MODIFY ACTION %v`, tb.QualifiedName(), UnescapeString(newAction))
}

// Suspend returns the sql that will suspend the alert.
func (tb *AlertBuilder) Suspend() string {
	return fmt.Sprintf(`ALTER ALERT %v SUSPEND`, tb.QualifiedName())
}

// Resume returns the sql that will resume the alert.
func (tb *AlertBuilder) Resume() string {
	return fmt.Sprintf(`ALTER ALERT %v RESUME`, tb.QualifiedName())
}

// Drop returns the sql that will remove the alert.
func (tb *AlertBuilder) Drop() string {
	return fmt.Sprintf(`DROP ALERT %v`, tb.QualifiedName())
}

// Describe returns the sql that will describe a alert.
func (tb *AlertBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE ALERT %v`, tb.QualifiedName())
}

// Show returns the sql that will show a alert.
func (tb *AlertBuilder) Show() string {
	return fmt.Sprintf(`SHOW ALERTS LIKE '%v' IN SCHEMA "%v"."%v"`, EscapeString(tb.name), EscapeString(tb.db), EscapeString(tb.schema))
}

// SetDisabled disables the alert builder.
func (tb *AlertBuilder) SetDisabled() *AlertBuilder {
	tb.disabled = true
	return tb
}

// IsDisabled returns if the alert builder is disabled.
func (tb *AlertBuilder) IsDisabled() bool {
	return tb.disabled
}

type Alert struct {
	CreatedOn    string  `db:"created_on"`
	Name         string  `db:"name"`
	DatabaseName string  `db:"database_name"`
	SchemaName   string  `db:"schema_name"`
	Owner        string  `db:"owner"`
	Comment      *string `db:"comment"`
	Warehouse    string  `db:"warehouse"`
	Schedule     string  `db:"schedule"`
	State        string  `db:"state"` // suspended, started
	Condition    string  `db:"condition"`
	Action       string  `db:"action"`
}

func (t *Alert) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"`, EscapeString(t.DatabaseName), EscapeString(t.SchemaName), EscapeString(t.Name))
}

func (t *Alert) Suspend() string {
	return fmt.Sprintf(`ALTER alert %v SUSPEND`, t.QualifiedName())
}

func (t *Alert) Resume() string {
	return fmt.Sprintf(`ALTER alert %v RESUME`, t.QualifiedName())
}

func (t *Alert) IsEnabled() bool {
	return strings.ToLower(t.State) == "started"
}

// ScanAlert turns a sql row into an alert object.
func ScanAlert(row *sqlx.Row) (*Alert, error) {
	t := &Alert{}
	e := row.StructScan(t)
	return t, e
}

func ListAlerts(databaseName string, schemaName string, db *sql.DB) ([]Alert, error) {
	stmt := fmt.Sprintf(`SHOW ALERTS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Alert{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no alerts found")
			return nil, nil
		}
		return dbs, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}

func WaitResumeAlert(db *sql.DB, name string, database string, schema string) error {
	builder := NewAlertBuilder(name, database, schema)

	// try to resume the alert, and verify that it was resumed.
	// if it's not resumed then try again up until a maximum of 5 times
	for i := 0; i < 5; i++ {
		q := builder.Resume()
		if err := Exec(db, q); err != nil {
			return fmt.Errorf("error resuming alert %v err = %w", name, err)
		}

		q = builder.Show()
		row := QueryRow(db, q)
		t, err := ScanAlert(row)
		if err != nil {
			return err
		}
		if t.IsEnabled() {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("unable to resume alert %v after 5 attempts", name)
}
