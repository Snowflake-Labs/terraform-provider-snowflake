package resources

import (
	"context"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NestedIntValueAccountObjectComputedIf is NestedValueComputedIf,
// but dedicated for account level objects with integer properties.
func NestedIntValueAccountObjectComputedIf(key string, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return NestedValueComputedIf(
		key,
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return client.Parameters.ShowAccountParameter(context.Background(), parameter)
		},
		func(v any) string { return strconv.Itoa(v.(int)) },
	)
}

// NestedStringValueAccountObjectComputedIf is NestedValueComputedIf,
// but dedicated for account level objects with string properties.
func NestedStringValueAccountObjectComputedIf(key string, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return NestedValueComputedIf(
		key,
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return client.Parameters.ShowAccountParameter(context.Background(), parameter)
		},
		func(v any) string { return v.(string) },
	)
}

// NestedBoolValueAccountObjectComputedIf is NestedValueComputedIf,
// but dedicated for account level objects with bool properties.
func NestedBoolValueAccountObjectComputedIf(key string, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return NestedValueComputedIf(
		key,
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return client.Parameters.ShowAccountParameter(context.Background(), parameter)
		},
		func(v any) string {
			return strconv.FormatBool(v.(bool))
		},
	)
}

// NestedValueComputedIf internally calls schema.ResourceDiff.SetNewComputed whenever the inner function returns true.
// It's main purpose was to use it with hierarchical values that are marked with Computed and Optional. Such values should
// be recomputed whenever the value is not in the configuration and the remote value is not equal to the value in state.
func NestedValueComputedIf(key string, showParam func(client *sdk.Client) (*sdk.Parameter, error), valueToString func(v any) string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		configValue, ok := d.GetRawConfig().AsValueMap()[key]
		if ok && len(configValue.AsValueSlice()) == 1 {
			return false
		}

		client := meta.(*provider.Context).Client

		param, err := showParam(client)
		if err != nil {
			return false
		}

		stateValue := d.Get(key).([]any)
		if len(stateValue) != 1 {
			return false
		}

		return param.Value != valueToString(stateValue[0].(map[string]any)["value"])
	})
}
