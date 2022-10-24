package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_TagAssociation(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_id", fmt.Sprintf("%s|%s|%s", accName, accName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_value", "finance"),
				),
			},
		},
	})
}

func tagAssociationConfig(n string) string {
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
	allowed_values = ["finance", "hr"]
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_association" "test" {
	object_identifier {
		name = snowflake_database.test.name
	  }
	object_type = "DATABASE"
	tag_id = snowflake_tag.test.id
	tag_value = "finance"
}
`, n)
}
