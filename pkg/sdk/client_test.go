package sdk

import (
	"context"
	"os"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func TestClient_newClient(t *testing.T) {
	config := &gosnowflake.Config{}
	t.Run("uses env vars if values are missing", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "TEST_ACCOUNT", "TEST_USER", "abcd1234", "ACCOUNTADMIN")
		t.Cleanup(cleanupEnvVars)
		_, err := NewClient(config)
		require.ErrorIs(t, err, ErrAccountIsEmpty)
	})

	t.Run("with default config", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)
	})
}

func TestClient_ping(t *testing.T) {
	client, err := NewDefaultClient()
	require.NoError(t, err)
	err = client.Ping()
	require.NoError(t, err)
}

func TestClient_defaultConfig(t *testing.T) {
	t.Run("with no environment variables", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "", "", "", "")
		t.Cleanup(cleanupEnvVars)
		config := DefaultConfig()
		require.Equal(t, "", config.Account)
		require.Equal(t, "", config.User)
		require.Equal(t, "", config.Password)
		require.Equal(t, "", config.Role)
	})

	t.Run("with environment variables", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "TEST_ACCOUNT", "TEST_USER", "abcd1234", "ACCOUNTADMIN")
		t.Cleanup(cleanupEnvVars)
		config := DefaultConfig()
		require.Equal(t, "TEST_ACCOUNT", config.Account)
		require.Equal(t, "TEST_USER", config.User)
		require.Equal(t, "abcd1234", config.Password)
		require.Equal(t, "ACCOUNTADMIN", config.Role)
	})
}

func TestClient_exec(t *testing.T) {
	client, err := NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()
	_, err = client.exec(ctx, "SELECT 1")
	require.NoError(t, err)
}

func TestClient_query(t *testing.T) {
	client, err := NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()
	rows := []struct {
		One int `db:"ONE"`
	}{}
	err = client.query(ctx, &rows, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.NotNil(t, rows)
	require.Equal(t, 1, len(rows))
	require.Equal(t, 1, rows[0].One)
}

func TestClient_queryOne(t *testing.T) {
	client, err := NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()
	row := struct {
		One int `db:"ONE"`
	}{}
	err = client.queryOne(ctx, &row, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.Equal(t, 1, row.One)
}

func setupEnvVars(t *testing.T, account, user, password, role string) func() {
	t.Helper()
	orginalAccount := os.Getenv("SNOWFLAKE_ACCOUNT")
	orginalUser := os.Getenv("SNOWFLAKE_USER")
	originalPassword := os.Getenv("SNOWFLAKE_PASSWORD")
	originalRole := os.Getenv("SNOWFLAKE_ROLE")

	os.Setenv("SNOWFLAKE_ACCOUNT", account)
	os.Setenv("SNOWFLAKE_USER", user)
	os.Setenv("SNOWFLAKE_PASSWORD", password)
	os.Setenv("SNOWFLAKE_ROLE", role)

	return func() {
		os.Setenv("SNOWFLAKE_ACCOUNT", orginalAccount)
		os.Setenv("SNOWFLAKE_USER", orginalUser)
		os.Setenv("SNOWFLAKE_PASSWORD", originalPassword)
		os.Setenv("SNOWFLAKE_ROLE", originalRole)
	}
}
