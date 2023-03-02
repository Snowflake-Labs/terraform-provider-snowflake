package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagGrantConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "database_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "schema_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "tag_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_tag_grant.test", "privilege", "APPLY"),
				),
			},
		},
	})
}

func tagGrantConfig(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_database" "test" {
		name = "%v"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_schema" "test" {
		name = "%v"
		database = snowflake_database.test.name
		comment = "Terraform acceptance test"
	}

	resource "snowflake_role" "test" {
		name = "%v"
	}

	resource "snowflake_tag" "test" {
		name = "%v"
		database = snowflake_database.test.name
		schema = snowflake_schema.test.name
		allowed_values = []
	}

	resource "snowflake_tag_grant" "test" {
		tag_name = snowflake_tag.test.name
		database_name = snowflake_database.test.name
		roles         = [snowflake_role.test.name]
		schema_name   = snowflake_schema.test.name
		privilege = "APPLY"
		
	}
	`, name, name, name, name)
}
