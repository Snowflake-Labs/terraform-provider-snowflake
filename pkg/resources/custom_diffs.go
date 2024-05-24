package resources

import (
	"context"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NestedIntValueAccountObjectComputedIf(key string, parameter sdk.AccountParameter) schema.CustomizeDiffFunc {
	return NestedValueComputedIf(
		key,
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return client.Parameters.ShowAccountParameter(context.Background(), parameter)
		},
		func(v any) string { return strconv.Itoa(v.(int)) },
	)
}

func NestedValueComputedIf(key string, showParam func(client *sdk.Client) (*sdk.Parameter, error), valueToString func(v any) string) schema.CustomizeDiffFunc {
	return customdiff.ComputedIf(key, func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		configValue := d.GetRawConfig().AsValueMap()[key].AsValueSlice()

		if len(configValue) == 1 {
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
