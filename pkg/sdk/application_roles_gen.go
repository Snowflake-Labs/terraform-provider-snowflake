package sdk

import (
	"context"
	"time"
)

type ApplicationRoles interface {
	Create(ctx context.Context, request *CreateApplicationRoleRequest) error
	Alter(ctx context.Context, request *AlterApplicationRoleRequest) error
	Drop(ctx context.Context, request *DropApplicationRoleRequest) error
	Show(ctx context.Context, request *ShowApplicationRoleRequest) ([]ApplicationRole, error)
	Grant(ctx context.Context, request *GrantApplicationRoleRequest) error
	Revoke(ctx context.Context, request *RevokeApplicationRoleRequest) error
}

// CreateApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-application-role.
type CreateApplicationRoleOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	applicationRole bool                    `ddl:"static" sql:"APPLICATION ROLE"`
	IfNotExists     *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	Comment         *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-application-role.
type AlterApplicationRoleOptions struct {
	alter           bool                     `ddl:"static" sql:"ALTER"`
	applicationRole bool                     `ddl:"static" sql:"APPLICATION ROLE"`
	IfExists        *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier  `ddl:"identifier"`
	RenameTo        *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetComment      *string                  `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment    *bool                    `ddl:"keyword" sql:"UNSET COMMENT"`
}

// DropApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-application-role.
type DropApplicationRoleOptions struct {
	drop            bool                    `ddl:"static" sql:"DROP"`
	applicationRole bool                    `ddl:"static" sql:"APPLICATION ROLE"`
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
}

// ShowApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-application-roles.
type ShowApplicationRoleOptions struct {
	show                          bool                      `ddl:"static" sql:"SHOW"`
	applicationRolesInApplication bool                      `ddl:"static" sql:"APPLICATION ROLES IN APPLICATION"`
	ApplicationName               AccountObjectIdentifier   `ddl:"identifier"`
	LimitFrom                     *LimitFromApplicationRole `ddl:"keyword" sql:"LIMIT"`
}

type applicationRoleDbRow struct {
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	Owner         string    `db:"owner"`
	Comment       string    `db:"comment"`
	OwnerRoleType string    `db:"owner_role_type"`
}

type ApplicationRole struct {
	CreatedOn     time.Time
	Name          string
	Owner         string
	Comment       string
	OwnerRoleTYpe string
}

type LimitFromApplicationRole struct {
	Rows int     `ddl:"keyword"`
	From *string `ddl:"parameter,single_quotes,no_equals" sql:"FROM"`
}

// GrantApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/grant-application-roles.
type GrantApplicationRoleOptions struct {
	grant           bool                    `ddl:"static" sql:"GRANT"`
	applicationRole bool                    `ddl:"static" sql:"APPLICATION ROLE"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	GrantTo         ApplicationGrantOptions `ddl:"keyword" sql:"TO"`
}

type ApplicationGrantOptions struct {
	ParentRole      *AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
	ApplicationRole *AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION ROLE"`
	Application     *AccountObjectIdentifier `ddl:"identifier" sql:"APPLICATION"`
}

// RevokeApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/revoke-application-roles.
type RevokeApplicationRoleOptions struct {
	revoke          bool                    `ddl:"static" sql:"REVOKE"`
	applicationRole bool                    `ddl:"static" sql:"APPLICATION ROLE"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	RevokeFrom      ApplicationGrantOptions `ddl:"keyword" sql:"FROM"`
}
