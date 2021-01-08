package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_MaskingPolicy(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MASKING_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccMaskingPolicy")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "value_data_type", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "masking_expression", "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"),
					resource.TestCheckResourceAttr("snowflake_masking_policy.test", "return_data_type", "VARCHAR(16777216)"),
				),
			},
		},
	})
}

func maskingPolicyConfig(n string) string {
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
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	value_data_type = "VARCHAR"
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR(16777216)"
	comment = "Terraform acceptance test"
}
`, n, n, n)
}
