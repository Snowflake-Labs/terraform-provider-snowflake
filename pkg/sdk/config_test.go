package sdk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFile(t *testing.T) {
	c := `
	[default]
	account='TEST_ACCOUNT'
	user='TEST_USER'
	password='abcd1234'
	role='ACCOUNTADMIN'

	[securityadmin]
	account='TEST_ACCOUNT'
	user='TEST_USER'
	password='abcd1234'
	role='SECURITYADMIN'
	`
	configPath := testFile(t, "config", []byte(c))
	t.Setenv(snowflakeenvs.ConfigPath, configPath)

	m, err := loadConfigFile()
	require.NoError(t, err)
	assert.Equal(t, "TEST_ACCOUNT", m["default"].Account)
	assert.Equal(t, "TEST_USER", m["default"].User)
	assert.Equal(t, "abcd1234", m["default"].Password)
	assert.Equal(t, "ACCOUNTADMIN", m["default"].Role)
	assert.Equal(t, "TEST_ACCOUNT", m["securityadmin"].Account)
	assert.Equal(t, "TEST_USER", m["securityadmin"].User)
	assert.Equal(t, "abcd1234", m["securityadmin"].Password)
	assert.Equal(t, "SECURITYADMIN", m["securityadmin"].Role)
}

func TestProfileConfig(t *testing.T) {
	c := `
	[securityadmin]
	account='TEST_ACCOUNT'
	user='TEST_USER'
	password='abcd1234'
	role='SECURITYADMIN'
	`
	configPath := testFile(t, "config", []byte(c))

	t.Run("with found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("securityadmin")
		require.NoError(t, err)
		assert.Equal(t, "TEST_ACCOUNT", config.Account)
		assert.Equal(t, "TEST_USER", config.User)
		assert.Equal(t, "abcd1234", config.Password)
		assert.Equal(t, "SECURITYADMIN", config.Role)
	})

	t.Run("with not found profile", func(t *testing.T) {
		t.Setenv(snowflakeenvs.ConfigPath, configPath)

		config, err := ProfileConfig("orgadmin")
		require.NoError(t, err)
		require.Nil(t, config)
	})

	t.Run("with not found config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config, err := ProfileConfig("orgadmin")
		require.Error(t, err)
		require.Nil(t, config)
	})
}

func Test_MergeConfig(t *testing.T) {
	createConfig := func(user string, password string, account string, region string) *gosnowflake.Config {
		return &gosnowflake.Config{
			User:     user,
			Password: password,
			Account:  account,
			Region:   region,
		}
	}

	t.Run("merge configs", func(t *testing.T) {
		config1 := createConfig("user", "password", "account", "")
		config2 := createConfig("user2", "", "", "region2")

		config := MergeConfig(config1, config2)

		require.Equal(t, "user", config.User)
		require.Equal(t, "password", config.Password)
		require.Equal(t, "account", config.Account)
		require.Equal(t, "region2", config.Region)
		require.Equal(t, "", config.Role)

		require.Equal(t, config1, config)
		require.Equal(t, "user", config1.User)
		require.Equal(t, "password", config1.Password)
		require.Equal(t, "account", config1.Account)
		require.Equal(t, "region2", config1.Region)
		require.Equal(t, "", config1.Role)
	})

	t.Run("merge configs inverted", func(t *testing.T) {
		config1 := createConfig("user", "password", "account", "")
		config2 := createConfig("user2", "", "", "region2")

		config := MergeConfig(config2, config1)

		require.Equal(t, "user2", config.User)
		require.Equal(t, "password", config.Password)
		require.Equal(t, "account", config.Account)
		require.Equal(t, "region2", config.Region)
		require.Equal(t, "", config.Role)

		require.Equal(t, config2, config)
		require.Equal(t, "user2", config2.User)
		require.Equal(t, "password", config2.Password)
		require.Equal(t, "account", config2.Account)
		require.Equal(t, "region2", config2.Region)
		require.Equal(t, "", config2.Role)
	})
}

func testFile(t *testing.T, filename string, dat []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), filename)
	err := os.WriteFile(path, dat, 0o600)
	require.NoError(t, err)
	return path
}
