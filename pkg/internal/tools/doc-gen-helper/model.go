package main

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
