package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_NetworkPolicies(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	allowedIP := NewIPRequest("123.0.0.1")
	blockedIP := NewIPRequest("125.0.0.1")
	blockedIP2 := NewIPRequest("124.0.0.1")
	defaultCreateRequest := func() *CreateNetworkPolicyRequest {
		id := randomAccountObjectIdentifier(t)
		comment := "some_comment"
		return NewCreateNetworkPolicyRequest(id).
			WithOrReplace(Bool(true)).
			WithAllowedIpList([]IPRequest{*allowedIP}).
			WithBlockedIpList([]IPRequest{*blockedIP, *blockedIP2}).
			WithComment(&comment)
	}

	findNetworkPolicy := func(nps []NetworkPolicy, name string) *NetworkPolicy {
		for _, v := range nps {
			if v.Name == name {
				return &v
			}
		}
		return nil
	}

	t.Run("Create", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, req.name.Name())
		require.NotNil(t, np)
		assert.Equal(t, *req.Comment, np.Comment)
		assert.Equal(t, len(req.AllowedIpList), np.EntriesInAllowedIpList)
		assert.Equal(t, len(req.BlockedIpList), np.EntriesInBlockedIpList)
	})

	t.Run("Alter - set allowed ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).
			WithSet(NewNetworkPolicySetRequest().WithAllowedIpList([]IPRequest{{IP: "123.0.0.1"}, {IP: "125.0.0.1"}})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, req.name.Name())
		assert.Equal(t, 2, np.EntriesInAllowedIpList)
	})

	t.Run("Alter - set blocked ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).
			WithSet(NewNetworkPolicySetRequest().WithBlockedIpList([]IPRequest{{IP: "123.0.0.1"}})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, req.name.Name())
		assert.Equal(t, 1, np.EntriesInBlockedIpList)
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		alteredComment := "altered_comment"
		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).
			WithSet(NewNetworkPolicySetRequest().WithComment(&alteredComment)))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, req.name.Name())
		assert.Equal(t, alteredComment, np.Comment)
	})

	t.Run("Alter - unset comment", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).WithUnsetComment(Bool(true)))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, req.name.Name())
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
		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).WithRenameTo(&newID))
		require.NoError(t, err)
		altered = true
		t.Cleanup(func() {
			if altered {
				err = client.NetworkPolicies.Drop(ctx, NewDropNetworkPolicyRequest(newID))
				require.NoError(t, err)
			}
		})

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)

		np := findNetworkPolicy(nps, newID.Name())
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

		desc, err := client.NetworkPolicies.Describe(ctx, req.name)
		require.NoError(t, err)

		assert.Equal(t, 2, len(desc))
		assert.Contains(t, desc, NetworkPolicyDescription{Name: "ALLOWED_IP_LIST", Value: allowedIP.IP})
		assert.Contains(t, desc, NetworkPolicyDescription{Name: "BLOCKED_IP_LIST", Value: fmt.Sprintf("%s,%s", blockedIP.IP, blockedIP2.IP)})
	})
}
