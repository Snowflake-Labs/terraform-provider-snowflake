package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var standardDatabaseSchema = map[string]*schema.Schema{
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
	"storage_serialization_policy": nestedPropertyWithInnerModifier(
		schema.TypeString,
		fmt.Sprintf("Specifies the storage serialization policy for Iceberg tables that use Snowflake as the catalog. Valid options are: %v. COMPATIBLE: Snowflake performs encoding and compression of data files that ensures interoperability with third-party compute engines. OPTIMIZED: Snowflake performs encoding and compression of data files that ensures the best table performance within Snowflake.", sdk.AsStringList(sdk.AllStorageSerializationPolicies)),
		func(inner *schema.Schema) {
			inner.ValidateDiagFunc = StringInSlice(sdk.AsStringList(sdk.AllStorageSerializationPolicies), true)
			inner.DiffSuppressFunc = func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue) || (d.Get(k).(string) == string(sdk.StorageSerializationPolicyOptimized) && newValue == "")
			}
		},
	),
	"log_level": nestedPropertyWithInnerModifier(
		schema.TypeString,
		fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
		func(inner *schema.Schema) {
			inner.ValidateDiagFunc = StringInSlice(sdk.AsStringList(sdk.AllLogLevels), true)
			inner.DiffSuppressFunc = func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue) || (d.Get(k).(string) == string(sdk.LogLevelOff) && newValue == "")
			}
		},
	),
	"trace_level": nestedPropertyWithInnerModifier(
		schema.TypeString,
		fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", sdk.AsStringList(sdk.AllTraceLevels)),
		func(inner *schema.Schema) {
			inner.ValidateDiagFunc = StringInSlice(sdk.AsStringList(sdk.AllTraceLevels), true)
			inner.DiffSuppressFunc = func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return strings.EqualFold(oldValue, newValue) || (d.Get(k).(string) == string(sdk.TraceLevelOff) && newValue == "")
			}
		},
	),
	"replication": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Configures replication for a given database. When specified, this database will be promoted to serve as a primary database for replication. A primary database can be replicated in one or more accounts, allowing users in those accounts to query objects in each secondary (i.e. replica) database.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable_for_account": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "Entry to enable replication and optionally failover for a given account identifier.",
					MinItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_identifier": {
								Type:     schema.TypeString,
								Required: true,
								// TODO(SNOW-1438810): Add account identifier validator
								Description: "Specifies account identifier for which replication should be enabled. The account identifiers should be in the form of `\"<organization_name>\".\"<account_name>\"`.",
							},
							"with_failover": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Specifies if failover should be enabled for the specified account identifier",
							},
						},
					},
				},
				"ignore_edition_check": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Allows replicating data to accounts on lower editions in either of the following scenarios: " +
						"1. The primary database is in a Business Critical (or higher) account but one or more of the accounts approved for replication are on lower editions. Business Critical Edition is intended for Snowflake accounts with extremely sensitive data. " +
						"2. The primary database is in a Business Critical (or higher) account and a signed business associate agreement is in place to store PHI data in the account per HIPAA and HITRUST regulations, but no such agreement is in place for one or more of the accounts approved for replication, regardless if they are Business Critical (or higher) accounts. " +
						"Both scenarios are prohibited by default in an effort to help prevent account administrators for Business Critical (or higher) accounts from inadvertently replicating sensitive data to accounts on lower editions.",
				},
			},
		},
	},
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

		CustomizeDiff: customdiff.All(
			NestedIntValueAccountObjectComputedIf("data_retention_time_in_days", sdk.AccountParameterDataRetentionTimeInDays),
			NestedIntValueAccountObjectComputedIf("max_data_extension_time_in_days", sdk.AccountParameterMaxDataExtensionTimeInDays),
			NestedStringValueAccountObjectComputedIf("external_volume", sdk.AccountParameterExternalVolume),
			NestedStringValueAccountObjectComputedIf("catalog", sdk.AccountParameterCatalog),
			NestedBoolValueAccountObjectComputedIf("replace_invalid_characters", sdk.AccountParameterReplaceInvalidCharacters),
			NestedStringValueAccountObjectComputedIf("default_ddl_collation", sdk.AccountParameterDefaultDDLCollation),
			NestedStringValueAccountObjectComputedIf("storage_serialization_policy", sdk.AccountParameterStorageSerializationPolicy),
			NestedStringValueAccountObjectComputedIf("log_level", sdk.AccountParameterLogLevel),
			NestedStringValueAccountObjectComputedIf("trace_level", sdk.AccountParameterTraceLevel),
		),

		Description: "Represents a standard database. If replication configuration is specified, the database is promoted to serve as a primary database for replication.",
		Schema:      standardDatabaseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	dataRetentionTimeInDays, _ := GetPropertyOfFirstNestedObjectByValueKey[int](d, "data_retention_time_in_days")
	maxDataExtensionTimeInDays, _ := GetPropertyOfFirstNestedObjectByValueKey[int](d, "max_data_extension_time_in_days")
	replaceInvalidCharacters, _ := GetPropertyOfFirstNestedObjectByValueKey[bool](d, "replace_invalid_characters")
	defaultDdlCollation, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "default_ddl_collation")

	var externalVolume *sdk.AccountObjectIdentifier
	if externalVolumeRaw, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "external_volume"); externalVolumeRaw != nil {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(*externalVolumeRaw))
	}

	var catalog *sdk.AccountObjectIdentifier
	if catalogRaw, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "catalog"); catalogRaw != nil {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(*catalogRaw))
	}

	var storageSerializationPolicy *sdk.StorageSerializationPolicy
	if storageSerializationPolicyRaw, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "storage_serialization_policy"); storageSerializationPolicyRaw != nil {
		storageSerializationPolicy = sdk.Pointer(sdk.StorageSerializationPolicy(*storageSerializationPolicyRaw))
	}

	var logLevel *sdk.LogLevel
	if logLevelRaw, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "log_level"); logLevelRaw != nil {
		logLevel = sdk.Pointer(sdk.LogLevel(*logLevelRaw))
	}

	var traceLevel *sdk.TraceLevel
	if traceLevelRaw, _ := GetPropertyOfFirstNestedObjectByValueKey[string](d, "trace_level"); traceLevelRaw != nil {
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

	var diags diag.Diagnostics

	if v, ok := d.GetOk("replication"); ok {
		replicationConfiguration := v.([]any)[0].(map[string]any)

		var ignoreEditionCheck *bool
		if v, ok := replicationConfiguration["ignore_edition_check"]; ok {
			ignoreEditionCheck = sdk.Pointer(v.(bool))
		}

		if enableForAccounts, ok := replicationConfiguration["enable_for_account"]; ok {
			enableForAccountList := enableForAccounts.([]any)

			if len(enableForAccountList) > 0 {
				replicationForAccounts := make([]sdk.AccountIdentifier, 0)
				failoverForAccounts := make([]sdk.AccountIdentifier, 0)

				for _, enableForAccount := range enableForAccountList {
					accountConfig := enableForAccount.(map[string]any)
					accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(accountConfig["account_identifier"].(string))

					replicationForAccounts = append(replicationForAccounts, accountIdentifier)
					if v, ok := accountConfig["with_failover"]; ok && v.(bool) {
						failoverForAccounts = append(failoverForAccounts, accountIdentifier)
					}
				}

				if len(replicationForAccounts) > 0 {
					err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
						EnableReplication: &sdk.EnableReplication{
							ToAccounts:         replicationForAccounts,
							IgnoreEditionCheck: ignoreEditionCheck,
						},
					})
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  err.Error(),
						})
					}
				}

				if len(failoverForAccounts) > 0 {
					err = client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
						EnableFailover: &sdk.EnableFailover{
							ToAccounts: failoverForAccounts,
						},
					})
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  err.Error(),
						})
					}
				}
			}
		}
	}

	return append(diags, ReadStandardDatabase(ctx, d, meta)...)
}

func UpdateStandardDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	var databaseSetRequest sdk.DatabaseSet
	var databaseUnsetRequest sdk.DatabaseUnset

	if d.HasChange("data_retention_time_in_days") {
		dataRetentionObject, ok := d.GetOk("data_retention_time_in_days")
		if ok && len(dataRetentionObject.([]any)) > 0 {
			dataRetentionTimeInDays, err := GetPropertyOfFirstNestedObjectByValueKey[int](d, "data_retention_time_in_days")
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
			maxDataExtensionTimeInDays, err := GetPropertyOfFirstNestedObjectByValueKey[int](d, "max_data_extension_time_in_days")
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
			externalVolume, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "external_volume")
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
			catalog, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "catalog")
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
			replaceInvalidCharacters, err := GetPropertyOfFirstNestedObjectByValueKey[bool](d, "replace_invalid_characters")
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
			defaultDdlCollation, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "default_ddl_collation")
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
			storageSerializationPolicy, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "storage_serialization_policy")
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
			logLevel, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "log_level")
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
			traceLevel, err := GetPropertyOfFirstNestedObjectByValueKey[string](d, "trace_level")
			if err != nil {
				return diag.FromErr(err)
			}
			databaseSetRequest.TraceLevel = sdk.Pointer(sdk.TraceLevel(*traceLevel))
		} else {
			databaseUnsetRequest.TraceLevel = sdk.Bool(true)
		}
	}

	if d.HasChange("replication") {
		before, after := d.GetChange("replication")

		var (
			accountsToEnableReplication  []sdk.AccountIdentifier
			accountsToDisableReplication []sdk.AccountIdentifier
			accountsToEnableFailover     []sdk.AccountIdentifier
			accountsToDisableFailover    []sdk.AccountIdentifier

			// maps represent replication configuration by having sdk.AccountIdentifier
			// as a key (implicitly enabling replication), and failover as an option value
			beforeReplicationFailoverConfigurationMap = make(map[sdk.AccountIdentifier]bool)
			afterReplicationFailoverConfigurationMap  = make(map[sdk.AccountIdentifier]bool)
		)

		fillReplicationMap := func(replicationConfigs []any, replicationFailoverMap map[sdk.AccountIdentifier]bool) {
			for _, replicationConfigurationMap := range replicationConfigs {
				replicationConfiguration := replicationConfigurationMap.(map[string]any)
				for _, enableForAccountMap := range replicationConfiguration["enable_for_account"].([]any) {
					enableForAccount := enableForAccountMap.(map[string]any)
					accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(enableForAccount["account_identifier"].(string))
					replicationFailoverMap[accountIdentifier] = enableForAccount["with_failover"].(bool)
				}
			}
		}
		fillReplicationMap(before.([]any), beforeReplicationFailoverConfigurationMap)
		fillReplicationMap(after.([]any), afterReplicationFailoverConfigurationMap)

		for accountIdentifier := range beforeReplicationFailoverConfigurationMap {
			if _, ok := afterReplicationFailoverConfigurationMap[accountIdentifier]; !ok {
				// Entry removed -> only replication needs to be disabled, because failover will be disabled implicitly
				// (in Snowflake you cannot have failover enabled when replication is disabled).
				accountsToDisableReplication = append(accountsToDisableReplication, accountIdentifier)
			}
		}

		for accountIdentifier, withFailover := range afterReplicationFailoverConfigurationMap {
			if beforeWithFailover, ok := beforeReplicationFailoverConfigurationMap[accountIdentifier]; !ok {
				// New entry, enable replication and failover if set to true
				accountsToEnableReplication = append(accountsToEnableReplication, accountIdentifier)
				if withFailover {
					accountsToEnableFailover = append(accountsToEnableFailover, accountIdentifier)
				}
				// Existing entry (check for possible failover modifications)
			} else if beforeWithFailover != withFailover {
				if withFailover {
					accountsToEnableFailover = append(accountsToEnableFailover, accountIdentifier)
				} else {
					accountsToDisableFailover = append(accountsToDisableFailover, accountIdentifier)
				}
			}
		}

		if len(accountsToEnableReplication) > 0 {
			err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
				EnableReplication: &sdk.EnableReplication{
					ToAccounts:         accountsToEnableReplication,
					IgnoreEditionCheck: sdk.Bool(d.Get("replication.0.ignore_edition_check").(bool)),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(accountsToEnableFailover) > 0 {
			err := client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
				EnableFailover: &sdk.EnableFailover{
					ToAccounts: accountsToEnableFailover,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(accountsToDisableReplication) > 0 {
			err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
				DisableReplication: &sdk.DisableReplication{
					ToAccounts: accountsToDisableReplication,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(accountsToDisableFailover) > 0 {
			err := client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
				DisableFailover: &sdk.DisableFailover{
					ToAccounts: accountsToDisableFailover,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
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
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			Set: &databaseSetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (databaseUnsetRequest != sdk.DatabaseUnset{}) {
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
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

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	currentAccountIdentifier := sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)
	replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, &sdk.ShowReplicationDatabasesOptions{
		WithPrimary: sdk.Pointer(sdk.NewExternalObjectIdentifier(currentAccountIdentifier, id)),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(replicationDatabases) == 1 {
		replicationAllowedToAccounts := make([]sdk.AccountIdentifier, 0)
		failoverAllowedToAccounts := make([]sdk.AccountIdentifier, 0)

		for _, allowedAccount := range strings.Split(replicationDatabases[0].ReplicationAllowedToAccounts, ",") {
			allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(strings.TrimSpace(allowedAccount))
			if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
				continue
			}
			replicationAllowedToAccounts = append(replicationAllowedToAccounts, allowedAccountIdentifier)
		}

		for _, allowedAccount := range strings.Split(replicationDatabases[0].FailoverAllowedToAccounts, ",") {
			allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(strings.TrimSpace(allowedAccount))
			if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
				continue
			}
			failoverAllowedToAccounts = append(failoverAllowedToAccounts, allowedAccountIdentifier)
		}

		enableForAccount := make([]map[string]any, 0)
		for _, allowedAccount := range replicationAllowedToAccounts {
			enableForAccount = append(enableForAccount, map[string]any{
				"account_identifier": allowedAccount.FullyQualifiedName(),
				"with_failover":      slices.Contains(failoverAllowedToAccounts, allowedAccount),
			})
		}

		var ignoreEditionCheck *bool
		if v, ok := d.GetOk("replication.0.ignore_edition_check"); ok {
			ignoreEditionCheck = sdk.Bool(v.(bool))
		}

		if len(enableForAccount) == 0 && ignoreEditionCheck == nil {
			err := d.Set("replication", []any{})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := d.Set("replication", []any{
				map[string]any{
					"enable_for_account":   enableForAccount,
					"ignore_edition_check": ignoreEditionCheck,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	for _, parameter := range parameters {
		switch parameter.Key {
		case "MAX_DATA_EXTENSION_TIME_IN_DAYS":
			maxDataExtensionTimeInDays, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "max_data_extension_time_in_days", maxDataExtensionTimeInDays); err != nil {
				return diag.FromErr(err)
			}
		case "EXTERNAL_VOLUME":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "external_volume", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "CATALOG":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "catalog", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "DEFAULT_DDL_COLLATION":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "default_ddl_collation", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "LOG_LEVEL":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "log_level", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "TRACE_LEVEL":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "trace_level", parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case "REPLACE_INVALID_CHARACTERS":
			boolValue, err := strconv.ParseBool(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "replace_invalid_characters", boolValue); err != nil {
				return diag.FromErr(err)
			}
		case "STORAGE_SERIALIZATION_POLICY":
			if err := SetPropertyOfFirstNestedObjectByValueKey(d, "storage_serialization_policy", parameter.Value); err != nil {
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
