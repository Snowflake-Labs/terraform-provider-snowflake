package client

import (
	"database/sql"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// TestNewClientWithoutInstrumentedSQL checks if the client is initialized with the different driver implementation.
// This is dependent on the NO_INSTRUMENTED_SQL env variable setting. That's why it was extracted to another file.
// To run this test use: `make test-client` command.
func TestNewClientWithoutInstrumentedSQL(t *testing.T) {
	if os.Getenv("NO_INSTRUMENTED_SQL") == "" {
		t.Skip("Skipping TestNewClientWithoutInstrumentedSQL, because NO_INSTRUMENTED_SQL is not set")
	}

	t.Run("registers snowflake-not-instrumented driver", func(t *testing.T) {
		config := sdk.DefaultConfig()
		_, err := sdk.NewClient(config)
		require.NoError(t, err)

		assert.NotContains(t, sql.Drivers(), "snowflake-instrumented")
		assert.Contains(t, sql.Drivers(), "snowflake-not-instrumented")
	})
}
