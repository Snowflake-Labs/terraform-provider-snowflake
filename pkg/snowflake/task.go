package snowflake

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TaskBuilder struct for building the query
type TaskBuilder struct {
	name              string
	schema            string
	database          string
	warehouse         string
	schedule          string
	scheduleSet       bool
	timeout           int
	timeoutSet        bool
	comment           string
	commentSet        bool
	predecessor       string
	predecessorSet    bool
	conditional       string
	conditionalSet    bool
	definition        string
	enabled           bool
	sessionParameters map[string]interface{}
}

// QualifiedName prepends name with the db and schema when specified
func (tb *TaskBuilder) QualifiedName() string {
	return buildFullyQualifiedTaskName(tb.name, tb.schema, tb.database)
}

// QualifiedPredecessorName prepends name with the db and schema when specified
func (tb *TaskBuilder) QualifiedPredecessorName() string {
	return buildFullyQualifiedTaskName(tb.predecessor, tb.schema, tb.database)
}

// Task returns a pointer to a TaskBuilder
func Task(name, schema, database string) *TaskBuilder {
	return &TaskBuilder{
		name:     name,
		schema:   schema,
		database: database,
	}
}

// WithWarehouse adds the warehouse property to the TaskBuilder
func (tb *TaskBuilder) WithWarehouse(warehouse string) *TaskBuilder {
	tb.warehouse = warehouse
	return tb
}

// WithSchedule sets the schedule on the TaskBuilder
func (tb *TaskBuilder) WithSchedule(schedule string) *TaskBuilder {
	tb.schedule = schedule
	tb.scheduleSet = true
	return tb
}

// WithUserTaskTimeout sets the user task timeout in ms for the TaskBuilder
func (tb *TaskBuilder) WithUserTaskTimeout(timeout int) *TaskBuilder {
	tb.timeoutSet = true
	tb.timeout = timeout
	return tb
}

// WithComment sets the comment for the TaskBuilder
func (tb *TaskBuilder) WithComment(comment string) *TaskBuilder {
	tb.comment = comment
	tb.commentSet = true
	return tb
}

// WithPredecessor sets the task whose completion will trigger the execution of this task on the TaskBuilder
func (tb *TaskBuilder) WithPredecessor(predecessor string) *TaskBuilder {
	tb.predecessor = predecessor
	tb.predecessorSet = true
	return tb
}

// WithConditional sets the conditional that determines if a task will run on the TaskBuilder
func (tb *TaskBuilder) WithConditional(conditional string) *TaskBuilder {
	tb.conditional = conditional
	tb.conditionalSet = true
	return tb
}

// WithSessionParameters adds session parameters to the TaskBuilder
func (tb *TaskBuilder) WithSessionParameters(params map[string]interface{}) *TaskBuilder {
	tb.sessionParameters = params
	return tb
}

// WithSQL Adds the task sql statement to the TaskBuilder
func (tb *TaskBuilder) WithSQL(definition string) *TaskBuilder {
	tb.definition = definition
	return tb
}

// IsEnabled Sets the status to enabled or not
func (tb *TaskBuilder) IsEnabled(enabled bool) *TaskBuilder {
	tb.enabled = enabled
	return tb
}

