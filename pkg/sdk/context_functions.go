package sdk

import (
	"context"
	"database/sql"
)

type ContextFunctions interface {
	CurrentSession(ctx context.Context) (string, error)
	CurrentDatabase(ctx context.Context) (string, error)
	CurrentSchema(ctx context.Context) (string, error)
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
