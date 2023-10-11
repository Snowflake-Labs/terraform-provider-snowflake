package sdk

import (
	"context"
	"database/sql"
)

// All the contents of this file were added to be able to use them outside the sdk package (i.e. integration tests package).
// It was easier to do it that way, so that we do not include big rename changes in the first moving PR.
// For each of them we will have to decide what do we do:
// - do we expose the field/method
// - do we keep the workaround
// - do we copy the code (e.g. ExecForTests)
// - do we move the code to other place and use it from both places
// - something else.
// This will be handled in subsequent PRs, so that the main difficulty (moving) is already merged.

// ExecForTests is an exact copy of exec (that is unexported), that some integration tests/helpers were using
// TODO: remove after we have all usages covered by SDK (for now it means implementing stages, tables, and tags)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}
