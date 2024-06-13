package resources

import (
	"context"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AccountObjectStringValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return value.(string) })
}

func AccountObjectIntValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return strconv.Itoa(value.(int)) })
}

func AccountObjectBoolValueComputedIf(key string, params []*sdk.Parameter, parameterLevel sdk.ParameterType, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(key, params, parameterLevel, parameter, func(value any) string { return strconv.FormatBool(value.(bool)) })
}

func ValueComputedIf(key string, parameters []*sdk.Parameter, objectParameterLevel sdk.ParameterType, accountParameter sdk.AccountParameter, valueToString func(v any) string) schema.CustomizeDiffFunc {
	var parameterValue *string
	var parameterLevel *sdk.ParameterType

	for _, parameter := range parameters {
		if parameter.Key == string(accountParameter) {
			parameterLevel = &parameter.Level
			parameterValue = &parameter.Value
			break
		}
	}

	condition := func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
		configValue, ok := d.GetRawConfig().AsValueMap()[key]

		if parameterLevel == nil || parameterValue == nil {
			log.Printf("[ERROR] ValueComputedIf, parameter %s not found", accountParameter)
			return false
		}

		// For cases where currently set value (in the config) is equal to the parameter, but not set on the right level.
		// The parameter is set somewhere higher in the hierarchy, and we need to "forcefully" set the value to
		// perform the actual set on Snowflake (and set the parameter on the correct level).
		if *parameterLevel != objectParameterLevel && !configValue.IsNull() && *parameterValue == valueToString(d.Get(key)) {
			return true
		}

		// For all other cases, if a parameter is set in the configuration, we can ignore parts needed for Computed fields.
		if ok && !configValue.IsNull() {
			return false
		}

		// If the configuration is not set, perform SetNewComputed for cases like:
		// 1. Check if the parameter value differs from the one saved in state (if they differ, we'll update the computed value).
		// 2. Check if the parameter level is set on the same level for the given object (if they're the same, we'll trigger the unset as it has been updated in Snowflake).
		return *parameterValue != valueToString(d.Get(key)) || *parameterLevel == objectParameterLevel
	}

	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		if condition(ctx, d, meta) {
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
