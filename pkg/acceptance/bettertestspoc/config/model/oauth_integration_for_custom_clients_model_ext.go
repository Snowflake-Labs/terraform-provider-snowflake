package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *OauthIntegrationForCustomClientsModel) WithBlockedRolesList(blockedRoles ...string) *OauthIntegrationForCustomClientsModel {
	blockedRolesListStringVariables := make([]tfconfig.Variable, len(blockedRoles))
	for i, v := range blockedRoles {
		blockedRolesListStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.BlockedRolesList = tfconfig.SetVariable(blockedRolesListStringVariables...)
	return t
}

func (t *OauthIntegrationForCustomClientsModel) WithPreAuthorizedRoles(roles ...sdk.AccountObjectIdentifier) *OauthIntegrationForCustomClientsModel {
	t.PreAuthorizedRolesList = tfconfig.SetVariable(
		collections.Map(roles, func(role sdk.AccountObjectIdentifier) tfconfig.Variable {
			return tfconfig.StringVariable(role.Name())
		})...,
	)
	return t
}
