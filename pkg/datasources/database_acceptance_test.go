package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Database(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: database(databaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_database.t", "comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "owner"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "retention_time"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "is_current"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "is_default"),
				),
			},
		},
	})
}

func database(databaseName, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_database" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_database "t" {
			depends_on = [snowflake_database.test_database]
			name = "%v"
		}
	`, databaseName, comment, databaseName)
}
