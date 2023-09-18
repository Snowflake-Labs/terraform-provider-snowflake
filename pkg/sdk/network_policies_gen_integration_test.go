package sdk

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt_NetworkPolicies(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	defaultCreateRequest := func() *CreateNetworkPolicyRequest {
		id := randomAccountObjectIdentifier(t)
		comment := "some_comment"
		allowedIP := NewIPRequest("123.0.0.1")
		blockedIP := NewIPRequest("321.0.0.1")
		blockedIP2 := NewIPRequest("124.0.0.1")
		return NewCreateNetworkPolicyRequest(id).
			WithOrReplace(Bool(true)).
			WithAllowedIpList([]IPRequest{*allowedIP}).
			WithBlockedIpList([]IPRequest{*blockedIP, *blockedIP2}).
			WithComment(&comment)
	}

	t.Run("Create", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, req.name.Name(), nps[0].Name)
		assert.Equal(t, req.Comment, nps[0].Comment)
		assert.Equal(t, len(req.AllowedIpList), nps[0].EntriesInAllowedIpList)
		assert.Equal(t, len(req.BlockedIpList), nps[0].EntriesInBlockedIpList)
	})

	t.Run("Alter - set allowed ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).
			WithSet(NewNetworkPolicySetRequest().WithAllowedIpList([]IPRequest{})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, 0, nps[0].EntriesInAllowedIpList)
	})

	t.Run("Alter - set blocked ip list", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).
			WithSet(NewNetworkPolicySetRequest().WithBlockedIpList([]IPRequest{})))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, 0, nps[0].EntriesInBlockedIpList)
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
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, alteredComment, nps[0].Comment)
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
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, 0, len(nps[0].Comment))
	})

	t.Run("Alter - rename", func(t *testing.T) {
		req := defaultCreateRequest()
		err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		require.NoError(t, err)
		t.Cleanup(dropNetworkPolicy)

		newID := randomAccountObjectIdentifier(t)
		err = client.NetworkPolicies.Alter(ctx, NewAlterNetworkPolicyRequest(req.name).WithRenameTo(&newID))
		require.NoError(t, err)

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		require.NoError(t, err)
		assert.Equal(t, 1, len(nps)) // TODO Can be more - filter and see if contains
		assert.Equal(t, newID.Name(), nps[0].Name)
	})

	t.Run("Describe", func(t *testing.T) {
		//req := defaultCreateRequest()
		//err, dropNetworkPolicy := createNetworkPolicy(t, client, req)
		//require.NoError(t, err)
		//t.Cleanup(dropNetworkPolicy)
		//
		//desc, err := client.NetworkPolicies.Describe(ctx, req.name)
		//require.NoError(t, err)
		//
		//assert.Equal(t, 3, len(desc))
	})
}
