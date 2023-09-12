package example2

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator2"

//go:generate go run ../main2.go

var (
	dbRoleName = g.NewStruct("DatabaseRoleRename").
			WithFields(g.DatabaseObjectIdentifier("Name")).
			WithValidations(g.NewValidation(g.ValidIdentifier, "Name"))

	nestedThirdLevel = g.NewStruct("NestedThirdLevel").
				WithFields(g.DatabaseObjectIdentifier("Field")).
				WithValidations(g.NewValidation(g.AtLeastOneValueSet, "Field"))

	dbRoleSet = g.NewStruct("DatabaseRoleSet").
			WithFields(
			g.OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)).WithRequired(true),
			g.NewStructField(nestedThirdLevel, "NestedThirdLevel", g.KindOfPointer(nestedThirdLevel.Name), g.Tags().List().NoParentheses().SQL("NESTED").Build()),
		)

	dbRoleUnset = g.NewStruct("DatabaseRoleUnset").
			WithFields(g.OptionalSQL("COMMENT")).
			WithValidations(g.NewValidation(g.AtLeastOneValueSet, "Comment"))
)

var DatabaseRole = g.NewInterface(
	"DatabaseRoles",
	"DatabaseRole",
	g.NewOperation(
		"Create",
		"https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
		g.NewStruct("CreateDatabaseRoleOptions").
			WithFields(
				g.Create(),
				g.OrReplace(),
				g.SQL("DATABASE ROLE"),
				g.IfNotExists(),
				g.DatabaseObjectIdentifier("name"),
				g.OptionalTextAssignment("COMMENT", nil),
			).
			WithValidations(
				g.NewValidation(g.ValidIdentifier, "name"),
				g.NewValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
			),
	),
	g.NewOperation(
		"Alter",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
		g.NewStruct("AlterDatabaseRoleOptions").
			WithFields(
				g.Alter(),
				g.SQL("DATABASE ROLE"),
				g.IfExists(),
				g.DatabaseObjectIdentifier("name"),
				g.NewStructField(dbRoleName, "Rename", g.KindOfPointer(dbRoleName.Name), g.Tags().List().NoParentheses().Build()),
				g.NewStructField(dbRoleSet, "Set", g.KindOfPointer(dbRoleSet.Name), g.Tags().List().NoParentheses().Build()),
				g.NewStructField(dbRoleUnset, "Unset", g.KindOfPointer(dbRoleUnset.Name), g.Tags().List().NoParentheses().Build()),
			).
			WithValidations(
				g.NewValidation(g.ValidIdentifier, "name"),
				g.NewValidation(g.ExactlyOneValueSet, "Rename", "Set", "Unset"),
			),
	),
)
