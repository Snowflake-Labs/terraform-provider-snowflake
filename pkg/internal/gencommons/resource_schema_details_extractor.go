package gencommons

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type ResourceSchemaDetails struct {
	Name   string
	Fields []SchemaField
}

func (s ResourceSchemaDetails) ObjectName() string {
	return s.Name
}

type SchemaField struct {
	Name string
}

func ExtractResourceSchemaDetails(schema map[string]*schema.Schema) ResourceSchemaDetails {
	return ResourceSchemaDetails{}
}
