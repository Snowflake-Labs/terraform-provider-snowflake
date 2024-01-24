package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var NotificationIntegrationAllowedRecipientDef = g.NewQueryStruct("NotificationIntegrationAllowedRecipient").
	Text("Email", g.KeywordOptions().SingleQuotes().Required())

// TODO [SNOW-1016561]: all integrations reuse almost the same show, drop, and describe. For now we are copying it. Consider reusing in linked issue.
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
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-notification-integration",
		g.NewQueryStruct("AlterNotificationIntegration").
			Alter().
			SQL("NOTIFICATION INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("NotificationIntegrationSet").
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					OptionalQueryStructField(
						"SetPushParams",
						g.NewQueryStruct("SetPushParams").
							OptionalQueryStructField(
								"SetAmazonPush",
								g.NewQueryStruct("SetAmazonPush").
									TextAssignment("AWS_SNS_TOPIC_ARN", g.ParameterOptions().SingleQuotes().Required()).
									TextAssignment("AWS_SNS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							OptionalQueryStructField(
								"SetGooglePush",
								g.NewQueryStruct("SetGooglePush").
									TextAssignment("GCP_PUBSUB_SUBSCRIPTION_NAME", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							OptionalQueryStructField(
								"SetAzurePush",
								g.NewQueryStruct("SetAzurePush").
									TextAssignment("AZURE_STORAGE_QUEUE_PRIMARY_URI", g.ParameterOptions().SingleQuotes().Required()).
									TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()),
								g.KeywordOptions(),
							).
							WithValidation(g.ExactlyOneValueSet, "SetAmazonPush", "SetGooglePush", "SetAzurePush"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"SetEmailParams",
						g.NewQueryStruct("SetEmailParams").
							ListAssignment("ALLOWED_RECIPIENTS", "NotificationIntegrationAllowedRecipient", g.ParameterOptions().Parentheses().Required()).
							WithValidation(g.ValidateValueSet, "AllowedRecipients"),
						g.KeywordOptions(),
					).
					OptionalComment().
					WithValidation(g.ConflictingFields, "SetPushParams", "SetEmailParams").
					WithValidation(g.AtLeastOneValueSet, "Enabled", "SetPushParams", "SetEmailParams", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("NotificationIntegrationUnset").
					OptionalSQL("ALLOWED_RECIPIENTS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "AllowedRecipients", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropNotificationIntegration").
			Drop().
			SQL("NOTIFICATION INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("showNotificationIntegrationsDbRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("NotificationIntegration").
			Text("Name").
			Text("NotificationType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowNotificationIntegrations").
			Show().
			SQL("NOTIFICATION INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("descNotificationIntegrationsDbRow").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("NotificationIntegrationProperty").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescribeNotificationIntegration").
			Describe().
			SQL("NOTIFICATION INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
