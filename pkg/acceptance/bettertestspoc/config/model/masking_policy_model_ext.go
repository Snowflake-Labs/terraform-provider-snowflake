package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (p *MaskingPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *MaskingPolicyModel {
	maps := make([]tfconfig.Variable, len(argument))
	for i, v := range argument {
		maps[i] = tfconfig.MapVariable(map[string]tfconfig.Variable{
			"name": tfconfig.StringVariable(v.Name),
			"type": tfconfig.StringVariable(string(v.Type)),
		})
	}
	p.Argument = tfconfig.SetVariable(maps...)
	return p
}
