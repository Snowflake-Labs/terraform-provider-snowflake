package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnObject_Database_ToAccountRole(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Database_IdentifiersWithDots(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Schema_ToAccountRole(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()

	schemaName := acc.TestClient().Ids.Alpha()
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "SCHEMA"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|SCHEMA|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeSchema, accountRoleName, schemaFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Schema_ToDatabaseRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	databaseRoleName := acc.TestClient().Ids.Alpha()
	databaseRoleFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName).FullyQualifiedName()

	configVariables := config.Variables{
		"database_role_name": config.StringVariable(databaseRoleName),
		"database_name":      config.StringVariable(databaseName),
		"schema_name":        config.StringVariable(schemaName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToDatabaseRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "SCHEMA"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToDatabaseRole|%s||OnObject|SCHEMA|%s", databaseRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							DatabaseRole: sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName),
						},
					}, sdk.ObjectTypeSchema, databaseRoleName, schemaFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Schema_ToDatabaseRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Table_ToAccountRole(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	schemaName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
		"table_name":        config.StringVariable(tableName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", tableFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TABLE|%s", accountRoleFullyQualifiedName, tableFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleName, tableFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnObject_Table_ToDatabaseRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()

	databaseRoleName := acc.TestClient().Ids.Alpha()
	databaseRoleFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName).FullyQualifiedName()

	configVariables := config.Variables{
		"database_role_name": config.StringVariable(databaseRoleName),
		"database_name":      config.StringVariable(databaseName),
		"schema_name":        config.StringVariable(schemaName),
		"table_name":         config.StringVariable(tableName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToDatabaseRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", tableFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToDatabaseRole|%s||OnObject|TABLE|%s", databaseRoleFullyQualifiedName, tableFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							DatabaseRole: sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName),
						},
					}, sdk.ObjectTypeTable, databaseRoleName, tableFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Table_ToDatabaseRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAll_InDatabase_ToAccountRole(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	schemaName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	secondTableName := acc.TestClient().Ids.Alpha()
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()
	secondTableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, secondTableName).FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
		"table_name":        config.StringVariable(tableName),
		"second_table_name": config.StringVariable(secondTableName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InDatabase_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|TABLES|InDatabase|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleName, tableFullyQualifiedName, secondTableFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InDatabase_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAll_InSchema_ToAccountRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	tableName := acc.TestClient().Ids.Alpha()
	secondTableName := acc.TestClient().Ids.Alpha()
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()
	secondTableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, secondTableName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
		"table_name":        config.StringVariable(tableName),
		"second_table_name": config.StringVariable(secondTableName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InSchema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.in_schema", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|TABLES|InSchema|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeTable, accountRoleName, tableFullyQualifiedName, secondTableFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll_InSchema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnFuture_InDatabase_ToAccountRole(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InDatabase_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnFuture|TABLES|InDatabase|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						Future: sdk.Bool(true),
						In: &sdk.ShowGrantsIn{
							Database: sdk.Pointer(databaseId),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf(`"%s"."<TABLE>"`, databaseName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InDatabase_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnFuture_InSchema_ToAccountRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
		"schema_name":       config.StringVariable(schemaName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InSchema_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.in_schema", schemaFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnFuture|TABLES|InSchema|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						Future: sdk.Bool(true),
						In: &sdk.ShowGrantsIn{
							Schema: sdk.Pointer(sdk.NewDatabaseObjectIdentifier(databaseName, schemaName)),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf(`"%s"."%s"."<TABLE>"`, databaseName, schemaName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture_InSchema_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_InvalidConfiguration_EmptyObjectType(t *testing.T) {
	configVariables := config.Variables{
		"account_role_name": config.StringVariable(acc.TestClient().Ids.Alpha()),
		"database_name":     config.StringVariable(acc.TestClient().Ids.Alpha()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/InvalidConfiguration_EmptyObjectType"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("expected on.0.object_type to be one of"),
			},
		},
	})
}

func TestAcc_GrantOwnership_InvalidConfiguration_MultipleTargets(t *testing.T) {
	configVariables := config.Variables{
		"account_role_name": config.StringVariable(acc.TestClient().Ids.Alpha()),
		"database_name":     config.StringVariable(acc.TestClient().Ids.Alpha()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/InvalidConfiguration_MultipleTargets"),
				ConfigVariables: configVariables,
				ExpectError:     regexp.MustCompile("only one of `on.0.all,on.0.future,on.0.object_name`"),
			},
		},
	})
}

func TestAcc_GrantOwnership_TargetObjectRemovedOutsideTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"
	var cleanupDatabase func()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanupDatabase = acc.TestClient().Database.CreateDatabaseWithName(t, databaseName)
					t.Cleanup(cleanupDatabase)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoDatabaseResource"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				PreConfig: func() {
					currentRole := acc.TestClient().Context.CurrentRole(t)
					acc.TestClient().Role.GrantOwnershipOnAccountObject(t, currentRole, databaseId, sdk.ObjectTypeDatabase)
					cleanupDatabase()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoDatabaseResource"),
				ConfigVariables: configVariables,
				// The error occurs in Create operation indicating the Read operation couldn't find the grant and set the resource as removed.
				ExpectError: regexp.MustCompile("An error occurred during grant ownership"),
			},
		},
	})
}

