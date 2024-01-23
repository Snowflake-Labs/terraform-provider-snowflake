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
					PredefinedQueryStructField("apiProvider", "string", g.StaticOptions().SQL("STORAGE_PROVIDER = google_api_gateway")).
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
	)
