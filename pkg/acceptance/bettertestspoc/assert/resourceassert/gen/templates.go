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
	DefinitionTemplate, _     = template.New("definitionTemplate").Parse(definitionTemplateContent)

	//go:embed templates/assertions.tmpl
	assertionsTemplateContent string
	AssertionsTemplate, _     = template.New("assertionsTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
	)).Parse(assertionsTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, AssertionsTemplate}
)
