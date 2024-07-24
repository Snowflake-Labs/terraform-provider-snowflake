package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var schemaSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the schema.",
		ForceNew:    true,
	},
	"with_managed_access": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies a schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	strings.ToLower(string(sdk.ObjectParameterPipeExecutionPaused)): {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Specifies whether to pause a running pipe, primarily in preparation for transferring ownership of the pipe to a different role.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the schema.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaParametersSchema,
		},
	},
}

// Schema returns a pointer to the resource representing a schema.
func Schema() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: CreateContextSchema,
		ReadContext:   ReadContextSchema,
		UpdateContext: UpdateContextSchema,
		DeleteContext: DeleteContextSchema,

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "comment", "with_managed_access", "is_transient"),
			ComputedIfAnyAttributeChanged(ParametersAttributeName,
				strings.ToLower(string(sdk.ObjectParameterMaxConcurrencyLevel)),
				strings.ToLower(string(sdk.ObjectParameterStatementQueuedTimeoutInSeconds)),
				strings.ToLower(string(sdk.ObjectParameterStatementTimeoutInSeconds)),
				strings.ToLower(string(sdk.ObjectParameterDataRetentionTimeInDays)),
				strings.ToLower(string(sdk.ObjectParameterMaxDataExtensionTimeInDays)),
				strings.ToLower(string(sdk.ObjectParameterExternalVolume)),
				strings.ToLower(string(sdk.ObjectParameterCatalog)),
				strings.ToLower(string(sdk.ObjectParameterReplaceInvalidCharacters)),
				strings.ToLower(string(sdk.ObjectParameterDefaultDDLCollation)),
				strings.ToLower(string(sdk.ObjectParameterStorageSerializationPolicy)),
				strings.ToLower(string(sdk.ObjectParameterLogLevel)),
				strings.ToLower(string(sdk.ObjectParameterTraceLevel)),
				strings.ToLower(string(sdk.ObjectParameterSuspendTaskAfterNumFailures)),
				strings.ToLower(string(sdk.ObjectParameterTaskAutoRetryAttempts)),
				strings.ToLower(string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize)),
				strings.ToLower(string(sdk.ObjectParameterUserTaskTimeoutMs)),
				strings.ToLower(string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds)),
				strings.ToLower(string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase)),
				strings.ToLower(string(sdk.ObjectParameterEnableConsoleOutput)),
				strings.ToLower(string(sdk.ObjectParameterPipeExecutionPaused)),
			),
			ParametersCustomDiff(
				schemaParametersProvider,
				parameter{sdk.AccountParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterExternalVolume, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterCatalog, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterDefaultDDLCollation, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterLogLevel, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterTraceLevel, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeSchema},
				parameter{sdk.AccountParameterPipeExecutionPaused, valueTypeBool, sdk.ParameterTypeSchema},
			),
		),

		Schema: helpers.MergeMaps(schemaSchema, DatabaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v093SchemaStateUpgrader,
			},
		},
	}
}

func schemaParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)
	warehouseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Schema: id,
		},
	})
	if err != nil {
		return nil, err
	}
	return warehouseParameters, nil
}

func CreateContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	id := sdk.NewDatabaseObjectIdentifier(database, name)

	dataRetentionTimeInDays,
		maxDataExtensionTimeInDays,
		externalVolume,
		catalog,
		replaceInvalidCharacters,
		defaultDDLCollation,
		storageSerializationPolicy,
		logLevel,
		traceLevel,
		suspendTaskAfterNumFailures,
		taskAutoRetryAttempts,
		userTaskManagedInitialWarehouseSize,
		userTaskTimeoutMs,
		userTaskMinimumTriggerIntervalInSeconds,
		quotedIdentifiersIgnoreCase,
		enableConsoleOutput,
		err := GetAllDatabaseParameters(d)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &sdk.CreateSchemaOptions{
		Transient:                               GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "is_transient"),
		WithManagedAccess:                       GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "with_managed_access"),
		DataRetentionTimeInDays:                 dataRetentionTimeInDays,
		MaxDataExtensionTimeInDays:              maxDataExtensionTimeInDays,
		ExternalVolume:                          externalVolume,
		Catalog:                                 catalog,
		ReplaceInvalidCharacters:                replaceInvalidCharacters,
		DefaultDDLCollation:                     defaultDDLCollation,
		StorageSerializationPolicy:              storageSerializationPolicy,
		LogLevel:                                logLevel,
		TraceLevel:                              traceLevel,
		SuspendTaskAfterNumFailures:             suspendTaskAfterNumFailures,
		TaskAutoRetryAttempts:                   taskAutoRetryAttempts,
		UserTaskManagedInitialWarehouseSize:     userTaskManagedInitialWarehouseSize,
		UserTaskTimeoutMs:                       userTaskTimeoutMs,
		UserTaskMinimumTriggerIntervalInSeconds: userTaskMinimumTriggerIntervalInSeconds,
		QuotedIdentifiersIgnoreCase:             quotedIdentifiersIgnoreCase,
		EnableConsoleOutput:                     enableConsoleOutput,
		PipeExecutionPaused:                     GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "pipe_execution_paused"),
		Comment:                                 GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"),
	}
	if strings.EqualFold(strings.TrimSpace(name), "PUBLIC") {
		opts.OrReplace = sdk.Pointer(true)
	}
	if err := client.Schemas.Create(ctx, id, opts); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(database, name))

	return ReadContextSchema(ctx, d, meta)
}

func ReadContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)

	_, err := client.Databases.ShowByID(ctx, id.DatabaseId())
	if err != nil {
		d.SetId("")
	}

	schema, err := client.Schemas.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query secondary database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("DatabaseName: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", schema.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_transient", schema.IsTransient()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("with_managed_access", schema.IsManagedAccess()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", schema.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", schema.Comment); err != nil {
		return diag.FromErr(err)
	}

	schemaParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Schema: id,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := HandleDatabaseParameterRead(d, schemaParameters); diags != nil {
		return diags
	}
	pipeExecutionPaused, err := collections.FindOne(schemaParameters, func(property *sdk.Parameter) bool {
		return property.Key == "PIPE_EXECUTION_PAUSED"
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to find schema PIPE_EXECUTION_PAUSED parameter, err = %w", err))
	}
	value, err := strconv.ParseBool((*pipeExecutionPaused).Value)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pipe_execution_paused", value); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SchemaToSchema(schema)}); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ParametersAttributeName, []map[string]any{schemas.SchemaParametersToSchema(schemaParameters)}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)
	client := meta.(*provider.Context).Client

	if d.HasChange("name") {
		newId := sdk.NewDatabaseObjectIdentifier(d.Get("database").(string), d.Get("name").(string))
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			NewName: sdk.Pointer(newId),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	if d.HasChange("with_managed_access") {
		var err error
		if _, ok := d.GetOk("with_managed_access"); ok {
			err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				EnableManagedAccess: sdk.Pointer(true),
			})
		} else {
			err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				DisableManagedAccess: sdk.Pointer(true),
			})
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("error changing management state on %v err = %w", d.Id(), err))
		}
	}

	set := new(sdk.SchemaSet)
	unset := new(sdk.SchemaUnset)

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			set.Comment = &comment
		} else {
			unset.Comment = sdk.Bool(true)
		}
	}

	if updateParamDiags := HandleSchemaParametersChanges(d, set, unset); len(updateParamDiags) > 0 {
		return updateParamDiags
	}
	if (*set != sdk.SchemaSet{}) {
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Set: set,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.SchemaUnset{}) {
		err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
			Unset: unset,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSchema(ctx, d, meta)
}

func HandleSchemaParametersChanges(d *schema.ResourceData, set *sdk.SchemaSet, unset *sdk.SchemaUnset) diag.Diagnostics {
	return JoinDiags(
		handleValuePropertyChange[int](d, "data_retention_time_in_days", &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleValuePropertyChange[int](d, "max_data_extension_time_in_days", &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleValuePropertyChangeWithMapping[string](d, "external_volume", &set.ExternalVolume, &unset.ExternalVolume, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleValuePropertyChangeWithMapping[string](d, "catalog", &set.Catalog, &unset.Catalog, func(value string) (sdk.AccountObjectIdentifier, error) {
			return sdk.NewAccountObjectIdentifier(value), nil
		}),
		handleValuePropertyChange[bool](d, "pipe_execution_paused", &set.PipeExecutionPaused, &unset.PipeExecutionPaused),
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

func DeleteContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)
	client := meta.(*provider.Context).Client

	err := client.Schemas.Drop(ctx, id, &sdk.DropSchemaOptions{IfExists: sdk.Pointer(true)})
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting schema",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
