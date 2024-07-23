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

type ResourceConfigBuilderModel struct {
	Name string
	PreambleModel
}

func (m ResourceConfigBuilderModel) SomeFunc() {
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails gencommons.ResourceSchemaDetails) ResourceConfigBuilderModel {
	//attributes := make([]ResourceAttributeAssertionModel, 0)
	//for _, attr := range resourceSchemaDetails.Attributes {
	//	if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
	//		continue
	//	}
	//	attributes = append(attributes, ResourceAttributeAssertionModel{
	//		Name: attr.Name,
	//		// TODO: add attribute type logic; allow type safe assertions, not only strings
	//		AttributeType: "string",
	//	})
	//}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceConfigBuilderModel{
		Name: resourceSchemaDetails.ObjectName(),
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