func TestAcc_GrantOwnership_AccountRoleRemovedOutsideTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	_, cleanupAccountRole := acc.TestClient().Role.CreateRoleWithName(t, accountRoleName)
	t.Cleanup(cleanupAccountRole)

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoRoleResource"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
			{
				PreConfig: func() {
					cleanupAccountRole()
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_Database_ToAccountRole_NoRoleResource"),
				ConfigVariables: configVariables,
				// The error occurs in Create operation indicating the Read operation couldn't find the grant and set the resource as removed.
				ExpectError: regexp.MustCompile("An error occurred during grant ownership"),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnMaterializedView(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	materializedViewName := acc.TestClient().Ids.Alpha()
	materializedViewFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, materializedViewName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name":      config.StringVariable(accountRoleName),
		"database_name":          config.StringVariable(databaseName),
		"schema_name":            config.StringVariable(schemaName),
		"table_name":             config.StringVariable(tableName),
		"materialized_view_name": config.StringVariable(materializedViewName),
		"warehouse_name":         config.StringVariable(acc.TestWarehouseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_MaterializedView_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "MATERIALIZED VIEW"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", materializedViewFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|MATERIALIZED VIEW|%s", accountRoleFullyQualifiedName, materializedViewFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypeMaterializedView, accountRoleName, materializedViewFullyQualifiedName),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_MaterializedView_ToAccountRole"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_RoleBasedAccessControlUseCase(t *testing.T) {
	t.Skip("Will be un-skipped in SNOW-1313849")

	accountRoleName := acc.TestClient().Ids.Alpha()
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	userId := acc.TestClient().Context.CurrentUser(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// We have to make it in two steps, because provider blocks cannot contain depends_on meta-argument
			// that are needed to grant the role to the current user before it can be used.
			// Additionally, only the Config field can specify a configuration with custom provider blocks.
			{
				Config: roleBasedAccessControlUseCaseConfig(accountRoleName, databaseName, userId.Name(), schemaName, false),
			},
			{
				Config: roleBasedAccessControlUseCaseConfig(accountRoleName, databaseName, userId.Name(), schemaName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func roleBasedAccessControlUseCaseConfig(accountRoleName string, databaseName string, userName string, schemaName string, withSecondaryProvider bool) string {
	baseConfig := fmt.Sprintf(`
resource "snowflake_role" "test" {
  name = "%[1]s"
}

resource "snowflake_database" "test" {
  name = "%[2]s"
}

resource "snowflake_grant_ownership" "test" {
  account_role_name = snowflake_role.test.name
  on {
    object_type = "DATABASE"
    object_name = snowflake_database.test.name
  }
}

resource "snowflake_grant_account_role" "test" {
  role_name = snowflake_role.test.name
  user_name = "%[3]s"
}
`, accountRoleName, databaseName, userName)

	secondaryProviderConfig := fmt.Sprintf(`
provider "snowflake" {
  profile = "default"
  alias = "secondary"
  role = snowflake_role.test.name
}

resource "snowflake_schema" "test" {
  depends_on = [snowflake_grant_ownership.test, snowflake_grant_account_role.test]
  provider = snowflake.secondary
  database = snowflake_database.test.name
  name     = "%s"
}
`, schemaName)

	if withSecondaryProvider {
		return fmt.Sprintf("%s\n%s", baseConfig, secondaryProviderConfig)
	}

	return baseConfig
}

func TestAcc_GrantOwnership_MoveOwnershipOutsideTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	otherAccountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	otherAccountRoleName := otherAccountRoleId.Name()

	configVariables := config.Variables{
		"account_role_name":       config.StringVariable(accountRoleName),
		"other_account_role_name": config.StringVariable(otherAccountRoleName),
		"database_name":           config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/MoveResourceOwnershipOutsideTerraform"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				PreConfig: func() {
					moveResourceOwnershipToAccountRole(t, sdk.ObjectTypeDatabase, databaseId, otherAccountRoleId)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/MoveResourceOwnershipOutsideTerraform"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeDatabase,
								Name:       databaseId,
							},
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_ForceOwnershipTransferOnCreate(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	databaseFullyQualifiedName := databaseId.FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	newDatabaseOwningAccountRoleName := acc.TestClient().Ids.Alpha()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(newDatabaseOwningAccountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					role, roleCleanup := acc.TestClient().Role.CreateRoleWithName(t, accountRoleName)
					t.Cleanup(roleCleanup)
					_, newRoleCleanup := acc.TestClient().Role.CreateRoleWithName(t, newDatabaseOwningAccountRoleName)
					t.Cleanup(newRoleCleanup)
					database, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithName(t, databaseName)
					t.Cleanup(databaseCleanup)
					acc.TestClient().Role.GrantOwnershipOnAccountObject(t, role.ID(), database.ID(), sdk.ObjectTypeDatabase)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/ForceOwnershipTransferOnCreate"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", newDatabaseOwningAccountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|\"%s\"||OnObject|DATABASE|%s", newDatabaseOwningAccountRoleName, databaseFullyQualifiedName)),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnPipe(t *testing.T) {
	stageName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	pipeName := acc.TestClient().Ids.Alpha()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()
	pipeFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, pipeName).FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"stage":             config.StringVariable(stageName),
		"table":             config.StringVariable(tableName),
		"pipe":              config.StringVariable(pipeName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnPipe"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", sdk.ObjectTypePipe.String()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", pipeFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|PIPE|%s", accountRoleFullyQualifiedName, pipeFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypePipe,
								Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(pipeFullyQualifiedName),
							},
						},
					}, sdk.ObjectTypePipe, accountRoleName, pipeFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAllPipes(t *testing.T) {
	stageName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
	pipeName := acc.TestClient().Ids.Alpha()
	secondPipeName := acc.TestClient().Ids.Alpha()
	pipeFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, pipeName).FullyQualifiedName()
	secondPipeFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, secondPipeName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName).FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"stage":             config.StringVariable(stageName),
		"table":             config.StringVariable(tableName),
		"pipe":              config.StringVariable(pipeName),
		"second_pipe":       config.StringVariable(secondPipeName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAllPipes"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnAll|PIPES|InSchema|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					}, sdk.ObjectTypePipe, accountRoleName, pipeFullyQualifiedName, secondPipeFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnTask(t *testing.T) {
	taskName := acc.TestClient().Ids.Alpha()
	taskFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, taskName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"task":              config.StringVariable(taskName),
		"warehouse":         config.StringVariable(acc.TestWarehouseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnTask"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", sdk.ObjectTypeTask.String()),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", taskFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|TASK|%s", accountRoleFullyQualifiedName, taskFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeTask,
								Name:       sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(taskFullyQualifiedName),
							},
						},
					}, sdk.ObjectTypeTask, accountRoleName, taskFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAllTasks(t *testing.T) {
	taskName := acc.TestClient().Ids.Alpha()
	secondTaskName := acc.TestClient().Ids.Alpha()
	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()
	schemaFullyQualifiedName := acc.TestClient().Ids.SchemaId().FullyQualifiedName()
	taskFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, taskName).FullyQualifiedName()
	secondTaskFullyQualifiedName := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, secondTaskName).FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database":          config.StringVariable(acc.TestDatabaseName),
		"schema":            config.StringVariable(acc.TestSchemaName),
		"task":              config.StringVariable(taskName),
		"second_task":       config.StringVariable(secondTaskName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAllTasks"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s|REVOKE|OnAll|TASKS|InSchema|%s", accountRoleFullyQualifiedName, schemaFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						To: &sdk.ShowGrantsTo{
							Role: accountRoleId,
						},
					},
						sdk.ObjectTypeTask, accountRoleName, taskFullyQualifiedName, secondTaskFullyQualifiedName),
				),
			},
		},
	})
}

