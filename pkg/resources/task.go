package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"golang.org/x/exp/slices"
)

var taskSchema = map[string]*schema.Schema{
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies if the task should be started (enabled) after creation or should remain suspended (default).",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the task; must be unique for the database and schema in which the task is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the task.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the task.",
		ForceNew:    true,
	},
	"warehouse": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The warehouse the task will use. Omit this parameter to use Snowflake-managed compute resources for runs of this task. (Conflicts with user_task_managed_initial_warehouse_size)",
		ForceNew:      false,
		ConflictsWith: []string{"user_task_managed_initial_warehouse_size"},
	},
	"schedule": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The schedule for periodically running the task. This can be a cron or interval in minutes. (Conflict with after)",
		ConflictsWith: []string{"after"},
	},
	"session_parameters": {
		Type:        schema.TypeMap,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies session parameters to set for the session when the task runs. A task supports all session parameters.",
	},
	"user_task_timeout_ms": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 86400000),
		Description:  "Specifies the time limit on a single run of the task before it times out (in milliseconds).",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the task.",
	},
	"after": {
		Type:          schema.TypeList,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		Description:   "Specifies one or more predecessor tasks for the current task. Use this option to create a DAG of tasks or add this task to an existing DAG. A DAG is a series of tasks that starts with a scheduled root task and is linked together by dependencies.",
		ConflictsWith: []string{"schedule"},
	},
	"when": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported.",
	},
	"sql_statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Any single SQL statement, or a call to a stored procedure, executed when the task runs.",
		ForceNew:         false,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"user_task_managed_initial_warehouse_size": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"XSMALL", "X-SMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "X-LARGE", "XXLARGE", "X2LARGE", "2X-LARGE",
		}, true),
		Description:   "Specifies the size of the compute resources to provision for the first run of the task, before a task history is available for Snowflake to determine an ideal size. Once a task has successfully completed a few runs, Snowflake ignores this parameter setting. (Conflicts with warehouse)",
		ConflictsWith: []string{"warehouse"},
	},
	"error_integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the notification integration used for error notifications.",
	},
	"allow_overlapping_execution": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "By default, Snowflake ensures that only one instance of a particular DAG is allowed to run at a time, setting the parameter value to TRUE permits DAG runs to overlap.",
	},
}

// difference find keys in 'a' but not in 'b'.
func difference(a, b map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})
	for k := range a {
		if _, ok := b[k]; !ok {
			diff[k] = a[k]
		}
	}
	return diff
}

