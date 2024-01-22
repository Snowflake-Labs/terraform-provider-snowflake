package sdk

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_NewClient(t *testing.T) {
	t.Run("with default config", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)
	})

	t.Run("uses env vars if values are missing", func(t *testing.T) {
		cleanupEnvVars := setupEnvVars(t, "TEST_ACCOUNT", "TEST_USER", "abcd1234", "ACCOUNTADMIN", "")
		t.Cleanup(cleanupEnvVars)
		config := EnvConfig()
		_, err := NewClient(config)
		require.Error(t, err)
	})

	t.Run("registers snowflake-instrumented driver", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)

		assert.ElementsMatch(t, sql.Drivers(), []string{"snowflake-instrumented", "snowflake"})
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

func TestClient_NewClientDriverLoggingLevel(t *testing.T) {
	t.Run("get default gosnowflake driver logging level", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)

		var expected string
		if os.Getenv("GITHUB_ACTIONS") != "" {
			expected = "fatal"
		} else {
			expected = "error"
		}
		assert.Equal(t, expected, gosnowflake.GetLogger().GetLogLevel())
	})

	t.Run("set gosnowflake driver logging level with config", func(t *testing.T) {
		config := DefaultConfig()
		config.Tracing = "trace"
		_, err := NewClient(config)
		require.NoError(t, err)

		assert.Equal(t, "trace", gosnowflake.GetLogger().GetLogLevel())
	})
}
