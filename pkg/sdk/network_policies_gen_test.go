package sdk

import (
	"testing"
)

func TestNetworkPolicies_Create(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	allowedNetworkRule := RandomSchemaObjectIdentifier()
	blockedNetworkRule := RandomSchemaObjectIdentifier()
	// Minimal valid CreateNetworkPolicyOptions
	defaultOpts := func() *CreateNetworkPolicyOptions {
		return &CreateNetworkPolicyOptions{
			OrReplace:              Bool(true),
			name:                   id,
			AllowedIpList:          []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}},
			BlockedIpList:          []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}},
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
			Comment:                String("some_comment"),
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
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE NETWORK POLICY %s ALLOWED_NETWORK_RULE_LIST = (%s) BLOCKED_NETWORK_RULE_LIST = (%s) ALLOWED_IP_LIST = ('123.0.0.1', '321.0.0.1') BLOCKED_IP_LIST = ('123.0.0.1', '321.0.0.1') COMMENT = 'some_comment'", opts.name.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
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

	t.Run("validation: exactly one field from [opts.Set opts.UnsetComment opts.RenameTo opts.Add opts.Remove] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNetworkPolicyOptions", "Set", "UnsetComment", "RenameTo", "Add", "Remove"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AllowedIpList opts.Set.BlockedIpList opts.Set.Comment opts.Set.AllowedNetworkRuleList opts.Set.BlockedNetworkRuleList] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkPolicyOptions.Set", "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
	})

	t.Run("validation: exactly one field from [opts.Add.AllowedNetworkRuleList opts.Add.BlockedNetworkRuleList] should be present", func(t *testing.T) {
		allowedNetworkRule := RandomSchemaObjectIdentifier()
		blockedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNetworkPolicyOptions.Add", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
	})

	t.Run("validation: exactly one field from [opts.Remove.AllowedNetworkRuleList opts.Remove.BlockedNetworkRuleList] should be present", func(t *testing.T) {
		allowedNetworkRule := RandomSchemaObjectIdentifier()
		blockedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Remove = &RemoveNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNetworkPolicyOptions.Remove", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
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

	t.Run("set allowed network rule list", func(t *testing.T) {
		allowedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("set blocked network rule list", func(t *testing.T) {
		blockedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
	})

	t.Run("add allowed network rule", func(t *testing.T) {
		allowedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s ADD ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("add blocked network rule", func(t *testing.T) {
		blockedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s ADD BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
	})

	t.Run("remove allowed network rule", func(t *testing.T) {
		allowedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Remove = &RemoveNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s REMOVE ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("remove blocked network rule", func(t *testing.T) {
		blockedNetworkRule := RandomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Remove = &RemoveNetworkRule{
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s REMOVE BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
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
