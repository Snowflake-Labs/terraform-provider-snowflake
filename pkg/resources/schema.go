package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
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
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the schema; must be unique for the database in which the schema is created.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The database in which to create the schema.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"with_managed_access": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      booleanStringFieldDescription("Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
		}),
	},
	"is_transient": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      booleanStringFieldDescription("Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss."),
		ValidateDiagFunc: validateBooleanString,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("options", func(x any) any {
			return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
		}),
	},
	strings.ToLower(string(sdk.ObjectParameterPipeExecutionPaused)): {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Specifies whether to pause a running pipe, primarily in preparation for transferring ownership of the pipe to a different role. For more information, see [PIPE_EXECUTION_PAUSED](https://docs.snowflake.com/en/sql-reference/parameters#pipe-execution-paused).",
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
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SCHEMA` for the given object. In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs/resources/grant_ownership) on all objects in the schema.",
		Elem: &schema.Resource{
			Schema: schemas.SchemaDescribeSchema,
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
		ReadContext:   ReadContextSchema(true),
		UpdateContext: UpdateContextSchema,
		DeleteContext: DeleteContextSchema,
		Description:   "Resource used to manage schema objects. For more information, check [schema documentation](https://docs.snowflake.com/en/sql-reference/sql/create-schema).",

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "name", "comment", "with_managed_access", "is_transient"),
			ComputedIfAnyAttributeChanged(DescribeOutputAttributeName, "name"),
			ComputedIfAnyAttributeChanged(ParametersAttributeName,
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
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterExternalVolume, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterCatalog, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterDefaultDDLCollation, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterLogLevel, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterTraceLevel, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeSchema},
				parameter[sdk.ObjectParameter]{sdk.ObjectParameterPipeExecutionPaused, valueTypeBool, sdk.ParameterTypeSchema},
			),
		),

		Schema: helpers.MergeMaps(schemaSchema, DatabaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: ImportSchema,
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

func ImportSchema(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting schema import")
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)

	s, err := client.Schemas.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("name", s.Name); err != nil {
		return nil, err
	}

	if err := d.Set("database", s.DatabaseName); err != nil {
		return nil, err
	}

	if err := d.Set("comment", s.Comment); err != nil {
		return nil, err
	}

	if err := d.Set("is_transient", booleanStringFromBool(s.IsTransient())); err != nil {
		return nil, err
	}

	if err := d.Set("with_managed_access", booleanStringFromBool(s.IsManagedAccess())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func schemaParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)
	return client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Schema: id,
		},
	})
}

func CreateContextSchema(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	id := sdk.NewDatabaseObjectIdentifier(database, name)

	opts := &sdk.CreateSchemaOptions{
		Comment: GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"),
	}
	if parametersCreateDiags := handleSchemaParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

	if v := d.Get("is_transient").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Transient = sdk.Bool(parsed)
	}
	if v := d.Get("with_managed_access").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		opts.WithManagedAccess = sdk.Bool(parsed)
	}
	if strings.EqualFold(strings.TrimSpace(name), "PUBLIC") {
		opts.OrReplace = sdk.Pointer(true)
	}
	if err := client.Schemas.Create(ctx, id, opts); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create schema.",
				Detail:   fmt.Sprintf("schema name: %s, err: %s", id.FullyQualifiedName(), err),
			},
		}
	}
	d.SetId(helpers.EncodeSnowflakeID(database, name))

	return ReadContextSchema(false)(ctx, d, meta)
}

func ReadContextSchema(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id := helpers.DecodeSnowflakeID(d.Id()).(sdk.DatabaseObjectIdentifier)

		_, err := client.Databases.ShowByID(ctx, id.DatabaseId())
		if err != nil {
			log.Printf("[DEBUG] database %s for schema %s not found", id.DatabaseId().Name(), id.Name())
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("database name: %s, Err: %s", id.DatabaseId(), err),
				},
			}
		}

		schema, err := client.Schemas.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query schema. Marking the resource as removed.",
						Detail:   fmt.Sprintf("schema name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		if err := d.Set("name", schema.Name); err != nil {
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
		if err := d.Set(strings.ToLower(string(sdk.ObjectParameterPipeExecutionPaused)), value); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				showMapping{"options", "is_transient", schema.IsTransient(), booleanStringFromBool(schema.IsTransient()), func(x any) any {
					return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "TRANSIENT")
				}},
				showMapping{"options", "with_managed_access", schema.IsManagedAccess(), booleanStringFromBool(schema.IsManagedAccess()), func(x any) any {
					return slices.Contains(sdk.ParseCommaSeparatedStringArray(x.(string), false), "MANAGED ACCESS")
				}},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, schemaSchema, []string{
			"is_transient",
			"with_managed_access",
		}); err != nil {
			return diag.FromErr(err)
		}

		describeResult, err := client.Schemas.Describe(ctx, schema.ID())
		if err != nil {
			log.Printf("[DEBUG] describing schema: %s, err: %s", id.FullyQualifiedName(), err)
		} else {
			if err = d.Set(DescribeOutputAttributeName, schemas.SchemaDescriptionToSchema(describeResult)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SchemaToSchema(schema)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ParametersAttributeName, []map[string]any{schemas.SchemaParametersToSchema(schemaParameters)}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
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
		if v := d.Get("with_managed_access").(string); v != BooleanDefault {
			var err error
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			if parsed {
				err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					EnableManagedAccess: sdk.Pointer(true),
				})
			} else {
				err = client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
					DisableManagedAccess: sdk.Pointer(true),
				})
			}
			if err != nil {
				return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
			}
		} else {
			// managed access can not be UNSET to a default value
			if err := client.Schemas.Alter(ctx, id, &sdk.AlterSchemaOptions{
				DisableManagedAccess: sdk.Pointer(true),
			}); err != nil {
				return diag.FromErr(fmt.Errorf("error handling with_managed_access on %v err = %w", d.Id(), err))
			}
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

	if updateParamDiags := handleSchemaParametersChanges(d, set, unset); len(updateParamDiags) > 0 {
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

	return ReadContextSchema(false)(ctx, d, meta)
}

func handleSchemaParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateSchemaOptions) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ObjectParameterDataRetentionTimeInDays, &createOpts.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &createOpts.MaxDataExtensionTimeInDays),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterExternalVolume, &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterCatalog, &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.ObjectParameterPipeExecutionPaused, &createOpts.PipeExecutionPaused),
		handleParameterCreate(d, sdk.ObjectParameterReplaceInvalidCharacters, &createOpts.ReplaceInvalidCharacters),
		handleParameterCreate(d, sdk.ObjectParameterDefaultDDLCollation, &createOpts.DefaultDDLCollation),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterLogLevel, &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterTraceLevel, &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.ObjectParameterEnableConsoleOutput, &createOpts.EnableConsoleOutput),
	)
}

func handleSchemaParametersChanges(d *schema.ResourceData, set *sdk.SchemaSet, unset *sdk.SchemaUnset) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ObjectParameterDataRetentionTimeInDays, &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterExternalVolume, &set.ExternalVolume, &unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterCatalog, &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, sdk.ObjectParameterPipeExecutionPaused, &set.PipeExecutionPaused, &unset.PipeExecutionPaused),
		handleParameterUpdate(d, sdk.ObjectParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdate(d, sdk.ObjectParameterDefaultDDLCollation, &set.DefaultDDLCollation, &unset.DefaultDDLCollation),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterLogLevel, &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterTraceLevel, &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskTimeoutMs, &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, sdk.ObjectParameterEnableConsoleOutput, &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
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
