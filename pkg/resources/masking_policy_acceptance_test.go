package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_MaskingPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"
	comment2 := "Terraform acceptance test 2"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyConfig(accName, accName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.#", "1"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.name", "val"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.type", "VARCHAR"),
				),
			},
			// change comment
			{
				Config: maskingPolicyConfig(accName, accName, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", comment2),
				),
			},
			// rename
			{
				Config: maskingPolicyConfig(accName, accName2, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName2),
				),
			},
			// change body and unset comment
			{
				Config: maskingPolicyConfigMultiline(accName, accName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case \n\twhen current_role() in ('ROLE_A') then \n\t\tval \n\twhen is_role_in_session( 'ROLE_B' ) then \n\t\t'ABC123'\n\telse\n\t\t'******'\nend"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", ""),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_masking_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func maskingPolicyConfig(n string, name string, comment string) string {
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

resource "snowflake_masking_policy" "test" {
	name = "%s"
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
	comment = "%s"
}
`, n, n, name, comment)
}

func maskingPolicyConfigMultiline(n string, name string) string {
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
	
	resource "snowflake_masking_policy" "test" {
		name = "%s"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		signature {
			column {
				name = "val"
				type = "VARCHAR"
			}
		}
		masking_expression = <<-EOF
			case 
				when current_role() in ('ROLE_A') then 
					val 
				when is_role_in_session( 'ROLE_B' ) then 
					'ABC123'
				else
					'******'
			end
    	EOF
		return_data_type = "VARCHAR"
	}
	`, n, n, name)
}
