package client

import (
	"database/sql"
	"os"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClientWithoutInstrumentedSQL checks if the client is initialized with the different driver implementation.
// This is dependent on the SF_TF_NO_INSTRUMENTED_SQL env variable setting. That's why it was extracted to another file.
// To run this test use: `make test-client` command.
func TestNewClientWithoutInstrumentedSQL(t *testing.T) {
	if os.Getenv(snowflakeenvs.NoInstrumentedSql) == "" {
		t.Skipf("Skipping TestNewClientWithoutInstrumentedSQL, because %s is not set", snowflakeenvs.NoInstrumentedSql)
	}

	t.Run("registers snowflake-not-instrumented driver", func(t *testing.T) {
		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)

		assert.NotContains(t, sql.Drivers(), "snowflake-instrumented")
		assert.Contains(t, sql.Drivers(), "snowflake")
	})
}

func TestNewClientWithDebugLoggingSetFromEnv(t *testing.T) {
	t.Run("set gosnowflake driver logging to debug", func(t *testing.T) {
		if os.Getenv(snowflakeenvs.GosnowflakeLogLevel) == "" {
			t.Skipf("Skipping TestNewClientWithDebugLoggingSet, because %s is not set", snowflakeenvs.GosnowflakeLogLevel)
		}

		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)

		assert.Equal(t, "debug", gosnowflake.GetLogger().GetLogLevel())
	})
}
