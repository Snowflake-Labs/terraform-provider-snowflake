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
				Config: passwordPolicyConfig(accName, 10, 30, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "30"),
				),
			},
			{
				Config: passwordPolicyConfig(accName, 20, 50, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "50"),
				),
			},
			/*
					todo: fix once comments are working again for password policies
					query CREATE PASSWORD POLICY IF NOT EXISTS "T_Kn1bY6?2kx"."}k*3DrsXP:w9TRK#4wtS"."9ec016f6-ce74-0c94-2bd5-dc46547dbeff" PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 5 COMMENT = 'test comment' err 001420 (22023): SQL compilation error: invalid property 'COMMENT' for 'PASSWORD_POLICY'
				{
					Config: passwordPolicyConfig(accName, 20, 50, ""),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", ""),
					),
				},
			*/
			{
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func passwordPolicyConfig(s string, minLength int, maxLength int, comment string) string {
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
		max_length = %d
		or_replace = true
	}
	`, s, s, s, minLength, maxLength)
}

func TestAcc_PasswordPolicyMaxAgeDays(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Creation sets zero properly
			{
				Config: passwordPolicyDefaultMaxageDaysConfig(accName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			{
				Config: passwordPolicyDefaultMaxageDaysConfig(accName, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "10"),
				),
			},
			// Update sets zero properly
			{
				Config: passwordPolicyDefaultMaxageDaysConfig(accName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
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

func passwordPolicyDefaultMaxageDaysConfig(s string, maxAgeDays int) string {
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
		database     = snowflake_database.test.name
		schema       = snowflake_schema.test.name
		name         = "%v"
		max_age_days = %d
	}
	`, s, s, s, maxAgeDays)
}
