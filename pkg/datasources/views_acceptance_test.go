package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1423486): Fix using warehouse in all tests.
func TestAcc_Views(t *testing.T) {
	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: views(viewId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_views.v", "database", viewId.DatabaseName()),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "schema", viewId.SchemaName()),
					resource.TestCheckResourceAttrSet("data.snowflake_views.v", "views.#"),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "views.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_views.v", "views.0.name", viewId.Name()),
				),
			},
		},
	})
}

func views(viewId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
	resource "snowflake_unsafe_execute" "use_warehouse" {
		execute = "USE WAREHOUSE \"%v\""
		revert  = "SELECT 1"
	}

	resource snowflake_view "v"{
		name 	 = "%v"
		schema 	 = "%v"
		database = "%v"
		statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"
		depends_on = [snowflake_unsafe_execute.use_warehouse]
	}

	data snowflake_views "v" {
		database = snowflake_view.v.database
		schema = snowflake_view.v.schema
		depends_on = [snowflake_view.v]
	}
	`, acc.TestWarehouseName, viewId.Name(), viewId.SchemaName(), viewId.DatabaseName())
}
