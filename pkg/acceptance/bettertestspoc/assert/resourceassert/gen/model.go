package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceAssertionsModel struct {
	Name string
	PreambleModel
}

func (m ResourceAssertionsModel) SomeFunc() {
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails gencommons.ResourceSchemaDetails) ResourceAssertionsModel {
	return ResourceAssertionsModel{}
}
