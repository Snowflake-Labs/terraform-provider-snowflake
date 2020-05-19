package snowflake

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TaskBuilder abstracts the creation of sql queries for a snowflake task
type TaskBuilder struct {
	name                 string
	db                   string
	schema               string
	warehouse            string
	schedule             string
	session_parameters   []string
	user_task_timeout_ms int
	comment              string
	after                string
	when                 string
	sql_statement        string
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
func (tb *TaskBuilder) WithSessionParameters(params []string) *TaskBuilder {
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

// WithDepedency adds an after task dependency to the TaskBuilder
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
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL that will create a new task
func (tb *TaskBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TASK %v`, tb.QualifiedName()))
	q.WriteString(fmt.Sprintf(` WAREHOUSE = %v`, tb.warehouse))

	if tb.schedule != "" {
		q.WriteString(fmt.Sprintf(` SCHEDULE = '%v'`, tb.schedule))
	}

	if len(tb.session_parameters) > 0 {
		q.WriteString(fmt.Sprintf(` %v`, strings.Join(tb.session_parameters, ", ")))
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, tb.comment))
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
		q.WriteString(fmt.Sprintf(` AS %v`, tb.sql_statement))
	}

	return q.String()
}

// getSessionParameters gets the actual parameter
func getSessionParameters(params []string) []string {
	out := make([]string, 0)
	for _, p := range params {
		s := strings.Split(p, "=")
		out = append(out, strings.TrimSpace(s[0]))
	}
	return out
}

// ChangeWarehouse returns the sql that will change the warehouse for the task.
func (tb *TaskBuilder) ChangeWarehouse(newWh string) string {
	return fmt.Sprintf(`ALTER TASK %v SET WAREHOUSE = %v`, tb.QualifiedName(), newWh)
}

// ChangeSchedule returns the sql that will change the schedule for the task.
func (tb *TaskBuilder) ChangeSchedule(newSchedule string) string {
	return fmt.Sprintf(`ALTER TASK %v SET SCHEDULE = '%v'`, tb.QualifiedName(), newSchedule)
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
	return fmt.Sprintf(`ALTER TASK %v SET COMMENT = '%v'`, tb.QualifiedName(), newComment)
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
func (tb *TaskBuilder) AddSessionParameters(params []string) string {
	return fmt.Sprintf(`ALTER TASK %v SET %v`, tb.QualifiedName(), strings.Join(params, ", "))
}

// RemoveSessionParameters returns the sql that will remove the session parameters for the task
func (tb *TaskBuilder) RemoveSessionParameters(params []string) string {
	p := getSessionParameters(params)
	log.Println(p)
	log.Println(len(p))
	return fmt.Sprintf(`ALTER TASK %v UNSET %v`, tb.QualifiedName(), strings.Join(p, ", "))
}

// ChangeCondition returns the sql that will update the when condition for the task.
func (tb *TaskBuilder) ChangeCondition(newCondition string) string {
	return fmt.Sprintf(`ALTER TASK %v MODIFY WHEN %v`, tb.QualifiedName(), newCondition)
}

// ChangeSqlStatement returns the sql that will update the sql the task executes.
func (tb *TaskBuilder) ChangeSqlStatement(newStatement string) string {
	return fmt.Sprintf(`ALTER TASK %v MODIFY AS %v`, tb.QualifiedName(), newStatement)
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
	return fmt.Sprintf(`SHOW TASKS LIKE '%v' IN DATABASE "%v"`, tb.name, tb.db)
}

type task struct {
	CreatedOn    string  `db:"created_on"`
	Name         string  `db:"name"`
	DatabaseName string  `db:"database_name"`
	SchemaName   string  `db:"schema_name"`
	Owner        string  `db:"owner"`
	Comment      *string `db:"comment"`
	Warehouse    string  `db:"warehouse"`
	Schedule     *string `db:"schedule"`
	Predecessors *string `db:"predecessors"`
	State        string  `db:"state"`
	Definition   string  `db:"definition"`
	Condition    *string `db:"condition"`
}

// ScanTask turns a sql row into a task object
func ScanTask(row *sqlx.Row) (*task, error) {
	t := &task{}
	e := row.StructScan(t)
	return t, e
}
