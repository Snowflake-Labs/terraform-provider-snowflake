package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-TODO]: create topics to perform integration tests
// auto google: https://docs.snowflake.com/en/user-guide/data-load-snowpipe-auto-gcs#creating-the-pub-sub-topic
// auto azure: https://docs.snowflake.com/en/user-guide/data-load-snowpipe-auto-azure#create-a-storage-queue
// push amazon: https://docs.snowflake.com/en/user-guide/data-load-snowpipe-errors-sns#step-1-creating-an-amazon-sns-topic
func TestInt_NotificationIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const gcpPubsubSubscriptionName = "projects/project-1234/subscriptions/sub2"
	const gcpPubsubTopicName = "projects/project-1234/topics/top2"
	const azureStorageQueuePrimaryUri = "azure://great-bucket/great-path/"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureEventGridTopicEndpoint = "https://apim-hello-world.azure-api.net/dev"
	const awsSnsTopicArn = "arn:aws:sns:us-east-2:123456789012:MyTopic"
	const awsSnsRoleArn = "arn:aws:iam::000000000001:/role/test"

	assertNotificationIntegration := func(t *testing.T, s *sdk.NotificationIntegration, name sdk.AccountObjectIdentifier, notificationType string, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, true, s.Enabled)
		assert.Equal(t, notificationType, s.NotificationType)
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

	createNotificationIntegrationAutoAzureRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithAutomatedDataLoadsParams(sdk.NewAutomatedDataLoadsParamsRequest().WithAzureAutomatedDataLoad(sdk.NewAzureAutomatedDataLoadRequest(azureStorageQueuePrimaryUri, azureTenantId)))
	}

	createNotificationIntegrationPushAmazonRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithAmazonPush(sdk.NewAmazonPushRequest(awsSnsTopicArn, awsSnsRoleArn)))
	}

	createNotificationIntegrationPushGoogleRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithGooglePush(sdk.NewGooglePushRequest(gcpPubsubTopicName)))
	}

	createNotificationIntegrationPushAzureRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithPushNotificationParams(sdk.NewPushNotificationParamsRequest().WithAzurePush(sdk.NewAzurePushRequest(azureEventGridTopicEndpoint, azureTenantId)))
	}

	createNotificationIntegrationEmailRequest := func(t *testing.T) *sdk.CreateNotificationIntegrationRequest {
		t.Helper()
		id := sdk.RandomAccountObjectIdentifier()

		// TODO [SNOW-1007539]: use email of our service user
		return sdk.NewCreateNotificationIntegrationRequest(id, true).
			WithEmailParams(sdk.NewEmailParamsRequest().WithAllowedRecipients([]sdk.NotificationIntegrationAllowedRecipient{{Email: "artur.sawicki@snowflake.com"}}))
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
		request := createNotificationIntegrationAutoGoogleRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - GCP_PUBSUB", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "INBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "GCP_PUBSUB_SUBSCRIPTION_NAME", Type: "String", Value: gcpPubsubSubscriptionName, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("create and describe notification integration - auto azure", func(t *testing.T) {
		request := createNotificationIntegrationAutoAzureRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AZURE_STORAGE_QUEUE", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_STORAGE_QUEUE_PRIMARY_URI", Type: "String", Value: azureStorageQueuePrimaryUri, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("create and describe notification integration - push amazon", func(t *testing.T) {
		request := createNotificationIntegrationPushAmazonRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AWS_SNS", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "AWS_SNS", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_TOPIC_ARN", Type: "String", Value: awsSnsTopicArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AWS_SNS_ROLE_ARN", Type: "String", Value: awsSnsRoleArn, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	// TODO [SNOW-]: check the error 001422 (22023): SQL compilation error: invalid value 'OUTBOUND' for property 'Direction'
	t.Run("create and describe notification integration - push google", func(t *testing.T) {
		t.Skip("Skipping because of the error: 001422 (22023): SQL compilation error: invalid value 'OUTBOUND' for property 'Direction'")
		request := createNotificationIntegrationPushGoogleRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - GCP_PUBSUB", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "GCP_PUBSUB", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "GCP_PUBSUB_TOPIC_NAME", Type: "String", Value: gcpPubsubTopicName, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	// TODO [SNOW-]: check the error 001008 (22023): SQL compilation error: invalid value [QUEUE - AZURE_EVENT_GRID] for parameter 'Integration Type'
	t.Run("create and describe notification integration - push azure", func(t *testing.T) {
		t.Skip("Skipping because of the error: 001008 (22023): SQL compilation error: invalid value [QUEUE - AZURE_EVENT_GRID] for parameter 'Integration Type'")
		request := createNotificationIntegrationPushAzureRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "QUEUE - AZURE_EVENT_GRID", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "NOTIFICATION_PROVIDER", Type: "String", Value: "GCP_PUBSUB", Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "DIRECTION", Type: "String", Value: "OUTBOUND", Default: "INBOUND"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_EVENT_GRID_TOPIC_ENDPOINT", Type: "String", Value: azureEventGridTopicEndpoint, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "AZURE_TENANT_ID", Type: "String", Value: azureTenantId, Default: ""})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("create and describe notification integration - email", func(t *testing.T) {
		request := createNotificationIntegrationEmailRequest(t)

		integration := createNotificationIntegrationWithRequest(t, request)

		assertNotificationIntegration(t, integration, request.GetName(), "EMAIL", "")

		details, err := client.NotificationIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "ALLOWED_RECIPIENTS", Type: "List", Value: "artur.sawicki@snowflake.com", Default: "[]"})
		assert.Contains(t, details, sdk.NotificationIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
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
