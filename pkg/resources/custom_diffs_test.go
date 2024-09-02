package resources_test

import (
	"context"
	"strings"
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

func TestParameterValueComputedIf(t *testing.T) {
	createProviderConfig := func(parameterLevel sdk.ParameterType, parameterValue sdk.LogLevel) *schema.Provider {
		customDiff := resources.ParameterValueComputedIf(
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
		return createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		}, customDiff)
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

func createProviderWithValuePropertyAndCustomDiff(t *testing.T, valueSchema *schema.Schema, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					"value": valueSchema,
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

func calculateDiffFromAttributes(t *testing.T, providerConfig *schema.Provider, oldValue map[string]string, newValue map[string]any) *terraform.InstanceDiff {
	t.Helper()
	diff, err := providerConfig.ResourcesMap["test"].Diff(
		context.Background(),
		&terraform.InstanceState{
			Attributes: oldValue,
		},
		&terraform.ResourceConfig{
			Config: newValue,
		},
		&provider.Context{Client: acc.Client(t)},
	)
	require.NoError(t, err)
	return diff
}

func TestForceNewIfChangeToEmptyString(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": "foo",
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value": "foo",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value": "bar",
			},
			rawConfigValue: map[string]any{
				"value": "foo",
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.ForceNewIfChangeToEmptyString(
				"value",
			)
			provider := createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				provider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func TestForceNewIfChangeToEmptySlice(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value.#": "1",
				"value.0": "foo",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value.#": "2",
				"value.0": "foo",
				"value.1": "bar",
			},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.ForceNewIfChangeToEmptySlice[any](
				"value",
			)
			provider := createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				provider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func TestForceNewIfChangeToEmptySet(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value.#": "1",
				// The Sets are using hashes to generate an index for a given value.
				// In this case: 2577344683 == hash("CREATE DATABASE").
				"value.2577344683": "CREATE DATABASE",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value.#": "2",
				"value.0": "foo",
				"value.1": "bar",
			},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := calculateDiffFromAttributes(t,
				createProviderWithValuePropertyAndCustomDiff(t,
					&schema.Schema{
						Type: schema.TypeSet,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					resources.ForceNewIfChangeToEmptySet(
						"value",
					),
				),
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_ComputedIfAnyAttributeChanged(t *testing.T) {
	testSuppressFunc := schema.SchemaDiffSuppressFunc(func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		return strings.Trim(oldValue, `"`) == strings.Trim(newValue, `"`)
	})
	testSchema := map[string]*schema.Schema{
		"value_with_diff_suppress": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: testSuppressFunc,
		},
		"value_without_diff_suppress": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"computed_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	testCustomDiff := resources.ComputedIfAnyAttributeChanged(
		testSchema,
		"computed_value",
		"value_with_diff_suppress",
		"value_without_diff_suppress",
	)
	testProvider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema:        testSchema,
				CustomizeDiff: testCustomDiff,
			},
		},
	}

	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		expectDiff     bool
	}{
		{
			name: "no change on both fields",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - suppressed (quotes in config added)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - suppressed (quotes in config removed)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - not suppressed (value change)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "bar",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: true,
		},
		{
			name: "change on field with diff suppression - not suppressed (value and quotes changed)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "bar",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression, suppressed change on field with diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression, not suppressed change on field with diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"bar\"",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			if tt.expectDiff {
				require.NotNil(t, diff)
				require.NotNil(t, diff.Attributes["computed_value"])
				assert.True(t, diff.Attributes["computed_value"].NewComputed)
			} else {
				require.Nil(t, diff)
			}
		})
	}

	t.Run("attributes not found in schema, both fields changed", func(t *testing.T) {
		otherTestSchema := map[string]*schema.Schema{
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: testSuppressFunc,
			},
			"computed_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}
		otherTestCustomDiff := resources.ComputedIfAnyAttributeChanged(
			otherTestSchema,
			"computed_value",
			"value_with_diff_suppress",
			"value_without_diff_suppress",
		)
		otherTestProvider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"test": {
					Schema:        testSchema,
					CustomizeDiff: otherTestCustomDiff,
				},
			},
		}

		diff := calculateDiffFromAttributes(
			t,
			otherTestProvider,
			map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			map[string]any{
				"value_with_diff_suppress":    "\"bar\"",
				"value_without_diff_suppress": "bar",
			},
		)

		require.NotNil(t, diff)
		assert.Nil(t, diff.Attributes["computed_value"])
	})
}
