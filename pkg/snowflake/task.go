package snowflake

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TaskBuilder abstracts the creation of sql queries for a snowflake task.
type TaskBuilder struct {
	name                                string
	db                                  string
	schema                              string
	warehouse                           string
	schedule                            string
	sessionParameters                   map[string]interface{}
	userTaskTimeoutMS                   int
	comment                             string
	after                               []string
	when                                string
	SQLStatement                        string
	disabled                            bool
	userTaskManagedInitialWarehouseSize string
	errorIntegration                    string
	allowOverlappingExecution           bool
}

// GetFullName prepends db and schema to in parameter.
func (tb *TaskBuilder) GetFullName(name string) string {
	var n strings.Builder

	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, tb.db, tb.schema, name))

	return n.String()
}

// QualifiedName prepends the db and schema and escapes everything nicely.
func (tb *TaskBuilder) QualifiedName() string {
	return tb.GetFullName(tb.name)
}

// Name returns the name of the task.
func (tb *TaskBuilder) Name() string {
	return tb.name
}

// WithWarehouse adds a warehouse to the TaskBuilder.
func (tb *TaskBuilder) WithWarehouse(s string) *TaskBuilder {
	tb.warehouse = s
	return tb
}

// WithSchedule adds a schedule to the TaskBuilder.
func (tb *TaskBuilder) WithSchedule(s string) *TaskBuilder {
	tb.schedule = s
	return tb
}

// WithSessionParameters adds session parameters to the TaskBuilder.
func (tb *TaskBuilder) WithSessionParameters(params map[string]interface{}) *TaskBuilder {
	tb.sessionParameters = params
	return tb
}

// WithComment adds a comment to the TaskBuilder.
func (tb *TaskBuilder) WithComment(c string) *TaskBuilder {
	tb.comment = c
	return tb
}

// WithAllowOverlappingExecution set the ALLOW_OVERLAPPING_EXECUTION on the TaskBuilder.
func (tb *TaskBuilder) WithAllowOverlappingExecution(flag bool) *TaskBuilder {
	tb.allowOverlappingExecution = flag
	return tb
}

// WithTimeout adds a timeout to the TaskBuilder.
func (tb *TaskBuilder) WithTimeout(t int) *TaskBuilder {
	tb.userTaskTimeoutMS = t
	return tb
}

// WithAfter adds after task dependencies to the TaskBuilder.
func (tb *TaskBuilder) WithAfter(after []string) *TaskBuilder {
	tb.after = after
	return tb
}

// WithCondition adds a WHEN condition to the TaskBuilder.
func (tb *TaskBuilder) WithCondition(when string) *TaskBuilder {
	tb.when = when
	return tb
}

// WithStatement adds a sql statement to the TaskBuilder.
func (tb *TaskBuilder) WithStatement(sql string) *TaskBuilder {
	tb.SQLStatement = sql
	return tb
}

// WithInitialWarehouseSize adds an initial warehouse size to the TaskBuilder.
func (tb *TaskBuilder) WithInitialWarehouseSize(initialWarehouseSize string) *TaskBuilder {
	tb.userTaskManagedInitialWarehouseSize = initialWarehouseSize
	return tb
}

// WithErrorIntegration adds ErrorIntegration specification to the TaskBuilder.
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
func NewTaskBuilder(name, db, schema string) *TaskBuilder {
	return &TaskBuilder{
		name:     name,
		db:       db,
		schema:   schema,
		disabled: false, // helper for when started root or standalone task gets suspended
	}
}

