package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceName = "snowflake_database_role.test_db_role"
	dbName       = "db_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	dbRoleName   = "db_role_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment      = "dummy"
	comment2     = "test comment"
)

func TestAcc_DatabaseRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoleConfig(dbName, dbRoleName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", dbName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: databaseRoleConfig(dbName, dbRoleName, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", dbName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment2),
				),
			},
		},
	})
}

func databaseRoleConfig(dbName string, dbRoleName string, comment string) string {
	s := `
resource "snowflake_database" "test_db" {
	name = "%s"
}

resource "snowflake_database_role" "test_db_role" {
	name     	  = "%s"
	database  	  = snowflake_database.test_db.name
	comment       = "%s"
}
	`
	return fmt.Sprintf(s, dbName, dbRoleName, comment)
}
