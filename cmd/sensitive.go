package main

import (
	"encoding/csv"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"slices"
	"strconv"
)

type Schema struct {
	ResourceName string
	SchemaMap    map[string]*schema.Schema
}

type StringField struct {
	ResourceName string
	FieldName    string
	IsSensitive  bool
}

func main() {
	file, err := os.OpenFile("sensitive.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
		fields = append(fields, extractStringFieldsFromSchemaMap(s.ResourceName, s.SchemaMap)...)
	}

	return fields
}

func extractStringFieldsFromSchemaMap(resourceName string, schemaMap map[string]*schema.Schema) []StringField {
	fields := make([]StringField, 0)
	for fieldName, v := range schemaMap {
		switch v.Type {
		case schema.TypeString, schema.TypeMap:
			fields = append(fields, StringField{
				ResourceName: resourceName,
				FieldName:    fieldName,
				IsSensitive:  v.Sensitive,
			})
		case schema.TypeList, schema.TypeSet:
			switch elem := v.Elem.(type) {
			case *schema.Schema:
				fields = append(fields, StringField{
					ResourceName: resourceName,
					FieldName:    fieldName,
					IsSensitive:  v.Sensitive,
				})
			case *schema.Resource:
				fields = append(fields, extractStringFieldsFromSchemaMap(resourceName, elem.Schema)...)
			}
		}
	}
	return fields
}

func filterFields(fields []StringField) []StringField {
	filteredFields := make([]StringField, 0)
	fieldsNamesToFilter := []string{
		"comment",
		"created_on",
	}

	for _, field := range fields {
		if !slices.Contains(fieldsNamesToFilter, field.FieldName) {
			filteredFields = append(filteredFields, field)
		}
	}

	return filteredFields
}

func writeFields(writer *csv.Writer, fields []StringField) {
	if err := writer.Write([]string{"ResourceName", "FieldName", "IsSensitive"}); err != nil {
		log.Fatal(err)
	}

	for _, field := range fields {
		if err := writer.Write([]string{field.ResourceName, field.FieldName, strconv.FormatBool(field.IsSensitive)}); err != nil {
			log.Fatal(err)
		}
	}
}
