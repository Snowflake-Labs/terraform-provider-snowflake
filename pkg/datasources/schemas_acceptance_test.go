package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schemas_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(id)
	statement := "SELECT ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	columnNames := []string{"ROLE_NAME"}

	schemaModel := model.Schema("test", id.DatabaseName(), id.Name()).
		WithComment(comment).
		WithIsTransient(datasources.BooleanTrue).
		WithWithManagedAccess(datasources.BooleanTrue)
	viewModel := model.View("test", viewId.DatabaseName(), viewId.Name(), viewId.SchemaName(), statement).
		WithColumnNames(columnNames...).
		WithDependsOn(schemaModel.ResourceReference())
	schemasModel := datasourcemodel.Schemas("test").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithLimit(1).
		WithDependsOn(schemaModel.ResourceReference(), viewModel.ResourceReference())
	schemasModelWithoutAdditional := datasourcemodel.Schemas("test").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithLimit(1).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(schemaModel.ResourceReference(), viewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schemaModel, viewModel, schemasModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.show_output.0.is_current"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.show_output.0.options", "TRANSIENT, MANAGED ACCESS"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.parameters.#", "1"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.data_retention_time_in_days.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.max_data_extension_time_in_days.0.value"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.parameters.0.external_volume.0.value", ""),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.parameters.0.catalog.0.value", ""),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.replace_invalid_characters.0.value"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.parameters.0.default_ddl_collation.0.value", ""),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.storage_serialization_policy.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.log_level.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.trace_level.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.suspend_task_after_num_failures.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.task_auto_retry_attempts.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.user_task_managed_initial_warehouse_size.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.user_task_minimum_trigger_interval_in_seconds.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.quoted_identifiers_ignore_case.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.enable_console_output.0.value"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.parameters.0.pipe_execution_paused.0.value"),

					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.created_on"),
					resource.TestCheckResourceAttrSet(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.name"),
					resource.TestCheckResourceAttr(schemasModel.DatasourceReference(), "schemas.0.describe_output.0.kind", "VIEW"),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModel, viewModel, schemasModelWithoutAdditional),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttrSet(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.is_current"),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttrSet(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.owner"),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.options", "TRANSIENT, MANAGED ACCESS"),
					resource.TestCheckResourceAttrSet(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.retention_time"),
					resource.TestCheckResourceAttrSet(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.show_output.0.owner_role_type"),

					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(schemasModelWithoutAdditional.DatasourceReference(), "schemas.0.parameters.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Schemas_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	database2, database2Cleanup := acc.TestClient().Database.DatabaseWithParametersSet(t)
	t.Cleanup(database2Cleanup)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "1")
	idTwo := acc.TestClient().Ids.RandomDatabaseObjectIdentifierWithPrefix(prefix + "2")
	idThree := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	idFour := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(database2.ID())

	schemaModel1 := model.Schema("test_1", idOne.DatabaseName(), idOne.Name())
	schemaModel2 := model.Schema("test_2", idTwo.DatabaseName(), idTwo.Name())
	schemaModel3 := model.Schema("test_3", idThree.DatabaseName(), idThree.Name())
	schemaModel4 := model.Schema("test_4", idFour.DatabaseName(), idFour.Name())
	schemasModelLike := datasourcemodel.Schemas("test1").
		WithLike(idOne.Name()).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())
	schemasModelStartsWith := datasourcemodel.Schemas("test2").
		WithStartsWith(prefix).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())
	schemasModelLimit := datasourcemodel.Schemas("test3").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())
	schemasModelIn := datasourcemodel.Schemas("test4").
		WithIn(idFour.DatabaseId()).
		WithStartsWith(idFour.Name()).
		WithDependsOn(schemaModel1.ResourceReference(), schemaModel2.ResourceReference(), schemaModel3.ResourceReference(), schemaModel4.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Schema),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelLike),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelLike.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelLike.DatasourceReference(), "schemas.0.show_output.0.name", idOne.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelStartsWith.DatasourceReference(), "schemas.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelLimit.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelLimit.DatasourceReference(), "schemas.0.show_output.0.name", idOne.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, schemaModel1, schemaModel2, schemaModel3, schemaModel4, schemasModelIn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(schemasModelIn.DatasourceReference(), "schemas.#", "1"),
					resource.TestCheckResourceAttr(schemasModelIn.DatasourceReference(), "schemas.0.show_output.0.name", idFour.Name()),
				),
			},
		},
	})
}

func TestAcc_Schemas_SchemaNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Schemas/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one schema"),
			},
		},
	})
}

func TestAcc_Schemas_BadCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      schemasDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func schemasDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_schemas" "test" {
  in {
    database = "%s"
    application = "foo"
    application_package = "bar"
  }
}
`, acc.TestDatabaseName)
}
