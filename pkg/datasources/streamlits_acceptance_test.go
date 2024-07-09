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
	configVariables := config.Variables{
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
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.show_output.0.owner", snowflakeroles.Accountadmin.Name()),
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
					resource.TestCheckResourceAttr("data.snowflake_streamlits.test", "streamlits.0.describe_output.0.external_access_integrations.0", integrationId.Name()),
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
