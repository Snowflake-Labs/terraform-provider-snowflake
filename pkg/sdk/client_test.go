package sdk

import (
	"context"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func TestClient_newClient(t *testing.T) {
	config := &gosnowflake.Config{}
	t.Run("uses env vars if values are missing", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "TEST_ACCOUNT", "TEST_USER", "abcd1234", "ACCOUNTADMIN", "")
		t.Cleanup(cleanupEnvVars)
		_, err := NewClient(config)
		require.ErrorIs(t, err, ErrAccountIsEmpty)
	})

	t.Run("with default config", func(t *testing.T) {
		config, err := DefaultConfig()
		require.NoError(t, err)
		_, err = NewClient(config)
		require.NoError(t, err)
	})
}

func TestClient_ping(t *testing.T) {
	client := testClient(t)
	err := client.Ping()
	require.NoError(t, err)
}

func TestClient_close(t *testing.T) {
	client := testClient(t)
	err := client.Close()
	require.NoError(t, err)
}

func TestClient_exec(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	_, err := client.exec(ctx, "SELECT 1")
	require.NoError(t, err)
}

func TestClient_query(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	rows := []struct {
		One int `db:"ONE"`
	}{}
	err := client.query(ctx, &rows, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.NotNil(t, rows)
	require.Equal(t, 1, len(rows))
	require.Equal(t, 1, rows[0].One)
}

func TestClient_queryOne(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	row := struct {
		One int `db:"ONE"`
	}{}
	err := client.queryOne(ctx, &row, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.Equal(t, 1, row.One)
}
