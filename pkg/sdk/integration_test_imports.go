package sdk

import (
	"context"
	"database/sql"
)

// TODO: do not export this method (it was just a quick workaround for PoC)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// TODO: temporary solution to move integration tests in separate package
var ErrObjectNotExistOrAuthorized = errObjectNotExistOrAuthorized
var ErrDifferentDatabase = errDifferentDatabase

// TODO: temporary; used in integration test helper
func (r *CreateNetworkPolicyRequest) GetName() AccountObjectIdentifier {
	return r.name
}
