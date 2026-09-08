package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (o *OpenflowRuntimesModel) WithLimit(rows int) *OpenflowRuntimesModel {
	return o.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (o *OpenflowRuntimesModel) WithInAccount() *OpenflowRuntimesModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (o *OpenflowRuntimesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *OpenflowRuntimesModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (o *OpenflowRuntimesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *OpenflowRuntimesModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}
