package gen

import (
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/gencommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	Name           string
	AttributeType  string
	Required       bool
	VariableMethod string
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails gencommons.ResourceSchemaDetails) ResourceConfigBuilderModel {
	attributes := make([]ResourceConfigBuilderAttributeModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}

		// TODO: support the rest of attribute types
		var attributeType string
		var variableMethod string
		switch attr.AttributeType {
		case schema.TypeBool:
			attributeType = "bool"
			variableMethod = "BoolVariable"
		case schema.TypeInt:
			attributeType = "int"
			variableMethod = "IntegerVariable"
		case schema.TypeFloat:
			attributeType = "float"
			variableMethod = "FloatVariable"
		case schema.TypeString:
			attributeType = "string"
			variableMethod = "StringVariable"
		}

		attributes = append(attributes, ResourceConfigBuilderAttributeModel{
			Name:           attr.Name,
			AttributeType:  attributeType,
			Required:       attr.Required,
			VariableMethod: variableMethod,
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
