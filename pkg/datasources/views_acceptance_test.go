package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1423486): Fix using warehouse in all tests and remove unsetting testenvs.ConfigureClientOnce.
func TestAcc_Views(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	viewNamePrefix := acc.TestClient().Ids.Alpha()
	viewName := viewNamePrefix + "1" + acc.TestClient().Ids.Alpha()
	viewName2 := viewNamePrefix + "2" + acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: views(acc.TestDatabaseName, acc.TestSchemaName, schemaId.Name(), viewName, viewName2, viewNamePrefix+"%"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.name", viewName2),
					resource.TestCheckResourceAttrSet("data.snowflake_views.in_schema", "views.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.kind", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.reserved", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.database_name", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_views.in_schema", "views.0.show_output.0.owner"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("data.snowflake_views.in_schema", "views.0.show_output.0.text"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.is_secure", "false"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.is_materialized", "false"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.owner_role_type", "ROLE"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.show_output.0.change_tracking", "OFF"),

					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.name", "ROLE_NAME"),
					resource.TestCheckResourceAttrSet("data.snowflake_views.in_schema", "views.0.describe_output.0.type"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.kind", "COLUMN"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.is_nullable", "true"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.default", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.is_primary", "false"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.is_unique", "false"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.check", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.expression", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.comment", ""),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.policy_name", ""),
					resource.TestCheckNoResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.0.policy_domain"),
					resource.TestCheckResourceAttr("data.snowflake_views.in_schema", "views.0.describe_output.1.name", "ROLE_OWNER"),

					resource.TestCheckResourceAttr("data.snowflake_views.filtering", "views.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_views.filtering", "views.0.show_output.0.name", viewName2),
				),
			},
		},
	})
}

func views(databaseName, defaultSchemaName, schemaName, view1Name, view2Name, viewPrefix string) string {
	return fmt.Sprintf(`
	resource snowflake_schema "test" {
		database = "%[1]v"
		name = "%[3]v" 
	}

	resource snowflake_view "v1"{
		database = "%[1]v"
		schema 	 = "%[2]v"
		name 	 = "%[4]v"
		statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"
		column {
			column_name = "ROLE_NAME"
		}
		column {
			column_name = "ROLE_OWNER"
		}
	}

	resource snowflake_view "v2"{
		database = snowflake_schema.test.database
		schema = snowflake_schema.test.name
		name 	 = "%[5]v"
		statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"
		column {
			column_name = "ROLE_NAME"
		}
		column {
			column_name = "ROLE_OWNER"
		}
	}

	data snowflake_views "in_schema" {
		depends_on = [ snowflake_view.v1, snowflake_view.v2 ]
		in {
			schema = snowflake_schema.test.fully_qualified_name
		}
	}

	data snowflake_views "filtering" {
		depends_on = [ snowflake_view.v1, snowflake_view.v2 ]
		in {
			database = snowflake_schema.test.database
		}
		like = "%[6]v"
		starts_with = trimsuffix("%[6]v", "%%")
		limit {
			rows = 1
			from = snowflake_view.v1.name 
		}
	}
	`, databaseName, defaultSchemaName, schemaName, view1Name, view2Name, viewPrefix)
}
