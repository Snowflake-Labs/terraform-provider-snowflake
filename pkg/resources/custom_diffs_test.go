package resources_test

import (
	"context"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueComputedIf(t *testing.T) {
	createProviderConfig := func(parameterLevel sdk.ParameterType, parameterValue sdk.LogLevel) *schema.Provider {
		customDiff := resources.ValueComputedIf(
			"value",
			[]*sdk.Parameter{
				{
					Key:   string(sdk.AccountParameterLogLevel),
					Level: parameterLevel,
					Value: string(parameterValue),
				},
			},
			sdk.ParameterTypeDatabase,
			sdk.AccountParameterLogLevel,
			func(v any) string { return v.(string) },
		)
		return createProviderWithValuePropertyAndCustomDiff(t, schema.TypeString, customDiff)
	}

	t.Run("config: true - state: true - level: different - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: different - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: same - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeDatabase, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: same - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeDatabase, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: different - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: different - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: same - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: same - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	// Tests for filled config and empty state were not added as the only way
	// of getting into this situation would be in create operation for which custom diffs are skipped.
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
