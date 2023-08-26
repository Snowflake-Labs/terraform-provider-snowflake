package sdk

import (
	"context"
	"fmt"
)

var (
	_ validatable = new(AlterSessionOptions)
	_ validatable = new(ShowParametersOptions)
)

type Sessions interface {
	// Parameters
	AlterSession(ctx context.Context, opts *AlterSessionOptions) error

	// Context
	UseWarehouse(ctx context.Context, warehouse AccountObjectIdentifier) error
	UseDatabase(ctx context.Context, database AccountObjectIdentifier) error
	UseSchema(ctx context.Context, schema DatabaseObjectIdentifier) error
}

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client *Client
}

type AlterSessionOptions struct {
	alter   bool          `ddl:"static" sql:"ALTER"`
	session bool          `ddl:"static" sql:"SESSION"`
	Set     *SessionSet   `ddl:"keyword" sql:"SET"`
	Unset   *SessionUnset `ddl:"keyword" sql:"UNSET"`
}

func (opts *AlterSessionOptions) validate() error {
	if everyValueNil(opts.Set, opts.Unset) {
		return fmt.Errorf("either SET or UNSET must be set")
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			return err
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			return err
		}
	}
	return nil
}

type SessionSet struct {
	SessionParameters *SessionParameters `ddl:"list"`
}

func (v *SessionSet) validate() error {
	if err := v.SessionParameters.validate(); err != nil {
		return err
	}
	return nil
}

type SessionUnset struct {
	SessionParametersUnset *SessionParametersUnset `ddl:"list"`
}

func (v *SessionUnset) validate() error {
	if err := v.SessionParametersUnset.validate(); err != nil {
		return err
	}
	return nil
}

func (v *sessions) AlterSession(ctx context.Context, opts *AlterSessionOptions) error {
	if opts == nil {
		opts = &AlterSessionOptions{}
	}
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

// Context
func (v *sessions) UseWarehouse(ctx context.Context, warehouse AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE WAREHOUSE %s`, warehouse.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

func (v *sessions) UseDatabase(ctx context.Context, database AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE DATABASE %s`, database.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}

func (v *sessions) UseSchema(ctx context.Context, schema DatabaseObjectIdentifier) error {
	sql := fmt.Sprintf(`USE SCHEMA %s`, schema.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}
