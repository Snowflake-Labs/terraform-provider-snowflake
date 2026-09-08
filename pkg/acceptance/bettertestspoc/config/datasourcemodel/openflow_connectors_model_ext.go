package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (o *OpenflowConnectorsModel) WithLimit(rows int) *OpenflowConnectorsModel {
	return o.WithLimitValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"rows": tfconfig.IntegerVariable(rows),
		}),
	)
}

func (o *OpenflowConnectorsModel) WithInAccount() *OpenflowConnectorsModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"account": tfconfig.BoolVariable(true),
		}),
	)
}

func (o *OpenflowConnectorsModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *OpenflowConnectorsModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (o *OpenflowConnectorsModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *OpenflowConnectorsModel {
	return o.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}
