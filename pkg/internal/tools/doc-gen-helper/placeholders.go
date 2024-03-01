package main

import "fmt"

func placeholder(name string) string {
	return fmt.Sprintf(`%%%%doc-gen-helper:%s%%%%`, name)
}

var (
	deprecatedResourcesPlaceholder   = placeholder("deprecated-resources")
	deprecatedDatasourcesPlaceholder = placeholder("deprecated-datasources")
)
