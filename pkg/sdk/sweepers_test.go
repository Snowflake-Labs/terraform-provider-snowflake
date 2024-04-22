package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/stretchr/testify/assert"
)

func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("sweep after tests", func(t *testing.T) {
		client := testClient(t)
		secondaryClient := testSecondaryClient(t)

		err := SweepAfterIntegrationTests(client, random.IntegrationTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterIntegrationTests(secondaryClient, random.IntegrationTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(client, random.AcceptanceTestsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(secondaryClient, random.AcceptanceTestsSuffix)
		assert.NoError(t, err)
	})
}
