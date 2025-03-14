package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (s *SchemasModel) WithLimit(rows int) *SchemasModel {
	return s.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}
