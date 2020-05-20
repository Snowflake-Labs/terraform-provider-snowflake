package resources

import (
	"database/sql"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"
)

var taskSchema = map[string]*schema.Schema{
	"enabled": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the task, must be unique for this schema",
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the task.",
	},
	"database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the task.",
	},
	"owner": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"warehouse": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"sql": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"schedule": &schema.Schema{
		Type:     schema.TypeString,
		Required: false,
		Optional: true,
	},
	"user_task_timeout_ms": &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 86400000),
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the comment for the task",
	},
	"after": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, data *schema.ResourceData) bool {
			t := snowflake.Task(new, data.Get("schema").(string), data.Get("database").(string))
			return old == t.QualifiedName()
		},
	},
	"when": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
}

// Task returns a pointer to the resource representing a task
func Task() *schema.Resource {
	return &schema.Resource{
		Schema: taskSchema,

		Create: CreateTask,
		Read:   ReadTask,
		Delete: DeleteTask,
		Update: UpdateTask,
	}
}

// CreateTask implements schema.CreateFunc
func CreateTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)
	name := data.Get("name").(string)

	builder := snowflake.Task(name, schema, database)

	activateRoot := false
	var rootBuilder *snowflake.TaskBuilder
	if data.Get("after").(string) != "" {
		root, err := getRootTask(data.Get("after").(string), schema, database, db)
		rootBuilder = snowflake.Task(root.TaskName, schema, database)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the root task: %v", rootBuilder.QualifiedName())
		}

		if root.IsEnabled() {
			activateRoot = true
			err := deactivateTask(rootBuilder, db)
			if err != nil {
				return errors.Wrapf(err, "Failed to deactivate root task: %v", rootBuilder.QualifiedName())
			}
		}
	}

	builder.WithSQL(data.Get("sql").(string))
	builder.WithWarehouse(data.Get("warehouse").(string))

	if v, ok := data.GetOk("schedule"); ok {
		builder.WithSchedule(v.(string))
	}

	if v, ok := data.GetOk("after"); ok {
		builder.WithPredecessor(v.(string))
	}

	if v, ok := data.GetOk("user_task_timeout_ms"); ok {
		builder.WithUserTaskTimeout(v.(int))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := data.GetOk("when"); ok {
		builder.WithConditional(v.(string))
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return err
	}

	// ensure correct state of the task
	enabled := data.Get("enabled").(bool)
	if enabled {
		builder.IsEnabled(enabled)

		q = builder.ChangeState()
		err = snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "Failed to activate task: %v", builder.QualifiedName())
		}
	}

	if activateRoot && rootBuilder != nil {
		err := activateTask(rootBuilder, db)
		if err != nil {
			return errors.Wrapf(err, "failed to reactivate task: %v", rootBuilder.QualifiedName())
		}
	}

	return ReadTask(data, meta)
}

// ReadTask implements schema.ReadFunc
func ReadTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	schema := data.Get("schema").(string)
	database := data.Get("database").(string)

	builder := snowflake.Task(name, schema, database)

	q := builder.Show()

	row := snowflake.QueryRow(db, q)
	task, err := snowflake.ScanTask(row)
	if err != nil {
		return err
	}

	data.SetId(task.TaskID)

	err = data.Set("database", task.DatabaseName)
	if err != nil {
		return err
	}

	err = data.Set("schema", task.SchemaName)
	if err != nil {
		return err
	}

	err = data.Set("name", task.TaskName)
	if err != nil {
		return err
	}

	err = data.Set("owner", task.Owner)
	if err != nil {
		return err
	}

	if task.Comment.String != "" {
		err = data.Set("comment", task.Comment.String)
		if err != nil {
			return err
		}
	}

	err = data.Set("warehouse", task.Warehouse)
	if err != nil {
		return err
	}

	if task.Schedule.String != "" {
		err = data.Set("schedule", task.Schedule.String)
		if err != nil {
			return err
		}
	}

	if task.Predecessor.String != "" {
		err = data.Set("after", task.Predecessor.String)
		if err != nil {
			return err
		}
	}

	if task.Condition.String != "" {
		err = data.Set("when", task.Condition.String)
		if err != nil {
			return err
		}
	}

	err = data.Set("enabled", strings.ToLower(task.State) == "started")
	if err != nil {
		return err
	}

	err = data.Set("sql", task.Definition)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask implements schema.DeleteFunc
func DeleteTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	schema := data.Get("schema").(string)
	database := data.Get("database").(string)

	builder := snowflake.Task(name, schema, database)

	activateRoot := false
	var rootBuilder *snowflake.TaskBuilder
	if data.Get("after").(string) != "" {
		root, err := getRootTask(data.Get("after").(string), schema, database, db)
		rootBuilder = snowflake.Task(root.TaskName, schema, database)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the root task: %v", rootBuilder.QualifiedName())
		}

		if root.IsEnabled() {
			activateRoot = true
			err := deactivateTask(rootBuilder, db)
			if err != nil {
				return errors.Wrapf(err, "Failed to deactivate root task: %v", rootBuilder.QualifiedName())
			}
		}
	}

	q := builder.Drop()
	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error delete task: %v", data.Id())
	}

	data.SetId("")

	if activateRoot && rootBuilder != nil {
		err := activateTask(rootBuilder, db)
		if err != nil {
			return errors.Wrapf(err, "failed to reactivate task: %v", rootBuilder.QualifiedName())
		}
	}

	return nil
}

