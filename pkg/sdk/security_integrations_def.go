package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

type OauthSecurityIntegrationUseSecondaryRolesOption string

const (
	OauthSecurityIntegrationUseSecondaryRolesImplicit OauthSecurityIntegrationUseSecondaryRolesOption = "IMPLICIT"
	OauthSecurityIntegrationUseSecondaryRolesNone     OauthSecurityIntegrationUseSecondaryRolesOption = "NONE"
)

type OauthSecurityIntegrationClientTypeOption string

const (
	OauthSecurityIntegrationClientTypePublic       OauthSecurityIntegrationClientTypeOption = "PUBLIC"
	OauthSecurityIntegrationClientTypeConfidential OauthSecurityIntegrationClientTypeOption = "CONFIDENTIAL"
)

type OauthSecurityIntegrationClientOption string

const (
	OauthSecurityIntegrationClientLooker         OauthSecurityIntegrationClientOption = "LOOKER"
	OauthSecurityIntegrationClientTableauDesktop OauthSecurityIntegrationClientOption = "TABLEAU_DESKTOP"
	OauthSecurityIntegrationClientTableauServer  OauthSecurityIntegrationClientOption = "TABLEAU_SERVER"
)

type ScimSecurityIntegrationScimClientOption string

const (
	ScimSecurityIntegrationScimClientOkta    ScimSecurityIntegrationScimClientOption = "OKTA"
	ScimSecurityIntegrationScimClientAzure   ScimSecurityIntegrationScimClientOption = "AZURE"
	ScimSecurityIntegrationScimClientGeneric ScimSecurityIntegrationScimClientOption = "GENERIC"
)

type ScimSecurityIntegrationRunAsRoleOption string

const (
	ScimSecurityIntegrationRunAsRoleOktaProvisioner        ScimSecurityIntegrationRunAsRoleOption = "OKTA_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleAadProvisioner         ScimSecurityIntegrationRunAsRoleOption = "AAD_PROVISIONER"
	ScimSecurityIntegrationRunAsRoleGenericScimProvisioner ScimSecurityIntegrationRunAsRoleOption = "GENERIC_SCIM_PROVISIONER"
)

var (
	userDomainDef             = g.NewQueryStruct("UserDomain").Text("Domain", g.KeywordOptions().SingleQuotes().Required())
	emailPatternDef           = g.NewQueryStruct("EmailPattern").Text("Pattern", g.KeywordOptions().SingleQuotes().Required())
	preAuthorizedRolesListDef = g.NewQueryStruct("PreAuthorizedRolesList").
					List("PreAuthorizedRolesList", "AccountObjectIdentifier", g.ListOptions().MustParentheses())
	blockedRolesListDef = g.NewQueryStruct("BlockedRolesList").
				List("BlockedRolesList", "AccountObjectIdentifier", g.ListOptions().MustParentheses())
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

var snowflakeOauthPartnerIntegrationSetDef = g.NewQueryStruct("OauthPartnerIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalAssignment(
		"OAUTH_USE_SECONDARY_ROLES",
		g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
		g.ParameterOptions(),
	).
	OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthRedirectUri", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles",
		"BlockedRolesList", "Comment")

var snowflakeOauthPartnerIntegrationUnsetDef = g.NewQueryStruct("OauthPartnerIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("OAUTH_USE_SECONDARY_ROLES").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthUseSecondaryRoles")

var snowflakeOauthCustomIntegrationSetDef = g.NewQueryStruct("OauthCustomIntegrationSet").
	OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
	OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("OAUTH_ALLOW_NON_TLS_REDIRECT_URI", g.ParameterOptions()).
	OptionalBooleanAssignment("OAUTH_ENFORCE_PKCE", g.ParameterOptions()).
	OptionalAssignment(
		"OAUTH_USE_SECONDARY_ROLES",
		g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
		g.ParameterOptions(),
	).
	OptionalQueryStructField("PreAuthorizedRolesList", preAuthorizedRolesListDef, g.ParameterOptions().SQL("PRE_AUTHORIZED_ROLES_LIST").Parentheses()).
	OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
	OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
	OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
	OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
	OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes()).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthRedirectUri", "OauthAllowNonTlsRedirectUri", "OauthEnforcePkce", "PreAuthorizedRolesList",
		"BlockedRolesList", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "NetworkPolicy", "OauthClientRsaPublicKey",
		"OauthClientRsaPublicKey2", "Comment")

