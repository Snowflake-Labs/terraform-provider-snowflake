package gen

import (
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceConfigBuilderModel struct {
	Name       string
	Attributes []ResourceConfigBuilderAttributeModel
	PreambleModel
}

func (m ResourceConfigBuilderModel) SomeFunc() {
}

type ResourceConfigBuilderAttributeModel struct {
	Name          string
	AttributeType string
	Required      bool
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails gencommons.ResourceSchemaDetails) ResourceConfigBuilderModel {
	attributes := make([]ResourceConfigBuilderAttributeModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}
		attributes = append(attributes, ResourceConfigBuilderAttributeModel{
			Name: attr.Name,
			// TODO: set attribute type to a proper value
			AttributeType: "string",
			Required:      attr.Required,
		})
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceConfigBuilderModel{
		Name:       resourceSchemaDetails.ObjectName(),
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
