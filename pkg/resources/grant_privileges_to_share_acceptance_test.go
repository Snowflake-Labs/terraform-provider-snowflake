package resources_test

import (
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantPrivilegesToShare_OnDatabase(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName.Name()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnSchema(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), acc.TestClient().Ids.Alpha())
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
			"schema":   config.StringVariable(schemaName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_schema", schemaName.FullyQualifiedName()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

// TODO(SNOW-1021686): Add on_function test

func TestAcc_GrantPrivilegesToShare_OnTable(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), acc.TestClient().Ids.Alpha())
	tableName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), acc.TestClient().Ids.Alpha())
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
			"schema":   config.StringVariable(schemaName.Name()),
			"on_table": config.StringVariable(tableName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_table", tableName.FullyQualifiedName()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnAllTablesInSchema(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), acc.TestClient().Ids.Alpha())
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
			"schema":   config.StringVariable(schemaName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_all_tables_in_schema", schemaName.FullyQualifiedName()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnView(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), acc.TestClient().Ids.Alpha())
	tableName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), acc.TestClient().Ids.Alpha())
	viewName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), acc.TestClient().Ids.Alpha())
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
			"schema":   config.StringVariable(schemaName.Name()),
			"on_table": config.StringVariable(tableName.Name()),
			"on_view":  config.StringVariable(viewName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeSelect.String()),
					resource.TestCheckResourceAttr(resourceName, "on_view", viewName.FullyQualifiedName()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnTag(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := sdk.NewDatabaseObjectIdentifier(databaseName.Name(), acc.TestClient().Ids.Alpha())
	tagName := sdk.NewSchemaObjectIdentifier(databaseName.Name(), schemaName.Name(), acc.TestClient().Ids.Alpha())
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
			"schema":   config.StringVariable(schemaName.Name()),
			"on_tag":   config.StringVariable(tagName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeRead.String()),
					resource.TestCheckResourceAttr(resourceName, "on_tag", tagName.FullyQualifiedName()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnPrivilegeUpdate(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool, privileges []sdk.ObjectPrivilege) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase"),
				ConfigVariables: configVariables(true, []sdk.ObjectPrivilege{
					sdk.ObjectPrivilegeUsage,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName.Name()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_OnDatabaseWithReferenceUsagePrivilege(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func(withGrant bool) config.Variables {
		variables := config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
		}
		if withGrant {
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
					resource.TestCheckResourceAttr(resourceName, "to_share", shareName.Name()),
					resource.TestCheckResourceAttr(resourceName, "privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", sdk.ObjectPrivilegeReferenceUsage.String()),
					resource.TestCheckResourceAttr(resourceName, "on_database", databaseName.Name()),
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
				Check:           acc.CheckSharePrivilegesRevoked(t),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoPrivileges(t *testing.T) {
	databaseName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func() config.Variables {
		return config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"database": config.StringVariable(databaseName.Name()),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnDatabase_NoPrivileges"),
				ConfigVariables: configVariables(),
				ExpectError:     regexp.MustCompile(`The argument "privileges" is required, but no definition was found.`),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToShare_NoOnOption(t *testing.T) {
	shareName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func() config.Variables {
		return config.Variables{
			"to_share": config.StringVariable(shareName.Name()),
			"privileges": config.ListVariable(
				config.StringVariable(sdk.ObjectPrivilegeReferenceUsage.String()),
			),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/NoOnOption"),
				ConfigVariables: configVariables(),
				ExpectError:     regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2621 doesn't apply to this resource
func TestAcc_GrantPrivilegesToShare_RemoveShareOutsideTerraform(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	shareName := acc.TestClient().Ids.Alpha()

	configVariables := config.Variables{
		"to_share": config.StringVariable(shareName),
		"database": config.StringVariable(databaseName),
		"privileges": config.ListVariable(
			config.StringVariable(sdk.ObjectPrivilegeUsage.String()),
		),
	}

	var shareCleanup func()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, shareCleanup = acc.TestClient().Share.CreateShareWithName(t, shareName)
					t.Cleanup(shareCleanup)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnCustomShare"),
				ConfigVariables: configVariables,
			},
			{
				PreConfig:       func() { shareCleanup() },
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantPrivilegesToShare/OnCustomShare"),
				ConfigVariables: configVariables,
				// The error occurs in the Create operation, indicating the Read operation removed the resource from the state in the previous step.
				ExpectError: regexp.MustCompile("An error occurred when granting privileges to share"),
			},
		},
	})
}
