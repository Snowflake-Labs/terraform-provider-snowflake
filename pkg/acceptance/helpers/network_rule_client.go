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

func (c *NetworkRuleClient) Create(t *testing.T) *sdk.NetworkRule {
	t.Helper()
	return c.CreateWithName(t, c.ids.Alpha())
}

func (c *NetworkRuleClient) CreateWithName(t *testing.T, name string) *sdk.NetworkRule {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.NewSchemaObjectIdentifier(name))
}

func (c *NetworkRuleClient) CreateWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.NetworkRule {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = c.client().Drop(ctx, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))
	})

	networkRule, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, networkRule)

	return networkRule
}
