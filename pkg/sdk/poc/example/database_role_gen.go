package example

import "context"

type DatabaseRoles interface {
	Create(ctx context.Context, request *CreateDatabaseRoleRequest) error
	Alter(ctx context.Context, request *AlterDatabaseRoleRequest) error
}

// CreateDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database-role.
type CreateDatabaseRoleOptions struct {
	create       bool                     `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	databaseRole bool                     `ddl:"static" sql:"DATABASE ROLE"`
	IfNotExists  *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         DatabaseObjectIdentifier `ddl:"identifier"`
	Comment      *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database-role.
type AlterDatabaseRoleOptions struct {
	alter        bool                     `ddl:"static" sql:"ALTER"`
	databaseRole bool                     `ddl:"static" sql:"DATABASE ROLE"`
	IfExists     *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name         DatabaseObjectIdentifier `ddl:"identifier"`
	Rename       *DatabaseRoleRename      `ddl:"list,no_parentheses" sql:"RENAME TO"`
	Set          *DatabaseRoleSet         `ddl:"list,no_parentheses" sql:"SET"`
	Unset        *DatabaseRoleUnset       `ddl:"list,no_parentheses" sql:"UNSET"`
}

type DatabaseRoleRename struct {
	Name DatabaseObjectIdentifier `ddl:"identifier"`
}

type DatabaseRoleSet struct {
	Comment          string            `ddl:"parameter,single_quotes" sql:"COMMENT"`
	NestedThirdLevel *NestedThirdLevel `ddl:"list,no_parentheses" sql:"NESTED"`
}

type NestedThirdLevel struct {
	Field DatabaseObjectIdentifier `ddl:"identifier"`
}

type DatabaseRoleUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}
