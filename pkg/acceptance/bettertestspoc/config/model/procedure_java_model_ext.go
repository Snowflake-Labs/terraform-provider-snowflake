package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (f *ProcedureJavaModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureJavaModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func ProcedureJavaBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	procedureDefinition string,
) *ProcedureJavaModel {
	return ProcedureJava(resourceName, id.DatabaseName(), handler, id.Name(), returnType.ToSql(), "11", id.SchemaName(), "1.14.0").
		WithProcedureDefinition(procedureDefinition)
}

func ProcedureJavaBasicStaged(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	stageLocation string,
	pathOnStage string,
) *ProcedureJavaModel {
	return ProcedureJava(resourceName, id.DatabaseName(), handler, id.Name(), returnType.ToSql(), "11", id.SchemaName(), "1.14.0").
		WithImport(stageLocation, pathOnStage)
}

func (f *ProcedureJavaModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureJavaModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}

func (f *ProcedureJavaModel) WithArgumentWithDefaultValue(argName string, argDataType datatypes.DataType, value string) *ProcedureJavaModel {
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

func (f *ProcedureJavaModel) WithImport(stageLocation string, pathOnStage string) *ProcedureJavaModel {
	return f.WithImportsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}

func (f *ProcedureJavaModel) WithTargetPathParts(stageLocation string, pathOnStage string) *ProcedureJavaModel {
	return f.WithTargetPathValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"stage_location": tfconfig.StringVariable(stageLocation),
				"path_on_stage":  tfconfig.StringVariable(pathOnStage),
			},
		),
	)
}
