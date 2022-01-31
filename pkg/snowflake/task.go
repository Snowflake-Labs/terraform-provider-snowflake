package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TaskBuilder abstracts the creation of sql queries for a snowflake task
type TaskBuilder struct {
	name                                     string
	db                                       string
	schema                                   string
	warehouse                                string
	schedule                                 string
	session_parameters                       map[string]interface{}
	user_task_timeout_ms                     int
	comment                                  string
	after                                    string
	when                                     string
	sql_statement                            string
	disabled                                 bool
	user_task_managed_initial_warehouse_size string
	errorIntegration                         string
}

// GetFullName prepends db and schema to in parameter
func (tb *TaskBuilder) GetFullName(in string) string {
	var n strings.Builder

	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, tb.db, tb.schema, in))

	return n.String()
}

// QualifiedName prepends the db and schema and escapes everything nicely
func (tb *TaskBuilder) QualifiedName() string {
	return tb.GetFullName(tb.name)
}

// Name returns the name of the task
func (tb *TaskBuilder) Name() string {
	return tb.name
}

// WithWarehouse adds a warehouse to the TaskBuilder
func (tb *TaskBuilder) WithWarehouse(s string) *TaskBuilder {
	tb.warehouse = s
	return tb
}

// WithSchedule adds a schedule to the TaskBuilder
func (tb *TaskBuilder) WithSchedule(s string) *TaskBuilder {
	tb.schedule = s
	return tb
}

// WithSessionParameters adds session parameters to the TaskBuilder
func (tb *TaskBuilder) WithSessionParameters(params map[string]interface{}) *TaskBuilder {
	tb.session_parameters = params
	return tb
}

// WithComment adds a comment to the TaskBuilder
func (tb *TaskBuilder) WithComment(c string) *TaskBuilder {
	tb.comment = c
	return tb
}

// WithTimeout adds a timeout to the TaskBuilder
func (tb *TaskBuilder) WithTimeout(t int) *TaskBuilder {
	tb.user_task_timeout_ms = t
	return tb
}

// WithDependency adds an after task dependency to the TaskBuilder
func (tb *TaskBuilder) WithDependency(after string) *TaskBuilder {
	tb.after = after
	return tb
}

// WithCondition adds a when condition to the TaskBuilder
func (tb *TaskBuilder) WithCondition(when string) *TaskBuilder {
	tb.when = when
	return tb
}

// WithStatement adds a sql statement to the TaskBuilder
func (tb *TaskBuilder) WithStatement(sql string) *TaskBuilder {
	tb.sql_statement = sql
	return tb
}

// WithInitialWarehouseSize adds an initial warehouse size to the TaskBuilder
func (tb *TaskBuilder) WithInitialWarehouseSize(initialWarehouseSize string) *TaskBuilder {
	tb.user_task_managed_initial_warehouse_size = initialWarehouseSize
	return tb
}

/// WithErrorIntegration adds ErrorIntegration specification to the TaskBuilder
func (tb *TaskBuilder) WithErrorIntegration(s string) *TaskBuilder {
	tb.errorIntegration = s
	return tb
}

// Task returns a pointer to a Builder that abstracts the DDL operations for a task.
//
// Supported DDL operations are:
//   - CREATE TASK
//   - ALTER TASK
//   - DROP TASK
//   - DESCRIBE TASK
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/tasks-intro.html#task-ddl)
func Task(name, db, schema string) *TaskBuilder {
	return &TaskBuilder{
		name:     name,
		db:       db,
		schema:   schema,
		disabled: false, // helper for when started root or standalone task gets supspended
	}
}

