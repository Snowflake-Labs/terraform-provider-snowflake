package sdk

import "context"

type DatabaseRoles interface {
	// Create creates a database role.
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateDatabaseRoleOptions) error
}

// CreateDatabaseRoleOptions contains options for creating a new database role or replace an existing database role in the system.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/create-database-role.
type CreateDatabaseRoleOptions struct {
	create       bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace    *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	databaseRole bool                   `ddl:"static" sql:"DATABASE ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists  *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}
