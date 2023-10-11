package sdk

import (
	"context"
	"database/sql"

	"github.com/snowflakedb/gosnowflake"
)

// All the contents of this file were added to be able to use them outside the sdk package (i.e. integration tests package).
// It was easier to do it that way, so that we do not include big rename changes in the first moving PR.
// For each of them we will have to decide what do we do:
// - do we expose the field/method (e.g. ValidObjectIdentifier)
// - do we keep the workaround
// - do we copy the code (e.g. ExecForTests)
// - do we move the code to other place and use it from both places (e.g. findOne)
// - something else.
// This will be handled in subsequent PRs, so that the main difficulty (moving) is already merged.

// ExecForTests is an exact copy of exec (that is unexported), that some integration tests/helpers were using
// TODO: remove after we have all usages covered by SDK (for now it means implementing stages, tables, and tags)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	ctx = context.WithValue(ctx, snowflakeAccountLocatorContextKey, c.accountLocator)
	result, err := c.db.ExecContext(ctx, sql)
	return result, decodeDriverError(err)
}

// ValidObjectIdentifier is just a delegate to existing unexported validObjectidentifier
func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	return validObjectidentifier(objectIdentifier)
}

// GetName is just an accessor to unexported name field
func (r *CreateNetworkPolicyRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// GetName is just an accessor to unexported name field
func (s *CreateRoleRequest) GetName() AccountObjectIdentifier {
	return s.name
}

// GetName is just an accessor to unexported name field
func (r *CreateTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// GetName is just an accessor to unexported name field
func (r *CloneTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// GetColumns is just an accessor to unexported name field
func (s *CreateExternalTableRequest) GetColumns() []*ExternalTableColumnRequest {
	return s.columns
}

// GetAccountLocator is an accessor to unexported accountLocator, which is needed in some tests
func (c *Client) GetAccountLocator() string {
	return c.accountLocator
}

// GetConfig is an accessor to unexported config, which is needed in some tests
func (c *Client) GetConfig() *gosnowflake.Config {
	return c.config
}

// FindOne just delegates to our util findOne from SDK
func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	return findOne(collection, condition)
}
