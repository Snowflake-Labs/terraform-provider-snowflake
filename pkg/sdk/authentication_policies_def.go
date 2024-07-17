package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var AuthenticationMethodsOptionDef = g.NewQueryStruct("AuthenticationMethods").Text("Method", g.KeywordOptions().SingleQuotes())
var MfaAuthenticationMethodsOptionDef = g.NewQueryStruct("MfaAuthenticationMethods").Text("Method", g.KeywordOptions().SingleQuotes())
var ClientTypesOptionDef = g.NewQueryStruct("ClientTypes").Text("ClientType", g.KeywordOptions().SingleQuotes())
var SecurityIntegrationsOptionDef = g.NewQueryStruct("SecurityIntegrationsOption").Text("Name", g.KeywordOptions().SingleQuotes())

var (
	AuthenticationPoliciesDef = g.NewInterface(
		"AuthenticationPolicies",
		"AuthenticationPolicy",
		g.KindOfT[SchemaObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy",
			g.NewQueryStruct("CreateAuthenticationPolicy").
				Create().
				OrReplace().
				SQL("AUTHENTICATION POLICY").
				Name().
				ListAssignment("AUTHENTICATION_METHODS", "AuthenticationMethods", g.ParameterOptions().Parentheses()).
				ListAssignment("MFA_AUTHENTICATION_METHODS", "MfaAuthenticationMethods", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("MFA_ENROLLMENT", g.ParameterOptions()).
				ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
				ListAssignment("SECURITY_INTEGRATIONS", "SecurityIntegrationsOption", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
			AuthenticationMethodsOptionDef,
			MfaAuthenticationMethodsOptionDef,
			ClientTypesOptionDef,
			SecurityIntegrationsOptionDef,
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-authentication-policy",
			g.NewQueryStruct("AlterAuthenticationPolicy").
				Alter().
				SQL("AUTHENTICATION POLICY").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					g.NewQueryStruct("AuthenticationPolicySet").
						ListAssignment("AUTHENTICATION_METHODS", "AuthenticationMethods", g.ParameterOptions().Parentheses()).
						ListAssignment("MFA_AUTHENTICATION_METHODS", "MfaAuthenticationMethods", g.ParameterOptions().Parentheses()).
						OptionalTextAssignment("MFA_ENROLLMENT", g.ParameterOptions().SingleQuotes()).
						ListAssignment("CLIENT_TYPES", "ClientTypes", g.ParameterOptions().Parentheses()).
						ListAssignment("SECURITY_INTEGRATIONS", "SecurityIntegrationsOption", g.ParameterOptions().Parentheses()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					g.NewQueryStruct("AuthenticationPolicyUnset").
						OptionalSQL("CLIENT_TYPES").
						OptionalSQL("AUTHENTICATION_METHODS").
						OptionalSQL("SECURITY_INTEGRATIONS").
						OptionalSQL("MFA_AUTHENTICATION_METHODS").
						OptionalSQL("MFA_ENROLLMENT").
						OptionalSQL("COMMENT").
						WithValidation(g.AtLeastOneValueSet, "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment"),
					g.ListOptions().NoParentheses().SQL("UNSET"),
				).
				Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RenameTo").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy",
			g.NewQueryStruct("DropAuthenticationPolicy").
				Drop().
				SQL("AUTHENTICATION POLICY").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies",
			g.DbStruct("showAuthenticationPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("database_name", "string").
				Field("schema_name", "string").
				Field("owner", "string").
				Field("owner_role_type", "string").
				Field("options", "string"),
			g.PlainStruct("AuthenticationPolicy").
				Field("CreatedOn", "string").
				Field("Name", "string").
				Field("Comment", "string").
				Field("DatabaseName", "string").
				Field("SchemaName", "string").
				Field("Owner", "string").
				Field("OwnerRoleType", "string").
				Field("Options", "string"),
			g.NewQueryStruct("ShowAuthenticationPolicies").
				Show().
				SQL("AUTHENTICATION POLICIES").
				OptionalLike().
				OptionalIn().
				OptionalStartsWith().
				OptionalLimit(),
		).
		ShowByIdOperation().
		DescribeOperation(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy",
			g.DbStruct("describeAuthenticationPolicyDBRow").
				Field("property", "string").
				Field("value", "string"),
			g.PlainStruct("AuthenticationPolicyDescription").
				Field("Property", "string").
				Field("Value", "string"),
			g.NewQueryStruct("DescribeAuthenticationPolicy").
				Describe().
				SQL("AUTHENTICATION POLICY").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