// Task returns a pointer to the resource representing a task.
func Task() *schema.Resource {
	return &schema.Resource{
		Create: CreateTask,
		Read:   ReadTask,
		Update: UpdateTask,
		Delete: DeleteTask,

		Schema: taskSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadTask implements schema.ReadFunc.
func ReadTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	taskId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	task, err := client.Tasks.ShowByID(ctx, taskId)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] task (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("enabled", task.IsStarted()); err != nil {
		return err
	}

	if err := d.Set("name", task.Name); err != nil {
		return err
	}

	if err := d.Set("database", task.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema", task.SchemaName); err != nil {
		return err
	}

	if err := d.Set("warehouse", task.Warehouse); err != nil {
		return err
	}

	if err := d.Set("schedule", task.Schedule); err != nil {
		return err
	}

	if err := d.Set("comment", task.Comment); err != nil {
		return err
	}

	if err := d.Set("allow_overlapping_execution", task.AllowOverlappingExecution); err != nil {
		return err
	}

	if err := d.Set("error_integration", task.ErrorIntegration); err != nil {
		return err
	}

	predecessors := make([]string, len(task.Predecessors))
	for i, p := range task.Predecessors {
		predecessors[i] = p.Name()
	}
	if err := d.Set("after", predecessors); err != nil {
		return err
	}

	if err := d.Set("when", task.Condition); err != nil {
		return err
	}

	if err := d.Set("sql_statement", task.Definition); err != nil {
		return err
	}

	// TODO: do task parameters
	/*
		q = builder.ShowParameters()
		paramRows, err := snowflake.Query(db, q)
		if err != nil {
			return err
		}
		params, err := snowflake.ScanTaskParameters(paramRows)
		if err != nil {
			return err
		}

		if len(params) > 0 {
			sessionParameters := map[string]interface{}{}
			fieldParameters := map[string]interface{}{
				"user_task_managed_initial_warehouse_size": "",
			}

			for _, param := range params {
				log.Printf("[TRACE] %+v\n", param)

				if param.Level != "TASK" {
					continue
				}
				switch param.Key {
				case "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE":
					fieldParameters["user_task_managed_initial_warehouse_size"] = param.Value
				case "USER_TASK_TIMEOUT_MS":
					timeout, err := strconv.ParseInt(param.Value, 10, 64)
					if err != nil {
						return err
					}

					fieldParameters["user_task_timeout_ms"] = timeout
				default:
					sessionParameters[param.Key] = param.Value
				}
			}

			if err := d.Set("session_parameters", sessionParameters); err != nil {
				return err
			}

			for key, value := range fieldParameters {
				// lintignore:R001
				err = d.Set(key, value)
				if err != nil {
					return err
				}
			}
		}
	*/

	return nil
}

// CreateTask implements schema.CreateFunc.
func CreateTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)

	sqlStatement := d.Get("sql_statement").(string)

	taskId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	createRequest := sdk.NewCreateTaskRequest(taskId, sqlStatement)

	// Set optionals
	if v, ok := d.GetOk("warehouse"); ok {
		warehouseId := sdk.NewAccountObjectIdentifier(v.(string))
		createRequest.WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithWarehouse(&warehouseId))
	}

	if v, ok := d.GetOk("user_task_managed_initial_warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return err
		}
		createRequest.WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(&size))
	}

	if v, ok := d.GetOk("schedule"); ok {
		createRequest.WithSchedule(sdk.Pointer(v.(string)))
	}

	if _, ok := d.GetOk("session_parameters"); ok {
		//TODO
		//builder.WithSessionParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("user_task_timeout_ms"); ok {
		createRequest.WithUserTaskTimeoutMs(sdk.Pointer(v.(int)))
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.Pointer(v.(string)))
	}

	if v, ok := d.GetOk("allow_overlapping_execution"); ok {
		createRequest.WithAllowOverlappingExecution(sdk.Pointer(v.(bool)))
	}

	if v, ok := d.GetOk("error_integration"); ok {
		createRequest.WithErrorIntegration(sdk.Pointer(v.(string)))
	}

	if v, ok := d.GetOk("after"); ok {
		after := expandStringList(v.([]interface{}))
		precedingTasks := make([]sdk.SchemaObjectIdentifier, 0, len(after))
		for _, dep := range after {
			precedingTaskId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, dep)
			rootTasks, err := client.Tasks.GetRootTasks(ctx, precedingTaskId)
			if err != nil {
				return err
			}
			for _, rootTask := range rootTasks {
				// if a root task is started, then it needs to be suspended before the child tasks can be created
				if rootTask.IsStarted() {
					err := suspendTask(ctx, client, rootTask.ID())
					if err != nil {
						return err
					}

					// TODO: wait here for task suspension?
					// resume the task after modifications are complete as long as it is not a standalone task
					if !(rootTask.Name == name) {
						defer func() { _ = resumeTask(ctx, client, rootTask.ID()) }()
					}
				}
			}
			precedingTasks = append(precedingTasks, precedingTaskId)
		}
		createRequest.WithAfter(precedingTasks)
	}

	if v, ok := d.GetOk("when"); ok {
		createRequest.WithWhen(sdk.Pointer(v.(string)))
	}

	if err := client.Tasks.Create(ctx, createRequest); err != nil {
		return fmt.Errorf("error creating task %s err = %w", taskId.FullyQualifiedName(), err)
	}

	d.SetId(helpers.EncodeSnowflakeID(taskId))

	enabled := d.Get("enabled").(bool)
	if enabled {
		if err := waitForTaskStart(ctx, client, taskId); err != nil {
			log.Printf("[WARN] failed to resume task %s", name)
		}
	}

	return ReadTask(d, meta)
}

func waitForTaskStart(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	err := resumeTask(ctx, client, id)
	if err != nil {
		return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err)
	}
	return helpers.Retry(5, 5*time.Second, func() (error, bool) {
		task, err := client.Tasks.ShowByID(ctx, id)
		if err != nil {
			return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err), false
		}
		if !task.IsStarted() {
			return nil, false
		}
		return nil, true
	})
}

func suspendTask(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(sdk.Pointer(true)))
	if err != nil {
		log.Printf("[WARN] failed to suspend task %s", id.FullyQualifiedName())
	}
	return err
}

func resumeTask(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithResume(sdk.Pointer(true)))
	if err != nil {
		log.Printf("[WARN] failed to resume task %s", id.FullyQualifiedName())
	}
	return err
}

func resumeTaskOld(db *sql.DB, rootTask *snowflake.Task) {
	q := rootTask.Resume()
	if err := snowflake.Exec(db, q); err != nil {
		log.Printf("[WARN] failed to resume task %s", rootTask.Name)
	}
}

