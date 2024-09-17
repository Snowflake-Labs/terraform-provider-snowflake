package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (r *RowAccessPolicyModel) WithArgument(argument []sdk.TableColumnSignature) *RowAccessPolicyModel {
	maps := make([]config.Variable, len(argument))
	for i, v := range argument {
		maps[i] = config.MapVariable(map[string]config.Variable{
			"name": config.StringVariable(v.Name),
			"type": config.StringVariable(string(v.Type)),
		})
	}
	r.Argument = tfconfig.SetVariable(maps...)
	return r
}
