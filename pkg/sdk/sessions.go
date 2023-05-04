package sdk

import (
	"context"
	"fmt"
)

type Sessions interface {
	UseDatabase(ctx context.Context, database AccountObjectIdentifier) error
	UseSchema(ctx context.Context, schema SchemaIdentifier) error
}

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client  *Client
	builder *sqlBuilder
}

func (c *sessions) UseDatabase(ctx context.Context, database AccountObjectIdentifier) error {
	sql := fmt.Sprintf(`USE DATABASE %s`, database.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

func (c *sessions) UseSchema(ctx context.Context, schema SchemaIdentifier) error {
	sql := fmt.Sprintf(`USE SCHEMA %s`, schema.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}
