package resources_test

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"slices"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnObject_Database_ToAccountRole(t *testing.T) {
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

func TestAcc_GrantOwnership_OnObject_Schema_ToAccountRole(t *testing.T) {
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
