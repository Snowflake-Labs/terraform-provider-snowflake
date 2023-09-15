package generator

import "text/template"

var PackageTemplate, _ = template.New("packageTemplate").Parse(`
package {{ . }}
`)

var InterfaceTemplate, _ = template.New("interfaceTemplate").Parse(`
import "context"

type {{ .Name }} interface {
	{{- range .Operations }}
		{{- if and (eq .Name "Show") .ShowMapping }}
			{{ .Name }}(ctx context.Context, request *{{ .OptsField.DtoDecl }}) ([]{{ .ShowMapping.To.Name }}, error)
		{{- else if and (eq .Name "Describe") .DescribeMapping }}
			{{ .Name }}(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .DescribeMapping.To.Name }}, error)
		{{ else }}
			{{ .Name }}(ctx context.Context, request *{{ .OptsField.DtoDecl }}) error
		{{- end -}}
	{{ end }}
}
`)

var OptionsTemplate, _ = template.New("optionsTemplate").Parse(`
// {{ .OptsField.KindNoPtr }} is based on {{ .Doc }}.
type {{ .OptsField.KindNoPtr }} struct {
	{{- range .OptsField.Fields }}
			{{ .Name }} {{ .Kind }} {{ .TagsPrintable }}
	{{- end }}
}
`)

// TODO: merge with template above? (requires moving Doc to field)
var StructTemplate, _ = template.New("structTemplate").Parse(`
type {{ .KindNoPtr }} struct {
	{{- range .Fields }}
			{{ .Name }} {{ .Kind }} {{ .TagsPrintable }}
	{{- end }}
}
`)

var DtoTemplate, _ = template.New("dtoTemplate").Parse(`
{{ define "DTO_STRUCT" }}
	type {{ .DtoDecl }} struct {
		{{- range .Fields }}
			{{- if .ShouldBeInDto }}
			{{ .Name }} {{ .DtoKind }} {{ if .Required }}// required{{ end }}
			{{- end }}
		{{- end }}
	}

	{{- range .Fields }}
		{{ if .IsStruct }}
		{{ template "DTO_STRUCT" . }}
		{{ end }}
	{{- end }}
{{ end }}

//go:generate go run ./dto-builder-generator/main.go

var (
	{{- range .Operations }}
	_ optionsProvider[{{ .OptsField.KindNoPtr }}] = new({{ .OptsField.DtoDecl }})
	{{- end }}
)

{{- range .Operations }}
	{{ template "DTO_STRUCT" .OptsField }}
{{- end }}
`)

var ImplementationTemplate, _ = template.New("implementationTemplate").Parse(`
{{ define "MAPPING" -}}
	&{{ .KindNoPtr }}{
		{{- range .Fields }}
			{{- if .ShouldBeInDto }}
			{{ if .IsStruct }}{{ else }}{{ .Name }}: r{{ .Path }},{{ end -}}
			{{- end -}}
		{{- end }}
	}
	{{- range .Fields }}
		{{- if .ShouldBeInDto }}
			{{- if .IsStruct }}
				if r{{ .Path }} != nil {
					opts{{ .Path }} = {{ template "MAPPING" . -}}
				}
			{{- end -}}
		{{ end -}}
	{{ end }}
{{ end }}
{{ define "MAPPING_FUNC" }}
	func (r {{ .From.Name }}) {{ .MappingFuncName }}() *{{ .To.KindNoPtr }} {
		// TODO: Mapping
		return &{{ .To.KindNoPtr }}{}
	}
{{ end }}
import "context"

{{ $impl := .NameLowerCased }}
var _ {{ .Name }} = (*{{ $impl }})(nil)

type {{ $impl }} struct {
	client *Client
}
{{ range .Operations }}
	{{ if and (eq .Name "Show") .ShowMapping }}
		func (v *{{ $impl }}) Show(ctx context.Context, request *{{ .OptsField.DtoDecl }}) ([]{{ .ShowMapping.To.Name }}, error) {
			opts := request.toOpts()
			dbRows, err := validateAndQuery[{{ .ShowMapping.From.Name }}](v.client, ctx, opts)
			if err != nil {
				return nil, err
			}
			resultList := convertRows[{{ .ShowMapping.From.Name }}, {{ .ShowMapping.To.Name }}](dbRows)
			return resultList, nil
		}
	{{ else if and (eq .Name "Describe") .DescribeMapping }}
		func (v *{{ $impl }}) Describe(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .DescribeMapping.To.Name }}, error) {
			opts := &{{ .OptsField.Name }}{
				// TODO enforce this convention in the DSL (field "name" is queryStruct identifier)
				name: id,
			}
			result, err := validateAndQueryOne[{{ .DescribeMapping.From.Name }}](v.client, ctx, opts)
			if err != nil {
				return nil, err
			}
			return result.convert(), nil
		}
	{{ else }}
		func (v *{{ $impl }}) {{ .Name }}(ctx context.Context, request *{{ .OptsField.DtoDecl }}) error {
			opts := request.toOpts()
			return validateAndExec(v.client, ctx, opts)
		}
	{{ end }}
{{ end }}

{{ range .Operations }}
	func (r *{{ .OptsField.DtoDecl }}) toOpts() *{{ .OptsField.KindNoPtr }} {
		opts := {{ template "MAPPING" .OptsField -}}
		return opts
	}
	{{ if .ShowMapping }}
		{{ template "MAPPING_FUNC" .ShowMapping }}
	{{ end }}
	{{ if .DescribeMapping }}
		{{ template "MAPPING_FUNC" .DescribeMapping }}
	{{ end }}
{{ end }}
`)

