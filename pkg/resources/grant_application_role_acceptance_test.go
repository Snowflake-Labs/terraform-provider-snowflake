package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantApplicationRole_accountRole(t *testing.T) {
	applicationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	parentAccountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)
	randomName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database_name":            config.StringVariable(acc.TestDatabaseName),
			"schema_name":              config.StringVariable(acc.TestSchemaName),
			"parent_account_role_name": config.StringVariable(parentAccountRoleName),
			"application_name":         config.StringVariable(applicationName),
			"application_package_name": config.StringVariable(randomName),
			"stage_name":               config.StringVariable(randomName),
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
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
	applicationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	applicationName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_application_role.g"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)
	randomName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database_name":            config.StringVariable(acc.TestDatabaseName),
			"schema_name":              config.StringVariable(acc.TestSchemaName),
			"application_name":         config.StringVariable(applicationName),
			"application_name2":        config.StringVariable(applicationName2),
			"application_package_name": config.StringVariable(randomName),
			"stage_name":               config.StringVariable(randomName),
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
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
