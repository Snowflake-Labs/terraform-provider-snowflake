package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var secretsSecurityIntegrationScopeDef = g.NewQueryStruct("SecurityIntegrationScope").Text("Scope", g.KeywordOptions().SingleQuotes().Required())

var SecretsDef = g.NewInterface(
	"Secrets",
	"Secret",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateWithOAuthClientCredentialsFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateOAuthWithClientCredentialsFlow").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Type", "string", g.StaticOptions().SQL("TYPE = OAUTH2")).
		Identifier("SecurityIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION").Required()).
		ListAssignment("OAUTH_SCOPES", "SecurityIntegrationScope", g.ParameterOptions().Parentheses()).
		OptionalComment().
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	secretsSecurityIntegrationScopeDef,
).CustomOperation(
	"CreateWithOAuthAuthorizationCode",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthAuthorizationCode").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Type", "string", g.StaticOptions().SQL("TYPE = OAUTH2")).
		TextAssignment("OAUTH_REFRESH_TOKEN", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		TextAssignment("OAUTH_REFRESH_TOKEN_EXPIRY_TIME", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		Identifier("SecurityIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION")).
		OptionalComment().
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	secretsSecurityIntegrationScopeDef,
).CustomOperation(
	"CreateWithBasicAuthentication",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithBasicAuthentication").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Type", "string", g.StaticOptions().SQL("TYPE = PASSWORD")).
		TextAssignment("USERNAME", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		TextAssignment("PASSWORD", g.ParameterOptions().NoParentheses().SingleQuotes().Required()).
		OptionalComment().
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateWithGenericString",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithGenericString").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Type", "string", g.StaticOptions().SQL("TYPE = GENERIC_STRING")).
		TextAssignment("SECRET_STRING", g.ParameterOptions().SingleQuotes().Required()).
		OptionalComment().
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
)
