package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSchemaGrant(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SHARE_TESTS"); ok {
		t.Skip("Skipping TestAccSchemaGrant")
	}

	sName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	shareName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: schemaGrantConfig(sName, roleName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "schema_name", sName),
					resource.TestCheckResourceAttr("snowflake_schema_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_schema_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func schemaGrantConfig(n, role, share string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%v"
}

resource "snowflake_schema" "test" {
  name      = "%v"
  database  = snowflake_database.test.name
  comment   = "Terraform acceptance test"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_share" "test" {
  name     = "%v"
  accounts = ["PC37737"]
}

resource "snowflake_database_grant" "test" {
  database_name = snowflake_schema.test.database
  shares        = [snowflake_share.test.name]
}

resource "snowflake_schema_grant" "test" {
  schema_name   = snowflake_schema.test.name
  database_name = snowflake_schema.test.database
  roles         = [snowflake_role.test.name]
  shares        = [snowflake_share.test.name]

  depends_on = [snowflake_database_grant.test]
}
`, n, n, role, share)
}
