package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UnsafeExecute() *schema.Resource {
	unsafeExecute := Execute()
	unsafeExecute.Description = "Experimental resource allowing execution of ANY SQL statement. It may destroy resources if used incorrectly. It may behave incorrectly combined with other resources. Use at your own risk."
	unsafeExecute.DeprecationMessage = "This resource is deprecated and will be removed in a future major version release. Please use snowflake_execute instead."
	return unsafeExecute
}
