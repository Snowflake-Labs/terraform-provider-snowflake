package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	userDomainDef   = g.NewQueryStruct("UserDomain").Text("Domain", g.KeywordOptions().SingleQuotes().Required())
	emailPatternDef = g.NewQueryStruct("EmailPattern").Text("Pattern", g.KeywordOptions().SingleQuotes().Required())
)

func createSecurityIntegrationOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Create().
		OrReplace().
		SQL("SECURITY INTEGRATION").
		IfNotExists().
		Name()
	qs = apply(qs)
	return qs.
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists")
}

func alterSecurityIntegrationOperation(structName string, apply func(qs *g.QueryStruct) *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(structName).
		Alter().
		SQL("SECURITY INTEGRATION").
		IfExists().
		Name()
	qs = apply(qs)
	return qs.
		NamedList("SET TAG", g.KindOfT[TagAssociation]()).
		NamedList("UNSET TAG", g.KindOfT[ObjectIdentifier]()).
		WithValidation(g.ValidIdentifier, "name")
}

var saml2IntegrationSetDef = g.NewQueryStruct("SAML2IntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_ISSUER", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_SSO_URL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_PROVIDER", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_X509_CERT", g.ParameterOptions().SingleQuotes()).
	ListAssignment("ALLOWED_USER_DOMAINS", "UserDomain", g.ParameterOptions().Parentheses()).
	ListAssignment("ALLOWED_EMAIL_PATTERNS", "EmailPattern", g.ParameterOptions().Parentheses()).
	OptionalTextAssignment("SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_ENABLE_SP_INITIATED", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_X509_CERT", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_SIGN_REQUEST", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_REQUESTED_NAMEID_FORMAT", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_POST_LOGOUT_REDIRECT_URL", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SAML2_FORCE_AUTHN", g.ParameterOptions()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_ISSUER_URL", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("SAML2_SNOWFLAKE_ACS_URL", g.ParameterOptions().SingleQuotes()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider", "Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns",
		"Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl",
		"Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment")

var saml2IntegrationUnsetDef = g.NewQueryStruct("SAML2IntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("SAML2_FORCE_AUTHN").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Saml2ForceAuthn")

var scimIntegrationSetDef = g.NewQueryStruct("SCIMIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword", "Comment")

var scimIntegrationUnsetDef = g.NewQueryStruct("SCIMIntegrationUnset").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("SYNC_PASSWORD").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "NetworkPolicy", "SyncPassword", "Comment")

var SecurityIntegrationsDef = g.NewInterface(
	"SecurityIntegrations",
	"SecurityIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CustomOperation(
		"CreateSAML2",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2",
		createSecurityIntegrationOperation("CreateSAML2Integration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SAML2")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				TextAssignment("SAML2_ISSUER", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("SAML2_SSO_URL", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("SAML2_PROVIDER", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("SAML2_X509_CERT", g.ParameterOptions().Required().SingleQuotes()).
				ListAssignment("ALLOWED_USER_DOMAINS", "UserDomain", g.ParameterOptions().Parentheses()).
				ListAssignment("ALLOWED_EMAIL_PATTERNS", "EmailPattern", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_ENABLE_SP_INITIATED", g.ParameterOptions()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_X509_CERT", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_SIGN_REQUEST", g.ParameterOptions()).
				OptionalTextAssignment("SAML2_REQUESTED_NAMEID_FORMAT", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("SAML2_POST_LOGOUT_REDIRECT_URL", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("SAML2_FORCE_AUTHN", g.ParameterOptions()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_ISSUER_URL", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("SAML2_SNOWFLAKE_ACS_URL", g.ParameterOptions().SingleQuotes())
		}),
		userDomainDef,
		emailPatternDef,
	).
	CustomOperation(
		"CreateSCIM",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim",
		createSecurityIntegrationOperation("CreateSCIMIntegration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SCIM")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				TextAssignment("SCIM_CLIENT", g.ParameterOptions().Required().SingleQuotes()).
				TextAssignment("RUN_AS_ROLE", g.ParameterOptions().Required().SingleQuotes()).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions())
		}),
	).
	CustomOperation(
		"AlterSAML2Integration",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2",
		alterSecurityIntegrationOperation("AlterSAML2Integration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				saml2IntegrationSetDef,
				g.KeywordOptions().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				saml2IntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).OptionalSQL("REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY")
		}),
	).
	CustomOperation(
		"AlterSCIMIntegration",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim",
		alterSecurityIntegrationOperation("AlterSCIMIntegration", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				scimIntegrationSetDef,
				g.KeywordOptions().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				scimIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			)
		}),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-integration",
		g.NewQueryStruct("DropSecurityIntegration").
			Drop().
			SQL("SECURITY INTEGRATION").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-integration",
		g.DbStruct("securityIntegrationDescRow").
			Field("property", "string").
			Field("property_type", "string").
			Field("property_value", "string").
			Field("property_default", "string"),
		g.PlainStruct("SecurityIntegrationProperty").
			Field("Name", "string").
			Field("Type", "string").
			Field("Value", "string").
			Field("Default", "string"),
		g.NewQueryStruct("DescSecurityIntegration").
			Describe().
			SQL("SECURITY INTEGRATION").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-integrations",
		g.DbStruct("securityIntegrationShowRow").
			Text("name").
			Text("type").
			Text("category").
			Bool("enabled").
			OptionalText("comment").
			Time("created_on"),
		g.PlainStruct("SecurityIntegration").
			Text("Name").
			Text("IntegrationType").
			Text("Category").
			Bool("Enabled").
			Text("Comment").
			Time("CreatedOn"),
		g.NewQueryStruct("ShowSecurityIntegrations").
			Show().
			SQL("SECURITY INTEGRATIONS").
			OptionalLike(),
	).
	ShowByIdOperation()
