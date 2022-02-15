package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDatabases(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{

		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: databases(databaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.#"),
					resource.TestCheckResourceAttr("data.snowflake_databases.t", "databases.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_databases.t", "databases.0.name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_databases.t", "databases.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.owner"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.retention_time"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.options"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.origin"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.is_current"),
					resource.TestCheckResourceAttrSet("data.snowflake_databases.t", "databases.0.is_default"),
				),
			},
		},
	})
}

func databases(databaseName, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_database" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_databases "t" {
			depends_on = [snowflake_database.test_database]
		}
	`, databaseName, comment)
}
