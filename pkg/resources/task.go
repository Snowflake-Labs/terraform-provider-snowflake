package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO: Go through descriptions

var taskSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the task."),
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the task."),
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the task; must be unique for the database and schema in which the task is created."),
	},
	"enabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("state", func(state any) any { return state.(string) == string(sdk.TaskStateStarted) }),
		Description:      booleanStringFieldDescription("Specifies if the task should be started (enabled) after creation or should remain suspended (default)."),
	},
	"warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "The warehouse the task will use. Omit this parameter to use Snowflake-managed compute resources for runs of this task. (Conflicts with user_task_managed_initial_warehouse_size)",
		ConflictsWith:    []string{"user_task_managed_initial_warehouse_size"},
	},
	"user_task_managed_initial_warehouse_size": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
		DiffSuppressFunc: SuppressIfAny(
			NormalizeAndCompare(sdk.ToWarehouseSize),
			IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskManagedInitialWarehouseSize))),
		),
		Description:   fmt.Sprintf("Specifies the size of the compute resources to provision for the first run of the task, before a task history is available for Snowflake to determine an ideal size. Once a task has successfully completed a few runs, Snowflake ignores this parameter setting. Valid values are (case-insensitive): %s. (Conflicts with warehouse)", possibleValuesListed(sdk.ValidWarehouseSizesString)),
		ConflictsWith: []string{"warehouse"},
	},
	"schedule": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("schedule"),
		Description:      "The schedule for periodically running the task. This can be a cron or interval in minutes. (Conflict with finalize and after)",
		ConflictsWith:    []string{"finalize", "after"},
	},
	"config": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("config"),
		// TODO: it could be retrieved with system function and show/desc (which should be used?)
		// TODO: Doc request: there's no schema for JSON config format
		Description: "Specifies a string representation of key value pairs that can be accessed by all tasks in the task graph. Must be in JSON format.",
	},
	"allow_overlapping_execution": {
		Type:             schema.TypeBool,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("allow_overlapping_execution"),
		Description:      booleanStringFieldDescription("By default, Snowflake ensures that only one instance of a particular DAG is allowed to run at a time, setting the parameter value to TRUE permits DAG runs to overlap."),
	},
	"session_parameters": {
		// TODO: Description and validation
		Type:     schema.TypeList, // TODO: make it actual schema (check user)
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"a": {},
				// TODO:
			},
		},
		Optional:    true,
		Description: "Specifies session parameters to set for the session when the task runs. A task supports all session parameters.",
	},
	"user_task_timeout_ms": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateFunc:     validation.IntAtLeast(0),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskTimeoutMs))),
		Description:      "Specifies the time limit on a single run of the task before it times out (in milliseconds).",
	},
	"suspend_task_after_num_failures": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateFunc:     validation.IntAtLeast(0),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterSuspendTaskAfterNumFailures))),
		Description:      "Specifies the number of consecutive failed task runs after which the current task is suspended automatically. The default is 0 (no automatic suspension).",
	},
	"error_integration": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("error_integration")),
		Description:      "Specifies the name of the notification integration used for error notifications.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the task.",
	},
	"finalize": {
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(
			suppressIdentifierQuoting,
			IgnoreChangeToCurrentSnowflakeValueInShow("task_relations.0.finalizer"),
		),
		ConflictsWith: []string{"schedule", "after"},
	},
	"task_auto_retry_attempts": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateFunc:     validation.IntAtLeast(0),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterTaskAutoRetryAttempts))),
		Description:      "Specifies the number of automatic task graph retry attempts. If any task graphs complete in a FAILED state, Snowflake can automatically retry the task graphs from the last task in the graph that failed.",
	},
	"user_task_minimum_trigger_interval_in_seconds": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateFunc:     validation.IntAtLeast(15),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds))),
		Description:      "Defines how frequently a task can execute in seconds. If data changes occur more often than the specified minimum, changes will be grouped and processed together.",
	},
	"after": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			DiffSuppressFunc: suppressIdentifierQuoting,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		Optional: true,
		// TODO: Check
		// Cannot use IgnoreChangeToCurrentSnowflakeValueInShow because output from predecessors may be ordered
		Description:   "Specifies one or more predecessor tasks for the current task. Use this option to create a DAG of tasks or add this task to an existing DAG. A DAG is a series of tasks that starts with a scheduled root task and is linked together by dependencies.",
		ConflictsWith: []string{"schedule", "finalize"},
	},
	"when": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInShow("condition")),
		Description:      "Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported.",
	},
	"sql_statement": { // TODO: Test all possibilities of this field (procedure, procedural logic, single sql statement)
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         false,
		DiffSuppressFunc: SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInShow("definition")),
		Description:      "Any single SQL statement, or a call to a stored procedure, executed when the task runs.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW TASKS` for the given task.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTaskSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN TASK` for the given task.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTaskParametersSchema,
		},
	},
}

