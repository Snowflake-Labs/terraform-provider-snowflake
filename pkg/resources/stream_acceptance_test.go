package resources_test

import (
	"fmt"
<<<<<<< HEAD
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
=======
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
				Config: streamConfig(accName, false),
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
<<<<<<< HEAD
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "stream_on_table")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", true),
				),
			},
=======
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: streamConfig(accName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", true),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				ResourceName:      "snowflake_stream.test_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
		},
	})
}

<<<<<<< HEAD
func streamConfig(name string) string {
=======
func streamConfig(name string, append_only bool) string {
	append_only_config := ""
	if append_only {
		append_only_config = "append_only = true"
	}

>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
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
<<<<<<< HEAD
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "stream_on_table"
	comment  = "Terraform acceptance test"
=======
	database        = snowflake_database.test_database.name
	schema          = snowflake_schema.test_schema.name
	name            = "STREAM_ON_TABLE"
	comment         = "Terraform acceptance test"
	change_tracking = true

>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
<<<<<<< HEAD
		type = "VARCHAR"
=======
		type = "VARCHAR(16777216)"
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
	}
}

resource "snowflake_stream" "test_stream" {
<<<<<<< HEAD
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	append_only = true
	
}
`
	return fmt.Sprintf(s, name, name, name)
=======
	database    = snowflake_database.test_database.name
	schema      = snowflake_schema.test_schema.name
	name        = "%s"
	comment     = "Terraform acceptance test"
	on_table    = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	%s
}
`
	return fmt.Sprintf(s, name, name, name, append_only_config)
>>>>>>> be74d18f7f46c07cc6e4849460ef3eb859a5d53c
}
