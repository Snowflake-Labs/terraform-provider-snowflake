package main

import (
	"bytes"
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
			nameRelativeLink := docs.RelativeLink(key, filepath.Join("docs", "resources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := docs.GetDeprecatedResourceReplacement(resource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("docs", "resources", path))
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
			nameRelativeLink := docs.RelativeLink(key, filepath.Join("docs", "data-sources", strings.Replace(key, "snowflake_", "", 1)))

			replacement, path, _ := docs.GetDeprecatedResourceReplacement(datasource.DeprecationMessage)
			var replacementRelativeLink string
			if replacement != "" && path != "" {
				replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("docs", "data-sources", path))
			}

			deprecatedDatasources = append(deprecatedDatasources, DeprecatedDatasource{
				NameRelativeLink:        nameRelativeLink,
				ReplacementRelativeLink: replacementRelativeLink,
			})
		}
	}

	err := printTo(DeprecatedResourcesTemplate, DeprecatedResourcesContext{deprecatedResources}, filepath.Join(additionalExamplesPath, deprecatedResourcesFilename))
	if err != nil {
		log.Fatal(err)
	}

	err = printTo(DeprecatedDatasourcesTemplate, DeprecatedDatasourcesContext{deprecatedDatasources}, filepath.Join(additionalExamplesPath, deprecatedDatasourcesFilename))
	if err != nil {
		log.Fatal(err)
	}
}

func printTo(template *template.Template, model any, filepath string) error {
	var writer bytes.Buffer
	err := template.Execute(&writer, model)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, writer.Bytes(), 0o600)
}
