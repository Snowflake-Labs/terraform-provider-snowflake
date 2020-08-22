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
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					checkBool("snowflake_table.test_table", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_table.test_table", "is_managed", false),
				),
			},
		},
	})
}

func tableConfig(name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%v"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%v"
	comment  = "Terraform acceptance test"
	columns  = {
		"column1" = "VARIANT",
		"column2" = "VARCHAR"
	}
}
`
	return fmt.Sprintf(s, name)
}
