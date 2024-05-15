package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func createApp(t *testing.T, name string) *sdk.Application {
	t.Helper()

	stage, cleanupStage := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)

	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_GrantApplicationRole/manifest.yml")
	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_GrantApplicationRole/setup.sql")

	applicationPackage, cleanupApplicationPackage := acc.TestClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanupApplicationPackage)

	acc.TestClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), "v1")

	application, cleanupApplication := acc.TestClient().Application.CreateApplicationWithID(t, sdk.NewAccountObjectIdentifier(name), applicationPackage.ID(), "v1")
	t.Cleanup(cleanupApplication)
	return application
}

func TestAcc_GrantApplicationRole_accountRole(t *testing.T) {
	applicationName := acc.TestClient().Ids.Alpha()
	parentAccountRoleName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"parent_account_role_name": config.StringVariable(parentAccountRoleName),
			"application_name":         config.StringVariable(applicationName),
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			createApp(t, applicationName)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.TestAccCheckGrantApplicationRoleDestroy,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/account_role"),
				ConfigVariables:          m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "application_role_name", applicationRoleNameFullyQualified),
					resource.TestCheckResourceAttr(resourceName, "parent_account_role_name", fmt.Sprintf("\"%s\"", parentAccountRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ACCOUNT_ROLE|"%v"`, applicationName, applicationRoleName, parentAccountRoleName)),
				),
			},
			// test import
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/account_role"),
				ConfigVariables:          m(),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}

func TestAcc_GrantApplicationRole_application(t *testing.T) {
	applicationName := acc.TestClient().Ids.Alpha()
	applicationName2 := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"application_name":  config.StringVariable(applicationName),
			"application_name2": config.StringVariable(applicationName2),
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			createApp(t, applicationName)
			createApp(t, applicationName2)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.TestAccCheckGrantApplicationRoleDestroy,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/application"),
				ConfigVariables:          m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "application_role_name", applicationRoleNameFullyQualified),
					resource.TestCheckResourceAttr(resourceName, "application_name", fmt.Sprintf("\"%s\"", applicationName2)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|APPLICATION|"%v"`, applicationName, applicationRoleName, applicationName2)),
				),
			},
			// test import
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/application"),
				ConfigVariables:          m(),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}
