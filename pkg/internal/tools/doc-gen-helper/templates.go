package main

import "text/template"

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`<!-- Section of deprecated resources -->
{{if gt (len .Resources) 0}} ## Currently deprecated resources {{end}}

{{ range .Resources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`,
)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`<!-- Section of deprecated data sources -->
{{if gt (len .Datasources) 0}} ## Currently deprecated data sources {{end}}

{{ range .Datasources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`,
)
