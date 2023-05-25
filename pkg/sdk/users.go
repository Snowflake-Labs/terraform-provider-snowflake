package sdk

import (
	"context"
)

type Users interface {
	// Create creates a user.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *UserCreateOptions) error
	// Alter modifies an existing user
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *UserAlterOptions) error
	// Drop removes a user.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *UserDropOptions) error
	// Describe returns the details of a user.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error)
	// Show returns a list of users.
	Show(ctx context.Context, opts *UserShowOptions) ([]*User, error)
	// ShowByID returns a user by ID
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error)
}

var _ Users = (*users)(nil)

type users struct {
	client *Client
}

type User struct{}

func (v *User) ID() AccountObjectIdentifier {
	return AccountObjectIdentifier{}
}

func (v *User) ObjectType() ObjectType {
	return ObjectTypeUser
}

// UserCreateOptions contains options for creating a user.
type UserCreateOptions struct{}

func (v *users) Create(ctx context.Context, id AccountObjectIdentifier, opts *UserCreateOptions) error {
	return nil
}

// UserAlterOptions contains options for altering a user.
type UserAlterOptions struct{}

func (v *users) Alter(ctx context.Context, id AccountObjectIdentifier, opts *UserAlterOptions) error {
	return nil
}

// UserDropOptions contains options for dropping a user.
type UserDropOptions struct{}

func (v *users) Drop(ctx context.Context, id AccountObjectIdentifier, opts *UserDropOptions) error {
	return nil
}

// UserDetails contains details about a user.
type UserDetails struct{}

func (v *users) Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error) {
	return nil, nil
}

// UserShowOptions contains options for listing users.
type UserShowOptions struct{}

func (v *users) Show(ctx context.Context, opts *UserShowOptions) ([]*User, error) {
	return nil, nil
}

func (v *users) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	return nil, nil
}
