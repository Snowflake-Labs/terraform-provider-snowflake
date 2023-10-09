package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TableColumnMaskingPolicyApplication(t *testing.T) {
	database := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyApplicationTestConfig(database),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_column_masking_policy_application.mpa", "table", fmt.Sprintf(`"%s"."test_schema"."table"`, database)),
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

func maskingPolicyApplicationTestConfig(database string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "test_schema"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_masking_policy" "test" {
	name               = "mypolicy"
	database           = snowflake_database.test.name
	schema             = snowflake_schema.test.name
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
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
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
}`,
		database)
}
