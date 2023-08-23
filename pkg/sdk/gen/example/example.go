package example

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/gen/generator"

var DatabaseRole = generator.Interface{
	Name:         "DatabaseRoles",
	NameSingular: "DatabaseRole",
	Operations: []*generator.Operation{
		{
			Name:            "Create",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
			OptsStructFields: []*generator.Field{
				{
					Name: "create",
					Kind: "bool",
					Tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"CREATE"},
					},
				},
				{
					Name: "OrReplace",
					Kind: "*bool",
					Tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"OR REPLACE"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					Tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfNotExists",
					Kind: "*bool",
					Tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF NOT EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "*bool",
					Tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Comment",
					Kind: "*string",
					Tags: map[string][]string{
						"ddl": {"parameter", "single_quotes"},
						"sql": {"COMMENT"},
					},
				},
			},
			Validations: []*generator.Validation{
				{
					Type:       generator.ValidIdentifier,
					FieldNames: []string{"name"},
				},
				{
					Type:       generator.ConflictingFields,
					FieldNames: []string{"OrReplace", "IfNotExists"},
				},
			},
		},
		{
			Name:            "Alter",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
			OptsStructFields: []*generator.Field{
				{
					Name: "alter",
					Kind: "bool",
					Tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"ALTER"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					Tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfExists",
					Kind: "*bool",
					Tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "DatabaseObjectIdentifier",
					Tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Rename",
					Kind: "*DatabaseRoleRename",
					Tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"RENAME TO"},
					},
					Fields: []*generator.Field{
						{
							Name: "Name",
							Kind: "DatabaseObjectIdentifier",
							Tags: map[string][]string{
								"ddl": {"identifier"},
							},
						},
					},
					Validations: []*generator.Validation{
						{
							Type:       generator.ValidIdentifier,
							FieldNames: []string{"Name"},
						},
					},
				},
				{
					Name: "Set",
					Kind: "*DatabaseRoleSet",
					Tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"SET"},
					},
					Fields: []*generator.Field{
						{
							Name: "Comment",
							Kind: "string",
							Tags: map[string][]string{
								"ddl": {"parameter", "single_quotes"},
								"sql": {"COMMENT"},
							},
						},
					},
				},
				{
					Name: "Unset",
					Kind: "*DatabaseRoleUnset",
					Tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"UNSET"},
					},
					Fields: []*generator.Field{
						{
							Name: "Comment",
							Kind: "bool",
							Tags: map[string][]string{
								"ddl": {"keyword"},
								"sql": {"COMMENT"},
							},
						},
					},
					Validations: []*generator.Validation{
						{
							Type:       generator.AtLeastOneValueSet,
							FieldNames: []string{"Comment"},
						},
					},
				},
			},
			Validations: []*generator.Validation{
				{
					Type:       generator.ValidIdentifier,
					FieldNames: []string{"name"},
				},
				{
					Type:       generator.ExactlyOneValueSet,
					FieldNames: []string{"Rename", "Set", "Unset"},
				},
			},
		},
	},
}
