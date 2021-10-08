package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Tag(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tagConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func tagConfig(n string) string {
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

resource "snowflake_tag" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	comment = "Terraform acceptance test"
}
`, n, n, n)
}
