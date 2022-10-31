package sdk

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Roles describes all the roles related methods that the
// Snowflake API supports.
type Roles interface {
	// List all the roles by pattern.
	List(ctx context.Context, options RoleListOptions) ([]*Role, error)
	// Create a new role with the given options.
	Create(ctx context.Context, options RoleCreateOptions) (*Role, error)
	// Read an role by its name.
	Read(ctx context.Context, role string) (*Role, error)
	// Update attributes of an existing role.
	Update(ctx context.Context, role string, options RoleUpdateOptions) (*Role, error)
	// Delete an role by its name.
	Delete(ctx context.Context, role string) error
}

// roles implements Roles
type roles struct {
	client *Client
}

// Role represents a Snowflake role.
type Role struct {
	Name            string
	CreatedOn       time.Time
	IsDefault       bool
	IsCurrent       bool
	IsInherited     bool
	AssignedToUsers int32
	GrantedToRoles  int32
	GrantedRoles    int32
	Owner           string
	Comment         string
}

type roleEntity struct {
	Name            sql.NullString `db:"name"`
	CreatedOn       sql.NullTime   `db:"created_on"`
	IsDefault       sql.NullBool   `db:"is_default"`
	IsCurrent       sql.NullBool   `db:"is_current"`
	IsInherited     sql.NullBool   `db:"is_inherited"`
	AssignedToUsers sql.NullInt32  `db:"assigned_to_users"`
	GrantedToRoles  sql.NullInt32  `db:"granted_to_roles"`
	GrantedRoles    sql.NullInt32  `db:"granted_roles"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
}

func (e *roleEntity) toRole() *Role {
	return &Role{
		Name:            e.Name.String,
		CreatedOn:       e.CreatedOn.Time,
		IsDefault:       e.IsDefault.Bool,
		IsCurrent:       e.IsCurrent.Bool,
		IsInherited:     e.IsInherited.Bool,
		AssignedToUsers: e.AssignedToUsers.Int32,
		GrantedToRoles:  e.GrantedToRoles.Int32,
		GrantedRoles:    e.GrantedRoles.Int32,
		Owner:           e.Owner.String,
		Comment:         e.Comment.String,
	}
}

// RoleListOptions represents the options for listing roles.
type RoleListOptions struct {
	Pattern string
}

func (o RoleListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

type RoleCreateOptions struct {
}

type RoleUpdateOptions struct {
}

type RoleProperties struct {
	// Optional: Specifies a comment for the role.
	Comment *string
}
