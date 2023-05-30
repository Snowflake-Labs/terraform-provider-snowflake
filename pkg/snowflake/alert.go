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
	name                        string
	db                          string
	schema                      string
	warehouse                   string
	alertScheduleInterval       int
	alertScheduleCronExpression string
	alertScheduleTimeZone       string
	comment                     string
	condition                   string
	action                      string
	disabled                    bool
}

// GetFullName prepends db and schema to in parameter.
func (builder *AlertBuilder) GetFullName(name string) string {
	var n strings.Builder
	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, builder.db, builder.schema, name))
	return n.String()
}

// QualifiedName prepends the db and schema and escapes everything nicely.
func (builder *AlertBuilder) QualifiedName() string {
	return builder.GetFullName(builder.name)
}

// Name returns the name of the alert.
func (builder *AlertBuilder) Name() string {
	return builder.name
}

// WithWarehouse adds a warehouse to the AlertBuilder.
func (builder *AlertBuilder) WithWarehouse(s string) *AlertBuilder {
	builder.warehouse = s
	return builder
}

// WithAlertScheduleCronExpression adds cron expression to alert schedule.
func (builder *AlertBuilder) WithAlertScheduleCronExpression(alertScheduleCronExpression string) *AlertBuilder {
	builder.alertScheduleCronExpression = alertScheduleCronExpression
	return builder
}

// WithAlertScheduleTimeZone adds a timezone to alert schedule.
func (builder *AlertBuilder) WithAlertScheduleTimeZone(alertScheduleTimeZone string) *AlertBuilder {
	builder.alertScheduleTimeZone = alertScheduleTimeZone
	return builder
}

// WithAlertScheduleInterval adds an interval to alert schedule.
func (builder *AlertBuilder) WithAlertScheduleInterval(alertScheduleInterval int) *AlertBuilder {
	builder.alertScheduleInterval = alertScheduleInterval
	return builder
}

// WithComment adds a comment to the AlertBuilder.
func (builder *AlertBuilder) WithComment(c string) *AlertBuilder {
	builder.comment = c
	return builder
}

// WithCondition adds a condition to the AlertBuilder.
func (builder *AlertBuilder) WithCondition(condition string) *AlertBuilder {
	builder.condition = condition
	return builder
}

// WithAction adds a sql statement to the AlertBuilder.
func (builder *AlertBuilder) WithAction(action string) *AlertBuilder {
	builder.action = action
	return builder
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
func (builder *AlertBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` ALERT %v`, builder.QualifiedName()))
	q.WriteString(fmt.Sprintf(` WAREHOUSE = "%v"`, EscapeString(builder.warehouse)))

	if builder.alertScheduleCronExpression != "" {
		q.WriteString(fmt.Sprintf(" SCHEDULE = 'USING CRON %v", builder.alertScheduleCronExpression))
		if builder.alertScheduleTimeZone != "" {
			q.WriteString(fmt.Sprintf(" %v", builder.alertScheduleTimeZone))
		}
		q.WriteString("'")
	}
	if builder.alertScheduleInterval > 0 {
		q.WriteString(fmt.Sprintf(" SCHEDULE = '%v MINUTE'", builder.alertScheduleInterval))
	}

	if builder.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(builder.comment)))
	}

	q.WriteString(fmt.Sprintf(` IF (EXISTS ( %v ))`, EscapeString(builder.condition)))
	q.WriteString(fmt.Sprintf(` THEN %v`, EscapeString(builder.action)))

	return q.String()
}

// ChangeWarehouse returns the sql that will change the warehouse for the alert.
func (builder *AlertBuilder) ChangeWarehouse(newWh string) string {
	return fmt.Sprintf(`ALTER alert %v SET WAREHOUSE = "%v"`, builder.QualifiedName(), EscapeString(newWh))
}

// RemoveSchedule returns the sql that will remove the schedule for the alert.
func (builder *AlertBuilder) RemoveSchedule() string {
	return fmt.Sprintf(`ALTER ALERT %v UNSET SCHEDULE`, builder.QualifiedName())
}

// ChangeAlertCronSchedule returns the sql that will change the cron schedule for the alert.
func (builder *AlertBuilder) ChangeAlertCronSchedule(alertScheduleCronExpression string, alertScheduleTimeZone string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER ALERT %v SET SCHEDULE = 'USING CRON %v`, builder.QualifiedName(), alertScheduleCronExpression))

	if alertScheduleTimeZone != "" {
		q.WriteString(fmt.Sprintf(` %v`, alertScheduleTimeZone))
	}
	q.WriteString(`'`)
	builder.alertScheduleCronExpression = alertScheduleCronExpression
	builder.alertScheduleTimeZone = alertScheduleTimeZone
	return q.String()
}

