package main

import (
	"encoding/csv"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/maps"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Schema struct {
	ResourceName string
	SchemaMap    map[string]*schema.Schema
}

type StringField struct {
	ResourceName string
	FieldName    string
	IsSensitive  bool
	IsComputed   bool
}

func NewStringField(resourceName, fieldName string, isSensitive, isComputed bool) StringField {
	return StringField{
		ResourceName: resourceName,
		FieldName:    fieldName,
		IsSensitive:  isSensitive,
		IsComputed:   isComputed,
	}
}

var fieldsNamesToFilter = []string{
	// values used in many resources that are not sensitive
	"database",
	"schema",
	"name",
	"comment",
	"created_on",
	"fully_qualified_name",

	// only used in data sources (in filtering)
	"like",
	"starts_with",
	"from",
}

func main() {
	file, err := os.OpenFile("cmd/sensitive.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	schemas := buildSchemas()
	log.Printf("Has %d schemas", len(schemas))

	fields := extractStringFields(schemas)
	log.Printf("Has %d fields", len(fields))

	filteredFields := filterFields(fields)
	log.Printf("Has %d fields after filtering", len(filteredFields))

	writeFields(writer, filteredFields)

	writer.Flush()
}

func buildSchemas() []Schema {
	p := provider.Provider()

	schemas := make([]Schema, 0)
	schemas = append(schemas, Schema{
		ResourceName: "provider",
		SchemaMap:    p.Schema,
	})

	for k, v := range p.ResourcesMap {
		schemas = append(schemas, Schema{
			ResourceName: k,
			SchemaMap:    v.Schema,
		})
	}

	for k, v := range p.DataSourcesMap {
		schemas = append(schemas, Schema{
			ResourceName: k,
			SchemaMap:    v.Schema,
		})
	}

	return schemas
}

func extractStringFields(schemas []Schema) []StringField {
	fields := make([]StringField, 0)

	for _, s := range schemas {
		fields = append(fields, extractStringFieldsFromSchemaMap(s.ResourceName, "", s.SchemaMap)...)
	}

	return fields
}

func extractStringFieldsFromSchemaMap(resourceName string, parentName string, schemaMap map[string]*schema.Schema) []StringField {
	fields := make([]StringField, 0)
	for fieldName, v := range schemaMap {
		switch v.Type {
		case schema.TypeString, schema.TypeMap:
			fields = append(fields, NewStringField(resourceName, parentName+fieldName, v.Sensitive, v.Computed))
		case schema.TypeList, schema.TypeSet:
			switch elem := v.Elem.(type) {
			case *schema.Schema:
				fields = append(fields, NewStringField(resourceName, parentName+fieldName, v.Sensitive, v.Computed))
			case *schema.Resource:
				if slices.ContainsFunc(maps.Keys(elem.Schema), func(name string) bool { return slices.Contains([]string{"key", "value", "default"}, name) }) {
					fields = append(fields, NewStringField(resourceName, parentName+fieldName, v.Sensitive, v.Computed))
				} else {
					// check recursively
					var parent string
					if parentName != "" {
						parent = parentName + "." + fieldName + "."
					} else {
						parent = fieldName + "."
					}
					fields = append(fields, extractStringFieldsFromSchemaMap(resourceName, parent, elem.Schema)...)
				}
			}
		}
	}
	return fields
}

func filterFields(fields []StringField) []StringField {
	filteredFields := make([]StringField, 0)

	for _, field := range fields {
		fieldNameParts := strings.Split(field.FieldName, ".")
		lastFieldNamePart := fieldNameParts[len(fieldNameParts)-1]
		if !slices.Contains(fieldsNamesToFilter, lastFieldNamePart) {
			filteredFields = append(filteredFields, field)
		}
	}

	return filteredFields
}

func writeFields(writer *csv.Writer, fields []StringField) {
	if err := writer.Write([]string{"ResourceName", "FieldName", "IsSensitive", "IsComputed"}); err != nil {
		log.Fatal(err)
	}

	for _, field := range fields {
		if err := writer.Write([]string{field.ResourceName, field.FieldName, strconv.FormatBool(field.IsSensitive), strconv.FormatBool(field.IsComputed)}); err != nil {
			log.Fatal(err)
		}
	}
}
