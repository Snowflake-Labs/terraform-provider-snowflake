package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type ScimSecurityIntegrationScimClientOption string

var (
	ScimSecurityIntegrationScimClientOkta    ScimSecurityIntegrationScimClientOption = "OKTA"
	ScimSecurityIntegrationScimClientAzure   ScimSecurityIntegrationScimClientOption = "AZURE"
	ScimSecurityIntegrationScimClientGeneric ScimSecurityIntegrationScimClientOption = "GENERIC"
)

type ScimSecurityIntegrationRunAsRoleOption string

var (
	ScimSecurityIntegrationRunAsRoleOktaProvisioner        ScimSecurityIntegrationRunAsRoleOption = "OKTA_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleAadProvisioner         ScimSecurityIntegrationRunAsRoleOption = "AAD_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner ScimSecurityIntegrationRunAsRoleOption = "GENERIC_SCIM_PROVISIONER"
)

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
		Name().
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name")
	qs = apply(qs)
	return qs
}

var saml2IntegrationSetDef = g.NewQueryStruct("Saml2IntegrationSet").
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

var saml2IntegrationUnsetDef = g.NewQueryStruct("Saml2IntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("SAML2_FORCE_AUTHN").
	OptionalSQL("SAML2_REQUESTED_NAMEID_FORMAT").
	OptionalSQL("SAML2_POST_LOGOUT_REDIRECT_URL").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment")

var scimIntegrationSetDef = g.NewQueryStruct("ScimIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword", "Comment")

var scimIntegrationUnsetDef = g.NewQueryStruct("ScimIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("SYNC_PASSWORD").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "NetworkPolicy", "SyncPassword", "Comment")

var SecurityIntegrationsDef = g.NewInterface(
	"SecurityIntegrations",
	"SecurityIntegration",
	g.KindOfT[AccountObjectIdentifier](),
).
	CustomOperation(
		"CreateSaml2",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2",
		createSecurityIntegrationOperation("CreateSaml2", func(qs *g.QueryStruct) *g.QueryStruct {
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
		"CreateScim",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim",
		createSecurityIntegrationOperation("CreateScim", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = SCIM")).
				BooleanAssignment("ENABLED", g.ParameterOptions().Required()).
				Assignment(
					"SCIM_CLIENT",
					g.KindOfT[ScimSecurityIntegrationScimClientOption](),
					g.ParameterOptions().SingleQuotes().Required(),
				).
				Assignment(
					"RUN_AS_ROLE",
					g.KindOfT[ScimSecurityIntegrationRunAsRoleOption](),
					g.ParameterOptions().SingleQuotes().Required(),
				).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalBooleanAssignment("SYNC_PASSWORD", g.ParameterOptions())
		}),
	).
	CustomOperation(
		"AlterSaml2",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2",
		alterSecurityIntegrationOperation("AlterSaml2", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				saml2IntegrationSetDef,
				g.KeywordOptions().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				saml2IntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).OptionalSQL("REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY").
				WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterScim",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-scim",
		alterSecurityIntegrationOperation("AlterScim", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				scimIntegrationSetDef,
				g.KeywordOptions().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				scimIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
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
