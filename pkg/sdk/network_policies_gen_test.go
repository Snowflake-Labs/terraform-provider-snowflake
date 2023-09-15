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

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
}

func TestNetworkPolicies_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *AlterNetworkPolicyOptions {
		return &AlterNetworkPolicyOptions{
			name: id,
		}
	}

	// TODO: remove me
	_ = defaultOpts()

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set opts.UnsetComment opts.RenameTo] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("Set", "UnsetComment", "RenameTo"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.AllowedIpList opts.Set.BlockedIpList opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AllowedIpList", "BlockedIpList", "Comment"))
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

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
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

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})
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

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
}
