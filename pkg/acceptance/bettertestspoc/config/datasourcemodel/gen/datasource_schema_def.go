package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatasourceSchemaDef struct {
	name   string
	schema map[string]*schema.Schema
}

// TODO [SNOW-1501905]: rename ResourceSchemaDetails (because it is used for the datasources and provider too)
func GetDatasourceSchemaDetails() []genhelpers.ResourceSchemaDetails {
	allDatasourcesSchemas := allDatasourcesSchemaDefs
	allDatasourcesSchemasDetails := make([]genhelpers.ResourceSchemaDetails, len(allDatasourcesSchemas))
	for idx, s := range allDatasourcesSchemas {
		allDatasourcesSchemasDetails[idx] = genhelpers.ExtractResourceSchemaDetails(s.name, s.schema)
	}
	return allDatasourcesSchemasDetails
}

var allDatasourcesSchemaDefs = []DatasourceSchemaDef{
	{
		name:   "Accounts",
		schema: datasources.Accounts().Schema,
	},
	{
		name:   "Database",
		schema: datasources.Database().Schema,
	},
	{
		name:   "DatabaseRole",
		schema: datasources.DatabaseRole().Schema,
	},
	{
		name:   "DatabaseRoles",
		schema: datasources.DatabaseRoles().Schema,
	},
	{
		name:   "Databases",
		schema: datasources.Databases().Schema,
	},
	{
		name:   "Grants",
		schema: datasources.Grants().Schema,
	},
}
