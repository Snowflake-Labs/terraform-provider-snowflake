package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableColumnMaskingPolicyApplication(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyApplicationTestConfig(acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_column_masking_policy_application.mpa", "table", fmt.Sprintf(`"%s"."%s"."table"`, acc.TestDatabaseName, acc.TestSchemaName)),
				),
			},
			{
				ResourceName:      "snowflake_table_column_masking_policy_application.mpa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func maskingPolicyApplicationTestConfig(databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_masking_policy" "test" {
	name               = "mypolicy"
	database           = "%s"
	schema             = "%s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type   = "VARCHAR"
	comment            = "Terraform acceptance test"
}

resource "snowflake_table" "table" {
	database = "%s"
	schema   = "%s"
	name     = "table"

	column {
	  name     = "secret"
	  type     = "VARCHAR(16777216)"
	}

	lifecycle {
		ignore_changes = [column[0].masking_policy]
	}
}

resource "snowflake_table_column_masking_policy_application" "mpa" {
	table          = snowflake_table.table.qualified_name
	column         = "secret"
	masking_policy = snowflake_masking_policy.test.qualified_name
}`, databaseName, schemaName, databaseName, schemaName)
}
