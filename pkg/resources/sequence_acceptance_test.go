package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Sequence(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accRename := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: sequenceConfig(accName, accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`%v.%v.%v`, accName, accName, accName)),
				),
			},
			// Set comment and rename
			{
				Config: sequenceConfigWithComment(accName, accRename, "look at me I am a comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accRename),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", "look at me I am a comment"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`%v.%v.%v`, accName, accName, accRename)),
				),
			},
			// Unset comment and set increment
			{
				Config: sequenceConfigWithIncrement(accName, accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "increment", "32"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`%v.%v.%v`, accName, accName, accName)),
				),
			},
		},
	})
}

func sequenceConfigWithIncrement(name, sequenceName string) string {
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

resource "snowflake_sequence" "test_sequence" {
	database   = snowflake_database.test_database.name
	schema     = snowflake_schema.test_schema.name
	name       = "%s"
  increment = 32
}
`
	return fmt.Sprintf(s, name, name, sequenceName)
}
func sequenceConfig(name, sequenceName string) string {
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

resource "snowflake_sequence" "test_sequence" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
}
`
	return fmt.Sprintf(s, name, name, sequenceName)
}

func sequenceConfigWithComment(name, sequenceName, comment string) string {
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

resource "snowflake_sequence" "test_sequence" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
  comment  = "%s"
}
`
	return fmt.Sprintf(s, name, name, sequenceName, comment)
}
