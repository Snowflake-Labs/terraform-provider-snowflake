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

var sharedDatabaseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the database; must be unique for your account.",
	},
	"from_share": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "A fully qualified path to a share from which the database will be created. A fully qualified path follows the format of `\"<organization_name>\".\"<account_name>\".\"<share_name>\"`.",
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the database as transient. Transient databases do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	"external_volume": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Description:      "The database parameter that specifies the default external volume to use for Iceberg tables.",
	},
	"catalog": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Description:      "The database parameter that specifies the default catalog to use for Iceberg tables.",
	},
	"default_ddl_collation": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a default collation specification for all schemas and tables added to the database. It can be overridden on schema or table level. For more information, see [collation specification](https://docs.snowflake.com/en/sql-reference/collation#label-collation-specification).",
	},
	"log_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllLogLevels), true),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(string) == "OFF" && newValue == ""
		},
		Description: fmt.Sprintf("Specifies the severity level of messages that should be ingested and made available in the active event table. Valid options are: %v. Messages at the specified level (and at more severe levels) are ingested. For more information, see [LOG_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-log-level).", sdk.AsStringList(sdk.AllLogLevels)),
	},
	"trace_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: StringInSlice(sdk.AsStringList(sdk.AllTraceLevels), true),
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			return d.Get(k).(string) == "OFF" && newValue == ""
		},
		Description: fmt.Sprintf("Controls how trace events are ingested into the event table. Valid options are: %v. For information about levels, see [TRACE_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters.html#label-trace-level).", sdk.AsStringList(sdk.AllTraceLevels)),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
}

func SharedDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateSharedDatabase,
		UpdateContext: UpdateSharedDatabase,
		ReadContext:   ReadSharedDatabase,
		DeleteContext: DeleteSharedDatabase,
		Description:   "A shared database creates a database from a share provided by another Snowflake account. For more information about shares, see [Introduction to Secure Data Sharing](https://docs.snowflake.com/en/user-guide/data-sharing-intro).",

		Schema: sharedDatabaseSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
	externalShareId := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(d.Get("from_share").(string))

	var externalVolume *sdk.AccountObjectIdentifier
	if v, ok := d.GetOk("external_volume"); ok {
		externalVolume = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	var catalog *sdk.AccountObjectIdentifier
	if v, ok := d.GetOk("catalog"); ok {
		catalog = sdk.Pointer(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	var logLevel *sdk.LogLevel
	if v, ok := d.GetOk("log_level"); ok {
		logLevel = sdk.Pointer(sdk.LogLevel(v.(string)))
	}

	var traceLevel *sdk.TraceLevel
	if v, ok := d.GetOk("trace_level"); ok {
		traceLevel = sdk.Pointer(sdk.TraceLevel(v.(string)))
	}

	err := client.Databases.CreateShared(ctx, id, externalShareId, &sdk.CreateSharedDatabaseOptions{
		Transient:           GetPropertyAsPointer[bool](d, "is_transient"),
		ExternalVolume:      externalVolume,
		Catalog:             catalog,
		DefaultDDLCollation: GetPropertyAsPointer[string](d, "default_ddl_collation"),
		LogLevel:            logLevel,
		TraceLevel:          traceLevel,
		Comment:             GetPropertyAsPointer[string](d, "comment"),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadSharedDatabase(ctx, d, meta)
}

func UpdateSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newName := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			NewName: &newName,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newName))
		id = newName
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					Comment: &comment,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					Comment: sdk.Bool(true),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadSharedDatabase(ctx, d, meta)
}

func ReadSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	database, err := client.Databases.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query shared database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("DatabaseName: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	secondaryDatabaseParameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{
			Database: id,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", database.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("from_share", sdk.NewExternalObjectIdentifierFromFullyQualifiedName(database.Origin).FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_transient", database.Transient); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", database.Comment); err != nil {
		return diag.FromErr(err)
	}

	for _, secondaryDatabaseParameter := range secondaryDatabaseParameters {
		switch secondaryDatabaseParameter.Key {
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

func DeleteSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
