package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_DatabaseRole(t *testing.T) {
	resourceName := "snowflake_database_role.test_db_role"
	dbRoleName := "db_role_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "dummy"
	comment2 := "test comment"

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: databaseRoleConfig(dbRoleName, acc.TestDatabaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: databaseRoleConfig(dbRoleName, acc.TestDatabaseName, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", dbRoleName),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment2),
				),
			},
		},
	})
}

func databaseRoleConfig(dbRoleName string, databaseName string, comment string) string {
	s := `
resource "snowflake_database_role" "test_db_role" {
	name     	  = "%s"
	database  	  = "%s"
	comment       = "%s"
}
	`
	return fmt.Sprintf(s, dbRoleName, databaseName, comment)
}
