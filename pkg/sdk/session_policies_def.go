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
		)
)
