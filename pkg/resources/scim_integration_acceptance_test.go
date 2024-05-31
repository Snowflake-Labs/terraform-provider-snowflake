package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ScimIntegration_basic(t *testing.T) {
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	role, roleCleanup := acc.TestClient().Role.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(snowflakeroles.GenericScimProvisioner).WithOrReplace(true))
	t.Cleanup(roleCleanup)
	acc.TestClient().Role.GrantRoleToCurrentRole(t, role.ID())
	role2, role2Cleanup := acc.TestClient().Role.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(snowflakeroles.OktaProvisioner).WithOrReplace(true))
	t.Cleanup(role2Cleanup)
	acc.TestClient().Role.GrantRoleToCurrentRole(t, role2.ID())

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func(enabled bool, scimClient sdk.ScimSecurityIntegrationScimClientOption, runAsRole *sdk.Role, complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":        config.StringVariable(id.Name()),
			"enabled":     config.BoolVariable(enabled),
			"scim_client": config.StringVariable(string(scimClient)),
			"run_as_role": config.StringVariable(runAsRole.ID().Name()),
		}
		if complete {
			c["sync_password"] = config.BoolVariable(false)
			c["network_policy_name"] = config.StringVariable(networkPolicy.Name)
			c["comment"] = config.StringVariable("foo")
		}
		return c
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(false, sdk.ScimSecurityIntegrationScimClientGeneric, role, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role2.Name),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", networkPolicy.Name),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables:   m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, true),
				ResourceName:      "snowflake_scim_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role2.Name),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", ""),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "created_on"),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_complete(t *testing.T) {
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)
	role, roleCleanup := acc.TestClient().Role.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(snowflakeroles.GenericScimProvisioner).WithOrReplace(true))
	t.Cleanup(roleCleanup)
	acc.TestClient().Role.GrantRoleToCurrentRole(t, role.ID())
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientGeneric)),
			"sync_password":       config.BoolVariable(false),
			"network_policy_name": config.StringVariable(networkPolicy.Name),
			"run_as_role":         config.StringVariable(role.ID().Name()),
			"comment":             config.StringVariable("foo"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", networkPolicy.Name),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "created_on"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_scim_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ScimIntegration_invalid(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable("invalid"),
			"sync_password":       config.BoolVariable(false),
			"network_policy_name": config.StringVariable("invalid"),
			"run_as_role":         config.StringVariable("invalid"),
			"comment":             config.StringVariable("foo"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(),
				ExpectError: helpers.MatchStringsRegex([]string{
					`expected scim_client to be one of ["OKTA" "AZURE" "GENERIC"], got invalid`,
					`expected run_as_role to be one of ["OKTA_PROVISIONER" "AAD_PROVISIONER" "GENERIC_SCIM_PROVISIONER"], got invalid`,
				}),
			},
		},
	})
}
