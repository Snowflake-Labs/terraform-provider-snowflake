package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_NetworkPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	allowedIP := sdk.NewIPRequest("123.0.0.1")
	blockedIP := sdk.NewIPRequest("125.0.0.1")
	blockedIP2 := sdk.NewIPRequest("124.0.0.1")

	databaseTest, schemaTest := testDb(t), testSchema(t)
	createNetworkRuleHandle := func(t *testing.T, client *sdk.Client) sdk.SchemaObjectIdentifier {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.AlphaN(4))
		err := client.NetworkRules.Create(ctx, sdk.NewCreateNetworkRuleRequest(id, sdk.NetworkRuleTypeIpv4, []sdk.NetworkRuleValue{}, sdk.NetworkRuleModeIngress))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.NetworkRules.Drop(ctx, sdk.NewDropNetworkRuleRequest(id))
			require.NoError(t, err)
		})
		return id
	}

	defaultCreateRequest := func() *sdk.CreateNetworkPolicyRequest {
		id := sdk.RandomAccountObjectIdentifier()
		comment := "some_comment"
		return sdk.NewCreateNetworkPolicyRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithAllowedIpList([]sdk.IPRequest{*allowedIP}).
			WithBlockedIpList([]sdk.IPRequest{*blockedIP, *blockedIP2}).
			WithComment(&comment)
	}

	findNetworkPolicy := func(nps []sdk.NetworkPolicy, name string) (*sdk.NetworkPolicy, error) {
		return collections.FindOne[sdk.NetworkPolicy](nps, func(t sdk.NetworkPolicy) bool {
			return t.Name == name
		})
	}

	t.Run("Create", func(t *testing.T) {
		req := defaultCreateRequest()
		allowedNetworkRule := createNetworkRuleHandle(t, client)
		blockedNetworkRule := createNetworkRuleHandle(t, client)
		req = req.WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{allowedNetworkRule})
		req = req.WithBlockedNetworkRuleList([]sdk.SchemaObjectIdentifier{blockedNetworkRule})

		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, *req.Comment, np.Comment)
		assert.Equal(t, len(req.AllowedIpList), np.EntriesInAllowedIpList)
		assert.Equal(t, len(req.BlockedIpList), np.EntriesInBlockedIpList)
		assert.Equal(t, len(req.AllowedNetworkRuleList), np.EntriesInAllowedNetworkRules)
		assert.Equal(t, len(req.BlockedNetworkRuleList), np.EntriesInBlockedNetworkRules)
	})

	t.Run("Alter - set allowed ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithAllowedIpList([]sdk.IPRequest{{IP: "123.0.0.1"}, {IP: "125.0.0.1"}})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, 2, np.EntriesInAllowedIpList)
	})

	t.Run("Alter - set blocked ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithBlockedIpList([]sdk.IPRequest{{IP: "123.0.0.1"}})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, 1, np.EntriesInBlockedIpList)
	})

	t.Run("Alter - set allowed network rule list", func(t *testing.T) {
		allowedNetworkRule := createNetworkRuleHandle(t, client)
		allowedNetworkRule2 := createNetworkRuleHandle(t, client)
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{allowedNetworkRule, allowedNetworkRule2})))
		require.NoError(t, err)

		np, err := client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 2, np.EntriesInAllowedNetworkRules)
	})

	t.Run("Alter - set blocked network rule list", func(t *testing.T) {
		blockedNetworkRule := createNetworkRuleHandle(t, client)

		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithBlockedNetworkRuleList([]sdk.SchemaObjectIdentifier{blockedNetworkRule})))
		require.NoError(t, err)

		np, err := client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 1, np.EntriesInBlockedNetworkRules)
	})

	t.Run("Alter - add and remove allowed network rule list", func(t *testing.T) {
		allowedNetworkRule := createNetworkRuleHandle(t, client)

		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithAdd(sdk.NewAddNetworkRuleRequest().WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{allowedNetworkRule})))
		require.NoError(t, err)

		np, err := client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 1, np.EntriesInAllowedNetworkRules)

		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithRemove(sdk.NewRemoveNetworkRuleRequest().WithAllowedNetworkRuleList([]sdk.SchemaObjectIdentifier{allowedNetworkRule})))
		require.NoError(t, err)

		np, err = client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 0, np.EntriesInAllowedNetworkRules)
	})

	t.Run("Alter - add and remove blocked network rule list", func(t *testing.T) {
		blockedNetworkRule := createNetworkRuleHandle(t, client)

		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithAdd(sdk.NewAddNetworkRuleRequest().WithBlockedNetworkRuleList([]sdk.SchemaObjectIdentifier{blockedNetworkRule})))
		require.NoError(t, err)

		np, err := client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 1, np.EntriesInBlockedNetworkRules)

		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithRemove(sdk.NewRemoveNetworkRuleRequest().WithBlockedNetworkRuleList([]sdk.SchemaObjectIdentifier{blockedNetworkRule})))
		require.NoError(t, err)

		np, err = client.NetworkPolicies.ShowByID(ctx, req.GetName())
		require.NoError(t, err)
		assert.Equal(t, 0, np.EntriesInAllowedNetworkRules)
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		alteredComment := "altered_comment"
		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithComment(&alteredComment)))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, alteredComment, np.Comment)
	})

	t.Run("Alter - unset comment", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).WithUnsetComment(sdk.Bool(true)))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, "", np.Comment)
	})

	t.Run("Alter - rename", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		newID := sdk.RandomAccountObjectIdentifier()
		err := client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).WithRenameTo(&newID))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().NetworkPolicy.DropNetworkPolicyFunc(t, newID))

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, newID.Name())
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), np.Name)
		assert.Equal(t, *req.Comment, np.Comment)
		assert.Equal(t, len(req.AllowedIpList), np.EntriesInAllowedIpList)
		assert.Equal(t, len(req.BlockedIpList), np.EntriesInBlockedIpList)
	})

	t.Run("Describe", func(t *testing.T) {
		req := defaultCreateRequest()
		_, dropNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, req)
		t.Cleanup(dropNetworkPolicy)

		desc, err := client.NetworkPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)

		assert.Equal(t, 2, len(desc))
		assert.Contains(t, desc, sdk.NetworkPolicyDescription{Name: "ALLOWED_IP_LIST", Value: allowedIP.IP})
		assert.Contains(t, desc, sdk.NetworkPolicyDescription{Name: "BLOCKED_IP_LIST", Value: fmt.Sprintf("%s,%s", blockedIP.IP, blockedIP2.IP)})
	})
}
