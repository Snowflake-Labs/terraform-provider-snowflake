//aaago:build exclude

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
)

func main() {
	if len(os.Args) < 2 {
		log.Panic("Requires path as a first arg")
	}

	path := os.Args[1]
	templatesPath := filepath.Join(path, "templates")
	currentPath := filepath.Join(path, "pkg", "internal", "tools", "doc-gen-helper")

	deprecatedResourcesPlaceholder := "%%doc-gen-helper:deprecated-resources%%"
	deprecatedDatasourcesPlaceholder := "%%doc-gen-helper:deprecated-datasources%%"

	indexTemplateRaw, err := os.ReadFile(filepath.Join(currentPath, "index.md.tmpl"))
	if err != nil {
		log.Panicf("Could not open index template file, %v", err)
	}
	indexTemplateContents := string(indexTemplateRaw)

	fmt.Printf("Will save template in %v\n", templatesPath)

	deprecatedResources := make([]DeprecatedResource, 0)
	for key, resource := range provider.Provider().ResourcesMap {
		if resource.DeprecationMessage != "" {
			deprecatedResources = append(deprecatedResources, DeprecatedResource{Name: key})
		}
	}

	deprecatedDatasources := make([]DeprecatedDatasource, 0)
	for key, datasource := range provider.Provider().DataSourcesMap {
		if datasource.DeprecationMessage != "" {
			deprecatedDatasources = append(deprecatedDatasources, DeprecatedDatasource{Name: key})
		}
	}

	var deprecatedResourcesBuffer bytes.Buffer
	printTo(&deprecatedResourcesBuffer, DeprecatedResourcesTemplate, DeprecatedResourcesContext{deprecatedResources})

	var deprecatedDatasourcesBuffer bytes.Buffer
	printTo(&deprecatedDatasourcesBuffer, DeprecatedDatasourcesTemplate, DeprecatedDatasourcesContext{deprecatedDatasources})

	indexTemplateContents = strings.ReplaceAll(indexTemplateContents, deprecatedResourcesPlaceholder, deprecatedResourcesBuffer.String())
	indexTemplateContents = strings.ReplaceAll(indexTemplateContents, deprecatedDatasourcesPlaceholder, deprecatedDatasourcesBuffer.String())

	err = os.WriteFile(filepath.Join(templatesPath, "updated-index.md.tmpl"), []byte(indexTemplateContents), 0o600)
	if err != nil {
		log.Panicln(err)
	}
}

func printTo(writer io.Writer, template *template.Template, model any) {
	err := template.Execute(writer, model)
	if err != nil {
		log.Panicln(err)
	}
}

type DeprecatedResourcesContext struct {
	Resources []DeprecatedResource
}

type DeprecatedResource struct {
	Name string
}

type DeprecatedDatasourcesContext struct {
	Datasources []DeprecatedDatasource
}

type DeprecatedDatasource struct {
	Name string
}

var DeprecatedResourcesTemplate, _ = template.New("deprecatedResourcesTemplate").Parse(
	`{{ range .Resources -}}
	- {{ .Name }}
{{ end -}}`)

var DeprecatedDatasourcesTemplate, _ = template.New("deprecatedDatasourcesTemplate").Parse(
	`{{ range .Datasources -}}
	- {{ .Name }}
{{ end -}}`)
