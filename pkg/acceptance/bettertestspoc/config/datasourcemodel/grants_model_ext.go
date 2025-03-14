package datasourcemodel

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func GrantsOnAccount(
	datasourceName string,
) *GrantsModel {
	return Grants(datasourceName).
		WithGrantsOnValue(
			tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"account": tfconfig.BoolVariable(true),
			}),
		)
}
