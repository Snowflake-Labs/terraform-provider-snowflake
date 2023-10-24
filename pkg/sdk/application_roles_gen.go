package sdk

import (
	"context"
	"time"
)

// ApplicationRoles is an interface that allows for querying application roles.
// It does not allow for other DDL queries (CREATE, ALTER, DROP, ...) to be called, because they are not possible
// to be called from the program level. Application roles are a special case where they're only usable
// inside application context (e.g. setup.sql). Right now, they can be only manipulated from the program context
// by applying debug_mode parameter to the application, but it's a hacky solution and even with that you're limited with GRANT and REVOKE options.
// That's why we're only exposing SHOW operations, because only they are the only allowed operations to be called from the program context.
type ApplicationRoles interface {
	Show(ctx context.Context, request *ShowApplicationRoleRequest) ([]ApplicationRole, error)
	ShowByID(ctx context.Context, request *ShowByIDApplicationRoleRequest) (*ApplicationRole, error)
}

// ShowApplicationRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-application-roles.
type ShowApplicationRoleOptions struct {
	show                          bool                    `ddl:"static" sql:"SHOW"`
	applicationRolesInApplication bool                    `ddl:"static" sql:"APPLICATION ROLES IN APPLICATION"`
	ApplicationName               AccountObjectIdentifier `ddl:"identifier"`
	Limit                         *LimitFrom              `ddl:"keyword" sql:"LIMIT"`
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
	OwnerRoleType string
}
