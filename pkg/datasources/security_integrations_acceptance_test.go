package datasources_test

import (
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1348100]: add other security integrations when they are ready
// TODO [SNOW-1348100]: test specific describe properties
func TestAcc_SecurityIntegrations_Scim(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	configVariables := config.Variables{
		"name":    config.StringVariable(id.Name()),
		"comment": config.StringVariable(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.0.enabled.0.value", "false"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.integration_type", "SCIM - GENERIC"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.category", sdk.SecurityIntegrationCategory),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.enabled", "false"),
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_security_integrations.test", "security_integrations.0.show_output.0.created_on"),

					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	commonVariables := config.Variables{
		"name_1": config.StringVariable(idOne.Name()),
		"name_2": config.StringVariable(idTwo.Name()),
		"name_3": config.StringVariable(idThree.Name()),
	}

	likeConfig := config.Variables{
		"like": config.StringVariable(idOne.Name()),
	}
	maps.Copy(likeConfig, commonVariables)

	likeConfig2 := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}
	maps.Copy(likeConfig2, commonVariables)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ScimSecurityIntegration),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_security_integrations.test", "security_integrations.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SecurityIntegrations_SecurityIntegrationNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SecurityIntegrations/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one security integration"),
			},
		},
	})
}
