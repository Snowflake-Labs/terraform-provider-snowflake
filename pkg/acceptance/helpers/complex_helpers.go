package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func (c *TestClient) SetUpLegacyServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, pass string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := c.Ids.RandomAccountObjectIdentifier()
	_, userCleanup := c.User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Type:     sdk.Pointer(sdk.UserTypeLegacyService),
		Password: sdk.String(pass),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := c.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	c.Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, c.Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, c.Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}

// TODO [this PR]: merge with the above function
func (c *TestClient) SetUpServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, publicKey string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := c.Ids.RandomAccountObjectIdentifier()
	_, userCleanup := c.User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Type:         sdk.Pointer(sdk.UserTypeService),
		RSAPublicKey: sdk.String(publicKey),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := c.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	c.Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, c.Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, c.Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}

func (c *TestClient) SetUpTemporaryServiceUser(t *testing.T) *TmpServiceUserConfig {
	accountId := c.Context.CurrentAccountId(t)
	warehouseId := c.Ids.SnowflakeWarehouseId()

	privateKey, publicKey, _ := random.GenerateRSAKeyPair(t)
	tmpUserId, tmpRoleId := c.SetUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	profile := random.AlphaN(6)
	toml := TomlConfigForServiceUser(t, profile, tmpUserId, tmpRoleId, warehouseId, accountId, privateKey)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))

	return &TmpServiceUserConfig{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		TmpUserConfig: TmpUserConfig{
			Profile: profile,
			Path:    configPath,
			UserId:  tmpUserId,
			RoleId:  tmpRoleId,
		},
	}
}

type TmpUserConfig struct {
	Profile string
	Path    string
	UserId  sdk.AccountObjectIdentifier
	RoleId  sdk.AccountObjectIdentifier
}

type TmpServiceUserConfig struct {
	PublicKey  string
	PrivateKey string
	TmpUserConfig
}

type TmpLegacyServiceUserConfig struct {
	Pass string
	TmpUserConfig
}
