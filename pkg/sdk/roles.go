package sdk

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Roles interface {
	// Create creates a role.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateRoleOptions) error
	// Alter modifies an existing role
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterRoleOptions) error
	// Drop removes a role.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropRoleOptions) error
	// Show returns a list of roles.
	Show(ctx context.Context, opts *ShowRoleOptions) ([]*Role, error)
	// ShowByID returns a user by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error)
	// Grant grants privileges on a role.
	Grant(ctx context.Context, id AccountObjectIdentifier, opts *RoleGrantOptions) error
	// Revoke revokes privileges on a role.
	Revoke(ctx context.Context, id AccountObjectIdentifier, opts *RoleRevokeOptions) error
	// Use sets the active role for the current session.
	Use(ctx context.Context, id AccountObjectIdentifier) error
	// UseSecondary specifies the active/current secondary roles for the session.
	UseSecondary(ctx context.Context, opts SecondaryRoleOption) error
}

var _ Roles = (*roles)(nil)

type roles struct {
	client *Client
}

type Role struct {
	CreatedOn       time.Time
	Name            string
	IsDefault       bool
	IsCurrent       bool
	IsInherited     bool
	AssignedToUsers int
	GrantedToRoles  int
	GrantedRoles    int
	Owner           string
	Comment         string
}

func (v *Role) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Role) ObjectType() ObjectType {
	return ObjectTypeRole
}

type roleDBRow struct {
	CreatedOn       time.Time      `db:"created_on"`
	Name            string         `db:"name"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	IsInherited     sql.NullString `db:"is_inherited"`
	AssignedToUsers int            `db:"assigned_to_users"`
	GrantedToRoles  int            `db:"granted_to_roles"`
	GrantedRoles    int            `db:"granted_roles"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
}

func (row *roleDBRow) toRole() *Role {
	role := Role{
		CreatedOn:       row.CreatedOn,
		Name:            row.Name,
		AssignedToUsers: row.AssignedToUsers,
		GrantedToRoles:  row.GrantedToRoles,
		GrantedRoles:    row.GrantedRoles,
	}
	if row.IsDefault.Valid {
		role.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		role.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.IsInherited.Valid {
		role.IsInherited = row.IsInherited.String == "Y"
	}
	if row.Owner.Valid {
		role.Owner = row.Owner.String
	}
	if row.Comment.Valid {
		role.Comment = row.Comment.String
	}

	return &role
}

// CreateRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-role
type CreateRoleOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	role        bool                    `ddl:"static" sql:"ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
	Comment     *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag         []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

func (opts *CreateRoleOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		return errors.New("IF NOT EXISTS and OR REPLACE are incompatible.")
	}
	return nil
}

func (v *roles) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateRoleOptions) error {
	if opts == nil {
		opts = &CreateRoleOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// AlterRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-role
type AlterRoleOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"` //lint:ignore U1000 This is used in the ddl tag
	role     bool                    `ddl:"static" sql:"ROLE"`  //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`

	// One of
	RenameTo AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set      *RoleSet                `ddl:"list,no_parentheses" sql:"SET"`
	Unset    *RoleUnset              `ddl:"list,no_parentheses" sql:"UNSET"`
}

func (opts *AlterRoleOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return errors.New("invalid object identifier")
	}
	if everyValueNil(opts.RenameTo, opts.Set, opts.Unset) {
		return errors.New("No alter action specified")
	}
	if !exactlyOneValueSet(opts.RenameTo, opts.Set, opts.Unset) {
		return errors.New("you can use one action at a time (RENAME TO, SET or UNSET)")
	}
	return nil
}

type RoleSet struct {
	Comment *string          `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag     []TagAssociation `ddl:"keyword" sql:"TAG"`
}

type RoleUnset struct {
	Comment *bool              `ddl:"keyword" sql:"COMMENT"`
	Tag     []ObjectIdentifier `ddl:"keyword" sql:"TAG"`
}

func (v *roles) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterRoleOptions) error {
	if opts == nil {
		return errors.New("alter alert options cannot be empty")
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// TODO - DropRoleOptions vs DropRoleOptions
// DropRoleOptions contains options for dropping a role.
type DropRoleOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"` //lint:ignore U1000 This is used in the ddl tag
	roles    bool                    `ddl:"static" sql:"ROLE"` //lint:ignore U1000 This is used in the ddl tag
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"` //lint:ignore U1000 This is used in the ddl tag
}

func (opts *DropRoleOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *roles) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropRoleOptions) error {
	if opts == nil {
		opts = &DropRoleOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// ShowRoleOptions contains options for listing roles.
type ShowRoleOptions struct {
	show  bool  `ddl:"static" sql:"SHOW"`  //lint:ignore U1000 This is used in the ddl tag
	roles bool  `ddl:"static" sql:"ROLES"` //lint:ignore U1000 This is used in the ddl tag
	Like  *Like `ddl:"keyword" sql:"LIKE"`
}

func (v *roles) Show(ctx context.Context, opts *ShowRoleOptions) ([]*Role, error) {
	if opts == nil {
		opts = &ShowRoleOptions{}
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []roleDBRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	roles := make([]*Role, len(rows))
	for i, row := range rows {
		roles[i] = row.toRole()
	}
	return roles, nil
}

func (v *roles) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error) {
	roles, err := v.client.Roles.Show(ctx, &ShowRoleOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.ID() == id {
			return role, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type RoleGrantOptions struct {
	grant bool                    `ddl:"static" sql:"GRANT"` //lint:ignore U1000 This is used in the ddl tag
	role  bool                    `ddl:"static" sql:"ROLE"`  //lint:ignore U1000 This is used in the ddl tag
	name  AccountObjectIdentifier `ddl:"identifier"`         //lint:ignore U1000 This is used in the ddl tag
	Grant GrantRole               `ddl:"keyword,no_parentheses" sql:"TO"`
}

type GrantRole struct {
	// one of
	Role *AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
	User *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
}

func (opts *RoleGrantOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.Grant) {
		return errors.New("Granted option should be set")
	}
	if !exactlyOneValueSet(opts.Grant.Role, opts.Grant.User) {
		return errors.New("Only one granted option should be set")
	}
	return nil
}

func (v *roles) Grant(ctx context.Context, id AccountObjectIdentifier, opts *RoleGrantOptions) error {
	if opts == nil {
		opts = &RoleGrantOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type RoleRevokeOptions struct {
	revoke bool                    `ddl:"static" sql:"REVOKE"` //lint:ignore U1000 This is used in the ddl tag
	role   bool                    `ddl:"static" sql:"ROLE"`   //lint:ignore U1000 This is used in the ddl tag
	name   AccountObjectIdentifier `ddl:"identifier"`
	Revoke RevokeRole              `ddl:"keyword,no_parentheses" sql:"FROM"`
}

func (opts *RoleRevokeOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	if !valueSet(opts.Revoke) {
		return errors.New("Revoked option should be set")
	}
	if !exactlyOneValueSet(opts.Revoke.Role, opts.Revoke.User) {
		return errors.New("Only one revoked option should be set")
	}
	return nil
}

type RevokeRole struct {
	User *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	Role *AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
}

func (v *roles) Revoke(ctx context.Context, id AccountObjectIdentifier, opts *RoleRevokeOptions) error {
	if opts == nil {
		opts = &RoleRevokeOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *roles) Use(ctx context.Context, id AccountObjectIdentifier) error {
	return v.client.Sessions.UseRole(ctx, id)
}

func (v *roles) UseSecondary(ctx context.Context, opt SecondaryRoleOption) error {
	return v.client.Sessions.UseSecondaryRoles(ctx, opt)
}
