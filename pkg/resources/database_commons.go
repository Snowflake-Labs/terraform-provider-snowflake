package resources

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	DatabaseParametersSchema              = make(map[string]*schema.Schema)
	SharedDatabaseParametersSchema        = make(map[string]*schema.Schema)
	sharedDatabaseNotApplicableParameters = []sdk.ObjectParameter{
		sdk.ObjectParameterDataRetentionTimeInDays,
		sdk.ObjectParameterMaxDataExtensionTimeInDays,
	}
	DatabaseParametersCustomDiff = ParametersCustomDiff(
		databaseParametersProvider,
		parameter[sdk.AccountParameter]{sdk.AccountParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterExternalVolume, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterCatalog, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterDefaultDDLCollation, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterLogLevel, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterTraceLevel, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeDatabase},
		parameter[sdk.AccountParameter]{sdk.AccountParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeDatabase},
	)
)

func databaseParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)
	databaseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: id,
		},
	})
	if err != nil {
		return nil, err
	}
	return databaseParameters, nil
}

func init() {
	databaseParameterFields := []struct {
		Name         sdk.ObjectParameter
		Type         schema.ValueType
		Description  string
		DiffSuppress schema.SchemaDiffSuppressFunc
		ValidateDiag schema.SchemaValidateDiagFunc
	}{
		{
			Name:        sdk.ObjectParameterDataRetentionTimeInDays,
			Type:        schema.TypeInt,
			Description: "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
			// Choosing higher range (for the standard edition or transient databases, the maximum number is 1)
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 90)),
		},
		{
			Name:        sdk.ObjectParameterDefaultDDLCollation,
			Type:        schema.TypeString,
			Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
		},
		{
			Name:         sdk.ObjectParameterCatalog,
			Type:         schema.TypeString,
			Description:  "The database parameter that specifies the default catalog to use for Iceberg tables. For more information, see [CATALOG](https://docs.snowflake.com/en/sql-reference/parameters#catalog).",
			ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		{
			Name:         sdk.ObjectParameterExternalVolume,
			Type:         schema.TypeString,
			Description:  "The database parameter that specifies the default external volume to use for Iceberg tables. For more information, see [EXTERNAL_VOLUME](https://docs.snowflake.com/en/sql-reference/parameters#external-volume).",
			ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		{
			Name:         sdk.ObjectParameterLogLevel,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllLogLevels), true),
			DiffSuppress: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue)
			},
		},
		{
			Name:         sdk.ObjectParameterTraceLevel,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", sdk.AsStringList(sdk.AllTraceLevels)),
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllTraceLevels), true),
			DiffSuppress: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue)
			},
		},
		{
			Name:         sdk.ObjectParameterMaxDataExtensionTimeInDays,
			Type:         schema.TypeInt,
			Description:  "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale. For a detailed description of this parameter, see [MAX_DATA_EXTENSION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters.html#label-max-data-extension-time-in-days).",
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 90)),
		},
		{
			Name:        sdk.ObjectParameterReplaceInvalidCharacters,
			Type:        schema.TypeBool,
			Description: "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog. For more information, see [REPLACE_INVALID_CHARACTERS](https://docs.snowflake.com/en/sql-reference/parameters#replace-invalid-characters).",
		},
		{
			Name:         sdk.ObjectParameterStorageSerializationPolicy,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("The storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake. For more information, see [STORAGE_SERIALIZATION_POLICY](https://docs.snowflake.com/en/sql-reference/parameters#storage-serialization-policy).", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllStorageSerializationPolicies), true),
			DiffSuppress: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue)
			},
		},
		{
			Name:         sdk.ObjectParameterSuspendTaskAfterNumFailures,
			Type:         schema.TypeInt,
			Description:  "How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending. For more information, see [SUSPEND_TASK_AFTER_NUM_FAILURES](https://docs.snowflake.com/en/sql-reference/parameters#suspend-task-after-num-failures).",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterTaskAutoRetryAttempts,
			Type:         schema.TypeInt,
			Description:  "Maximum automatic retries allowed for a user task. For more information, see [TASK_AUTO_RETRY_ATTEMPTS](https://docs.snowflake.com/en/sql-reference/parameters#task-auto-retry-attempts).",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterUserTaskManagedInitialWarehouseSize,
			Type:         schema.TypeString,
			Description:  "The initial size of warehouse to use for managed warehouses in the absence of history. For more information, see [USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE](https://docs.snowflake.com/en/sql-reference/parameters#user-task-managed-initial-warehouse-size).",
			ValidateDiag: sdkValidation(sdk.ToWarehouseSize),
			DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize),
		},
		{
			Name:         sdk.ObjectParameterUserTaskTimeoutMs,
			Type:         schema.TypeInt,
			Description:  "User task execution timeout in milliseconds. For more information, see [USER_TASK_TIMEOUT_MS](https://docs.snowflake.com/en/sql-reference/parameters#user-task-timeout-ms).",
			ValidateDiag: validation.ToDiagFunc(validation.IntBetween(0, 86400000)),
		},
		{
			Name:        sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds,
			Type:        schema.TypeInt,
			Description: "Minimum amount of time between Triggered Task executions in seconds.",
			// TODO(DOC-2511): ValidateDiag: Not documented
		},
		{
			Name:        sdk.ObjectParameterQuotedIdentifiersIgnoreCase,
			Type:        schema.TypeBool,
			Description: "If true, the case of quoted identifiers is ignored. For more information, see [QUOTED_IDENTIFIERS_IGNORE_CASE](https://docs.snowflake.com/en/sql-reference/parameters#quoted-identifiers-ignore-case).",
		},
		{
			Name:        sdk.ObjectParameterEnableConsoleOutput,
			Type:        schema.TypeBool,
			Description: "If true, enables stdout/stderr fast path logging for anonymous stored procedures.",
		},
	}

	for _, field := range databaseParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		DatabaseParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
		}

		if !slices.Contains(sharedDatabaseNotApplicableParameters, field.Name) {
			SharedDatabaseParametersSchema[fieldName] = &schema.Schema{
				Type:             field.Type,
				Description:      field.Description,
				ForceNew:         true,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: field.ValidateDiag,
				DiffSuppressFunc: field.DiffSuppress,
			}
		}
	}
}

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
	dataRetentionTimeInDays = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "data_retention_time_in_days")
	maxDataExtensionTimeInDays = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "max_data_extension_time_in_days")
	if externalVolumeRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "external_volume"); externalVolumeRaw != nil {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*externalVolumeRaw))
	}
	if catalogRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "catalog"); catalogRaw != nil {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*catalogRaw))
	}
	replaceInvalidCharacters = GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "replace_invalid_characters")
	defaultDDLCollation = GetConfigPropertyAsPointerAllowingZeroValue[string](d, "default_ddl_collation")
	if storageSerializationPolicyRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "storage_serialization_policy"); storageSerializationPolicyRaw != nil {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*storageSerializationPolicyRaw))
	}
	if logLevelRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "log_level"); logLevelRaw != nil {
		logLevel = sdk.Pointer(sdk.LogLevel(*logLevelRaw))
	}
	if traceLevelRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "trace_level"); traceLevelRaw != nil {
		traceLevel = sdk.Pointer(sdk.TraceLevel(*traceLevelRaw))
	}
	suspendTaskAfterNumFailures = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "suspend_task_after_num_failures")
	taskAutoRetryAttempts = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "task_auto_retry_attempts")
	if userTaskManagedInitialWarehouseSizeRaw := GetConfigPropertyAsPointerAllowingZeroValue[string](d, "user_task_managed_initial_warehouse_size"); userTaskManagedInitialWarehouseSizeRaw != nil {
		var warehouseSize sdk.WarehouseSize
		if warehouseSize, err = sdk.ToWarehouseSize(*userTaskManagedInitialWarehouseSizeRaw); err != nil {
			return
		}
		userTaskManagedInitialWarehouseSize = sdk.Pointer(warehouseSize)
	}
	userTaskTimeoutMs = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "user_task_timeout_ms")
	userTaskMinimumTriggerIntervalInSeconds = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "user_task_minimum_trigger_interval_in_seconds")
	quotedIdentifiersIgnoreCase = GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "quoted_identifiers_ignore_case")
	enableConsoleOutput = GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "enable_console_output")
	return
}

