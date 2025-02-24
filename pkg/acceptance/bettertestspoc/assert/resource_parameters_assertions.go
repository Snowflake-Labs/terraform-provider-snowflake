package assert

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	parametersPrefix      = "parameters.0."
	parametersValueSuffix = ".0.value"
	parametersLevelSuffix = ".0.level"
)

func ResourceParameterBoolValueSet[T ~string](parameterName T, expected bool) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, strconv.FormatBool(expected))
}

func ResourceParameterBoolValueNotSet[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValueNotSet(parameterName)
}

func ResourceParameterBoolValuePresent[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValuePresent(parameterName)
}

func ResourceParameterIntValueSet[T ~string](parameterName T, expected int) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, strconv.Itoa(expected))
}

func ResourceParameterIntValueNotSet[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValueNotSet(parameterName)
}

func ResourceParameterIntValuePresent[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValuePresent(parameterName)
}

func ResourceParameterStringUnderlyingValueSet[T ~string, U ~string](parameterName T, expected U) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, string(expected))
}

func ResourceParameterStringUnderlyingValueNotSet[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValueNotSet(parameterName)
}

func ResourceParameterStringUnderlyingValuePresent[T ~string](parameterName T) ResourceAssertion {
	return ResourceParameterValuePresent(parameterName)
}

func ResourceParameterValueSet[T ~string](parameterName T, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersValueSuffix, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceParameterValueNotSet[T ~string](parameterName T) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersValueSuffix, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

func ResourceParameterValuePresent[T ~string](parameterName T) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersValueSuffix, resourceAssertionType: resourceAssertionTypeValuePresent}
}

func ResourceParameterLevelSet[T ~string](parameterName T, parameterType sdk.ParameterType) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersLevelSuffix, expectedValue: string(parameterType), resourceAssertionType: resourceAssertionTypeValueSet}
}
