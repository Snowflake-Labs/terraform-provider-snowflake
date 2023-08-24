package example

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ../main.go

var _ = DatabaseRole

var DatabaseRole = g.NewInterface("DatabaseRoles", "DatabaseRole", "DatabaseObjectIdentifier").WithOperations(
	[]*g.Operation{
		g.NewOperation("Create", "https://docs.snowflake.com/en/sql-reference/sql/create-database-role").WithOptsField(
			g.NewField("<should be updated programmatically>", "<should be updated programmatically>", nil).
				WithFields([]*g.Field{
					g.NewField("create", "bool", map[string][]string{"ddl": {"static"}, "sql": {"CREATE"}}),
					g.NewField("OrReplace", "*bool", map[string][]string{"ddl": {"keyword"}, "sql": {"OR REPLACE"}}),
					g.NewField("databaseRole", "bool", map[string][]string{"ddl": {"static"}, "sql": {"DATABASE ROLE"}}),
					g.NewField("IfNotExists", "*bool", map[string][]string{"ddl": {"keyword"}, "sql": {"IF NOT EXISTS"}}),
					g.NewField("name", "DatabaseObjectIdentifier", map[string][]string{"ddl": {"identifier"}}).WithRequired(true),
					g.NewField("Comment", "*string", map[string][]string{"ddl": {"parameter", "single_quotes"}, "sql": {"COMMENT"}}),
				}).
				WithValidations([]*g.Validation{
					g.NewValidation(g.ValidIdentifier, []string{"name"}),
					g.NewValidation(g.ConflictingFields, []string{"OrReplace", "IfNotExists"}),
				}),
		),
		g.NewOperation("Alter", "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role").WithOptsField(
			g.NewField("<should be updated programmatically>", "<should be updated programmatically>", nil).
				WithFields([]*g.Field{
					g.NewField("alter", "bool", map[string][]string{"ddl": {"static"}, "sql": {"ALTER"}}),
					g.NewField("databaseRole", "bool", map[string][]string{"ddl": {"static"}, "sql": {"DATABASE ROLE"}}),
					g.NewField("IfExists", "*bool", map[string][]string{"ddl": {"keyword"}, "sql": {"IF EXISTS"}}),
					g.NewField("name", "DatabaseObjectIdentifier", map[string][]string{"ddl": {"identifier"}}).WithRequired(true),
					g.NewField("Rename", "*DatabaseRoleRename", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"RENAME TO"}}).
						WithFields([]*g.Field{
							g.NewField("Name", "DatabaseObjectIdentifier", map[string][]string{"ddl": {"identifier"}}).WithRequired(true),
						}).
						WithValidations([]*g.Validation{
							g.NewValidation(g.ValidIdentifier, []string{"Name"}),
						}),
					g.NewField("Set", "*DatabaseRoleSet", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"SET"}}).
						WithFields([]*g.Field{
							g.NewField("Comment", "string", map[string][]string{"ddl": {"parameter", "single_quotes"}, "sql": {"COMMENT"}}).WithRequired(true),
						}),
					g.NewField("Unset", "*DatabaseRoleUnset", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"UNSET"}}).
						WithFields([]*g.Field{
							g.NewField("Comment", "bool", map[string][]string{"ddl": {"keyword"}, "sql": {"COMMENT"}}).WithRequired(true),
						}).
						WithValidations([]*g.Validation{
							g.NewValidation(g.AtLeastOneValueSet, []string{"Comment"}),
						}),
				}).
				WithValidations([]*g.Validation{
					g.NewValidation(g.ValidIdentifier, []string{"name"}),
					g.NewValidation(g.ExactlyOneValueSet, []string{"Rename", "Set", "Unset"}),
				}),
		),
	},
)
