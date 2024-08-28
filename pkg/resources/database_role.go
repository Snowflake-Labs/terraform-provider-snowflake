package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database role."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the database role."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database role.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW DATABASE ROLES` for the given database role. Note that this value will be only recomputed whenever comment field changes.",
		Elem: &schema.Resource{
			Schema: schemas.ShowDatabaseRoleSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func DatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateDatabaseRole,
		ReadContext:   ReadDatabaseRole,
		UpdateContext: UpdateDatabaseRole,
		DeleteContext: DeleteDatabaseRole,

		Schema: databaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportName[sdk.DatabaseObjectIdentifier],
		},

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChangedWithSuppressDiff(ShowOutputAttributeName, suppressIdentifierQuoting, "name"),
		),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},
	}
}

func ReadDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Databases.ShowByID(ctx, id.DatabaseId())
	if err != nil {
		log.Printf("[DEBUG] database %s for database role %s not found", id.DatabaseId().Name(), id.Name())
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to query database. Marking the resource as removed.",
				Detail:   fmt.Sprintf("database role name: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Database role not found; marking it as removed",
					Detail:   fmt.Sprintf("Database role name: %s, err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set("comment", databaseRole.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.DatabaseRoleToSchema(databaseRole)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func CreateDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	roleName := d.Get("name").(string)
	id := sdk.NewDatabaseObjectIdentifier(databaseName, roleName)
	createRequest := sdk.NewCreateDatabaseRoleRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	err := client.DatabaseRoles.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadDatabaseRole(ctx, d, meta)
}

func UpdateDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), d.Get("name").(string))

		err = client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("comment") {
		newComment := d.Get("comment").(string)
		err := client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithSet(*sdk.NewDatabaseRoleSetRequest(newComment)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadDatabaseRole(ctx, d, meta)
}

func DeleteDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
