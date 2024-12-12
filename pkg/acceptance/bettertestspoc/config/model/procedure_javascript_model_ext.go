package model

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (f *ProcedureJavascriptModel) MarshalJSON() ([]byte, error) {
	type Alias ProcedureJavascriptModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func ProcedureJavascriptBasicInline(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	procedureDefinition string,
) *ProcedureJavascriptModel {
	return ProcedureJavascript(resourceName, id.DatabaseName(), id.Name(), procedureDefinition, returnType.ToSql(), id.SchemaName())
}

func (f *ProcedureJavascriptModel) WithArgument(argName string, argDataType datatypes.DataType) *ProcedureJavascriptModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}
