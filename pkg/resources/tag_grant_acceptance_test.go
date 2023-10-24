package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TagGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagGrantConfig(accName, "APPLY", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "tag_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "privilege", "APPLY"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: tagGrantConfig(accName, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "tag_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			{
				ResourceName:      "snowflake_tag_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func tagGrantConfig(name string, privilege string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "test" {
		name = "%v"
	}

	resource "snowflake_tag" "test" {
		name = "%v"
		database = "%s"
		schema = "%s"
		allowed_values = []
	}

	resource "snowflake_tag_grant" "test" {
		tag_name = snowflake_tag.test.name
		database_name = "%s"
		roles         = [snowflake_role.test.name]
		schema_name   = "%s"
		privilege = "%s"

	}
	`, name, name, databaseName, schemaName, databaseName, schemaName, privilege)
}
