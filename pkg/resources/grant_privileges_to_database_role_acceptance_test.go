package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

// TODO Use cases to cover in acc tests
// - basic - check create, read and destroy
// 		- grant privileges on database
// - update - check update of privileges
// 		- privileges
//		- privileges to all_privileges
//		- all_privileges to privilege
// - import - check import
// 		- different paths to parse (on database, on schema, on schema object)

func TestAcc_GrantPrivilegesToDatabaseRole_basic(t *testing.T) {
	name := "test_database_role_name"
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable("CREATE SCHEMA"),
			config.StringVariable("MODIFY"),
			config.StringVariable("USAGE"),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, name) },
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "CREATE SCHEMA"),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", "MODIFY"),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", "USAGE"),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestDatabaseName),
				),
			},
		},
	})
}

func createDatabaseRoleOutsideTerraform(t *testing.T, name string) {
	t.Helper()
	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	databaseRoleId := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name)
	if err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(databaseRoleId)); err != nil {
		t.Fatal(fmt.Errorf("error database role (%s): %w", databaseRoleId.FullyQualifiedName(), err))
	}
	t.Cleanup(func() {
		if err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(databaseRoleId)); err != nil {
			t.Fatal(fmt.Errorf("failed to drop database role (%s): %w", databaseRoleId.FullyQualifiedName(), err))
		}
	})
}

func testAccCheckDatabaseRolePrivilegesRevoked(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_grant_privileges_to_database_role" {
			continue
		}
		ctx := context.Background()

		id := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["database_role_name"])
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: id,
			},
		})
		if err != nil {
			return err
		}
		var grantedPrivileges []string
		for _, grant := range grants {
			grantedPrivileges = append(grantedPrivileges, grant.Privilege)
		}
		if len(grantedPrivileges) > 0 {
			return fmt.Errorf("database role (%s) still grants , granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
		}
	}
	return nil
}
