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

func TestAcc_View_basic(t *testing.T) {
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

	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement)

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
			"comment":                       config.StringVariable("Terraform test resource"),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: accconfig.FromModel(t, viewModel),
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
					HasCommentString("Terraform test resource"),
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
					HasCommentString("Terraform test resource"),
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
					HasCommentString("Terraform test resource"),
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
					HasCommentString("Terraform test resource"),
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
						HasCommentString("Terraform test resource").
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
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
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel.WithIsRecursive("true")),
				Check: assert.AssertThat(t, resourceassert.ViewResource(t, "snowflake_view.test").
					HasNameString(id.Name()).
					HasStatementString(statement).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasIsRecursiveString("true")),
			},
			{
				Config:       accconfig.FromModel(t, viewModel.WithIsRecursive("true")),
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

func TestAcc_View_complete(t *testing.T) {
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
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithComment("foo")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModel(t, viewModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "foo"),
				),
			},
			// rename with one param changed
			{
				Config: accconfig.FromModel(t, model.View("test", newId.DatabaseName(), newId.Name(), newId.SchemaName(), statement).WithComment("foo")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_view.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "foo"),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithIsSecure("true").WithOrReplace(false).WithCopyGrants(false)

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
				Config: accconfig.FromModel(t, viewModel),
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
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(true).WithOrReplace(true)),
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
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement).WithIsSecure("true").WithOrReplace(true).WithCopyGrants(true)

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
				Config: accconfig.FromModel(t, viewModel),
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
				Config: accconfig.FromModel(t, viewModel.WithCopyGrants(false)),
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
	tableName := acc.TestClient().Ids.Alpha()
	viewName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithGrants(acc.TestDatabaseName, acc.TestSchemaName, tableName, viewName, `\"name\"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: viewConfigWithGrants(acc.TestDatabaseName, acc.TestSchemaName, tableName, viewName, "*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAcc_View_copyGrants(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()
	query := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.View),
		Steps: []resource.TestStep{
			{
				Config:      viewConfigWithCopyGrants(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true),
				ExpectError: regexp.MustCompile("all of `copy_grants,or_replace` must be specified"),
			},
			{
				Config: viewConfigWithCopyGrantsAndOrReplace(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
				),
			},
			{
				Config: viewConfigWithOrReplace(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
				),
			},
		},
	})
}

func TestAcc_View_Issue2640(t *testing.T) {
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
				Config: viewConfigWithMultilineUnionStatement(acc.TestDatabaseName, acc.TestSchemaName, id.Name(), part1, part2),
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

func TestAcc_view_migrateFromVersion094(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_view.test"
	statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	viewModel := model.View("test", id.DatabaseName(), id.Name(), id.SchemaName(), statement)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: viewv094WithTags(id, tag.SchemaName, tag.Name, "foo", statement),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tag.0.name", tag.Name),
					resource.TestCheckResourceAttr(resourceName, "tag.0.value", "foo"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModel(t, viewModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckNoResourceAttr(resourceName, "tag.#"),
				),
			},
		},
	})
}

func viewv094WithTags(id sdk.SchemaObjectIdentifier, tagSchema, tagName, tagValue, statement string) string {
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

func viewConfigWithGrants(databaseName string, schemaName string, tableName string, viewName string, selectStatement string) string {
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
  depends_on = [snowflake_table.table]
  name = "%[4]s"
  comment = "created by terraform"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "select %[5]s from \"%[1]s\".\"%[2]s\".\"${snowflake_table.table.name}\""
  or_replace = true
  copy_grants = true
  is_secure = true
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
  depends_on = [snowflake_grant_privileges_to_account_role.grant, snowflake_view.test]
  grants_on {
    object_name = "\"%[1]s\".\"%[2]s\".\"${snowflake_view.test.name}\""
    object_type = "VIEW"
  }
}
	`, databaseName, schemaName, tableName, viewName, selectStatement)
}

func viewConfigWithCopyGrants(databaseName string, schemaName string, name string, selectStatement string, copyGrants bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  copy_grants = %[5]t
}
	`, databaseName, schemaName, name, selectStatement, copyGrants)
}

func viewConfigWithCopyGrantsAndOrReplace(databaseName string, schemaName string, name string, selectStatement string, copyGrants bool, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  copy_grants = %[5]t
  or_replace = %[6]t
}
	`, databaseName, schemaName, name, selectStatement, copyGrants, orReplace)
}

func viewConfigWithOrReplace(databaseName string, schemaName string, name string, selectStatement string, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  or_replace = %[5]t
}
	`, databaseName, schemaName, name, selectStatement, orReplace)
}

func viewConfigWithMultilineUnionStatement(databaseName string, schemaName string, name string, part1 string, part2 string) string {
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
}
	`, databaseName, schemaName, name, part1, part2)
}
