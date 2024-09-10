package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantDatabaseRole_databaseRole(t *testing.T) {
	databaseRoleName := acc.TestClient().Ids.Alpha()
	parentDatabaseRoleName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|DATABASE ROLE|"%v"."%v"`, acc.TestDatabaseName, databaseRoleName, acc.TestDatabaseName, parentDatabaseRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_databaseRoleMixedQuoting(t *testing.T) {
	databaseRoleName := acc.TestClient().Ids.Alpha()
	parentDatabaseRoleName := strings.ToUpper(acc.TestClient().Ids.Alpha())
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|DATABASE ROLE|"%v"."%v"`, acc.TestDatabaseName, databaseRoleName, acc.TestDatabaseName, parentDatabaseRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_issue2402(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	databaseRoleName := acc.TestClient().Ids.Alpha()
	parentDatabaseRoleName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(databaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantDatabaseRole/issue2402"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|DATABASE ROLE|"%v"."%v"`, databaseName, databaseRoleName, databaseName, parentDatabaseRoleName)),
				),
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_accountRole(t *testing.T) {
	databaseRoleName := acc.TestClient().Ids.Alpha()
	parentRoleName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":           config.StringVariable(acc.TestDatabaseName),
			"database_role_name": config.StringVariable(databaseRoleName),
			"parent_role_name":   config.StringVariable(parentRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", fmt.Sprintf("%v", parentRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ROLE|"%v"`, acc.TestDatabaseName, databaseRoleName, parentRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2410 is fixed
func TestAcc_GrantDatabaseRole_share(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	shareId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	configVariables := func() config.Variables {
		return config.Variables{
			"database":           config.StringVariable(databaseId.Name()),
			"database_role_name": config.StringVariable(databaseRoleId.Name()),
			"share_name":         config.StringVariable(shareId.Name()),
		}
	}
	resourceName := "snowflake_grant_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables: configVariables(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "share_name", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`%v|%v|%v`, databaseRoleId.FullyQualifiedName(), "SHARE", shareId.FullyQualifiedName())),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables:   configVariables(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_shareWithDots(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	shareId := acc.TestClient().Ids.RandomAccountObjectIdentifierContaining(".")

	configVariables := func() config.Variables {
		return config.Variables{
			"database":           config.StringVariable(databaseId.Name()),
			"database_role_name": config.StringVariable(databaseRoleId.Name()),
			"share_name":         config.StringVariable(shareId.Name()),
		}
	}
	resourceName := "snowflake_grant_database_role.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckGrantDatabaseRoleDestroy(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables: configVariables(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "share_name", shareId.Name()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`%v|%v|%v`, databaseRoleId.FullyQualifiedName(), "SHARE", shareId.FullyQualifiedName())),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables:   configVariables(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	parentRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: grantDatabaseRoleBasicConfigQuoted(databaseId.Name(), roleId.Name(), parentRoleId.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantDatabaseRoleBasicConfigQuoted(databaseId.Name(), roleId.Name(), parentRoleId.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantDatabaseRoleBasicConfigQuoted(databaseName string, roleName string, parentRoleName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%[1]s"
}

resource "snowflake_database_role" "role" {
  database = snowflake_database.test.name
  name = "%[2]s"
}

resource "snowflake_database_role" "parent_role" {
  database = snowflake_database.test.name
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name        = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.role.name}\""
  parent_database_role_name = "\"${snowflake_database.test.name}\".\"${snowflake_database_role.parent_role.name}\""
}
`, databaseName, roleName, parentRoleName)
}

func grantDatabaseRoleBasicConfigUnquoted(databaseName string, roleName string, parentRoleName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%[1]s"
}

resource "snowflake_database_role" "role" {
  database = snowflake_database.test.name
  name = "%[2]s"
}

resource "snowflake_database_role" "parent_role" {
  database = snowflake_database.test.name
  name = "%[3]s"
}

resource "snowflake_grant_database_role" "test" {
  database_role_name        = "${snowflake_database.test.name}.${snowflake_database_role.role.name}"
  parent_database_role_name = "${snowflake_database.test.name}.${snowflake_database_role.parent_role.name}"
}
`, databaseName, roleName, parentRoleName)
}

func TestAcc_GrantDatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	roleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	parentRoleId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: grantDatabaseRoleBasicConfigUnquoted(databaseId.Name(), roleId.Name(), parentRoleId.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "database_role_name", fmt.Sprintf("%s.%s", roleId.DatabaseName(), roleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "parent_database_role_name", fmt.Sprintf("%s.%s", roleId.DatabaseName(), parentRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantDatabaseRoleBasicConfigUnquoted(databaseId.Name(), roleId.Name(), parentRoleId.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "database_role_name", fmt.Sprintf("%s.%s", roleId.DatabaseName(), roleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "parent_database_role_name", fmt.Sprintf("%s.%s", roleId.DatabaseName(), parentRoleId.Name())),
					resource.TestCheckResourceAttr("snowflake_grant_database_role.test", "id", fmt.Sprintf(`%s|DATABASE ROLE|%s`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}
