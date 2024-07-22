package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Role() *schema.Resource {
	accountRole := AccountRole()
	accountRole.DeprecationMessage = "This resource is deprecated and will be removed in a future major version release. Please use snowflake_account_role instead."
	return accountRole
}