func Task() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateTask,
		UpdateContext: UpdateTask,
		ReadContext:   ReadTask(true),
		DeleteContext: DeleteTask,

		Schema: taskSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext, // TODO: Import
		},
	}
}

func CreateTask(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	req := sdk.NewCreateTaskRequest(id, d.Get("sql_statement").(string))

	if v, ok := d.GetOk("warehouse"); ok {
		warehouseId, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(warehouseId))
	}

	if v, ok := d.GetOk("user_task_managed_initial_warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(size))
	}

	if v, ok := d.GetOk("schedule"); ok {
		req.WithSchedule(v.(string)) // TODO: What about cron, how do we track changed (only through show)
	}

	if v, ok := d.GetOk("config"); ok {
		req.WithConfig(v.(string))
	}

	if v, ok := d.GetOk("allow_overlapping_execution"); ok {
		req.WithAllowOverlappingExecution(v.(bool))
	}

	if v, ok := d.GetOk("session_parameters"); ok {
		sessionParameters, err := sdk.GetSessionParametersFrom(v.(map[string]any))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithSessionParameters(*sessionParameters)
	}

	if v := d.Get("user_task_timeout_ms"); v != IntDefault {
		req.WithUserTaskTimeoutMs(v.(int))
	}

	if v := d.Get("suspend_task_after_num_failures"); v != IntDefault {
		req.WithSuspendTaskAfterNumFailures(v.(int))
	}

	// TODO: Decide on name (error_notification_integration ?)
	if v, ok := d.GetOk("error_integration"); ok {
		notificationIntegrationId, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithErrorNotificationIntegration(notificationIntegrationId)
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("finalize"); ok {
		rootTaskId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithFinalize(rootTaskId)
	}

	if v := d.Get("task_auto_retry_attempts"); v != IntDefault {
		req.WithTaskAutoRetryAttempts(v.(int))
	}

	if v := d.Get("user_task_minimum_trigger_interval_in_seconds"); v != IntDefault {
		req.WithUserTaskMinimumTriggerIntervalInSeconds(v.(int))
	}

	if v, ok := d.GetOk("after"); ok { // TODO: Should after take in task names or fully qualified names?
		after := expandStringList(v.([]interface{}))
		precedingTasks := make([]sdk.SchemaObjectIdentifier, 0)
		for _, dep := range after {
			precedingTaskId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, dep)
			tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, precedingTaskId, id) // TODO: What if this fails and only half of the tasks are suspended?
			defer func() {
				if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
					log.Printf("[WARN] failed to resume tasks: %s", err)
				}
			}()
			if err != nil {
				return diag.FromErr(err)
			}
			precedingTasks = append(precedingTasks, precedingTaskId)
		}
		req.WithAfter(precedingTasks)
	}

	if v, ok := d.GetOk("when"); ok {
		req.WithWhen(v.(string))
	}

	if err := client.Tasks.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	// TODO: State upgrader for "id"
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if v := d.Get("enabled").(string); v != BooleanDefault {
		enabled, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		if enabled {
			if err := waitForTaskStart(ctx, client, id); err != nil {
				log.Printf("[WARN] failed to resume task %s", name)
			}
		}
	}

	return ReadTask(false)(ctx, d, meta)
}

