// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowPasswordPolicySchema represents output of SHOW query for the single PasswordPolicy.
var ShowPasswordPolicySchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
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
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowPasswordPolicySchema

func PasswordPolicyToSchema(passwordPolicy *sdk.PasswordPolicy) map[string]any {
	passwordPolicySchema := make(map[string]any)
	passwordPolicySchema["created_on"] = passwordPolicy.CreatedOn.String()
	passwordPolicySchema["name"] = passwordPolicy.Name
	passwordPolicySchema["database_name"] = passwordPolicy.DatabaseName
	passwordPolicySchema["schema_name"] = passwordPolicy.SchemaName
	passwordPolicySchema["kind"] = passwordPolicy.Kind
	passwordPolicySchema["owner"] = passwordPolicy.Owner
	passwordPolicySchema["comment"] = passwordPolicy.Comment
	passwordPolicySchema["owner_role_type"] = passwordPolicy.OwnerRoleType
	return passwordPolicySchema
}

var _ = PasswordPolicyToSchema