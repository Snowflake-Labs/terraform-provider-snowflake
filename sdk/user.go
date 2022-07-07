package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// Compile-time proof of interface implementation.
var _ Users = (*users)(nil)

type Users interface {
	List(ctx context.Context, options *UserListOptions) ([]*User, error)
	//Create(ctx context.Context, options *UserCreateOptions) (*User, error)
	//Read(ctx context.Context, userID string) (*User, error)
	//Update(ctx context.Context, options *UserUpdateOptions) (*User, error)
	//Delete(ctx context.Context, userID string) error
}

// users implements Users
type users struct {
	client *Client
}

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

type userRecord struct {
	Comment               sql.NullString `db:"comment"`
	DefaultNamespace      sql.NullString `db:"default_namespace"`
	DefaultRole           sql.NullString `db:"default_role"`
	DefaultSecondaryRoles sql.NullString `db:"default_secondary_roles"`
	DefaultWarehouse      sql.NullString `db:"default_warehouse"`
	Disabled              bool           `db:"disabled"`
	DisplayName           sql.NullString `db:"display_name"`
	Email                 sql.NullString `db:"email"`
	FirstName             sql.NullString `db:"first_name"`
	HasRsaPublicKey       bool           `db:"has_rsa_public_key"`
	LastName              sql.NullString `db:"last_name"`
	LoginName             sql.NullString `db:"login_name"`
	Name                  sql.NullString `db:"name"`
}

func (r *userRecord) toUser() *User {
	return &User{
		Comment:               r.Comment.String,
		DefaultNamespace:      r.DefaultNamespace.String,
		DefaultRole:           r.DefaultRole.String,
		DefaultSecondaryRoles: strings.Split(r.DefaultSecondaryRoles.String, ","),
		DefaultWarehouse:      r.DefaultWarehouse.String,
		Disabled:              r.Disabled,
		DisplayName:           r.DisplayName.String,
		Email:                 r.Email.String,
		FirstName:             r.FirstName.String,
		HasRsaPublicKey:       r.HasRsaPublicKey,
		LastName:              r.LastName.String,
		LoginName:             r.LoginName.String,
		Name:                  r.Name.String,
	}
}

// UserListOptions represents the options for listing teams.
type UserListOptions struct {
	Pattern string
}

func (o *UserListOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if o.Pattern == "" {
		return errors.Errorf("pattern must not be an empty string")
	}

	return nil
}

type UserCreateOptions struct {
	Comment               *string
	DefaultNamespace      *string
	DefaultRole           *string
	DefaultSecondaryRoles []string
	DefaultWarehouse      *string
	Disabled              *bool
	DisplayName           *string
	Email                 *string
	FirstName             *string
	HasRsaPublicKey       *bool
	LastName              *string
	LoginName             *string
	Name                  *string
}

type UserUpdateOptions struct {
	Comment               string
	DefaultNamespace      string
	DefaultRole           string
	DefaultSecondaryRoles []string
	DefaultWarehouse      *string
	Disabled              *bool
	DisplayName           *string
	Email                 *string
	FirstName             *string
	HasRsaPublicKey       *bool
	LastName              *string
	LoginName             *string
	Name                  *string
}

func (u *users) List(ctx context.Context, options *UserListOptions) ([]*User, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf(`SHOW USERS like '%s'`, options.Pattern)
	rows, err := u.client.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userRecordList := []userRecord{}
	err = sqlx.StructScan(rows, &userRecordList)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(err, "no users found")
		}
		return nil, errors.Wrapf(err, "unable to scan row for %s", stmt)
	}

	var userList []*User
	for _, userRecord := range userRecordList {
		fmt.Printf("%+v\n", userRecord)
		userList = append(userList, userRecord.toUser())
	}

	return userList, nil
}