var snowflakeOauthCustomIntegrationUnsetDef = g.NewQueryStruct("OauthCustomIntegrationUnset").
	OptionalSQL("ENABLED").
	OptionalSQL("OAUTH_USE_SECONDARY_ROLES").
	OptionalSQL("NETWORK_POLICY").
	OptionalSQL("OAUTH_CLIENT_RSA_PUBLIC_KEY").
	OptionalSQL("OAUTH_CLIENT_RSA_PUBLIC_KEY_2").
	WithValidation(g.AtLeastOneValueSet, "Enabled", "OauthUseSecondaryRoles", "NetworkPolicy", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2")

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
	OptionalSQL("SAML2_FORCE_AUTHN").
	OptionalSQL("SAML2_REQUESTED_NAMEID_FORMAT").
	OptionalSQL("SAML2_POST_LOGOUT_REDIRECT_URL").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet, "Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment")

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
		"CreateOauthPartner",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake",
		createSecurityIntegrationOperation("CreateOauthPartner", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = OAUTH")).
				Assignment(
					"OAUTH_CLIENT",
					g.KindOfT[OauthSecurityIntegrationClientOption](),
					g.ParameterOptions().Required(),
				).
				OptionalTextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalAssignment(
					"OAUTH_USE_SECONDARY_ROLES",
					g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
					g.ParameterOptions(),
				).
				OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses())
		}),
		preAuthorizedRolesListDef,
		blockedRolesListDef,
	).
	CustomOperation(
		"CreateOauthCustom",
		"https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake",
		createSecurityIntegrationOperation("CreateOauthCustom", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.
				PredefinedQueryStructField("integrationType", "string", g.StaticOptions().SQL("TYPE = OAUTH")).
				PredefinedQueryStructField("oauthClient", "string", g.StaticOptions().SQL("OAUTH_CLIENT = CUSTOM")).
				Assignment(
					"OAUTH_CLIENT_TYPE",
					g.KindOfT[OauthSecurityIntegrationClientTypeOption](),
					g.ParameterOptions().Required().SingleQuotes(),
				).
				TextAssignment("OAUTH_REDIRECT_URI", g.ParameterOptions().Required().SingleQuotes()).
				OptionalBooleanAssignment("ENABLED", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ALLOW_NON_TLS_REDIRECT_URI", g.ParameterOptions()).
				OptionalBooleanAssignment("OAUTH_ENFORCE_PKCE", g.ParameterOptions()).
				OptionalAssignment(
					"OAUTH_USE_SECONDARY_ROLES",
					g.KindOfT[OauthSecurityIntegrationUseSecondaryRolesOption](),
					g.ParameterOptions(),
				).
				OptionalQueryStructField("PreAuthorizedRolesList", preAuthorizedRolesListDef, g.ParameterOptions().SQL("PRE_AUTHORIZED_ROLES_LIST").Parentheses()).
				OptionalQueryStructField("BlockedRolesList", blockedRolesListDef, g.ParameterOptions().SQL("BLOCKED_ROLES_LIST").Parentheses()).
				OptionalBooleanAssignment("OAUTH_ISSUE_REFRESH_TOKENS", g.ParameterOptions()).
				OptionalNumberAssignment("OAUTH_REFRESH_TOKEN_VALIDITY", g.ParameterOptions()).
				OptionalIdentifier("NetworkPolicy", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("OAUTH_CLIENT_RSA_PUBLIC_KEY_2", g.ParameterOptions().SingleQuotes())
		}),
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
		"AlterOauthPartner",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake",
		alterSecurityIntegrationOperation("AlterOauthPartner", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				snowflakeOauthPartnerIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				snowflakeOauthPartnerIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterOauthCustom",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake",
		alterSecurityIntegrationOperation("AlterOauthCustom", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				snowflakeOauthCustomIntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
			).OptionalQueryStructField(
				"Unset",
				snowflakeOauthCustomIntegrationUnsetDef,
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags")
		}),
	).
	CustomOperation(
		"AlterSaml2",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2",
		alterSecurityIntegrationOperation("AlterSaml2", func(qs *g.QueryStruct) *g.QueryStruct {
			return qs.OptionalQueryStructField(
				"Set",
				saml2IntegrationSetDef,
				g.ListOptions().NoParentheses().SQL("SET"),
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
				g.ListOptions().NoParentheses().SQL("SET"),
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