// Create returns the SQL that will create a new task
func (tb *TaskBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TASK %v`, tb.QualifiedName()))

	if tb.warehouse != "" {
		q.WriteString(fmt.Sprintf(` WAREHOUSE = "%v"`, EscapeString(tb.warehouse)))
	} else {
		if tb.user_task_managed_initial_warehouse_size != "" {
			q.WriteString(fmt.Sprintf(` USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = '%v'`, EscapeString(tb.user_task_managed_initial_warehouse_size)))
		}
	}

	if tb.schedule != "" {
		q.WriteString(fmt.Sprintf(` SCHEDULE = '%v'`, EscapeString(tb.schedule)))
	}

	if len(tb.session_parameters) > 0 {
		sp := make([]string, 0)
		sortedKeys := make([]string, 0)
		for k := range tb.session_parameters {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			sp = append(sp, EscapeString(fmt.Sprintf(`%v = "%v"`, k, tb.session_parameters[k])))
		}
		q.WriteString(fmt.Sprintf(` %v`, strings.Join(sp, ", ")))
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	if tb.errorIntegration != "" {
		q.WriteString(fmt.Sprintf(` ERROR_INTEGRATION = '%v'`, EscapeString(tb.errorIntegration)))
	}

	if tb.user_task_timeout_ms > 0 {
		q.WriteString(fmt.Sprintf(` USER_TASK_TIMEOUT_MS = %v`, tb.user_task_timeout_ms))
	}

	if tb.after != "" {
		q.WriteString(fmt.Sprintf(` AFTER %v`, tb.GetFullName(tb.after)))
	}

	if tb.when != "" {
		q.WriteString(fmt.Sprintf(` WHEN %v`, tb.when))
	}

	if tb.sql_statement != "" {
		q.WriteString(fmt.Sprintf(` AS %v`, UnescapeString(tb.sql_statement)))
	}

	return q.String()
}

// ChangeWarehouse returns the sql that will change the warehouse for the task.
func (tb *TaskBuilder) ChangeWarehouse(newWh string) string {
	return fmt.Sprintf(`ALTER TASK %v SET WAREHOUSE = "%v"`, tb.QualifiedName(), EscapeString(newWh))
}

// SwitchWarehouseToManaged returns the sql that will switch to managed warehouse.
func (tb *TaskBuilder) SwitchWarehouseToManaged() string {
	return fmt.Sprintf(`ALTER TASK %v SET WAREHOUSE = null`, tb.QualifiedName())
}

// SwitchManagedWithInitialSize returns the sql that will update warehouse initial size .
func (tb *TaskBuilder) SwitchManagedWithInitialSize(initialSize string) string {
	return fmt.Sprintf(`ALTER TASK %v SET USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = '%v'`, tb.QualifiedName(), EscapeString(initialSize))
}

// ChangeSchedule returns the sql that will change the schedule for the task.
func (tb *TaskBuilder) ChangeSchedule(newSchedule string) string {
	return fmt.Sprintf(`ALTER TASK %v SET SCHEDULE = '%v'`, tb.QualifiedName(), EscapeString(newSchedule))
}

// RemoveSchedule returns the sql that will remove the schedule for the task.
func (tb *TaskBuilder) RemoveSchedule() string {
	return fmt.Sprintf(`ALTER TASK %v UNSET SCHEDULE`, tb.QualifiedName())
}

// ChangeTimeout returns the sql that will change the user task timeout for the task.
func (tb *TaskBuilder) ChangeTimeout(newTimeout int) string {
	return fmt.Sprintf(`ALTER TASK %v SET USER_TASK_TIMEOUT_MS = %v`, tb.QualifiedName(), newTimeout)
}

// RemoveTimeout returns the sql that will remove the user task timeout for the task.
func (tb *TaskBuilder) RemoveTimeout() string {
	return fmt.Sprintf(`ALTER TASK %v UNSET USER_TASK_TIMEOUT_MS`, tb.QualifiedName())
}

// ChangeComment returns the sql that will change the comment for the task.
func (tb *TaskBuilder) ChangeComment(newComment string) string {
	return fmt.Sprintf(`ALTER TASK %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(newComment))
}

// RemoveComment returns the sql that will remove the comment for the task.
func (tb *TaskBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TASK %v UNSET COMMENT`, tb.QualifiedName())
}

// AddDependency returns the sql that will add the after dependency for the task.
func (tb *TaskBuilder) AddDependency(after string) string {
	return fmt.Sprintf(`ALTER TASK %v ADD AFTER %v`, tb.QualifiedName(), tb.GetFullName(after))
}

// RemoveDependency returns the sql that will remove the after dependency for the task.
func (tb *TaskBuilder) RemoveDependency(after string) string {
	return fmt.Sprintf(`ALTER TASK %v REMOVE AFTER %v`, tb.QualifiedName(), tb.GetFullName(after))
}

// AddSessionParameters returns the sql that will remove the session parameters for the task
func (tb *TaskBuilder) AddSessionParameters(params map[string]interface{}) string {
	p := make([]string, 0)
	sortedKeys := make([]string, 0)
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		p = append(p, EscapeString(fmt.Sprintf(`%v = "%v"`, k, params[k])))
	}

	return fmt.Sprintf(`ALTER TASK %v SET %v`, tb.QualifiedName(), strings.Join(p, ", "))
}

// RemoveSessionParameters returns the sql that will remove the session parameters for the task
func (tb *TaskBuilder) RemoveSessionParameters(params map[string]interface{}) string {
	sortedKeys := make([]string, 0)
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	return fmt.Sprintf(`ALTER TASK %v UNSET %v`, tb.QualifiedName(), strings.Join(sortedKeys, ", "))
}

// ChangeCondition returns the sql that will update the when condition for the task.
func (tb *TaskBuilder) ChangeCondition(newCondition string) string {
	return fmt.Sprintf(`ALTER TASK %v MODIFY WHEN %v`, tb.QualifiedName(), newCondition)
}

// ChangeSqlStatement returns the sql that will update the sql the task executes.
func (tb *TaskBuilder) ChangeSqlStatement(newStatement string) string {
	return fmt.Sprintf(`ALTER TASK %v MODIFY AS %v`, tb.QualifiedName(), UnescapeString(newStatement))
}

// Suspend returns the sql that will suspend the task.
func (tb *TaskBuilder) Suspend() string {
	return fmt.Sprintf(`ALTER TASK %v SUSPEND`, tb.QualifiedName())
}

// Resume returns the sql that will resume the task.
func (tb *TaskBuilder) Resume() string {
	return fmt.Sprintf(`ALTER TASK %v RESUME`, tb.QualifiedName())
}

// Drop returns the sql that will remove the task.
func (tb *TaskBuilder) Drop() string {
	return fmt.Sprintf(`DROP TASK %v`, tb.QualifiedName())
}

// Describe returns the sql that will describe a task.
func (tb *TaskBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE TASK %v`, tb.QualifiedName())
}

