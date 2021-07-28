package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Stream(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streamConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
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
	database        = snowflake_database.test_database.name
	schema          = snowflake_schema.test_schema.name
	name            = "STREAM_ON_TABLE"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_stream" "test_stream" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
}
`
	return fmt.Sprintf(s, name, name, name)
}