func HandleDatabaseParametersChanges(d *schema.ResourceData, set *sdk.DatabaseSet, unset *sdk.DatabaseUnset) diag.Diagnostics {
	return JoinDiags(
		handleValuePropertyChange[int](d, "data_retention_time_in_days", &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleValuePropertyChange[int](d, "max_data_extension_time_in_days", &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleValuePropertyChangeWithMapping[string](d, "external_volume", &set.ExternalVolume, &unset.ExternalVolume, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleValuePropertyChangeWithMapping[string](d, "catalog", &set.Catalog, &unset.Catalog, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleValuePropertyChange[bool](d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleValuePropertyChange[string](d, "default_ddl_collation", &set.DefaultDDLCollation, &unset.DefaultDDLCollation),
		handleValuePropertyChangeWithMapping[string](d, "storage_serialization_policy", &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleValuePropertyChangeWithMapping[string](d, "log_level", &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleValuePropertyChangeWithMapping[string](d, "trace_level", &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleValuePropertyChange[int](d, "suspend_task_after_num_failures", &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleValuePropertyChange[int](d, "task_auto_retry_attempts", &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleValuePropertyChangeWithMapping[string](d, "user_task_managed_initial_warehouse_size", &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleValuePropertyChange[int](d, "user_task_timeout_ms", &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleValuePropertyChange[int](d, "user_task_minimum_trigger_interval_in_seconds", &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleValuePropertyChange[bool](d, "quoted_identifiers_ignore_case", &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleValuePropertyChange[bool](d, "enable_console_output", &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
	)
}

// handleValuePropertyChange calls internally handleValuePropertyChangeWithMapping with identity mapping
func handleValuePropertyChange[T any](d *schema.ResourceData, key string, setField **T, unsetField **bool) diag.Diagnostics {
	return handleValuePropertyChangeWithMapping[T, T](d, key, setField, unsetField, func(value T) (T, error) { return value, nil })
}

// handleValuePropertyChangeWithMapping checks schema.ResourceData for change in key's value. If there's a change detected
// (or unknown value that basically indicates diff.SetNewComputed was called on the key), it checks if the value is set in the configuration.
// If the value is set, setField (representing setter for a value) is set to the new planned value applying mapping beforehand in cases where enum values,
// identifiers, etc. have to be set. Otherwise, unsetField is populated.
func handleValuePropertyChangeWithMapping[T, R any](d *schema.ResourceData, key string, setField **R, unsetField **bool, mapping func(value T) (R, error)) diag.Diagnostics {
	if d.HasChange(key) || !d.GetRawPlan().AsValueMap()[key].IsKnown() {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			mappedValue, err := mapping(d.Get(key).(T))
			if err != nil {
				return diag.FromErr(err)
			}
			*setField = sdk.Pointer(mappedValue)
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
