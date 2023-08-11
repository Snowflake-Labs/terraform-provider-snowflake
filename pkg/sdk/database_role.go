package sdk

import (
	"context"
	"database/sql"
)

type DatabaseRoles interface {
	Create(ctx context.Context, id DatabaseObjectIdentifier, opts *CreateDatabaseRoleOptions) error
	Alter(ctx context.Context, id DatabaseObjectIdentifier, opts *AlterDatabaseRoleOptions) error
	Drop(ctx context.Context, id DatabaseObjectIdentifier) error
	Show(ctx context.Context, opts *ShowDatabaseRoleOptions) ([]*DatabaseRole, error)
	ShowByID(ctx context.Context, id DatabaseObjectIdentifier) (*DatabaseRole, error)
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
	alter        bool                     `ddl:"static" sql:"ALTER"`         //lint:ignore U1000 This is used in the ddl tag
	databaseRole bool                     `ddl:"static" sql:"DATABASE ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists     *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name         DatabaseObjectIdentifier `ddl:"identifier"`

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

// DropDatabaseRoleOptions contains options for removing the specified database role.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/drop-database-role.
type DropDatabaseRoleOptions struct {
	drop         bool                     `ddl:"static" sql:"DROP"`          //lint:ignore U1000 This is used in the ddl tag
	databaseRole bool                     `ddl:"static" sql:"DATABASE ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists     *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name         DatabaseObjectIdentifier `ddl:"identifier"`
}

// ShowDatabaseRoleOptions contains options for showing database roles in given database.
// At the time of writing LIKE is not visible in the docs, but it works.
//
// Based on https://docs.snowflake.com/en/sql-reference/sql/show-database-roles.
type ShowDatabaseRoleOptions struct {
	show          bool                    `ddl:"static" sql:"SHOW"`           //lint:ignore U1000 This is used in the ddl tag
	databaseRoles bool                    `ddl:"static" sql:"DATABASE ROLES"` //lint:ignore U1000 This is used in the ddl tag
	Like          *Like                   `ddl:"keyword" sql:"LIKE"`
	in            bool                    `ddl:"static" sql:"IN DATABASE"` //lint:ignore U1000 This is used in the ddl tag
	database      AccountObjectIdentifier `ddl:"identifier"`
}

// databaseRoleDBRow is used to decode the result of a SHOW DATABASE ROLES query.
type databaseRoleDBRow struct {
	CreatedOn              string         `db:"created_on"`
	Name                   string         `db:"name"`
	IsDefault              sql.NullString `db:"is_default"`
	IsCurrent              sql.NullString `db:"is_current"`
	IsInherited            sql.NullString `db:"is_inherited"`
	GrantedToRoles         sql.NullString `db:"granted_to_roles"`
	GrantedToDatabaseRoles sql.NullString `db:"granted_to_database_roles"`
	GrantedDatabaseRoles   sql.NullString `db:"granted_database_roles"`
	Owner                  string         `db:"owner"`
	Comment                sql.NullString `db:"comment"`
	OwnerRoleType          sql.NullString `db:"owner_role_type"`
}

// DatabaseRole is a user-friendly result for a SHOW DATABASE ROLES query.
// At the time of writing there is no format specified in the docs.
type DatabaseRole struct {
	CreatedOn              string
	Name                   string
	IsDefault              bool
	IsCurrent              bool
	IsInherited            bool
	GrantedToRoles         string
	GrantedToDatabaseRoles string
	GrantedDatabaseRoles   string
	Owner                  string
	Comment                string
	OwnerRoleType          string
}

func (row databaseRoleDBRow) toDatabaseRole() *DatabaseRole {
	databaseRole := DatabaseRole{
		CreatedOn: row.CreatedOn,
		Name:      row.Name,
		Owner:     row.Owner,
	}
	if row.IsDefault.Valid {
		databaseRole.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		databaseRole.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.IsInherited.Valid {
		databaseRole.IsInherited = row.IsInherited.String == "Y"
	}
	if row.GrantedToRoles.Valid {
		databaseRole.GrantedToRoles = row.GrantedToRoles.String
	}
	if row.GrantedToDatabaseRoles.Valid {
		databaseRole.GrantedToDatabaseRoles = row.GrantedToDatabaseRoles.String
	}
	if row.GrantedDatabaseRoles.Valid {
		databaseRole.GrantedDatabaseRoles = row.GrantedDatabaseRoles.String
	}
	if row.Comment.Valid {
		databaseRole.Comment = row.Comment.String
	}
	if row.OwnerRoleType.Valid {
		databaseRole.OwnerRoleType = row.OwnerRoleType.String
	}
	return &databaseRole
}
