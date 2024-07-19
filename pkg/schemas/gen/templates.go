package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate, _     = template.New("schemaTemplate").Parse(schemaTemplateContent)

	//go:embed templates/to_schema_mapper.tmpl
	toSchemaMapperTemplateContent string
	ToSchemaMapperTemplate, _     = template.New("toSchemaMapperTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.RunMapper,
	)).Parse(toSchemaMapperTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SchemaTemplate, ToSchemaMapperTemplate}
)
