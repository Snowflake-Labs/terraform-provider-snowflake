package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: passwordPolicyConfig(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", "mypolicy"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
				),
			},
			{
				Config: passwordPolicyConfig(20),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
				),
			},
		},
	})
}

func passwordPolicyConfig(minLength int) string {
	s := `
resource "snowflake_password_policy" "pa" {
	database   = "TEST_DB"
	schema     = "PUBLIC"
	name       = "mypolicy"
	min_length = %v
	comment    = "this is a test resource"
}
`
	return fmt.Sprintf(s, minLength)
}
