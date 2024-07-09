package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlit_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", r.BooleanDefault),
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
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "sync_password", r.BooleanDefault),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_Streamlit_complete(t *testing.T) {
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	integrationId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	ctx := context.Background()
	acc.TestAccPreCheck(t)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.CreateNetworkRule(t)
	t.Cleanup(networkRuleCleanup)
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	err := acc.Client(t).Sessions.UseSchema(ctx, schemaId)
	require.NoError(t, err)
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	_, err = acc.Client(t).ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = TRUE`, integrationId.Name(), networkRule.ID().Name()))
	require.NoError(t, err)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                     config.StringVariable(databaseId.Name()),
			"schema":                       config.StringVariable(schemaId.Name()),
			"stage":                        config.StringVariable(stage.ID().Name()),
			"name":                         config.StringVariable(id.Name()),
			"main_file":                    config.StringVariable("foo"),
			"query_warehouse":              config.StringVariable(warehouse.ID().Name()),
			"external_access_integrations": config.SetVariable(config.StringVariable(integrationId.Name())),
			"title":                        config.StringVariable("foo"),
			"comment":                      config.StringVariable("foo"),
		}
	}
	rootLocation := fmt.Sprintf(`@"%s"."%s".%s`, databaseId.Name(), schemaId.Name(), stage.Name)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "root_location", rootLocation),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.0", integrationId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", snowflakeroles.Accountadmin.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.url_id"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.root_location", rootLocation),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", integrationId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.external_access_secrets"),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables:   m(),
				ResourceName:      "snowflake_streamlit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_Streamlit_InvalidQueryWarehouse(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	client := acc.Client(t)
	ctx := context.Background()
	acc.TestAccPreCheck(t)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.CreateNetworkRule(t)
	t.Cleanup(networkRuleCleanup)
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION "%s" ALLOWED_NETWORK_RULES = (%s) ENABLED = FALSE`, integrationId.Name(), networkRule.ID().FullyQualifiedName()))
	require.NoError(t, err)
	schemaId := acc.TestClient().Ids.SchemaId()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"schema":                       config.StringVariable(schemaId.FullyQualifiedName()),
			"name":                         config.StringVariable(id.Name()),
			"root_location":                config.StringVariable("foo"),
			"main_file":                    config.StringVariable("foo"),
			"query_warehouse":              config.StringVariable("invalid"),
			"external_access_integrations": config.SetVariable(config.StringVariable(integrationId.FullyQualifiedName())),
			"title":                        config.StringVariable("foo"),
			"comment":                      config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile(`Invalid identifier type`),
			},
		},
	})
	defer func() {
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`DROP EXTERNAL ACCESS INTEGRATION "%s"`, integrationId.Name()))
		require.NoError(t, err)
	}()
}
