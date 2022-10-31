package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Compile-time proof of interface implementation.
var _ Users = (*users)(nil)

// Users describes all the users related methods that the
// Snowflake API supports.
type Users interface {
	// List all the users by pattern.
	List(ctx context.Context, options UserListOptions) ([]*User, error)
	// Create a new user with the given options.
	Create(ctx context.Context, options UserCreateOptions) (*User, error)
	// Read an user by its name.
	Read(ctx context.Context, user string) (*User, error)
	// Update attributes of an existing user.
	Update(ctx context.Context, user string, options UserUpdateOptions) (*User, error)
	// Delete an user by its name.
	Delete(ctx context.Context, user string) error
}

// users implements Users
type users struct {
	client *Client
}

// User represents a Snowflake user.
type User struct {
	Comment               string
	DefaultNamespace      string
	DefaultRole           string
	DefaultSecondaryRoles []string
	DefaultWarehouse      string
	Disabled              bool
	DisplayName           string
	Email                 string
	FirstName             string
	HasRsaPublicKey       bool
	LastName              string
	LoginName             string
	Name                  string
}

type userEntity struct {
	Name                  sql.NullString `db:"name"`
	Comment               sql.NullString `db:"comment"`
	DefaultNamespace      sql.NullString `db:"default_namespace"`
	DefaultRole           sql.NullString `db:"default_role"`
	DefaultSecondaryRoles sql.NullString `db:"default_secondary_roles"`
	DefaultWarehouse      sql.NullString `db:"default_warehouse"`
	Disabled              sql.NullBool   `db:"disabled"`
	DisplayName           sql.NullString `db:"display_name"`
	Email                 sql.NullString `db:"email"`
	FirstName             sql.NullString `db:"first_name"`
	HasRsaPublicKey       sql.NullBool   `db:"has_rsa_public_key"`
	LastName              sql.NullString `db:"last_name"`
	LoginName             sql.NullString `db:"login_name"`
}

func (e *userEntity) toUser() *User {
	var roles []string
	if e.DefaultSecondaryRoles.Valid {
		value := strings.Trim(e.DefaultSecondaryRoles.String, "[]")
		roles = strings.Split(value, ",")
	}
	return &User{
		Comment:               e.Comment.String,
		DefaultNamespace:      e.DefaultNamespace.String,
		DefaultRole:           e.DefaultRole.String,
		DefaultSecondaryRoles: roles,
		DefaultWarehouse:      e.DefaultWarehouse.String,
		Disabled:              e.Disabled.Bool,
		DisplayName:           e.DisplayName.String,
		Email:                 e.Email.String,
		FirstName:             e.FirstName.String,
		HasRsaPublicKey:       e.HasRsaPublicKey.Bool,
		LastName:              e.LastName.String,
		LoginName:             e.LoginName.String,
		Name:                  e.Name.String,
	}
}

// UserListOptions represents the options for listing users.
type UserListOptions struct {
	Pattern string
}

func (o UserListOptions) validate() error {
	if o.Pattern == "" {
		return errors.New("pattern must not be empty")
	}
	return nil
}

type UserProperties struct {
	// Optional: Name that the user enters to log into the system.
	// Login names for users must be unique across your entire account.
	LoginName *string

	// Optional: Name displayed for the user in the Snowflake web interface.
	DisplayName *string

	// Optional: First, middle, and last name of the user.
	FirstName  *string
	MiddleName *string
	LastName   *string

	// Optional: Email address for the user.
	Email *string

	// Optional: Specifies whether the user is forced to change their password on next login (including their first/initial login) into the system.
	MustChangePassword *bool

	// Optional: Specifies whether the user is disabled
	Disabled *bool

	// Optional: Specifies the virtual warehouse that is active by default for the user’s session upon login.
	DefaultWarehouse *string

	// Optional: Specifies the namespace (database only or database and schema) that is active by default for the user’s session upon login
	DefaultNamespace *string

	// Optional: Specifies the primary role that is active by default for the user’s session upon login
	DefaultRole *string

	// Optional: Specifies the set of secondary roles that are active for the user’s session upon login
	DefaultSecondaryRoles *[]string

	// Optional: Specifies a comment for the user.
	Comment *string
}

