package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SecretClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSecretClient(context *TestClientContext, idsGenerator *IdsGenerator) *SecretClient {
	return &SecretClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SecretClient) client() sdk.Secrets {
	return c.context.client.Secrets
}

func (c *SecretClient) CreateWithBasicFlow(t *testing.T, id sdk.SchemaObjectIdentifier, username, password string) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)

	err := c.client().CreateWithBasicAuthentication(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.DropFunc(t, id)
}

func (c *SecretClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	}
}

func (c *SecretClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Secret, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
