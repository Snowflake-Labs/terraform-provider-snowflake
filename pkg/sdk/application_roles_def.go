package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var ApplicationRolesDef = g.NewInterface(
	"ApplicationRoles",
	"ApplicationRole",
	g.KindOfT[DatabaseObjectIdentifier](),
).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-application-roles",
		g.DbStruct("applicationRoleDbRow").
			Field("created_on", "time.Time").
			Field("name", "string").
			Field("owner", "string").
			Field("comment", "string").
			Field("owner_role_type", "string"),
		g.PlainStruct("ApplicationRole").
			Field("CreatedOn", "time.Time").
			Field("Name", "string").
			Field("Owner", "string").
			Field("Comment", "string").
			Field("OwnerRoleType", "string"),
		g.QueryStruct("ShowApplicationRoles").
			Show().
			SQL("APPLICATION ROLES IN APPLICATION").
			Identifier("ApplicationName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions()).
			OptionalLimitFrom().
			WithValidation(g.ValidIdentifier, "ApplicationName"),
	)