// Show returns the sql that will show a task.
func (tb *TaskBuilder) Show() string {
	return fmt.Sprintf(`SHOW TASKS LIKE '%v' IN SCHEMA "%v"."%v"`, EscapeString(tb.name), EscapeString(tb.db), EscapeString(tb.schema))
}

// ShowParameters returns the query to show the session parameters for the task
func (tb *TaskBuilder) ShowParameters() string {
	return fmt.Sprintf(`SHOW PARAMETERS IN TASK %v`, tb.QualifiedName())
}

// SetDisabled disables the task builder
func (tb *TaskBuilder) SetDisabled() *TaskBuilder {
	tb.disabled = true
	return tb
}

// IsDisabled returns if the task builder is disabled
func (tb *TaskBuilder) IsDisabled() bool {
	return tb.disabled
}

// ChangeErrorIntegration return SQL query that will update the error_integration on the task.
func (tb *TaskBuilder) ChangeErrorIntegration(c string) string {
	return fmt.Sprintf(`ALTER TASK %v SET ERROR_INTEGRATION = %v`, tb.QualifiedName(), EscapeString(c))
}

// RemoveErrorIntegration returns the SQL query that will remove the error_integration on the task.
func (tb *TaskBuilder) RemoveErrorIntegration() string {
	return fmt.Sprintf(`ALTER TASK %v UNSET ERROR_INTEGRATION`, tb.QualifiedName())
}

type task struct {
	Id               string         `db:"id"`
	CreatedOn        string         `db:"created_on"`
	Name             string         `db:"name"`
	DatabaseName     string         `db:"database_name"`
	SchemaName       string         `db:"schema_name"`
	Owner            string         `db:"owner"`
	Comment          *string        `db:"comment"`
	Warehouse        *string        `db:"warehouse"`
	Schedule         *string        `db:"schedule"`
	Predecessors     *string        `db:"predecessors"`
	State            string         `db:"state"`
	Definition       string         `db:"definition"`
	Condition        *string        `db:"condition"`
	ErrorIntegration sql.NullString `db:"error_integration"`
}

func (t *task) IsEnabled() bool {
	return strings.ToLower(t.State) == "started"
}

func (t *task) GetPredecessorName() string {
	if t.Predecessors == nil {
		return ""
	}

	pre := strings.Split(*t.Predecessors, ".")
	name, err := strconv.Unquote(pre[len(pre)-1])
	if err != nil {
		return pre[len(pre)-1]
	}
	return name
}

// ScanTask turns a sql row into a task object
func ScanTask(row *sqlx.Row) (*task, error) {
	t := &task{}
	e := row.StructScan(t)
	return t, e
}

// taskParams struct to represent a row of parameters
type taskParams struct {
	Key          string `db:"key"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Level        string `db:"level"`
	Description  string `db:"description"`
}

// ScanTaskParameters takes a database row and converts it to a task parameter pointer
func ScanTaskParameters(rows *sqlx.Rows) ([]*taskParams, error) {
	t := []*taskParams{}

	for rows.Next() {
		r := &taskParams{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		t = append(t, r)

	}
	return t, nil
}

func ListTasks(databaseName string, schemaName string, db *sql.DB) ([]task, error) {
	stmt := fmt.Sprintf(`SHOW TASKS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []task{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tasks found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
