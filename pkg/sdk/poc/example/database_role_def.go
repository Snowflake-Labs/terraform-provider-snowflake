package example

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ../main.go

var (
	dbRoleRename = g.QueryStruct("DatabaseRoleRename").
		// Fields
		Identifier("Name", g.KindOfT[DatabaseObjectIdentifier]()).
		// Validations
		WithValidation(g.ValidIdentifier, "Name")

	nestedThirdLevel = g.QueryStruct("NestedThirdLevel").
		// Fields
		Identifier("Field", g.KindOfT[DatabaseObjectIdentifier]()).
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Field")

	dbRoleSet = g.QueryStruct("DatabaseRoleSet").
		// Fields
		TextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()). // TODO Required
		QueryStructField(nestedThirdLevel, "NestedThirdLevel", g.KindOfPointer(nestedThirdLevel.Name), g.ListOptions().NoParens().SQL("NESTED"))

	dbRoleUnset = g.QueryStruct("DatabaseRoleUnset").
		// Fields
		Text("Comment", g.KeywordOptions().SQL("COMMENT")). // TODO Required
		QueryStructField(nestedThirdLevel, "NestedThirdLevel", g.KindOfPointer(nestedThirdLevel.Name), g.ListOptions().NoParens().SQL("NESTED")).
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Comment")
)

var DatabaseRole = g.NewInterface(
	"DatabaseRoles",
	"DatabaseRole",
	"DatabaseObjectIdentifier",
	g.NewOperation("Create", "https://docs.snowflake.com/en/sql-reference/sql/create-database-role").
		WithOptionsStruct(
			g.QueryStruct("CreateDatabaseRole").
				// Fields
				Create().
				OrReplace().
				SQL("DATABASE ROLE").
				IfNotExists().
				Identifier("name", g.KindOfT[DatabaseObjectIdentifier]()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
		),
	g.NewOperation("Alter", "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role").
		WithOptionsStruct(
			g.QueryStruct("AlterDatabaseRole").
				// Fields
				Alter().
				SQL("DATABASE ROLE").
				IfExists().
				Identifier("name", g.KindOfT[DatabaseObjectIdentifier]()).
				QueryStructField(dbRoleRename, "Rename", g.KindOfPointer(dbRoleRename.Name), g.ListOptions().NoParens().SQL("RENAME TO")).
				QueryStructField(dbRoleSet, "Set", g.KindOfPointer(dbRoleSet.Name), g.ListOptions().NoParens().SQL("SET")).
				QueryStructField(dbRoleUnset, "Unset", g.KindOfPointer(dbRoleUnset.Name), g.ListOptions().NoParens().SQL("UNSET")).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "Rename", "Set", "Unset"),
		),
)
