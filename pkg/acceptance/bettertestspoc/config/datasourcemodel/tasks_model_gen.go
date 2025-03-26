// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type TasksModel struct {
	In             tfconfig.Variable `json:"in,omitempty"`
	Like           tfconfig.Variable `json:"like,omitempty"`
	Limit          tfconfig.Variable `json:"limit,omitempty"`
	RootOnly       tfconfig.Variable `json:"root_only,omitempty"`
	StartsWith     tfconfig.Variable `json:"starts_with,omitempty"`
	Tasks          tfconfig.Variable `json:"tasks,omitempty"`
	WithParameters tfconfig.Variable `json:"with_parameters,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Tasks(
	datasourceName string,
) *TasksModel {
	t := &TasksModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.Tasks)}
	return t
}

func TasksWithDefaultMeta() *TasksModel {
	t := &TasksModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.Tasks)}
	return t
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (t *TasksModel) MarshalJSON() ([]byte, error) {
	type Alias TasksModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(t),
		DependsOn:                 t.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (t *TasksModel) WithDependsOn(values ...string) *TasksModel {
	t.SetDependsOn(values...)
	return t
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// in attribute type is not yet supported, so WithIn can't be generated

func (t *TasksModel) WithLike(like string) *TasksModel {
	t.Like = tfconfig.StringVariable(like)
	return t
}

// limit attribute type is not yet supported, so WithLimit can't be generated

func (t *TasksModel) WithRootOnly(rootOnly bool) *TasksModel {
	t.RootOnly = tfconfig.BoolVariable(rootOnly)
	return t
}

func (t *TasksModel) WithStartsWith(startsWith string) *TasksModel {
	t.StartsWith = tfconfig.StringVariable(startsWith)
	return t
}

// tasks attribute type is not yet supported, so WithTasks can't be generated

func (t *TasksModel) WithWithParameters(withParameters bool) *TasksModel {
	t.WithParameters = tfconfig.BoolVariable(withParameters)
	return t
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (t *TasksModel) WithInValue(value tfconfig.Variable) *TasksModel {
	t.In = value
	return t
}

func (t *TasksModel) WithLikeValue(value tfconfig.Variable) *TasksModel {
	t.Like = value
	return t
}

func (t *TasksModel) WithLimitValue(value tfconfig.Variable) *TasksModel {
	t.Limit = value
	return t
}

func (t *TasksModel) WithRootOnlyValue(value tfconfig.Variable) *TasksModel {
	t.RootOnly = value
	return t
}

func (t *TasksModel) WithStartsWithValue(value tfconfig.Variable) *TasksModel {
	t.StartsWith = value
	return t
}

func (t *TasksModel) WithTasksValue(value tfconfig.Variable) *TasksModel {
	t.Tasks = value
	return t
}

func (t *TasksModel) WithWithParametersValue(value tfconfig.Variable) *TasksModel {
	t.WithParameters = value
	return t
}
