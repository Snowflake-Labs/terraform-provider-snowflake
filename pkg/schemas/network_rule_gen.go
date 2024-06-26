// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowNetworkRuleSchema represents output of SHOW query for the single NetworkRule.
var ShowNetworkRuleSchema = map[string]*schema.Schema{
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
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"mode": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"entries_in_value_list": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowNetworkRuleSchema

func NetworkRuleToSchema(networkRule *sdk.NetworkRule) map[string]any {
	networkRuleSchema := make(map[string]any)
	networkRuleSchema["created_on"] = networkRule.CreatedOn.String()
	networkRuleSchema["name"] = networkRule.Name
	networkRuleSchema["database_name"] = networkRule.DatabaseName
	networkRuleSchema["schema_name"] = networkRule.SchemaName
	networkRuleSchema["owner"] = networkRule.Owner
	networkRuleSchema["comment"] = networkRule.Comment
	networkRuleSchema["type"] = string(networkRule.Type)
	networkRuleSchema["mode"] = string(networkRule.Mode)
	networkRuleSchema["entries_in_value_list"] = networkRule.EntriesInValueList
	networkRuleSchema["owner_role_type"] = networkRule.OwnerRoleType
	return networkRuleSchema
}

var _ = NetworkRuleToSchema
