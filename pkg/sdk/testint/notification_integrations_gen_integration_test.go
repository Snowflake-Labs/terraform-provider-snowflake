package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-TODO]: create topics to perform integration tests
func TestInt_NotificationIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const gcpPubsubSubscriptionName = "TODO"

	assertNotificationIntegration := func(t *testing.T, s *sdk.NotificationIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, true, s.Enabled)
		assert.Equal(t, "EXTERNAL_API", s.NotificationType)
		assert.Equal(t, "NOTIFICATION", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	cleanupNotificationIntegrationProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.NotificationIntegrations.Drop(ctx, sdk.NewDropNotificationIntegrationRequest(id))
			require.NoError(t, err)
		}
	}

	createNotificationIntegrationAutoGoogleRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithAutomatedDataLoadsParams(sdk.NewAutomatedDataLoadsParamsRequest().WithGoogleAutomatedDataLoad(sdk.NewGoogleAutomatedDataLoadRequest(gcpPubsubSubscriptionName)))
	}

	createNotificationIntegrationWithRequest := func(t *testing.T, request *sdk.CreateNotificationIntegrationRequest) *sdk.NotificationIntegration {
		t.Helper()
		id := request.GetName()

		err := client.NotificationIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupNotificationIntegrationProvider(id))

		integration, err := client.NotificationIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	t.Run("create and describe notification integration - auto google", func(t *testing.T) {
		t.Skipf("Skip until we create pub/sub topic (read more in %s; issue: SNOW-TODO)", "https://docs.snowflake.com/en/user-guide/data-load-snowpipe-auto-gcs#creating-the-pub-sub-topic")

		request := createNotificationIntegrationAutoGoogleRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create and describe notification integration - auto azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create and describe notification integration - push amazon", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create and describe notification integration - push google", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create and describe notification integration - push azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create and describe notification integration - email", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: auto", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: push amazon", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: push google", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: push azure", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: email", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter notification integration: set and unset tags", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop notification integration: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop notification integration: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show notification integration: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show notification integration: with options", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe notification integration: non-existing", func(t *testing.T) {
		// TODO: fill me
	})
}
