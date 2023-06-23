package sdk

import (
	"context"
	"fmt"
)

type Roles interface {
	// Create creates a role.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *RoleCreateOptions) error
	// Alter modifies an existing role
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *RoleAlterOptions) error
	// Drop removes a role.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *RoleDropOptions) error
	// Show returns a list of roles.
	Show(ctx context.Context, opts *RoleShowOptions) ([]*Role, error)
	// ShowByID returns a user by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error)
	// Grant grants privileges on a role.
	// Grant(ctx context.Context, id AccountObjectIdentifier, opts *RoleGrantOptions) error
	// Revoke revokes privileges on a role.
	// Revoke(ctx context.Context, id AccountObjectIdentifier, opts *RoleRevokeOptions) error
	// Use sets the active role for the current session.
	// Use(ctx context.Context, id AccountObjectIdentifier) error
	// UseSecondary specifies the active/current secondary roles for the session.
	// UseSecondary(ctx context.Context, opts *RolesUseSecondaryOptions) error
}

var _ Roles = (*roles)(nil)

type roles struct {
	client *Client
}

type Role struct {
	Name string `db:"name"`
}

func (v *Role) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Role) ObjectType() ObjectType {
	return ObjectTypeRole
}

// RoleCreateOptions contains options for creating a role.
type RoleCreateOptions struct{}

func (v *roles) Create(ctx context.Context, id AccountObjectIdentifier, opts *RoleCreateOptions) error {
	sql := fmt.Sprintf(`CREATE ROLE %v`, id)
	_, err := v.client.exec(ctx, sql)
	return err
}

// RoleAlterOptions contains options for altering a user.
type RoleAlterOptions struct{}

func (v *roles) Alter(ctx context.Context, id AccountObjectIdentifier, opts *RoleAlterOptions) error {
	return nil
}

// RoleDropOptions contains options for dropping a role.
type RoleDropOptions struct{}

func (v *roles) Drop(ctx context.Context, id AccountObjectIdentifier, opts *RoleDropOptions) error {
	sql := fmt.Sprintf(`DROP ROLE %v`, id.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

// RoleShowOptions contains options for listing roles.
type RoleShowOptions struct{}

func (v *roles) Show(ctx context.Context, opts *RoleShowOptions) ([]*Role, error) {
	var rows []Role
	sql := `SHOW ROLES`
	err := v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}

	roles := make([]*Role, len(rows))
	for i, row := range rows {
		var role Role = row
		roles[i] = &role
	}

	return roles, nil
}

func (v *roles) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error) {
	return nil, nil
}
