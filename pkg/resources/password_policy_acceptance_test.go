package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "snowflake_password_policy" "pa" {
					database   = "TEST_DB"
					schema     = "PUBLIC"
					name       = "mypolicy"
					min_length = 10
					comment    = "this is a test resource"
					or_replace = true
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", "mypolicy"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
				),
			},
			{
				Config: `
				resource "snowflake_password_policy" "pa" {
					database   = "TEST_DB"
					schema     = "PUBLIC"
					name       = "mypolicy"
					min_length = 20
					comment    = "this is a test resource"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
				),
			},
			{
				Config: `
				resource "snowflake_password_policy" "pa" {
					database   = "TEST_DB"
					schema     = "PUBLIC"
					name       = "mypolicy"
					min_length = 20
				}
				`,
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
