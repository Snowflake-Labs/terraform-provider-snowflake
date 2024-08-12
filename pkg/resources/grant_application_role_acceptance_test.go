package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1431726]: Move to helpers
func createApp(t *testing.T) *sdk.Application {
	t.Helper()

	stage, cleanupStage := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)

	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_GrantApplicationRole/manifest.yml")
	acc.TestClient().Stage.PutOnStage(t, stage.ID(), "TestAcc_GrantApplicationRole/setup.sql")

	applicationPackage, cleanupApplicationPackage := acc.TestClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanupApplicationPackage)

	acc.TestClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), "v1")

	application, cleanupApplication := acc.TestClient().Application.CreateApplication(t, applicationPackage.ID(), "v1")
	t.Cleanup(cleanupApplication)
	return application
}

func TestAcc_GrantApplicationRole_accountRole(t *testing.T) {
	parentRole, cleanupParentRole := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(cleanupParentRole)
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := testvars.ApplicationRole1

	acc.TestAccPreCheck(t)
	app := createApp(t)
	applicationRoleNameFullyQualified := sdk.NewDatabaseObjectIdentifier(app.Name, applicationRoleName).FullyQualifiedName()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"parent_account_role_name": config.StringVariable(parentRole.Name.Name()),
			"application_name":         config.StringVariable(app.Name),
			"application_role_name":    config.StringVariable(applicationRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "parent_account_role_name", parentRole.Name.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ACCOUNT_ROLE|"%v"`, app.Name, applicationRoleName, parentRole.Name.Name())),
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
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := testvars.ApplicationRole1

	acc.TestAccPreCheck(t)
	app := createApp(t)
	app2 := createApp(t)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"application_name":      config.StringVariable(app.Name),
			"application_name2":     config.StringVariable(app2.Name),
			"application_role_name": config.StringVariable(applicationRoleName),
		}
	}
	applicationRoleNameFullyQualified := sdk.NewDatabaseObjectIdentifier(app.Name, applicationRoleName).FullyQualifiedName()
	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "application_name", fmt.Sprintf("\"%s\"", app2.Name)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|APPLICATION|"%v"`, app.Name, applicationRoleName, app2.Name)),
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
