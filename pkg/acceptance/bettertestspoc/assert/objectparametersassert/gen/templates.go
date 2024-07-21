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
	snowflakeObjectAssertionsDefinitionTemplateContent string
	SnowflakeObjectAssertionsDefinitionTemplate, _     = template.New("snowflakeObjectAssertionsDefinitionTemplate").Parse(snowflakeObjectAssertionsDefinitionTemplateContent)

	//go:embed templates/generic_checks.tmpl
	genericChecksTemplateContent string
	GenericChecksTemplate, _     = template.New("genericChecksTemplateContent").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
	)).Parse(genericChecksTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SnowflakeObjectAssertionsDefinitionTemplate, GenericChecksTemplate}
)
