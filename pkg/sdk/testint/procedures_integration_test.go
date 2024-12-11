package testint

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1850370]: 'ExtendedIn' struct for procedures not support keyword "CLASS" now
// TODO [SNOW-1850370]: Call/CreateAndCall methods were not updated before V1 because we are not using them
func TestInt_Procedures(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	secretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()

	networkRule, networkRuleCleanup := testClientHelper().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	secret, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
	t.Cleanup(secretCleanup)

	externalAccessIntegration, externalAccessIntegrationCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t, networkRule.ID(), secret.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	tmpJavaProcedure := testClientHelper().CreateSampleJavaProcedureAndJar(t)
	tmpPythonFunction := testClientHelper().CreateSamplePythonFunctionAndModule(t)

	assertParametersSet := func(t *testing.T, procedureParametersAssert *objectparametersassert.ProcedureParametersAssert) {
		t.Helper()
		assertions.AssertThatObject(t, procedureParametersAssert.
			// TODO [SNOW-1850370]: every value end with invalid value [OFF] for parameter 'AUTO_EVENT_LOGGING'
			// HasAutoEventLogging(sdk.AutoEventLoggingTracing).
			HasEnableConsoleOutput(true).
			HasLogLevel(sdk.LogLevelWarn).
			HasMetricLevel(sdk.MetricLevelAll).
			HasTraceLevel(sdk.TraceLevelAlways),
		)
	}

	t.Run("create procedure for Java - inline minimal", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		definition := testClientHelper().Procedure.SampleJavaDefinition(t, className, funcName, argName)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

		request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedure, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, procedure).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, procedure.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, procedure.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Java - inline full", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		definition := testClientHelper().Procedure.SampleJavaDefinition(t, className, funcName, argName)
		jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
		targetPath := fmt.Sprintf("@~/%s", jarName)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0"),
			*sdk.NewProcedurePackageRequest("com.snowflake:telemetry:0.1.0"),
		}

		request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpJavaProcedure.JarLocation())}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithTargetPath(targetPath).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))
		t.Cleanup(testClientHelper().Stage.RemoveFromUserStageFunc(t, jarName))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			// TODO [SNOW-1850370]: apparently external access integrations and secrets are not filled out correctly for procedures
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaProcedure.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPath(targetPath).
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Java - staged minimal", func(t *testing.T) {
		dataType := tmpJavaProcedure.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaProcedure.JavaHandler()
		importPath := tmpJavaProcedure.JarLocation()
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0"),
			*sdk.NewProcedurePackageRequest("com.snowflake:telemetry:0.1.0"),
		}

		requestStaged := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(importPath)})

		err := client.Procedures.CreateForJava(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Java - staged full", func(t *testing.T) {
		dataType := tmpJavaProcedure.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaProcedure.JavaHandler()
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0"),
			*sdk.NewProcedurePackageRequest("com.snowflake:telemetry:0.1.0"),
		}

		requestStaged := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, handler).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpJavaProcedure.JarLocation())}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}})

		err := client.Procedures.CreateForJava(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaProcedure.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("11").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Javascript - inline minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "d"
		definition := testClientHelper().Procedure.SampleJavascriptDefinition(t, argName)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)

		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), dataType, definition).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument})

		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Javascript - inline full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "d"
		definition := testClientHelper().Procedure.SampleJavascriptDefinition(t, argName)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), dataType, definition).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNotNull(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithExecuteAs(sdk.ExecuteAsCaller).
			WithComment("comment")

		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("JAVASCRIPT").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Python - inline minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeNumber_36_2
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		funcName := "dump"
		definition := testClientHelper().Procedure.SamplePythonDefinition(t, funcName, argName)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("snowflake-snowpark-python==1.14.0"),
		}
		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, funcName).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")).
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(funcName).
			HasRuntimeVersion("3.8").
			HasPackages(`['snowflake-snowpark-python==1.14.0']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Python - inline full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeNumber_36_2
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		funcName := "dump"
		definition := testClientHelper().Procedure.SamplePythonDefinition(t, funcName, argName)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("snowflake-snowpark-python==1.14.0"),
			*sdk.NewProcedurePackageRequest("absl-py==0.10.0"),
		}

		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, funcName).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpPythonFunction.PythonModuleLocation())}).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithExecuteAs(sdk.ExecuteAsCaller).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL").
			HasLanguage("PYTHON").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(funcName).
			HasRuntimeVersion("3.8").
			HasPackages(`['snowflake-snowpark-python==1.14.0','absl-py==0.10.0']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Python - staged minimal", func(t *testing.T) {
		dataType := testdatatypes.DataTypeVarchar_100
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("snowflake-snowpark-python==1.14.0"),
		}
		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, tmpPythonFunction.PythonHandler()).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpPythonFunction.PythonModuleLocation())})

		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")).
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion("3.8").
			HasPackages(`['snowflake-snowpark-python==1.14.0']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Python - staged full", func(t *testing.T) {
		dataType := testdatatypes.DataTypeVarchar_100
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "i"
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("snowflake-snowpark-python==1.14.0"),
			*sdk.NewProcedurePackageRequest("absl-py==0.10.0"),
		}

		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, tmpPythonFunction.PythonHandler()).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpPythonFunction.PythonModuleLocation())}).
			WithExecuteAs(sdk.ExecuteAsCaller)

		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(strings.ReplaceAll(dataType.ToSql(), " ", "")+" NOT NULL").
			HasLanguage("PYTHON").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpPythonFunction.PythonModuleLocation())).
			HasHandler(tmpPythonFunction.PythonHandler()).
			HasRuntimeVersion("3.8").
			HasPackages(`['snowflake-snowpark-python==1.14.0','absl-py==0.10.0']`).
			HasTargetPathNil().
			HasInstalledPackagesNotEmpty().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Scala - inline minimal", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		definition := testClientHelper().Procedure.SampleScalaDefinition(t, className, funcName, argName)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

		request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, handler).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(`[]`).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Scala - inline full", func(t *testing.T) {
		className := "TestFunc"
		funcName := "echoVarchar"
		argName := "x"
		dataType := testdatatypes.DataTypeVarchar_100

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		definition := testClientHelper().Procedure.SampleScalaDefinition(t, className, funcName, argName)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		handler := fmt.Sprintf("%s.%s", className, funcName)
		jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
		targetPath := fmt.Sprintf("@~/%s", jarName)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0"),
			*sdk.NewProcedurePackageRequest("com.snowflake:telemetry:0.1.0"),
		}

		request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, handler).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpJavaProcedure.JarLocation())}).
			WithTargetPath(targetPath).
			WithExecuteAs(sdk.ExecuteAsCaller).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithProcedureDefinitionWrapped(definition)

		err := client.Procedures.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))
		t.Cleanup(testClientHelper().Stage.RemoveFromUserStageFunc(t, jarName))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SCALA").
			HasBody(definition).
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaProcedure.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPath(targetPath).
			HasInstalledPackagesNil().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Scala - staged minimal", func(t *testing.T) {
		dataType := tmpJavaProcedure.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		handler := tmpJavaProcedure.JavaHandler()
		importPath := tmpJavaProcedure.JarLocation()
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0")}

		requestStaged := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, handler).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(importPath)})

		err := client.Procedures.CreateForScala(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorCalledOnNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorVolatile)).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImports(fmt.Sprintf(`[%s]`, importPath)).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for Scala - staged full", func(t *testing.T) {
		dataType := tmpJavaProcedure.ArgType
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		argName := "x"
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		handler := tmpJavaProcedure.JavaHandler()

		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:1.14.0"),
			*sdk.NewProcedurePackageRequest("com.snowflake:telemetry:0.1.0"),
		}

		requestStaged := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, handler).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithComment("comment").
			WithExecuteAs(sdk.ExecuteAsCaller).
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecrets([]sdk.SecretReference{{VariableName: "abc", Name: secretId}}).
			WithImports([]sdk.ProcedureImportRequest{*sdk.NewProcedureImportRequest(tmpJavaProcedure.JarLocation())})

		err := client.Procedures.CreateForScala(ctx, requestStaged)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SCALA").
			HasBodyNil().
			HasNullHandling(string(sdk.NullInputBehaviorReturnsNullInput)).
			HasVolatility(string(sdk.ReturnResultsBehaviorImmutable)).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}).
			HasImports(fmt.Sprintf(`[%s]`, tmpJavaProcedure.JarLocation())).
			HasHandler(handler).
			HasRuntimeVersion("2.12").
			HasPackages(`[com.snowflake:snowpark:1.14.0,com.snowflake:telemetry:0.1.0]`).
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for SQL - inline minimal", func(t *testing.T) {
		argName := "x"
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Procedure.SampleSqlDefinition(t)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType)
		returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument})

		err := client.Procedures.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription(sdk.DefaultProcedureComment).
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(dataType.ToSql()).
			HasLanguage("SQL").
			HasBody(definition).
			HasNullHandlingNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("OWNER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("create procedure for SQL - inline full", func(t *testing.T) {
		argName := "x"
		dataType := testdatatypes.DataTypeFloat
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))

		definition := testClientHelper().Procedure.SampleSqlDefinition(t)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(dataType).
			WithNotNull(true)
		returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest(argName, dataType)
		request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCopyGrants(true).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithReturnResultsBehavior(sdk.ReturnResultsBehaviorImmutable).
			WithExecuteAs(sdk.ExecuteAsCaller).
			WithComment("comment")

		err := client.Procedures.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		function, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureFromObject(t, function).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasSchemaName(id.SchemaName()).
			HasIsBuiltin(false).
			HasIsAggregate(false).
			HasIsAnsi(false).
			HasMinNumArguments(1).
			HasMaxNumArguments(1).
			HasArgumentsOld([]sdk.DataType{sdk.LegacyDataTypeFrom(dataType)}).
			HasArgumentsRaw(fmt.Sprintf(`%[1]s(%[2]s) RETURN %[2]s`, function.ID().Name(), dataType.ToLegacyDataTypeSql())).
			HasDescription("comment").
			HasCatalogName(id.DatabaseName()).
			HasIsTableFunction(false).
			HasValidForClustering(false).
			HasIsSecure(false).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, function.ID()).
			HasSignature(fmt.Sprintf(`(%s %s)`, argName, dataType.ToLegacyDataTypeSql())).
			HasReturns(fmt.Sprintf(`%s NOT NULL`, dataType.ToSql())).
			HasLanguage("SQL").
			HasBody(definition).
			// TODO [SNOW-1348103]: null handling and volatility are not returned and is present in create syntax
			HasNullHandlingNil().
			HasVolatilityNil().
			HasVolatilityNil().
			HasExternalAccessIntegrationsNil().
			HasSecretsNil().
			HasImportsNil().
			HasHandlerNil().
			HasRuntimeVersionNil().
			HasPackagesNil().
			HasTargetPathNil().
			HasInstalledPackagesNil().
			HasExecuteAs("CALLER"),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Java: returns table", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#specifying-return-column-names-and-types
		name := "filter_by_role"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
			import com.snowflake.snowpark_java.*;
			public class Filter {
				public DataFrame filterByRole(Session session, String tableName, String role) {
					DataFrame table = session.table(tableName);
					DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
					return filteredRows;
				}
			}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("table_name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, "Filter.filterByRole").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Javascript", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript#basic-examples
		name := "stproc1"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeFloat)

		definition := `
				var sql_command = "INSERT INTO stproc_test_table1 (num_col1) VALUES (" + FLOAT_PARAM1 + ")";
				try {
					snowflake.execute (
						{sqlText: sql_command}
					);
					return "Succeeded."; // Return a success/error indicator.
				}
				catch (err)  {
					return "Failed: " + err; // Return a success/error indicator.
				}`
		argument := sdk.NewProcedureArgumentRequest("FLOAT_PARAM1", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), nil, definition).
			WithResultDataTypeOld(sdk.DataTypeString).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorStrict)).
			WithExecuteAs(*sdk.ExecuteAsPointer(sdk.ExecuteAsCaller))
		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Javascript: no arguments", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript#basic-examples
		name := "sp_pi"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name)

		definition := `return 3.1415926;`
		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), nil, definition).WithResultDataTypeOld(sdk.DataTypeFloat).WithNotNull(true).WithOrReplace(true)
		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Scala: returns result data type", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala#reading-a-dynamically-specified-file-with-snowflakefile
		name := "file_reader_scala_proc_snowflakefile"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR)

		definition := `
			import java.io.InputStream
			import java.nio.charset.StandardCharsets
			import com.snowflake.snowpark_java.types.SnowflakeFile
			import com.snowflake.snowpark_java.Session
			object FileReader {
				def execute(session: Session, fileName: String): String = {
					var input: InputStream = SnowflakeFile.newInstance(fileName).getInputStream()
					return new String(input.readAllBytes(), StandardCharsets.UTF_8)
				}
			}`
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVARCHAR)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest("input", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, "FileReader.execute").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Scala: returns table", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala#specifying-return-column-names-and-types
		name := "filter_by_role"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
			import com.snowflake.snowpark.functions._
			import com.snowflake.snowpark._
			object Filter {
				def filterByRole(session: Session, tableName: String, role: String): DataFrame = {
					val table = session.table(tableName)
					val filteredRows = table.filter(col("role") === role)
					return filteredRows
				}
			}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("table_name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, "Filter.filterByRole").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Python: returns result data type", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-python#running-concurrent-tasks-with-worker-processes
		name := "joblib_multiprocessing_proc"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeInt)

		definition := `
import joblib
from math import sqrt
def joblib_multiprocessing(session, i):
	result = joblib.Parallel(n_jobs=-1)(joblib.delayed(sqrt)(i ** 2) for i in range(10))
	return str(result)`

		dt := sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeString)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest("i", nil).WithArgDataTypeOld(sdk.DataTypeInt)
		packages := []sdk.ProcedurePackageRequest{
			*sdk.NewProcedurePackageRequest("snowflake-snowpark-python"),
			*sdk.NewProcedurePackageRequest("joblib"),
		}
		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, "joblib_multiprocessing").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for Python: returns table", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-python#specifying-return-column-names-and-types
		name := "filterByRole"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
from snowflake.snowpark.functions import col
def filter_by_role(session, table_name, role):
	df = session.table(table_name)
	return df.filter(col("role") == role)`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("table_name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("snowflake-snowpark-python")}
		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, "filter_by_role").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for SQL: returns result data type", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-snowflake-scripting
		name := "output_message"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR)

		definition := `
		BEGIN
			RETURN message;
		END;`

		dt := sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVARCHAR)
		returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*dt).WithNotNull(true)
		argument := sdk.NewProcedureArgumentRequest("message", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			// Suddenly this is erroring out, when it used to not have an problem. Must be an error with the Snowflake API.
			// Created issue in docs-discuss channel. https://snowflake.slack.com/archives/C6380540P/p1707511734666249
			// Error:      	Received unexpected error:
			// 001003 (42000): SQL compilation error:
			// syntax error line 1 at position 210 unexpected 'NULL'.
			// syntax error line 1 at position 215 unexpected 'ON'.
			// WithNullInputBehavior(sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument})
		err := client.Procedures.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	// TODO [SNOW-1348103]: adjust or remove
	t.Run("create procedure for SQL: returns table", func(t *testing.T) {
		t.Skipf("Skipped for now; left as inspiration for resource rework as part of SNOW-1348103")

		name := "find_invoice_by_id"
		id := testClientHelper().Ids.NewSchemaObjectIdentifierWithArguments(name, sdk.DataTypeVARCHAR)

		definition := `
		DECLARE
			res RESULTSET DEFAULT (SELECT * FROM invoices WHERE id = :id);
		BEGIN
			RETURN TABLE(res);
		END;`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld("INTEGER")
		column2 := sdk.NewProcedureColumnRequest("price", nil).WithColumnDataTypeOld("NUMBER(12,2)")
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2})
		returns := sdk.NewProcedureSQLReturnsRequest().WithTable(*returnsTable)
		argument := sdk.NewProcedureArgumentRequest("id", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			// SNOW-1051627 todo: uncomment once null input behavior working again
			// WithNullInputBehavior(sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorReturnsNullInput)).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument})
		err := client.Procedures.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, id))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("show parameters", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterLogLevel, sdk.Object{ObjectType: sdk.ObjectTypeProcedure, Name: id})
		require.NoError(t, err)
		assert.Equal(t, string(sdk.LogLevelOff), param.Value)

		parameters, err := client.Parameters.ShowParameters(ctx, &sdk.ShowParametersOptions{
			In: &sdk.ParametersIn{
				Procedure: id,
			},
		})
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParametersPrefetched(t, id, parameters).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		// check that ShowParameters on procedure level works too
		parameters, err = client.Procedures.ShowParameters(ctx, id)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParametersPrefetched(t, id, parameters).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("alter procedure: rename", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		nid := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(id.ArgumentDataTypes()...)

		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithRenameTo(nid.SchemaObjectId()))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Procedure.DropProcedureFunc(t, nid))

		_, err = client.Procedures.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Procedures.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter procedure: set and unset all for Java", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateJava(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		assertions.AssertThatObject(t, objectassert.Procedure(t, id).
			HasName(id.Name()).
			HasDescription(sdk.DefaultProcedureComment),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		request := sdk.NewAlterProcedureRequest(id).WithSet(*sdk.NewProcedureSetRequest().
			WithExternalAccessIntegrations([]sdk.AccountObjectIdentifier{externalAccessIntegration}).
			WithSecretsList(*sdk.NewSecretsListRequest([]sdk.SecretReference{{VariableName: "abc", Name: secretId}})).
			// TODO [SNOW-1850370]: every value end with invalid value [OFF] for parameter 'AUTO_EVENT_LOGGING'
			// WithAutoEventLogging(sdk.AutoEventLoggingAll).
			WithEnableConsoleOutput(true).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithComment("new comment"),
		)

		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Procedure(t, id).
			HasName(id.Name()).
			HasDescription("new comment"),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExactlyExternalAccessIntegrations(externalAccessIntegration).
			HasExactlySecrets(map[string]sdk.SchemaObjectIdentifier{"abc": secretId}),
		)

		assertParametersSet(t, objectparametersassert.ProcedureParameters(t, id))

		unsetRequest := sdk.NewAlterProcedureRequest(id).WithUnset(*sdk.NewProcedureUnsetRequest().
			WithExternalAccessIntegrations(true).
			// WithAutoEventLogging(true).
			WithEnableConsoleOutput(true).
			WithLogLevel(true).
			WithMetricLevel(true).
			WithTraceLevel(true).
			WithComment(true),
		)

		err = client.Procedures.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Procedure(t, id).
			HasName(id.Name()).
			HasDescription(sdk.DefaultProcedureComment).
			// both nil, because they are always nil in SHOW for procedures
			HasExternalAccessIntegrationsNil().
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExternalAccessIntegrationsNil().
			// TODO [SNOW-1850370]: apparently UNSET external access integrations cleans out secrets in the describe but leaves it in SHOW
			HasSecretsNil(),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		unsetSecretsRequest := sdk.NewAlterProcedureRequest(id).WithSet(*sdk.NewProcedureSetRequest().
			WithSecretsList(*sdk.NewSecretsListRequest([]sdk.SecretReference{})),
		)

		err = client.Procedures.Alter(ctx, unsetSecretsRequest)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasSecretsNil(),
		)
	})

	t.Run("alter procedure: set and unset all for SQL", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		request := sdk.NewAlterProcedureRequest(id).WithSet(*sdk.NewProcedureSetRequest().
			// WithAutoEventLogging(sdk.AutoEventLoggingTracing).
			WithEnableConsoleOutput(true).
			WithLogLevel(sdk.LogLevelWarn).
			WithMetricLevel(sdk.MetricLevelAll).
			WithTraceLevel(sdk.TraceLevelAlways).
			WithComment("new comment"),
		)

		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Procedure(t, id).
			HasName(id.Name()).
			HasDescription("new comment"),
		)

		assertParametersSet(t, objectparametersassert.ProcedureParameters(t, id))

		unsetRequest := sdk.NewAlterProcedureRequest(id).WithUnset(*sdk.NewProcedureUnsetRequest().
			// WithAutoEventLogging(true).
			WithEnableConsoleOutput(true).
			WithLogLevel(true).
			WithMetricLevel(true).
			WithTraceLevel(true).
			WithComment(true),
		)

		err = client.Procedures.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.Procedure(t, id).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDescription(sdk.DefaultProcedureComment),
		)

		assertions.AssertThatObject(t, objectparametersassert.ProcedureParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("alter procedure: set execute as", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExecuteAs("OWNER"),
		)

		err := client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithExecuteAs(*sdk.ExecuteAsPointer(sdk.ExecuteAsCaller)))
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExecuteAs("CALLER"),
		)

		err = client.Procedures.Alter(ctx, sdk.NewAlterProcedureRequest(id).WithExecuteAs(*sdk.ExecuteAsPointer(sdk.ExecuteAsOwner)))
		require.NoError(t, err)

		assertions.AssertThatObject(t, objectassert.ProcedureDetails(t, id).
			HasExecuteAs("OWNER"),
		)
	})

	t.Run("show procedure: without like", func(t *testing.T) {
		p1, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		p2, pCleanup2 := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup2)

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(procedures), 1)
		require.Contains(t, procedures, *p1)
		require.Contains(t, procedures, *p2)
	})

	t.Run("show procedure: with like", func(t *testing.T) {
		p1, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		p2, pCleanup2 := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup2)

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest().WithLike(sdk.Like{Pattern: &p1.Name}))
		require.NoError(t, err)

		require.Equal(t, 1, len(procedures))
		require.Contains(t, procedures, *p1)
		require.NotContains(t, procedures, *p2)
	})

	t.Run("show procedure: no matches", func(t *testing.T) {
		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest().
			WithIn(sdk.ExtendedIn{In: sdk.In{Schema: testClientHelper().Ids.SchemaId()}}).
			WithLike(sdk.Like{Pattern: sdk.String(NonExistingSchemaObjectIdentifier.Name())}))
		require.NoError(t, err)
		require.Equal(t, 0, len(procedures))
	})

	t.Run("describe procedure: for SQL", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateSql(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		details, err := client.Procedures.Describe(ctx, id)
		require.NoError(t, err)
		assert.Len(t, details, 5)

		pairs := make(map[string]*string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		assert.Equal(t, "(x FLOAT)", *pairs["signature"])
		assert.Equal(t, "FLOAT", *pairs["returns"])
		assert.Equal(t, "SQL", *pairs["language"])
		assert.Equal(t, "\nBEGIN\n\tRETURN 3.141592654::FLOAT;\nEND;\n", *pairs["body"])
		assert.Equal(t, "OWNER", *pairs["execute as"])
	})

	t.Run("describe procedure: for Java", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateJava(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		details, err := client.Procedures.Describe(ctx, id)
		require.NoError(t, err)
		assert.Len(t, details, 12)

		pairs := make(map[string]*string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		assert.Equal(t, "(x VARCHAR)", *pairs["signature"])
		assert.Equal(t, "VARCHAR(100)", *pairs["returns"])
		assert.Equal(t, "JAVA", *pairs["language"])
		assert.NotEmpty(t, *pairs["body"])
		assert.Equal(t, string(sdk.NullInputBehaviorCalledOnNullInput), *pairs["null handling"])
		assert.Equal(t, string(sdk.VolatileTableKind), *pairs["volatility"])
		assert.Nil(t, pairs["external_access_integration"])
		assert.Nil(t, pairs["secrets"])
		assert.Equal(t, "[]", *pairs["imports"])
		assert.Equal(t, "TestFunc.echoVarchar", *pairs["handler"])
		assert.Equal(t, "11", *pairs["runtime_version"])
		assert.Equal(t, "OWNER", *pairs["execute as"])
	})

	t.Run("drop procedure for SQL", func(t *testing.T) {
		p, pCleanup := testClientHelper().Procedure.CreateJava(t)
		t.Cleanup(pCleanup)
		id := p.ID()

		err := client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id))
		require.NoError(t, err)
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		dataType := testdatatypes.DataTypeFloat
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(id1.Name(), schema.ID(), sdk.LegacyDataTypeFrom(dataType))

		_, pCleanup1 := testClientHelper().Procedure.CreateSqlWithIdentifierAndArgument(t, id1.SchemaObjectId(), dataType, testClientHelper().Procedure.SampleSqlDefinition(t))
		t.Cleanup(pCleanup1)
		_, pCleanup2 := testClientHelper().Procedure.CreateSqlWithIdentifierAndArgument(t, id2.SchemaObjectId(), dataType, testClientHelper().Procedure.SampleSqlDefinition(t))
		t.Cleanup(pCleanup2)

		e1, err := client.Procedures.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Procedures.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show procedure by id - same name, different arguments", func(t *testing.T) {
		dataType := testdatatypes.DataTypeFloat
		name := testClientHelper().Ids.Alpha()

		id1 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.LegacyDataTypeFrom(dataType))
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeVARCHAR)

		e := testClientHelper().Procedure.CreateWithIdentifier(t, id1)
		testClientHelper().Procedure.CreateWithIdentifier(t, id2)

		es, err := client.Procedures.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	// This test shows behavior of detailed types (e.g. VARCHAR(20) and NUMBER(10, 0)) on Snowflake side for procedures.
	// For SHOW, data type is generalized both for argument and return type (to e.g. VARCHAR and NUMBER).
	// FOR DESCRIBE, data type is generalized for argument and works weirdly for the return type: type is generalized to the canonical one, but we also get the attributes.
	for _, tc := range []string{
		"NUMBER(36, 5)",
		"NUMBER(36)",
		"NUMBER",
		"DECIMAL",
		"INTEGER",
		"FLOAT",
		"DOUBLE",
		"VARCHAR",
		"VARCHAR(20)",
		"CHAR",
		"CHAR(10)",
		"TEXT",
		"BINARY",
		"BINARY(1000)",
		"VARBINARY",
		"BOOLEAN",
		"DATE",
		"DATETIME",
		"TIME",
		"TIMESTAMP_LTZ",
		"TIMESTAMP_NTZ",
		"TIMESTAMP_TZ",
		"VARIANT",
		"OBJECT",
		"ARRAY",
		"GEOGRAPHY",
		"GEOMETRY",
		"VECTOR(INT, 16)",
		"VECTOR(FLOAT, 8)",
	} {
		tc := tc
		t.Run(fmt.Sprintf("procedure returns non detailed data types of arguments for %s", tc), func(t *testing.T) {
			procName := "add"
			argName := "A"
			dataType, err := datatypes.ParseDataType(tc)
			require.NoError(t, err)
			args := []sdk.ProcedureArgumentRequest{
				*sdk.NewProcedureArgumentRequest(argName, dataType),
			}
			oldDataType := sdk.LegacyDataTypeFrom(dataType)
			idWithArguments := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(oldDataType)

			packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("snowflake-snowpark-python")}
			definition := fmt.Sprintf("def add(%[1]s): %[1]s", argName)

			err = client.Procedures.CreateForPython(ctx, sdk.NewCreateForPythonProcedureRequest(
				idWithArguments.SchemaObjectId(),
				*sdk.NewProcedureReturnsRequest().WithResultDataType(*sdk.NewProcedureReturnsResultDataTypeRequest(dataType)),
				"3.8",
				packages,
				procName,
			).
				WithArguments(args).
				WithProcedureDefinitionWrapped(definition),
			)
			require.NoError(t, err)

			procedure, err := client.Procedures.ShowByID(ctx, idWithArguments)
			require.NoError(t, err)
			assert.Equal(t, []sdk.DataType{oldDataType}, procedure.ArgumentsOld)
			assert.Equal(t, fmt.Sprintf("%[1]s(%[2]s) RETURN %[2]s", idWithArguments.Name(), oldDataType), procedure.ArgumentsRaw)

			details, err := client.Procedures.Describe(ctx, idWithArguments)
			require.NoError(t, err)
			pairs := make(map[string]string)
			for _, detail := range details {
				pairs[detail.Property] = *detail.Value
			}
			assert.Equal(t, fmt.Sprintf("(%s %s)", argName, oldDataType), pairs["signature"])
			assert.Equal(t, dataType.Canonical(), pairs["returns"])
		})
	}
}

