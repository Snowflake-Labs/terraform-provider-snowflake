package gen

import (
	"strings"
	"text/template"

	_ "embed"
)

// TODO [SNOW-1501905]: extract common funcs
var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate, _     = template.New("schemaTemplate").Funcs(template.FuncMap{
		"firstLetterLowercase": func(in string) string { return strings.ToLower(in[:1]) + in[1:] },
	}).Parse(schemaTemplateContent)

	//go:embed templates/to_schema_mapper.tmpl
	toSchemaMapperTemplateContent string
	ToSchemaMapperTemplate, _     = template.New("toSchemaMapperTemplate").Funcs(template.FuncMap{
		"firstLetterLowercase": func(in string) string { return strings.ToLower(in[:1]) + in[1:] },
		"runMapper":            func(mapper Mapper, in ...string) string { return mapper(strings.Join(in, "")) },
	}).Parse(toSchemaMapperTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SchemaTemplate, ToSchemaMapperTemplate}
)
