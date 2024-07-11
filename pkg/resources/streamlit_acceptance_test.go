package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlit_basic(t *testing.T) {
	acc.TestAccPreCheck(t)
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	ctx := context.Background()
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	err := acc.Client(t).Sessions.UseSchema(ctx, schemaId)
	require.NoError(t, err)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	integrationId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.CreateNetworkRule(t)
	t.Cleanup(networkRuleCleanup)
	_, err = acc.Client(t).ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = TRUE`, integrationId.Name(), networkRule.ID().Name()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = acc.Client(t).ExecForTests(ctx, fmt.Sprintf(`DROP EXTERNAL ACCESS INTEGRATION %s`, integrationId.Name()))
		require.NoError(t, err)
	})
	rootLocation := fmt.Sprintf("@%s", stage.ID().FullyQualifiedName())
	rootLocationWithCatalog := fmt.Sprintf("%s/foo", rootLocation)
	m := func(complete bool) map[string]config.Variable {
		c := map[string]config.Variable{
			"database":  config.StringVariable(databaseId.Name()),
			"schema":    config.StringVariable(schemaId.Name()),
			"stage":     config.StringVariable(stage.ID().FullyQualifiedName()),
			"name":      config.StringVariable(id.Name()),
			"main_file": config.StringVariable("foo"),
		}
		if complete {
			c["directory_location"] = config.StringVariable("foo")
			c["query_warehouse"] = config.StringVariable(warehouse.ID().Name())
			c["external_access_integrations"] = config.SetVariable(config.StringVariable(integrationId.Name()))
			c["title"] = config.StringVariable("foo")
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
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "main_file", "foo"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", ""),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", snowflakeroles.Accountadmin.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.query_warehouse", ""),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.url_id"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.title", ""),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.root_location", rootLocation),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.query_warehouse", ""),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.#", "0"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.external_access_secrets"),
				),
			},
			// import - without optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(false),
				ResourceName:    "snowflake_streamlit.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "database", databaseId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "schema", schemaId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "stage", stage.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "main_file", "foo"),
				),
			},
			// set optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.root_location", rootLocationWithCatalog),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", integrationId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.external_access_secrets")),
			},
			// import - complete
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(true),
				ResourceName:    "snowflake_streamlit.test",
				ImportState:     true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "database", databaseId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "schema", schemaId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "stage", stage.ID().FullyQualifiedName()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "directory_location", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "main_file", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "query_warehouse", warehouse.ID().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "external_access_integrations.#", "1"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "external_access_integrations.0", integrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "title", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "comment", "foo")),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().Streamlit.Update(t, sdk.NewAlterStreamlitRequest(id).WithSet(
						*sdk.NewStreamlitSetRequest().
							WithRootLocation(fmt.Sprintf("@%s/bar", stage.ID().FullyQualifiedName())),
					))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_streamlit.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.root_location", rootLocationWithCatalog),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", integrationId.Name()),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.external_access_secrets")),
			},
			// unset
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "directory_location", ""),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "main_file", "foo")),
			},
		},
	})
}

func TestAcc_Streamlit_complete(t *testing.T) {
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	ctx := context.Background()
	acc.TestAccPreCheck(t)
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	err := acc.Client(t).Sessions.UseSchema(ctx, schemaId)
	require.NoError(t, err)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	integrationId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.CreateNetworkRule(t)
	t.Cleanup(networkRuleCleanup)
	_, err = acc.Client(t).ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = TRUE`, integrationId.Name(), networkRule.ID().Name()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = acc.Client(t).ExecForTests(ctx, fmt.Sprintf(`DROP EXTERNAL ACCESS INTEGRATION %s`, integrationId.Name()))
		require.NoError(t, err)
	})
	rootLocation := fmt.Sprintf("@%s/foo", stage.ID().FullyQualifiedName())
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                     config.StringVariable(databaseId.Name()),
			"schema":                       config.StringVariable(schemaId.Name()),
			"stage":                        config.StringVariable(stage.ID().FullyQualifiedName()),
			"directory_location":           config.StringVariable("foo"),
			"name":                         config.StringVariable(id.Name()),
			"main_file":                    config.StringVariable("foo"),
			"query_warehouse":              config.StringVariable(warehouse.ID().Name()),
			"external_access_integrations": config.SetVariable(config.StringVariable(integrationId.Name())),
			"title":                        config.StringVariable("foo"),
			"comment":                      config.StringVariable("foo"),
		}
	}
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "stage", stage.ID().FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "directory_location", "foo"),
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
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "describe_output.0.user_packages.#"),
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

func TestAcc_Streamlit_Rename(t *testing.T) {
	acc.TestAccPreCheck(t)
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	m := func(name string) map[string]config.Variable {
		return map[string]config.Variable{
			"database":  config.StringVariable(databaseId.Name()),
			"schema":    config.StringVariable(schemaId.Name()),
			"stage":     config.StringVariable(stage.ID().FullyQualifiedName()),
			"name":      config.StringVariable(name),
			"main_file": config.StringVariable("foo"),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.NetworkPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(id.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(newId.Name()),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_streamlit.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", newId.Name()),
				),
			},
		},
	})
}

func TestAcc_Streamlit_InvalidStage(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	schemaId := acc.TestClient().Ids.SchemaId()
	databaseId := acc.TestClient().Ids.DatabaseId()

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"schema":          config.StringVariable(schemaId.FullyQualifiedName()),
			"database":        config.StringVariable(databaseId.FullyQualifiedName()),
			"name":            config.StringVariable(id.Name()),
			"stage":           config.StringVariable("foo"),
			"main_file":       config.StringVariable("foo"),
			"query_warehouse": config.StringVariable("foo"),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basic"),
				ConfigVariables: m(),
				ExpectError:     regexp.MustCompile(`Invalid identifier type`),
			},
		},
	})
}
