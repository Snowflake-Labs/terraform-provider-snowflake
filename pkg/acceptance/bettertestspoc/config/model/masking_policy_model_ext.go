package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

func (p *MaskingPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *MaskingPolicyModel {
	maps := make([]config.Variable, len(argument))
	for i, v := range argument {
		maps[i] = config.MapVariable(map[string]config.Variable{
			"name": config.StringVariable(v.Name),
			"type": config.StringVariable(string(v.Type)),
		})
	}
	p.Argument = tfconfig.SetVariable(maps...)
	return p
}
