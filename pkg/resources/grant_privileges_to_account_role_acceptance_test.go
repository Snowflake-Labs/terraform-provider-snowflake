package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_GrantPrivilegesToAccountRole_OnAccount(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateRole)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE,CREATE ROLE|OnAccount", roleFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccount_PrivilegesReversed(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
		),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.GlobalPrivilegeCreateDatabase)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.GlobalPrivilegeCreateRole)),
					resource.TestCheckResourceAttr(resourceName, "on_account", "true"),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE,CREATE ROLE|OnAccount", roleFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAccountObject(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", databaseName),
					// TODO (SNOW-999049): Even if the identifier is passed as a non-escaped value it will be escaped in the identifier and later on in the CRUD operations (right now, it's "only" read, which can cause behavior similar to always_apply)
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE DATABASE ROLE,CREATE SCHEMA|OnAccountObject|DATABASE|\"%s\"", roleFullyQualifiedName, acc.TestDatabaseName)),
				),
			},
			{
				// TODO (SNOW-999049): this fails, because after import object_name identifier is escaped (was unescaped)
				// 	Make grant_privileges_to_account_role and grant_privileges_to_account_role identifiers accept
				//  quoted and unquoted identifiers.
				// ConfigPlanChecks: resource.ConfigPlanChecks{
				//	PostApplyPostRefresh: []plancheck.PlanCheck{
				//		plancheck.ExpectEmptyPlan(),
				//	},
				// },
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This proves that infinite plan is not produced as in snowflake_grant_privileges_to_role.
// More details can be found in the fix pr https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2364.
func TestAcc_GrantPrivilegesToApplicationRole_OnAccountObject_InfinitePlan(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject_InfinitePlan"),
				ConfigVariables: config.Variables{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchema(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	schemaName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchema_ExactlyOneOf(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema_ExactlyOneOf"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnAllSchemasInDatabase(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnAllSchemasInDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnFutureSchemasInDatabase(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateTable)),
			config.StringVariable(string(sdk.SchemaPrivilegeModify)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.future_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnFutureSchemasInDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnFutureSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	tblName := "test_database_role_table_name"
	tableName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tblName).FullyQualifiedName()
	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tblName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", tableName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnObject|TABLE|%s", roleFullyQualifiedName, tableName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnObject_OwnershipPrivilege(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"name":       config.StringVariable(name),
		"table_name": config.StringVariable(tableName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectOwnership)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnAll|TABLES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAllPipes(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
		),
		"database":          config.StringVariable(databaseName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypePipes)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MONITOR|OnSchemaObject|OnAll|PIPES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_InDatabase(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnFuture|TABLES|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO [SNOW-1272222]: fix the test when it starts working on Snowflake side
func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	t.Skip("Fix after it starts working on Snowflake side, reference: SNOW-1272222")
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported feature 'STREAMLIT'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(databaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnAll|STREAMLITS|InDatabase|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_UpdatePrivileges(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(roleFullyQualifiedName),
			"database": config.StringVariable(databaseName),
		}
		if allPrivileges {
			configVariables["all_privileges"] = config.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]config.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = config.StringVariable(string(privilege))
			}
			configVariables["privileges"] = config.ListVariable(configPrivileges...)
		}
		return configVariables
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeModify,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeMonitor,
					sdk.AccountObjectPrivilegeUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeModify,
					sdk.AccountObjectPrivilegeMonitor,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MODIFY,MONITOR|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_UpdatePrivileges_SnowflakeChecked(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	name := roleId.Name()
	schemaName := "test_database_role_schema_name"
	configVariables := func(allPrivileges bool, privileges []string, schemaName string) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(roleId.FullyQualifiedName()),
			"database": config.StringVariable(acc.TestDatabaseName),
		}
		if allPrivileges {
			configVariables["all_privileges"] = config.BoolVariable(allPrivileges)
		}
		if len(privileges) > 0 {
			configPrivileges := make([]config.Variable, len(privileges))
			for i, privilege := range privileges {
				configPrivileges[i] = config.StringVariable(privilege)
			}
			configVariables["privileges"] = config.ListVariable(configPrivileges...)
		}
		if len(schemaName) > 0 {
			configVariables["schema_name"] = config.StringVariable(schemaName)
		}
		return configVariables
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				}, ""),
				Check: queriedAccountRolePrivilegesEqualTo(
					roleId,
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
				ConfigVariables: configVariables(true, []string{}, ""),
				Check: queriedAccountRolePrivilegesContainAtLeast(
					roleId,
					sdk.AccountObjectPrivilegeCreateDatabaseRole.String(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
					sdk.AccountObjectPrivilegeUsage.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				}, ""),
				Check: queriedAccountRolePrivilegesEqualTo(
					roleId,
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
				ConfigVariables: configVariables(false, []string{
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				}, schemaName),
				Check: queriedAccountRolePrivilegesEqualTo(
					roleId,
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_AlwaysApply(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(roleFullyQualifiedName),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(databaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivileges(t *testing.T) {
	sharedDatabaseName := acc.TestClient().Ids.Alpha()
	shareName := acc.TestClient().Ids.Alpha()
	roleName := acc.TestClient().Ids.Alpha()
	secondaryAccountName := acc.SecondaryTestClient().Context.CurrentAccount(t)
	configVariables := config.Variables{
		"role_name":            config.StringVariable(roleName),
		"shared_database_name": config.StringVariable(sharedDatabaseName),
		"share_name":           config.StringVariable(shareName),
		"account_name":         config.StringVariable(secondaryAccountName),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.AccountObjectPrivilegeImportedPrivileges.String()),
		),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createSharedDatabaseOnSecondaryAccount(t, sharedDatabaseName, shareName) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivileges"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivileges"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1998 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ImportedPrivilegesOnSnowflakeDatabase(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"role_name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.AccountObjectPrivilegeImportedPrivileges.String()),
		),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
				ConfigVariables: configVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", "\"SNOWFLAKE\""),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.AccountObjectPrivilegeImportedPrivileges.String()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/ImportedPrivilegesOnSnowflakeDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO(SNOW-1213622): Add test for custom applications using on_account_object.object_type = "DATABASE"

func TestAcc_GrantPrivilegesToAccountRole_MultiplePartsInRoleName(t *testing.T) {
	nameBytes := []byte(acc.TestClient().Ids.Alpha())
	nameBytes[3] = '.'
	nameBytes[6] = '.'
	name := string(nameBytes)
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.GlobalPrivilegeCreateDatabase)),
			config.StringVariable(string(sdk.GlobalPrivilegeCreateRole)),
		),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, name)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccount"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", name),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2533 is fixed
func TestAcc_GrantPrivilegesToAccountRole_OnExternalVolume(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	externalVolumeName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"name":            config.StringVariable(roleFullyQualifiedName),
		"external_volume": config.StringVariable(externalVolumeName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
		),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
					cleanupExternalVolume := createExternalVolume(t, externalVolumeName)
					t.Cleanup(cleanupExternalVolume)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnExternalVolume"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_type", "EXTERNAL VOLUME"),
					resource.TestCheckResourceAttr(resourceName, "on_account_object.0.object_name", externalVolumeName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|USAGE|OnAccountObject|EXTERNAL VOLUME|\"%s\"", roleFullyQualifiedName, externalVolumeName)),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2651
// TODO [SNOW-1270457]: This seems to be a Snowflake error, we are waiting for the confirmation. Alter the test when the behavior is fixed. Update the resource documentation (section known issues).
func TestAcc_GrantPrivilegesToAccountRole_MLPrivileges(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	configVariables := config.Variables{
		"name": config.StringVariable(roleFullyQualifiedName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	schemaName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", roleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SNOWFLAKE.ML.ANOMALY_DETECTION,CREATE SNOWFLAKE.ML.FORECAST|OnSchema|OnSchema|%s", roleFullyQualifiedName, schemaName)),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ChangeWithGrantOptionsOutsideOfTerraform_WithGrantOptions(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	tableName := acc.TestClient().Ids.Alpha()

	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tableName),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnTableToAccountRole(
						t,
						roleId,
						sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName),
						[]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeTruncate},
						false,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToAccountRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	tableName := acc.TestClient().Ids.Alpha()

	configVariables := config.Variables{
		"name":       config.StringVariable(roleFullyQualifiedName),
		"table_name": config.StringVariable(tableName),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.SchemaObjectPrivilegeTruncate.String()),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnTableToAccountRole(
						t,
						roleId,
						sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName),
						[]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeTruncate},
						true,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
			},
		},
	})
}

func revokeAndGrantPrivilegesOnTableToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	tableName sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	client := acc.Client(t)
	ctx := context.Background()
	err := client.Grants.RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: sdk.ObjectTypeTable,
					Name:       tableName,
				},
			},
		},
		accountRoleId,
		new(sdk.RevokePrivilegesFromAccountRoleOptions),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Grants.GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: sdk.ObjectTypeTable,
					Name:       tableName,
				},
			},
		},
		accountRoleId,
		&sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveGrantedObjectOutsideTerraform(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	var databaseCleanup func()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup = acc.TestClient().Database.CreateDatabaseWithName(t, databaseName)
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToAccountRole_RemoveAccountRoleOutsideTerraform(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"name":     config.StringVariable(roleFullyQualifiedName),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateDatabaseRole)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	var roleCleanup func()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, dbCleanup := acc.TestClient().Database.CreateDatabaseWithName(t, databaseName)
					t.Cleanup(dbCleanup)
					_, roleCleanup = acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { roleCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/OnAccountObject"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to account role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2689 is fixed
func TestAcc_GrantPrivilegesToAccountRole_AlwaysApply_SetAfterCreate(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleName := roleId.Name()
	roleFullyQualifiedName := roleId.FullyQualifiedName()
	databaseName := acc.TestClient().Ids.DatabaseId().FullyQualifiedName()
	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(roleFullyQualifiedName),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(databaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}
	resourceName := "snowflake_grant_privileges_to_account_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, roleName)
					t.Cleanup(roleCleanup)
				},
				ConfigDirectory:    acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToAccountRole/AlwaysApply"),
				ConfigVariables:    configVariables(true),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnAccountObject|DATABASE|%s", roleFullyQualifiedName, databaseName)),
				),
			},
		},
	})
}