// ChangeAlertIntervalSchedule returns the sql that will change the schedule's interval for the alert.
func (builder *AlertBuilder) ChangeAlertIntervalSchedule(alertScheduleInterval int) string {
	s := fmt.Sprintf(`ALTER ALERT %v SET SCHEDULE = '%v MINUTE'`, builder.QualifiedName(), alertScheduleInterval)
	builder.alertScheduleInterval = alertScheduleInterval
	return s
}

// ChangeComment returns the sql that will change the comment for the alert.
func (builder *AlertBuilder) ChangeComment(newComment string) string {
	return fmt.Sprintf(`ALTER alert %v SET COMMENT = '%v'`, builder.QualifiedName(), EscapeString(newComment))
}

// RemoveComment returns the sql that will remove the comment for the alert.
func (builder *AlertBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER alert %v UNSET COMMENT`, builder.QualifiedName())
}

// ChangeCondition returns the sql that will update the WHEN condition for the alert.
func (builder *AlertBuilder) ChangeCondition(newCondition string) string {
	return fmt.Sprintf(`ALTER ALERT %v MODIFY CONDITION EXISTS ( %v )`, builder.QualifiedName(), newCondition)
}

// ChangeAction returns the sql that will update the sql the alert executes.
func (builder *AlertBuilder) ChangeAction(newAction string) string {
	return fmt.Sprintf(`ALTER ALERT %v MODIFY ACTION %v`, builder.QualifiedName(), UnescapeString(newAction))
}

// Suspend returns the sql that will suspend the alert.
func (builder *AlertBuilder) Suspend() string {
	return fmt.Sprintf(`ALTER ALERT %v SUSPEND`, builder.QualifiedName())
}

// Resume returns the sql that will resume the alert.
func (builder *AlertBuilder) Resume() string {
	return fmt.Sprintf(`ALTER ALERT %v RESUME`, builder.QualifiedName())
}

// Drop returns the sql that will remove the alert.
func (builder *AlertBuilder) Drop() string {
	return fmt.Sprintf(`DROP ALERT %v`, builder.QualifiedName())
}

// Describe returns the sql that will describe a alert.
func (builder *AlertBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE ALERT %v`, builder.QualifiedName())
}

// Show returns the sql that will show a alert.
func (builder *AlertBuilder) Show() string {
	return fmt.Sprintf(`SHOW ALERTS LIKE '%v' IN SCHEMA "%v"."%v"`, EscapeString(builder.name), EscapeString(builder.db), EscapeString(builder.schema))
}

// SetDisabled disables the alert builder.
func (builder *AlertBuilder) SetDisabled() *AlertBuilder {
	builder.disabled = true
	return builder
}

// IsDisabled returns if the alert builder is disabled.
func (builder *AlertBuilder) IsDisabled() bool {
	return builder.disabled
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

func (t *Alert) IsSuspended() bool {
	return strings.ToLower(t.State) == "suspended"
}

// ScanAlert turns a sql row into an alert object.
func ScanAlert(row *sqlx.Row) (*Alert, error) {
	t := &Alert{}
	e := row.StructScan(t)
	return t, e
}

func ListAlerts(databaseName, schemaName, pattern string, db *sql.DB) ([]Alert, error) {
	stmt := strings.Builder{}
	stmt.WriteString("SHOW ALERTS")
	if pattern != "" {
		stmt.WriteString(fmt.Sprintf(` LIKE '%v'`, pattern))
	}
	if schemaName != "" && databaseName == "" {
		stmt.WriteString(fmt.Sprintf(` IN SCHEMA '%v'`, schemaName))
	}
	if databaseName != "" {
		if schemaName == "" {
			stmt.WriteString(fmt.Sprintf(` IN DATABASE %v`, databaseName))
		} else {
			stmt.WriteString(fmt.Sprintf(` IN SCHEMA %v.%v`, databaseName, schemaName))
		}
	}

	rows, err := Query(db, stmt.String())
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
		return dbs, fmt.Errorf("unable to scan row for %s err = %w", stmt.String(), err)
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

func WaitSuspendAlert(db *sql.DB, name string, database string, schema string) error {
	builder := NewAlertBuilder(name, database, schema)

	// try to suspend the alert, and verify that it was suspended.
	// if it's not suspended then try again up until a maximum of 5 times
	for i := 0; i < 5; i++ {
		q := builder.Suspend()
		if err := Exec(db, q); err != nil {
			return fmt.Errorf("error suspending alert %v err = %w", name, err)
		}

		q = builder.Show()
		row := QueryRow(db, q)
		alert, err := ScanAlert(row)
		if err != nil {
			return err
		}
		if alert.IsSuspended() {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("unable to suspend alert %v after 5 attempts", name)
}
