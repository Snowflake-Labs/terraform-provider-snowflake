package sdk

import (
	"context"
	"database/sql"

	"github.com/snowflakedb/gosnowflake"
)

// TODO: do not export this method (it was just a quick workaround to move integration tests in separate package)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// TODO: do not export this method (it was just a quick workaround to move integration tests in separate package)
func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	return validObjectidentifier(objectIdentifier)
}

// TODO: temporary solution to move integration tests in separate package
var ErrObjectNotExistOrAuthorized = errObjectNotExistOrAuthorized
var ErrDifferentDatabase = errDifferentDatabase

// TODO: temporary; used in integration test helper
func (r *CreateNetworkPolicyRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// TODO: temporary; used in integration test helper
func (s *CreateRoleRequest) GetName() AccountObjectIdentifier {
	return s.name
}

// TODO: temporary; used in integration tests
func (r *CreateTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// TODO: temporary; used in integration tests
func (r *CloneTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// TODO: temporary; used in integration tests
func (s *CreateExternalTableRequest) GetColumns() []*ExternalTableColumnRequest {
	return s.columns
}

// TODO: temporary; used in integration tests
func (c *Client) GetAccountLocator() string {
	return c.accountLocator
}

// TODO: temporary; used in integration tests
func (c *Client) GetConfig() *gosnowflake.Config {
	return c.config
}

// TODO: temporary; used in integration tests
func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	return findOne(collection, condition)
}
