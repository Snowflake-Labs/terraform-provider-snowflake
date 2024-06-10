package resources

import (
	"context"
	"log"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AccountObjectStringValueComputedIf(key string, params []*sdk.Parameter, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(
		key,
		params,
		parameter,
		func(value any) string { return value.(string) },
		func(value string) (string, error) { return value, nil },
	)
}

func AccountObjectIntValueComputedIf(key string, params []*sdk.Parameter, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(
		key,
		params,
		parameter,
		func(value any) string { return strconv.Itoa(value.(int)) },
		strconv.Atoi,
	)
}

func AccountObjectBoolValueComputedIf(key string, params []*sdk.Parameter, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return ValueComputedIf(
		key,
		params,
		parameter,
		func(value any) string { return strconv.FormatBool(value.(bool)) },
		strconv.ParseBool,
	)
}

func ValueComputedIf[T any](key string, parameters []*sdk.Parameter, accountParameter sdk.AccountParameter, valueToString func(v any) string, valueFromString func(value string) (T, error)) schema.CustomizeDiffFunc {
	var parameterValue *string
	for _, parameter := range parameters {
		if parameter.Key == string(accountParameter) {
			parameterValue = &parameter.Value
			break
		}
	}

	condition := func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
		configValue, ok := d.GetRawConfig().AsValueMap()[key]
		if ok && !configValue.IsNull() {
			return false
		}

		if parameterValue == nil {
			log.Printf("[ERROR] ValueComputedIf, parameter %s not found", accountParameter)
			return false
		}

		return *parameterValue != valueToString(d.Get(key))
	}

	return func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
		if condition(ctx, d, meta) {
			if *parameterValue == "" {
				return d.SetNew(key, "<null>")
			} else {
				value, err := valueFromString(*parameterValue)
				if err != nil {
					return err
				}
				return d.SetNew(key, value)
			}
		}

		return nil
	}
}
