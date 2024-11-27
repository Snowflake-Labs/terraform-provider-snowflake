package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

// TODO [SNOW-1827324]: add TestClient ref to each specific client, so that we enhance specific client and not the base one

func (c *TestClient) TempTomlConfigForServiceUser(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUser(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.PrivateKey)
	})
}

func (c *TestClient) TempTomlConfigForServiceUserWithEncryptedKey(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUserWithEncryptedKey(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.EncryptedPrivateKey, serviceUser.Pass)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForServiceUser(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlIncorrectConfigForServiceUser(t, profile, serviceUser.AccountId)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForServiceUserWithEncryptedKey(t *testing.T, serviceUser *TmpServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForServiceUserWithEncryptedKey(t, profile, serviceUser.UserId, serviceUser.RoleId, serviceUser.WarehouseId, serviceUser.AccountId, serviceUser.EncryptedPrivateKey, "incorrect pass")
	})
}

func (c *TestClient) TempTomlConfigForLegacyServiceUser(t *testing.T, legacyServiceUser *TmpLegacyServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForLegacyServiceUser(t, profile, legacyServiceUser.UserId, legacyServiceUser.RoleId, legacyServiceUser.WarehouseId, legacyServiceUser.AccountId, legacyServiceUser.Pass)
	})
}

func (c *TestClient) TempIncorrectTomlConfigForLegacyServiceUser(t *testing.T, legacyServiceUser *TmpLegacyServiceUser) *TmpTomlConfig {
	t.Helper()
	return c.StoreTempTomlConfig(t, func(profile string) string {
		return TomlConfigForLegacyServiceUser(t, profile, legacyServiceUser.UserId, legacyServiceUser.RoleId, legacyServiceUser.WarehouseId, legacyServiceUser.AccountId, "incorrect pass")
	})
}

func (c *TestClient) StoreTempTomlConfig(t *testing.T, tomlProvider func(string) string) *TmpTomlConfig {
	t.Helper()

	profile := random.AlphaN(6)
	toml := tomlProvider(profile)
	configPath := testhelpers.TestFile(t, random.AlphaN(10), []byte(toml))
	return &TmpTomlConfig{
		Profile: profile,
		Path:    configPath,
	}
}

type TmpTomlConfig struct {
	Profile string
	Path    string
}
