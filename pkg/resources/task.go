package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

const (
	taskIDDelimiter = '|'
)

var taskSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the task; must be unique for the database and schema in which the task is created.",
		ForceNew:    true,
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the task.",
		ForceNew:    true,
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the task.",
		ForceNew:    true,
	},
	"warehouse": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse the task will use.",
		ForceNew:    false,
	},
	"schedule": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The schedule for periodically running the task. This can be a cron or interval in minutes.",
	},
	"session_parameters": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies a comma-separated list of session parameters to set for the session when the task runs. A task supports all session parameters.",
	},
	"user_task_timeout_ms": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Specifies the time limit on a single run of the task before it times out (in milliseconds).",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the task.",
	},
	"after": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the predecessor task in the same database and schema of the current task. When a run of the predecessor task finishes successfully, it triggers this task (after a brief lag).",
	},
	"when": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported.",
	},
	"sql_statement": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Any single SQL statement, or a call to a stored procedure, executed when the task runs.",
		ForceNew:    false,
	},
}

type taskID struct {
	DatabaseName string
	SchemaName   string
	TaskName     string
}

//String() takes in a taskID object and returns a pipe-delimited string:
//DatabaseName|SchemaName|TaskName
func (t *taskID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = taskIDDelimiter
	dataIdentifiers := [][]string{{t.DatabaseName, t.SchemaName, t.TaskName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strTaskID := strings.TrimSpace(buf.String())
	return strTaskID, nil
}

// taskIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|TaskName
// and returns a taskID object
func taskIDFromString(stringID string) (*taskID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = pipeIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per task")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	taskResult := &taskID{
		DatabaseName: lines[0][0],
		SchemaName:   lines[0][1],
		TaskName:     lines[0][2],
	}
	return taskResult, nil
}

// Task returns a pointer to the resource representing a task
func Task() *schema.Resource {
	return &schema.Resource{
		Create: CreateTask,
		Read:   ReadTask,
		Update: UpdateTask,
		Delete: DeleteTask,
		Exists: TaskExists,

		Schema: taskSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// ReadTask implements schema.ReadFunc
func ReadTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(data.Id())
	if err != nil {
		return err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	q := snowflake.Task(name, database, schema).Show()
	row := snowflake.QueryRow(db, q)
	t, err := snowflake.ScanTask(row)
	if err != nil {
		return err
	}

	err = data.Set("name", t.Name)
	if err != nil {
		return err
	}

	err = data.Set("database", t.DatabaseName)
	if err != nil {
		return err
	}

	err = data.Set("schema", t.SchemaName)
	if err != nil {
		return err
	}

	err = data.Set("warehouse", t.Warehouse)
	if err != nil {
		return err
	}

	err = data.Set("schedule", t.Schedule)
	if err != nil {
		return err
	}

	err = data.Set("comment", t.Comment)
	if err != nil {
		return err
	}

	if t.Predecessors != nil {
		pre := strings.Split(*t.Predecessors, ".")
		err = data.Set("after", pre[len(pre)-1])
		if err != nil {
			return err
		}
	}

	err = data.Set("when", t.Condition)
	if err != nil {
		return err
	}

	err = data.Set("sql_statement", t.Definition)
	if err != nil {
		return err
	}

	return nil
}

// CreateTask implements schema.CreateFunc
func CreateTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	dbSchema := data.Get("schema").(string)
	name := data.Get("name").(string)
	warehouse := data.Get("warehouse").(string)
	sql := data.Get("sql_statement").(string)

	builder := snowflake.Task(name, database, dbSchema)
	builder.WithWarehouse(warehouse)
	builder.WithStatement(sql)

	// Set optionals
	if v, ok := data.GetOk("schedule"); ok {
		builder.WithSchedule(v.(string))
	}

	if v, ok := data.GetOk("session_parameters"); ok {
		builder.WithSessionParameters(expandStringList(v.(*schema.Set).List()))
	}

	if v, ok := data.GetOk("user_task_timeout_ms"); ok {
		builder.WithTimeout(v.(int))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("after"); ok {
		builder.WithDependency(v.(string))
	}

	if v, ok := data.GetOk("when"); ok {
		builder.WithCondition(v.(string))
	}

	q := builder.Create()
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating task %v", name)
	}

	q = builder.Resume()
	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error resuming task %v", name)
	}

	taskID := &taskID{
		DatabaseName: database,
		SchemaName:   dbSchema,
		TaskName:     name,
	}
	dataIDInput, err := taskID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadTask(data, meta)
}

// UpdateTask implements schema.UpdateFunc
func UpdateTask(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(data.Id())
	if err != nil {
		return err
	}

	database := taskID.DatabaseName
	dbSchema := taskID.SchemaName
	name := taskID.TaskName

	builder := snowflake.Task(name, database, dbSchema)

	if data.HasChange("warehouse") {
		_, new := data.GetChange("warehouse")
		q := builder.ChangeWarehouse(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating warehouse on task %v", data.Id())
		}
		data.SetPartial("warehouse")
	}

	if data.HasChange("schedule") {
		var q string
		old, new := data.GetChange("schedule")
		if old != "" && new == "" {
			q = builder.RemoveSchedule()
		} else {
			q = builder.ChangeSchedule(new.(string))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating schedule on task %v", data.Id())
		}
		data.SetPartial("schedule")
	}

	if data.HasChange("user_task_timeout_ms") {
		var q string
		old, new := data.GetChange("user_task_timeout_ms")
		if old.(int) > 0 && new.(int) == 0 {
			q = builder.RemoveTimeout()
		} else {
			q = builder.ChangeTimeout(new.(int))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating user task timeout on task %v", data.Id())
		}
		data.SetPartial("user_task_timeout_ms")
	}

	if data.HasChange("comment") {
		var q string
		old, new := data.GetChange("comment")
		if old != "" && new == "" {
			q = builder.RemoveComment()
		} else {
			q = builder.ChangeComment(new.(string))
		}
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating comment on task %v", data.Id())
		}
		data.SetPartial("comment")
	}

	if data.HasChange("after") {
		var (
			q   string
			err error
		)
		old, new := data.GetChange("after")

		if old != "" {
			q = builder.RemoveDependency(old.(string))
			err = snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error removing old after dependency from task %v", data.Id())
			}
		}

		if new != "" {
			q = builder.AddDependency(new.(string))
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error adding after dependency on task %v", data.Id())
			}
		}

		// Resume task after changing dependency
		q = builder.Resume()
		err = snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error resuming task %v", data.Id())
		}
		data.SetPartial("after")
	}

	if data.HasChange("session_parameters") {
		var q string
		o, n := data.GetChange("session_parameters")

		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove := expandStringList(os.Difference(ns).List())
		add := expandStringList(ns.Difference(os).List())

		if len(remove) > 0 {
			q = builder.RemoveSessionParameters(remove)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error removing session_parameters on task %v", data.Id())
			}
		}

		if len(add) > 0 {
			q = builder.AddSessionParameters(add)
			err := snowflake.Exec(db, q)
			if err != nil {
				return errors.Wrapf(err, "error adding session_parameters to task %v", data.Id())
			}
		}

		data.SetPartial("session_parameters")
	}

	if data.HasChange("when") {
		_, new := data.GetChange("when")
		q := builder.ChangeCondition(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating when condition on task %v", data.Id())
		}
		data.SetPartial("when")
	}

	if data.HasChange("sql_statement") {
		_, new := data.GetChange("sql_statement")
		q := builder.ChangeSqlStatement(new.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating sql statement on task %v", data.Id())
		}
		data.SetPartial("sql_statement")
	}

	return ReadTask(data, meta)
}

// DeleteTask implements schema.DeleteFunc
func DeleteTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(data.Id())
	if err != nil {
		return err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	q := snowflake.Task(name, database, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting task %v", data.Id())
	}

	data.SetId("")

	return nil

}

// TaskExists implements schema.ExistsFunc
func TaskExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	taskID, err := taskIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	database := taskID.DatabaseName
	schema := taskID.SchemaName
	name := taskID.TaskName

	q := snowflake.Task(name, database, schema).Show()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
