package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1423486): Fix using warehouse in all tests and remove unsetting testenvs.ConfigureClientOnce
// TODO(next pr): cleanup setting warehouse with unsafe_execute
func TestAcc_View_basic(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	rowAccessPolicy, rowAccessPolicyCleanup := acc.TestClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, sdk.DataTypeVARCHAR)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := acc.TestClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	rowAccessPolicy2, rowAccessPolicy2Cleanup := acc.TestClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, sdk.DataTypeVARCHAR)
	t.Cleanup(rowAccessPolicy2Cleanup)

	aggregationPolicy2, aggregationPolicy2Cleanup := acc.TestClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicy2Cleanup)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	otherStatement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"
	comment := "Terraform test resource'"

	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement)
	viewModelWithDependency := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})

	// generators currently don't handle lists, so use the old way
	basicUpdate := func(rap, ap sdk.SchemaObjectIdentifier, statement string) config.Variables {
		return config.Variables{
			"name":                          config.StringVariable(id.Name()),
			"database":                      config.StringVariable(id.DatabaseName()),
			"schema":                        config.StringVariable(id.SchemaName()),
			"statement":                     config.StringVariable(statement),
			"row_access_policy":             config.StringVariable(rap.FullyQualifiedName()),
			"row_access_policy_on":          config.ListVariable(config.StringVariable("ROLE_NAME")),
			"aggregation_policy":            config.StringVariable(ap.FullyQualifiedName()),
			"aggregation_policy_entity_key": config.ListVariable(config.StringVariable("ROLE_NAME")),
			"comment":                       config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: accconfig.FromModel(t, viewModelWithDependency) + useWarehouseConfig(acc.TestWarehouseName),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName())),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModel(t, viewModel),
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeSnowflakeID(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasStatementString(statement)),
			},
			// set policies externally
			{
				PreConfig: func() {
					acc.TestClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []sdk.Column{{Value: "ROLE_NAME"}})))
					acc.TestClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithSetAggregationPolicy(*sdk.NewViewSetAggregationPolicyRequest(aggregationPolicy)))
				},
				Config: accconfig.FromModel(t, viewModel),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "0")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "0")),
				),
			},
			// set other fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy.ID(), aggregationPolicy, statement),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},
			// change policies
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy2.ID(), aggregationPolicy2, statement),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},
			// change statement and policies
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy.ID(), aggregationPolicy, otherStatement),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},
			// change statements externally
			{
				PreConfig: func() {
					acc.TestClient().View.RecreateView(t, id, statement)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy.ID(), aggregationPolicy, otherStatement),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},
			// unset policies externally
			{
				PreConfig: func() {
					acc.TestClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithDropAllRowAccessPolicies(true))
					acc.TestClient().View.Alter(t, sdk.NewAlterViewRequest(id).WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest()))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy.ID(), aggregationPolicy, otherStatement),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},

			// import - with optionals
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/basic_update"),
				ConfigVariables: basicUpdate(rowAccessPolicy.ID(), aggregationPolicy, otherStatement),
				ResourceName:    "snowflake_view.test",
				ImportState:     true,
				ImportStateCheck: assert.AssertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeSnowflakeID(id)).
						HasNameString(id.Name()).
						HasStatementString(otherStatement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasIsSecureString("false").
						HasIsTemporaryString("false").
						HasChangeTrackingString("false"),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.entity_key.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.entity_key.0", "ROLE_NAME")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.on.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.on.0", "ROLE_NAME")),
				),
			},
			// unset
			{
				Config:       accconfig.FromModel(t, viewModel.WithStatement(otherStatement)),
				ResourceName: "snowflake_view.test",
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(""),
					assert.Check(resource.TestCheckNoResourceAttr("snowflake_view.test", "aggregation_policy.#")),
					assert.Check(resource.TestCheckNoResourceAttr("snowflake_view.test", "row_access_policy.#")),
				),
			},
			// recreate - change is_recursive
			{
				Config: accconfig.FromModel(t, viewModel.WithIsRecursive("true")),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(otherStatement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("").
					HasIsRecursiveString("true").
					HasIsTemporaryString("default").
					HasChangeTrackingString("default"),
					assert.Check(resource.TestCheckNoResourceAttr("snowflake_view.test", "aggregation_policy.#")),
					assert.Check(resource.TestCheckNoResourceAttr("snowflake_view.test", "row_access_policy.#")),
				),
			},
		},
	})
}

