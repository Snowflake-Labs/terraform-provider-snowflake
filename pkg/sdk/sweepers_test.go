package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/require"
)

func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("all sweepers in primary account", func(t *testing.T) {
		client := testClient(t)
		err := SweepAll(client)
		require.NoError(t, err)
	})

	t.Run("all sweepers in secondary account", func(t *testing.T) {
		client := testSecondaryClient(t)
		err := SweepAll(client)
		require.NoError(t, err)
	})
}
