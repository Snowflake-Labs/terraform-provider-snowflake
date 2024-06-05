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

type Property interface {
	GetName() string
	GetDefault() string
}

func BoolComputedIf(key, property string, describe func(client *sdk.Client, id sdk.AccountObjectIdentifier) ([]Property, error)) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		configValue := d.GetRawConfig().AsValueMap()[key]
		if !configValue.IsNull() {
			return false
		}

		client := meta.(*provider.Context).Client

		props, err := describe(client, sdk.NewAccountObjectIdentifier(d.Id()))
		if err != nil {
			return false
		}
		var def string
		for _, v := range props {
			if v.GetName() == property {
				def = v.GetDefault()
				break
			}
		}
		if def == "" {
			return false
		}
		stateValue := d.Get(key).(bool)
		return def != strconv.FormatBool(stateValue)
	})
}
