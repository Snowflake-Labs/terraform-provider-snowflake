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
	"data_retention_time_in_days": {
		Type:        schema.TypeInt,
		Computed:    true,
		Optional:    true,
		Description: "Specifies the number of days for which Time Travel actions (CLONE and UNDROP) can be performed on the database, as well as specifying the default Time Travel retention time for all schemas created in the database. For more details, see [Understanding & Using Time Travel](https://docs.snowflake.com/en/user-guide/data-time-travel).",
	},
	"max_data_extension_time_in_days": {
		Type:        schema.TypeInt,
		Computed:    true,
		Optional:    true,
		Description: "Object parameter that specifies the maximum number of days for which Snowflake can extend the data retention period for tables in the database to prevent streams on the tables from becoming stale. For a detailed description of this parameter, see [MAX_DATA_EXTENSION_TIME_IN_DAYS](https://docs.snowflake.com/en/sql-reference/parameters.html#label-max-data-extension-time-in-days).",
	},
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
	"default_ddl_collation": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
	},
	"log_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: StringInSlice([]string{}, true),                                                                                                                                                                                                                                                                                                                                // TODO: enum
		Description:      fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", []string{}), // TODO:
	},
	"trace_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: StringInSlice([]string{}, true),                                                                                                                                                                                                                // TODO: enum
		Description:      fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", []string{}), // TODO:
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

	var dataRetentionTimeInDays *int
	if v, ok := d.GetOk("data_retention_time_in_days"); ok {
		dataRetentionTimeInDays = sdk.Int(v.(int))
	}

	err := client.Databases.CreateSecondary(ctx, secondaryDatabaseId, primaryDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
		DataRetentionTimeInDays: dataRetentionTimeInDays,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	/*
		"name"
		"as_replica_of"
		"is_transient"
		"data_retention_time_in_days"
		"max_data_extension_time_in_days"
		"external_volume"
		"catalog"
		"default_ddl_collation"
		"log_level"
		"trace_level"
		"comment"
	*/

	d.SetId(helpers.EncodeSnowflakeID(secondaryDatabaseId))

	return ReadSecondaryDatabase(ctx, d, meta)
}

func UpdateSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	secondaryDatabaseId := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newName := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			NewName: newName,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newName))
		secondaryDatabaseId = newName
	}

	runSet := false
	databaseSetRequest := sdk.DatabaseSet{}

	runUnset := false
	databaseUnsetRequest := sdk.DatabaseUnset{}

	if d.HasChange("data_retention_time_in_days") {
	}

	if d.HasChange("max_data_extension_time_in_days") {
	}

	if d.HasChange("external_volume") {
	}

	if d.HasChange("catalog") {
	}

	if d.HasChange("default_ddl_collation") {
	}

	if d.HasChange("log_level") {
	}

	if d.HasChange("trace_level") {
	}

	if d.HasChange("comment") {
	}

	if runSet {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Set: &databaseSetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if runUnset {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Unset: &databaseUnsetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	/*
		"name"
		"as_replica_of"
		"is_transient"
		"data_retention_time_in_days"
		"max_data_extension_time_in_days"
		"external_volume"
		"catalog"
		"default_ddl_collation"
		"log_level"
		"trace_level"
		"comment"
	*/

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

	if err := d.Set("is_transient", secondaryDatabase.Transient); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("data_retention_time_in_days", secondaryDatabase.RetentionTime); err != nil {
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
			if err := d.Set("max_data_extension_time_in_days", maxDataExtensionTimeInDays); err != nil {
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
