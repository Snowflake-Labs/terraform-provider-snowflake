package resources_test

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestNestedValueComputedIf(t *testing.T) {
	customDiff := resources.ValueComputedIf[string](
		"value",
		[]*sdk.Parameter{
			{
				Key:   string(sdk.AccountParameterLogLevel),
				Value: string(sdk.LogLevelInfo),
			},
		},
		sdk.AccountParameterLogLevel,
		func(v any) string { return v.(string) },
		func(v string) string { return v },
	)
	providerConfig := createProviderWithValuePropertyAndCustomDiff(t, schema.TypeString, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelDebug),
		})
		assert.Equal(t, string(sdk.LogLevelInfo), diff.Attributes["value"].New)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})
}

func createProviderWithValuePropertyAndCustomDiff(t *testing.T, valueType schema.ValueType, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     valueType,
						Computed: true,
						Optional: true,
					},
				},
				CustomizeDiff: customDiffFunc,
			},
		},
	}
}

func calculateDiff(t *testing.T, providerConfig *schema.Provider, rawConfigValue cty.Value, stateValue map[string]any) *terraform.InstanceDiff {
	t.Helper()
	diff, err := providerConfig.ResourcesMap["test"].Diff(
		context.Background(),
		&terraform.InstanceState{
			RawConfig: rawConfigValue,
		},
		&terraform.ResourceConfig{
			Config: stateValue,
		},
		&provider.Context{Client: acc.Client(t)},
	)
	require.NoError(t, err)
	return diff
}
