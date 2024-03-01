package main

import "text/template"

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`{{ range .Resources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`{{ range .Datasources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`)