func TestAcc_View_recursive(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel.WithIsRecursive("true")) + useWarehouseConfig(acc.TestWarehouseName),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasIsRecursiveString("true")),
			},
			{
				Config:       accconfig.FromModel(t, viewModel.WithIsRecursive("true")) + useWarehouseConfig(acc.TestWarehouseName),
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeSnowflakeID(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasStatementString(statement).
						HasIsRecursiveString("true")),
			},
		},
	})
}

func TestAcc_View_temporary(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	// we use one configured client, so a temporary view should be visible after creation
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel.WithIsTemporary("true")) + useWarehouseConfig(acc.TestWarehouseName),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasIsTemporaryString("true")),
			},
		},
	})
}

func TestAcc_View_complete(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	// use a simple table to test change_tracking, otherwise it fails with: Change tracking is not supported on queries with joins of type '[LEFT_OUTER_JOIN]'
	table, tableCleanup := acc.TestClient().Table.CreateTable(t)
	t.Cleanup(tableCleanup)
	statement := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())
	rowAccessPolicy, rowAccessPolicyCleanup := acc.TestClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, sdk.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := acc.TestClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":                          config.StringVariable(id.Name()),
			"database":                      config.StringVariable(id.DatabaseName()),
			"schema":                        config.StringVariable(id.SchemaName()),
			"comment":                       config.StringVariable("Terraform test resource"),
			"is_secure":                     config.BoolVariable(true),
			"is_temporary":                  config.BoolVariable(false),
			"or_replace":                    config.BoolVariable(false),
			"copy_grants":                   config.BoolVariable(false),
			"change_tracking":               config.BoolVariable(true),
			"row_access_policy":             config.StringVariable(rowAccessPolicy.ID().FullyQualifiedName()),
			"row_access_policy_on":          config.ListVariable(config.StringVariable("ID")),
			"aggregation_policy":            config.StringVariable(aggregationPolicy.FullyQualifiedName()),
			"aggregation_policy_entity_key": config.ListVariable(config.StringVariable("ID")),
			"statement":                     config.StringVariable(statement),
			"warehouse":                     config.StringVariable(acc.TestWarehouseName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/complete"),
				ConfigVariables: m(),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString("Terraform test resource").
					HasIsSecureString("true").
					HasIsTemporaryString("false").
					HasChangeTrackingString("true"),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "aggregation_policy.0.entity_key.0", "ID")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.#", "1")),
					assert.Check(resource.TestCheckResourceAttr("snowflake_view.test", "row_access_policy.0.on.0", "ID")),
					resourceshowoutputassert.ViewShowOutput(t, "snowflake_view.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("Terraform test resource").
						HasIsSecure(true).
						HasChangeTracking("ON"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View/complete"),
				ConfigVariables: m(),
				ResourceName:    "snowflake_view.test",
				ImportState:     true,
				ImportStateCheck: assert.AssertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeSnowflakeID(id)).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("Terraform test resource").
						HasIsSecureString("true").
						HasIsTemporaryString("false").HasChangeTrackingString("true"),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.policy_name", aggregationPolicy.FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.entity_key.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "aggregation_policy.0.entity_key.0", "ID")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.policy_name", rowAccessPolicy.ID().FullyQualifiedName())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.on.#", "1")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "row_access_policy.0.on.0", "ID")),
				),
			},
		},
	})
}

func TestAcc_View_Rename(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithComment("foo").WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})
	newViewModel := model.View("test", newId.DatabaseName(), newId.Name(), newId.SchemaName(), statement).WithComment("foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_view.test", "fully_qualified_name", id.FullyQualifiedName()),
				),
			},
			// rename with one param changed
			{
				Config: accconfig.FromModel(t, newViewModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "foo"),
					resource.TestCheckResourceAttr("snowflake_view.test", "fully_qualified_name", newId.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithIsSecure("true").WithOrReplace(false).WithCopyGrants(false).
		WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true),
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					}),
				),
			},
			// Checks that copy_grants changes don't trigger a drop
			{
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(true).WithOrReplace(true)) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrantsReversed(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithIsSecure("true").WithOrReplace(true).WithCopyGrants(true).
		WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			{
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(false)) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_ViewCopyGrantsStatementUpdate(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	viewId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithGrants(viewId, tableId, `\"name\"`) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: viewConfigWithGrants(viewId, tableId, "*") + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAcc_View_copyGrants(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModel(t, viewModel.WithCopyGrants(true)) + useWarehouseConfig(acc.TestWarehouseName),
				ExpectError: regexp.MustCompile("all of `copy_grants,or_replace` must be specified"),
			},
			{
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(true).WithOrReplace(true)) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
				),
			},
			{
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(false).WithOrReplace(true)) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
				),
			},
		},
	})
}

