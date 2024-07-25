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
)

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase_PrivilegesReversed(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeModify)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchema(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	schemaName := acc.TestClient().Ids.SchemaId().FullyQualifiedName()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), schemaName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchema_ExactlyOneOf(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema_ExactlyOneOf"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnAllSchemasInDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnAllSchemasInDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnFutureSchemasInDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnFutureSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.future_schemas_in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnFutureSchemasInDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnFutureSchemasInDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnObject(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	databaseRoleName := databaseRoleId.Name()

	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableName := tableId.Name()

	resourceName := "snowflake_grant_privileges_to_database_role.test"
	configVariables := config.Variables{
		"name":       config.StringVariable(databaseRoleName),
		"table_name": config.StringVariable(tableName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleName)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", tableId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnObject|TABLE|%s", databaseRoleId.FullyQualifiedName(), tableId.FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnObject_OwnershipPrivilege(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	tableName := "test_database_role_table_name"
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnAll|TABLES|InDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAllPipes(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypePipes)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MONITOR|OnSchemaObject|OnAll|PIPES|InDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFuture_InDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnFuture|TABLES|InDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO [SNOW-1272222]: fix the test when it starts working on Snowflake side
func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnFuture_Streamlits_InDatabase(t *testing.T) {
	t.Skip("Fix after it starts working on Snowflake side, reference: SNOW-1272222")
	name := acc.TestClient().Ids.Alpha()
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported feature 'STREAMLIT'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnAll|STREAMLITS|InDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(databaseRoleId.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
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
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeCreateSchema,
					sdk.AccountObjectPrivilegeModify,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/privileges"),
				ConfigVariables: configVariables(false, []sdk.AccountObjectPrivilege{
					sdk.AccountObjectPrivilegeModify,
					sdk.AccountObjectPrivilegeMonitor,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "false"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MODIFY,MONITOR|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges_SnowflakeChecked(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(allPrivileges bool, privileges []string, schemaName string) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(databaseRoleId.Name()),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				}, ""),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					databaseRoleId,
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
				ConfigVariables: configVariables(true, []string{}, ""),
				Check: queriedPrivilegesToDatabaseRoleContainAtLeast(
					databaseRoleId,
					sdk.AccountObjectPrivilegeCreateDatabaseRole.String(),
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
					sdk.AccountObjectPrivilegeUsage.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				}, ""),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					databaseRoleId,
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
				ConfigVariables: configVariables(false, []string{
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				}, schemaId.Name()),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					databaseRoleId,
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_AlwaysApply(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(databaseRoleId.Name()),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(acc.TestDatabaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

// proved https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2651
func TestAcc_GrantPrivilegesToDatabaseRole_MLPrivileges(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", acc.TestClient().Ids.SchemaId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SNOWFLAKE.ML.ANOMALY_DETECTION,CREATE SNOWFLAKE.ML.FORECAST|OnSchema|OnSchema|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.SchemaId().FullyQualifiedName())),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithGrantOptions(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	name := databaseRoleId.Name()
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(true),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
						t, databaseRoleId,
						acc.TestClient().Ids.DatabaseId(),
						[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
						false,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2459 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithoutGrantOptions(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	name := databaseRoleId.Name()
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig: func() {
					revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
						t, databaseRoleId,
						acc.TestClient().Ids.DatabaseId(),
						[]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeCreateSchema},
						true,
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToDatabaseRole_RemoveGrantedObjectOutsideTerraform(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	configVariables := config.Variables{
		"name":     config.StringVariable(name),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup = acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(databaseCleanup)
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleInDatabaseWithName(t, databaseId, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to database role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToDatabaseRole_RemoveDatabaseRoleOutsideTerraform(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	configVariables := config.Variables{
		"name":     config.StringVariable(name),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"with_grant_option": config.BoolVariable(true),
	}

	var databaseRoleCleanup func()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, dbCleanup := acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(dbCleanup)
					_, databaseRoleCleanup = acc.TestClient().DatabaseRole.CreateDatabaseRoleInDatabaseWithName(t, databaseId, name)
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { databaseRoleCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to database role"),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2689 is fixed
func TestAcc_GrantPrivilegesToDatabaseRole_AlwaysApply_SetAfterCreate(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(databaseRoleId.Name()),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(acc.TestDatabaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory:    acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables:    configVariables(true),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2960
func TestAcc_GrantPrivilegesToDatabaseRole_CreateNotebooks(t *testing.T) {
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()

	configVariables := config.Variables{
		"name": config.StringVariable(databaseRoleId.Name()),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateNotebook)),
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
		CheckDestroy: acc.CheckDatabaseRolePrivilegesRevoked(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRoleWithName(t, databaseRoleId.Name())
					t.Cleanup(databaseRoleCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateNotebook)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE NOTEBOOK|OnSchema|OnAllSchemasInDatabase|%s", databaseRoleId.FullyQualifiedName(), acc.TestClient().Ids.DatabaseId().FullyQualifiedName())),
				),
			},
		},
	})
}

func queriedPrivilegesToDatabaseRoleEqualTo(databaseRoleName sdk.DatabaseObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return queriedPrivilegesEqualTo(func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error) {
		return client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleName,
			},
		})
	}, privileges...)
}

func queriedPrivilegesToDatabaseRoleContainAtLeast(databaseRoleName sdk.DatabaseObjectIdentifier, privileges ...string) func(s *terraform.State) error {
	return queriedPrivilegesContainAtLeast(func(client *sdk.Client, ctx context.Context) ([]sdk.Grant, error) {
		return client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleName,
			},
		})
	}, databaseRoleName, privileges...)
}

func revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
	t *testing.T,
	databaseRoleName sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	client := acc.Client(t)
	ctx := context.Background()
	err := client.Grants.RevokePrivilegesFromDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleName,
		new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Grants.GrantPrivilegesToDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleName,
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
