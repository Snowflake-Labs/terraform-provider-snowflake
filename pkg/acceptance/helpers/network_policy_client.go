package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
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

func (c *NetworkPolicyClient) CreateNetworkPolicy(t *testing.T) (*sdk.NetworkPolicy, func()) {
	t.Helper()
	return c.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(sdk.RandomAccountObjectIdentifier()))
}

func (c *NetworkPolicyClient) CreateNetworkPolicyWithRequest(t *testing.T, request *sdk.CreateNetworkPolicyRequest) (*sdk.NetworkPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	networkPolicy, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return networkPolicy, c.DropNetworkPolicyFunc(t, request.GetName())
}

func (c *NetworkPolicyClient) DropNetworkPolicyFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropNetworkPolicyRequest(id))
		require.NoError(t, err)
	}
}
