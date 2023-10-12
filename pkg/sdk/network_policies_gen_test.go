package sdk

import (
	"testing"
)

func TestNetworkPolicies_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid CreateNetworkPolicyOptions
	defaultOpts := func() *CreateNetworkPolicyOptions {
		return &CreateNetworkPolicyOptions{
			OrReplace:     Bool(true),
			name:          id,
			AllowedIpList: []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}},
			BlockedIpList: []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}},
			Comment:       String("some_comment"),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE NETWORK POLICY %s ALLOWED_IP_LIST = ('123.0.0.1', '321.0.0.1') BLOCKED_IP_LIST = ('123.0.0.1', '321.0.0.1') COMMENT = 'some_comment'", opts.name.FullyQualifiedName())
	})
}

func TestNetworkPolicies_Alter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid AlterNetworkPolicyOptions
	defaultOpts := func() *AlterNetworkPolicyOptions {
		return &AlterNetworkPolicyOptions{
			name:     id,
			IfExists: Bool(true),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.UnsetComment opts.RenameTo] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("Set", "UnsetComment", "RenameTo"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AllowedIpList opts.Set.BlockedIpList opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AllowedIpList", "BlockedIpList", "Comment"))
	})

	t.Run("set allowed ip list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			AllowedIpList: []IP{{IP: "123.0.0.1"}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName())
	})

	t.Run("set blocked ip list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedIpList: []IP{{IP: "123.0.0.1"}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			Comment: String("some_comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET COMMENT = 'some_comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("rename to", func(t *testing.T) {
		opts := defaultOpts()
		newName := RandomAccountObjectIdentifier()
		opts.RenameTo = &newName
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s RENAME TO %s", id.FullyQualifiedName(), newName.FullyQualifiedName())
	})
}

func TestNetworkPolicies_Drop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DropNetworkPolicyOptions
	defaultOpts := func() *DropNetworkPolicyOptions {
		return &DropNetworkPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP NETWORK POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestNetworkPolicies_Show(t *testing.T) {
	// Minimal valid ShowNetworkPolicyOptions
	defaultOpts := func() *ShowNetworkPolicyOptions {
		return &ShowNetworkPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW NETWORK POLICIES")
	})
}

func TestNetworkPolicies_Describe(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	// Minimal valid DescribeNetworkPolicyOptions
	defaultOpts := func() *DescribeNetworkPolicyOptions {
		return &DescribeNetworkPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNetworkPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE NETWORK POLICY %s", id.FullyQualifiedName())
	})
}