// UserCreateOptions represents the options for creating an user.
type UserCreateOptions struct {
	*UserProperties

	// Required: Identifier for the user; must be unique for your account.
	Name string
	// Optional: The password for the user must be enclosed in single or double quotes
	Password *string
}

func (o UserCreateOptions) validate() error {
	if o.Name == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

// UserUpdateOptions represents the options for updating an user.
type UserUpdateOptions struct {
	*UserProperties
}

func (u *users) List(ctx context.Context, options UserListOptions) ([]*User, error) {
	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("validate list options: %w", err)
	}

	query := fmt.Sprintf(`SHOW USERS LIKE '%s'`, options.Pattern)
	rows, err := u.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	entities := []*User{}
	for rows.Next() {
		var entity userEntity
		if err := rows.StructScan(&entity); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		entities = append(entities, entity.toUser())
	}
	return entities, nil
}

func (u *users) Read(ctx context.Context, user string) (*User, error) {
	query := fmt.Sprintf(`SHOW USERS LIKE '%s'`, user)
	rows, err := u.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("do query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	var entity userEntity
	if err := rows.StructScan(&entity); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}
	return entity.toUser(), nil
}

func (u *users) formatUserProperties(properties *UserProperties) string {
	var s string
	if properties.LoginName != nil {
		s = s + " login_name='" + *properties.LoginName + "'"
	}
	if properties.DisplayName != nil {
		s = s + " display_name='" + *properties.DisplayName + "'"
	}
	if properties.FirstName != nil {
		s = s + " first_name='" + *properties.FirstName + "'"
	}
	if properties.MiddleName != nil {
		s = s + " middle_name='" + *properties.MiddleName + "'"
	}
	if properties.LastName != nil {
		s = s + " last_name='" + *properties.LastName + "'"
	}
	if properties.Email != nil {
		s = s + " email='" + *properties.Email + "'"
	}
	if properties.MustChangePassword != nil {
		s = s + fmt.Sprintf(" must_change_password=%t", *properties.MustChangePassword)
	}
	if properties.Disabled != nil {
		s = s + fmt.Sprintf(" disabled=%t", *properties.Disabled)
	}
	if properties.DefaultWarehouse != nil {
		s = s + " default_warehouse='" + *properties.DefaultWarehouse + "'"
	}
	if properties.DefaultNamespace != nil {
		s = s + " default_namespace='" + *properties.DefaultNamespace + "'"
	}
	if properties.DefaultRole != nil {
		s = s + " default_role='" + *properties.DefaultRole + "'"
	}
	if properties.DefaultSecondaryRoles != nil {
		roles := []string{}
		for _, role := range *properties.DefaultSecondaryRoles {
			roles = append(roles, "'"+role+"'")
		}
		s = s + " default_secondary_roles=(" + strings.Join(roles, ",") + ")"
	}
	if properties.Comment != nil {
		s = s + " comment='" + *properties.Comment + "'"
	}
	return s
}

func (u *users) Update(ctx context.Context, user string, opts UserUpdateOptions) (*User, error) {
	if user == "" {
		return nil, errors.New("name must not be empty")
	}
	query := fmt.Sprintf("ALTER USER %s SET", user)
	if opts.UserProperties != nil {
		query = query + u.formatUserProperties(opts.UserProperties)
	}
	if _, err := u.client.Exec(ctx, query); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	return u.Read(ctx, user)
}

func (u *users) Create(ctx context.Context, opts UserCreateOptions) (*User, error) {
	if err := opts.validate(); err != nil {
		return nil, fmt.Errorf("validate create options: %w", err)
	}
	query := fmt.Sprintf("CREATE USER %s", opts.Name)
	if opts.UserProperties != nil {
		query = query + u.formatUserProperties(opts.UserProperties)
	}
	if _, err := u.client.Exec(ctx, query); err != nil {
		return nil, fmt.Errorf("db exec: %w", err)
	}
	return u.Read(ctx, opts.Name)
}

func (u *users) Delete(ctx context.Context, user string) error {
	query := fmt.Sprintf(`DROP USER %s`, user)
	if _, err := u.client.Exec(ctx, query); err != nil {
		return fmt.Errorf("db exec: %w", err)
	}
	return nil
}
