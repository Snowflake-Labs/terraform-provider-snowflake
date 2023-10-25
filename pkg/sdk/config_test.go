package sdk

import (
	"os"
	"path/filepath"
	"testing"

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
	cleanupEnvVars := setupEnvVars(t, "", "", "", "", configPath)
	t.Cleanup(cleanupEnvVars)
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
		cleanupEnvVars := setupEnvVars(t, "", "", "", "", configPath)
		t.Cleanup(cleanupEnvVars)
		config, err := ProfileConfig("securityadmin")
		require.NoError(t, err)
		assert.Equal(t, "TEST_ACCOUNT", config.Account)
		assert.Equal(t, "TEST_USER", config.User)
		assert.Equal(t, "abcd1234", config.Password)
		assert.Equal(t, "SECURITYADMIN", config.Role)
	})

	t.Run("with not found profile", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "", "", "", "", configPath)
		t.Cleanup(cleanupEnvVars)
		config, err := ProfileConfig("orgadmin")
		require.NoError(t, err)
		require.Nil(t, config)
	})
}

func TestEnvConfig(t *testing.T) {
	t.Run("with no environment variables", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "", "", "", "", "")
		t.Cleanup(cleanupEnvVars)
		config := EnvConfig()
		assert.Equal(t, "", config.Account)
		assert.Equal(t, "", config.User)
		assert.Equal(t, "", config.Password)
		assert.Equal(t, "", config.Role)
	})

	t.Run("with environment variables", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "TEST_ACCOUNT", "TEST_USER", "abcd1234", "ACCOUNTADMIN", "")
		t.Cleanup(cleanupEnvVars)
		config := EnvConfig()
		assert.Equal(t, "TEST_ACCOUNT", config.Account)
		assert.Equal(t, "TEST_USER", config.User)
		assert.Equal(t, "abcd1234", config.Password)
		assert.Equal(t, "ACCOUNTADMIN", config.Role)
	})
}

func testFile(t *testing.T, filename string, dat []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), filename)
	err := os.WriteFile(path, dat, 0o600)
	require.NoError(t, err)
	return path
}

func setupEnvVars(t *testing.T, account, user, password, role, configPath string) func() {
	t.Helper()
	orginalAccount := os.Getenv("SNOWFLAKE_ACCOUNT")
	orginalUser := os.Getenv("SNOWFLAKE_USER")
	originalPassword := os.Getenv("SNOWFLAKE_PASSWORD")
	originalRole := os.Getenv("SNOWFLAKE_ROLE")
	originalPath := os.Getenv("SNOWFLAKE_CONFIG_PATH")

	os.Setenv("SNOWFLAKE_ACCOUNT", account)
	os.Setenv("SNOWFLAKE_USER", user)
	os.Setenv("SNOWFLAKE_PASSWORD", password)
	os.Setenv("SNOWFLAKE_ROLE", role)
	os.Setenv("SNOWFLAKE_CONFIG_PATH", configPath)

	return func() {
		os.Setenv("SNOWFLAKE_ACCOUNT", orginalAccount)
		os.Setenv("SNOWFLAKE_USER", orginalUser)
		os.Setenv("SNOWFLAKE_PASSWORD", originalPassword)
		os.Setenv("SNOWFLAKE_ROLE", originalRole)
		os.Setenv("SNOWFLAKE_CONFIG_PATH", originalPath)
	}
}
