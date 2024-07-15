package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalAccessIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalAccessIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalAccessIntegrationClient {
	return &ExternalAccessIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalAccessIntegrationClient) client() *sdk.Client {
	return c.context.client
}

func (c *ExternalAccessIntegrationClient) CreateExternalAccessIntegration(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = TRUE`, id.Name(), networkRuleId.Name()))
	require.NoError(t, err)
	return id, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) DropExternalAccessIntegrationFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP EXTERNAL ACCESS INTEGRATION IF EXISTS %s`, id.Name()))
		require.NoError(t, err)
	}
}
