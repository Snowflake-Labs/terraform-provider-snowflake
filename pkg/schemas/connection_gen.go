// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowConnectionSchema represents output of SHOW query for the single Connection.
var ShowConnectionSchema = map[string]*schema.Schema{
	"region_group": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"snowflake_region": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_primary": {
		Type:     schema.TypeBool,
		Computed: true,
		ForceNew: true,
	},
	"primary": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"failover_allowed_to_accounts": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"connection_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_locator": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowConnectionSchema

func ConnectionToSchema(connection *sdk.Connection) map[string]any {
	connectionSchema := make(map[string]any)
	if connection.RegionGroup != nil {
		connectionSchema["region_group"] = connection.RegionGroup
	}
	connectionSchema["snowflake_region"] = connection.SnowflakeRegion
	connectionSchema["created_on"] = connection.CreatedOn.String()
	connectionSchema["account_name"] = connection.AccountName
	connectionSchema["name"] = connection.Name
	if connection.Comment != nil {
		connectionSchema["comment"] = connection.Comment
	}
	connectionSchema["is_primary"] = connection.IsPrimary
	connectionSchema["primary"] = connection.Primary.FullyQualifiedName()
	var allowedAccounts []string
	for _, accountId := range connection.FailoverAllowedToAccounts {
		allowedAccounts = append(allowedAccounts, accountId.Name())
	}
	connectionSchema["failover_allowed_to_accounts"] = allowedAccounts
	connectionSchema["connection_url"] = connection.ConnectionUrl
	connectionSchema["organization_name"] = connection.OrganizationName
	connectionSchema["account_locator"] = connection.AccountLocator
	return connectionSchema
}

var _ = ConnectionToSchema
