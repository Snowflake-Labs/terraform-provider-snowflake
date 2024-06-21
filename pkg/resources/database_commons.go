package resources

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO(SNOW-1480106): Change to smaller and safer return type
func GetAllDatabaseParameters(d *schema.ResourceData) (
	dataRetentionTimeInDays *int,
	maxDataExtensionTimeInDays *int,
	externalVolume *sdk.AccountObjectIdentifier,
	catalog *sdk.AccountObjectIdentifier,
	replaceInvalidCharacters *bool,
	defaultDDLCollation *string,
	storageSerializationPolicy *sdk.StorageSerializationPolicy,
	logLevel *sdk.LogLevel,
	traceLevel *sdk.TraceLevel,
	suspendTaskAfterNumFailures *int,
	taskAutoRetryAttempts *int,
	userTaskManagedInitialWarehouseSize *sdk.WarehouseSize,
	userTaskTimeoutMs *int,
	userTaskMinimumTriggerIntervalInSeconds *int,
	quotedIdentifiersIgnoreCase *bool,
	enableConsoleOutput *bool,
	err error,
) {
	if v := GetPropertyAsPointer[int](d, "data_retention_time_in_days"); v != nil && *v != -1 {
		dataRetentionTimeInDays = v
	}
	if v := GetPropertyAsPointer[int](d, "max_data_extension_time_in_days"); v != nil && *v != -1 {
		maxDataExtensionTimeInDays = v
	}
	if v := GetPropertyAsPointer[int](d, "suspend_task_after_num_failures"); v != nil && *v != -1 {
		suspendTaskAfterNumFailures = v
	}
	if v := GetPropertyAsPointer[int](d, "task_auto_retry_attempts"); v != nil && *v != -1 {
		taskAutoRetryAttempts = v
	}
	if v := GetPropertyAsPointer[int](d, "user_task_timeout_ms"); v != nil && *v != -1 {
		userTaskTimeoutMs = v
	}
	if v := GetPropertyAsPointer[int](d, "user_task_minimum_trigger_interval_in_seconds"); v != nil && *v != -1 {
		userTaskMinimumTriggerIntervalInSeconds = v
	}

	if v := GetPropertyAsPointer[string](d, "replace_invalid_characters"); v != nil && *v != "unknown" {
		var parsedValue bool
		parsedValue, err = strconv.ParseBool(*v)
		if err != nil {
			return
		}
		replaceInvalidCharacters = &parsedValue
	}
	if v := GetPropertyAsPointer[string](d, "quoted_identifiers_ignore_case"); v != nil && *v != "unknown" {
		var parsedValue bool
		parsedValue, err = strconv.ParseBool(*v)
		if err != nil {
			return
		}
		quotedIdentifiersIgnoreCase = &parsedValue
	}
	if v := GetPropertyAsPointer[string](d, "enable_console_output"); v != nil && *v != "unknown" {
		var parsedValue bool
		parsedValue, err = strconv.ParseBool(*v)
		if err != nil {
			return
		}
		enableConsoleOutput = &parsedValue
	}

	if v := GetPropertyAsPointer[string](d, "external_volume"); v != nil && *v != "" && *v != "unknown" {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*v))
	}
	if v := GetPropertyAsPointer[string](d, "catalog"); v != nil && *v != "" && *v != "unknown" {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*v))
	}
	if v := GetPropertyAsPointer[string](d, "default_ddl_collation"); v != nil && *v != "" && *v != "unknown" {
		defaultDDLCollation = v
	}
	if v := GetPropertyAsPointer[string](d, "storage_serialization_policy"); v != nil && *v != "" && *v != "unknown" {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*v))
	}
	if v := GetPropertyAsPointer[string](d, "log_level"); v != nil && *v != "" && *v != "unknown" {
		logLevel = sdk.Pointer(sdk.LogLevel(*v))
	}
	if v := GetPropertyAsPointer[string](d, "trace_level"); v != nil && *v != "" && *v != "unknown" {
		traceLevel = sdk.Pointer(sdk.TraceLevel(*v))
	}
	if v := GetPropertyAsPointer[string](d, "user_task_managed_initial_warehouse_size"); v != nil && *v != "" && *v != "unknown" {
		var warehouseSize sdk.WarehouseSize
		if warehouseSize, err = sdk.ToWarehouseSize(*v); err != nil {
			return
		}
		userTaskManagedInitialWarehouseSize = sdk.Pointer(warehouseSize)
	}

	return
}

