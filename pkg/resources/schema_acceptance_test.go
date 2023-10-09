package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Schema(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaConfig(databaseName, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func TestAcc_SchemaRename(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	oldSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaConfig(databaseName, oldSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
			{
				Config: schemaConfig(databaseName, newSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", databaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func schemaConfig(databaseName string, schemaName string) string {
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
`, databaseName, schemaName)
}
