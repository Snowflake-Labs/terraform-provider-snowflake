package resources

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StringParameterValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return value.(string) })
}

func IntParameterValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return strconv.Itoa(value.(int)) })
}

func BoolParameterValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ParameterValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return strconv.FormatBool(value.(bool)) })
}

func ParameterValueComputedIf(key string, parameters []*sdk.Parameter, objectParameterLevel sdk.ParameterType, accountParameter sdk.AccountParameter, valueToString func(v any) string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		foundParameter, err := collections.FindOne(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == string(accountParameter) })
		if err != nil {
			log.Printf("[WARN] failed to find account parameter: %s", accountParameter)
			return nil
		}
		parameter := *foundParameter

		configValue, ok := d.GetRawConfig().AsValueMap()[key]

		// For cases where currently set value (in the config) is equal to the parameter, but not set on the right level.
		// The parameter is set somewhere higher in the hierarchy, and we need to "forcefully" set the value to
		// perform the actual set on Snowflake (and set the parameter on the correct level).
		if ok && !configValue.IsNull() && parameter.Level != objectParameterLevel && parameter.Value == valueToString(d.Get(key)) {
			return d.SetNewComputed(key)
		}

		// For all other cases, if a parameter is set in the configuration, we can ignore parts needed for Computed fields.
		if ok && !configValue.IsNull() {
			return nil
		}

		// If the configuration is not set, perform SetNewComputed for cases like:
		// 1. Check if the parameter value differs from the one saved in state (if they differ, we'll update the computed value).
		// 2. Check if the parameter is set on the object level (if so, it means that it was set externally, and we have to unset it).
		if parameter.Value != valueToString(d.Get(key)) || parameter.Level == objectParameterLevel {
			return d.SetNewComputed(key)
		}

		return nil
	}
}

// ForceNewIfChangeToEmptySlice sets a ForceNew for a list field which was set to an empty value.
func ForceNewIfChangeToEmptySlice[T any](key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.([]T), newValue.([]T)
		return len(oldList) > 0 && len(newList) == 0
	})
}

// ForceNewIfChangeToEmptySet sets a ForceNew for a list field which was set to an empty value.
func ForceNewIfChangeToEmptySet(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.(*schema.Set).List(), newValue.(*schema.Set).List()
		return len(oldList) > 0 && len(newList) == 0
	})
}

// ForceNewIfChangeToEmptyString sets a ForceNew for a string field which was set to an empty value.
func ForceNewIfChangeToEmptyString(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldString, newString := oldValue.(string), newValue.(string)
		return len(oldString) > 0 && len(newString) == 0
	})
}

// TODO [follow-up PR]: test
func ComputedIfAnyAttributeChanged(key string, changedAttributeKeys ...string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
		var result bool
		for _, changedKey := range changedAttributeKeys {
			if diff.HasChange(changedKey) {
				log.Printf("[DEBUG] ComputedIfAnyAttributeChanged: changed key: %s\n", changedKey)
			}
			result = result || diff.HasChange(changedKey)
		}
		return result
	})
}

// ForceNewIfChangeToEmptyString sets a ForceNew for a string field which was set to an empty value.
func ForceNewIfChangeToEmptyString(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldString, newString := oldValue.(string), newValue.(string)
		return len(oldString) > 0 && len(newString) == 0
	})
}

// ForceNewIfChangeToEmptySet sets a ForceNew for a set field which was set to an empty value.
func ForceNewIfChangeToEmptySet(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.(*schema.Set).List(), newValue.(*schema.Set).List()
		return len(oldList) > 0 && len(newList) == 0
	})
}

type parameter struct {
	parameterName sdk.AccountParameter
	valueType     valueType
	parameterType sdk.ParameterType
}

type valueType string

const (
	valueTypeInt    valueType = "int"
	valueTypeBool   valueType = "bool"
	valueTypeString valueType = "string"
)

type ResourceIdProvider interface {
	Id() string
}

func ParametersCustomDiff(parametersProvider func(context.Context, ResourceIdProvider, any) ([]*sdk.Parameter, error), parameters ...parameter) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		if d.Id() == "" {
			return nil
		}

		params, err := parametersProvider(ctx, d, meta)
		if err != nil {
			return err
		}

		diffFunctions := make([]schema.CustomizeDiffFunc, len(parameters))
		for idx, p := range parameters {
			var diffFunction schema.CustomizeDiffFunc
			switch p.valueType {
			case valueTypeInt:
				diffFunction = IntParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			case valueTypeBool:
				diffFunction = BoolParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			case valueTypeString:
				diffFunction = StringParameterValueComputedIf(strings.ToLower(string(p.parameterName)), params, p.parameterType, p.parameterName)
			}
			diffFunctions[idx] = diffFunction
		}

		return customdiff.All(diffFunctions...)(ctx, d, meta)
	}
}
