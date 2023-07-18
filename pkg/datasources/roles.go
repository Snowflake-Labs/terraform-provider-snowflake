package datasources

import (
	"context"
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RolesPatternKey = "pattern"
	RolesRolesKey   = "roles"
)

var rolesSchema = map[string]*schema.Schema{
	RolesPatternKey: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	RolesRolesKey: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the roles which you can view across your entire account, including the system-defined roles and any custom roles that exist.",
		Elem:        roleSchema,
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
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	d.SetId("roles_read")

	pattern := d.Get(RolesPatternKey).(string)
	roleList, err := client.Roles.Show(ctx, &sdk.RoleShowOptions{
		Like: &sdk.Like{
			Pattern: &pattern,
		},
	})
	if err != nil {
		log.Println("[DEBUG] failed to list roles")
		d.SetId("")
		return nil
	}
	if len(roleList) != 0 {
		log.Printf("[DEBUG] no roles found in account (%s)", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] found roles: %v", roleList)

	roles := []map[string]interface{}{}
	for _, role := range roleList {
		if role != nil {
			roleMap := map[string]interface{}{}
			roleMap[RoleNameKey] = role.Name
			roleMap[RoleCommentKey] = role.Comment
			roles = append(roles, roleMap)
		}
	}

	return d.Set(RolesRolesKey, roles)
}