var TestFuncTemplate, _ = template.New("testFuncTemplate").Parse(`
{{ define "VALIDATION_TEST" }}
	{{ $field := . }}
	{{- range .Validations }}
		t.Run("{{ .TodoComment $field }}", func(t *testing.T) {
			opts := defaultOpts()
			// TODO: fill me
			assertOptsInvalidJoinedErrors(t, opts, {{ .Error }})
		})
	{{ end -}}
{{ end }}

{{ define "VALIDATIONS" }}
	{{- template "VALIDATION_TEST" . -}}
	{{ range .Fields }}
		{{- if .HasAnyValidationInSubtree }}
			{{- template "VALIDATIONS" . -}}
		{{ end -}}
	{{- end -}}
{{ end }}

import "testing"

{{ range .Operations }}
	func Test{{ .ObjectInterface.Name }}_{{ .Name }}(t *testing.T) {
		id := random{{ .ObjectInterface.IdentifierKind }}(t)

		defaultOpts := func() *{{ .OptsField.KindNoPtr }} {
			return &{{ .OptsField.KindNoPtr }}{
				name: id,
			}
		}

		// TODO: remove me
		_ = defaultOpts()

		t.Run("validation: nil options", func(t *testing.T) {
			var opts *{{ .OptsField.KindNoPtr }} = nil
			assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
		})

		{{- template "VALIDATIONS" .OptsField -}}
	}
{{ end }}
`)

var ValidationsImplTemplate, _ = template.New("validationsImplTemplate").Parse(`
{{ define "VALIDATIONS" }}
	{{- $field := . -}}
	{{- range .Validations }}
		if {{ .Condition $field }} {
			errs = append(errs, {{ .Error $field }})
		}
	{{- end -}}
	{{- range .Fields }}
		{{- if .HasAnyValidationInSubtree }}
			if valueSet(opts{{ .Path }}) {
				{{- template "VALIDATIONS" . }}
			}
		{{- end -}}
	{{- end -}}
{{ end }}

import "errors"

var (
{{- range .Operations }}
	_ validatable = new({{ .OptsField.KindNoPtr }})
{{- end }}
)
{{ range .Operations }}
	func (opts *{{ .OptsField.KindNoPtr }}) validate() error {
		if opts == nil {
			return errors.Join(errNilOptions)
		}
		var errs []error
		{{- template "VALIDATIONS" .OptsField }}
		return errors.Join(errs...)
	}
{{ end }}
`)

var IntegrationTestsTemplate, _ = template.New("integrationTestsTemplate").Parse(`
import "testing"

func TestInt_{{ .Name }}(t *testing.T) {
	// TODO: prepare resources

	{{ range .Operations }}
	t.Run("{{ .Name }}", func(t *testing.T) {
		// TODO: fill me
	})
	{{ end -}}
}
`)
