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
	attributes := make([]ResourceAttributeAssertionModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		attributes = append(attributes, ResourceAttributeAssertionModel{
			Name: attr.Name,
			// TODO: add attribute type logic; allow type safe assertions, not only strings
			AttributeType: "string",
		})
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceAssertionsModel{
		Name:       resourceSchemaDetails.ObjectName(),
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
