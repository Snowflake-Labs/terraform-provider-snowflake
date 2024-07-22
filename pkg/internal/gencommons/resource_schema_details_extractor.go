package gencommons

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceSchemaDetails struct {
	Name       string
	Attributes []SchemaAttribute
}

func (s ResourceSchemaDetails) ObjectName() string {
	return s.Name
}

type SchemaAttribute struct {
	Name          string
	AttributeType schema.ValueType
}

func ExtractResourceSchemaDetails(name string, schema map[string]*schema.Schema) ResourceSchemaDetails {
	orderedAttributeNames := make([]string, 0)
	for key := range schema {
		orderedAttributeNames = append(orderedAttributeNames, key)
	}
	slices.Sort(orderedAttributeNames)

	attributes := make([]SchemaAttribute, 0)
	for _, k := range orderedAttributeNames {
		s := schema[k]
		attributes = append(attributes, SchemaAttribute{
			Name:          k,
			AttributeType: s.Type,
		})
	}

	return ResourceSchemaDetails{
		Name:       name,
		Attributes: attributes,
	}
}
