package sdk

import "context"

type DatabaseRoles interface {
	Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateDatabaseRoleOptions) error
	Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterDatabaseRoleOptions) error
}

// CreateDatabaseRoleOptions contains options for creating a new database role or replace an existing database role in the system.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-database-role.
type CreateDatabaseRoleOptions struct {
	create       bool                     `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace    *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	databaseRole bool                     `ddl:"static" sql:"DATABASE ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists  *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         DatabaseObjectIdentifier `ddl:"identifier"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterDatabaseRoleOptions contains options for modifying a limited set of properties for an existing database role object.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/alter-database-role.
type AlterDatabaseRoleOptions struct {
	alter    bool                     `ddl:"static" sql:"ALTER"`         //lint:ignore U1000 This is used in the ddl tag
	role     bool                     `ddl:"static" sql:"DATABASE ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     DatabaseObjectIdentifier `ddl:"identifier"`

	// One of
	Rename *DatabaseRoleRename `ddl:"list,no_parentheses" sql:"RENAME TO"`
	Set    *DatabaseRoleSet    `ddl:"list,no_parentheses" sql:"SET"`
	Unset  *DatabaseRoleUnset  `ddl:"list,no_parentheses" sql:"UNSET"`
}

type DatabaseRoleRename struct {
	Name DatabaseObjectIdentifier `ddl:"identifier"`
}

type DatabaseRoleSet struct {
	Comment string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type DatabaseRoleUnset struct {
	Comment bool `ddl:"keyword" sql:"COMMENT"`
}
