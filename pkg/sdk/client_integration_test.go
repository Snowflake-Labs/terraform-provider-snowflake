package sdk

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
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

	t.Run("with missing config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config := DefaultConfig()
		_, err = NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("with incorrect config", func(t *testing.T) {
		config, err := ProfileConfig(testprofiles.IncorrectUserAndPassword)
		require.NoError(t, err)
		require.NotNil(t, config)

		_, err = NewClient(config)
		require.ErrorContains(t, err, "Incorrect username or password was specified")
	})

	t.Run("with missing config - should not care about correct env variables", func(t *testing.T) {
		config, err := ProfileConfig(testprofiles.Default)
		require.NoError(t, err)
		require.NotNil(t, config)

		account := config.Account
		t.Setenv(snowflakeenvs.Account, account)

		dir, err := os.UserHomeDir()
		require.NoError(t, err)
		t.Setenv(snowflakeenvs.ConfigPath, dir)

		config = DefaultConfig()
		_, err = NewClient(config)
		require.ErrorContains(t, err, "260000: account is empty")
	})

	t.Run("registers snowflake-instrumented driver", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)

		assert.ElementsMatch(t, sql.Drivers(), []string{"snowflake-instrumented", "snowflake"})
	})
}

func TestClient_ping(t *testing.T) {
	client := defaultTestClient(t)
	err := client.Ping()
	require.NoError(t, err)
}

func TestClient_close(t *testing.T) {
	client := defaultTestClient(t)
	err := client.Close()
	require.NoError(t, err)
}

func TestClient_exec(t *testing.T) {
	client := defaultTestClient(t)
	ctx := context.Background()
	_, err := client.exec(ctx, "SELECT 1")
	require.NoError(t, err)
}

func TestClient_query(t *testing.T) {
	client := defaultTestClient(t)
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
	client := defaultTestClient(t)
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
