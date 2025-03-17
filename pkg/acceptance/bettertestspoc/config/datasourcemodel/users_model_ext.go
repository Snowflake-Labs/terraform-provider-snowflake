package datasourcemodel

import tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

func (u *UsersModel) WithLimitRowsAndFrom(rows int, from string) *UsersModel {
	return u.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
			"from": tfconfig.StringVariable(from),
		}),
	)
}
