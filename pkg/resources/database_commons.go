package resources

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	DatabaseParametersSchema              = make(map[string]*schema.Schema)
	SharedDatabaseParametersSchema        = make(map[string]*schema.Schema)
	sharedDatabaseNotApplicableParameters = []sdk.ObjectParameter{
		sdk.ObjectParameterDataRetentionTimeInDays,
		sdk.ObjectParameterMaxDataExtensionTimeInDays,
	}
	DatabaseParametersCustomDiff = func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		if d.Id() == "" {
			return nil
		}

		client := meta.(*provider.Context).Client
		params, err := client.Parameters.ShowParameters(context.Background(), &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{
				Database: helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier),
			},
		})
		if err != nil {
			return err
		}

		return customdiff.All(
			IntParameterValueComputedIf("data_retention_time_in_days", params, sdk.ParameterTypeDatabase, sdk.AccountParameterDataRetentionTimeInDays),
			IntParameterValueComputedIf("max_data_extension_time_in_days", params, sdk.ParameterTypeDatabase, sdk.AccountParameterMaxDataExtensionTimeInDays),
			StringParameterValueComputedIf("external_volume", params, sdk.ParameterTypeDatabase, sdk.AccountParameterExternalVolume),
			StringParameterValueComputedIf("catalog", params, sdk.ParameterTypeDatabase, sdk.AccountParameterCatalog),
			BoolParameterValueComputedIf("replace_invalid_characters", params, sdk.ParameterTypeDatabase, sdk.AccountParameterReplaceInvalidCharacters),
			StringParameterValueComputedIf("default_ddl_collation", params, sdk.ParameterTypeDatabase, sdk.AccountParameterDefaultDDLCollation),
			StringParameterValueComputedIf("storage_serialization_policy", params, sdk.ParameterTypeDatabase, sdk.AccountParameterStorageSerializationPolicy),
			StringParameterValueComputedIf("log_level", params, sdk.ParameterTypeDatabase, sdk.AccountParameterLogLevel),
			StringParameterValueComputedIf("trace_level", params, sdk.ParameterTypeDatabase, sdk.AccountParameterTraceLevel),
			IntParameterValueComputedIf("suspend_task_after_num_failures", params, sdk.ParameterTypeDatabase, sdk.AccountParameterSuspendTaskAfterNumFailures),
			IntParameterValueComputedIf("task_auto_retry_attempts", params, sdk.ParameterTypeDatabase, sdk.AccountParameterTaskAutoRetryAttempts),
			StringParameterValueComputedIf("user_task_managed_initial_warehouse_size", params, sdk.ParameterTypeDatabase, sdk.AccountParameterUserTaskManagedInitialWarehouseSize),
			IntParameterValueComputedIf("user_task_timeout_ms", params, sdk.ParameterTypeDatabase, sdk.AccountParameterUserTaskTimeoutMs),
			IntParameterValueComputedIf("user_task_minimum_trigger_interval_in_seconds", params, sdk.ParameterTypeDatabase, sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds),
			BoolParameterValueComputedIf("quoted_identifiers_ignore_case", params, sdk.ParameterTypeDatabase, sdk.AccountParameterQuotedIdentifiersIgnoreCase),
			BoolParameterValueComputedIf("enable_console_output", params, sdk.ParameterTypeDatabase, sdk.AccountParameterEnableConsoleOutput),
		)(ctx, d, meta)
	}
)

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
			Description:  "The database parameter that specifies the default catalog to use for Iceberg tables.",
			ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		{
			Name:         sdk.ObjectParameterExternalVolume,
			Type:         schema.TypeString,
			Description:  "The database parameter that specifies the default external volume to use for Iceberg tables.",
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
			Description: "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog.",
		},
		{
			Name:         sdk.ObjectParameterStorageSerializationPolicy,
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("The storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllStorageSerializationPolicies), true),
			DiffSuppress: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue)
			},
		},
		{
			Name:         sdk.ObjectParameterSuspendTaskAfterNumFailures,
			Type:         schema.TypeInt,
			Description:  "How many times a task must fail in a row before it is automatically suspended. 0 disables auto-suspending.",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterTaskAutoRetryAttempts,
			Type:         schema.TypeInt,
			Description:  "Maximum automatic retries allowed for a user task.",
			ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		{
			Name:         sdk.ObjectParameterUserTaskManagedInitialWarehouseSize,
			Type:         schema.TypeString,
			Description:  "The initial size of warehouse to use for managed warehouses in the absence of history.",
			ValidateDiag: sdkValidation(sdk.ToWarehouseSize),
			DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize),
		},
		{
			Name:         sdk.ObjectParameterUserTaskTimeoutMs,
			Type:         schema.TypeInt,
			Description:  "User task execution timeout in milliseconds.",
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
			Description: "If true, the case of quoted identifiers is ignored.",
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
	dataRetentionTimeInDays = GetPropertyAsPointer[int](d, "data_retention_time_in_days")
	maxDataExtensionTimeInDays = GetPropertyAsPointer[int](d, "max_data_extension_time_in_days")
	if externalVolumeRaw := GetPropertyAsPointer[string](d, "external_volume"); externalVolumeRaw != nil {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*externalVolumeRaw))
	}
	if catalogRaw := GetPropertyAsPointer[string](d, "catalog"); catalogRaw != nil {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*catalogRaw))
	}
	replaceInvalidCharacters = GetPropertyAsPointer[bool](d, "replace_invalid_characters")
	defaultDDLCollation = GetPropertyAsPointer[string](d, "default_ddl_collation")
	if storageSerializationPolicyRaw := GetPropertyAsPointer[string](d, "storage_serialization_policy"); storageSerializationPolicyRaw != nil {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*storageSerializationPolicyRaw))
	}
	if logLevelRaw := GetPropertyAsPointer[string](d, "log_level"); logLevelRaw != nil {
		logLevel = sdk.Pointer(sdk.LogLevel(*logLevelRaw))
	}
	if traceLevelRaw := GetPropertyAsPointer[string](d, "trace_level"); traceLevelRaw != nil {
		traceLevel = sdk.Pointer(sdk.TraceLevel(*traceLevelRaw))
	}
	suspendTaskAfterNumFailures = GetPropertyAsPointer[int](d, "suspend_task_after_num_failures")
	taskAutoRetryAttempts = GetPropertyAsPointer[int](d, "task_auto_retry_attempts")
	if userTaskManagedInitialWarehouseSizeRaw := GetPropertyAsPointer[string](d, "user_task_managed_initial_warehouse_size"); userTaskManagedInitialWarehouseSizeRaw != nil {
		var warehouseSize sdk.WarehouseSize
		if warehouseSize, err = sdk.ToWarehouseSize(*userTaskManagedInitialWarehouseSizeRaw); err != nil {
			return
		}
		userTaskManagedInitialWarehouseSize = sdk.Pointer(warehouseSize)
	}
	userTaskTimeoutMs = GetPropertyAsPointer[int](d, "user_task_timeout_ms")
	userTaskMinimumTriggerIntervalInSeconds = GetPropertyAsPointer[int](d, "user_task_minimum_trigger_interval_in_seconds")
	quotedIdentifiersIgnoreCase = GetPropertyAsPointer[bool](d, "quoted_identifiers_ignore_case")
	enableConsoleOutput = GetPropertyAsPointer[bool](d, "enable_console_output")
	return
}