func HandleDatabaseParametersChanges(d *schema.ResourceData, set *sdk.DatabaseSet, unset *sdk.DatabaseUnset) diag.Diagnostics {
	return JoinDiags(
		handleIntValueChange(d, "data_retention_time_in_days", &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleIntValueChange(d, "max_data_extension_time_in_days", &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleIntValueChange(d, "suspend_task_after_num_failures", &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleIntValueChange(d, "task_auto_retry_attempts", &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleIntValueChange(d, "user_task_timeout_ms", &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleIntValueChange(d, "user_task_minimum_trigger_interval_in_seconds", &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),

		handleBoolValueChange(d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleBoolValueChange(d, "quoted_identifiers_ignore_case", &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleBoolValueChange(d, "enable_console_output", &set.EnableConsoleOutput, &unset.EnableConsoleOutput),

		handleStringValueChange(d, "external_volume", &set.ExternalVolume, &unset.ExternalVolume, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleStringValueChange(d, "catalog", &set.Catalog, &unset.Catalog, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleStringValueChange(d, "default_ddl_collation", &set.DefaultDDLCollation, &unset.DefaultDDLCollation, func(value string) (string, error) { return value, nil }),
		handleStringValueChange(d, "storage_serialization_policy", &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, func(value string) (sdk.StorageSerializationPolicy, error) {
			return sdk.StorageSerializationPolicy(value), nil
		}),
		handleStringValueChange(d, "log_level", &set.LogLevel, &unset.LogLevel, func(value string) (sdk.LogLevel, error) { return sdk.LogLevel(value), nil }),
		handleStringValueChange(d, "trace_level", &set.TraceLevel, &unset.TraceLevel, func(value string) (sdk.TraceLevel, error) { return sdk.TraceLevel(value), nil }),
		handleStringValueChange(d, "user_task_managed_initial_warehouse_size", &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
	)
}

func handleIntValueChange(d *schema.ResourceData, key string, setField **int, unsetField **bool) diag.Diagnostics {
	if d.HasChange(key) {
		if v := d.Get(key).(int); v != -1 {
			*setField = sdk.Int(v)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func handleBoolValueChange(d *schema.ResourceData, key string, setField **bool, unsetField **bool) diag.Diagnostics {
	if d.HasChange(key) {
		if v := d.Get(key).(string); v != "unknown" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			*setField = sdk.Bool(parsed)
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func handleStringValueChange[T any](d *schema.ResourceData, key string, setField **T, unsetField **bool, mapping func(value string) (T, error)) diag.Diagnostics {
	if d.HasChange(key) {
		if v := d.Get(key); v != "unknown" {
			mappedValue, err := mapping(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			*setField = &mappedValue
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func HandleDatabaseParameterRead(d *schema.ResourceData, databaseParameters []*sdk.Parameter) diag.Diagnostics {
	for _, parameter := range databaseParameters {
		switch parameter.Key {
		case
			string(sdk.ObjectParameterDataRetentionTimeInDays),
			string(sdk.ObjectParameterMaxDataExtensionTimeInDays),
			string(sdk.ObjectParameterSuspendTaskAfterNumFailures),
			string(sdk.ObjectParameterTaskAutoRetryAttempts),
			string(sdk.ObjectParameterUserTaskTimeoutMs),
			string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds):
			value, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterExternalVolume),
			string(sdk.ObjectParameterCatalog),
			string(sdk.ObjectParameterDefaultDDLCollation),
			string(sdk.ObjectParameterStorageSerializationPolicy),
			string(sdk.ObjectParameterLogLevel),
			string(sdk.ObjectParameterTraceLevel),
			string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize):
			if err := d.Set(strings.ToLower(parameter.Key), parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterReplaceInvalidCharacters),
			string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase),
			string(sdk.ObjectParameterEnableConsoleOutput):
			value, err := strconv.ParseBool(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