// Create returns the SQL statement required to create a task
func (tb *TaskBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE TASK %v `, tb.QualifiedName()))

	if tb.warehouse != "" {
		q.WriteString(fmt.Sprintf(`WAREHOUSE = "%v" `, EscapeString(tb.warehouse)))
	}

	if tb.schedule != "" {
		q.WriteString(fmt.Sprintf(`SCHEDULE = '%v' `, tb.schedule))
	}

	if len(tb.sessionParameters) > 0 {
		sp := make([]string, 0)
		sortedKeys := make([]string, 0)
		for k := range tb.sessionParameters {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			sp = append(sp, fmt.Sprintf("%v = '%v'", k, tb.sessionParameters[k]))
		}
		q.WriteString(fmt.Sprintf(`%v `, strings.Join(sp, ", ")))
	}

	if tb.timeoutSet {
		q.WriteString(fmt.Sprintf(`USER_TASK_TIMEOUT_MS = %v `, tb.timeout))
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(`COMMENT = '%v' `, EscapeString(tb.comment)))
	}

	if tb.predecessor != "" {
		q.WriteString(fmt.Sprintf(`AFTER %v `, tb.QualifiedPredecessorName()))
	}

	if tb.conditional != "" {
		q.WriteString(fmt.Sprintf(`WHEN %v `, tb.conditional))
	}

	q.WriteString(fmt.Sprintf(`AS %v`, tb.definition))

	return q.String()
}

// ChangeState Returns the query to alter the task state based on the TaskBuilder
func (tb *TaskBuilder) ChangeState() string {
	state := "SUSPEND"
	if tb.enabled {
		state = "RESUME"
	}
	return fmt.Sprintf("ALTER TASK %v %v", tb.QualifiedName(), state)
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
		p = append(p, fmt.Sprintf(`%v = '%v'`, k, params[k]))
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

// ChangeWarehouseAndSchedule method that returns a query to update the warehouse and schedule of a task
func (tb *TaskBuilder) ChangeWarehouseAndSchedule() string {
	// TODO handle full removal of SCHEDULE...
	q := strings.Builder{}

	q.WriteString(fmt.Sprintf(`ALTER TASK %v SET `, tb.QualifiedName()))

	if tb.warehouse != "" {
		q.WriteString(fmt.Sprintf(`WAREHOUSE = '%v' `, tb.warehouse))
	}

	if tb.commentSet {
		if tb.comment == "" {
			q.WriteString("COMMENT = NULL ")
		} else {
			q.WriteString(fmt.Sprintf(`COMMENT = '%v' `, tb.comment))
		}
	}

	if tb.scheduleSet { // Schedule has been set we need to at least remove it if not update
		if tb.schedule == "" {
			q.WriteString("SCHEDULE = NULL ")
		} else {
			q.WriteString(fmt.Sprintf(`SCHEDULE = '%v' `, tb.schedule))
		}
	}

	return q.String()
}

// UpdateConditional returns the query to update the tasks conditional based on the TaskBuilder Settings
func (tb *TaskBuilder) UpdateConditional() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf("ALTER TASK %v MODIFY WHEN ", tb.QualifiedName()))
	if tb.conditional == "" {
		q.WriteString(fmt.Sprintf("NULL"))
	} else {
		q.WriteString(fmt.Sprintf("%v", tb.conditional))
	}

	return q.String()
}

// UpdateSQL returns the query to update the task sql command based on the TaskBuilder pointer
func (tb *TaskBuilder) UpdateSQL() string {
	return fmt.Sprintf("ALTER TASK %v MODIFY AS %v", tb.QualifiedName(), tb.definition)
}

// RemovePredecessor returns a query to remove the current after clause
func (tb *TaskBuilder) RemovePredecessor() string {
	return fmt.Sprintf(`ALTER TASK %v REMOVE AFTER %v`, tb.QualifiedName(), tb.QualifiedPredecessorName())
}

// SetPredecessor returns the query to set a currently null after clause on the task
func (tb *TaskBuilder) SetPredecessor() string {
	return fmt.Sprintf(`ALTER TASK %v ADD AFTER %v`, tb.QualifiedName(), tb.QualifiedPredecessorName())
}

// Show returns the SQL query that will show the task
func (tb *TaskBuilder) Show() string {
	return fmt.Sprintf(`SHOW TASKS LIKE '%v' IN "%v"."%v"`, tb.name, tb.database, tb.schema)
}

// ShowParameters returns the query to show the session parameters for the task
func (tb *TaskBuilder) ShowParameters() string {
	return fmt.Sprintf(`SHOW PARAMETERS IN TASK %v`, tb.QualifiedName())
}

// Drop Returns the Query to drop/delete a task
func (tb *TaskBuilder) Drop() string {
	return fmt.Sprintf("DROP TASK %v", tb.QualifiedName())
}

// TaskRow Struct to represent tasks row in the database for easy reading
type TaskRow struct {
	CreatedOn    string         `db:"created_on"`
	TaskName     string         `db:"name"`
	TaskID       string         `db:"id"`
	DatabaseName string         `db:"database_name"`
	SchemaName   string         `db:"schema_name"`
	Owner        string         `db:"owner"`
	Comment      sql.NullString `db:"comment"`
	Warehouse    string         `db:"warehouse"`
	Schedule     sql.NullString `db:"schedule"`
	Predecessor  sql.NullString `db:"predecessors"`
	State        string         `db:"state"`
	Definition   string         `db:"definition"`
	Condition    sql.NullString `db:"condition"`
}

// TaskSessionParameterRow  Struct to represent a row of parameters
type TaskSessionParameterRow struct {
	Key          string `db:"key"`
	Value        string `db:"value"`
	DefaultValue string `db:"default"`
	Level        string `db:"level"`
	Description  string `db:"description"`
}

// ScanTaskParameters takes a database row and converts it to a task parameter pointer
func ScanTaskParameters(rows *sqlx.Rows) ([]*TaskSessionParameterRow, error) {
	t := []*TaskSessionParameterRow{}

	for rows.Next() {
		r := &TaskSessionParameterRow{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		t = append(t, r)

	}
	return t, nil
}

// ScanTask Takes a database row and converts it to a task pointer object
func ScanTask(row *sqlx.Row) (*TaskRow, error) {
	t := &TaskRow{}
	e := row.StructScan(t)
	return t, e
}

// IsEnabled returns boolean for whether or not the task is enabled
func (tr *TaskRow) IsEnabled() bool {
	return strings.ToLower(tr.State) == "started"
}

// QualifiedName returns the fully qualified name
func (tr *TaskRow) QualifiedName() string {
	return buildFullyQualifiedTaskName(tr.TaskName, tr.SchemaName, tr.DatabaseName)
}

// QualifiedPredecessorName returns the fully quallified predecessor name with database and schema
func (tr *TaskRow) QualifiedPredecessorName() string {
	return buildFullyQualifiedTaskName(tr.Predecessor.String, tr.SchemaName, tr.DatabaseName)
}

func buildFullyQualifiedTaskName(name, schema, database string) string {
	if name == "" || schema == "" || database == "" {
		return ""
	}
	n := strings.Builder{}
	// Snowflake will strip quotes if the name/schema/database is all caps so lets handle that case
	if strings.ToUpper(database) == database && database != "" {
		n.WriteString(database)
	} else {
		n.WriteString(fmt.Sprintf("\"%v\"", database))
	}

	if strings.ToUpper(schema) == schema && schema != "" {
		n.WriteString(fmt.Sprintf(".%v", schema))
	} else {
		n.WriteString(fmt.Sprintf(".\"%v\"", schema))
	}

	if strings.ToUpper(name) == name && name != "" {
		n.WriteString(fmt.Sprintf(".%v", name))
	} else {
		n.WriteString(fmt.Sprintf(".\"%v\"", name))
	}

	return n.String()
}
