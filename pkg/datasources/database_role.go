package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseRoleSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the database role from.",
	},
	"role": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Database role name.",
	},
	"comment": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The comment on the role",
	},
	"owner": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The owner of the role",
	},
}

// DatabaseRole Snowflake Database Role resource.
func DatabaseRole() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabaseRole,
		Schema: databaseRoleSchema,
	}
}

// ReadDatabaseRole Reads the database role metadata information.
func ReadDatabaseRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	d.SetId("database_role_read")

	databaseName := d.Get("database").(string)
	roleName := d.Get("role").(string)

	ctx := context.Background()
	showRequest := sdk.NewShowDatabaseRoleRequest(sdk.NewAccountObjectIdentifier(databaseName)).WithLike(roleName)
	extractedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
	if err != nil {
		log.Printf("[DEBUG] unable to show database role %s in db (%s)", roleName, databaseName)
		d.SetId("")
		return err
	}

	for _, databaseRole := range extractedDatabaseRoles {
		err = d.Set("comment", databaseRole.Comment)
		if err != nil {
			return err
		}
		err = d.Set("owner", databaseRole.Owner)
		if err != nil {
			return err
		}
		break
	}

	return err
}
