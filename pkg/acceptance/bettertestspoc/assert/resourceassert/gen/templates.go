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
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Parse(definitionTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate}
)
