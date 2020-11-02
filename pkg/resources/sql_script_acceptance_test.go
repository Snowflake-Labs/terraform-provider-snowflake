package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSqlScript(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SQL_SCRIPT_TESTS"); ok {
		t.Skip("Skipping TestAccSqlScript")
	}

	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: sqlScriptConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sql_script.test_script", "name", accName),
				),
			},
		},
	})
}

func sqlScriptConfig(name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name = "%s"
}

resource "snowflake_schema" "test_schema" {
	name = "%s"
	database = snowflake_database.database.name
}

resource "snowflake_table" "test_table" {
	name = "%s"
	database = snowflake_database.database.name
	schema = snowflake_schema.schema.name
	column {
		name = "column1"
		type = "VARIANT"
	}
}

resource "snowflake_role" "test_role" {
	name = "%s"
}

resource "snowflake_sql_script" "test_script" {
	depends_on = [
		snowflake_table.table,
	]
	name = "%s"
	lifecycle_commands {
		create = join("", ["GRANT ALL ON ALL TABLES IN DATABASE ", snowflake_database.database.name, " TO ROLE ", snowflake_role.role.name, ";"])
		delete = join("", ["REVOKE ALL ON ALL TABLES IN DATABASE ", snowflake_database.database.name, " FROM ROLE ", snowflake_role.role.name, ";"])
	}
}
`
	return fmt.Sprintf(s, name, name, name, name, name)
}
