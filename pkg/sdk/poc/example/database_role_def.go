package example

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ../main.go

// TODO do we need this ?
var _ = DatabaseRole

var DatabaseRole = g.NewInterface(
	"DatabaseRoles",
	"DatabaseRole",
	"DatabaseObjectIdentifier", // TODO do we need this
	g.NewOperation("Create", "https://docs.snowflake.com/en/sql-reference/sql/create-database-role").
		WithOptionsStruct(
			// TODO why do we need this thing vvv (Should this be NewOptsStruct ???) - Field represents Field or Struct ?
			g.NewOptionsStruct().
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
							// TODO g.NewField("Comment", "string", map[string][]string{"ddl": {"parameter", "single_quotes"}, "sql": {"COMMENT"}}).WithRequired(true),
							g.OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)).WithRequired(true),
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
							g.OptionalSQL("COMMENT"),
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
