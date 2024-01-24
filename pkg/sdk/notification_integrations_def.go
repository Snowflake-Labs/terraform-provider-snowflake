package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var NotificationIntegrationAllowedRecipientDef = g.NewQueryStruct("NotificationIntegrationAllowedRecipient").
	Text("Email", g.KeywordOptions().SingleQuotes().Required())

var NotificationIntegrationsDef = g.NewInterface(
	"NotificationIntegrations",
	"NotificationIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration",
		g.NewQueryStruct("CreateNotificationIntegration").
			Create().
			OrReplace().
			SQL("NOTIFICATION INTEGRATION").
			IfNotExists().
			Name().
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalQueryStructField(
				"AutomatedDataLoadsParams",
				g.NewQueryStruct("AutomatedDataLoadsParams").
					PredefinedQueryStructField("notificationType", "string", g.StaticOptions().SQL("TYPE = QUEUE")).
					OptionalQueryStructField(
						"GoogleAutomatedDataLoad",
						g.NewQueryStruct("GoogleAutomatedDataLoad").
							PredefinedQueryStructField("notificationProvider", "string", g.StaticOptions().SQL("NOTIFICATION_PROVIDER = GCP_PUBSUB")).
							TextAssignment("GCP_PUBSUB_SUBSCRIPTION_NAME", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureAutomatedDataLoad",
						g.NewQueryStruct("AzureAutomatedDataLoad").
							PredefinedQueryStructField("notificationProvider", "string", g.StaticOptions().SQL("NOTIFICATION_PROVIDER = AZURE_STORAGE_QUEUE")).
							TextAssignment("AZURE_STORAGE_QUEUE_PRIMARY_URI", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					WithValidation(g.ExactlyOneValueSet, "GoogleAutomatedDataLoad", "AzureAutomatedDataLoad"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"PushNotificationParams",
				g.NewQueryStruct("PushNotificationParams").
					PredefinedQueryStructField("direction", "string", g.StaticOptions().SQL("DIRECTION = OUTBOUND")).
					PredefinedQueryStructField("notificationType", "string", g.StaticOptions().SQL("TYPE = QUEUE")).
					OptionalQueryStructField(
						"AmazonPush",
						g.NewQueryStruct("AmazonPush").
							PredefinedQueryStructField("notificationProvider", "string", g.StaticOptions().SQL("NOTIFICATION_PROVIDER = AWS_SNS")).
							TextAssignment("AWS_SNS_TOPIC_ARN", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AWS_SNS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"GooglePush",
						g.NewQueryStruct("GooglePush").
							PredefinedQueryStructField("notificationProvider", "string", g.StaticOptions().SQL("NOTIFICATION_PROVIDER = GCP_PUBSUB")).
							TextAssignment("GCP_PUBSUB_TOPIC_NAME", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzurePush",
						g.NewQueryStruct("AzurePush").
							PredefinedQueryStructField("notificationProvider", "string", g.StaticOptions().SQL("NOTIFICATION_PROVIDER = AZURE_EVENT_GRID")).
							TextAssignment("AZURE_EVENT_GRID_TOPIC_ENDPOINT", g.ParameterOptions().SingleQuotes().Required()).
							TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
						g.KeywordOptions(),
					).
					WithValidation(g.ExactlyOneValueSet, "AmazonPush", "GooglePush", "AzurePush"),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"EmailParams",
				g.NewQueryStruct("EmailParams").
					PredefinedQueryStructField("notificationType", "string", g.StaticOptions().SQL("TYPE = EMAIL")).
					ListAssignment("ALLOWED_RECIPIENTS", "NotificationIntegrationAllowedRecipient", g.ParameterOptions().Parentheses()),
				g.KeywordOptions(),
			).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "AutomatedDataLoadsParams", "PushNotificationParams", "EmailParams"),
		NotificationIntegrationAllowedRecipientDef,
	)
