package model

import (
	"encoding/json"

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

func FunctionJavaWithId(
	resourceName string,
	id sdk.SchemaObjectIdentifierWithArguments,
	returnType datatypes.DataType,
	handler string,
	functionDefinition string,
) *FunctionJavaModel {
	return FunctionJava(resourceName, id.DatabaseName(), functionDefinition, handler, id.Name(), returnType.ToSql(), id.SchemaName())
}
