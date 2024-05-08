package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schemas(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
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
