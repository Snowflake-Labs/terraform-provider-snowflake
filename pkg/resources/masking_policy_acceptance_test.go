package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_MaskingPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"
	comment2 := "Terraform acceptance test 2"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyConfig(accName, accName, comment, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", comment),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.#", "1"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.name", "val"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.type", "VARCHAR"),
				),
			},
			// rename + change comment
			{
				Config: maskingPolicyConfig(accName, accName2, comment2, acc.TestDatabaseName, acc.TestSchemaName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_masking_policy.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName2),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", comment2),
				),
			},
			// change body and unset comment
			{
				Config: maskingPolicyConfigMultiline(accName, accName2, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case\n\twhen current_role() in ('ROLE_A') then\n\t\tval\n\twhen is_role_in_session( 'ROLE_B' ) then\n\t\t'ABC123'\n\telse\n\t\t'******'\nend"),
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

func maskingPolicyConfig(n string, name string, comment string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	name = "%s"
	database = "%s"
	schema = "%s"
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
`, name, databaseName, schemaName, comment)
}

func maskingPolicyConfigMultiline(n string, name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_masking_policy" "test" {
		name = "%s"
		database = "%s"
		schema = "%s"
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
	`, name, databaseName, schemaName)
}

func TestAcc_MaskingPolicyMultiColumns(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyConfigMultiColumn(accName, accName, acc.TestDatabaseName, acc.TestSchemaName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.#", "1"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.name", "val"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.0.type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.1.name", "val2"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "signature.0.column.1.type", "VARCHAR"),
				),
			},
		},
	})
}

func maskingPolicyConfigMultiColumn(n string, name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	name = "%s"
	database = "%s"
	schema = "%s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}

		column {
			name = "val2"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR"
}
`, name, databaseName, schemaName)
}
