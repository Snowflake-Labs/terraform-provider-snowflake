package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func v091ToWarehouseSize(s string) (sdk.WarehouseSize, error) {
	s = strings.ToUpper(s)
	switch s {
	case "XSMALL", "X-SMALL":
		return sdk.WarehouseSizeXSmall, nil
	case "SMALL":
		return sdk.WarehouseSizeSmall, nil
	case "MEDIUM":
		return sdk.WarehouseSizeMedium, nil
	case "LARGE":
		return sdk.WarehouseSizeLarge, nil
	case "XLARGE", "X-LARGE":
		return sdk.WarehouseSizeXLarge, nil
	case "XXLARGE", "X2LARGE", "2X-LARGE", "2XLARGE":
		return sdk.WarehouseSizeXXLarge, nil
	case "XXXLARGE", "X3LARGE", "3X-LARGE", "3XLARGE":
		return sdk.WarehouseSizeXXXLarge, nil
	case "X4LARGE", "4X-LARGE", "4XLARGE":
		return sdk.WarehouseSizeX4Large, nil
	case "X5LARGE", "5X-LARGE", "5XLARGE":
		return sdk.WarehouseSizeX5Large, nil
	case "X6LARGE", "6X-LARGE", "6XLARGE":
		return sdk.WarehouseSizeX6Large, nil
	default:
		return "", fmt.Errorf("invalid warehouse size: %s", s)
	}
}

// v091WarehouseSizeStateUpgrader is needed because we are removing incorrect mapped values from sdk.ToWarehouseSize (like 2XLARGE, 3XLARGE, ...)
// Result of:
// - https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/1873
// - https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/1946
// - https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1889#issuecomment-1631149585
func v091WarehouseSizeStateUpgrader(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldWarehouseSize := rawState["warehouse_size"].(string)
	if oldWarehouseSize == "" {
		return rawState, nil
	}

	warehouseSize, err := v091ToWarehouseSize(oldWarehouseSize)
	if err != nil {
		return nil, err
	}
	rawState["warehouse_size"] = string(warehouseSize)

	// TODO: clear wait_for_provisioning and test

	return rawState, nil
}
