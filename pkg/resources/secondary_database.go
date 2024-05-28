package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secondaryDatabaseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the database; must be unique for your account. As a best practice for [Database Replication and Failover](https://docs.snowflake.com/en/user-guide/db-replication-intro), it is recommended to give each secondary database the same name as its primary database. This practice supports referencing fully-qualified objects (i.e. '<db>.<schema>.<object>') by other objects in the same database, such as querying a fully-qualified table name in a view. If a secondary database has a different name from the primary database, then these object references would break in the secondary database.",
	},
	"as_replica_of": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "A fully qualified path to a database to create a replica from. A fully qualified path follows the format of `\"<organization_name>\".\"<account_name>\".\"<database_name>\"`.",
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
	// TODO: Below parameters should be nested properties
	"external_volume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Description:      "The database parameter that specifies the default external volume to use for Iceberg tables.",
	},
	"catalog": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Description:      "The database parameter that specifies the default catalog to use for Iceberg tables.",
	},
	"replace_invalid_characters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�) in query results for an Iceberg table. You can only set this parameter for tables that use an external Iceberg catalog.",
	},
	"default_ddl_collation": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
	},
	"storage_serialization_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllStorageSerializationPolicies), true),
		Description:      fmt.Sprintf("Specifies the storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(string) == string(sdk.StorageSerializationPolicyOptimized) && newValue == ""
		},
	},
	"log_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllLogLevels), true),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(string) == string(sdk.LogLevelOff) && newValue == ""
		},
		Description: fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
	},
	"trace_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllTraceLevels), true),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(string) == string(sdk.TraceLevelOff) && newValue == ""
		},
		Description: fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", sdk.AsStringList(sdk.AllTraceLevels)),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
}

func SecondaryDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateSecondaryDatabase,
		UpdateContext: UpdateSecondaryDatabase,
		ReadContext:   ReadSecondaryDatabase,
		DeleteContext: DeleteSecondaryDatabase,
		Description:   "A secondary database creates a replica of an existing primary database (i.e. a secondary database). For more information about database replication, see [Introduction to database replication across multiple accounts](https://docs.snowflake.com/en/user-guide/db-replication-intro).",

		CustomizeDiff: customdiff.All(
			NestedIntValueAccountObjectComputedIf("data_retention_time_in_days", sdk.AccountParameterDataRetentionTimeInDays),
			NestedIntValueAccountObjectComputedIf("max_data_extension_time_in_days", sdk.AccountParameterMaxDataExtensionTimeInDays),
		),

		Schema: secondaryDatabaseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	secondaryDatabaseId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
	primaryDatabaseId := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(d.Get("as_replica_of").(string))

	dataRetentionTimeInDays, _ := GetPropertyOfFirstNestedObjectByKey[int](d, "data_retention_time_in_days", "value")
	maxDataExtensionTimeInDays, _ := GetPropertyOfFirstNestedObjectByKey[int](d, "max_data_extension_time_in_days", "value")

	var externalVolume *sdk.AccountObjectIdentifier
	if v, ok := d.GetOk("external_volume"); ok {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	var catalog *sdk.AccountObjectIdentifier
	if v, ok := d.GetOk("catalog"); ok {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	var storageSerializationPolicy *sdk.StorageSerializationPolicy
	if v, ok := d.GetOk("storage_serialization_policy"); ok {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(v.(string)))
	}

	var logLevel *sdk.LogLevel
	if v, ok := d.GetOk("log_level"); ok {
		logLevel = sdk.Pointer(sdk.LogLevel(v.(string)))
	}

	var traceLevel *sdk.TraceLevel
	if v, ok := d.GetOk("trace_level"); ok {
		traceLevel = sdk.Pointer(sdk.TraceLevel(v.(string)))
	}

	err := client.Databases.CreateSecondary(ctx, secondaryDatabaseId, primaryDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
		Transient:                  GetPropertyAsPointer[bool](d, "is_transient"),
		DataRetentionTimeInDays:    dataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: maxDataExtensionTimeInDays,
		ExternalVolume:             externalVolume,
		Catalog:                    catalog,
		ReplaceInvalidCharacters:   GetPropertyAsPointer[bool](d, "replace_invalid_characters"),
		DefaultDDLCollation:        GetPropertyAsPointer[string](d, "default_ddl_collation"),
		StorageSerializationPolicy: storageSerializationPolicy,
		LogLevel:                   logLevel,
		TraceLevel:                 traceLevel,
		Comment:                    GetPropertyAsPointer[string](d, "comment"),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(secondaryDatabaseId))

	return ReadSecondaryDatabase(ctx, d, meta)
}

func UpdateSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		externalVolume := d.Get("external_volume").(string)
		if len(externalVolume) > 0 {
			databaseSetRequest.ExternalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(externalVolume))
		} else {
			databaseUnsetRequest.ExternalVolume = sdk.Bool(true)
		}
	}

	if d.HasChange("catalog") {
		catalog := d.Get("catalog").(string)
		if len(catalog) > 0 {
			databaseSetRequest.Catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(catalog))
		} else {
			databaseUnsetRequest.Catalog = sdk.Bool(true)
		}
	}

	if d.HasChange("replace_invalid_characters") {
		if d.Get("replace_invalid_characters").(bool) {
			databaseSetRequest.ReplaceInvalidCharacters = sdk.Bool(true)
		} else {
			databaseUnsetRequest.ReplaceInvalidCharacters = sdk.Bool(true)
		}
	}

	if d.HasChange("default_ddl_collation") {
		defaultDdlCollation := d.Get("default_ddl_collation").(string)
		if len(defaultDdlCollation) > 0 {
			databaseSetRequest.DefaultDDLCollation = &defaultDdlCollation
		} else {
			databaseUnsetRequest.DefaultDDLCollation = sdk.Bool(true)
		}
	}

	if d.HasChange("storage_serialization_policy") {
		storageSerializationPolicy := d.Get("storage_serialization_policy").(string)
		if len(storageSerializationPolicy) > 0 {
			databaseSetRequest.StorageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(storageSerializationPolicy))
		} else {
			databaseUnsetRequest.StorageSerializationPolicy = sdk.Bool(true)
		}
	}

	if d.HasChange("log_level") {
		logLevel := d.Get("log_level").(string)
		if len(logLevel) > 0 {
			databaseSetRequest.LogLevel = sdk.Pointer(sdk.LogLevel(logLevel))
		} else {
			databaseUnsetRequest.LogLevel = sdk.Bool(true)
		}
	}

	if d.HasChange("trace_level") {
		traceLevel := d.Get("trace_level").(string)
		if len(traceLevel) > 0 {
			databaseSetRequest.TraceLevel = sdk.Pointer(sdk.TraceLevel(traceLevel))
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

	return ReadSecondaryDatabase(ctx, d, meta)
}

func ReadSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	secondaryDatabaseId := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	secondaryDatabase, err := client.Databases.ShowByID(ctx, secondaryDatabaseId)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query secondary database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("DatabaseName: %s, Err: %s", secondaryDatabaseId.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	secondaryDatabaseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: secondaryDatabaseId,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, &sdk.ShowReplicationDatabasesOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(secondaryDatabaseId.Name()),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var replicationPrimaryDatabase *sdk.ReplicationDatabase
	for _, replicationDatabase := range replicationDatabases {
		replicationDatabase := replicationDatabase
		if !replicationDatabase.IsPrimary &&
			replicationDatabase.AccountLocator == client.GetAccountLocator() &&
			replicationDatabase.Name == secondaryDatabaseId.Name() {
			replicationPrimaryDatabase = &replicationDatabase
		}
	}
	if replicationPrimaryDatabase == nil {
		return diag.FromErr(fmt.Errorf("could not find replication database for %s", secondaryDatabaseId.Name()))
	}

	if err := d.Set("name", secondaryDatabase.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("as_replica_of", sdk.NewExternalObjectIdentifierFromFullyQualifiedName(replicationPrimaryDatabase.PrimaryDatabase).FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_transient", secondaryDatabase.Transient); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("data_retention_time_in_days", []any{map[string]any{"value": secondaryDatabase.RetentionTime}}); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", secondaryDatabase.Comment); err != nil {
		return diag.FromErr(err)
	}

	for _, secondaryDatabaseParameter := range secondaryDatabaseParameters {
		switch secondaryDatabaseParameter.Key {
		case "MAX_DATA_EXTENSION_TIME_IN_DAYS":
			maxDataExtensionTimeInDays, err := strconv.Atoi(secondaryDatabaseParameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("max_data_extension_time_in_days", []any{map[string]any{"value": maxDataExtensionTimeInDays}}); err != nil {
				return diag.FromErr(err)
			}
		case "EXTERNAL_VOLUME":
			if err := d.Set("external_volume", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "CATALOG":
			if err := d.Set("catalog", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "DEFAULT_DDL_COLLATION":
			if err := d.Set("default_ddl_collation", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "LOG_LEVEL":
			if err := d.Set("log_level", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "TRACE_LEVEL":
			if err := d.Set("trace_level", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "REPLACE_INVALID_CHARACTERS":
			boolValue, err := strconv.ParseBool(secondaryDatabaseParameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("replace_invalid_characters", boolValue); err != nil {
				return diag.FromErr(err)
			}
		case "STORAGE_SERIALIZATION_POLICY":
			if err := d.Set("storage_serialization_policy", secondaryDatabaseParameter.Value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func DeleteSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
