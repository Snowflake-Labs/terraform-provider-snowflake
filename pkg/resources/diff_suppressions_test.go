package resources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func Test_NormalizeAndCompare(t *testing.T) {
	genericNormalize := func(value string) (any, error) {
		switch value {
		case "ok", "ok1":
			return "ok", nil
		default:
			return nil, fmt.Errorf("incorrect value %s", value)
		}
	}

	t.Run("generic normalize", func(t *testing.T) {
		result := resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok1", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "nok", nil)
		assert.False(t, result)
	})

	t.Run("warehouse size", func(t *testing.T) {
		result := resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX4Large), nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "4X-LARGE", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX5Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "invalid", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "invalid", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)
	})
}
