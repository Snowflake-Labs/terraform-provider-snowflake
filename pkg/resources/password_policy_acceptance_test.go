package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: passwordPolicyConfig(accName, 10, 30, "this is a test resource", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "30"),
				),
			},
			{
				Config: passwordPolicyConfig(accName, 20, 50, "this is a test resource", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "50"),
				),
			},
			{
				Config: passwordPolicyConfig(accName, 20, 50, "", acc.TestDatabaseName, acc.TestSchemaName),
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

func passwordPolicyConfig(s string, minLength int, maxLength int, comment string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_password_policy" "pa" {
		name       = "%v"
		database   = "%s"
		schema     = "%s"
		min_length = %d
		max_length = %d
		or_replace = true
	}
	`, s, databaseName, schemaName, minLength, maxLength)
}

func TestAcc_PasswordPolicyMaxAgeDays(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Creation sets zero properly
			{
				Config: passwordPolicyDefaultMaxAgeDaysConfig(accName, acc.TestDatabaseName, acc.TestSchemaName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			{
				Config: passwordPolicyDefaultMaxAgeDaysConfig(accName, acc.TestDatabaseName, acc.TestSchemaName, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "10"),
				),
			},
			// Update sets zero properly
			{
				Config: passwordPolicyDefaultMaxAgeDaysConfig(accName, acc.TestDatabaseName, acc.TestSchemaName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			// Unsets properly
			{
				Config: passwordPolicyDefaultConfigWithoutMaxAgeDays(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "90"),
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

func passwordPolicyDefaultMaxAgeDaysConfig(s string, databaseName string, schemaName string, maxAgeDays int) string {
	return fmt.Sprintf(`
	resource "snowflake_password_policy" "pa" {
		name         = "%v"
		database   = "%s"
		schema     = "%s"
		max_age_days = %d
	}
	`, s, databaseName, schemaName, maxAgeDays)
}

func passwordPolicyDefaultConfigWithoutMaxAgeDays(s string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_password_policy" "pa" {
		name         = "%v"
		database   = "%s"
		schema     = "%s"
	}
	`, s, databaseName, schemaName)
}
