package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

type Sessions interface {
	// Parameters
	AlterSession(ctx context.Context, opts *AlterSessionOptions) error
	ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error)
	ShowAccountParameter(ctx context.Context, parameter AccountParameter) (*Parameter, error)
	ShowSessionParameter(ctx context.Context, parameter SessionParameter) (*Parameter, error)
	ShowUserParameter(ctx context.Context, parameter UserParameter, user AccountObjectIdentifier) (*Parameter, error)
	ShowObjectParameter(ctx context.Context, parameter ObjectParameter, objectType ObjectType, objectID Identifier) (*Parameter, error)

	// Context
	UseWarehouse(ctx context.Context, warehouse AccountObjectIdentifier) error
	UseDatabase(ctx context.Context, database AccountObjectIdentifier) error
	UseSchema(ctx context.Context, schema SchemaIdentifier) error
}

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client *Client
}

type AlterSessionOptions struct {
	alter   bool          `ddl:"static" sql:"ALTER"`   //lint:ignore U1000 This is used in the ddl tag
	session bool          `ddl:"static" sql:"SESSION"` //lint:ignore U1000 This is used in the ddl tag
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

type ShowParametersOptions struct {
	show       bool          `ddl:"static" sql:"SHOW"`       //lint:ignore U1000 This is used in the ddl tag
	parameters bool          `ddl:"static" sql:"PARAMETERS"` //lint:ignore U1000 This is used in the ddl tag
	Like       *Like         `ddl:"keyword" sql:"LIKE"`
	In         *ParametersIn `ddl:"keyword" sql:"IN"`
}

func (opts *ShowParametersOptions) validate() error {
	if valueSet(opts.In) {
		if err := opts.In.validate(); err != nil {
			return err
		}
	}
	return nil
}

type ParametersIn struct {
	Session   *bool                   `ddl:"keyword" sql:"SESSION"`
	Account   *bool                   `ddl:"keyword" sql:"ACCOUNT"`
	User      AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	Warehouse AccountObjectIdentifier `ddl:"identifier" sql:"WAREHOUSE"`
	Database  AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Schema    SchemaIdentifier        `ddl:"identifier" sql:"SCHEMA"`
	Task      SchemaObjectIdentifier  `ddl:"identifier" sql:"TASK"`
	Table     SchemaObjectIdentifier  `ddl:"identifier" sql:"TABLE"`
}

func (v *ParametersIn) validate() error {
	if ok := anyValueSet(v.Session, v.Account, v.User, v.Warehouse, v.Database, v.Schema, v.Task, v.Table); !ok {
		return fmt.Errorf("at least one IN parameter must be set")
	}
	return nil
}

type ParameterType string

const (
	ParameterTypeAccount ParameterType = "ACCOUNT"
	ParameterTypeUser    ParameterType = "USER"
	ParameterTypeSession ParameterType = "SESSION"
	ParameterTypeObject  ParameterType = "OBJECT"
)

type Parameter struct {
	Key         string
	Value       string
	Default     string
	Level       ParameterType
	Description string
}

type parameterRow struct {
	Key         sql.NullString `db:"key"`
	Value       sql.NullString `db:"value"`
	Default     sql.NullString `db:"default"`
	Level       sql.NullString `db:"level"`
	Description sql.NullString `db:"description"`
}

func (row *parameterRow) toParameter() *Parameter {
	return &Parameter{
		Key:         row.Key.String,
		Value:       row.Value.String,
		Default:     row.Default.String,
		Level:       ParameterType(row.Level.String),
		Description: row.Description.String,
	}
}

func (v *sessions) ShowParameters(ctx context.Context, opts *ShowParametersOptions) ([]*Parameter, error) {
	if opts == nil {
		opts = &ShowParametersOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	rows := []parameterRow{}
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	parameters := make([]*Parameter, len(rows))
	for i, row := range rows {
		parameters[i] = row.toParameter()
	}
	return parameters, nil
}

func (v *sessions) ShowAccountParameter(ctx context.Context, parameter AccountParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Account: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *sessions) ShowSessionParameter(ctx context.Context, parameter SessionParameter) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			Session: Bool(true),
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *sessions) ShowUserParameter(ctx context.Context, parameter UserParameter, user AccountObjectIdentifier) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(parameter)),
		},
		In: &ParametersIn{
			User: user,
		},
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", parameter)
	}
	return parameters[0], nil
}

func (v *sessions) ShowObjectParameter(ctx context.Context, key ObjectParameter, objectType ObjectType, objectID Identifier) (*Parameter, error) {
	opts := &ShowParametersOptions{
		Like: &Like{
			Pattern: String(string(key)),
		},
		In: &ParametersIn{},
	}
	switch objectType {
	case ObjectTypeWarehouse:
		opts.In.Warehouse = objectID.(AccountObjectIdentifier)
	case ObjectTypeDatabase:
		opts.In.Database = objectID.(AccountObjectIdentifier)
	case ObjectTypeSchema:
		opts.In.Schema = objectID.(SchemaIdentifier)
	case ObjectTypeTask:
		opts.In.Task = objectID.(SchemaObjectIdentifier)
	case ObjectTypeTable:
		opts.In.Table = objectID.(SchemaObjectIdentifier)
	default:
		return nil, fmt.Errorf("unsupported object type %s", objectType)
	}
	parameters, err := v.ShowParameters(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(parameters) == 0 {
		return nil, fmt.Errorf("parameter %s not found", key)
	}
	return parameters[0], nil
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

func (v *sessions) UseSchema(ctx context.Context, schema SchemaIdentifier) error {
	sql := fmt.Sprintf(`USE SCHEMA %s`, schema.FullyQualifiedName())
	_, err := v.client.exec(ctx, sql)
	return err
}