// Create returns the SQL that will create a new task.
func (tb *TaskBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TASK %v`, tb.QualifiedName()))

	if tb.warehouse != "" {
		q.WriteString(fmt.Sprintf(` WAREHOUSE = "%v"`, EscapeString(tb.warehouse)))
	} else if tb.userTaskManagedInitialWarehouseSize != "" {
		q.WriteString(fmt.Sprintf(` USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = '%v'`, EscapeString(tb.userTaskManagedInitialWarehouseSize)))
	}

	if tb.schedule != "" {
		q.WriteString(fmt.Sprintf(` SCHEDULE = '%v'`, EscapeString(tb.schedule)))
	}

	if len(tb.sessionParameters) > 0 {
		sp := make([]string, 0)
		sortedKeys := make([]string, 0)
		for k := range tb.sessionParameters {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			sp = append(sp, EscapeString(fmt.Sprintf(`%v = "%v"`, k, tb.sessionParameters[k])))
		}
		q.WriteString(fmt.Sprintf(` %v`, strings.Join(sp, ", ")))
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	if tb.allowOverlappingExecution {
		q.WriteString(` ALLOW_OVERLAPPING_EXECUTION = TRUE`)
	}

	if tb.errorIntegration != "" {
		q.WriteString(fmt.Sprintf(` ERROR_INTEGRATION = '%v'`, EscapeString(tb.errorIntegration)))
	}

	if tb.userTaskTimeoutMS > 0 {
		q.WriteString(fmt.Sprintf(` USER_TASK_TIMEOUT_MS = %v`, tb.userTaskTimeoutMS))
	}

	if len(tb.after) > 0 {
		after := make([]string, 0)
		for _, a := range tb.after {
			after = append(after, tb.GetFullName(a))
		}
		q.WriteString(fmt.Sprintf(` AFTER %v`, strings.Join(after, ", ")))
	}

	if tb.when != "" {
		q.WriteString(fmt.Sprintf(` WHEN %v`, tb.when))
	}

	if tb.SQLStatement != "" {
		q.WriteString(fmt.Sprintf(` AS %v`, UnescapeString(tb.SQLStatement)))
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

// SetAllowOverlappingExecutionParameter returns the sql that will change the ALLOW_OVERLAPPING_EXECUTION for the task.
func (tb *TaskBuilder) SetAllowOverlappingExecutionParameter() string {
	return fmt.Sprintf(`ALTER TASK %v SET ALLOW_OVERLAPPING_EXECUTION = TRUE`, tb.QualifiedName())
}

// UnsetAllowOverlappingExecutionParameter returns the sql that will unset the ALLOW_OVERLAPPING_EXECUTION for the task.
func (tb *TaskBuilder) UnsetAllowOverlappingExecutionParameter() string {
	return fmt.Sprintf(`ALTER TASK %v UNSET ALLOW_OVERLAPPING_EXECUTION`, tb.QualifiedName())
}

// AddAfter returns the sql that will add the after dependency for the task.
func (tb *TaskBuilder) AddAfter(after []string) string {
	afterTasks := make([]string, 0)
	for _, a := range after {
		afterTasks = append(afterTasks, tb.GetFullName(a))
	}
	return fmt.Sprintf(`ALTER TASK %v ADD AFTER %v`, tb.QualifiedName(), strings.Join(afterTasks, ", "))
}

// RemoveAfter returns the sql that will remove the after dependency for the task.
func (tb *TaskBuilder) RemoveAfter(after []string) string {
	afterTasks := make([]string, 0)
	for _, a := range after {
		afterTasks = append(afterTasks, tb.GetFullName(a))
	}
	return fmt.Sprintf(`ALTER TASK %v REMOVE AFTER %v`, tb.QualifiedName(), strings.Join(afterTasks, ", "))
}

// AddSessionParameters returns the sql that will remove the session parameters for the task.
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

// RemoveSessionParameters returns the sql that will remove the session parameters for the task.
func (tb *TaskBuilder) RemoveSessionParameters(params map[string]interface{}) string {
	sortedKeys := make([]string, 0)
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	return fmt.Sprintf(`ALTER TASK %v UNSET %v`, tb.QualifiedName(), strings.Join(sortedKeys, ", "))
}

// ChangeCondition returns the sql that will update the WHEN condition for the task.
func (tb *TaskBuilder) ChangeCondition(newCondition string) string {
	return fmt.Sprintf(`ALTER TASK %v MODIFY WHEN %v`, tb.QualifiedName(), newCondition)
}

// ChangeSQLStatement returns the sql that will update the sql the task executes.
func (tb *TaskBuilder) ChangeSQLStatement(newStatement string) string {
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

// ShowParameters returns the query to show the session parameters for the task.
func (tb *TaskBuilder) ShowParameters() string {
	return fmt.Sprintf(`SHOW PARAMETERS IN TASK %v`, tb.QualifiedName())
}

// SetDisabled disables the task builder.
func (tb *TaskBuilder) SetDisabled() *TaskBuilder {
	tb.disabled = true
	return tb
}

// IsDisabled returns if the task builder is disabled.
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

func (tb *TaskBuilder) SetAllowOverlappingExecution() *TaskBuilder {
	tb.allowOverlappingExecution = true
	return tb
}

func (tb *TaskBuilder) IsAllowOverlappingExecution() bool {
	return tb.allowOverlappingExecution
}

type Task struct {
	ID                        string         `db:"id"`
	CreatedOn                 string         `db:"created_on"`
	Name                      string         `db:"name"`
	DatabaseName              string         `db:"database_name"`
	SchemaName                string         `db:"schema_name"`
	Owner                     string         `db:"owner"`
	Comment                   *string        `db:"comment"`
	Warehouse                 *string        `db:"warehouse"`
	Schedule                  *string        `db:"schedule"`
	Predecessors              *string        `db:"predecessors"`
	State                     string         `db:"state"`
	Definition                string         `db:"definition"`
	Condition                 *string        `db:"condition"`
	ErrorIntegration          sql.NullString `db:"error_integration"`
	AllowOverlappingExecution sql.NullString `db:"allow_overlapping_execution"`
}

func (t *Task) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"`, EscapeString(t.DatabaseName), EscapeString(t.SchemaName), EscapeString(t.Name))
}

