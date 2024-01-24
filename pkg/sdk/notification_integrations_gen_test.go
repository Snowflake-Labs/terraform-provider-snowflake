package sdk

import "testing"

// TODO: extract
const (
	gcpPubsubSubscriptionName   = "TODO"
	gcpPubsubTopicName          = "TODO"
	azureStorageQueuePrimaryUri = "TODO"
	azureEventGridTopicEndpoint = "TODO"
	awsSnsTopicArn              = "TODO"
)

func TestNotificationIntegrations_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateNotificationIntegrationOptions for AutomatedDataLoads
	defaultOptsAutomatedDataLoads := func() *CreateNotificationIntegrationOptions {
		return &CreateNotificationIntegrationOptions{
			name:    id,
			Enabled: true,
			AutomatedDataLoadsParams: &AutomatedDataLoadsParams{
				GoogleAutomatedDataLoad: &GoogleAutomatedDataLoad{
					GcpPubsubSubscriptionName: gcpPubsubSubscriptionName,
				},
			},
		}
	}

	// Minimal valid CreateNotificationIntegrationOptions for Push
	defaultOptsPush := func() *CreateNotificationIntegrationOptions {
		return &CreateNotificationIntegrationOptions{
			name:    id,
			Enabled: true,
			PushNotificationParams: &PushNotificationParams{
				AmazonPush: &AmazonPush{
					AwsSnsTopicArn: awsSnsTopicArn,
					AwsSnsRoleArn:  apiAwsRoleArn,
				},
			},
		}
	}

	// Minimal valid CreateNotificationIntegrationOptions for Email
	defaultOptsEmail := func() *CreateNotificationIntegrationOptions {
		return &CreateNotificationIntegrationOptions{
			name:        id,
			Enabled:     true,
			EmailParams: &EmailParams{},
		}
	}

	defaultOpts := defaultOptsEmail

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNotificationIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateNotificationIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from [opts.AutomatedDataLoadsParams opts.PushNotificationParams opts.EmailParams] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.EmailParams = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions", "AutomatedDataLoadsParams", "PushNotificationParams", "EmailParams"))
	})

	t.Run("validation: exactly one field from [opts.AutomatedDataLoadsParams opts.PushNotificationParams opts.EmailParams] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.PushNotificationParams = &PushNotificationParams{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions", "AutomatedDataLoadsParams", "PushNotificationParams", "EmailParams"))
	})

	t.Run("validation: exactly one field from [opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad] should be present", func(t *testing.T) {
		opts := defaultOptsAutomatedDataLoads()
		opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions.AutomatedDataLoadsParams", "GoogleAutomatedDataLoad", "AzureAutomatedDataLoad"))
	})

	t.Run("validation: exactly one field from [opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad] should be present - more present", func(t *testing.T) {
		opts := defaultOptsAutomatedDataLoads()
		opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad = &AzureAutomatedDataLoad{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions.AutomatedDataLoadsParams", "GoogleAutomatedDataLoad", "AzureAutomatedDataLoad"))
	})

	t.Run("validation: exactly one field from [opts.PushNotificationParams.AmazonPush opts.PushNotificationParams.GooglePush opts.PushNotificationParams.AzurePush] should be present", func(t *testing.T) {
		opts := defaultOptsPush()
		opts.PushNotificationParams.AmazonPush = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions.PushNotificationParams", "AmazonPush", "GooglePush", "AzurePush"))
	})

	t.Run("validation: exactly one field from [opts.PushNotificationParams.AmazonPush opts.PushNotificationParams.GooglePush opts.PushNotificationParams.AzurePush] should be present - more present", func(t *testing.T) {
		opts := defaultOptsPush()
		opts.PushNotificationParams.AzurePush = &AzurePush{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateNotificationIntegrationOptions.PushNotificationParams", "AmazonPush", "GooglePush", "AzurePush"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE NOTIFICATION INTEGRATION %s ENABLED = true TYPE = EMAIL`, id.FullyQualifiedName())
	})

	t.Run("all options - auto google", func(t *testing.T) {
		opts := defaultOptsAutomatedDataLoads()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = QUEUE NOTIFICATION_PROVIDER = GCP_PUBSUB GCP_PUBSUB_SUBSCRIPTION_NAME = '%s' COMMENT = 'some comment'", id.FullyQualifiedName(), gcpPubsubSubscriptionName)
	})

	t.Run("all options - auto azure", func(t *testing.T) {
		opts := defaultOptsAutomatedDataLoads()
		opts.AutomatedDataLoadsParams.GoogleAutomatedDataLoad = nil
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		opts.AutomatedDataLoadsParams.AzureAutomatedDataLoad = &AzureAutomatedDataLoad{
			AzureStorageQueuePrimaryUri: azureStorageQueuePrimaryUri,
			AzureTenantId:               azureTenantId,
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = QUEUE NOTIFICATION_PROVIDER = AZURE_STORAGE_QUEUE AZURE_STORAGE_QUEUE_PRIMARY_URI = '%s' AZURE_TENANT_ID = '%s' COMMENT = 'some comment'", id.FullyQualifiedName(), azureStorageQueuePrimaryUri, azureTenantId)
	})

	t.Run("all options - push amazon", func(t *testing.T) {
		opts := defaultOptsPush()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = AWS_SNS AWS_SNS_TOPIC_ARN = '%s' AWS_SNS_ROLE_ARN = '%s' COMMENT = 'some comment'", id.FullyQualifiedName(), awsSnsTopicArn, apiAwsRoleArn)
	})

	t.Run("all options - push google", func(t *testing.T) {
		opts := defaultOptsPush()
		opts.PushNotificationParams.AmazonPush = nil
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		opts.PushNotificationParams.GooglePush = &GooglePush{
			GcpPubsubTopicName: gcpPubsubTopicName,
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = GCP_PUBSUB GCP_PUBSUB_TOPIC_NAME = '%s' COMMENT = 'some comment'", id.FullyQualifiedName(), gcpPubsubTopicName)
	})

	t.Run("all options - push azure", func(t *testing.T) {
		opts := defaultOptsPush()
		opts.PushNotificationParams.AmazonPush = nil
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		opts.PushNotificationParams.AzurePush = &AzurePush{
			AzureEventGridTopicEndpoint: azureEventGridTopicEndpoint,
			AzureTenantId:               azureTenantId,
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true DIRECTION = OUTBOUND TYPE = QUEUE NOTIFICATION_PROVIDER = AZURE_EVENT_GRID AZURE_EVENT_GRID_TOPIC_ENDPOINT = '%s' AZURE_TENANT_ID = '%s' COMMENT = 'some comment'", id.FullyQualifiedName(), azureEventGridTopicEndpoint, azureTenantId)
	})

	t.Run("all options - email", func(t *testing.T) {
		email := "some.email@some.com"
		otherEmail := "some.other.email@some.com"

		opts := defaultOptsEmail()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		opts.EmailParams.AllowedRecipients = []NotificationIntegrationAllowedRecipient{
			{Email: email},
			{Email: otherEmail},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTIFICATION INTEGRATION IF NOT EXISTS %s ENABLED = true TYPE = EMAIL ALLOWED_RECIPIENTS = ('%s', '%s') COMMENT = 'some comment'", id.FullyQualifiedName(), email, otherEmail)
	})
}

func TestNotificationIntegrations_Alter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid AlterNotificationIntegrationOptions
	defaultOpts := func() *AlterNotificationIntegrationOptions {
		return &AlterNotificationIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterNotificationIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNotificationIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: conflicting fields for [opts.Set.SetPushParams opts.Set.SetEmailParams]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("AlterNotificationIntegrationOptions.Set", "SetPushParams", "SetEmailParams"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Enabled opts.Set.SetPushParams opts.Set.SetEmailParams opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNotificationIntegrationOptions.Set", "Enabled", "SetPushParams", "SetEmailParams", "Comment"))
	})

	t.Run("validation: exactly one field from [opts.Set.SetPushParams.SetAmazonPush opts.Set.SetPushParams.SetGooglePush opts.Set.SetPushParams.SetAzurePush] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNotificationIntegrationOptions.Set.SetPushParams", "SetAmazonPush", "SetGooglePush", "SetAzurePush"))
	})

	t.Run("validation: [opts.Set.SetEmailParams.AllowedRecipients] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("AlterNotificationIntegrationOptions.Set.SetEmailParams", "AllowedRecipients"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.AllowedRecipients opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNotificationIntegrationOptions.Unset", "AllowedRecipients", "Comment"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestNotificationIntegrations_Drop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DropNotificationIntegrationOptions
	defaultOpts := func() *DropNotificationIntegrationOptions {
		return &DropNotificationIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNotificationIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP NOTIFICATION INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestNotificationIntegrations_Show(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid ShowNotificationIntegrationOptions
	defaultOpts := func() *ShowNotificationIntegrationOptions {
		return &ShowNotificationIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNotificationIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW NOTIFICATION INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW NOTIFICATION INTEGRATIONS LIKE '%s'", id.Name())
	})
}

func TestNotificationIntegrations_Describe(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DescribeNotificationIntegrationOptions
	defaultOpts := func() *DescribeNotificationIntegrationOptions {
		return &DescribeNotificationIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNotificationIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE NOTIFICATION INTEGRATION %s", id.FullyQualifiedName())
	})
}
