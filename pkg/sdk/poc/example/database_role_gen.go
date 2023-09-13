package example

import "context"

type DatabaseRoles interface {
	Create(ctx context.Context, request *CreateDatabaseRoleRequest) error
	Alter(ctx context.Context, request *AlterDatabaseRoleRequest) error
}

// CreateDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-database-role.
type CreateDatabaseRoleOptions struct {
	OrReplace   *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	IfNotExists *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        DatabaseObjectIdentifier `ddl:"identifier"`
	Comment     *string                  `ddl:"parameter,single_quotes"`
}

// AlterDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-database-role.
type AlterDatabaseRoleOptions struct {
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`
	Rename   *DatabaseRoleRename      `ddl:"list,no_parentheses" sql:"RENAME TO"`
	Set      *DatabaseRoleSet         `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *DatabaseRoleUnset       `ddl:"list,no_parentheses" sql:"UNSET"`
}

type DatabaseRoleRename struct {
	Name DatabaseObjectIdentifier `ddl:"identifier"`
}

type DatabaseRoleSet struct {
	Comment          *string           `ddl:"parameter,single_quotes"`
	NestedThirdLevel *NestedThirdLevel `ddl:"list,no_parentheses" sql:"NESTED"`
}

type NestedThirdLevel struct {
	Field DatabaseObjectIdentifier `ddl:"identifier"`
}

type DatabaseRoleUnset struct {
	Comment          name              `ddl:"keyword" sql:"COMMENT"`
	NestedThirdLevel *NestedThirdLevel `ddl:"list,no_parentheses" sql:"NESTED"`
}

type NestedThirdLevel struct {
	Field DatabaseObjectIdentifier `ddl:"identifier"`
}
