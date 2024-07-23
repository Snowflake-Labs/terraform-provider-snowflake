package gen

import (
	"text/template"

	_ "embed"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate}
)
