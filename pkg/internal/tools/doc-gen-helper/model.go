package main

type DeprecatedResourcesContext struct {
	Resources []DeprecatedResource
}

type DeprecatedResource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}

type DeprecatedDatasourcesContext struct {
	Datasources []DeprecatedDatasource
}

type DeprecatedDatasource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}
