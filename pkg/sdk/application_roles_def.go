package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	LimitFromDef = g.QueryStruct("LimitFromApplicationRole").
			Number("Rows", g.KeywordOptions().Required()).
			OptionalTextAssignment("FROM", g.ParameterOptions().NoEquals().SingleQuotes())

	ApplicationGrantOptionsDef = g.QueryStruct("ApplicationGrantOptions").
					OptionalIdentifier("ParentRole", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
					OptionalIdentifier("ApplicationRole", g.KindOfTPointer[DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION ROLE")).
					OptionalIdentifier("Application", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION")).
					WithValidation(g.ExactlyOneValueSet, "ParentRole", "ApplicationRole", "Application").
					WithValidation(g.ValidIdentifierIfSet, "ParentRole").
					WithValidation(g.ValidIdentifierIfSet, "ApplicationRole").
					WithValidation(g.ValidIdentifierIfSet, "Application")

	ApplicationRolesDef = g.NewInterface(
		"ApplicationRoles",
		"ApplicationRole",
		g.KindOfT[DatabaseObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-application-role",
			g.QueryStruct("CreateApplicationRole").
				Create().
				OrReplace().
				SQL("APPLICATION ROLE").
				IfNotExists().
				Name().
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-application-role",
			g.QueryStruct("AlterApplicationRole").
				Alter().
				SQL("APPLICATION ROLE").
				IfExists().
				Name().
				OptionalIdentifier("RenameTo", g.KindOfTPointer[DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalSQL("UNSET COMMENT").
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetComment", "UnsetComment").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-application-role",
			g.QueryStruct("DropApplicationRole").
				Drop().
				SQL("APPLICATION ROLE").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
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
				OptionalQueryStructField("LimitFrom", LimitFromDef, g.KeywordOptions().SQL("LIMIT")).
				WithValidation(g.ValidIdentifier, "ApplicationName"),
		).
		GrantOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/grant-application-roles",
			g.QueryStruct("GrantApplicationRole").
				Grant().
				SQL("APPLICATION ROLE").
				Name().
				QueryStructField("GrantTo", ApplicationGrantOptionsDef, g.KeywordOptions().SQL("TO").Required()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		RevokeOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/revoke-application-roles",
			g.QueryStruct("RevokeApplicationRole").
				Revoke().
				SQL("APPLICATION ROLE").
				Name().
				QueryStructField("RevokeFrom", ApplicationGrantOptionsDef, g.KeywordOptions().SQL("FROM").Required()).
				WithValidation(g.ValidIdentifier, "name"),
		)
)
