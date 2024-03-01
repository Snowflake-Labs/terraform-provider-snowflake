package main

import "text/template"

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`{{ range .Resources -}}
	- {{ .Name }}
{{ end -}}`)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`{{ range .Datasources -}}
	- {{ .Name }}
{{ end -}}`)
