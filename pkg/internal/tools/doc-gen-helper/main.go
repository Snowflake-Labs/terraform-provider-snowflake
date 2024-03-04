//aaago:build exclude

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
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

	indexTemplateRaw, err := os.ReadFile(filepath.Join(currentPath, "index.md.tmpl"))
	if err != nil {
		log.Panicf("Could not open index template file, %v", err)
	}
	indexTemplateContents := string(indexTemplateRaw)

	orderedResources := make([]string, 0)
	for key, _ := range provider.Provider().ResourcesMap {
		orderedResources = append(orderedResources, key)
	}
	slices.Sort(orderedResources)

	deprecatedResources := make([]DeprecatedResource, 0)
	for _, key := range orderedResources {
		resource := provider.Provider().ResourcesMap[key]
		if resource.DeprecationMessage != "" {
			nameRelativeLink := provider.RelativeLink(key, filepath.Join("resources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := provider.GetDeprecatedResourceReplacement(resource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = provider.RelativeLink(replacement, filepath.Join("resources", path))
			}

			deprecatedResources = append(deprecatedResources, DeprecatedResource{
				NameRelativeLink:        nameRelativeLink,
				ReplacementRelativeLink: replacementRelativeLink,
			})
		}
	}

	orderedDatasources := make([]string, 0)
	for key, _ := range provider.Provider().DataSourcesMap {
		orderedDatasources = append(orderedDatasources, key)
	}
	slices.Sort(orderedDatasources)

	deprecatedDatasources := make([]DeprecatedDatasource, 0)
	for _, key := range orderedDatasources {
		datasource := provider.Provider().DataSourcesMap[key]
		if datasource.DeprecationMessage != "" {
			nameRelativeLink := provider.RelativeLink(key, filepath.Join("resources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := provider.GetDeprecatedResourceReplacement(datasource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = provider.RelativeLink(replacement, filepath.Join("resources", path))
			}

			deprecatedDatasources = append(deprecatedDatasources, DeprecatedDatasource{
				NameRelativeLink:        nameRelativeLink,
				ReplacementRelativeLink: replacementRelativeLink,
			})
		}
	}

	var deprecatedResourcesBuffer bytes.Buffer
	printTo(&deprecatedResourcesBuffer, DeprecatedResourcesTemplate, DeprecatedResourcesContext{deprecatedResources})

	var deprecatedDatasourcesBuffer bytes.Buffer
	printTo(&deprecatedDatasourcesBuffer, DeprecatedDatasourcesTemplate, DeprecatedDatasourcesContext{deprecatedDatasources})

	indexTemplateContents = strings.ReplaceAll(indexTemplateContents, deprecatedResourcesPlaceholder, deprecatedResourcesBuffer.String())
	indexTemplateContents = strings.ReplaceAll(indexTemplateContents, deprecatedDatasourcesPlaceholder, deprecatedDatasourcesBuffer.String())

	err = os.WriteFile(filepath.Join(templatesPath, "index.md.tmpl"), []byte(indexTemplateContents), 0o600)
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
