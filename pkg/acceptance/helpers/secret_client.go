package helpers

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

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
