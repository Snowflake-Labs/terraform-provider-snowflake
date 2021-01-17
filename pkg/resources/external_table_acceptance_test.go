package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExternalTable(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: externalTableConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func externalTableConfig(name string) string {
	s := `
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_stage" "test" {
	name = "%v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_external_table" "test_table" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "%v"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
    as = "($1:\"CreatedDate\"::timestamp)"
	}
	column {
		name = "column2"
		type = "VARCHAR"
    as = "($1:\"CreatedDate\"::timestamp)"
	}
  file_format = "TYPE = CSV"
  location = "@${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_stage.test.name}"
}
`
	return fmt.Sprintf(s, name, name, name, name)
}
