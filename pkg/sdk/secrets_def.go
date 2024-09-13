package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var secretsSecurityIntegrationScopeDef = g.NewQueryStruct("SecurityIntegrationScope").
	Text("Scope", g.KeywordOptions().SingleQuotes().Required())

var secretsOAuthScopes = g.NewQueryStruct("OAuthScopes").
	List("OAuthScopes", "SecurityIntegrationScope", g.ListOptions().MustParentheses())

var secretSet = g.NewQueryStruct("SecretSet").
	OptionalComment().
	OptionalQueryStructField("OAuthScopes", secretsOAuthScopes, g.ParameterOptions().MustParentheses().SQL("OAUTH_SCOPES")).
	OptionalTextAssignment("OAUTH_REFRESH_TOKEN", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_REFRESH_TOKEN_EXPIRY_TIME", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("USERNAME", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("PASSWORD", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SECRET_STRING", g.ParameterOptions().SingleQuotes())

// unset doest work, need to use "set comment = null"
// OptionalSQL("SET COMMENT = NULL")
var secretUnset = g.NewQueryStruct("SecretUnset").
	OptionalSQL("UNSET COMMENT")

var SecretsDef = g.NewInterface(
	"Secrets",
	"Secret",
	g.KindOfT[SchemaObjectIdentifier](),
).CustomOperation(
	"CreateWithOAuthClientCredentialsFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthClientCredentialsFlow").
		Create().
		OrReplace().
		SQL("SECRET").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Type", "string", g.StaticOptions().SQL("TYPE = OAUTH2")).
		Identifier("SecurityIntegration", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required().Equals().SQL("API_AUTHENTICATION").Required()).
		ListAssignment("OAUTH_SCOPES", "SecurityIntegrationScope", g.ParameterOptions().Parentheses()).
		//QueryStructField("OAuthScopes", secretsOAuthScopes, g.ParameterOptions().MustParentheses()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	secretsSecurityIntegrationScopeDef,
).CustomOperation(
	"CreateWithOAuthAuthorizationCodeFlow",
	"https://docs.snowflake.com/en/sql-reference/sql/create-secret",
	g.NewQueryStruct("CreateWithOAuthAuthorizationCodeFlow").
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
		WithValidation(g.ValidIdentifier, "name").
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
		WithValidation(g.ValidIdentifier, "name").
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
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-secret",
	g.NewQueryStruct("AlterSecret").
		Alter().
		SQL("SECRET").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			secretSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			secretUnset,
			g.KeywordOptions().SQL("UNSET"),
		).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
)
