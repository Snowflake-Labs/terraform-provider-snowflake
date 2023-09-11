package example

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ../main.go

var DatabaseRole = g.NewInterface(
	"DatabaseRoles",
	"DatabaseRole",
	"DatabaseObjectIdentifier",
	g.NewOperation("Create", "https://docs.snowflake.com/en/sql-reference/sql/create-database-role").
		WithOptionsStruct(
			g.NewOptionsStruct().
				WithFields(
					g.Create(),
					g.OrReplace(),
					g.SQL("DATABASE ROLE"),
					g.IfNotExists(),
					g.DatabaseObjectIdentifier("name"),
					g.OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)),
				).
				WithValidations(
					g.NewValidation(g.ValidIdentifier, "name"),
					g.NewValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
				),
		),
	g.NewOperation("Alter", "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role").
		WithOptionsStruct(
			g.NewOptionsStruct().
				WithFields(
					g.Alter(),
					g.SQL("DATABASE ROLE"),
					g.IfExists(),
					g.DatabaseObjectIdentifier("name"),
					g.NewField("Rename", "*DatabaseRoleRename", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"RENAME TO"}}).
						WithFields(
							g.DatabaseObjectIdentifier("Name"),
						).
						WithValidations(
							g.NewValidation(g.ValidIdentifier, "Name"),
						),
					g.NewField("Set", "*DatabaseRoleSet", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"SET"}}).
						WithFields(
							g.TextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)).WithRequired(true),
							g.NewField("NestedThirdLevel", "*NestedThirdLevel", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"NESTED"}}).
								WithFields(
									g.DatabaseObjectIdentifier("Field"),
								).
								WithValidations(
									g.NewValidation(g.AtLeastOneValueSet, "Field"),
								),
						),
					g.NewField("Unset", "*DatabaseRoleUnset", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"UNSET"}}).
						WithFields(
							g.OptionalSQL("COMMENT").WithRequired(true),
						).
						WithValidations(
							g.NewValidation(g.AtLeastOneValueSet, "Comment"),
						),
				).
				WithValidations(
					g.NewValidation(g.ValidIdentifier, "name"),
					g.NewValidation(g.ExactlyOneValueSet, "Rename", "Set", "Unset"),
				),
		),
)
