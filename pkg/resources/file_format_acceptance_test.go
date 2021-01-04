package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccFileFormat(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: fileFormatConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "format_type", "CSV"),
					resource.TestCheckResourceAttr("snowflake_file_format.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func fileFormatConfig(n string) string {
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

resource "snowflake_file_format" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	format_type = "CSV"
	comment = "Terraform acceptance test"
}
`, n, n, n)
}
