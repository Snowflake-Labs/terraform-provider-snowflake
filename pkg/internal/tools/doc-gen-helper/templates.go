package main

import "text/template"

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`## Currently deprecated resources

{{ range .Resources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`,
)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`## Currently deprecated datasources

{{ range .Datasources -}}
	- {{ .NameRelativeLink }}{{ if .ReplacementRelativeLink }} - use {{ .ReplacementRelativeLink }} instead{{ end }}
{{ end -}}`,
)
