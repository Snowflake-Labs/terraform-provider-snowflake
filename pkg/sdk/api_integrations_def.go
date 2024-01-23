package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type ApiIntegrationAwsApiProviderType string

var (
	ApiIntegrationAwsApiGateway           ApiIntegrationAwsApiProviderType = "aws_api_gateway"
	ApiIntegrationAwsPrivateApiGateway    ApiIntegrationAwsApiProviderType = "aws_private_api_gateway"
	ApiIntegrationAwsGovApiGateway        ApiIntegrationAwsApiProviderType = "aws_gov_api_gateway"
	ApiIntegrationAwsGovPrivateApiGateway ApiIntegrationAwsApiProviderType = "aws_gov_private_api_gateway"
)

var ApiIntegrationEndpointPrefixDef = g.NewQueryStruct("ApiIntegrationEndpointPrefix").Text("Path", g.KeywordOptions().SingleQuotes().Required())

var ApiIntegrationsDef = g.NewInterface(
	"ApiIntegrations",
	"ApiIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-api-integration",
		g.NewQueryStruct("CreateApiIntegration").
			Create().
			OrReplace().
			SQL("API INTEGRATION").
			IfNotExists().
			Name().
			OptionalQueryStructField(
				"S3ApiProviderParams",
				g.NewQueryStruct("S3ApiParams").
					Assignment("API_PROVIDER", g.KindOfT[ApiIntegrationAwsApiProviderType](), g.ParameterOptions().NoQuotes().Required()).
					TextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"AzureApiProviderParams",
				g.NewQueryStruct("AzureApiParams").
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("API_PROVIDER = azure_api_management")).
					TextAssignment("AZURE_TENANT_ID", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes().Required()).
					OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()),
				g.KeywordOptions(),
			).
			OptionalQueryStructField(
				"GCSApiProviderParams",
				g.NewQueryStruct("GCSApiParams").
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("API_PROVIDER = google_api_gateway")).
					TextAssignment("GOOGLE_AUDIENCE", g.ParameterOptions().SingleQuotes().Required()),
				g.KeywordOptions(),
			).
			ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses().Required()).
			ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
			BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithValidation(g.ExactlyOneValueSet, "S3ApiProviderParams", "AzureApiProviderParams", "GCSApiProviderParams"),
		ApiIntegrationEndpointPrefixDef,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-api-integration",
		g.NewQueryStruct("AlterApiIntegration").
			Alter().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("ApiIntegrationSet").
					OptionalQueryStructField(
						"S3Params",
						g.NewQueryStruct("SetS3ApiParams").
							OptionalTextAssignment("API_AWS_ROLE_ARN", g.ParameterOptions().SingleQuotes()).
							OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()).
							WithValidation(g.AtLeastOneValueSet, "ApiAwsRoleArn", "ApiKey"),
						g.KeywordOptions(),
					).
					OptionalQueryStructField(
						"AzureParams",
						g.NewQueryStruct("SetAzureApiParams").
							TextAssignment("AZURE_AD_APPLICATION_ID", g.ParameterOptions().SingleQuotes().Required()).
							OptionalTextAssignment("API_KEY", g.ParameterOptions().SingleQuotes()).
							WithValidation(g.AtLeastOneValueSet, "AzureAdApplicationId", "ApiKey"),
						g.KeywordOptions(),
					).
					OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
					ListAssignment("API_ALLOWED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					ListAssignment("API_BLOCKED_PREFIXES", "ApiIntegrationEndpointPrefix", g.ParameterOptions().Parentheses()).
					OptionalComment().
					WithValidation(g.ConflictingFields, "S3Params", "AzureParams").
					WithValidation(g.AtLeastOneValueSet, "S3Params", "AzureParams", "Enabled", "ApiAllowedPrefixes", "ApiBlockedPrefixes", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ApiIntegrationUnset").
					OptionalSQL("API_KEY").
					OptionalSQL("ENABLED").
					OptionalSQL("API_BLOCKED_PREFIXES").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ApiKey", "Enabled", "ApiBlockedPrefixes", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfExists", "SetTags").
			WithValidation(g.ConflictingFields, "IfExists", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropApiIntegration").
			Drop().
			SQL("API INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("showApiIntegrationsDbRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("ApiIntegration").
			Text("Name").
			Text("ApiType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowApiIntegrations").
			Show().
			SQL("API INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperation().
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("descApiIntegrationsDbRow").
			Text("property").
			Text("property_type").
			Text("property_value").
			Text("property_default"),
		g.PlainStruct("ApiIntegrationProperty").
			Text("Name").
			Text("Type").
			Text("Value").
			Text("Default"),
		g.NewQueryStruct("DescribeApiIntegration").
			Describe().
			SQL("API INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
