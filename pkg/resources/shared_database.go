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
	// TODO: Should it be imported (set in Read)?
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

func SharedDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateSharedDatabase,
		UpdateContext: UpdateSharedDatabase,
		ReadContext:   ReadSharedDatabase,
		DeleteContext: DeleteSharedDatabase,

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

	err := client.Databases.CreateShared(ctx, id, externalShareId, &sdk.CreateSharedDatabaseOptions{
		Comment: nil,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	/*
		"name"
		"from_share"
		"is_transient"
		"external_volume"
		"catalog"
		"default_ddl_collation"
		"log_level"
		"trace_level"
		"comment"
	*/

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadSharedDatabase(ctx, d, meta)
}

func UpdateSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	if d.HasChange("name") {
		newName := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			NewName: newName,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeSnowflakeID(newName))
		id = newName
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

	/*
		"name"
		"from_share"
		"is_transient"
		"external_volume"
		"catalog"
		"default_ddl_collation"
		"log_level"
		"trace_level"
		"comment"
	*/

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
