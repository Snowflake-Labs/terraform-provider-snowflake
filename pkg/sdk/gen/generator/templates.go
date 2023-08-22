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

var StructTemplate, _ = template.New("structTemplate").Parse(`
type {{.KindNoPtr}} struct {
	{{- range .Fields}}
			{{.Name}} {{.Kind}} {{.TagsPrintable}}
	{{- end}}
}
`)

var ImplementationTemplate, _ = template.New("implementationTemplate").Parse(`
{{$impl := .NameLowerCased}}
var _ {{.Name}} = (*{{$impl}})(nil)

type {{$impl}} struct {
	client *Client
}
{{range .Operations}}
func (v *{{$impl}}) {{.Name}}(ctx context.Context, opts *{{.OptsName}}) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}
{{end}}
`)

var TestFuncTemplate, _ = template.New("testFuncTemplate").Parse(`
func Test{{.ObjectInterface.Name}}_{{.Name}}(t *testing.T) {
	id := random{{.ObjectInterface.IdentifierKind}}(t)

	defaultOpts := func() *{{.OptsName}} {
		return &{{.OptsName}}{
			name: id,
		}
	}

	// TODO: fill me
}
`)
