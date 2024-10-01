package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"log"
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
	//"user_task_managed_initial_warehouse_size": {
	//	Type:             schema.TypeString,
	//	Optional:         true,
	//	ValidateDiagFunc: sdkValidation(sdk.ToWarehouseSize),
	//	DiffSuppressFunc: SuppressIfAny(
	//		NormalizeAndCompare(sdk.ToWarehouseSize),
	//		IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskManagedInitialWarehouseSize))),
	//	),
	//	Description:   fmt.Sprintf("Specifies the size of the compute resources to provision for the first run of the task, before a task history is available for Snowflake to determine an ideal size. Once a task has successfully completed a few runs, Snowflake ignores this parameter setting. Valid values are (case-insensitive): %s. (Conflicts with warehouse)", possibleValuesListed(sdk.ValidWarehouseSizesString)),
	//	ConflictsWith: []string{"warehouse"},
	//},
	"schedule": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("schedule"),
		Description:      "The schedule for periodically running the task. This can be a cron or interval in minutes. (Conflict with finalize and after)",
		ConflictsWith:    []string{"finalize", "after"},
	},
	"config": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeValueInShow("config"),
			func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// TODO: Trim left and right instead of replace all + extract
				return strings.ReplaceAll(oldValue, "$", "") == strings.ReplaceAll(newValue, "$", "")
			},
		),
		// TODO: it could be retrieved with system function and show/desc (which should be used?)
		// TODO: Doc request: there's no schema for JSON config format
		Description: "Specifies a string representation of key value pairs that can be accessed by all tasks in the task graph. Must be in JSON format.",
	},
	"allow_overlapping_execution": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("allow_overlapping_execution"),
		Description:      booleanStringFieldDescription("By default, Snowflake ensures that only one instance of a particular DAG is allowed to run at a time, setting the parameter value to TRUE permits DAG runs to overlap."),
	},
	//"user_task_timeout_ms": {
	//	Type:             schema.TypeInt,
	//	Optional:         true,
	//	Default:          IntDefault,
	//	ValidateFunc:     validation.IntAtLeast(0),
	//	DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskTimeoutMs))+".0.value"),
	//	Description:      "Specifies the time limit on a single run of the task before it times out (in milliseconds).",
	//},
	//"suspend_task_after_num_failures": {
	//	Type:             schema.TypeInt,
	//	Optional:         true,
	//	Default:          IntDefault,
	//	ValidateFunc:     validation.IntAtLeast(0),
	//	DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterSuspendTaskAfterNumFailures))+".0.value"),
	//	Description:      "Specifies the number of consecutive failed task runs after which the current task is suspended automatically. The default is 0 (no automatic suspension).",
	//},
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
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(
			suppressIdentifierQuoting,
			IgnoreChangeToCurrentSnowflakeValueInShow("task_relations.0.finalize"),
		),
		ConflictsWith: []string{"schedule", "after"},
	},
	//"task_auto_retry_attempts": {
	//	Type:             schema.TypeInt,
	//	Optional:         true,
	//	Default:          IntDefault,
	//	ValidateFunc:     validation.IntAtLeast(0),
	//	DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterTaskAutoRetryAttempts))+".0.value"),
	//	Description:      "Specifies the number of automatic task graph retry attempts. If any task graphs complete in a FAILED state, Snowflake can automatically retry the task graphs from the last task in the graph that failed.",
	//},
	"user_task_minimum_trigger_interval_in_seconds": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateFunc:     validation.IntAtLeast(15),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakePlainValueInOutput(ParametersAttributeName, strings.ToLower(string(sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds))+".0.value"),
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

		Schema: helpers.MergeMaps(taskSchema, taskParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: ImportTask,
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(taskSchema, ShowOutputAttributeName, "name", "enabled", "warehouse", "user_task_managed_initial_warehouse_size", "schedule", "config", "allow_overlapping_execution", "error_integration", "comment", "finalize", "after", "when"),
			ComputedIfAnyAttributeChanged(taskParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllTaskParameters), strings.ToLower)...),
			ComputedIfAnyAttributeChanged(taskSchema, FullyQualifiedNameAttributeName, "name"),
			taskParametersCustomDiff,
		),
	}
}

