package main

import "fmt"

func markdown(name string) string {
	return fmt.Sprintf(`%s.MD`, name)
}

var (
	deprecatedResourcesFilename   = markdown("deprecated_resources")
	deprecatedDatasourcesFilename = markdown("deprecated_datasources")
)
