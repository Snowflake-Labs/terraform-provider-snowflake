package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Schema(t *testing.T) {
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaConfig(schemaName, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", schemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func TestAcc_SchemaRename(t *testing.T) {
	oldSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: schemaConfig(oldSchemaName, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
			{
				Config: schemaConfig(newSchemaName, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func schemaConfig(schemaName string, databaseName string) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	name = "%v"
	database = "%s"
	comment = "Terraform acceptance test"
}
`, schemaName, databaseName)
}
