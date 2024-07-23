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

	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
	)).Parse(definitionTemplateContent)

	//go:embed templates/specific_checks.tmpl
	specificChecksTemplateContent string
	SpecificChecksTemplate, _     = template.New("specificChecksTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
	)).Parse(specificChecksTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, SpecificChecksTemplate}
)
