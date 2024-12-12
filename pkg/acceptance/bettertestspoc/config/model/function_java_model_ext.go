package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (f *FunctionJavaModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionJavaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func FunctionJavaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	functionDefinition string,
) *FunctionJavaModel {
	return FunctionJava(resourceName, id.DatabaseName(), handler, id.Name(), returnType.ToSql(), id.SchemaName()).WithFunctionDefinition(functionDefinition)
}

func FunctionJavaBasicStaged(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	stageLocation string,
	pathOnStage string,
) *FunctionJavaModel {
	return FunctionJava(resourceName, id.DatabaseName(), handler, id.Name(), returnType.ToSql(), id.SchemaName()).
		WithImport(stageLocation, pathOnStage)
}

func (f *FunctionJavaModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *FunctionJavaModel) WithArgumentWithDefaultValue(argName string, argDataType datatypes.DataType, value string) *FunctionJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":          tfconfig.StringVariable(argName),
				"arg_data_type":     tfconfig.StringVariable(argDataType.ToSql()),
				"arg_default_value": tfconfig.StringVariable(value),
			},
		),
	)
}

func (f *FunctionJavaModel) WithImport(stageLocation string, pathOnStage string) *FunctionJavaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *FunctionJavaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *FunctionJavaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}
