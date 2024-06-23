package gen

import (
	"strings"
	"text/template"

	_ "embed"
)

// TODO: extract common funcs
var (
	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate, _     = template.New("schemaTemplate").Funcs(template.FuncMap{
		"uppercase":            func(in string) string { return strings.ToUpper(in) },
		"lowercase":            func(in string) string { return strings.ToLower(in) },
		"firstLetterLowercase": func(in string) string { return strings.ToLower(in[:1]) + in[1:] },
	}).Parse(schemaTemplateContent)

	//go:embed templates/to_schema_mapper.tmpl
	toSchemaMapperTemplateContent string
	ToSchemaMapperTemplate, _     = template.New("toSchemaMapperTemplate").Funcs(template.FuncMap{
		"uppercase":            func(in string) string { return strings.ToUpper(in) },
		"lowercase":            func(in string) string { return strings.ToLower(in) },
		"firstLetterLowercase": func(in string) string { return strings.ToLower(in[:1]) + in[1:] },
		"runMapper":            func(mapper Mapper, in ...string) string { return mapper(strings.Join(in, "")) },
	}).Parse(toSchemaMapperTemplateContent)
)
