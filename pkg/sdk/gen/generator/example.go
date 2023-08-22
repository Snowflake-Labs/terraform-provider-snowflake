package generator

var DatabaseRoleInterface = Interface{
	Name:         "DatabaseRoles",
	nameSingular: "DatabaseRole",
	Operations: []*Operation{
		{
			Name:            "Create",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
			OptsStructFields: []*Field{
				{
					Name: "create",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"CREATE"},
					},
				},
				{
					Name: "OrReplace",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"OR REPLACE"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfNotExists",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF NOT EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Comment",
					Kind: "*string",
					tags: map[string][]string{
						"ddl": {"parameter", "single_quotes"},
						"sql": {"COMMENT"},
					},
				},
			},
		},
		{
			Name:            "Alter",
			ObjectInterface: nil,
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
			OptsStructFields: []*Field{
				{
					Name: "alter",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"ALTER"},
					},
				},
				{
					Name: "databaseRole",
					Kind: "bool",
					tags: map[string][]string{
						"ddl": {"static"},
						"sql": {"DATABASE ROLE"},
					},
				},
				{
					Name: "IfExists",
					Kind: "*bool",
					tags: map[string][]string{
						"ddl": {"keyword"},
						"sql": {"IF EXISTS"},
					},
				},
				{
					Name: "name",
					Kind: "DatabaseObjectIdentifier",
					tags: map[string][]string{
						"ddl": {"identifier"},
					},
				},
				{
					Name: "Rename",
					Kind: "*DatabaseRoleRename",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"RENAME TO"},
					},
				},
				{
					Name: "Set",
					Kind: "*DatabaseRoleSet",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"SET"},
					},
				},
				{
					Name: "Unset",
					Kind: "*DatabaseRoleUnset",
					tags: map[string][]string{
						"ddl": {"list,no_parentheses"},
						"sql": {"UNSET"},
					},
				},
			},
		},
	},
}
