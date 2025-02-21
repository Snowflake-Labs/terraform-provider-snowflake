package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_FunctionPython_InlineBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	funcName := "some_function"
	argName := "x"
	dataType := testdatatypes.DataTypeNumber_36_2

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	definition := acc.TestClient().Function.SamplePythonDefinition(t, funcName, argName)

	functionModel := model.FunctionPythonBasicInline("test", id, "3.8", dataType, funcName, definition).
		WithArgument(argName, dataType)
	functionModelRenamed := model.FunctionPythonBasicInline("test", idWithChangedNameButTheSameDataType, "3.8", dataType, funcName, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionPython),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasRuntimeVersionString("3.8").
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("PYTHON").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// REMOVE EXTERNALLY (CHECK RECREATION)
			{
				PreConfig: func() {
					acc.TestClient().Function.DropFunctionFunc(t, id)()
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arguments.0.arg_data_type", "is_aggregate", "is_secure", "null_input_behavior", "return_results_behavior"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedFunctionPythonResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "NUMBER(38, 0)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assert.AssertThat(t,
					resourceassert.FunctionPythonResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionPython_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	secretId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	secretId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	secret, secretCleanup := acc.TestClient().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
	t.Cleanup(secretCleanup)

	secret2, secret2Cleanup := acc.TestClient().Secret.CreateWithGenericString(t, secretId2, "test_secret_string_2")
	t.Cleanup(secret2Cleanup)

	externalAccessIntegration, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	externalAccessIntegration2, externalAccessIntegration2Cleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret2.ID())
	t.Cleanup(externalAccessIntegration2Cleanup)

	tmpPythonFunction := acc.TestClient().CreateSamplePythonFunctionAndModule(t)
	tmpPythonFunction2 := acc.TestClient().CreateSamplePythonFunctionAndModule(t)

	funcName := "some_function"
	argName := "x"
	dataType := testdatatypes.DataTypeNumber_36_2

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	definition := acc.TestClient().Function.SamplePythonDefinition(t, funcName, argName)

	functionModel := model.FunctionPythonBasicInline("test", id, "3.8", dataType, funcName, definition).
		WithIsSecure(r.BooleanFalse).
		WithArgument(argName, dataType).
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some comment").
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction.ModuleName + ".py"},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction2.ModuleName + ".py"},
		).
		WithPackages("numpy", "pandas").
		WithExternalAccessIntegrations(externalAccessIntegration, externalAccessIntegration2).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"abc": secretId,
			"def": secretId2,
		})

	functionModelUpdateWithoutRecreation := model.FunctionPythonBasicInline("test", id, "3.8", dataType, funcName, definition).
		WithIsSecure(r.BooleanFalse).
		WithArgument(argName, dataType).
		WithNullInputBehavior(string(sdk.NullInputBehaviorCalledOnNullInput)).
		WithReturnResultsBehavior(string(sdk.ReturnResultsBehaviorVolatile)).
		WithComment("some other comment").
		WithImports(
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction.ModuleName + ".py"},
			sdk.NormalizedPath{StageLocation: "~", PathOnStage: tmpPythonFunction2.ModuleName + ".py"},
		).
		WithPackages("numpy", "pandas").
		WithExternalAccessIntegrations(externalAccessIntegration).
		WithSecrets(map[string]sdk.SchemaObjectIdentifier{
			"def": secretId2,
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionPython),
		Steps: []resource.TestStep{
			// CREATE WITH ALL
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionPythonResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasRuntimeVersionString("3.8").
						HasFunctionDefinitionString(definition).
						HasCommentString("some comment").
						HasFunctionLanguageString("PYTHON").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "secrets.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "external_access_integrations.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "packages.#", "2")),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
			// IMPORT
			{
				ResourceName:            functionModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_aggregate", "arguments.0.arg_data_type"},
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedFunctionPythonResource(t, id.FullyQualifiedName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_name", argName)),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_data_type", "NUMBER(38, 0)")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "arguments.0.arg_default_value", "")),
				),
			},
			// UPDATE WITHOUT RECREATION
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModelUpdateWithoutRecreation.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, functionModelUpdateWithoutRecreation),
				Check: assert.AssertThat(t,
					resourceassert.FunctionPythonResource(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanFalse).
						HasRuntimeVersionString("3.8").
						HasFunctionDefinitionString(definition).
						HasCommentString("some other comment").
						HasFunctionLanguageString("PYTHON").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_variable_name", "def")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "secrets.0.secret_id", secretId2.FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "external_access_integrations.0", externalAccessIntegration.Name())),
					assert.Check(resource.TestCheckResourceAttr(functionModelUpdateWithoutRecreation.ResourceReference(), "packages.#", "2")),
					resourceshowoutputassert.FunctionShowOutput(t, functionModelUpdateWithoutRecreation.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}
