package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccViews(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: views(databaseName, schemaName, viewName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_views.v", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_views.v", "views.#"),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "views.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "views.0.name", viewName),
				),
			},
		},
	})
}

func views(databaseName string, schemaName string, viewName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_view "v"{
		name 	 = "%v"
		database = snowflake_schema.s.database
		schema 	 = snowflake_schema.s.name
		statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"
	}

	data snowflake_views "v" {
		database = snowflake_view.v.database
		schema = snowflake_view.v.schema
	}
	`, databaseName, schemaName, viewName)
}
