package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

// TODO: Add refresh in secondary database in Read

var databaseV1Schema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the database; must be unique for your account. As a best practice for [Database Replication and Failover](https://docs.snowflake.com/en/user-guide/db-replication-intro), it is recommended to give each secondary database the same name as its primary database. This practice supports referencing fully-qualified objects (i.e. '<db>.<schema>.<object>') by other objects in the same database, such as querying a fully-qualified table name in a view. If a secondary database has a different name from the primary database, then these object references would break in the secondary database.",
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the database as transient. Transient databases do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	"data_retention_time_in_days": nestedProperty(
		schema.TypeInt,
		"Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
	),
	"max_data_extension_time_in_days": nestedProperty(
		schema.TypeInt,
		"Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale. For a detailed description of this parameter, see [MAX_DATA_EXTENSION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters.html#label-max-data-extension-time-in-days).",
	),
	"external_volume": nestedPropertyWithInnerModifier(
		schema.TypeString,
		"The database parameter that specifies the default external volume to use for Iceberg tables.",
		func(inner *schema.Schema) {
			inner.ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
		},
	),
	"catalog": nestedPropertyWithInnerModifier(
		schema.TypeString,
		"The database parameter that specifies the default catalog to use for Iceberg tables.",
		func(inner *schema.Schema) {
			inner.ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
		},
	),
	"replace_invalid_characters": nestedProperty(
		schema.TypeBool,
		"Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog.",
	),
	"default_ddl_collation": nestedProperty(
		schema.TypeString,
		"Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
	),
	"storage_serialization_policy": nestedProperty(
		schema.TypeString,
		fmt.Sprintf("Specifies the storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
	),
	"log_level": nestedProperty(
		schema.TypeString,
		fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
	),
	"trace_level": nestedProperty(
		schema.TypeString,
		fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", sdk.AsStringList(sdk.AllTraceLevels)),
	),
	// TODO: failover (should It be integration into "replicate" field or its own)
	//"replicate": {
	//	Type:     schema.TypeList,
	//	MaxItems: 1,
	//	Elem: &schema.Resource{
	//		Schema: map[string]*schema.Schema{
	//			"accounts_enabled_for_replication": {
	//				Type: schema.TypeList,
	//				Elem: &schema.Schema{
	//					Type: schema.TypeString,
	//				},
	//				Description: "TODO",
	//			},
	//			"accounts_disabled_for_replication": {
	//				Type: schema.TypeList,
	//				Elem: &schema.Schema{
	//					Type: schema.TypeString,
	//				},
	//				Description: "TODO",
	//			},
	//			"ignore_edition_check": {
	//				Type: schema.TypeBool,
	//			},
	//		},
	//	},
	//	Computed:    true,
	//	Optional:    true,
	//	Description: "TODO",
	//},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
}

func StandardDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateStandardDatabase,
		ReadContext:   ReadStandardDatabase,
		DeleteContext: DeleteStandardDatabase,
		UpdateContext: UpdateStandardDatabase,
		// TODO CustomDiffs

		// TODO: Desc
		Description: "",
		Schema:      databaseV1Schema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	dataRetentionTimeInDays, _ := GetPropertyOfFirstNestedObjectByKey[int](d, "data_retention_time_in_days", "value")
	maxDataExtensionTimeInDays, _ := GetPropertyOfFirstNestedObjectByKey[int](d, "max_data_extension_time_in_days", "value")
	replaceInvalidCharacters, _ := GetPropertyOfFirstNestedObjectByKey[bool](d, "replace_invalid_characters", "value")
	defaultDdlCollation, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "default_ddl_collation", "value")

	var externalVolume *sdk.AccountObjectIdentifier
	if externalVolumeRaw, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "external_volume", "value"); externalVolumeRaw != nil {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*externalVolumeRaw))
	}

	var catalog *sdk.AccountObjectIdentifier
	if catalogRaw, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "catalog", "value"); catalogRaw != nil {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*catalogRaw))
	}

	var storageSerializationPolicy *sdk.StorageSerializationPolicy
	if storageSerializationPolicyRaw, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "storage_serialization_policy", "value"); storageSerializationPolicyRaw != nil {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*storageSerializationPolicyRaw))
	}

	var logLevel *sdk.LogLevel
	if logLevelRaw, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "log_level", "value"); logLevelRaw != nil {
		logLevel = sdk.Pointer(sdk.LogLevel(*logLevelRaw))
	}

	var traceLevel *sdk.TraceLevel
	if traceLevelRaw, _ := GetPropertyOfFirstNestedObjectByKey[string](d, "trace_level", "value"); traceLevelRaw != nil {
		traceLevel = sdk.Pointer(sdk.TraceLevel(*traceLevelRaw))
	}

	err := client.Databases.Create(ctx, id, &sdk.CreateDatabaseOptions{
		Transient:                  GetPropertyAsPointer[bool](d, "is_transient"),
		DataRetentionTimeInDays:    dataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: maxDataExtensionTimeInDays,
		ExternalVolume:             externalVolume,
		Catalog:                    catalog,
		ReplaceInvalidCharacters:   replaceInvalidCharacters,
		DefaultDDLCollation:        defaultDdlCollation,
		StorageSerializationPolicy: storageSerializationPolicy,
		LogLevel:                   logLevel,
		TraceLevel:                 traceLevel,
		Comment:                    GetPropertyAsPointer[string](d, "comment"),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	// TODO: Alter Replication and failover
	//err = client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
	//	EnableReplication:  nil,
	//	DisableReplication: nil,
	//	Refresh:            nil,
	//})
	//if err != nil {
	//	// TODO: Return error or warning ?
	//	return diag.FromErr(err)
	//}
	//
	//err = client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
	//	EnableFailover:  nil,
	//	DisableFailover: nil,
	//	Primary:         nil,
	//})
	//if err != nil {
	//	// TODO: Return error or warning ?
	//	return diag.FromErr(err)
	//}

	return ReadStandardDatabase(ctx, d, meta)
}

func UpdateStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	secondaryDatabaseId := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newId))
		secondaryDatabaseId = newId
	}

	var databaseSetRequest sdk.DatabaseSet
	var databaseUnsetRequest sdk.DatabaseUnset

	if d.HasChange("data_retention_time_in_days") {
		dataRetentionObject, ok := d.GetOk("data_retention_time_in_days")
		if ok && len(dataRetentionObject.([]any)) > 0 {
			dataRetentionTimeInDays, err := GetPropertyOfFirstNestedObjectByKey[int](d, "data_retention_time_in_days", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.DataRetentionTimeInDays = dataRetentionTimeInDays
		} else {
			databaseUnsetRequest.DataRetentionTimeInDays = sdk.Bool(true)
		}
	}

	if d.HasChange("max_data_extension_time_in_days") {
		maxDataExtensionTimeInDaysObject, ok := d.GetOk("max_data_extension_time_in_days")
		if ok && len(maxDataExtensionTimeInDaysObject.([]any)) > 0 {
			maxDataExtensionTimeInDays, err := GetPropertyOfFirstNestedObjectByKey[int](d, "max_data_extension_time_in_days", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.MaxDataExtensionTimeInDays = maxDataExtensionTimeInDays
		} else {
			databaseUnsetRequest.MaxDataExtensionTimeInDays = sdk.Bool(true)
		}
	}

	if d.HasChange("external_volume") {
		externalVolumeObject, ok := d.GetOk("external_volume")
		if ok && len(externalVolumeObject.([]any)) > 0 {
			externalVolume, err := GetPropertyOfFirstNestedObjectByKey[string](d, "external_volume", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.ExternalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*externalVolume))
		} else {
			databaseUnsetRequest.ExternalVolume = sdk.Bool(true)
		}
	}

	if d.HasChange("catalog") {
		catalogObject, ok := d.GetOk("catalog")
		if ok && len(catalogObject.([]any)) > 0 {
			catalog, err := GetPropertyOfFirstNestedObjectByKey[string](d, "catalog", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.Catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*catalog))
		} else {
			databaseUnsetRequest.Catalog = sdk.Bool(true)
		}
	}

	if d.HasChange("replace_invalid_characters") {
		replaceInvalidCharactersObject, ok := d.GetOk("replace_invalid_characters")
		if ok && len(replaceInvalidCharactersObject.([]any)) > 0 {
			replaceInvalidCharacters, err := GetPropertyOfFirstNestedObjectByKey[bool](d, "replace_invalid_characters", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.ReplaceInvalidCharacters = sdk.Bool(*replaceInvalidCharacters)
		} else {
			databaseUnsetRequest.ReplaceInvalidCharacters = sdk.Bool(true)
		}
	}

	if d.HasChange("default_ddl_collation") {
		defaultDdlCollationObject, ok := d.GetOk("default_ddl_collation")
		if ok && len(defaultDdlCollationObject.([]any)) > 0 {
			defaultDdlCollation, err := GetPropertyOfFirstNestedObjectByKey[string](d, "default_ddl_collation", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.DefaultDDLCollation = defaultDdlCollation
		} else {
			databaseUnsetRequest.DefaultDDLCollation = sdk.Bool(true)
		}
	}

	if d.HasChange("storage_serialization_policy") {
		storageSerializationPolicyObject, ok := d.GetOk("storage_serialization_policy")
		if ok && len(storageSerializationPolicyObject.([]any)) > 0 {
			storageSerializationPolicy, err := GetPropertyOfFirstNestedObjectByKey[string](d, "storage_serialization_policy", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.StorageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*storageSerializationPolicy))
		} else {
			databaseUnsetRequest.StorageSerializationPolicy = sdk.Bool(true)
		}
	}

	if d.HasChange("log_level") {
		logLevelObject, ok := d.GetOk("log_level")
		if ok && len(logLevelObject.([]any)) > 0 {
			logLevel, err := GetPropertyOfFirstNestedObjectByKey[string](d, "log_level", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.LogLevel = sdk.Pointer(sdk.LogLevel(*logLevel))
		} else {
			databaseUnsetRequest.LogLevel = sdk.Bool(true)
		}
	}

	if d.HasChange("trace_level") {
		traceLevelObject, ok := d.GetOk("trace_level")
		if ok && len(traceLevelObject.([]any)) > 0 {
			traceLevel, err := GetPropertyOfFirstNestedObjectByKey[string](d, "trace_level", "value")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.TraceLevel = sdk.Pointer(sdk.TraceLevel(*traceLevel))
		} else {
			databaseUnsetRequest.TraceLevel = sdk.Bool(true)
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			databaseSetRequest.Comment = &comment
		} else {
			databaseUnsetRequest.Comment = sdk.Bool(true)
		}
	}

	if (databaseSetRequest != sdk.DatabaseSet{}) {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Set: &databaseSetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (databaseUnsetRequest != sdk.DatabaseUnset{}) {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Unset: &databaseUnsetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadStandardDatabase(ctx, d, meta)
}

func ReadStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	database, err := client.Databases.ShowByID(ctx, id)
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

	if err := d.Set("name", database.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_transient", database.Transient); err != nil {
		return diag.FromErr(err)
	}

	if err := SetPropertyOfFirstNestedObjectByKey(d, "data_retention_time_in_days", "value", database.RetentionTime); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", database.Comment); err != nil {
		return diag.FromErr(err)
	}

	parameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: id,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	//replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, &sdk.ShowReplicationDatabasesOptions{
	//	Like: &sdk.Like{
	//		Pattern: sdk.String(secondaryDatabaseId.Name()),
	//	},
	//})
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	// TODO: Set enabled/disabled accounts for replication/failover

	for _, parameter := range parameters {
		switch parameter.Key {
		case "MAX_DATA_EXTENSION_TIME_IN_DAYS":
			maxDataExtensionTimeInDays, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			// TODO: I think "value" could be assumed and the option for other value names should be opt-in
			if err := SetPropertyOfFirstNestedObjectByKey(d, "max_data_extension_time_in_days", "value", maxDataExtensionTimeInDays); err != nil {
				return diag.FromErr(err)
			}
		case "EXTERNAL_VOLUME":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "external_volume", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "CATALOG":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "catalog", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "DEFAULT_DDL_COLLATION":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "default_ddl_collation", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "LOG_LEVEL":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "log_level", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "TRACE_LEVEL":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "trace_level", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "REPLACE_INVALID_CHARACTERS":
			boolValue, err := strconv.ParseBool(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := SetPropertyOfFirstNestedObjectByKey(d, "replace_invalid_characters", "value", boolValue); err != nil {
				return diag.FromErr(err)
			}
		case "STORAGE_SERIALIZATION_POLICY":
			if err := SetPropertyOfFirstNestedObjectByKey(d, "storage_serialization_policy", "value", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func DeleteStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	err := client.Databases.Drop(ctx, id, &sdk.DropDatabaseOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
