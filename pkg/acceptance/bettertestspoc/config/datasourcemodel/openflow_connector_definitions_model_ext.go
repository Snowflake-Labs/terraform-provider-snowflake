package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (o *OpenflowConnectorDefinitionsModel) WithLimit(rows int) *OpenflowConnectorDefinitionsModel {
	return o.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}
