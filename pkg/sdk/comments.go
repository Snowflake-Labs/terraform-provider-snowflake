package sdk

import (
	"context"
)

type Comments interface {
	Set(ctx context.Context, opts *SetCommentOpts) error
	SetColumn(ctx context.Context, opts *SetColumnCommentOpts) error
}

type comments struct {
	client *Client
}

var _ Comments = (*comments)(nil)

type SetCommentOpts struct {
	comment    bool             `ddl:"static" db:"COMMENT"`
	IfExists   *bool            `ddl:"keyword" db:"IF EXISTS"`
	on         bool             `ddl:"static" db:"ON"`
	ObjectType ObjectType       `ddl:"keyword"`
	ObjectName ObjectIdentifier `ddl:"identifier"`
	Value      *string          `ddl:"parameter,single_quotes,no_equals" db:"IS"`
}

func (opts *SetCommentOpts) validate() error {
	return nil
}

func (c *comments) Set(ctx context.Context, opts *SetCommentOpts) error {
	if opts == nil {
		opts = &SetCommentOpts{}
	}
	// opts.name = name
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}

type SetColumnCommentOpts struct {
	comment  bool             `ddl:"static" db:"COMMENT"`
	IfExists *bool            `ddl:"keyword" db:"IF EXISTS"`
	on       bool             `ddl:"static" db:"ON"`
	Column   ObjectIdentifier `ddl:"identifier" db:"COLUMN"`
	Value    *string          `ddl:"parameter,single_quotes,no_equals" db:"IS"`
}

func (opts *SetColumnCommentOpts) validate() error {
	return nil
}

func (c *comments) SetColumn(ctx context.Context, opts *SetColumnCommentOpts) error {
	if opts == nil {
		opts = &SetColumnCommentOpts{}
	}
	// We only want to render table.column, not the fully qualified name with database and schema.
	if v, ok := opts.Column.(TableColumnIdentifier); ok {
		opts.Column = NewSchemaIdentifier(v.tableName, v.columnName)
	}
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}