// UpdateTask implements schema.UpdateFunc
func UpdateTask(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	schema := data.Get("schema").(string)
	database := data.Get("database").(string)

	builder := snowflake.Task(name, schema, database)

	var rootBuilder *snowflake.TaskBuilder

	rootNode := false
	activateRoot := false
	if data.Get("after").(string) == "" {
		rootNode = true

		if data.HasChange("enabled") {
			prevState, _ := data.GetChange("enabled")
			if prevState.(bool) { // task is currently active
				activateRoot = false
				err := deactivateTask(builder, db)
				if err != nil {
					return errors.Wrapf(err, "Failed to deactivate task: %v", builder.QualifiedName())
				}
			} else {
				activateRoot = true
			}
		} else if data.Get("enabled").(bool) {
			activateRoot = true
			err := deactivateTask(builder, db)
			if err != nil {
				return errors.Wrapf(err, "Failed to deactivate task: %v", builder.QualifiedName())
			}
		}
	} else {
		// child element of a tree need to find and suspend root task
		currentNodePredecessor := data.Get("after").(string)
		var err error
		root, err := getRootTask(currentNodePredecessor, schema, database, db)
		if err != nil {
			return err
		}

		if root.IsEnabled() {
			activateRoot = true
			rootBuilder = snowflake.Task(root.TaskName, schema, database)
			err := deactivateTask(rootBuilder, db)
			if err != nil {
				return errors.Wrapf(err, "Failed to deactivate task: %v", rootBuilder.QualifiedName())
			}
		}
	}

	if data.HasChange("after") {
		curAfter, newAfter := data.GetChange("after")

		if curAfter != nil {
			builder.WithPredecessor(curAfter.(string))
			q := builder.RemovePredecessor()
			err := snowflake.Exec(db, q)
			if err != nil {
				errors.Wrapf(err, "Failed to remove previous after: %v", builder.QualifiedName())
			}
		}

		if newAfter != nil {
			builder.WithPredecessor(newAfter.(string))
			q := builder.SetPredecessor()
			err := snowflake.Exec(db, q)
			if err != nil {
				errors.Wrapf(err, "Failed to set the after value: %v", builder.QualifiedName())
			}
		}
	}

	if data.HasChange("warehouse") || data.HasChange("schedule") {
		if data.HasChange("warehouse") {
			_, warehouse := data.GetChange("warehouse")
			builder.WithWarehouse(warehouse.(string))
		}

		if data.HasChange("schedule") {
			_, schedule := data.GetChange("schedule")
			builder.WithSchedule(schedule.(string))
		}

		q := builder.ChangeWarehouseAndSchedule()
		snowflake.Exec(db, q)
	}

	if data.HasChange("sql") {
		_, definition := data.GetChange("sql")

		builder.WithSQL(definition.(string))

		q := builder.UpdateSQL()
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "Failed to update the sql definition %v", builder.QualifiedName())
		}
	}

	if (activateRoot && rootNode) || data.Get("enabled").(bool) {
		err := activateTask(builder, db)
		if err != nil {
			return errors.Wrapf(err, "failed to activate task: %v", builder.QualifiedName())
		}
	} else if !rootNode && !data.Get("enabled").(bool) {
		err := deactivateTask(builder, db)
		if err != nil {
			return errors.Wrapf(err, "failed to deactivate task: %v", builder.QualifiedName())
		}
	}

	if activateRoot && rootBuilder != nil {
		err := activateTask(rootBuilder, db)
		if err != nil {
			return err
		}
	}

	// ensure state is correct
	return ReadTask(data, meta)
}

func activateTask(builder *snowflake.TaskBuilder, db *sql.DB) error {
	q := builder.IsEnabled(true).ChangeState()
	err := snowflake.Exec(db, q)
	if err != nil {
		return err
	}

	return nil
}

func deactivateTask(builder *snowflake.TaskBuilder, db *sql.DB) error {
	q := builder.IsEnabled(false).ChangeState()
	err := snowflake.Exec(db, q)
	if err != nil {
		return err
	}

	return nil
}

func splitQualifiedName(qualifiedName string) (name, schema, database string) {
	split := strings.Split(qualifiedName, ".")
	if len(split) != 3 {
		name = qualifiedName
		return
	}

	database = strings.Trim(split[0], "\\\"")
	schema = strings.Trim(split[1], "\\\"")
	name = strings.Trim(split[2], "\\\"")
	return
}

func getRootTask(currentPredecessor, schema, database string, db *sql.DB) (*snowflake.TaskRow, error) {
	predecessor, _, _ := splitQualifiedName(currentPredecessor)
	for {
		builder := snowflake.Task(predecessor, schema, database)
		q := builder.Show()
		row := snowflake.QueryRow(db, q)
		task, err := snowflake.ScanTask(row)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to locate the root node of: %v", currentPredecessor)
		}

		if task.Predecessor.String == "" {
			return task, nil
		}

		predecessor, _, _ = splitQualifiedName(task.Predecessor.String)
	}
}
