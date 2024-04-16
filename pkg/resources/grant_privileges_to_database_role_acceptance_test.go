package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToDatabaseRole_OnDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.AccountObjectPrivilegeUsage)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeModify)),
			config.StringVariable(string(sdk.AccountObjectPrivilegeCreateSchema)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(true),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.AccountObjectPrivilegeCreateSchema)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.AccountObjectPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "privileges.2", string(sdk.AccountObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|true|false|CREATE SCHEMA,MODIFY,USAGE|OnDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	schemaName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnSchema|%s", databaseRoleName, schemaName)),
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.all_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnAllSchemasInDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnFutureSchemasInDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateTable)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeModify)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.future_schemas_in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE TABLE,MODIFY|OnSchema|OnFutureSchemasInDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tblName := "test_database_role_table_name"
	configVariables := config.Variables{
		"name":       config.StringVariable(name),
		"table_name": config.StringVariable(tblName),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	tableName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tblName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.object_name", tableName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnObject|TABLE|%s", databaseRoleName, tableName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnObject"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported privilege 'OWNERSHIP'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_InDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnAll|TABLES|InDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeMonitor)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeMonitor)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypePipes)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MONITOR|OnSchemaObject|OnAll|PIPES|InDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeInsert)),
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUpdate)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeTables.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeInsert)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaObjectPrivilegeUpdate)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeTables)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|INSERT,UPDATE|OnSchemaObject|OnFuture|TABLES|InDatabase|%s", databaseRoleName, databaseName)),
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnFuture_InDatabase"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("Unsupported feature 'STREAMLIT'"),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_OnSchemaObject_OnAll_Streamlits_InDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaObjectPrivilegeUsage)),
		),
		"database":           config.StringVariable(acc.TestDatabaseName),
		"object_type_plural": config.StringVariable(sdk.PluralObjectTypeStreamlits.String()),
		"with_grant_option":  config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchemaObject_OnAll_InDatabase"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaObjectPrivilegeUsage)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.object_type_plural", string(sdk.PluralObjectTypeStreamlits)),
					resource.TestCheckResourceAttr(resourceName, "on_schema_object.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|USAGE|OnSchemaObject|OnAll|STREAMLITS|InDatabase|%s", databaseRoleName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := func(allPrivileges bool, privileges []sdk.AccountObjectPrivilege) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(name),
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,MODIFY|OnDatabase|%s", databaseRoleName, databaseName)),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SCHEMA,USAGE,MONITOR|OnDatabase|%s", databaseRoleName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges/all_privileges"),
				ConfigVariables: configVariables(true, []sdk.AccountObjectPrivilege{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_privileges", "true"),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|MODIFY,MONITOR|OnDatabase|%s", databaseRoleName, databaseName)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_UpdatePrivileges_SnowflakeChecked(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := "test_database_role_schema_name"
	configVariables := func(allPrivileges bool, privileges []string, schemaName string) config.Variables {
		configVariables := config.Variables{
			"name":     config.StringVariable(name),
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

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/privileges"),
				ConfigVariables: configVariables(false, []string{
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				}, ""),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					databaseRoleName,
					sdk.AccountObjectPrivilegeCreateSchema.String(),
					sdk.AccountObjectPrivilegeModify.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/all_privileges"),
				ConfigVariables: configVariables(true, []string{}, ""),
				Check: queriedPrivilegesToDatabaseRoleContainAtLeast(
					databaseRoleName,
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
					databaseRoleName,
					sdk.AccountObjectPrivilegeModify.String(),
					sdk.AccountObjectPrivilegeMonitor.String(),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/UpdatePrivileges_SnowflakeChecked/on_schema"),
				ConfigVariables: configVariables(false, []string{
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				}, schemaName),
				Check: queriedPrivilegesToDatabaseRoleEqualTo(
					databaseRoleName,
					sdk.SchemaPrivilegeCreateTask.String(),
					sdk.SchemaPrivilegeCreateExternalTable.String(),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToDatabaseRole_AlwaysApply(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := func(alwaysApply bool) config.Variables {
		return config.Variables{
			"name":           config.StringVariable(name),
			"all_privileges": config.BoolVariable(true),
			"database":       config.StringVariable(acc.TestDatabaseName),
			"always_apply":   config.BoolVariable(alwaysApply),
		}
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	databaseName := sdk.NewAccountObjectIdentifier(acc.TestDatabaseName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/AlwaysApply"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "always_apply", "true"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|true|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|ALL|OnDatabase|%s", databaseRoleName, databaseName)),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2651
// TODO [SNOW-1270457]: This seems to be a Snowflake error, we are waiting for the confirmation. Alter the test when the behavior is fixed. Update the resource documentation (section known issues).
func TestAcc_GrantPrivilegesToDatabaseRole_MLPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"name": config.StringVariable(name),
		"privileges": config.ListVariable(
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
			config.StringVariable(string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
		),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"with_grant_option": config.BoolVariable(false),
	}
	resourceName := "snowflake_grant_privileges_to_database_role.test"

	databaseRoleName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name).FullyQualifiedName()
	schemaName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig:       func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToDatabaseRole/OnSchema"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleName),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", string(sdk.SchemaPrivilegeCreateSnowflakeMlAnomalyDetection)),
					resource.TestCheckResourceAttr(resourceName, "privileges.1", string(sdk.SchemaPrivilegeCreateSnowflakeMlForecast)),
					resource.TestCheckResourceAttr(resourceName, "on_schema.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_schema.0.schema_name", schemaName),
					resource.TestCheckResourceAttr(resourceName, "with_grant_option", "false"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("%s|false|false|CREATE SNOWFLAKE.ML.ANOMALY_DETECTION,CREATE SNOWFLAKE.ML.FORECAST|OnSchema|OnSchema|%s", databaseRoleName, schemaName)),
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
func TestAcc_GrantPrivilegesToDatabaseRole_ChangeWithGrantOptionsOutsideOfTerraform_WithGrantOptions(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
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
						t, sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name),
						acc.TestDatabaseName,
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { createDatabaseRoleOutsideTerraform(t, acc.TestDatabaseName, name) },
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
						t, sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, name),
						acc.TestDatabaseName,
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					databaseCleanup = createDatabaseOutsideTerraform(t, databaseName)
					createDatabaseRoleOutsideTerraform(t, databaseName, name)
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
		CheckDestroy: testAccCheckDatabaseRolePrivilegesRevoked,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Cleanup(createDatabaseOutsideTerraform(t, databaseName))
					databaseRoleCleanup = createDatabaseRoleOutsideTerraform(t, databaseName, name)
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

func createDatabaseRoleOutsideTerraform(t *testing.T, databaseName string, name string) func() {
	t.Helper()
	client := acc.Client(t)
	ctx := context.Background()
	databaseRoleId := sdk.NewDatabaseObjectIdentifier(databaseName, name)
	if err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(databaseRoleId).WithOrReplace(true)); err != nil {
		t.Fatal(fmt.Errorf("error database role (%s): %w", databaseRoleId.FullyQualifiedName(), err))
	}

	return func() {
		if err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(databaseRoleId).WithIfExists(true)); err != nil {
			t.Fatal(fmt.Errorf("error database role (%s): %w", databaseRoleId.FullyQualifiedName(), err))
		}
	}
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

func testAccCheckDatabaseRolePrivilegesRevoked(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
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
			return fmt.Errorf("database role (%s) is still granted, granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
		}
	}
	return nil
}

func revokeAndGrantPrivilegesOnDatabaseToDatabaseRole(
	t *testing.T,
	databaseRoleName sdk.DatabaseObjectIdentifier,
	databaseName string,
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
			Database: sdk.Pointer(sdk.NewAccountObjectIdentifier(databaseName)),
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
			Database: sdk.Pointer(sdk.NewAccountObjectIdentifier(databaseName)),
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
