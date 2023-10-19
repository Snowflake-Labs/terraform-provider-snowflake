package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceName = "snowflake_database_role.test_db_role"
	dbRoleName   = "db_role_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment      = "dummy"
	comment2     = "test comment"
)

func TestAcc_DatabaseRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoleConfig(dbRoleName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: databaseRoleConfig(dbRoleName, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment2),
				),
			},
		},
	})
}

func databaseRoleConfig(dbRoleName string, comment string) string {
	s := `
resource "snowflake_database_role" "test_db_role" {
	name     	  = "%s"
	database  	  = "terraform_test_database"
	comment       = "%s"
}
	`
	return fmt.Sprintf(s, dbRoleName, comment)
}