func HandleDatabaseParametersChanges(d *schema.ResourceData, set *sdk.DatabaseSet, unset *sdk.DatabaseUnset) diag.Diagnostics {
	return JoinDiags(
		handleValuePropertyChange[int](d, "data_retention_time_in_days", &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleValuePropertyChange[int](d, "max_data_extension_time_in_days", &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleValuePropertyChangeWithMapping[string](d, "external_volume", &set.ExternalVolume, &unset.ExternalVolume, sdk.NewAccountObjectIdentifier),
		handleValuePropertyChangeWithMapping[string](d, "catalog", &set.Catalog, &unset.Catalog, sdk.NewAccountObjectIdentifier),
		handleValuePropertyChange[bool](d, "replace_invalid_characters", &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleValuePropertyChange[string](d, "default_ddl_collation", &set.DefaultDDLCollation, &unset.DefaultDDLCollation),
		handleValuePropertyChangeWithMapping[string](d, "storage_serialization_policy", &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, func(value string) sdk.StorageSerializationPolicy { return sdk.StorageSerializationPolicy(value) }),
		handleValuePropertyChangeWithMapping[string](d, "log_level", &set.LogLevel, &unset.LogLevel, func(value string) sdk.LogLevel { return sdk.LogLevel(value) }),
		handleValuePropertyChangeWithMapping[string](d, "trace_level", &set.TraceLevel, &unset.TraceLevel, func(value string) sdk.TraceLevel { return sdk.TraceLevel(value) }),
		handleValuePropertyChange[int](d, "suspend_task_after_num_failures", &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleValuePropertyChange[int](d, "task_auto_retry_attempts", &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleValuePropertyChangeWithMapping[string](d, "user_task_managed_initial_warehouse_size", &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, func(value string) sdk.WarehouseSize { return sdk.WarehouseSize(value) }),
		handleValuePropertyChange[int](d, "user_task_timeout_ms", &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleValuePropertyChange[int](d, "user_task_minimum_trigger_interval_in_seconds", &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleValuePropertyChange[bool](d, "quoted_identifiers_ignore_case", &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleValuePropertyChange[bool](d, "enable_console_output", &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
	)
}

// handleValuePropertyChange calls internally handleValuePropertyChangeWithMapping with identity mapping
func handleValuePropertyChange[T any](d *schema.ResourceData, key string, setField **T, unsetField **bool) diag.Diagnostics {
	return handleValuePropertyChangeWithMapping[T, T](d, key, setField, unsetField, func(value T) T { return value })
}

// handleValuePropertyChangeWithMapping checks schema.ResourceData for change in key's value. If there's a change detected
// (or unknown value that basically indicates diff.SetNewComputed was called on the key), it checks if the value is set in the configuration.
// If the value is set, setField (representing setter for a value) is set to the new planned value applying mapping beforehand in cases where enum values,
// identifiers, etc. have to be set. Otherwise, unsetField is populated.
func handleValuePropertyChangeWithMapping[T, R any](d *schema.ResourceData, key string, setField **R, unsetField **bool, mapping func(value T) R) diag.Diagnostics {
	if d.HasChange(key) || !d.GetRawPlan().AsValueMap()[key].IsKnown() {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			*setField = sdk.Pointer(mapping(d.Get(key).(T)))
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
