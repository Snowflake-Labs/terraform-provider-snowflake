package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SessionPolicyClient struct {
	context *TestClientContext
}

func NewSessionPolicyClient(context *TestClientContext) *SessionPolicyClient {
	return &SessionPolicyClient{
		context: context,
	}
}

func (c *SessionPolicyClient) client() sdk.SessionPolicies {
	return c.context.client.SessionPolicies
}
