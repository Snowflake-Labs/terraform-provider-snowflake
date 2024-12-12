package resources_test

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-1348103]: test import
// TODO [SNOW-1348103]: test external changes
// TODO [SNOW-1348103]: test changes of attributes separately

func TestAcc_FunctionJava_InlineBasic(t *testing.T) {
	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)
	idWithChangedNameButTheSameDataType := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	functionModelRenamed := model.FunctionJavaBasicInline("w", idWithChangedNameButTheSameDataType, dataType, handler, definition).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasImportsLength(0).
						HasTargetPathEmpty().
						HasNoRuntimeVersion().
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", "")),
				),
			},
			// RENAME
			{
				Config: config.FromModels(t, functionModelRenamed),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModelRenamed.ResourceReference()).
						HasNameString(idWithChangedNameButTheSameDataType.Name()).
						HasFullyQualifiedNameString(idWithChangedNameButTheSameDataType.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineEmptyArgs(t *testing.T) {
	className := "TestFunc"
	funcName := "echoVarchar"
	returnDataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes()

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinitionNoArgs(t, className, funcName)

	functionModel := model.FunctionJavaBasicInline("w", id, returnDataType, handler, definition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineBasicDefaultArg(t *testing.T) {
	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	defaultValue := "'hello'"

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgumentWithDefaultValue(argName, dataType, defaultValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_name", argName)),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_data_type", dataType.ToSql())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "arguments.0.arg_default_value", defaultValue)),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_InlineFull(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)
	// TODO [SNOW-1850370]: extract to helper
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithTargetPathParts(stage.ID().FullyQualifiedName(), jarName).
		WithRuntimeVersion("11")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasImportsLength(0).
						HasRuntimeVersionString("11").
						HasFunctionDefinitionString(definition).
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "target_path.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "target_path.0.path_on_stage", jarName)),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_StagedBasic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	tmpJavaFunction := acc.TestClient().CreateSampleJavaFunctionAndJarOnStage(t, stage)

	dataType := tmpJavaFunction.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaFunction.JavaHandler()

	functionModel := model.FunctionJavaBasicStaged("w", id, dataType, handler, stage.ID().FullyQualifiedName(), tmpJavaFunction.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// CREATE BASIC
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).
						HasNameString(id.Name()).
						HasIsSecureString(r.BooleanDefault).
						HasCommentString(sdk.DefaultFunctionComment).
						HasImportsLength(1).
						HasNoFunctionDefinition().
						HasFunctionLanguageString("JAVA").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "imports.0.stage_location", stage.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(functionModel.ResourceReference(), "imports.0.path_on_stage", tmpJavaFunction.JarName)),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).
						HasIsSecure(false),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_AllParameters(t *testing.T) {
	className := "TestFunc"
	funcName := "echoVarchar"
	argName := "x"
	dataType := testdatatypes.DataTypeVarchar_100
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := acc.TestClient().Function.SampleJavaDefinition(t, className, funcName, argName)

	functionModel := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType)
	functionModelWithAllParametersSet := model.FunctionJavaBasicInline("w", id, dataType, handler, definition).
		WithArgument(argName, dataType).
		WithEnableConsoleOutput(true).
		WithLogLevel(string(sdk.LogLevelWarn)).
		WithMetricLevel(string(sdk.MetricLevelAll)).
		WithTraceLevel(string(sdk.TraceLevelAlways))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.FunctionResourceParameters(t, functionModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName: functionModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedFunctionResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, functionModelWithAllParametersSet),
				Check: assert.AssertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.FunctionResourceParameters(t, functionModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// import when all parameters set
			{
				ResourceName: functionModelWithAllParametersSet.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceparametersassert.ImportedFunctionResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
			// unset all the parameters
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.FunctionResourceParameters(t, functionModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// destroy
			{
				Config:  config.FromModels(t, functionModel),
				Destroy: true,
			},
			// create with all parameters set
			{
				Config: config.FromModels(t, functionModelWithAllParametersSet),
				Check: assert.AssertThat(t,
					objectparametersassert.FunctionParameters(t, id).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
					resourceparametersassert.FunctionResourceParameters(t, functionModelWithAllParametersSet.ResourceReference()).
						HasEnableConsoleOutput(true).
						HasLogLevel(sdk.LogLevelWarn).
						HasMetricLevel(sdk.MetricLevelAll).
						HasTraceLevel(sdk.TraceLevelAlways),
				),
			},
		},
	})
}

func TestAcc_FunctionJava_handleExternalLanguageChange(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tmpJavaFunction := acc.TestClient().CreateSampleJavaFunctionAndJarOnUserStage(t)

	dataType := tmpJavaFunction.ArgType
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithArgumentsNewDataTypes(dataType)

	argName := "x"
	handler := tmpJavaFunction.JavaHandler()

	functionModel := model.FunctionJavaBasicStaged("w", id, dataType, handler, "~", tmpJavaFunction.JarName).
		WithArgument(argName, dataType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.FunctionJava),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, functionModel),
				Check: assert.AssertThat(t,
					objectassert.Function(t, id).HasLanguage("JAVA"),
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).HasNameString(id.Name()).HasFunctionLanguageString("JAVA"),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).HasLanguage("JAVA"),
				),
			},
			// change type externally by creating a new function with the exact same id but using different language
			{
				PreConfig: func() {
					acc.TestClient().Function.DropFunctionFunc(t, id)()
					acc.TestClient().Function.CreateScalaStaged(t, id, dataType, tmpJavaFunction.JarLocation(), handler)
					objectassert.Function(t, id).HasLanguage("SCALA")
				},
				Config: config.FromModels(t, functionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(functionModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t,
					objectassert.Function(t, id).HasLanguage("JAVA"),
					resourceassert.FunctionJavaResource(t, functionModel.ResourceReference()).HasNameString(id.Name()).HasFunctionLanguageString("JAVA"),
					resourceshowoutputassert.FunctionShowOutput(t, functionModel.ResourceReference()).HasLanguage("JAVA"),
				),
			},
		},
	})
}
