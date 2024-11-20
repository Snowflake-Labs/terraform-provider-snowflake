package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *TagModel) WithAllowedValues(allowedValues ...string) *TagModel {
	allowedValuesStringVariables := make([]tfconfig.Variable, len(allowedValues))
	for i, v := range allowedValues {
		allowedValuesStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.AllowedValues = tfconfig.SetVariable(allowedValuesStringVariables...)
	return t
}

func (t *TagModel) WithMaskingPolicies(maskingPolicies ...sdk.SchemaObjectIdentifier) *TagModel {
	maskingPoliciesStringVariables := make([]tfconfig.Variable, len(maskingPolicies))
	for i, v := range maskingPolicies {
		maskingPoliciesStringVariables[i] = tfconfig.StringVariable(v.FullyQualifiedName())
	}

	t.MaskingPolicies = tfconfig.SetVariable(maskingPoliciesStringVariables...)
	return t
}