func TestAcc_GrantOwnership_OnDatabaseRole(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()

	databaseRoleName := acc.TestClient().Ids.Alpha()
	databaseRoleFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName).FullyQualifiedName()

	accountRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountRoleName := accountRoleId.Name()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	configVariables := config.Variables{
		"account_role_name":  config.StringVariable(accountRoleName),
		"database_name":      config.StringVariable(databaseName),
		"database_role_name": config.StringVariable(databaseRoleName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject_DatabaseRole_ToAccountRole"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE ROLE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseRoleFullyQualifiedName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE ROLE|%s", accountRoleFullyQualifiedName, databaseRoleFullyQualifiedName)),
					checkResourceOwnershipIsGranted(&sdk.ShowGrantOptions{
						On: &sdk.ShowGrantsOn{
							Object: &sdk.Object{
								ObjectType: sdk.ObjectTypeDatabaseRole,
								Name:       sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleFullyQualifiedName),
							},
						},
					}, sdk.ObjectTypeRole, accountRoleName, databaseRoleFullyQualifiedName),
				),
			},
		},
	})
}

func moveResourceOwnershipToAccountRole(t *testing.T, objectType sdk.ObjectType, objectName sdk.ObjectIdentifier, accountRoleName sdk.AccountObjectIdentifier) {
	t.Helper()

	client := acc.Client(t)
	ctx := context.Background()
	err := client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectName,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: &accountRoleName,
		},
		new(sdk.GrantOwnershipOptions),
	)
	assert.NoError(t, err)
}

func checkResourceOwnershipIsGranted(opts *sdk.ShowGrantOptions, grantOn sdk.ObjectType, roleName string, objectNames ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		grants, err := client.Grants.Show(ctx, opts)
		if err != nil {
			return err
		}

		found := make([]string, 0)
		for _, grant := range grants {
			if grant.Privilege == "OWNERSHIP" &&
				(grant.GrantedOn == grantOn || grant.GrantOn == grantOn) &&
				grant.GranteeName.Name() == roleName &&
				slices.Contains(objectNames, grant.Name.FullyQualifiedName()) {
				found = append(found, grant.Name.FullyQualifiedName())
			}
		}

		if len(found) != len(objectNames) {
			return fmt.Errorf("unable to find ownership privilege on %s granted to %s, expected names: %v, found: %v", grantOn, roleName, objectNames, found)
		}

		return nil
	}
}
