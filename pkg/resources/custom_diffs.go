package resources

import (
	"context"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func BoolComputedIf(key string, getDefault func(client *sdk.Client, id sdk.AccountObjectIdentifier) (string, error)) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		configValue := d.GetRawConfig().AsValueMap()[key]
		if !configValue.IsNull() {
			return false
		}

		client := meta.(*provider.Context).Client

		def, err := getDefault(client, helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier))
		if err != nil {
			return false
		}
		stateValue := d.Get(key).(bool)
		return def != strconv.FormatBool(stateValue)
	})
}

// SetEmptyForceNewIfChange sets a ForceNew for a list field which was set to an empty value.
func SetEmptyForceNewIfChange[T any](key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldList, newList := oldValue.([]T), newValue.([]T)
		return len(oldList) > 0 && len(newList) == 0
	})
}

// SetEmptyForceNewIfChangeString sets a ForceNew for a string field which was set to an empty value.
func SetEmptyForceNewIfChangeString(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, func(ctx context.Context, oldValue, newValue, meta any) bool {
		oldString, newString := oldValue.(string), newValue.(string)
		return len(oldString) > 0 && len(newString) == 0
	})
}

// SetForceNewIfNull sets a ForceNew for key missing in RawConfig.
func SetForceNewIfNull(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		return d.GetRawConfig().AsValueMap()[key].IsNull()
	})
}

// TODO [follow-up PR]: test
func ComputedIfAnyAttributeChanged(key string, changedAttributeKeys ...string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
		var result bool
		for _, changedKey := range changedAttributeKeys {
			result = result || diff.HasChange(changedKey)
		}
		return result
	})
}
