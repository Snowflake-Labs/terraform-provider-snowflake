package gen

import "strings"

type ShowResultSchemaModel struct {
	Name         string
	SdkType      string
	SchemaFields []SchemaField
}

func ModelFromStructDetails(sdkStruct Struct) ShowResultSchemaModel {
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
