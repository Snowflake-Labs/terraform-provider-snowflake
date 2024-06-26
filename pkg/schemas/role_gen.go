// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowRoleSchema represents output of SHOW query for the single Role.
var ShowRoleSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_inherited": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"assigned_to_users": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"granted_to_roles": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"granted_roles": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowRoleSchema

func RoleToSchema(role *sdk.Role) map[string]any {
	roleSchema := make(map[string]any)
	roleSchema["created_on"] = role.CreatedOn.String()
	roleSchema["name"] = role.Name
	roleSchema["is_default"] = role.IsDefault
	roleSchema["is_current"] = role.IsCurrent
	roleSchema["is_inherited"] = role.IsInherited
	roleSchema["assigned_to_users"] = role.AssignedToUsers
	roleSchema["granted_to_roles"] = role.GrantedToRoles
	roleSchema["granted_roles"] = role.GrantedRoles
	roleSchema["owner"] = role.Owner
	roleSchema["comment"] = role.Comment
	return roleSchema
}

var _ = RoleToSchema
