package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PackagesPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewPackagesPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *PackagesPolicyClient {
	return &PackagesPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *PackagesPolicyClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.context.client.ExecUnsafe(context.Background(), fmt.Sprintf("CREATE PACKAGES POLICY %s LANGUAGE PYTHON", id.FullyQualifiedName()))
	require.NoError(t, err)

	return id, func() {
		_, err = c.context.client.ExecUnsafe(context.Background(), fmt.Sprintf("DROP PACKAGES POLICY IF EXISTS %s", id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
