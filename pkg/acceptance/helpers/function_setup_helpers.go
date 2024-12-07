package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827324]: add TestClient ref to each specific client, so that we enhance specific client and not the base one
func (c *TestClient) CreateSampleJavaFunctionAndJar(t *testing.T) *TmpFunction {
	t.Helper()
	ctx := context.Background()

	className := fmt.Sprintf("TestClassAbc%s", random.AlphaLowerN(3))
	funcName := fmt.Sprintf("echoVarchar%s", random.AlphaLowerN(3))
	argName := fmt.Sprintf("arg%s", random.AlphaLowerN(3))
	dataType := testdatatypes.DataTypeVarchar_100

	id := c.Ids.RandomSchemaObjectIdentifierWithArguments(sdk.LegacyDataTypeFrom(dataType))
	argument := sdk.NewFunctionArgumentRequest(argName, testdatatypes.DataTypeVarchar_100)
	dt := sdk.NewFunctionReturnsResultDataTypeRequest(testdatatypes.DataTypeVarchar_100)
	returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
	handler := fmt.Sprintf("%s.%s", className, funcName)
	definition := c.Function.SampleJavaDefinition(t, className, funcName, argName)
	jarName := fmt.Sprintf("tf-%d-%s.jar", time.Now().Unix(), random.AlphaN(5))
	targetPath := fmt.Sprintf("@~/%s", jarName)

	request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, handler).
		WithArguments([]sdk.FunctionArgumentRequest{*argument}).
		WithTargetPath(targetPath).
		WithFunctionDefinitionWrapped(definition)

	err := c.context.client.Functions.CreateForJava(ctx, request)
	require.NoError(t, err)
	t.Cleanup(c.Function.DropFunctionFunc(t, id))
	t.Cleanup(c.Stage.RemoveFromUserStageFunc(t, jarName))

	return &TmpFunction{
		FunctionId: id,
		ClassName:  className,
		FuncName:   funcName,
		ArgName:    argName,
		ArgType:    dataType,
		JarName:    jarName,
	}
}

type TmpFunction struct {
	FunctionId sdk.SchemaObjectIdentifierWithArguments
	ClassName  string
	FuncName   string
	ArgName    string
	ArgType    datatypes.DataType
	JarName    string
}

func (f *TmpFunction) JarLocation() string {
	return fmt.Sprintf("@~/%s", f.JarName)
}

func (f *TmpFunction) Handler() string {
	return fmt.Sprintf("%s.%s", f.ClassName, f.FuncName)
}
