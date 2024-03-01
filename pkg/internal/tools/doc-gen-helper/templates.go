package main

import "text/template"

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`{{ range .Resources -}}
	- {{ .Name }}{{ if .Replacement}} - use [{{ .Replacement }}]({{ .ReplacementPathRelative }}) instead{{ end }}
{{ end -}}`)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`{{ range .Datasources -}}
	- {{ .Name }}{{ if .Replacement}} - use [{{ .Replacement }}]({{ .ReplacementPathRelative }}) instead{{ end }}
{{ end -}}`)
