package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var (
	LimitFromDef = g.QueryStruct("LimitFromApplicationRole").
			Number("Rows", g.KeywordOptions()).
			OptionalText("From", g.KeywordOptions().SingleQuotes())

	ApplicationGrantOptionsDef = g.QueryStruct("ApplicationGrantOptions").
					OptionalIdentifier("ParentRole", g.KindOfTPointer[AccountObjectIdentifier](), nil).
					OptionalIdentifier("ApplicationRole", g.KindOfTPointer[AccountObjectIdentifier](), nil).
					OptionalIdentifier("Application", g.KindOfTPointer[AccountObjectIdentifier](), nil)

	ApplicationRolesDef = g.NewInterface(
		"ApplicationRoles",
		"ApplicationRole",
		g.KindOfT[AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-application-role",
			g.QueryStruct("CreateApplicationRole").
				Create().
				OrReplace().
				SQL("APPLICATION ROLE").
				IfNotExists().
				Name().
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-application-role",
			g.QueryStruct("AlterApplicationRole").
				Alter().
				SQL("APPLICATION ROLE").
				IfExists().
				Name().
				OptionalIdentifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalSQL("UNSET COMMENT"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-application-role",
			g.QueryStruct("DropApplicationRole").
				Drop().
				SQL("APPLICATION ROLE").
				IfExists().
				Name(),
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
				Field("OwnerRoleTYpe", "string"),
			g.QueryStruct("ShowApplicationRoles").
				Show().
				SQL("APPLICATION ROLES IN APPLICATION").
				Name().
				OptionalQueryStructField("LimitFrom", LimitFromDef, nil),
		).
		GrantOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/grant-application-roles",
			g.QueryStruct("GrantApplicationRole").
				Grant().
				SQL("APPLICATION ROLE").
				Name().
				QueryStructField("GrantTo", ApplicationGrantOptionsDef, nil),
		).
		RevokeOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/revoke-application-roles",
			g.QueryStruct("RevokeApplicationRole").
				Revoke().
				SQL("APPLICATION ROLE").
				Name().
				QueryStructField("RevokeFrom", ApplicationGrantOptionsDef, nil),
		)
)
