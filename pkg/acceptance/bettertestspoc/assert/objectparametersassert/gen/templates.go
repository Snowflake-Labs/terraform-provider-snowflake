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

	//go:embed templates/aggregated_generic_checks.tmpl
	aggregatedGenericChecksTemplateContent string
	AggregatedGenericChecksTemplate, _     = template.New("aggregatedGenericChecksTemplateContent").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
		gencommons.IsLastItem,
	)).Parse(aggregatedGenericChecksTemplateContent)

	//go:embed templates/value_checks.tmpl
	valueChecksTemplateContent string
	ValueChecksTemplate, _     = template.New("valueChecksTemplateContent").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.SnakeCaseToCamel,
	)).Parse(valueChecksTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SnowflakeObjectAssertionsDefinitionTemplate, GenericChecksTemplate, AggregatedGenericChecksTemplate, ValueChecksTemplate}
)
