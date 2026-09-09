package resources

import (
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestDirectoryTableOutputMapping_IgnoresLastRefreshedOnChanges(t *testing.T) {
	storedDescribeOutput := []any{
		map[string]any{
			"enable":            true,
			"auto_refresh":      false,
			"last_refreshed_on": "2026-09-08 12:00:00.000 -0700",
		},
	}

	directoryTable := sdk.StageDirectoryTable{
		Enable:          true,
		AutoRefresh:     false,
		LastRefreshedOn: sdk.String("2026-09-08 13:00:00.000 -0700"),
	}

	mapping := directoryTableOutputMapping(directoryTable)
	normalizedStored := mapping.normalizeFunc(storedDescribeOutput)

	require.True(t, reflect.DeepEqual(normalizedStored, mapping.valueToCompare))
}

func TestDirectoryTableOutputMapping_DetectsEnableChanges(t *testing.T) {
	storedDescribeOutput := []any{
		map[string]any{
			"enable":            false,
			"auto_refresh":      false,
			"last_refreshed_on": "2026-09-08 12:00:00.000 -0700",
		},
	}

	directoryTable := sdk.StageDirectoryTable{
		Enable:          true,
		AutoRefresh:     false,
		LastRefreshedOn: sdk.String("2026-09-08 12:00:00.000 -0700"),
	}

	mapping := directoryTableOutputMapping(directoryTable)
	normalizedStored := mapping.normalizeFunc(storedDescribeOutput)

	require.False(t, reflect.DeepEqual(normalizedStored, mapping.valueToCompare))
}
