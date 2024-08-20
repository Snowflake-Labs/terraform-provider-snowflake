package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the database role.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The database in which to create the database role.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database role.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// DatabaseRole returns a pointer to the resource representing a database role.
func DatabaseRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateDatabaseRole,
		Read:   ReadDatabaseRole,
		Update: UpdateDatabaseRole,
		Delete: DeleteDatabaseRole,

		Schema: databaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

// ReadDatabaseRole implements schema.ReadFunc.
func ReadDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] database role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}

	if err := d.Set("name", databaseRole.Name); err != nil {
		return err
	}

	if err := d.Set("database", id.DatabaseName()); err != nil {
		return err
	}

	if err := d.Set("comment", databaseRole.Comment); err != nil {
		return err
	}
	return nil
}

// CreateDatabaseRole implements schema.CreateFunc.
func CreateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	roleName := d.Get("name").(string)

	objectIdentifier := sdk.NewDatabaseObjectIdentifier(databaseName, roleName)
	createRequest := sdk.NewCreateDatabaseRoleRequest(objectIdentifier)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	ctx := context.Background()
	err := client.DatabaseRoles.Create(ctx, createRequest)
	if err != nil {
		return err
	}

	d.SetId(objectIdentifier.FullyQualifiedName())

	return ReadDatabaseRole(d, meta)
}

// UpdateDatabaseRole implements schema.UpdateFunc.
func UpdateDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	objectIdentifier, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("comment") {
		_, newVal := d.GetChange("comment")

		ctx := context.Background()
		alterRequest := sdk.NewAlterDatabaseRoleRequest(objectIdentifier).WithSetComment(newVal.(string))
		err := client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			return fmt.Errorf("error updating database role %v: %w", objectIdentifier.Name(), err)
		}
	}

	return ReadDatabaseRole(d, meta)
}

// DeleteDatabaseRole implements schema.DeleteFunc.
func DeleteDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	objectIdentifier, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	dropRequest := sdk.NewDropDatabaseRoleRequest(objectIdentifier)
	err = client.DatabaseRoles.Drop(ctx, dropRequest)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
