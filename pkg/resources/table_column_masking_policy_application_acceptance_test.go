package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TableColumnMaskingPolicyApplication(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyApplicationTestConfig(),
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

func maskingPolicyApplicationTestConfig() string {
	return `
resource "snowflake_masking_policy" "test" {
	name               = "mypolicy"
	database           = "terraform_test_database"
	schema             = "terraform_test_schema"
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
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
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
}`
}
