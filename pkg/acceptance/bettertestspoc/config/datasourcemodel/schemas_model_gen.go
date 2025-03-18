// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type SchemasModel struct {
	In             tfconfig.Variable `json:"in,omitempty"`
	Like           tfconfig.Variable `json:"like,omitempty"`
	Limit          tfconfig.Variable `json:"limit,omitempty"`
	Schemas        tfconfig.Variable `json:"schemas,omitempty"`
	StartsWith     tfconfig.Variable `json:"starts_with,omitempty"`
	WithDescribe   tfconfig.Variable `json:"with_describe,omitempty"`
	WithParameters tfconfig.Variable `json:"with_parameters,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Schemas(
	datasourceName string,
) *SchemasModel {
	s := &SchemasModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.Schemas)}
	return s
}

func SchemasWithDefaultMeta() *SchemasModel {
	s := &SchemasModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.Schemas)}
	return s
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (s *SchemasModel) MarshalJSON() ([]byte, error) {
	type Alias SchemasModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(s),
		DependsOn:                 s.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (s *SchemasModel) WithDependsOn(values ...string) *SchemasModel {
	s.SetDependsOn(values...)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// in attribute type is not yet supported, so WithIn can't be generated

func (s *SchemasModel) WithLike(like string) *SchemasModel {
	s.Like = tfconfig.StringVariable(like)
	return s
}

// limit attribute type is not yet supported, so WithLimit can't be generated

// schemas attribute type is not yet supported, so WithSchemas can't be generated

func (s *SchemasModel) WithStartsWith(startsWith string) *SchemasModel {
	s.StartsWith = tfconfig.StringVariable(startsWith)
	return s
}

func (s *SchemasModel) WithWithDescribe(withDescribe bool) *SchemasModel {
	s.WithDescribe = tfconfig.BoolVariable(withDescribe)
	return s
}

func (s *SchemasModel) WithWithParameters(withParameters bool) *SchemasModel {
	s.WithParameters = tfconfig.BoolVariable(withParameters)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *SchemasModel) WithInValue(value tfconfig.Variable) *SchemasModel {
	s.In = value
	return s
}

func (s *SchemasModel) WithLikeValue(value tfconfig.Variable) *SchemasModel {
	s.Like = value
	return s
}

func (s *SchemasModel) WithLimitValue(value tfconfig.Variable) *SchemasModel {
	s.Limit = value
	return s
}

func (s *SchemasModel) WithSchemasValue(value tfconfig.Variable) *SchemasModel {
	s.Schemas = value
	return s
}

func (s *SchemasModel) WithStartsWithValue(value tfconfig.Variable) *SchemasModel {
	s.StartsWith = value
	return s
}

func (s *SchemasModel) WithWithDescribeValue(value tfconfig.Variable) *SchemasModel {
	s.WithDescribe = value
	return s
}

func (s *SchemasModel) WithWithParametersValue(value tfconfig.Variable) *SchemasModel {
	s.WithParameters = value
	return s
}
