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

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
)

func main() {
	if len(os.Args) < 2 {
		log.Panic("Requires path as a first arg")
	}

	path := os.Args[1]
	additionalExamplesPath := filepath.Join(path, "examples", "additional")

	orderedResources := make([]string, 0)
	for key := range provider.Provider().ResourcesMap {
		orderedResources = append(orderedResources, key)
	}
	slices.Sort(orderedResources)

	deprecatedResources := make([]DeprecatedResource, 0)
	for _, key := range orderedResources {
		resource := provider.Provider().ResourcesMap[key]
		if resource.DeprecationMessage != "" {
			nameRelativeLink := docs.RelativeLink(key, filepath.Join("resources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := docs.GetDeprecatedResourceReplacement(resource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("resources", path))
			}

			deprecatedResources = append(deprecatedResources, DeprecatedResource{
				NameRelativeLink:        nameRelativeLink,
				ReplacementRelativeLink: replacementRelativeLink,
			})
		}
	}

	orderedDatasources := make([]string, 0)
	for key := range provider.Provider().DataSourcesMap {
		orderedDatasources = append(orderedDatasources, key)
	}
	slices.Sort(orderedDatasources)

	deprecatedDatasources := make([]DeprecatedDatasource, 0)
	for _, key := range orderedDatasources {
		datasource := provider.Provider().DataSourcesMap[key]
		if datasource.DeprecationMessage != "" {
			nameRelativeLink := docs.RelativeLink(key, filepath.Join("data-sources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := docs.GetDeprecatedResourceReplacement(datasource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("data-sources", path))
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

	err := os.WriteFile(filepath.Join(additionalExamplesPath, deprecatedResourcesFilename), deprecatedResourcesBuffer.Bytes(), 0o600)
	if err != nil {
		log.Panicln(err)
	}
	err = os.WriteFile(filepath.Join(additionalExamplesPath, deprecatedDatasourcesFilename), deprecatedDatasourcesBuffer.Bytes(), 0o600)
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
