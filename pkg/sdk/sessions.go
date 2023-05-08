package sdk

import (
	"context"
	"fmt"
)

type Sessions interface {
	// Context functions.
	UseWarehouse(ctx context.Context, warehouse AccountLevelIdentifier) error
	UseDatabase(ctx context.Context, database AccountLevelIdentifier) error
	UseSchema(ctx context.Context, schema SchemaIdentifier) error
}

var _ Sessions = (*sessions)(nil)

type sessions struct {
	client  *Client
	builder *sqlBuilder
}

func (c *sessions) UseWarehouse(ctx context.Context, warehouse AccountLevelIdentifier) error {
	sql := fmt.Sprintf(`USE WAREHOUSE %s`, warehouse.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

func (c *sessions) UseDatabase(ctx context.Context, database AccountLevelIdentifier) error {
	sql := fmt.Sprintf(`USE DATABASE %s`, database.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}

func (c *sessions) UseSchema(ctx context.Context, schema SchemaIdentifier) error {
	sql := fmt.Sprintf(`USE SCHEMA %s`, schema.FullyQualifiedName())
	_, err := c.client.exec(ctx, sql)
	return err
}
