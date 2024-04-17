package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagMaskingPolicyAssociationBasic(t *testing.T) {
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
				Config: tagAttachmentConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "masking_policy_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, accName)),
				),
			},
		},
	})
}

func TestAcc_TagMaskingPolicyAssociationsystem_functions_integration_testComplete(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_tag.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(name),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
			"comment":  config.StringVariable("Terraform acceptance test"),
		}
	}
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagMaskingPolicyAssociation/basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.0", "alv1"),
					resource.TestCheckResourceAttr(resourceName, "allowed_values.1", "alv2"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),

					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "masking_policy_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, name)),
				),
			},

			// test - change comment
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagMaskingPolicyAssociation/basic"),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
				),
			},

			// test - import
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_TagMaskingPolicyAssociation/basic"),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func tagAttachmentConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = "%[2]s"
	schema = "%[3]s"
	allowed_values = []
	comment = "Terraform acceptance test"
}

resource "snowflake_masking_policy" "test" {
	name = "%[1]v"
	database = "%[2]s"
	schema = "%[3]s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR(16777216)"
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_masking_policy_association" "test" {
	tag_id = snowflake_tag.test.id
	masking_policy_id = snowflake_masking_policy.test.id
}
`, n, databaseName, schemaName)
}
