package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (a *AccountModel) WithAdminUserTypeEnum(adminUserType sdk.UserType) *AccountModel {
	a.AdminUserType = tfconfig.StringVariable(string(adminUserType))
	return a
}

func (a *AccountModel) WithAdminRsaPublicKeyMultiline(adminRsaPublicKey string) *AccountModel {
	a.AdminRsaPublicKey = config.MultilineWrapperVariable(adminRsaPublicKey)
	return a
}
