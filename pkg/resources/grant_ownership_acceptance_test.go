package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnObject_Database_ToAccountRole(t *testing.T) {
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseName),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(5, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(5, acctest.CharSetAlpha) + "." + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeDatabase, accountRoleName, databaseName),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeSchema, accountRoleName, fmt.Sprintf("%s.%s", databaseName, schemaName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
					}, sdk.ObjectTypeSchema, databaseRoleName, fmt.Sprintf("%s.%s", databaseName, schemaName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf("%s.%s.%s", databaseName, schemaName, tableName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableFullyQualifiedName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName).FullyQualifiedName()

	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
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
					}, sdk.ObjectTypeTable, databaseRoleName, fmt.Sprintf("%s.%s.%s", databaseName, schemaName, tableName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf("%s.%s.%s", databaseName, schemaName, tableName), fmt.Sprintf("%s.%s.%s", databaseName, schemaName, secondTableName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf("%s.%s.%s", databaseName, schemaName, tableName), fmt.Sprintf("%s.%s.%s", databaseName, schemaName, secondTableName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
							Database: sdk.Pointer(sdk.NewAccountObjectIdentifier(databaseName)),
						},
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf("%s.<TABLE>", databaseName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaFullyQualifiedName := sdk.NewDatabaseObjectIdentifier(databaseName, schemaName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
					}, sdk.ObjectTypeTable, accountRoleName, fmt.Sprintf("%s.%s.<TABLE>", databaseName, schemaName)),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))),
		"database_name":     config.StringVariable(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))),
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
	t.Skip("will be unskipped in the following grant ownership prs")

	configVariables := config.Variables{
		"account_role_name": config.StringVariable(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))),
		"database_name":     config.StringVariable(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))),
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

func TestAcc_GrantOwnership_MoveOwnershipOutsideTerraform(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

	otherAccountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

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
					moveResourceOwnershipToAccountRole(t, sdk.ObjectTypeDatabase, sdk.NewAccountObjectIdentifier(databaseName), sdk.NewAccountObjectIdentifier(otherAccountRoleName))
				},
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
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/MoveResourceOwnershipOutsideTerraform"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_ForceOwnershipTransferOnCreate(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()

	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()

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
				PreConfig: func() {
					t.Cleanup(createAccountRole(t, accountRoleName))
					t.Cleanup(createDatabaseWithRoleAsOwner(t, accountRoleName, databaseName))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/ForceOwnershipTransferOnCreate"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/ForceOwnershipTransferOnCreate"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnPipe(t *testing.T) {
	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	stageName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	pipeName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()
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
						// TODO(): Fix this identifier
					}, sdk.ObjectTypePipe, accountRoleName, fmt.Sprintf("%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, pipeName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnPipe"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantOwnership_OnAllPipes(t *testing.T) {
	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	stageName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	pipeName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondPipeName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()
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
							Role: sdk.NewAccountObjectIdentifier(accountRoleName),
						},
						// TODO: Fix this identifier
					}, sdk.ObjectTypePipe, accountRoleName, fmt.Sprintf("%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, pipeName), fmt.Sprintf("%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, secondPipeName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAllPipes"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createDatabaseWithRoleAsOwner(t *testing.T, roleName string, databaseName string) func() {
	t.Helper()
	client, err := sdk.NewDefaultClient()
	assert.NoError(t, err)

	ctx := context.Background()
	databaseId := sdk.NewAccountObjectIdentifier(databaseName)
	assert.NoError(t, client.Databases.Create(ctx, databaseId, &sdk.CreateDatabaseOptions{}))

	err = client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       databaseId,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: sdk.Pointer(sdk.NewAccountObjectIdentifier(roleName)),
		},
		new(sdk.GrantOwnershipOptions),
	)
	assert.NoError(t, err)

	return func() {
		assert.NoError(t, client.Databases.Drop(ctx, databaseId, &sdk.DropDatabaseOptions{}))
	}
}

func createAccountRole(t *testing.T, name string) func() {
	t.Helper()
	client, err := sdk.NewDefaultClient()
	assert.NoError(t, err)

	ctx := context.Background()
	roleId := sdk.NewAccountObjectIdentifier(name)
	assert.NoError(t, client.Roles.Create(ctx, sdk.NewCreateRoleRequest(roleId)))

	return func() {
		assert.NoError(t, client.Roles.Drop(ctx, sdk.NewDropRoleRequest(roleId)))
	}
}

func moveResourceOwnershipToAccountRole(t *testing.T, objectType sdk.ObjectType, objectName sdk.ObjectIdentifier, accountRoleName sdk.AccountObjectIdentifier) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Grants.GrantOwnership(
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
				slices.Contains(objectNames, grant.Name.Name()) {
				found = append(found, grant.Name.Name())
			}
		}

		if len(found) != len(objectNames) {
			return fmt.Errorf("unable to find ownership privilege on %s granted to %s, expected names: %v, found: %v", grantOn, roleName, objectNames, found)
		}

		return nil
	}
}
