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

// TODO: test this custom diff func
func UpdateValueWithSnowflakeDefault(key string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		sfStateKey := key + "_sf"
		needsRefreshKey := key + "_sf_changed"

		configValue, ok := d.GetRawConfig().AsValueMap()[key]
		stateValue := d.Get(key).(string)
		_, needsRefresh := d.GetChange(needsRefreshKey)

		if needsRefresh.(bool) && stateValue == "" && (!ok || configValue.IsNull()) {
			err := d.SetNew(needsRefreshKey, false)
			if err != nil {
				return err
			}
			err = d.SetNewComputed(sfStateKey)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