func ImportTask(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	logging.DebugLogger.Printf("[DEBUG] Starting task import")
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	task, err := client.Tasks.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	if err = errors.Join(
		d.Set("enabled", booleanStringFromBool(task.State == sdk.TaskStateStarted)),
		d.Set("allow_overlapping_execution", booleanStringFromBool(task.AllowOverlappingExecution)),
	); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
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

	if v, ok := d.GetOk("schedule"); ok {
		req.WithSchedule(v.(string)) // TODO: What about cron, how do we track changed (only through show)
	}

	if v, ok := d.GetOk("config"); ok {
		req.WithConfig(v.(string))
	}

	if v := d.Get("allow_overlapping_execution").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithAllowOverlappingExecution(parsedBool)
	}

	//if v, ok := d.GetOk("session_parameters"); ok {
	//	sessionParameters, err := sdk.GetSessionParametersFrom(v.(map[string]any))
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//	req.WithSessionParameters(*sessionParameters)
	//}

	//if v := d.Get("user_task_timeout_ms"); v != IntDefault {
	//	req.WithUserTaskTimeoutMs(v.(int))
	//}
	//
	//if v := d.Get("suspend_task_after_num_failures"); v != IntDefault {
	//	req.WithSuspendTaskAfterNumFailures(v.(int))
	//}

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

	//if v := d.Get("task_auto_retry_attempts"); v != IntDefault {
	//	req.WithTaskAutoRetryAttempts(v.(int))
	//}
	//
	//if v := d.Get("user_task_minimum_trigger_interval_in_seconds"); v != IntDefault {
	//	req.WithUserTaskMinimumTriggerIntervalInSeconds(v.(int))
	//}

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

	if parameterCreateDiags := handleTaskParametersCreate(d, req); len(parameterCreateDiags) > 0 {
		return parameterCreateDiags
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

	set := sdk.NewTaskSetRequest()
	unset := sdk.NewTaskUnsetRequest()

	err = errors.Join(
		accountObjectIdentifierAttributeUpdate(d, "warehouse", &set.Warehouse, &unset.Warehouse),
		stringAttributeUpdate(d, "schedule", &set.Schedule, &unset.Schedule),
		stringAttributeUpdate(d, "config", &set.Config, &unset.Config),
		booleanStringAttributeUpdate(d, "allow_overlapping_execution", &set.AllowOverlappingExecution, &unset.AllowOverlappingExecution),
		accountObjectIdentifierAttributeUpdate(d, "error_integration", &set.ErrorNotificationIntegration, &unset.ErrorIntegration), // TODO: name inconsistency
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if updateDiags := handleTaskParametersUpdate(d, set, unset); len(updateDiags) > 0 {
		return updateDiags
	}

	if *set != (sdk.TaskSetRequest{}) {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if *unset != (sdk.TaskUnsetRequest{}) {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("when") {
		if v := d.Get("when"); v != "" {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithModifyWhen(v.(string))); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithRemoveWhen(true)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("sql_statement") {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithModifyAs(d.Get("sql_statement").(string))); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("after") {
		// TOOD: after
		// Making changes to after require suspending the current task
		// (the task will be brought up to the correct running state in the "enabled" check at the bottom of Update function).
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(true)); err != nil {
			return diag.FromErr(err)
		}

		oldAfter, newAfter := d.GetChange("after")
		addedTasks, removedTasks := ListDiff(
			expandStringList(oldAfter.(*schema.Set).List()),
			expandStringList(newAfter.(*schema.Set).List()),
		)

		// Order of commands matters:
		// The "after"s can only be added when the task doesn't have a "schedule".
		// That's why this ALTER has to be below regular ALTER SET/UNSET commands.
		if len(addedTasks) > 0 {
			addedTaskIds, err := collections.MapErr(addedTasks, sdk.ParseSchemaObjectIdentifier)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithAddAfter(addedTaskIds)); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedTasks) > 0 {
			removedTaskIds, err := collections.MapErr(removedTasks, sdk.ParseSchemaObjectIdentifier)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithRemoveAfter(removedTaskIds)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("enabled") {
		if v := d.Get("enabled").(string); v != BooleanDefault {
			enabled, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}

			if enabled {
				if waitForTaskStart(ctx, client, id) != nil {
					log.Printf("[WARN] failed to resume task %s", id.FullyQualifiedName())
				}
			} else {
				if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(true)); err != nil {
					return diag.FromErr(err)
				}
			}
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
				showMapping{"state", "enabled", string(task.State), booleanStringFromBool(task.State == sdk.TaskStateStarted), nil},
				showMapping{"allow_overlapping_execution", "allow_overlapping_execution", task.AllowOverlappingExecution, booleanStringFromBool(task.AllowOverlappingExecution), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}
		if err = setStateToValuesFromConfig(d, taskSchema, []string{
			"enabled",
			"allow_overlapping_execution",
		}); err != nil {
			return diag.FromErr(err)
		}

		errorIntegrationId := ""
		if task.ErrorIntegration != nil {
			errorIntegrationId = task.ErrorIntegration.Name()
		}

		finalizedTaskId := ""
		if task.TaskRelations.FinalizerTask != nil {
			finalizedTaskId = task.TaskRelations.FinalizerTask.FullyQualifiedName()
		}

		if errs := errors.Join(
			// TODO: handleTaskParametersRead(d, taskParameters)
			// TODO: Reorder
			d.Set("warehouse", task.Warehouse),
			d.Set("schedule", task.Schedule),
			d.Set("when", task.Condition),
			d.Set("config", task.Config),
			d.Set("error_integration", errorIntegrationId),
			d.Set("comment", task.Comment),
			d.Set("sql_statement", task.Definition),
			d.Set("after", collections.Map(task.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName)),
			d.Set("finalize", finalizedTaskId),
			handleTaskParameterRead(d, taskParameters),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.TaskToSchema(task)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.TaskParametersToSchema(taskParameters)}),
		); errs != nil {
			return diag.FromErr(errs)
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
