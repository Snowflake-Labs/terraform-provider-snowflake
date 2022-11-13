package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	ResourceRole  = "ROLE"
	ResourceRoles = "ROLES"
)

// Compile-time proof of interface implementation.
var _ Roles = (*roles)(nil)

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
	// Delete a role by its name.
	Delete(ctx context.Context, role string) error
	// Rename a role name.
	Rename(ctx context.Context, old string, new string) error
}

// roles implements Roles
type roles struct {
	client *Client
}

// Role represents a Snowflake role.
type Role struct {
	Name            string
	CreatedOn       time.Time
	IsDefault       string
	IsCurrent       string
	IsInherited     string
	AssignedToUsers int32
	GrantedToRoles  int32
	GrantedRoles    int32
	Owner           string
	Comment         string
}

type roleEntity struct {
	Name            sql.NullString `db:"name"`
	CreatedOn       sql.NullTime   `db:"created_on"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	IsInherited     sql.NullString `db:"is_inherited"`
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
		IsDefault:       e.IsDefault.String,
		IsCurrent:       e.IsCurrent.String,
		IsInherited:     e.IsInherited.String,
		AssignedToUsers: e.AssignedToUsers.Int32,
		GrantedToRoles:  e.GrantedToRoles.Int32,
		GrantedRoles:    e.GrantedRoles.Int32,
		Owner:           e.Owner.String,
		Comment:         e.Comment.String,
	}
}

// RoleListOptions represents the options for listing roles.
type RoleListOptions struct {
	// Required: Filters the command output by object name
	Pattern string

	// Optional: Limits the maximum number of rows returned
	Limit *int
}

func (o RoleListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

type RoleProperties struct {
	// Optional: Specifies a comment for the role.
	Comment *string
}

// RoleCreateOptions represents the options for creating a role.
type RoleCreateOptions struct {
	*RoleProperties

	// Required: Identifier for the role; must be unique for your account.
	Name string
}

func (o RoleCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

// RoleUpdateOptions represents the options for updating a role.
type RoleUpdateOptions struct {
	*RoleProperties
}

// List all the roles by pattern.
func (r *roles) List(ctx context.Context, options RoleListOptions) ([]*Role, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	sql := fmt.Sprintf("SHOW %s LIKE '%s'", ResourceRoles, options.Pattern)
	rows, err := r.client.query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*Role{}
	for rows.Next() {
		var entity roleEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toRole())
	}
	return entities, nil
}

// Read a role by its name.
func (r *roles) Read(ctx context.Context, role string) (*Role, error) {
	var entity roleEntity
	if err := r.client.read(ctx, ResourceRoles, role, &entity); err != nil {
		return nil, err
	}
	return entity.toRole(), nil
}

func (*roles) formatRoleProperties(properties *RoleProperties) string {
	var s string
	if properties.Comment != nil {
		s = s + " comment='" + *properties.Comment + "'"
	}
	return s
}

// Update attributes of an existing role.
func (r *roles) Update(ctx context.Context, role string, options RoleUpdateOptions) (*Role, error) {
	if role == "" {
		return nil, errors.New("role name must not be empty")
	}
	sql := fmt.Sprintf("ALTER %s %s SET", ResourceRole, role)
	if options.RoleProperties != nil {
		sql = sql + r.formatRoleProperties(options.RoleProperties)
	}
	if _, err := r.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity roleEntity
	if err := r.client.read(ctx, ResourceRoles, role, &entity); err != nil {
		return nil, err
	}
	return entity.toRole(), nil
}

// Create a new role with the given options.
func (r *roles) Create(ctx context.Context, options RoleCreateOptions) (*Role, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	sql := fmt.Sprintf("CREATE %s %s", ResourceRole, options.Name)
	if options.RoleProperties != nil {
		sql = sql + r.formatRoleProperties(options.RoleProperties)
	}
	if _, err := r.client.exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	var entity roleEntity
	if err := r.client.read(ctx, ResourceRoles, options.Name, &entity); err != nil {
		return nil, err
	}
	return entity.toRole(), nil
}

// Delete a role by its name.
func (r *roles) Delete(ctx context.Context, role string) error {
	return r.client.drop(ctx, ResourceRole, role)
}

// Rename a role name.
func (r *roles) Rename(ctx context.Context, old string, new string) error {
	return r.client.rename(ctx, ResourceRole, old, new)
}