func (t *Task) Suspend() string {
	return fmt.Sprintf(`ALTER TASK %v SUSPEND`, t.QualifiedName())
}

func (t *Task) Resume() string {
	return fmt.Sprintf(`ALTER TASK %v RESUME`, t.QualifiedName())
}

func (t *Task) IsEnabled() bool {
	return strings.ToLower(t.State) == "started"
}

func (t *Task) GetPredecessors() ([]string, error) {
	if t.Predecessors == nil {
		return []string{}, nil
	}

	// Since 2022_03, Snowflake returns this as a JSON array (even empty)
	var predecessorNames []string // nolint: prealloc  //todo: fixme
	if err := json.Unmarshal([]byte(*t.Predecessors), &predecessorNames); err == nil {
		for i, predecessorName := range predecessorNames {
			formattedName := predecessorName[strings.LastIndex(predecessorName, ".")+1:]
			formattedName = strings.Trim(formattedName, "\\\"")
			predecessorNames[i] = formattedName
		}
		return predecessorNames, nil
	}

	pre := strings.Split(*t.Predecessors, ".")
	for _, p := range pre {
		predecessorName, err := strconv.Unquote(p)
		if err != nil {
			return nil, err
		}
		predecessorNames = append(predecessorNames, predecessorName)
	}
	return predecessorNames, nil
}

// ScanTask turns a sql row into a task object.
func ScanTask(row *sqlx.Row) (*Task, error) {
	t := &Task{}
	e := row.StructScan(t)
	return t, e
}

// TaskParams struct to represent a row of parameters.
type TaskParams struct {
	Key          string `db:"key"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Level        string `db:"level"`
	Description  string `db:"description"`
}

// ScanTaskParameters takes a database row and converts it to a task parameter pointer.
func ScanTaskParameters(rows *sqlx.Rows) ([]*TaskParams, error) {
	t := []*TaskParams{}

	for rows.Next() {
		r := &TaskParams{}
		if err := rows.StructScan(r); err != nil {
			return nil, err
		}
		t = append(t, r)
	}
	return t, nil
}

func ListTasks(databaseName string, schemaName string, db *sql.DB) ([]Task, error) {
	stmt := fmt.Sprintf(`SHOW TASKS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Task{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no tasks found")
			return nil, nil
		}
		return dbs, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
