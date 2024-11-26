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

func (c *TestClient) SetUpTemporaryLegacyServiceUser(t *testing.T) *TmpLegacyServiceUser {
	warehouseId := c.Ids.SnowflakeWarehouseId()
	accountId := c.Context.CurrentAccountId(t)

	pass := random.Password()
	tmpUserId, tmpRoleId := c.SetUpLegacyServiceUserWithAccessToTestDatabaseAndWarehouse(t, pass)

	return &TmpLegacyServiceUser{
		Pass: pass,
		TmpUser: TmpUser{
			UserId:      tmpUserId,
			RoleId:      tmpRoleId,
			WarehouseId: warehouseId,
			AccountId:   accountId,
		},
	}
}

func (c *TestClient) SetUpTemporaryServiceUser(t *testing.T) *TmpServiceUser {
	warehouseId := c.Ids.SnowflakeWarehouseId()
	accountId := c.Context.CurrentAccountId(t)

	pass := random.Password()
	privateKey, encryptedKey, publicKey, _ := random.GenerateRSAKeyPair(t, pass)
	tmpUserId, tmpRoleId := c.SetUpServiceUserWithAccessToTestDatabaseAndWarehouse(t, publicKey)

	return &TmpServiceUser{
		PublicKey:           publicKey,
		PrivateKey:          privateKey,
		EncryptedPrivateKey: encryptedKey,
		Pass:                pass,
		TmpUser: TmpUser{
			UserId:      tmpUserId,
			RoleId:      tmpRoleId,
			WarehouseId: warehouseId,
			AccountId:   accountId,
		},
	}
}

func (c *TestClient) TempTomlConfigForServiceUser(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUser(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.PrivateKey)
	})
}

func (c *TestClient) TempTomlConfigForServiceUserWithEncryptedKey(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUserWithEncryptedKey(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.EncryptedPrivateKey, serviceUser.Pass)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForServiceUser(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlIncorrectConfigForServiceUser(t, profile, serviceUser.AccountId)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForServiceUserWithEncryptedKey(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUserWithEncryptedKey(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.EncryptedPrivateKey, "incorrect pass")
	})
}

func (c *TestClient) TempTomlConfigForLegacyServiceUser(t *testing.T, legacyServiceUser *TmpLegacyServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForLegacyServiceUser(t, profile, legacyServiceUser.UserId, legacyServiceUser.RoleId, legacyServiceUser.WarehouseId, legacyServiceUser.AccountId, legacyServiceUser.Pass)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForLegacyServiceUser(t *testing.T, legacyServiceUser *TmpLegacyServiceUser) *TmpTomlConfig {
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForLegacyServiceUser(t, profile, legacyServiceUser.UserId, legacyServiceUser.RoleId, legacyServiceUser.WarehouseId, legacyServiceUser.AccountId, "incorrect pass")
	})
}

func (c *TestClient) StoreTempTomlConfig(t *testing.T, tomlProvider func(string) string) *TmpTomlConfig {
	profile := random.AlphaN(6)
	toml := tomlProvider(profile)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))
	return &TmpTomlConfig{
		Profile: profile,
		Path:    configPath,
	}
}

type TmpUser struct {
	UserId      sdk.AccountObjectIdentifier
	RoleId      sdk.AccountObjectIdentifier
	WarehouseId sdk.AccountObjectIdentifier
	AccountId   sdk.AccountIdentifier
}

type TmpServiceUser struct {
	PublicKey           string
	PrivateKey          string
	EncryptedPrivateKey string
	Pass                string
	TmpUser
}

type TmpLegacyServiceUser struct {
	Pass string
	TmpUser
}

type TmpTomlConfig struct {
	Profile string
	Path    string
}
