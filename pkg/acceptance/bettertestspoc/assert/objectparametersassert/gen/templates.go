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
	GenericChecksTemplate, _     = template.New("genericChecksTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
	)).Parse(genericChecksTemplateContent)

	//go:embed templates/aggregated_generic_checks.tmpl
	aggregatedGenericChecksTemplateContent string
	AggregatedGenericChecksTemplate, _     = template.New("aggregatedGenericChecksTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
		gencommons.IsLastItem,
	)).Parse(aggregatedGenericChecksTemplateContent)

	//go:embed templates/specific_checks.tmpl
	specificChecksTemplateContent string
	SpecificChecksTemplate, _     = template.New("specificChecksTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
	)).Parse(specificChecksTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SnowflakeObjectAssertionsDefinitionTemplate, GenericChecksTemplate, AggregatedGenericChecksTemplate, SpecificChecksTemplate}
)