func UpdateTask(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	taskId := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, taskId, taskId)
	defer func() {
		if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
			log.Printf("[WARN] failed to resume tasks: %s", err)
		}
	}()
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("warehouse") {
		newWarehouse := d.Get("warehouse")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newWarehouse == "" {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithWarehouse(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithWarehouse(sdk.NewAccountObjectIdentifier(newWarehouse.(string))))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating warehouse on task %s err = %w", taskId.FullyQualifiedName(), err))
		}
	}

	if d.HasChange("user_task_managed_initial_warehouse_size") {
		newSize := d.Get("user_task_managed_initial_warehouse_size")
		warehouse := d.Get("warehouse")

		if warehouse == "" && newSize != "" {
			size, err := sdk.ToWarehouseSize(newSize.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			alterRequest := sdk.NewAlterTaskRequest(taskId).WithSet(*sdk.NewTaskSetRequest().WithUserTaskManagedInitialWarehouseSize(size))
			err = client.Tasks.Alter(ctx, alterRequest)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating user_task_managed_initial_warehouse_size on task %s", taskId.FullyQualifiedName()))
			}
		}
	}

	if d.HasChange("error_integration") {
		newErrorIntegration := d.Get("error_integration")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newErrorIntegration == "" {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithErrorIntegration(true))
		} else {
			newErrorIntegrationId, err := sdk.ParseAccountObjectIdentifier(newErrorIntegration.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithErrorNotificationIntegration(newErrorIntegrationId))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating error integration on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("after") {
		// making changes to after require suspending the current task
		// (the task will be brought up to the correct running state in the "enabled" check at the bottom of Update function).
		err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithSuspend(true))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error suspending task %s, err: %w", taskId.FullyQualifiedName(), err))
		}

		o, n := d.GetChange("after")
		oldAfter := expandStringList(o.([]interface{}))
		newAfter := expandStringList(n.([]interface{}))

		if len(newAfter) > 0 {
			// preemptively removing schedule because a task cannot have both after and schedule
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true))); err != nil {
				return diag.FromErr(fmt.Errorf("error updating schedule on task %s", taskId.FullyQualifiedName()))
			}
		}

		// Remove old dependencies that are not in new dependencies
		toRemove := make([]sdk.SchemaObjectIdentifier, 0)
		for _, dep := range oldAfter {
			if !slices.Contains(newAfter, dep) {
				toRemove = append(toRemove, sdk.NewSchemaObjectIdentifierInSchema(taskId.SchemaId(), dep))
			}
		}
		if len(toRemove) > 0 {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithRemoveAfter(toRemove)); err != nil {
				return diag.FromErr(fmt.Errorf("error removing after dependencies from task %s", taskId.FullyQualifiedName()))
			}
		}

		// Add new dependencies that are not in old dependencies
		toAdd := make([]sdk.SchemaObjectIdentifier, 0)
		for _, dep := range newAfter {
			if !slices.Contains(oldAfter, dep) {
				toAdd = append(toAdd, sdk.NewSchemaObjectIdentifierInSchema(taskId.SchemaId(), dep))
			}
		}
		if len(toAdd) > 0 {
			for _, depId := range toAdd {
				tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, depId, taskId)
				defer func() {
					if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
						log.Printf("[WARN] failed to resume tasks: %s", err)
					}
				}()
				if err != nil {
					return diag.FromErr(err)
				}
			}

			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithAddAfter(toAdd)); err != nil {
				return diag.FromErr(fmt.Errorf("error adding after dependencies from task %s", taskId.FullyQualifiedName()))
			}
		}
	}

	if d.HasChange("schedule") {
		newSchedule := d.Get("schedule")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newSchedule == "" {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithSchedule(newSchedule.(string)))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating schedule on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("user_task_timeout_ms") {
		o, n := d.GetChange("user_task_timeout_ms")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if o.(int) > 0 && n.(int) == 0 {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithUserTaskTimeoutMs(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithUserTaskTimeoutMs(n.(int)))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating user task timeout on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("suspend_task_after_num_failures") {
		o, n := d.GetChange("suspend_task_after_num_failures")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if o.(int) > 0 && n.(int) == 0 {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithSuspendTaskAfterNumFailures(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithSuspendTaskAfterNumFailures(n.(int)))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating suspend task after num failures on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("comment") {
		newComment := d.Get("comment")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if newComment == "" {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithComment(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithComment(newComment.(string)))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating comment on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("allow_overlapping_execution") {
		n := d.Get("allow_overlapping_execution")
		alterRequest := sdk.NewAlterTaskRequest(taskId)
		if n == "" {
			alterRequest.WithUnset(*sdk.NewTaskUnsetRequest().WithAllowOverlappingExecution(true))
		} else {
			alterRequest.WithSet(*sdk.NewTaskSetRequest().WithAllowOverlappingExecution(n.(bool)))
		}
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating allow overlapping execution on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("session_parameters") {
		o, n := d.GetChange("session_parameters")

		if o == nil {
			o = make(map[string]interface{})
		}
		if n == nil {
			n = make(map[string]interface{})
		}
		os := o.(map[string]any)
		ns := n.(map[string]any)

		remove := difference(os, ns)
		add := difference(ns, os)
		change := differentValue(os, ns)

		if len(remove) > 0 {
			sessionParametersUnset, err := sdk.GetSessionParametersUnsetFrom(remove)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithUnset(*sdk.NewTaskUnsetRequest().WithSessionParametersUnset(*sessionParametersUnset))); err != nil {
				return diag.FromErr(fmt.Errorf("error removing session_parameters on task %v err = %w", d.Id(), err))
			}
		}

		if len(add) > 0 {
			sessionParameters, err := sdk.GetSessionParametersFrom(add)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithSet(*sdk.NewTaskSetRequest().WithSessionParameters(*sessionParameters))); err != nil {
				return diag.FromErr(fmt.Errorf("error adding session_parameters to task %v err = %w", d.Id(), err))
			}
		}

		if len(change) > 0 {
			sessionParameters, err := sdk.GetSessionParametersFrom(change)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithSet(*sdk.NewTaskSetRequest().WithSessionParameters(*sessionParameters))); err != nil {
				return diag.FromErr(fmt.Errorf("error updating session_parameters in task %v err = %w", d.Id(), err))
			}
		}
	}

	if d.HasChange("when") {
		n := d.Get("when")
		alterRequest := sdk.NewAlterTaskRequest(taskId).WithModifyWhen(n.(string))
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating when condition on task %s", taskId.FullyQualifiedName()))
		}
	}

	if d.HasChange("sql_statement") {
		n := d.Get("sql_statement")
		alterRequest := sdk.NewAlterTaskRequest(taskId).WithModifyAs(n.(string))
		err := client.Tasks.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating sql statement on task %s", taskId.FullyQualifiedName()))
		}
	}

	enabled := d.Get("enabled").(bool)
	if enabled {
		if waitForTaskStart(ctx, client, taskId) != nil {
			log.Printf("[WARN] failed to resume task %s", taskId.FullyQualifiedName())
		}
	} else {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(taskId).WithSuspend(true)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to suspend task %s", taskId.FullyQualifiedName()))
		}
	}

	return ReadTask(false)(ctx, d, meta)
}

func ReadTask(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		task, err := client.Tasks.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query task. Marking the resource as removed.",
						Detail:   fmt.Sprintf("task name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		taskParameters, err := client.Tasks.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"", "", task.Config, task.Config, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = setStateToValuesFromConfig(d, taskSchema, []string{
				"abc",
			}); err != nil {
				return diag.FromErr(err)
			}
		}

		if errs := errors.Join(
			// TODO: handleTaskParametersRead(d, taskParameters)
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.TaskToSchema(task)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.TaskParametersToSchema(taskParameters)}),
		); errs != nil {
			return diag.FromErr(errs)
		}

		if err := d.Set("enabled", task.State == sdk.TaskStateStarted); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("warehouse", task.Warehouse); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("schedule", task.Schedule); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("comment", task.Comment); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("allow_overlapping_execution", task.AllowOverlappingExecution); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("error_integration", task.ErrorIntegration); err != nil {
			return diag.FromErr(err)
		}

		predecessors := make([]string, len(task.Predecessors))
		for i, p := range task.Predecessors {
			predecessors[i] = p.Name()
		}
		if err := d.Set("after", predecessors); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("when", task.Condition); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("sql_statement", task.Definition); err != nil {
			return diag.FromErr(err)
		}

		opts := &sdk.ShowParametersOptions{In: &sdk.ParametersIn{Task: id}}
		params, err := client.Parameters.ShowParameters(ctx, opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(params) > 0 {
			sessionParameters := make(map[string]any)
			fieldParameters := map[string]interface{}{
				"user_task_managed_initial_warehouse_size": "",
			}

			for _, param := range params {
				if param.Level != "TASK" {
					continue
				}
				switch param.Key {
				case "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE":
					fieldParameters["user_task_managed_initial_warehouse_size"] = param.Value
				case "USER_TASK_TIMEOUT_MS":
					timeout, err := strconv.ParseInt(param.Value, 10, 64)
					if err != nil {
						return diag.FromErr(err)
					}

					fieldParameters["user_task_timeout_ms"] = timeout
				case "SUSPEND_TASK_AFTER_NUM_FAILURES":
					num, err := strconv.ParseInt(param.Value, 10, 64)
					if err != nil {
						return diag.FromErr(err)
					}

					fieldParameters["suspend_task_after_num_failures"] = num
				default:
					sessionParameters[param.Key] = param.Value
				}
			}

			if err := d.Set("session_parameters", sessionParameters); err != nil {
				return diag.FromErr(err)
			}

			for key, value := range fieldParameters {
				// lintignore:R001
				err = d.Set(key, value)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		return nil
	}
}

func DeleteTask(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, id, id)
	defer func() {
		if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
			log.Printf("[WARN] failed to resume tasks: %s", err)
		}
	}()
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id).WithIfExists(true))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting task %s err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId("")
	return nil
}

func waitForTaskStart(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithResume(true))
	if err != nil {
		return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err)
	}
	return util.Retry(5, 5*time.Second, func() (error, bool) {
		task, err := client.Tasks.ShowByID(ctx, id)
		if err != nil {
			return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err), false
		}
		if task.State != sdk.TaskStateStarted {
			return nil, false
		}
		return nil, true
	})
}

// TODO: Remove functions below

// difference find keys in 'a' but not in 'b'.
func difference(a, b map[string]any) map[string]any {
	diff := make(map[string]any)
	for k := range a {
		if _, ok := b[k]; !ok {
			diff[k] = a[k]
		}
	}
	return diff
}

// differentValue find keys present both in 'a' and 'b' but having different values.
func differentValue(a, b map[string]any) map[string]any {
	diff := make(map[string]any)
	for k, va := range a {
		if vb, ok := b[k]; ok {
			if vb != va {
				diff[k] = vb
			}
		}
	}
	return diff
}
