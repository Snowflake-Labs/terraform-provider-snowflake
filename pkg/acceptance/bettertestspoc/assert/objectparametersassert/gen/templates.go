package gen

import (
	"text/template"

	_ "embed"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/definition.tmpl
	snowflakeObjectAssertionsDefinitionTemplateContent string
	SnowflakeObjectAssertionsDefinitionTemplate, _     = template.New("snowflakeObjectAssertionsDefinitionTemplate").Parse(snowflakeObjectAssertionsDefinitionTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, SnowflakeObjectAssertionsDefinitionTemplate}
)
