package datasources_test

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlits(t *testing.T) {
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	acc.TestAccPreCheck(t)
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	// use schema is needed because otherwise reference to external access integration fails.
	err := acc.Client(t).Sessions.UseSchema(context.Background(), schemaId)
	require.NoError(t, err)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.CreateNetworkRule(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)
	rootLocation := fmt.Sprintf("@%s/foo", stage.ID().FullyQualifiedName())
	configVariables := config.Variables{
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/optionals_set"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.show_output.0.url_id"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.title", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.root_location", rootLocation),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.main_file", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.external_access_secrets"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/optionals_unset"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.#", "1"),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.title", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.owner", snowflakeroles.Accountadmin.Name()),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.comment", "foo"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet("data.snowflake_streamlits.test", "streamlits.0.show_output.0.url_id"),
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Streamlits_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	acc.TestAccPreCheck(t)
	stage, stageCleanup := acc.TestClient().Stage.CreateStageInSchema(t, schemaId)
	t.Cleanup(stageCleanup)
	commonVariables := config.Variables{
		"name_1":    config.StringVariable(idOne.Name()),
		"name_2":    config.StringVariable(idTwo.Name()),
		"name_3":    config.StringVariable(idThree.Name()),
		"main_file": config.StringVariable("foo"),
		"schema":    config.StringVariable(schemaId.Name()),
		"stage":     config.StringVariable(stage.ID().FullyQualifiedName()),
		"database":  config.StringVariable(databaseId.Name()),
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
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/like"),
				ConfigVariables: likeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.#", "1"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/like"),
				ConfigVariables: likeConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Streamlits_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamlitsDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Streamlits_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamlitsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Streamlits_StreamlitNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one streamlit"),
			},
		},
	})
}

func streamlitsDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_streamlits" "test" {
  in {
    database = "%s"
    schema   = "%s"
  }
}
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func streamlitsDatasourceEmptyIn() string {
	return `
data "snowflake_streamlits" "test" {
  in {
  }
}
`
}
