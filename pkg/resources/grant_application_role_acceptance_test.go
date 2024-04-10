package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantApplicationRole_accountRole(t *testing.T) {
	applicationName := "int_test_app"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)
	parentRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_application_role.g"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":             config.StringVariable(applicationRoleNameFullyQualified),
			"parent_role_name": config.StringVariable(parentRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantApplicationRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", applicationRoleNameFullyQualified),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", parentRoleName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ROLE|"%v"`, applicationName, applicationRoleName, parentRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/account_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantApplicationRole_application(t *testing.T) {
	applicationName := "int_test_app"
	applicationRoleName := "app_role_1"
	applicationRoleNameFullyQualified := fmt.Sprintf("\"%s\".\"%s\"", applicationName, applicationRoleName)
	secondApplicationName := "int_test_other_app"
	resourceName := "snowflake_grant_application_role.g"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":             config.StringVariable(applicationRoleNameFullyQualified),
			"application_name": config.StringVariable(secondApplicationName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantApplicationRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/application"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", applicationRoleNameFullyQualified),
					resource.TestCheckResourceAttr(resourceName, "application_name", secondApplicationName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|APPLICATION|"%v"`, applicationName, applicationRoleName, secondApplicationName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantApplicationRole/application"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGrantApplicationRoleDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_grant_application_role" {
			continue
		}
		ctx := context.Background()
		id := rs.Primary.ID
		ids := strings.Split(id, "|")
		applicationRoleName := ids[0]
		objectType := ids[1]
		parentRoleName := ids[2]
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				ApplicationRole: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(applicationRoleName),
			},
		})
		if err != nil {
			continue
		}
		for _, grant := range grants {
			if grant.GrantedTo == sdk.ObjectType(objectType) {
				if grant.GranteeName.FullyQualifiedName() == parentRoleName {
					return fmt.Errorf("application role grant %v still exists", grant)
				}
			}
		}
	}
	return nil
}