func createSharedDatabaseOnSecondaryAccount(t *testing.T, databaseName string, shareName string) {
	t.Helper()

	database, databaseCleanup := acc.SecondaryTestClient().Database.CreateDatabaseWithName(t, databaseName)
	t.Cleanup(databaseCleanup)

	share, shareCleanup := acc.SecondaryTestClient().Share.CreateShareWithName(t, shareName)
	t.Cleanup(shareCleanup)

	acc.SecondaryTestClient().Role.GrantPrivilegeOnDatabaseToShare(t, database.ID(), share.ID())

	accountName := acc.TestClient().Context.CurrentAccount(t)
	accountId := sdk.NewAccountIdentifierFromAccountLocator(accountName)
	acc.SecondaryTestClient().Share.SetAccountOnShare(t, accountId, share.ID())
}

func queriedAccountRolePrivilegesEqualTo(roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return queriedPrivilegesEqualTo(func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error) {
		return client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleName,
			},
		})
	}, privileges...)
}

func queriedAccountRolePrivilegesContainAtLeast(roleName sdk.AccountObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return queriedPrivilegesContainAtLeast(func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error) {
		return client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: roleName,
			},
		})
	}, roleName, privileges...)
}

func createExternalVolume(t *testing.T, externalVolumeName string) func() {
	t.Helper()

	client := acc.Client(t)
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`create external volume "%s" storage_locations = (
    (
        name = 'test' 
        storage_provider = 's3' 
        storage_base_url = 's3://my_example_bucket/'
        storage_aws_role_arn = 'arn:aws:iam::123456789012:role/myrole'
        encryption=(type='aws_sse_kms' kms_key_id='1234abcd-12ab-34cd-56ef-1234567890ab')
    )
)
`, externalVolumeName))
	require.NoError(t, err)

	return func() {
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`drop external volume "%s"`, externalVolumeName))
		require.NoError(t, err)
	}
}
