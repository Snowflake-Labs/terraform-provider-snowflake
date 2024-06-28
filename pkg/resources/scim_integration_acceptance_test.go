package resources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	tfjson "github.com/hashicorp/terraform-json"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ScimIntegration_basic(t *testing.T) {
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)

	role, role2 := snowflakeroles.GenericScimProvisioner, snowflakeroles.OktaProvisioner
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	m := func(enabled bool, scimClient sdk.ScimSecurityIntegrationScimClientOption, runAsRole sdk.AccountObjectIdentifier, complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":        config.StringVariable(id.Name()),
			"scim_client": config.StringVariable(string(scimClient)),
			"run_as_role": config.StringVariable(runAsRole.Name()),
			"enabled":     config.BoolVariable(enabled),
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
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(false, sdk.ScimSecurityIntegrationScimClientGeneric, role, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name()),
					resource.TestCheckNoResourceAttr("snowflake_scim_integration.test", "network_policy"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "unknown"),
					resource.TestCheckNoResourceAttr("snowflake_scim_integration.test", "comment"),

					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.enabled.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.network_policy.0.value", ""),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.run_as_role.0.value", role.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.sync_password.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.comment.0.value", ""),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(false, sdk.ScimSecurityIntegrationScimClientGeneric, role, false),
				ResourceName:    "snowflake_scim_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "scim_client", "GENERIC"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "run_as_role", role.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "network_policy", ""),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "sync_password", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", ""),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role2.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", sdk.NewAccountObjectIdentifier(networkPolicy.Name).Name()), // TODO(SNOW-999049): Fix during identifiers rework
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.integration_type", "SCIM - OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.category", "SECURITY"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test", "show_output.0.created_on"),

					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.#", "1"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.enabled.0.value", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.network_policy.0.value", sdk.NewAccountObjectIdentifier(networkPolicy.Name).Name()), // TODO(SNOW-999049): Fix during identifiers rework
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.run_as_role.0.value", role2.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.sync_password.0.value", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.comment.0.value", "foo"),
				),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, true),
				ResourceName:    "snowflake_scim_integration.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "scim_client", "OKTA"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "run_as_role", role2.Name()),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "network_policy", sdk.NewAccountObjectIdentifier(networkPolicy.Name).Name()), // TODO(SNOW-999049): Fix during identifiers rework
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "sync_password", "false"),
					importchecks.TestCheckResourceAttrInstanceState(id.Name(), "comment", "foo"),
				),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(true, sdk.ScimSecurityIntegrationScimClientOkta, role2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "true"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "OKTA"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role2.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", ""),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "unknown"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_complete(t *testing.T) {
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientGeneric)),
			"sync_password":       config.BoolVariable(false),
			"network_policy_name": config.StringVariable(networkPolicy.Name),
			"run_as_role":         config.StringVariable(role.Name()),
			"comment":             config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", sdk.NewAccountObjectIdentifier(networkPolicy.Name).Name()), // TODO(SNOW-999049): Fix during identifiers rework
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", "foo"),
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
			"network_policy_name": config.StringVariable("foo"),
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
				ExpectError: helpers.MatchAllStringsInOrderNonOverlapping([]string{
					`expected scim_client to be one of ["OKTA" "AZURE" "GENERIC"], got invalid`,
					`expected run_as_role to be one of ["OKTA_PROVISIONER" "AAD_PROVISIONER" "GENERIC_SCIM_PROVISIONER"], got invalid`,
				}),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidIncomplete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name": config.StringVariable(id.Name()),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			`The argument "scim_client" is required, but no definition was found.`,
			`The argument "run_as_role" is required, but no definition was found.`,
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/invalid"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion091(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	resourceName := "snowflake_scim_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: scimIntegrationv091(id.Name(), role.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "provisioner_role", role.Name()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationv092(id.Name(), role.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_scim_integration.test", "name", tfjson.ActionUpdate, sdk.String(id.Name()), sdk.String(id.Name())),
						planchecks.ExpectChange("snowflake_scim_integration.test", "enabled", tfjson.ActionUpdate, sdk.String("true"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "scim_client", tfjson.ActionUpdate, sdk.String("GENERIC"), sdk.String("GENERIC")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "run_as_role", tfjson.ActionUpdate, sdk.String(role.Name()), sdk.String(role.Name())),
						planchecks.ExpectChange("snowflake_scim_integration.test", "network_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "sync_password", tfjson.ActionUpdate, nil, sdk.String("unknown")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "comment", tfjson.ActionUpdate, nil, nil),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "run_as_role", role.Name()),
				),
			},
		},
	})
}

func scimIntegrationv091(name, roleName string) string {
	s := `
resource "snowflake_scim_integration" "test" {
	name             = "%s"
	scim_client      = "%s"
	provisioner_role = "%s"
}
`
	return fmt.Sprintf(s, name, sdk.ScimSecurityIntegrationScimClientGeneric, roleName)
}

func scimIntegrationv092(name, roleName string) string {
	s := `
resource "snowflake_scim_integration" "test" {
	name             = "%s"
	enabled          = true
	scim_client      = "%s"
	run_as_role		 = "%s"
}
`
	return fmt.Sprintf(s, name, sdk.ScimSecurityIntegrationScimClientGeneric, roleName)
}
