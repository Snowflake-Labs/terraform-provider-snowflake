// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"context"
	"database/sql"
)

// All the contents of this file were added to be able to use them outside the sdk package (i.e. integration tests package).
// It was easier to do it that way, so that we do not include big rename changes in the first moving PR.

// ExecForTests is an exact copy of exec (that is unexported), that some integration tests/helpers were using
// TODO: remove after we have all usages covered by SDK (for now it means implementing stages, tables, and tags)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}
