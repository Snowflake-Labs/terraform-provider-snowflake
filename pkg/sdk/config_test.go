package sdk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
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

func testFile(t *testing.T, filename string, dat []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), filename)
	err := os.WriteFile(path, dat, 0o600)
	require.NoError(t, err)
	return path
}
