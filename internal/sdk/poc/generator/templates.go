// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package generator

import "text/template"

var PackageTemplate, _ = template.New("packageTemplate").Parse(`
package {{ . }}
`)

var InterfaceTemplate, _ = template.New("interfaceTemplate").
	Funcs(template.FuncMap{
		"deref": func(p *DescriptionMappingKind) string { return string(*p) },
	}).
	Parse(`
import "context"

type {{ .Name }} interface {
	{{- range .Operations }}
		{{- if and (eq .Name "Show") .ShowMapping }}
			{{ .Name }}(ctx context.Context, request *{{ .OptsField.DtoDecl }}) ([]{{ .ShowMapping.To.Name }}, error)
		{{- else if eq .Name "ShowByID" }}
			{{ .Name }}(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .ObjectInterface.NameSingular }}, error)
		{{- else if and (eq .Name "Describe") .DescribeMapping }}
			{{- if .DescribeKind }}
				{{- if eq (deref .DescribeKind) "single_value" }}
					{{ .Name }}(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .DescribeMapping.To.Name }}, error)
				{{- else if eq (deref .DescribeKind) "slice" }}
					{{ .Name }}(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) ([]{{ .DescribeMapping.To.Name }}, error)
				{{- end }}
			{{- end }}
		{{- else }}
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
//go:generate go run ./dto-builder-generator/main.go

var (
	{{- range .Operations }}
	{{- if .OptsField }}
	_ optionsProvider[{{ .OptsField.KindNoPtr }}] = new({{ .OptsField.DtoDecl }})
	{{- end }}
	{{- end }}
)
`)

var DtoDeclTemplate, _ = template.New("dtoTemplate").Parse(`
type {{ .DtoDecl }} struct {
	{{- range .Fields }}
		{{- if .ShouldBeInDto }}
		{{ .Name }} {{ .DtoKind }} {{ if .Required }}// required{{ end }}
		{{- end }}
	{{- end }}
}
`)

var ImplementationTemplate, _ = template.New("implementationTemplate").
	Funcs(template.FuncMap{
		"deref": func(p *DescriptionMappingKind) string { return string(*p) },
	}).
	Parse(`
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
					{{- if not .IsSlice }}
						opts{{ .Path }} = {{ template "MAPPING" . -}}
					{{- else }}
						s := make({{ .Kind }}, len(r{{ .Path }}))
						for i, v := range r{{ .Path }} {
							s[i] = {{ .KindNoSlice }}{
							  {{- range .Fields }}
								   {{ .Name }}: v.{{ .Name }},
							  {{- end }}
							}
						}
						opts{{ .Path }} = s
					{{ end -}}
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
import (
"context"

"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/internal/collections"
)

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
	{{ else if eq .Name "ShowByID" }}
		func (v *{{ $impl }}) ShowByID(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .ObjectInterface.NameSingular }}, error) {
			// TODO: adjust request if e.g. LIKE is supported for the resource
			{{ $impl }}, err := v.Show(ctx, NewShow{{ .ObjectInterface.NameSingular }}Request())
			if err != nil {
				return nil, err
			}
			return collections.FindOne({{ $impl }}, func(r {{ .ObjectInterface.NameSingular }}) bool { return r.Name == id.Name() })
		}
	{{ else if and (eq .Name "Describe") .DescribeMapping }}
		{{ if .DescribeKind }}
			{{ if eq (deref .DescribeKind) "single_value" }}
				func (v *{{ $impl }}) Describe(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) (*{{ .DescribeMapping.To.Name }}, error) {
					opts := &{{ .OptsField.Name }}{
						 name: id,
					}
					result, err := validateAndQueryOne[{{ .DescribeMapping.From.Name }}](v.client, ctx, opts)
					if err != nil {
						 return nil, err
					}
					return result.convert(), nil
				}
			{{ else if eq (deref .DescribeKind) "slice" }}
				func (v *{{ $impl }}) Describe(ctx context.Context, id {{ .ObjectInterface.IdentifierKind }}) ([]{{ .DescribeMapping.To.Name }}, error) {
					opts := &{{ .OptsField.Name }}{
						 name: id,
					}
					rows, err := validateAndQuery[{{ .DescribeMapping.From.Name}}](v.client, ctx, opts)
					if err != nil {
						 return nil, err
					}
					return convertRows[{{ .DescribeMapping.From.Name }}, {{ .DescribeMapping.To.Name }}](rows), nil
				}
			{{ end }}
		{{ end }}
	{{ else }}
		func (v *{{ $impl }}) {{ .Name }}(ctx context.Context, request *{{ .OptsField.DtoDecl }}) error {
			opts := request.toOpts()
			return validateAndExec(v.client, ctx, opts)
		}
	{{ end }}
{{ end }}

{{ range .Operations }}
	{{- if .OptsField }}
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
	{{- end}}
{{ end }}
`)

var TestFuncTemplate, _ = template.New("testFuncTemplate").Parse(`
{{ define "VALIDATION_TEST" }}
	{{ $field := . }}
	{{- range .Validations }}
		t.Run("{{ .TodoComment $field }}", func(t *testing.T) {
			opts := defaultOpts()
			// TODO: fill me
			assertOptsInvalidJoinedErrors(t, opts, {{ .ReturnedError $field }})
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
	{{- if .OptsField }}
	func Test{{ .ObjectInterface.Name }}_{{ .Name }}(t *testing.T) {
		id := Random{{ .ObjectInterface.IdentifierKind }}()

		// Minimal valid {{ .OptsField.KindNoPtr }}
		defaultOpts := func() *{{ .OptsField.KindNoPtr }} {
			return &{{ .OptsField.KindNoPtr }}{
				name: id,
			}
		}

		t.Run("validation: nil options", func(t *testing.T) {
			var opts *{{ .OptsField.KindNoPtr }} = nil
			assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
		})

		{{- template "VALIDATIONS" .OptsField }}

		t.Run("basic", func(t *testing.T) {
			opts := defaultOpts()
			// TODO: fill me
			assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
		})

		t.Run("all options", func(t *testing.T) {
			opts := defaultOpts()
			// TODO: fill me
			assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
		})
	}
	{{- end }}
{{ end }}
`)

var ValidationsImplTemplate, _ = template.New("validationsImplTemplate").Parse(`
{{ define "VALIDATIONS" }}
	{{- $field := . -}}
	{{- range .Validations }}
		if {{ .Condition $field }} {
			errs = append(errs, {{ .ReturnedError $field }})
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
	{{- if .OptsField }}
	_ validatable = new({{ .OptsField.KindNoPtr }})
	{{- end }}
{{- end }}
)
{{ range .Operations }}
	{{- if .OptsField }}
	func (opts *{{ .OptsField.KindNoPtr }}) validate() error {
		if opts == nil {
			return errors.Join(ErrNilOptions)
		}
		var errs []error
		{{- template "VALIDATIONS" .OptsField }}
		return errors.Join(errs...)
	}
	{{- end }}
{{ end }}
`)

var IntegrationTestsTemplate, _ = template.New("integrationTestsTemplate").Parse(`
import "testing"

func TestInt_{{ .Name }}(t *testing.T) {
	// TODO: prepare common resources

	{{ range .Operations }}
	t.Run("{{ .Name }}", func(t *testing.T) {
		// TODO: fill me
	})
	{{ end -}}
}
`)
