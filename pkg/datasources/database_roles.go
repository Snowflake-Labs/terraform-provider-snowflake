package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseRolesSchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The database from which to return the database roles from.",
	},
	"database_roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Lists all the database roles in a specified database.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Identifier for the role.",
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
			},
		},
	},
}

// DatabaseRoles Snowflake Database Roles resource.
func DatabaseRoles() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabaseRoles,
		Schema: databaseRolesSchema,
	}
}

// ReadDatabaseRoles Reads the database metadata information.
func ReadDatabaseRoles(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	d.SetId("database_roles_read")

	databaseName := d.Get("database").(string)

	ctx := context.Background()
	showRequest := sdk.NewShowDatabaseRoleRequest(sdk.NewAccountObjectIdentifier(databaseName))
	extractedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
	if err != nil {
		log.Printf("[DEBUG] unable to show database roles in db (%s)", databaseName)
		d.SetId("")
		return err
	}

	databaseRoles := make([]map[string]any, 0, len(extractedDatabaseRoles))
	for _, databaseRole := range extractedDatabaseRoles {
		databaseRoleMap := map[string]any{}

		databaseRoleMap["name"] = databaseRole.Name
		databaseRoleMap["comment"] = databaseRole.Comment
		databaseRoleMap["owner"] = databaseRole.Owner

		databaseRoles = append(databaseRoles, databaseRoleMap)
	}

	return d.Set("database_roles", databaseRoles)
}
