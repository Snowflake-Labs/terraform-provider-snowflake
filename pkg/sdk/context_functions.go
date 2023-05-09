package sdk

import (
	"context"
	"database/sql"
)

type ContextFunctions interface {
	// Session functions.
	CurrentSession(ctx context.Context) (string, error)

	// Session Object functions.
	CurrentDatabase(ctx context.Context) (string, error)
	CurrentSchema(ctx context.Context) (string, error)
	CurrentWarehouse(ctx context.Context) (string, error)
}

var _ ContextFunctions = (*contextFunctions)(nil)

type contextFunctions struct {
	client  *Client
	builder *sqlBuilder
}

func (c *contextFunctions) CurrentSession(ctx context.Context) (string, error) {
	s := &struct {
		CurrentSession string `db:"CURRENT_SESSION"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SESSION() as CURRENT_SESSION")
	if err != nil {
		return "", err
	}
	return s.CurrentSession, nil
}

func (c *contextFunctions) CurrentDatabase(ctx context.Context) (string, error) {
	s := &struct {
		CurrentDatabase sql.NullString `db:"CURRENT_DATABASE"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_DATABASE() as CURRENT_DATABASE")
	if err != nil {
		return "", err
	}
	if !s.CurrentDatabase.Valid {
		return "", nil
	}
	return s.CurrentDatabase.String, nil
}

func (c *contextFunctions) CurrentSchema(ctx context.Context) (string, error) {
	s := &struct {
		CurrentSchema sql.NullString `db:"CURRENT_SCHEMA"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_SCHEMA() as CURRENT_SCHEMA")
	if err != nil {
		return "", err
	}
	if !s.CurrentSchema.Valid {
		return "", nil
	}
	return s.CurrentSchema.String, nil
}

func (c *contextFunctions) CurrentWarehouse(ctx context.Context) (string, error) {
	s := &struct {
		CurrentWarehouse sql.NullString `db:"CURRENT_WAREHOUSE"`
	}{}
	err := c.client.queryOne(ctx, s, "SELECT CURRENT_WAREHOUSE() as CURRENT_WAREHOUSE")
	if err != nil {
		return "", err
	}
	if !s.CurrentWarehouse.Valid {
		return "", nil
	}
	return s.CurrentWarehouse.String, nil
}
