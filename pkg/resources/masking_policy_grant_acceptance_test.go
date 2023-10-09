package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_MaskingPolicyGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyGrantConfig(accName, "APPLY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "privilege", "APPLY"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: maskingPolicyGrantConfig(accName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_masking_policy_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func maskingPolicyGrantConfig(name string, privilege string) string {
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

	resource "snowflake_role" "test" {
		name = "%v"
	}

	resource "snowflake_masking_policy" "test" {
		name = "%v"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		signature {
			column {
				name = "val"
				type = "VARCHAR"
			}
		}
		masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
		return_data_type = "VARCHAR"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_masking_policy_grant" "test" {
		masking_policy_name = snowflake_masking_policy.test.name
		database_name = snowflake_database.test.name
		roles         = [snowflake_role.test.name]
		schema_name   = snowflake_schema.test.name
		privilege = "%s"
	}
	`, name, name, name, name, privilege)
}
