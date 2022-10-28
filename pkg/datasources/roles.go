package datasources

import (
	"database/sql"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

var rolesSchema = map[string]*schema.Schema{
	"roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "All roles in the account, including built in roles",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the role.",
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

// Roles Snowflake Roles resource.
func Roles() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRoles,
		Schema: rolesSchema,
	}
}

// ReadRoles Reads the database metadata information.
func ReadRoles(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	d.SetId("roles_read")

	listRoles, err := snowflake.ListRoles(db)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no roles found in account (%s)", d.Id())
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
		if !role.Name.Valid {
			continue
		}
		roleMap["name"] = role.Name.String
		roleMap["comment"] = role.Comment.String
		roleMap["owner"] = role.Owner.String
		roles = append(roles, roleMap)
	}

	err = d.Set("roles", roles)
	if err != nil {
		return err
	}
	return nil
}
