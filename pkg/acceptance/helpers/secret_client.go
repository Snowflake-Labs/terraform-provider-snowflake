package helpers

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

func (c *SecretClient) CreateSecreteWithBasicFlow(t *testing.T, id sdk.SchemaObjectIdentifier, username, password string) (*sdk.Secret, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)

	err := c.client().CreateWithBasicAuthentication(ctx, request)
	require.NoError(t, err)

	secret, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return secret, c.CleanupSecretFunc(t, id)
}

func (c *SecretClient) CleanupSecretFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ShowByID(ctx, id)
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}

		err = c.client().Drop(ctx, sdk.NewDropSecretRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	}
}
