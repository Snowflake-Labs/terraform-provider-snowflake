package example

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ../main.go

var (
	dbRoleRename = g.QueryStruct("DatabaseRoleRename").
		// Fields
		Identifier("Name", g.KindOfT[DatabaseObjectIdentifier](), nil).
		// Validations
		WithValidation(g.ValidIdentifier, "Name")

	nestedThirdLevel = g.QueryStruct("NestedThirdLevel").
		// Fields
		Identifier("Field", g.KindOfT[DatabaseObjectIdentifier](), nil).
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Field")

	dbRoleSet = g.QueryStruct("DatabaseRoleSet").
		// Fields
		TextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().Required()).
		QueryStructField("NestedThirdLevel", nestedThirdLevel, g.ListOptions().NoParens().SQL("NESTED"))

	dbRoleUnset = g.QueryStruct("DatabaseRoleUnset").
		// Fields
		Text("Comment", g.KeywordOptions().SQL("COMMENT").Required()).
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Comment")

	DatabaseRole = g.NewInterface(
		"DatabaseRoles",
		"DatabaseRole",
		"DatabaseObjectIdentifier",
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
			g.QueryStruct("CreateDatabaseRole").
				// Fields
				Create().
				OrReplace().
				SQL("DATABASE ROLE").
				IfNotExists().
				SelfIdentifier().
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
			g.QueryStruct("AlterDatabaseRole").
				// Fields
				Alter().
				SQL("DATABASE ROLE").
				IfExists().
				SelfIdentifier().
				QueryStructField("Rename", dbRoleRename, g.ListOptions().NoParens().SQL("RENAME TO")).
				QueryStructField("Set", dbRoleSet, g.ListOptions().NoParens().SQL("SET")).
				QueryStructField("Unset", dbRoleUnset, g.ListOptions().NoParens().SQL("UNSET")).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "Rename", "Set", "Unset"),
		)
)
