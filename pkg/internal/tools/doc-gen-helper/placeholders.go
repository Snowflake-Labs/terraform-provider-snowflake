package main

import "fmt"

func placeholder(name string) string {
	return fmt.Sprintf(`%%%%doc-gen-helper:%s%%%%`, name)
}

var deprecatedResourcesPlaceholder = placeholder("deprecated-resources")
var deprecatedDatasourcesPlaceholder = placeholder("deprecated-datasources")
