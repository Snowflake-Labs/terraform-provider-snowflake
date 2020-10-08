package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTable(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func tableConfig(name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR"
	}
}
`
	return fmt.Sprintf(s, name, name, name)
}
