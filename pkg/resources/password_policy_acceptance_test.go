package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: passwordPolicyConfig(accName, 10, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
				),
			},
			{
				Config: passwordPolicyConfig(accName, 20, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
				),
			},
			{
				Config: passwordPolicyConfig(accName, 20, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", ""),
				),
			},
			{
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func passwordPolicyConfig(s string, minLength int, comment string) string {
	return fmt.Sprintf(`
	resource "snowflake_database" "test" {
		name = "%v"
		comment = "Terraform acceptance test"
	  }
	  
	  resource "snowflake_schema" "test" {
		name = "%v"
		database = snowflake_database.test.name
		comment = "Terraform acceptance test"
	  }
	  
	resource "snowflake_password_policy" "pa" {
		database   = snowflake_database.test.name
		schema     = snowflake_schema.test.name
		name       = "%v"
		min_length = %d
		comment    = "%s"
		or_replace = true
	}
	`, s, s, s, minLength, comment)
}
