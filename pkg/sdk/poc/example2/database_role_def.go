package example2

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator2"

//go:generate go run ../main2.go

var (
	dbRoleName = g.NewStruct("DatabaseRoleRename").
			WithFields(g.DatabaseObjectIdentifier("Name"))
	nestedThirdLevel = g.NewStruct("NestedThirdLevel").
				WithFields(g.DatabaseObjectIdentifier("Field"))
	dbRoleSet = g.NewStruct("DatabaseRoleSet").
			WithFields(g.OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)).WithRequired(true))
	dbRoleUnset = g.NewStruct("DatabaseRoleUnset").
			WithFields(g.OptionalSQL("COMMENT"))
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
				g.NewField("Rename", g.KindOfPointer(dbRoleName.Name), map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"RENAME TO"}}),
				g.NewField("Set", g.KindOfPointer(dbRoleSet.Name), map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"SET"}}),
				g.NewField("Unset", g.KindOfPointer(dbRoleUnset.Name), map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"UNSET"}}),
			),
		dbRoleName,
		nestedThirdLevel,
		dbRoleSet,
	),

	//WithOptionsStruct(
	//	// TODO why do we need this thing vvv (Should this be NewOptsStruct ???) - Field represents Field or Struct ?
	//	g.NewOptionsStruct().
	//		WithFields(
	//			g.Create(),
	//			g.OrReplace(),
	//			g.SQL("DATABASE ROLE"),
	//			g.IfNotExists(),
	//			g.DatabaseObjectIdentifier("name"),
	//			g.OptionalTextAssignment("COMMENT", nil),
	//		).
	//		WithValidations(
	//			g.NewValidation(g.ValidIdentifier, "name"),
	//			g.NewValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	//		),
	//),
	//g.NewOperation("Alter", "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role").
	//	WithOptionsStruct(
	//		g.NewOptionsStruct().
	//			WithFields(
	//				g.Alter(),
	//				g.SQL("DATABASE ROLE"),
	//				g.IfExists(),
	//				g.DatabaseObjectIdentifier("name"),
	//				g.NewField("Rename", "*DatabaseRoleRename", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"RENAME TO"}}).
	//					WithFields(
	//						g.DatabaseObjectIdentifier("Name"),
	//					).
	//					WithValidations(
	//						g.NewValidation(g.ValidIdentifier, "Name"),
	//					),
	//				g.NewField("Set", "*DatabaseRoleSet", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"SET"}}).
	//					WithFields(
	//						// TODO g.NewField("Comment", "string", map[string][]string{"ddl": {"parameter", "single_quotes"}, "sql": {"COMMENT"}}).WithRequired(true),
	//						g.OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes(true)).WithRequired(true),
	//						g.NewField("NestedThirdLevel", "*NestedThirdLevel", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"NESTED"}}).
	//							WithFields(
	//								g.DatabaseObjectIdentifier("Field"),
	//							).
	//							WithValidations(
	//								g.NewValidation(g.AtLeastOneValueSet, "Field"),
	//							),
	//					),
	//				g.NewField("Unset", "*DatabaseRoleUnset", map[string][]string{"ddl": {"list,no_parentheses"}, "sql": {"UNSET"}}).
	//					WithFields(
	//						g.OptionalSQL("COMMENT"),
	//					).
	//					WithValidations(
	//						g.NewValidation(g.AtLeastOneValueSet, "Comment"),
	//					),
	//			).
	//			WithValidations(
	//				g.NewValidation(g.ValidIdentifier, "name"),
	//				g.NewValidation(g.ExactlyOneValueSet, "Rename", "Set", "Unset"),
	//			),
	//	),
)
