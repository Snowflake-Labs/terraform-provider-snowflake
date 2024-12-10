package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
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
			c["network_policy_name"] = config.StringVariable(networkPolicy.ID().Name())
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", ""),

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
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "scim_client", "GENERIC"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "run_as_role", role.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", ""),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "sync_password", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", ""),
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", networkPolicy.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "describe_output.0.network_policy.0.value", networkPolicy.ID().Name()),
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
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "fully_qualified_name", id.FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "enabled", "true"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "scim_client", "OKTA"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "run_as_role", role2.Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "network_policy", networkPolicy.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "sync_password", "false"),
					importchecks.TestCheckResourceAttrInstanceState(resourcehelpers.EncodeResourceIdentifier(id), "comment", "foo"),
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", r.BooleanDefault),
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
			"network_policy_name": config.StringVariable(networkPolicy.ID().Name()),
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", "GENERIC"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", networkPolicy.ID().Name()),
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

func TestAcc_ScimIntegration_completeAzure(t *testing.T) {
	networkPolicy, networkPolicyCleanup := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(networkPolicyCleanup)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientAzure)),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/completeAzure"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "enabled", "false"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "scim_client", string(sdk.ScimSecurityIntegrationScimClientAzure)),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "run_as_role", role.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "network_policy", networkPolicy.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", "foo"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_ScimIntegration/completeAzure"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_scim_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidScimClient(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable("invalid"),
			"sync_password":       config.BoolVariable(false),
			"network_policy_name": config.StringVariable("foo"),
			"run_as_role":         config.StringVariable(snowflakeroles.GenericScimProvisioner.Name()),
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
				ExpectError:     regexp.MustCompile(`invalid ScimSecurityIntegrationScimClientOption: INVALID`),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidRunAsRole(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"enabled":             config.BoolVariable(false),
			"scim_client":         config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientGeneric)),
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
				ExpectError:     regexp.MustCompile(`invalid ScimSecurityIntegrationRunAsRoleOption: INVALID`),
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

func TestAcc_ScimIntegration_InvalidCreateWithSyncPasswordWithAzure(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                config.StringVariable(id.Name()),
			"scim_client":         config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientAzure)),
			"run_as_role":         config.StringVariable(snowflakeroles.GenericScimProvisioner.Name()),
			"enabled":             config.BoolVariable(true),
			"sync_password":       config.BoolVariable(false),
			"network_policy_name": config.StringVariable(""),
			"comment":             config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			"can not CREATE scim integration with field `sync_password`",
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(),
			},
		},
	})
}

func TestAcc_ScimIntegration_InvalidUpdateWithSyncPasswordWithAzure(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"name":        config.StringVariable(id.Name()),
			"scim_client": config.StringVariable(string(sdk.ScimSecurityIntegrationScimClientAzure)),
			"run_as_role": config.StringVariable(snowflakeroles.GenericScimProvisioner.Name()),
			"enabled":     config.BoolVariable(true),
		}
		if complete {
			c["sync_password"] = config.BoolVariable(false)
			c["network_policy_name"] = config.StringVariable("")
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
		ErrorCheck: helpers.AssertErrorContainsPartsFunc(t, []string{
			"can not SET and UNSET field `sync_password`",
		}),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/basic"),
				ConfigVariables: m(false),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_ScimIntegration/complete"),
				ConfigVariables: m(true),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion092EnabledTrue(t *testing.T) {
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
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: scimIntegrationv092(id.Name(), role.Name(), sdk.ScimSecurityIntegrationScimClientGeneric),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "provisioner_role", role.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationv093(id.Name(), role.Name(), true, sdk.ScimSecurityIntegrationScimClientGeneric),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange("snowflake_scim_integration.test", "name", tfjson.ActionUpdate, sdk.String(id.Name()), sdk.String(id.Name())),
						planchecks.ExpectChange("snowflake_scim_integration.test", "enabled", tfjson.ActionUpdate, sdk.String("true"), sdk.String("true")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "scim_client", tfjson.ActionUpdate, sdk.String("GENERIC"), sdk.String("GENERIC")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "run_as_role", tfjson.ActionUpdate, sdk.String(role.Name()), sdk.String(role.Name())),
						planchecks.ExpectChange("snowflake_scim_integration.test", "network_policy", tfjson.ActionUpdate, sdk.String(""), sdk.String("")),
						planchecks.ExpectChange("snowflake_scim_integration.test", "sync_password", tfjson.ActionUpdate, nil, sdk.String(r.BooleanDefault)),
						planchecks.ExpectChange("snowflake_scim_integration.test", "comment", tfjson.ActionUpdate, nil, nil),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion092EnabledFalse(t *testing.T) {
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
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: scimIntegrationv092(id.Name(), role.Name(), sdk.ScimSecurityIntegrationScimClientGeneric),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "provisioner_role", role.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationv093(id.Name(), role.Name(), false, sdk.ScimSecurityIntegrationScimClientGeneric),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "run_as_role", role.Name()),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_migrateFromVersion093HandleSyncPassword(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	role := snowflakeroles.GenericScimProvisioner
	resourceName := "snowflake_scim_integration.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			// create resource with v0.92
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.92.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: scimIntegrationv092(id.Name(), role.Name(), sdk.ScimSecurityIntegrationScimClientAzure),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
				),
			},
			// migrate to v0.93 - there is a diff due to new field sync_password in state
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.93.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						planchecks.ExpectChange(resourceName, "sync_password", tfjson.ActionUpdate, nil, sdk.String(r.BooleanDefault)),
					},
				},
				Config: scimIntegrationv093(id.Name(), role.Name(), true, sdk.ScimSecurityIntegrationScimClientAzure),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
				),
				ExpectError: regexp.MustCompile("invalid property 'SYNC_PASSWORD' for 'INTEGRATION"),
			},
			// check with newest version - the value in state was set to boolean default, so there should be no diff
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationv093(id.Name(), role.Name(), true, sdk.ScimSecurityIntegrationScimClientAzure),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "sync_password", r.BooleanDefault),
				),
			},
		},
	})
}

func scimIntegrationv092(name, roleName string, scimClient sdk.ScimSecurityIntegrationScimClientOption) string {
	s := `
resource "snowflake_scim_integration" "test" {
	name             = "%s"
	scim_client      = "%s"
	provisioner_role = "%s"
}
`
	return fmt.Sprintf(s, name, scimClient, roleName)
}

func scimIntegrationv093(name, roleName string, enabled bool, scimClient sdk.ScimSecurityIntegrationScimClientOption) string {
	s := `
resource "snowflake_scim_integration" "test" {
	name             = "%s"
	enabled          = %t
	scim_client      = "%s"
	run_as_role		 = "%s"
}
`
	return fmt.Sprintf(s, name, enabled, scimClient, roleName)
}

func TestAcc_ScimIntegration_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: scimIntegrationBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationBasicConfig(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func TestAcc_ScimIntegration_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	quotedId := fmt.Sprintf(`\"%s\"`, id.Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             scimIntegrationBasicConfig(quotedId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "id", id.Name()),
				),
			},
			{
				PreConfig:                func() { acc.UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   scimIntegrationBasicConfig(quotedId),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_scim_integration.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_scim_integration.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "id", id.Name()),
				),
			},
		},
	})
}

func scimIntegrationBasicConfig(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_scim_integration" "test" {
	 name        = "%s"
	 scim_client = "GENERIC"
	 run_as_role = "%s"
	 enabled     = true
	}
	`, name, snowflakeroles.GenericScimProvisioner.Name())
}
