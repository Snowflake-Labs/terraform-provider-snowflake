package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlit_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

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
			c["external_access_integrations"] = config.SetVariable(config.StringVariable(externalAccessIntegrationId.Name()))
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
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
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "external_access_integrations.0", externalAccessIntegrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "title", "foo"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "comment", "foo")),
			},
			// change externally
			{
				PreConfig: func() {
					acc.TestClient().Streamlit.Update(t, sdk.NewAlterStreamlitRequest(id).WithSet(
						*sdk.NewStreamlitSetRequest().
							WithRootLocation(fmt.Sprintf("@%s/bar", stage.ID().FullyQualifiedName())).
							WithComment("bar"),
					))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/complete"),
				ConfigVariables: m(true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_streamlit.test", plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift("snowflake_streamlit.test", "directory_location", sdk.String("foo"), sdk.String("bar")),
						planchecks.ExpectChange("snowflake_streamlit.test", "directory_location", tfjson.ActionUpdate, sdk.String("bar"), sdk.String("foo")),
						planchecks.ExpectDrift("snowflake_streamlit.test", "comment", sdk.String("foo"), sdk.String("bar")),
						planchecks.ExpectChange("snowflake_streamlit.test", "comment", tfjson.ActionUpdate, sdk.String("bar"), sdk.String("foo")),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "database", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "schema", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "main_file", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "comment", "foo"),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
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
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

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
			"external_access_integrations": config.SetVariable(config.StringVariable(externalAccessIntegrationId.Name())),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "fully_qualified_name", id.FullyQualifiedName()),

					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_streamlit.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
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
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
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
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	m := func(name, comment string) map[string]config.Variable {
		return map[string]config.Variable{
			"database":  config.StringVariable(databaseId.Name()),
			"schema":    config.StringVariable(schemaId.Name()),
			"stage":     config.StringVariable(stage.ID().FullyQualifiedName()),
			"name":      config.StringVariable(name),
			"main_file": config.StringVariable("foo"),
			"comment":   config.StringVariable(comment),
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basicWithComment"),
				ConfigVariables: m(id.Name(), "foo"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlit/basicWithComment"),
				ConfigVariables: m(newId.Name(), "bar"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_streamlit.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "show_output.0.comment", "bar"),
					resource.TestCheckResourceAttr("snowflake_streamlit.test", "fully_qualified_name", newId.FullyQualifiedName()),
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
