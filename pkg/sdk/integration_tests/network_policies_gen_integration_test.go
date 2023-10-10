package sdk_integration_tests

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_NetworkPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	allowedIP := sdk.NewIPRequest("123.0.0.1")
	blockedIP := sdk.NewIPRequest("125.0.0.1")
	blockedIP2 := sdk.NewIPRequest("124.0.0.1")
	defaultCreateRequest := func() *sdk.CreateNetworkPolicyRequest {
		id := randomAccountObjectIdentifier(t)
		comment := "some_comment"
		return sdk.NewCreateNetworkPolicyRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithAllowedIpList([]sdk.IPRequest{*allowedIP}).
			WithBlockedIpList([]sdk.IPRequest{*blockedIP, *blockedIP2}).
			WithComment(&comment)
	}

	findNetworkPolicy := func(nps []sdk.NetworkPolicy, name string) (*sdk.NetworkPolicy, error) {
		return sdk.FindOne[sdk.NetworkPolicy](nps, func(t sdk.NetworkPolicy) bool {
			return t.Name == name
		})
	}

	t.Run("Create", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, *req.Comment, np.Comment)
		assert.Equal(t, len(req.AllowedIpList), np.EntriesInAllowedIpList)
		assert.Equal(t, len(req.BlockedIpList), np.EntriesInBlockedIpList)
	})

	t.Run("Alter - set allowed ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
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
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
			WithSet(sdk.NewNetworkPolicySetRequest().WithBlockedIpList([]sdk.IPRequest{{IP: "123.0.0.1"}})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, 1, np.EntriesInBlockedIpList)
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		alteredComment := "altered_comment"
		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).
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
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).WithUnsetComment(sdk.Bool(true)))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np, err := findNetworkPolicy(nps, req.GetName().Name())
		require.NoError(t, err)
		assert.Equal(t, "", np.Comment)
	})

	t.Run("Alter - rename", func(t *testing.T) {
		altered := false

		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			if !altered {
				dropNetworkPolicy()
			}
		})

		newID := randomAccountObjectIdentifier(t)
		err = client.NetworkPolicies.Alter(ctx, sdk.NewAlterNetworkPolicyRequest(req.GetName()).WithRenameTo(&newID))
		require.NoError(t, err)
		altered = true
		t.Cleanup(func() {
			if altered {
				err = client.NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(newID))
				require.NoError(t, err)
			}
		})

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
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		desc, err := client.NetworkPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)

		assert.Equal(t, 2, len(desc))
		assert.Contains(t, desc, sdk.NetworkPolicyDescription{Name: "ALLOWED_IP_LIST", Value: allowedIP.IP})
		assert.Contains(t, desc, sdk.NetworkPolicyDescription{Name: "BLOCKED_IP_LIST", Value: fmt.Sprintf("%s,%s", blockedIP.IP, blockedIP2.IP)})
	})
}
