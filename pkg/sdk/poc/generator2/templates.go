package generator2

import "text/template"

var PackageTemplate, _ = template.New("packageTemplate").Parse(`
package {{ . }}
`)

var InterfaceTemplate, _ = template.New("interfaceTemplate").Parse(`
import "context"

type {{ .Name }} interface {
	{{- range .Operations }}
		{{ .Name }}(ctx context.Context, request *{{ .Options.DtoDecl }}) error
	{{- end }}
}
`)

var OptionsTemplate, _ = template.New("optionsTemplate").Parse(`
// {{ .Options.Name }} is based on {{ .Doc }}.
type {{ .Options.Name }} struct {
	{{- range .Options.Fields }}
			{{ .Name }} {{ .Kind.Kind }} {{ .TagsPrintable }}
	{{- end }}
}
`)

// TODO: merge with template above? (requires moving Doc to field)
var StructTemplate, _ = template.New("structTemplate").Parse(`
type {{ .Name }} struct {
	{{- range .Fields }}
			{{ .Name }} {{ .Kind.Kind }} {{ .TagsPrintable }}
	{{- end }}
}
`)

var DtoTemplate, _ = template.New("dtoTemplate").Parse(`
{{ define "DTO_STRUCT" }}
type {{ .DtoDecl }} struct {
	{{- range .Fields }}
		{{- if .ShouldBeInDto }}
		{{ .Name }} {{ .Kind.Kind }} {{ if .Required }}// required{{ end }}
		{{- end }}
	{{- end }}
}

{{- range .Fields }}
	{{ if .Struct }}
	{{ template "DTO_STRUCT" .Struct }}
	{{ end }}
{{- end }}
{{ end }}

//go:generate go run ../../dto-builder-generator/main.go

var (
	{{- range .Operations }}
	_ optionsProvider[{{ .Options.Name }}] = new({{ .Options.DtoDecl }})
	{{- end }}
)

{{- range .Operations }}
	{{ template "DTO_STRUCT" .Options }}
{{- end }}
`)

var ImplementationTemplate, _ = template.New("implementationTemplate").Funcs(template.FuncMap{
	"arr": arr,
}).Parse(`
{{ define "MAPPING" -}}
	{{- $struct := index . 0 }}
	{{- $path := index . 1 -}}

	&{{ $struct.Name }}{
		{{- range $struct.Fields }}
		{{- if .ShouldBeInDto }}
		{{ .Name }}: r.{{ .Name }},
		{{- end }}
		{{- end }}
	}
	{{- range $struct.Fields }}
		{{- if .ShouldBeInDto }}
			{{- if .Struct }}
			{{- $pathWithName := printf "%s%s" $path .Name }}
	if r.{{ $pathWithName }} != nil {
		opts.{{ $pathWithName }} = {{ template "MAPPING" (arr .Struct (printf "%s." .Name)) }}
	}
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}
import "context"

{{ $impl := .NameLowerCased }}
var _ {{ .Name }} = (*{{ $impl }})(nil)

type {{ $impl }} struct {
	client *Client
}

{{ range .Operations }}
func (v *{{ $impl }}) {{ .Name }}(ctx context.Context, request *{{ .Options.DtoDecl }}) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}
{{ end }}

{{ range .Operations }}
func (r *{{ .Options.DtoDecl }}) toOpts() *{{ .Options.Name }} {
	opts := {{ template "MAPPING" (arr .Options "") }}
	return opts
}
{{ end }}
`)

var TestFuncTemplate, _ = template.New("testFuncTemplate").Parse(`
{{ define "VALIDATION_TEST" }}
	{{ $field := . }}
	{{- range .Validations }}
	{{ .TodoComment $field }}
	{{- end }}
{{ end }}

{{ define "VALIDATIONS" }}
	{{ template "VALIDATION_TEST" . }}
	{{- range .Fields }}
		{{ if .HasAnyValidationInSubtree }}
			{{ template "VALIDATIONS" . }}
		{{ end }}
	{{- end }}
{{ end }}

import "testing"

{{- range .Operations }}
func Test{{ .ObjectInterface.Name }}_{{ .Name }}(t *testing.T) {
	id := random{{ .ObjectInterface.IdentifierKind }}(t)

	defaultOpts := func() *{{ .OptsField.KindNoPtr }} {
		return &{{ .OptsField.KindNoPtr }}{
			name: id,
		}
	}
	// TODO: remove me
	_ = defaultOpts()

	// TODO: fill me
	{{ template "VALIDATIONS" .OptsField }}
}
{{ end }}
`)

var ValidationsImplTemplate, _ = template.New("validationsImplTemplate").Funcs(template.FuncMap{
	"arr": arr,
}).Parse(`
{{ define "VALIDATIONS" }}
	{{- $struct := index . 0 }}
	{{- $path := index . 1 -}}

	{{- range $struct.Validations }}
	if {{ .Condition $struct }} {
		errs = append(errs, {{ .Error }})
	}
	{{ end -}}
	{{- range $struct.Fields }}
	{{- if .Struct }}
	{{- $pathWithName := printf "%s%s" $path .Name }}
	if valueSet(opts.{{ $pathWithName }}) {
		{{- template "VALIDATIONS" (arr .Struct (printf "%s." .Name)) -}}
	}
	{{ end -}}
	{{ end -}}
{{ end }}

import "errors"

var (
{{- range .Operations }}
	_ validatable = new({{ .Options.Name }})
{{- end }}
)
{{ range .Operations }}
func (opts *{{ .Options.Name }}) validate() error {
	if opts == nil {
		return errors.Join(errNilOptions)
	}
	var errs []error
	{{ template "VALIDATIONS" (arr .Options "") }}
	return errors.Join(errs...)
}
{{ end }}
`)

var IntegrationTestsTemplate, _ = template.New("integrationTestsTemplate").Parse(`
import "testing"

func TestInt_{{ .Name }}(t *testing.T) {
	// TODO: fill me
}
`)

func arr(args ...any) []any {
	return args
}