// UpdateTask implements schema.UpdateFunc.
func UpdateTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	taskId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	rootTasks, err := client.Tasks.GetRootTasks(ctx, taskId)
	if err != nil {
		return err
	}
	for _, rootTask := range rootTasks {
		// if a root task is started, then it needs to be suspended before the child tasks can be created
		if rootTask.IsStarted() {
			err := suspendTask(ctx, client, rootTask.ID())
			if err != nil {
				return err
			}

			// TODO: wait here for task suspension?
			// resume the task after modifications are complete as long as it is not a standalone task
			if !(rootTask.Name == taskId.Name()) {
				defer func() { _ = resumeTask(ctx, client, rootTask.ID()) }()
			}
		}
	}

	if d.HasChange("warehouse") {
		newWarehouse := d.Get("warehouse")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newWarehouse == "" {
			alterRequest.WithUnset(sdk.NewTaskUnsetRequest().WithWarehouse(sdk.Pointer(true)))
		} else {
			alterRequest.WithSet(sdk.NewTaskSetRequest().WithWarehouse(sdk.Pointer(sdk.NewAccountObjectIdentifier(newWarehouse.(string)))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating warehouse on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("user_task_managed_initial_warehouse_size") {
		newSize := d.Get("user_task_managed_initial_warehouse_size")
		warehouse := d.Get("warehouse")

		if warehouse == "" && newSize != "" {
			// TODO: user_task_managed_initial_warehouse_size is not on the list in the docs https://docs.snowflake.com/en/sql-reference/sql/alter-task#syntax
			alterRequest := sdk.NewAlterTaskRequest(taskId).WithSet(sdk.NewTaskSetRequest())
			err := client.Tasks.Alter(ctx, alterRequest)
			if err != nil {
				return fmt.Errorf("error updating user_task_managed_initial_warehouse_size on task %s", taskId.FullyQualifiedName())
			}
		}
	}

	// TODO: error integration is not on the list in the docs https://docs.snowflake.com/en/sql-reference/sql/alter-task#syntax
	if d.HasChange("error_integration") {
		if errorIntegration, ok := d.GetOk("error_integration"); ok {
			_ = errorIntegration
			//setRequest.
		} else {
			_ = errorIntegration
			//unsetRequest.
		}
	}

	if d.HasChange("after") {
		// making changes to after require suspending the current task
		if err := suspendTask(ctx, client, taskId); err != nil {
			return fmt.Errorf("error suspending task %s", taskId.FullyQualifiedName())
		}

		o, n := d.GetChange("after")
		oldAfter := expandStringList(o.([]interface{}))
		newAfter := expandStringList(n.([]interface{}))

		if len(newAfter) > 0 {
			// preemptively removing schedule because a task cannot have both after and schedule
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithUnset(sdk.NewTaskUnsetRequest().WithSchedule(sdk.Pointer(true)))); err != nil {
				return fmt.Errorf("error updating schedule on task %s", taskId.FullyQualifiedName())
			}
		}

		// Remove old dependencies that are not in new dependencies
		var toRemove []sdk.SchemaObjectIdentifier
		for _, dep := range oldAfter {
			if !slices.Contains(newAfter, dep) {
				toRemove = append(toRemove, sdk.NewSchemaObjectIdentifier(taskId.DatabaseName(), taskId.SchemaName(), dep))
			}
		}
		if len(toRemove) > 0 {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithRemoveAfter(toRemove)); err != nil {
				return fmt.Errorf("error removing after dependencies from task %s", taskId.FullyQualifiedName())
			}
		}

		// Add new dependencies that are not in old dependencies
		var toAdd []sdk.SchemaObjectIdentifier
		for _, dep := range newAfter {
			if !slices.Contains(oldAfter, dep) {
				toAdd = append(toAdd, sdk.NewSchemaObjectIdentifier(taskId.DatabaseName(), taskId.SchemaName(), dep))
			}
		}
		if len(toAdd) > 0 {
			// need to suspend any new root tasks from dependencies before adding them
			for _, dep := range toAdd {
				rootTasks, err := client.Tasks.GetRootTasks(ctx, dep)
				if err != nil {
					return err
				}
				for _, rootTask := range rootTasks {
					// if a root task is started, then it needs to be suspended before the child tasks can be created
					if rootTask.IsStarted() {
						err := suspendTask(ctx, client, rootTask.ID())
						if err != nil {
							return err
						}

						// TODO: wait here for task suspension?
						// resume the task after modifications are complete as long as it is not a standalone task
						if !(rootTask.Name == taskId.Name()) {
							defer func() { _ = resumeTask(ctx, client, rootTask.ID()) }()
						}
					}
				}
			}
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithAddAfter(toAdd)); err != nil {
				return fmt.Errorf("error adding after dependencies from task %s", taskId.FullyQualifiedName())
			}
		}
	}

	if d.HasChange("schedule") {
		newSchedule := d.Get("schedule")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newSchedule == "" {
			alterRequest.WithUnset(sdk.NewTaskUnsetRequest().WithSchedule(sdk.Pointer(true)))
		} else {
			alterRequest.WithSet(sdk.NewTaskSetRequest().WithSchedule(sdk.Pointer(newSchedule.(string))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating schedule on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("user_task_timeout_ms") {
		o, n := d.GetChange("user_task_timeout_ms")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if o.(int) > 0 && n.(int) == 0 {
			alterRequest.WithUnset(sdk.NewTaskUnsetRequest().WithUserTaskTimeoutMs(sdk.Pointer(true)))
		} else {
			alterRequest.WithSet(sdk.NewTaskSetRequest().WithUserTaskTimeoutMs(sdk.Pointer(n.(int))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating user task timeout on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("comment") {
		newComment := d.Get("comment")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newComment == "" {
			alterRequest.WithUnset(sdk.NewTaskUnsetRequest().WithComment(sdk.Pointer(true)))
		} else {
			alterRequest.WithSet(sdk.NewTaskSetRequest().WithComment(sdk.Pointer(newComment.(string))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating comment on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("allow_overlapping_execution") {
		n := d.Get("allow_overlapping_execution")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if n == "" {
			alterRequest.WithUnset(sdk.NewTaskUnsetRequest().WithAllowOverlappingExecution(sdk.Pointer(true)))
		} else {
			alterRequest.WithSet(sdk.NewTaskSetRequest().WithAllowOverlappingExecution(sdk.Pointer(n.(bool))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating allow overlapping execution on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("session_parameters") {
		// TODO: handle session parameters
		//var q string
		//o, n := d.GetChange("session_parameters")
		//
		//if o == nil {
		//	o = make(map[string]interface{})
		//}
		//if n == nil {
		//	n = make(map[string]interface{})
		//}
		//os := o.(map[string]interface{})
		//ns := n.(map[string]interface{})
		//
		//remove := difference(os, ns)
		//add := difference(ns, os)
		//
		//if len(remove) > 0 {
		//	q = builder.RemoveSessionParameters(remove)
		//	if err := snowflake.Exec(db, q); err != nil {
		//		return fmt.Errorf("error removing session_parameters on task %v", d.Id())
		//	}
		//}
		//
		//if len(add) > 0 {
		//	q = builder.AddSessionParameters(add)
		//	if err := snowflake.Exec(db, q); err != nil {
		//		return fmt.Errorf("error adding session_parameters to task %v", d.Id())
		//	}
		//}
	}

	if d.HasChange("when") {
		n := d.Get("when")
		alterRequest := sdk.NewAlterTaskRequest(taskId).WithModifyWhen(sdk.Pointer(n.(string)))
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating when condition on task %s", taskId.FullyQualifiedName())
		}
	}

	if d.HasChange("sql_statement") {
		n := d.Get("sql_statement")
		alterRequest := sdk.NewAlterTaskRequest(taskId).WithModifyAs(sdk.Pointer(n.(string)))
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating sql statement on task %s", taskId.FullyQualifiedName())
		}
	}

	enabled := d.Get("enabled").(bool)
	if enabled {
		if waitForTaskStart(ctx, client, taskId) != nil {
			log.Printf("[WARN] failed to resume task %s", taskId.FullyQualifiedName())
		}
	} else {
		if suspendTask(ctx, client, taskId) != nil {
			return fmt.Errorf("[WARN] failed to suspend task %s", taskId.FullyQualifiedName())
		}
	}
	return ReadTask(d, meta)
}

// DeleteTask implements schema.DeleteFunc.
func DeleteTask(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	taskId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	rootTasks, err := client.Tasks.GetRootTasks(ctx, taskId)
	if err != nil {
		return err
	}
	for _, rootTask := range rootTasks {
		// if a root task is started, then it needs to be suspended before the child tasks can be created
		if rootTask.IsStarted() {
			err := suspendTask(ctx, client, rootTask.ID())
			if err != nil {
				return err
			}

			// TODO: wait here for task suspension?
			// resume the task after modifications are complete as long as it is not a standalone task
			if !(rootTask.Name == taskId.Name()) {
				defer func() { _ = resumeTask(ctx, client, rootTask.ID()) }()
			}
		}
	}

	dropRequest := sdk.NewDropTaskRequest(taskId)
	err = client.Tasks.Drop(ctx, dropRequest)
	if err != nil {
		return fmt.Errorf("error deleting task %s err = %w", taskId.FullyQualifiedName(), err)
	}

	d.SetId("")
	return nil
}
