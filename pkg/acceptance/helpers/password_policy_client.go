package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type PasswordPolicyClient struct {
	context *TestClientContext
}

func NewPasswordPolicyClient(context *TestClientContext) *PasswordPolicyClient {
	return &PasswordPolicyClient{
		context: context,
	}
}

func (c *PasswordPolicyClient) client() sdk.PasswordPolicies {
	return c.context.client.PasswordPolicies
}