func TestAcc_View_Issue2640(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	part1 := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	part2 := "SELECT ROLE_OWNER, ROLE_NAME FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	statement := fmt.Sprintf("%s\n\tunion\n%s\n", part1, part2)
	roleId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithMultilineUnionStatement(id, part1, part2) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", statement),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", acc.TestSchemaName),
				),
			},
			// try to import secure view without being its owner (proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2640)
			{
				PreConfig: func() {
					role, roleCleanup := acc.TestClient().Role.CreateRoleWithIdentifier(t, roleId)
					t.Cleanup(roleCleanup)
					acc.TestClient().Role.GrantOwnershipOnSchemaObject(t, role.ID(), id, sdk.ObjectTypeView, sdk.Revoke)
				},
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ExpectError:  regexp.MustCompile("`text` is missing; if the view is secure then the role used by the provider must own the view"),
			},
			// import with the proper role
			{
				PreConfig: func() {
					acc.TestClient().Role.GrantOwnershipOnSchemaObject(t, snowflakeroles.Accountadmin, id, sdk.ObjectTypeView, sdk.Revoke)
				},
				ResourceName: "snowflake_view.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t, assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeSnowflakeID(id), "name", id.Name())),
					resourceassert.ImportedViewResource(t, helpers.EncodeSnowflakeID(id)).
						HasNameString(id.Name()).
						HasStatementString(statement).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()),
				),
			},
		},
	})
}

func TestAcc_view_migrateFromVersion_0_94_1(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_view.test"
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithDependsOn([]string{"snowflake_unsafe_execute.use_warehouse"})

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: viewv_0_94_1_WithTags(id, tag.SchemaName, tag.Name, "foo", statement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.0.name", tag.Name),
					resource.TestCheckResourceAttr(resourceName, "tag.0.value", "foo"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModel(t, viewModel) + useWarehouseConfig(acc.TestWarehouseName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "tag.#"),
				),
			},
		},
	})
}

func useWarehouseConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_unsafe_execute" "use_warehouse" {
	execute			= "USE WAREHOUSE \"%s\""
	revert			= "SELECT 1"
}
`, name)
}

func viewv_0_94_1_WithTags(id sdk.SchemaObjectIdentifier, tagSchema, tagName, tagValue, statement string) string {
	s := `
resource "snowflake_view" "test" {
	name					= "%[1]s"
	database				= "%[2]s"
	schema				= "%[6]s"
	statement				= "%[7]s"
	tag {
		name = "%[4]s"
		value = "%[5]s"
		schema = "%[3]s"
		database = "%[2]s"
	}
}
`
	return fmt.Sprintf(s, id.Name(), id.DatabaseName(), tagSchema, tagName, tagValue, id.SchemaName(), statement)
}

func viewConfigWithGrants(viewId, tableId sdk.SchemaObjectIdentifier, selectStatement string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "table" {
  database = "%[1]s"
  schema = "%[2]s"
  name     = "%[3]s"

  column {
    name = "name"
    type = "text"
  }
}

resource "snowflake_view" "test" {
  name = "%[4]s"
  comment = "created by terraform"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "select %[5]s from \"%[1]s\".\"%[2]s\".\"${snowflake_table.table.name}\""
  or_replace = true
  copy_grants = true
  is_secure = true
  depends_on = [snowflake_unsafe_execute.use_warehouse, snowflake_table.table]
}

resource "snowflake_account_role" "test" {
  name = "test"
}

resource "snowflake_grant_privileges_to_account_role" "grant" {
  privileges        = ["SELECT"]
  account_role_name = snowflake_account_role.test.name
  on_schema_object {
    object_type = "VIEW"
    object_name = "\"%[1]s\".\"%[2]s\".\"${snowflake_view.test.name}\""
  }
}

data "snowflake_grants" "grants" {
  depends_on = [snowflake_grant_privileges_to_account_role.grant, snowflake_view.test, snowflake_unsafe_execute.use_warehouse]
  grants_on {
    object_name = "\"%[1]s\".\"%[2]s\".\"${snowflake_view.test.name}\""
    object_type = "VIEW"
  }
}
	`, viewId.DatabaseName(), viewId.SchemaName(), tableId.Name(), viewId.Name(), selectStatement)
}

func viewConfigWithMultilineUnionStatement(id sdk.SchemaObjectIdentifier, part1 string, part2 string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = <<-SQL
%[4]s
	union
%[5]s
SQL
  is_secure = true
  depends_on = [snowflake_unsafe_execute.use_warehouse]
}
	`, id.DatabaseName(), id.SchemaName(), id.Name(), part1, part2)
}
