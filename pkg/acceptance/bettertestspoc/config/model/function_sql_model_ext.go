package model

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (f *FunctionSqlModel) MarshalJSON() ([]byte, error) {
	type Alias FunctionSqlModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(f),
		DependsOn: f.DependsOn(),
	})
}

func FunctionSqlBasicInline(resourceName string, id sdk.SchemaObjectIdentifierWithArguments, functionDefinition string, returnType string) *FunctionSqlModel {
	f := &FunctionSqlModel{ResourceModelMeta: config.Meta(resourceName, resources.FunctionSql)}
	f.WithDatabase(id.DatabaseName())
	f.WithFunctionDefinition(functionDefinition)
	f.WithName(id.Name())
	f.WithReturnType(returnType)
	f.WithSchema(id.SchemaName())
	return f
}

func (f *FunctionSqlModel) WithArgument(argName string, argDataType datatypes.DataType) *FunctionSqlModel {
	return f.WithArgumentsValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"arg_name":      tfconfig.StringVariable(argName),
				"arg_data_type": tfconfig.StringVariable(argDataType.ToSql()),
			},
		),
	)
}
