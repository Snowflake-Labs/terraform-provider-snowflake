// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type WarehousesModel struct {
	Like           tfconfig.Variable `json:"like,omitempty"`
	Warehouses     tfconfig.Variable `json:"warehouses,omitempty"`
	WithDescribe   tfconfig.Variable `json:"with_describe,omitempty"`
	WithParameters tfconfig.Variable `json:"with_parameters,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Warehouses(
	datasourceName string,
) *WarehousesModel {
	w := &WarehousesModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.Warehouses)}
	return w
}

func WarehousesWithDefaultMeta() *WarehousesModel {
	w := &WarehousesModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.Warehouses)}
	return w
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (w *WarehousesModel) MarshalJSON() ([]byte, error) {
	type Alias WarehousesModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(w),
		DependsOn:                 w.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (w *WarehousesModel) WithDependsOn(values ...string) *WarehousesModel {
	w.SetDependsOn(values...)
	return w
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (w *WarehousesModel) WithLike(like string) *WarehousesModel {
	w.Like = tfconfig.StringVariable(like)
	return w
}

// warehouses attribute type is not yet supported, so WithWarehouses can't be generated

func (w *WarehousesModel) WithWithDescribe(withDescribe bool) *WarehousesModel {
	w.WithDescribe = tfconfig.BoolVariable(withDescribe)
	return w
}

func (w *WarehousesModel) WithWithParameters(withParameters bool) *WarehousesModel {
	w.WithParameters = tfconfig.BoolVariable(withParameters)
	return w
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (w *WarehousesModel) WithLikeValue(value tfconfig.Variable) *WarehousesModel {
	w.Like = value
	return w
}

func (w *WarehousesModel) WithWarehousesValue(value tfconfig.Variable) *WarehousesModel {
	w.Warehouses = value
	return w
}

func (w *WarehousesModel) WithWithDescribeValue(value tfconfig.Variable) *WarehousesModel {
	w.WithDescribe = value
	return w
}

func (w *WarehousesModel) WithWithParametersValue(value tfconfig.Variable) *WarehousesModel {
	w.WithParameters = value
	return w
}
