package sdk

import (
	"context"
)

type SessionPolicies interface {
	// Create creates a session policy.
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateSessionPolicyOptions) error
	// Alter modifies an existing session policy
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterSessionPolicyOptions) error
	// Drop removes a session policy.
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropSessionPolicyOptions) error
	// Show returns a list of session policy.
	Show(ctx context.Context) ([]*SessionPolicy, error)
	// ShowByID returns a session policy by ID
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicy, error)
	// Describe returns the details of a session policy.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicyDetails, error)
}

var _ SessionPolicies = (*sessionPolicies)(nil)

type sessionPolicies struct {
	client *Client
}

type SessionPolicy struct {
	Name         string
	DatabaseName string
	SchemaName   string
}

type sessionPolicyRow struct {
	Name         string `db:"name"`
	DatabaseName string `db:"database_name"`
	SchemaName   string `db:"schema_name"`
}

func (row *sessionPolicyRow) toSessionPolicy() *SessionPolicy {
	return &SessionPolicy{
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
	}
}

func (v *SessionPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *SessionPolicy) ObjectType() ObjectType {
	return ObjectTypeSessionPolicy
}

// CreateSessionPolicyOptions contains options for creating a session policy.
type CreateSessionPolicyOptions struct {
	create        bool                   `ddl:"static" sql:"CREATE"`         //lint:ignore U1000 This is used in the ddl tag
	sessionPolicy bool                   `ddl:"static" sql:"SESSION POLICY"` //lint:ignore U1000 This is used in the ddl tag
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *CreateSessionPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *sessionPolicies) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateSessionPolicyOptions) error {
	if opts == nil {
		opts = &CreateSessionPolicyOptions{}
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

// AlterSessionPolicyOptions contains options for altering a session policy.
type AlterSessionPolicyOptions struct{}

func (v *sessionPolicies) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterSessionPolicyOptions) error {
	return nil
}

// DropSessionPolicyOptions contains options for dropping a session policy.
type DropSessionPolicyOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`           //lint:ignore U1000 This is used in the ddl tag
	sessionPolicy bool                   `ddl:"static" sql:"SESSION POLICY"` //lint:ignore U1000 This is used in the ddl tag
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropSessionPolicyOptions) validate() error {
	if !validObjectidentifier(opts.name) {
		return ErrInvalidObjectIdentifier
	}
	return nil
}

func (v *sessionPolicies) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropSessionPolicyOptions) error {
	if opts == nil {
		opts = &DropSessionPolicyOptions{}
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

// sessionPolicyShowOptions contains options for listing session policies.
type sessionPolicyShowOptions struct {
	show            bool `ddl:"static" sql:"SHOW"`             //lint:ignore U1000 This is used in the ddl tag
	sessionPolicies bool `ddl:"static" sql:"SESSION POLICIES"` //lint:ignore U1000 This is used in the ddl tag
}

func (opts *sessionPolicyShowOptions) validate() error {
	return nil
}

func (v *sessionPolicies) Show(ctx context.Context) ([]*SessionPolicy, error) {
	opts := &sessionPolicyShowOptions{}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []*sessionPolicyRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	sessionPolicies := make([]*SessionPolicy, 0, len(rows))
	for _, row := range rows {
		sessionPolicies = append(sessionPolicies, row.toSessionPolicy())
	}
	return sessionPolicies, nil
}

func (v *sessionPolicies) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicy, error) {
	sessionPolicies, err := v.Show(ctx)
	if err != nil {
		return nil, err
	}
	for _, sessionPolicy := range sessionPolicies {
		if sessionPolicy.Name == id.Name() {
			return sessionPolicy, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type SessionPolicyDetails struct{}

func (v *sessionPolicies) Describe(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicyDetails, error) {
	return nil, nil
}
