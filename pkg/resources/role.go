package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Role() *schema.Resource {
	return &schema.Resource{
		CreateContext:      CreateAccountRole,
		ReadContext:        ReadAccountRole,
		DeleteContext:      DeleteAccountRole,
		UpdateContext:      UpdateAccountRole,
		DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_account_role instead.",

		Schema: accountRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
