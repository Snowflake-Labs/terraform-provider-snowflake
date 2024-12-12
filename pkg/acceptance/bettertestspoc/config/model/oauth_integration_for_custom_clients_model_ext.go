package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (t *OauthIntegrationForCustomClientsModel) WithBlockedRolesList(blockedRoles ...string) *OauthIntegrationForCustomClientsModel {
	blockedRolesListStringVariables := make([]tfconfig.Variable, len(blockedRoles))
	for i, v := range blockedRoles {
		blockedRolesListStringVariables[i] = tfconfig.StringVariable(v)
	}

	t.BlockedRolesList = tfconfig.SetVariable(blockedRolesListStringVariables...)
	return t
}
