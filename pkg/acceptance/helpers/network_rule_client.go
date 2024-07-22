package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type NetworkRuleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewNetworkRuleClient(context *TestClientContext, idsGenerator *IdsGenerator) *NetworkRuleClient {
	return &NetworkRuleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *NetworkRuleClient) client() sdk.NetworkRules {
	return c.context.client.NetworkRules
}

func (c *NetworkRuleClient) CreateNetworkRule(t *testing.T) (*sdk.NetworkRule, func()) {
	t.Helper()
	return c.CreateNetworkRuleWithRequest(t, sdk.NewCreateNetworkRuleRequest(c.ids.RandomSchemaObjectIdentifier(),
		sdk.NetworkRuleTypeHostPort,
		[]sdk.NetworkRuleValue{},
		sdk.NetworkRuleModeEgress,
	))
}

func (c *NetworkRuleClient) CreateNetworkRuleWithRequest(t *testing.T, request *sdk.CreateNetworkRuleRequest) (*sdk.NetworkRule, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	networkRule, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return networkRule, c.DropNetworkRuleFunc(t, request.GetName())
}

func (c *NetworkRuleClient) DropNetworkRuleFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
