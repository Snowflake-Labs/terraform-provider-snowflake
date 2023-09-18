package sdk

import "testing"

func TestNetworkPolicies_Create(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid CreateNetworkPolicyOptions
	defaultOpts := func() *CreateNetworkPolicyOptions {
		return &CreateNetworkPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestNetworkPolicies_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid AlterNetworkPolicyOptions
	defaultOpts := func() *AlterNetworkPolicyOptions {
		return &AlterNetworkPolicyOptions{
			name: id,
		}
	}

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

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestNetworkPolicies_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid DropNetworkPolicyOptions
	defaultOpts := func() *DropNetworkPolicyOptions {
		return &DropNetworkPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestNetworkPolicies_Show(t *testing.T) {
	// Minimal valid ShowNetworkPolicyOptions
	defaultOpts := func() *ShowNetworkPolicyOptions {
		return &ShowNetworkPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestNetworkPolicies_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid DescribeNetworkPolicyOptions
	defaultOpts := func() *DescribeNetworkPolicyOptions {
		return &DescribeNetworkPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
