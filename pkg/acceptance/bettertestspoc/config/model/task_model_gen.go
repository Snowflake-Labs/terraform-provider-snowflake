// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type TaskModel struct {
	After                                   tfconfig.Variable `json:"after,omitempty"`
	AllowOverlappingExecution               tfconfig.Variable `json:"allow_overlapping_execution,omitempty"`
	Comment                                 tfconfig.Variable `json:"comment,omitempty"`
	Config                                  tfconfig.Variable `json:"config,omitempty"`
	Database                                tfconfig.Variable `json:"database,omitempty"`
	Enabled                                 tfconfig.Variable `json:"enabled,omitempty"`
	ErrorIntegration                        tfconfig.Variable `json:"error_integration,omitempty"`
	Finalize                                tfconfig.Variable `json:"finalize,omitempty"`
	FullyQualifiedName                      tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name                                    tfconfig.Variable `json:"name,omitempty"`
	Schedule                                tfconfig.Variable `json:"schedule,omitempty"`
	Schema                                  tfconfig.Variable `json:"schema,omitempty"`
	SessionParameters                       tfconfig.Variable `json:"session_parameters,omitempty"`
	SqlStatement                            tfconfig.Variable `json:"sql_statement,omitempty"`
	SuspendTaskAfterNumFailures             tfconfig.Variable `json:"suspend_task_after_num_failures,omitempty"`
	TaskAutoRetryAttempts                   tfconfig.Variable `json:"task_auto_retry_attempts,omitempty"`
	UserTaskManagedInitialWarehouseSize     tfconfig.Variable `json:"user_task_managed_initial_warehouse_size,omitempty"`
	UserTaskMinimumTriggerIntervalInSeconds tfconfig.Variable `json:"user_task_minimum_trigger_interval_in_seconds,omitempty"`
	UserTaskTimeoutMs                       tfconfig.Variable `json:"user_task_timeout_ms,omitempty"`
	Warehouse                               tfconfig.Variable `json:"warehouse,omitempty"`
	When                                    tfconfig.Variable `json:"when,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Task(
	resourceName string,
	database string,
	name string,
	schema string,
	sqlStatement string,
) *TaskModel {
	t := &TaskModel{ResourceModelMeta: config.Meta(resourceName, resources.Task)}
	t.WithDatabase(database)
	t.WithName(name)
	t.WithSchema(schema)
	t.WithSqlStatement(sqlStatement)
	return t
}

func TaskWithDefaultMeta(
	database string,
	name string,
	schema string,
	sqlStatement string,
) *TaskModel {
	t := &TaskModel{ResourceModelMeta: config.DefaultMeta(resources.Task)}
	t.WithDatabase(database)
	t.WithName(name)
	t.WithSchema(schema)
	t.WithSqlStatement(sqlStatement)
	return t
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// after attribute type is not yet supported, so WithAfter can't be generated

func (t *TaskModel) WithAllowOverlappingExecution(allowOverlappingExecution bool) *TaskModel {
	t.AllowOverlappingExecution = tfconfig.BoolVariable(allowOverlappingExecution)
	return t
}

func (t *TaskModel) WithComment(comment string) *TaskModel {
	t.Comment = tfconfig.StringVariable(comment)
	return t
}

func (t *TaskModel) WithConfig(config string) *TaskModel {
	t.Config = tfconfig.StringVariable(config)
	return t
}

func (t *TaskModel) WithDatabase(database string) *TaskModel {
	t.Database = tfconfig.StringVariable(database)
	return t
}

func (t *TaskModel) WithEnabled(enabled string) *TaskModel {
	t.Enabled = tfconfig.StringVariable(enabled)
	return t
}

func (t *TaskModel) WithErrorIntegration(errorIntegration string) *TaskModel {
	t.ErrorIntegration = tfconfig.StringVariable(errorIntegration)
	return t
}

// finalize attribute type is not yet supported, so WithFinalize can't be generated

func (t *TaskModel) WithFullyQualifiedName(fullyQualifiedName string) *TaskModel {
	t.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return t
}

func (t *TaskModel) WithName(name string) *TaskModel {
	t.Name = tfconfig.StringVariable(name)
	return t
}

func (t *TaskModel) WithSchedule(schedule string) *TaskModel {
	t.Schedule = tfconfig.StringVariable(schedule)
	return t
}

func (t *TaskModel) WithSchema(schema string) *TaskModel {
	t.Schema = tfconfig.StringVariable(schema)
	return t
}

// session_parameters attribute type is not yet supported, so WithSessionParameters can't be generated

func (t *TaskModel) WithSqlStatement(sqlStatement string) *TaskModel {
	t.SqlStatement = tfconfig.StringVariable(sqlStatement)
	return t
}

func (t *TaskModel) WithSuspendTaskAfterNumFailures(suspendTaskAfterNumFailures int) *TaskModel {
	t.SuspendTaskAfterNumFailures = tfconfig.IntegerVariable(suspendTaskAfterNumFailures)
	return t
}

func (t *TaskModel) WithTaskAutoRetryAttempts(taskAutoRetryAttempts int) *TaskModel {
	t.TaskAutoRetryAttempts = tfconfig.IntegerVariable(taskAutoRetryAttempts)
	return t
}

func (t *TaskModel) WithUserTaskManagedInitialWarehouseSize(userTaskManagedInitialWarehouseSize string) *TaskModel {
	t.UserTaskManagedInitialWarehouseSize = tfconfig.StringVariable(userTaskManagedInitialWarehouseSize)
	return t
}

func (t *TaskModel) WithUserTaskMinimumTriggerIntervalInSeconds(userTaskMinimumTriggerIntervalInSeconds int) *TaskModel {
	t.UserTaskMinimumTriggerIntervalInSeconds = tfconfig.IntegerVariable(userTaskMinimumTriggerIntervalInSeconds)
	return t
}

func (t *TaskModel) WithUserTaskTimeoutMs(userTaskTimeoutMs int) *TaskModel {
	t.UserTaskTimeoutMs = tfconfig.IntegerVariable(userTaskTimeoutMs)
	return t
}

func (t *TaskModel) WithWarehouse(warehouse string) *TaskModel {
	t.Warehouse = tfconfig.StringVariable(warehouse)
	return t
}

func (t *TaskModel) WithWhen(when string) *TaskModel {
	t.When = tfconfig.StringVariable(when)
	return t
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (t *TaskModel) WithAfterValue(value tfconfig.Variable) *TaskModel {
	t.After = value
	return t
}

func (t *TaskModel) WithAllowOverlappingExecutionValue(value tfconfig.Variable) *TaskModel {
	t.AllowOverlappingExecution = value
	return t
}

func (t *TaskModel) WithCommentValue(value tfconfig.Variable) *TaskModel {
	t.Comment = value
	return t
}

func (t *TaskModel) WithConfigValue(value tfconfig.Variable) *TaskModel {
	t.Config = value
	return t
}

func (t *TaskModel) WithDatabaseValue(value tfconfig.Variable) *TaskModel {
	t.Database = value
	return t
}

func (t *TaskModel) WithEnabledValue(value tfconfig.Variable) *TaskModel {
	t.Enabled = value
	return t
}

func (t *TaskModel) WithErrorIntegrationValue(value tfconfig.Variable) *TaskModel {
	t.ErrorIntegration = value
	return t
}

func (t *TaskModel) WithFinalizeValue(value tfconfig.Variable) *TaskModel {
	t.Finalize = value
	return t
}

func (t *TaskModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *TaskModel {
	t.FullyQualifiedName = value
	return t
}

func (t *TaskModel) WithNameValue(value tfconfig.Variable) *TaskModel {
	t.Name = value
	return t
}

func (t *TaskModel) WithScheduleValue(value tfconfig.Variable) *TaskModel {
	t.Schedule = value
	return t
}

func (t *TaskModel) WithSchemaValue(value tfconfig.Variable) *TaskModel {
	t.Schema = value
	return t
}

func (t *TaskModel) WithSessionParametersValue(value tfconfig.Variable) *TaskModel {
	t.SessionParameters = value
	return t
}

func (t *TaskModel) WithSqlStatementValue(value tfconfig.Variable) *TaskModel {
	t.SqlStatement = value
	return t
}

func (t *TaskModel) WithSuspendTaskAfterNumFailuresValue(value tfconfig.Variable) *TaskModel {
	t.SuspendTaskAfterNumFailures = value
	return t
}

func (t *TaskModel) WithTaskAutoRetryAttemptsValue(value tfconfig.Variable) *TaskModel {
	t.TaskAutoRetryAttempts = value
	return t
}

func (t *TaskModel) WithUserTaskManagedInitialWarehouseSizeValue(value tfconfig.Variable) *TaskModel {
	t.UserTaskManagedInitialWarehouseSize = value
	return t
}

func (t *TaskModel) WithUserTaskMinimumTriggerIntervalInSecondsValue(value tfconfig.Variable) *TaskModel {
	t.UserTaskMinimumTriggerIntervalInSeconds = value
	return t
}

func (t *TaskModel) WithUserTaskTimeoutMsValue(value tfconfig.Variable) *TaskModel {
	t.UserTaskTimeoutMs = value
	return t
}

func (t *TaskModel) WithWarehouseValue(value tfconfig.Variable) *TaskModel {
	t.Warehouse = value
	return t
}

func (t *TaskModel) WithWhenValue(value tfconfig.Variable) *TaskModel {
	t.When = value
	return t
}
