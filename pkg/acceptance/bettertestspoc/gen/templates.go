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

	//go:embed templates/snowflake_object_assertions_definition.tmpl
	snowflakeObjectAssertionsDefinitionTemplateContent string
	SnowflakeObjectAssertionsDefinitionTemplate, _     = template.New("snowflakeObjectAssertionsDefinitionTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
	)).Parse(snowflakeObjectAssertionsDefinitionTemplateContent)

	//go:embed templates/snowflake_object_assertions.tmpl
	snowflakeObjectAssertionsTemplateContent string
	SnowflakeObjectAssertionsTemplate, _     = template.New("snowflakeObjectAssertionsTemplate").Funcs(gencommons.BuildTemplateFuncMap(
		gencommons.FirstLetterLowercase,
		gencommons.FirstLetter,
		gencommons.TypeWithoutPointer,
		gencommons.CamelToWords,
	)).Parse(snowflakeObjectAssertionsTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SnowflakeObjectAssertionsDefinitionTemplate, SnowflakeObjectAssertionsTemplate}
)
