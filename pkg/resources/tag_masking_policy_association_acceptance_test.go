package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TagMaskingPolicyAssociation(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAttachmentConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "masking_policy_id", fmt.Sprintf("%[1]v|%[1]v|%[1]v", accName)),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_association.test", "tag_id", fmt.Sprintf("%[1]v|%[1]v|%[1]v", accName)),
				),
			},
		},
	})
}

func tagAttachmentConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%[1]v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	allowed_values = []
	comment = "Terraform acceptance test"
}

resource "snowflake_masking_policy" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
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
`, n)
}
