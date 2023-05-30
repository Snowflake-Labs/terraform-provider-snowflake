package sdk

import (
	"context"
)

type Users interface {
	// Create creates a user.
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateUserOptions) error
	// Alter modifies an existing user
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterUserOptions) error
	// Drop removes a user.
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropUserOptions) error
	// Describe returns the details of a user.
	Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error)
	// Show returns a list of users.
	Show(ctx context.Context, opts *ShowUserOptions) ([]*User, error)
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

// CreateUserOptions contains options for creating a user.
type CreateUserOptions struct{}

func (v *users) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateUserOptions) error {
	return nil
}

// AlterUserOptions contains options for altering a user.
type AlterUserOptions struct{}

func (v *users) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterUserOptions) error {
	return nil
}

// DropUserOptions contains options for dropping a user.
type DropUserOptions struct{}

func (v *users) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropUserOptions) error {
	return nil
}

// UserDetails contains details about a user.
type UserDetails struct{}

func (v *users) Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error) {
	return nil, nil
}

// ShowUserOptions contains options for listing users.
type ShowUserOptions struct{}

func (v *users) Show(ctx context.Context, opts *ShowUserOptions) ([]*User, error) {
	return nil, nil
}

func (v *users) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	return nil, nil
}
