package generator

import "text/template"

var InterfaceTemplate, _ = template.New("interfaceTemplate").Parse(`
type {{.Name}} interface {
	{{- range .Operations}}
		{{.Name}}(ctx context.Context, opts *{{.OptsName}}) error
	{{- end}}
}
`)

var OptionsTemplate, _ = template.New("optionsTemplate").Parse(`
// {{.OptsName}} is based on {{.Doc}}.
type {{.OptsName}} struct {
	{{- range .OptsStructFields}}
			{{.Name}} {{.Kind}} {{.TagsPrintable}}
	{{- end}}
}
`)
