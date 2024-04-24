package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type NetworkPolicyClient struct {
	context *TestClientContext
}

func NewNetworkPolicyClient(context *TestClientContext) *NetworkPolicyClient {
	return &NetworkPolicyClient{
		context: context,
	}
}

func (c *NetworkPolicyClient) client() sdk.NetworkPolicies {
	return c.context.client.NetworkPolicies
}
