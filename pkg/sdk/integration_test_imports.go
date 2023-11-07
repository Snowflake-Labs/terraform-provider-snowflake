package sdk

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func ErrorsEqual(t *testing.T, expected error, actual error) {
	t.Helper()
	var expectedErr *Error
	var actualErr *Error
	if errors.As(expected, &expectedErr) && errors.As(actual, &actualErr) {
		expectedErrorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(expectedErr.Error(), "")
		errorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(actualErr.Error(), "")
		assert.Equal(t, expectedErrorWithoutFileInfo, errorWithoutFileInfo)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func ErrExactlyOneOf(structName string, fieldNames ...string) error {
	return errExactlyOneOf(structName, fieldNames...)
}

func ErrOneOf(structName string, fieldNames ...string) error {
	return errOneOf(structName, fieldNames...)
}
