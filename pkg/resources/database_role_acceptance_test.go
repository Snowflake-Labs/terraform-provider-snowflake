package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole(t *testing.T) {
	resourceName := "snowflake_database_role.test_db_role"
	dbRoleName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	comment2 := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.DatabaseRole),
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
