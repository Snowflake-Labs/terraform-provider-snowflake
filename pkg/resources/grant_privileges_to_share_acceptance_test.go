package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_GrantPrivilegesToShare_OnDatabase(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "database_name", databaseName.Name()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchema(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
			"schema":     config.StringVariable(schemaName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "schema_name", schemaName.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnSchema_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTable(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	tableName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
			"schema":     config.StringVariable(schemaName.Name()),
			"table_name": config.StringVariable(tableName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "table_name", tableName.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTable_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnAllTablesInSchema(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
			"schema":     config.StringVariable(schemaName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "all_tables_in_schema", schemaName.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnAllTablesInSchema_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnView(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	tableName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	viewName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
			"schema":     config.StringVariable(schemaName.Name()),
			"table_name": config.StringVariable(tableName.Name()),
			"view_name":  config.StringVariable(viewName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeSelect.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "view_name", viewName.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnView_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTag(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	tagName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
			"schema":     config.StringVariable(schemaName.Name()),
			"tag_name":   config.StringVariable(tagName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeRead.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeRead.String()),
					resource.TestCheckResourceAttr(resourceName, "tag_name", tagName.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnTag_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnPrivilegeUpdate(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool, privileges []sdk.ObjectPrivilege) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			if len(privileges) > 0 {
				configPrivileges := make([]config.Variable, len(privileges))
				for i, privilege := range privileges {
					configPrivileges[i] = config.StringVariable(privilege.String())
				}
				variables["privileges"] = config.ListVariable(configPrivileges...)
			}
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true, []sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeReferenceUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "database_name", databaseName.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true, []sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "database_name", databaseName.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true, []sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
				}),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				ConfigVariables: configVariables(false, []sdk.ObjectPrivilege{}),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDatabaseWithReferenceUsagePrivilege(t *testing.T) {
	databaseName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))
	shareName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)))

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	currentAccount, err := client.ContextFunctions.CurrentAccount(context.Background())
	require.NoError(t, err)

	shareAccountName := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(fmt.Sprintf("%s.%s", currentAccount, shareName.Name()))
	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"share_name": config.StringVariable(shareName.Name()),
			"database":   config.StringVariable(databaseName.Name()),
		}
		if withGrant {
			variables["share_account_name"] = config.StringVariable(shareAccountName.FullyQualifiedName())
			variables["privileges"] = config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeReferenceUsage.String()),
			)
		}
		return variables
	}
	resourceName := "snowflake_grant_privileges_to_share.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "share_name", shareAccountName.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "database_name", databaseName.Name()),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables:   configVariables(true),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoGrant"),
				ConfigVariables: configVariables(false),
				Check:           testAccCheckSharePrivilegesRevoked(),
			},
		},
	})
}

func testAccCheckSharePrivilegesRevoked() func(*terraform.State) error {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_share" {
				continue
			}
			db := acc.TestAccProvider.Meta().(*sql.DB)
			client := sdk.NewClientFromDB(db)
			ctx := context.Background()

			id := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["share_name"])
			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				To: &sdk.ShowGrantsTo{
					Share: sdk.NewAccountObjectIdentifier(id.Name()),
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
				return fmt.Errorf("share (%s) is still granted with privileges: %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}
