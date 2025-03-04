package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AccountModel) WithAdminUserTypeEnum(adminUserType sdk.UserType) *AccountModel {
	a.AdminUserType = tfconfig.StringVariable(string(adminUserType))
	return a
}
