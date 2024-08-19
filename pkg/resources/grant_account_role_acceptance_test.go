package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantAccountRole_accountRole(t *testing.T) {
	roleName := acc.TestClient().Ids.Alpha()
	parentRoleName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_account_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"role_name":        config.StringVariable(roleName),
			"parent_role_name": config.StringVariable(parentRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantAccountRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|ROLE|"%v"`, roleName, parentRoleName)),
				),
			},
			// import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantAccountRole/account_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantAccountRole_user(t *testing.T) {
	roleName := acc.TestClient().Ids.Alpha()
	userName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_account_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"role_name": config.StringVariable(roleName),
			"user_name": config.StringVariable(userName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckGrantAccountRoleDestroy(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantAccountRole/user"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"|USER|"%v"`, roleName, userName)),
				),
			},
			// import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantAccountRole/user"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantAccountRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	roleName := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	parentRoleName := acc.TestClient().Ids.RandomAccountObjectIdentifier()

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
				Config: grantAccountRoleBasicConfig(roleName.Name(), parentRoleName.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleName.FullyQualifiedName(), parentRoleName.FullyQualifiedName())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantAccountRoleBasicConfig(roleName.Name(), parentRoleName.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleName.FullyQualifiedName(), parentRoleName.FullyQualifiedName())),
				),
			},
		},
	})
}

func grantAccountRoleBasicConfig(roleName string, parentRoleName string) string {
	return fmt.Sprintf(`
resource "snowflake_account_role" "role" {
  name = "%s"
}

resource "snowflake_account_role" "parent_role" {
  name = "%s"
}

resource "snowflake_grant_account_role" "test" {
  role_name        = snowflake_account_role.role.name
  parent_role_name = snowflake_account_role.parent_role.name
}
`, roleName, parentRoleName)
}

func TestAcc_GrantAccountRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedRoleId := fmt.Sprintf(`\"%s\"`, roleId.Name())

	parentRoleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedParentRoleId := fmt.Sprintf(`\"%s\"`, parentRoleId.Name())

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
				ExpectError: regexp.MustCompile("Error: Provider produced inconsistent final plan"),
				Config:      grantAccountRoleBasicConfig(quotedRoleId, quotedParentRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "role_name", roleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "parent_role_name", parentRoleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   grantAccountRoleBasicConfig(quotedRoleId, quotedParentRoleId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "role_name", roleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "parent_role_name", parentRoleId.Name()),
					resource.TestCheckResourceAttr("snowflake_grant_account_role.test", "id", fmt.Sprintf(`%v|ROLE|%v`, roleId.FullyQualifiedName(), parentRoleId.FullyQualifiedName())),
				),
			},
		},
	})
}
