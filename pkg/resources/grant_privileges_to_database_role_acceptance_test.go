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

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase(t *testing.T) {
	name := "test_database_role_name"
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeModify)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
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
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchema(t *testing.T) {
	name := "test_database_role_name"
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
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
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnAllSchemasInDatabase(t *testing.T) {
	name := "test_database_role_name"
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
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
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject(t *testing.T) {
	name := "test_database_role_name"
	tableName := "test_database_role_table_name"
	configVariables := config.Variables{
		"name":       config.StringVariable(name),
		"table_name": config.StringVariable(tableName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
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
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName).FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
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
	if err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(databaseRoleId).WithOrReplace(true)); err != nil {
		t.Fatal(fmt.Errorf("error database role (%s): %w", databaseRoleId.FullyQualifiedName(), err))
	}
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
			// usage is the default privilege available after creation (it won't be revoked)
			if grant.Privilege != "USAGE" {
				grantedPrivileges = append(grantedPrivileges, grant.Privilege)
			}
		}
		if len(grantedPrivileges) > 0 {
			return fmt.Errorf("database role (%s) still grants , granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
		}
	}
	return nil
}
