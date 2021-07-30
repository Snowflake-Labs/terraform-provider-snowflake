package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSchemas(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: schemas(databaseName, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_schemas.s", "database", databaseName),
					resource.TestCheckResourceAttrSet("data.snowflake_schemas.s", "schemas.#"),
					resource.TestCheckResourceAttr("data.snowflake_schemas.s", "schemas.#", "3"),
				),
			},
		},
	})
}

func schemas(databaseName string, schemaName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	data snowflake_schemas "s" {
		database = snowflake_schema.s.database
		depends_on = [snowflake_schema.s]
	}
	`, databaseName, schemaName)
}
