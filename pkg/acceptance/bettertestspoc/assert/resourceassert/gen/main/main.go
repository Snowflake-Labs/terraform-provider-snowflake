//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func main() {
	gencommons.NewGenerator(
		getResourceSchemaDetails,
		gen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

type ResourceSchemaDef struct {
	name   string
	schema map[string]*schema.Schema
}

func getResourceSchemaDetails() []gencommons.ResourceSchemaDetails {
	allResourceSchemas := allResourceSchemaDefs
	allResourceSchemasDetails := make([]gencommons.ResourceSchemaDetails, len(allResourceSchemas))
	for idx, s := range allResourceSchemas {
		allResourceSchemasDetails[idx] = gencommons.ExtractResourceSchemaDetails(s.name, s.schema)
	}
	return allResourceSchemasDetails
}

func getFilename(_ gencommons.ResourceSchemaDetails, model gen.ResourceAssertionsModel) string {
	return gencommons.ToSnakeCase(model.Name) + "_resource" + "_gen.go"
}

var allResourceSchemaDefs = []ResourceSchemaDef{
	{
		name:   "Warehouse",
		schema: resources.Warehouse().Schema,
	},
}
