package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccStream(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streamConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "stream_on_table")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", true),
				),
			},
		},
	})
}

func streamConfig(name string) string {
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

resource "snowflake_table" "test_stream_on_table" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "stream_on_table"
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

resource "snowflake_stream" "test_stream" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	append_only = true
	
}
`
	return fmt.Sprintf(s, name, name, name)
}
