package resources

import (
	"context"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NestedIntValueAccountObjectComputedIf is NestedValueComputedIf,
// but dedicated for account level objects with integer-typed properties.
func NestedIntValueAccountObjectComputedIf(key string, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return NestedValueComputedIf(
		key,
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return client.Parameters.ShowAccountParameter(context.Background(), parameter)
		},
		func(v any) string { return strconv.Itoa(v.(int)) },
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

func NormalizeAndCompare[T comparable](normalize func(string) (T, error)) schema.SchemaDiffSuppressFunc {
	return func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
		oldNormalized, err := normalize(oldValue)
		if err != nil {
			return false
		}
		newNormalized, err := normalize(newValue)
		if err != nil {
			return false
		}
		return oldNormalized == newNormalized
	}
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
