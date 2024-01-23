package datasources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Roles() *schema.Resource {
	return &schema.Resource{
		ReadContext:        ReadAccountRoles,
		Schema:             accountRolesSchema,
		DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_account_roles instead.",
	}
}
