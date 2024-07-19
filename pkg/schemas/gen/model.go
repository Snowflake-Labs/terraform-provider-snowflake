package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

type ShowResultSchemaModel struct {
	Name         string
	SdkType      string
	SchemaFields []SchemaField
}

func (m ShowResultSchemaModel) SomeFunc() {
	return
}

func ModelFromStructDetails(sdkStruct gencommons.StructDetails) ShowResultSchemaModel {
	name, _ := strings.CutPrefix(sdkStruct.Name, "sdk.")
	schemaFields := make([]SchemaField, len(sdkStruct.Fields))
	for idx, field := range sdkStruct.Fields {
		schemaFields[idx] = MapToSchemaField(field)
	}

	return ShowResultSchemaModel{
		Name:         name,
		SdkType:      sdkStruct.Name,
		SchemaFields: schemaFields,
	}
}
