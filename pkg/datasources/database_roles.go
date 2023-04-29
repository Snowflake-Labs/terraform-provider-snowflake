package datasources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	d.SetId("database_roles_read")
	databaseName := d.Get("database").(string)

	listRoles, err := snowflake.ListDatabaseRoles(databaseName, db)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[DEBUG] no roles found in database (%s)", databaseName)
		d.SetId("")
		return nil
	} else if err != nil {
		log.Println("[DEBUG] failed to list roles")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] list roles: %v", listRoles)

	roles := []map[string]interface{}{}
	for _, role := range listRoles {
		roleMap := map[string]interface{}{}

		roleMap["name"] = role.Name
		roleMap["comment"] = role.Comment
		roleMap["owner"] = role.Owner
		roles = append(roles, roleMap)
	}

	if err := d.Set("database_roles", roles); err != nil {
		return err
	}
	return nil
}
