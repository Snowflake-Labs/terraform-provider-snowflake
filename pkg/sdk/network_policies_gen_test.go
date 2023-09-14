package sdk

import "testing"

func TestNetworkPolicies_Create(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *CreateNetworkPolicyOptions {
		return &CreateNetworkPolicyOptions{
			name: id,
		}
	}

	// TODO: remove me
	_ = defaultOpts()

	t.Run("validate valid identifier for [opts.name]", func(t *testing.T) {
		// TODO: fill me
	})

}

func TestNetworkPolicies_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *DropNetworkPolicyOptions {
		return &DropNetworkPolicyOptions{
			name: id,
		}
	}

	// TODO: remove me
	_ = defaultOpts()

}

func TestNetworkPolicies_Show(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *ShowNetworkPolicyOptions {
		return &ShowNetworkPolicyOptions{
			name: id,
		}
	}

	// TODO: remove me
	_ = defaultOpts()

}

func TestNetworkPolicies_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *DescribeNetworkPolicyOptions {
		return &DescribeNetworkPolicyOptions{
			name: id,
		}
	}

	// TODO: remove me
	_ = defaultOpts()

}
