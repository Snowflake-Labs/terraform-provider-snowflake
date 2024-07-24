package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func BasicWarehouseModel(
	name string,
	comment string,
) *WarehouseModel {
	return WarehouseWithDefaultMeta(name).WithComment(comment)
}

// TODO [SNOW-1501905]: currently config builder are generated from the resource schema, so there is no direct connection to the source enum (like sdk.WarehouseSize)
// For now, we can just add extension methods manually.
// Later, we could provide type overrides map or even SDK object to automatically match by name.
func (w *WarehouseModel) WithWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *WarehouseModel {
	return w.WithWarehouseSize(string(warehouseSize))
}
