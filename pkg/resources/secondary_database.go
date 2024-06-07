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

		CustomizeDiff: DatabaseParametersCustomDiff,
		Schema:        MergeMaps(secondaryDatabaseSchema, DatabaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateSecondaryDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	secondaryDatabaseId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
	primaryDatabaseId := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(d.Get("as_replica_of").(string))

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
		enableConsoleOutput := GetAllDatabaseParameters(d)

	err := client.Databases.CreateSecondary(ctx, secondaryDatabaseId, primaryDatabaseId, &sdk.CreateSecondaryDatabaseOptions{
		Transient:                               GetPropertyAsPointer[bool](d, "is_transient"),
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
		Comment:                                 GetPropertyAsPointer[string](d, "comment"),
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

	databaseSetRequest := new(sdk.DatabaseSet)
	databaseUnsetRequest := new(sdk.DatabaseUnset)

	if updateParamDiags := HandleDatabaseParameterChanges(d, databaseSetRequest, databaseUnsetRequest); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			databaseSetRequest.Comment = &comment
		} else {
			databaseUnsetRequest.Comment = sdk.Bool(true)
		}
	}

	if (*databaseSetRequest != sdk.DatabaseSet{}) {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Set: databaseSetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*databaseUnsetRequest != sdk.DatabaseUnset{}) {
		err := client.Databases.Alter(ctx, secondaryDatabaseId, &sdk.AlterDatabaseOptions{
			Unset: databaseUnsetRequest,
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

	if err := d.Set("comment", secondaryDatabase.Comment); err != nil {
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

	if diags := HandleDatabaseParameterRead(d, secondaryDatabaseParameters); diags != nil {
		return diags
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
