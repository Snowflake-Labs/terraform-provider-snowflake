package sdk

import "context"

func (ts *testSuite) createNetworkPolicy() (*NetworkPolicy, error) {
	options := NetworkPolicyCreateOptions{
		Name: "NETWORK_POLICY_TEST",
		AllowedIPList: []string{
			"192.168.1.0",
			"192.168.1.100",
		},
		NetworkPolicyProperties: &NetworkPolicyProperties{
			Comment:       String("test network policy"),
			BlockedIPList: StringSlice([]string{"192.168.1.99"}),
		},
	}
	return ts.client.NetworkPolicies.Create(context.Background(), options)
}

func (ts *testSuite) TestListNetworkPolicy() {
	networkPolicy, err := ts.createNetworkPolicy()
	ts.NoError(err)

	networkPolicies, err := ts.client.NetworkPolicies.List(context.Background())
	ts.NoError(err)
	ts.Greater(len(networkPolicies), 1)

	ts.NoError(ts.client.NetworkPolicies.Delete(context.Background(), networkPolicy.Name))
}

func (ts *testSuite) TestReadNetworkPolicy() {
	networkPolicy, err := ts.createNetworkPolicy()
	ts.NoError(err)

	entity, err := ts.client.NetworkPolicies.Read(context.Background(), networkPolicy.Name)
	ts.NoError(err)
	ts.Equal(networkPolicy.Name, entity.Name)
	ts.Equal(networkPolicy.AllowedIPList, entity.AllowedIPList)
	ts.Equal(networkPolicy.BlockedIPList, entity.BlockedIPList)
	ts.NoError(ts.client.NetworkPolicies.Delete(context.Background(), networkPolicy.Name))
}

func (ts *testSuite) TestUpdateNetworkPolicy() {
	networkPolicy, err := ts.createNetworkPolicy()
	ts.NoError(err)

	options := NetworkPolicyUpdateOptions{
		AllowedIPList: StringSlice([]string{"192.168.1.101"}),
		NetworkPolicyProperties: &NetworkPolicyProperties{
			BlockedIPList: StringSlice([]string{"192.168.1.100"}),
			Comment:       String("updated"),
		},
	}
	afterUpdate, err := ts.client.NetworkPolicies.Update(context.Background(), networkPolicy.Name, options)
	ts.NoError(err)
	ts.Equal(*options.NetworkPolicyProperties.Comment, afterUpdate.Comment)

	ts.NoError(ts.client.NetworkPolicies.Delete(context.Background(), networkPolicy.Name))
}

func (ts *testSuite) TestRenameNetworkPolicy() {
	networkPolicy, err := ts.createNetworkPolicy()
	ts.NoError(err)

	newNetworkPolicy := "NEW_USER_TEST"
	ts.NoError(ts.client.NetworkPolicies.Rename(context.Background(), networkPolicy.Name, newNetworkPolicy))
	ts.NoError(ts.client.NetworkPolicies.Delete(context.Background(), newNetworkPolicy))
}