func TestInt_CallProcedure(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseId, schemaId := testClientHelper().Ids.DatabaseId(), testClientHelper().Ids.SchemaId()
	cleanupProcedureHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createTableHandle := func(t *testing.T, table sdk.SchemaObjectIdentifier) {
		t.Helper()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE OR REPLACE TABLE %s (id NUMBER, name VARCHAR, role VARCHAR)`, table.FullyQualifiedName()))
		require.NoError(t, err)
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`INSERT INTO %s (id, name, role) VALUES (1, 'Alice', 'op'), (2, 'Bob', 'dev'), (3, 'Cindy', 'dev')`, table.FullyQualifiedName()))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP TABLE %s`, table.FullyQualifiedName()))
			require.NoError(t, err)
		})
	}

	// create a employees table
	tid := sdk.NewSchemaObjectIdentifier(databaseId.Name(), schemaId.Name(), "employees")
	createTableHandle(t, tid)

	createProcedureForSQLHandle := func(t *testing.T, cleanup bool) *sdk.Procedure {
		t.Helper()

		definition := `
		BEGIN
			RETURN MESSAGE;
		END;`
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR)
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVARCHAR)
		returns := sdk.NewProcedureSQLReturnsRequest().WithResultDataType(*dt).WithNotNull(true)
		argument := sdk.NewProcedureArgumentRequest("MESSAGE", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForSQLProcedureRequestDefinitionWrapped(id.SchemaObjectId(), *returns, definition).
			WithSecure(true).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithExecuteAs(*sdk.ExecuteAsPointer(sdk.ExecuteAsCaller))
		err := client.Procedures.CreateForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupProcedureHandle(id))
		}
		procedure, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)
		return procedure
	}

	t.Run("call procedure for SQL: argument positions", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)
		err := client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(f.ID().SchemaObjectId()).WithCallArguments([]string{"'hi'"}))
		require.NoError(t, err)
	})

	t.Run("call procedure for SQL: argument names", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)
		err := client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(f.ID().SchemaObjectId()).WithCallArguments([]string{"MESSAGE => 'hi'"}))
		require.NoError(t, err)
	})

	t.Run("call procedure for Java: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#omitting-return-column-names-and-types
		name := "filter_by_role"
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseId.Name(), schemaId.Name(), name, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
		import com.snowflake.snowpark_java.*;
		public class Filter {
			public DataFrame filterByRole(Session session, String name, String role) {
				DataFrame table = session.table(name);
				DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
				return filteredRows;
			}
		}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		request := sdk.NewCreateForJavaProcedureRequest(id.SchemaObjectId(), *returns, "11", packages, "Filter.filterByRole").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id))

		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		err = client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(id.SchemaObjectId()).WithCallArguments(ca))
		require.NoError(t, err)
	})

	t.Run("call procedure for Scala: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala#omitting-return-column-names-and-types
		name := "filter_by_role"
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseId.Name(), schemaId.Name(), name, sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
		import com.snowflake.snowpark.functions._
		import com.snowflake.snowpark._

		object Filter {
			def filterByRole(session: Session, name: String, role: String): DataFrame = {
				val table = session.table(name)
				val filteredRows = table.filter(col("role") === role)
				return filteredRows
			}
		}`
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		request := sdk.NewCreateForScalaProcedureRequest(id.SchemaObjectId(), *returns, "2.12", packages, "Filter.filterByRole").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id))

		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		err = client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(id.SchemaObjectId()).WithCallArguments(ca))
		require.NoError(t, err)
	})

	t.Run("call procedure for Javascript", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript#basic-examples
		name := "stproc1"
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseId.Name(), schemaId.Name(), name, sdk.DataTypeFloat)

		definition := `
		var sql_command = "INSERT INTO stproc_test_table1 (num_col1) VALUES (" + FLOAT_PARAM1 + ")";
		try {
			snowflake.execute (
				{sqlText: sql_command}
			);
			return "Succeeded."; // Return a success/error indicator.
		}
		catch (err)  {
			return "Failed: " + err; // Return a success/error indicator.
		}`
		arg := sdk.NewProcedureArgumentRequest("FLOAT_PARAM1", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), nil, definition).
			WithResultDataTypeOld(sdk.DataTypeString).
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorStrict)).
			WithExecuteAs(*sdk.ExecuteAsPointer(sdk.ExecuteAsOwner))
		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id))

		err = client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(id.SchemaObjectId()).WithCallArguments([]string{"5.14::FLOAT"}))
		require.NoError(t, err)
	})

	t.Run("call procedure for Javascript: no arguments", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript#basic-examples
		name := "sp_pi"
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseId.Name(), schemaId.Name(), name)

		definition := `return 3.1415926;`
		request := sdk.NewCreateForJavaScriptProcedureRequestDefinitionWrapped(id.SchemaObjectId(), nil, definition).WithResultDataTypeOld(sdk.DataTypeFloat).WithNotNull(true).WithOrReplace(true)
		err := client.Procedures.CreateForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id))

		err = client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(id.SchemaObjectId()))
		require.NoError(t, err)
	})

	t.Run("call procedure for Python: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-python#omitting-return-column-names-and-types
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseId.Name(), schemaId.Name(), "filterByRole", sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR)

		definition := `
from snowflake.snowpark.functions import col
def filter_by_role(session, name, role):
	df = session.table(name)
	return df.filter(col("role") == role)`
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("snowflake-snowpark-python")}
		request := sdk.NewCreateForPythonProcedureRequest(id.SchemaObjectId(), *returns, "3.8", packages, "filter_by_role").
			WithOrReplace(true).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinitionWrapped(definition)
		err := client.Procedures.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id))

		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		err = client.Procedures.Call(ctx, sdk.NewCallProcedureRequest(id.SchemaObjectId()).WithCallArguments(ca))
		require.NoError(t, err)
	})
}

func TestInt_CreateAndCallProcedures(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseId, schemaId := testClientHelper().Ids.DatabaseId(), testClientHelper().Ids.SchemaId()
	createTableHandle := func(t *testing.T, table sdk.SchemaObjectIdentifier) {
		t.Helper()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE OR REPLACE TABLE %s (id NUMBER, name VARCHAR, role VARCHAR)`, table.FullyQualifiedName()))
		require.NoError(t, err)
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`INSERT INTO %s (id, name, role) VALUES (1, 'Alice', 'op'), (2, 'Bob', 'dev'), (3, 'Cindy', 'dev')`, table.FullyQualifiedName()))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP TABLE %s`, table.FullyQualifiedName()))
			require.NoError(t, err)
		})
	}

	// create a employees table
	tid := sdk.NewSchemaObjectIdentifier(databaseId.Name(), schemaId.Name(), "employees")
	createTableHandle(t, tid)

	t.Run("create and call procedure for Java: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#omitting-return-column-names-and-types
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("filter_by_role")

		definition := `
		import com.snowflake.snowpark_java.*;
		public class Filter {
			public DataFrame filterByRole(Session session, String name, String role) {
				DataFrame table = session.table(name);
				DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
				return filteredRows;
			}
		}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		request := sdk.NewCreateAndCallForJavaProcedureRequest(name, *returns, "11", packages, "Filter.filterByRole", name).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinition(definition).
			WithCallArguments(ca)
		err := client.Procedures.CreateAndCallForJava(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for Scala: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-scala#omitting-return-column-names-and-types
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("filter_by_role")

		definition := `
		import com.snowflake.snowpark.functions._
		import com.snowflake.snowpark._

		object Filter {
			def filterByRole(session: Session, name: String, role: String): DataFrame = {
				val table = session.table(name)
				val filteredRows = table.filter(col("role") === role)
				return filteredRows
			}
		}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}
		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		request := sdk.NewCreateAndCallForScalaProcedureRequest(name, *returns, "2.12", packages, "Filter.filterByRole", name).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinition(definition).
			WithCallArguments(ca)
		err := client.Procedures.CreateAndCallForScala(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for Javascript", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-javascript#basic-examples
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("stproc1")

		definition := `
		var sql_command = "INSERT INTO stproc_test_table1 (num_col1) VALUES (" + FLOAT_PARAM1 + ")";
		try {
			snowflake.execute (
				{sqlText: sql_command}
			);
			return "Succeeded."; // Return a success/error indicator.
		}
		catch (err)  {
			return "Failed: " + err; // Return a success/error indicator.
		}`
		arg := sdk.NewProcedureArgumentRequest("FLOAT_PARAM1", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request := sdk.NewCreateAndCallForJavaScriptProcedureRequest(name, nil, definition, name).
			WithResultDataTypeOld(sdk.DataTypeString).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorStrict)).
			WithCallArguments([]string{"5.14::FLOAT"})
		err := client.Procedures.CreateAndCallForJavaScript(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for Javascript: no arguments", func(t *testing.T) {
		// https://docs.snowflake.com/en/sql-reference/sql/create-procedure#examples
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("sp_pi")

		definition := `return 3.1415926;`
		request := sdk.NewCreateAndCallForJavaScriptProcedureRequest(name, nil, definition, name).WithResultDataTypeOld(sdk.DataTypeFloat).WithNotNull(true)
		err := client.Procedures.CreateAndCallForJavaScript(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for SQL: argument positions", func(t *testing.T) {
		definition := `
		BEGIN
			RETURN MESSAGE;
		END;`

		name := testClientHelper().Ids.RandomAccountObjectIdentifier()
		dt := sdk.NewProcedureReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVARCHAR)
		returns := sdk.NewProcedureReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewProcedureArgumentRequest("MESSAGE", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateAndCallForSQLProcedureRequest(name, *returns, definition, name).
			WithArguments([]sdk.ProcedureArgumentRequest{*argument}).
			WithCallArguments([]string{"message => 'hi'"})
		err := client.Procedures.CreateAndCallForSQL(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for Python: returns table", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-python#omitting-return-column-names-and-types
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("filterByRole")
		definition := `
from snowflake.snowpark.functions import col
def filter_by_role(session, name, role):
	df = session.table(name)
	return df.filter(col("role") == role)`
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("snowflake-snowpark-python")}
		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		request := sdk.NewCreateAndCallForPythonProcedureRequest(name, *returns, "3.8", packages, "filter_by_role", name).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinition(definition).
			WithCallArguments(ca)
		err := client.Procedures.CreateAndCallForPython(ctx, request)
		require.NoError(t, err)
	})

	t.Run("create and call procedure for Java: returns table and with clause", func(t *testing.T) {
		// https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#omitting-return-column-names-and-types
		// TODO [SNOW-1348106]: make random with procedures rework
		name := sdk.NewAccountObjectIdentifier("filter_by_role")
		definition := `
		import com.snowflake.snowpark_java.*;
		public class Filter {
			public DataFrame filterByRole(Session session, String name, String role) {
				DataFrame table = session.table(name);
				DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
				return filteredRows;
			}
		}`
		column1 := sdk.NewProcedureColumnRequest("id", nil).WithColumnDataTypeOld(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest("name", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest("role", nil).WithColumnDataTypeOld(sdk.DataTypeVARCHAR)
		returnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		returns := sdk.NewProcedureReturnsRequest().WithTable(*returnsTable)
		arg1 := sdk.NewProcedureArgumentRequest("name", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		arg2 := sdk.NewProcedureArgumentRequest("role", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		packages := []sdk.ProcedurePackageRequest{*sdk.NewProcedurePackageRequest("com.snowflake:snowpark:latest")}

		ca := []string{fmt.Sprintf(`'%s'`, tid.FullyQualifiedName()), "'dev'"}
		// TODO [SNOW-1348106]: make random with procedures rework
		cte := sdk.NewAccountObjectIdentifier("records")
		statement := fmt.Sprintf(`(SELECT name, role FROM %s WHERE name = 'Bob')`, tid.FullyQualifiedName())
		clause := sdk.NewProcedureWithClauseRequest(cte, statement).WithCteColumns([]string{"name", "role"})
		request := sdk.NewCreateAndCallForJavaProcedureRequest(name, *returns, "11", packages, "Filter.filterByRole", name).
			WithArguments([]sdk.ProcedureArgumentRequest{*arg1, *arg2}).
			WithProcedureDefinition(definition).
			WithWithClause(*clause).
			WithCallArguments(ca)
		err := client.Procedures.CreateAndCallForJava(ctx, request)
		require.NoError(t, err)
	})
}
