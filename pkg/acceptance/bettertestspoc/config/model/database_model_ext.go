package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (d *DatabaseModel) WithReplication(accountIdentifier sdk.AccountIdentifier, withFailover bool, ignoreEditionCheck bool) *DatabaseModel {
	return d.WithReplicationValue(
		tfconfig.ObjectVariable(
			map[string]tfconfig.Variable{
				"enable_to_account": tfconfig.ObjectVariable(
					map[string]tfconfig.Variable{
						"account_identifier": tfconfig.StringVariable(accountIdentifier.FullyQualifiedName()),
						"with_failover":      tfconfig.BoolVariable(withFailover),
					},
				),
				"ignore_edition_check": tfconfig.BoolVariable(ignoreEditionCheck),
			},
		),
	)
}
