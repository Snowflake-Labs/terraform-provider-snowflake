package gen

import (
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceAssertionsModel struct {
	Name       string
	Attributes []ResourceAttributeAssertionModel
	PreambleModel
}

func (m ResourceAssertionsModel) SomeFunc() {
}

type ResourceAttributeAssertionModel struct {
	Name          string
	AttributeType string
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails gencommons.ResourceSchemaDetails) ResourceAssertionsModel {
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceAssertionsModel{
		Name: resourceSchemaDetails.ObjectName(),
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
