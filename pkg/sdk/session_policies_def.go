package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	SessionPoliciesDef = g.NewInterface(
		"SessionPolicies",
		"SessionPolicy",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-session-policy",
			g.QueryStruct("CreateSessionPolicy").
				Create().
				OrReplace().
				SQL("SESSION POLICY").
				IfNotExists().
				Name().
				OptionalIntAssignment("SESSION_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
				OptionalIntAssignment("SESSION_UI_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy",
			g.QueryStruct("AlterSessionPolicy").
				Alter().
				SQL("SESSION POLICY").
				IfExists().
				Name().
				Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				OptionalQueryStructField(
					"Set",
					g.QueryStruct("SessionPolicySet").
						OptionalIntAssignment("SESSION_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
						OptionalIntAssignment("SESSION_UI_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					g.QueryStruct("SessionPolicyUnset").
						OptionalSQL("SESSION_IDLE_TIMEOUT_MINS").
						OptionalSQL("SESSION_UI_IDLE_TIMEOUT_MINS").
						OptionalSQL("COMMENT").
						WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"),
					g.KeywordOptions().SQL("UNSET"),
				).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		)
)
