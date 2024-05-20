package sdk

import (
	"testing"
)

func TestNetworkPolicies_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()

	allowedNetworkRule := randomSchemaObjectIdentifier()
	blockedNetworkRule := randomSchemaObjectIdentifier()
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
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE NETWORK POLICY %s ALLOWED_NETWORK_RULE_LIST = (%s) BLOCKED_NETWORK_RULE_LIST = (%s) ALLOWED_IP_LIST = ('123.0.0.1', '321.0.0.1') BLOCKED_IP_LIST = ('123.0.0.1', '321.0.0.1') COMMENT = 'some_comment'", opts.name.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
	})
}

func TestNetworkPolicies_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()

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
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo opts.Add opts.Remove] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNetworkPolicyOptions", "Set", "Unset", "RenameTo", "Add", "Remove"))
	})

	t.Run("validation: at least one of the fields [opts.Set.AllowedIpList opts.Set.BlockedIpList opts.Set.Comment opts.Set.AllowedNetworkRuleList opts.Set.BlockedNetworkRuleList] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkPolicyOptions.Set", "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.AllowedIpList opts.Unset.BlockedIpList opts.Unset.Comment opts.Unset.AllowedNetworkRuleList opts.Unset.BlockedNetworkRuleList] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &NetworkPolicyUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkPolicyOptions.Unset", "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
	})

	t.Run("validation: exactly one field from [opts.Add.AllowedNetworkRuleList opts.Add.BlockedNetworkRuleList] should be present", func(t *testing.T) {
		allowedNetworkRule := randomSchemaObjectIdentifier()
		blockedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNetworkPolicyOptions.Add", "AllowedNetworkRuleList", "BlockedNetworkRuleList"))
	})

	t.Run("validation: exactly one field from [opts.Remove.AllowedNetworkRuleList opts.Remove.BlockedNetworkRuleList] should be present", func(t *testing.T) {
		allowedNetworkRule := randomSchemaObjectIdentifier()
		blockedNetworkRule := randomSchemaObjectIdentifier()
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
			AllowedIpList: &AllowedIPList{[]IP{{"123.0.0.1"}}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName())
	})

	t.Run("set empty ip list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			AllowedIpList: &AllowedIPList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_IP_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("set blocked ip list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedIpList: &BlockedIPList{[]IP{{"123.0.0.1"}}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName())
	})

	t.Run("set empty blocked ip list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedIpList: &BlockedIPList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_IP_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("set allowed network rule list", func(t *testing.T) {
		allowedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			AllowedNetworkRuleList: &AllowedNetworkRuleList{[]SchemaObjectIdentifier{allowedNetworkRule}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("set empty allowed network rule list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			AllowedNetworkRuleList: &AllowedNetworkRuleList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_NETWORK_RULE_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("set blocked network rule list", func(t *testing.T) {
		blockedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedNetworkRuleList: &BlockedNetworkRuleList{[]SchemaObjectIdentifier{blockedNetworkRule}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
	})

	t.Run("set empty blocked network rule list", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkPolicySet{
			BlockedNetworkRuleList: &BlockedNetworkRuleList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_NETWORK_RULE_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("unset single", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &NetworkPolicyUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("unset multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &NetworkPolicyUnset{
			AllowedNetworkRuleList: Bool(true),
			BlockedNetworkRuleList: Bool(true),
			AllowedIpList:          Bool(true),
			BlockedIpList:          Bool(true),
			Comment:                Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s UNSET ALLOWED_NETWORK_RULE_LIST, BLOCKED_NETWORK_RULE_LIST, ALLOWED_IP_LIST, BLOCKED_IP_LIST, COMMENT", id.FullyQualifiedName())
	})

	t.Run("add allowed network rule", func(t *testing.T) {
		allowedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s ADD ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("add blocked network rule", func(t *testing.T) {
		blockedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Add = &AddNetworkRule{
			BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s ADD BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRule.FullyQualifiedName())
	})

	t.Run("remove allowed network rule", func(t *testing.T) {
		allowedNetworkRule := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.Remove = &RemoveNetworkRule{
			AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRule},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s REMOVE ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRule.FullyQualifiedName())
	})

	t.Run("remove blocked network rule", func(t *testing.T) {
		blockedNetworkRule := randomSchemaObjectIdentifier()
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

	t.Run("rename to", func(t *testing.T) {
		opts := defaultOpts()
		newName := randomAccountObjectIdentifier()
		opts.RenameTo = &newName
		assertOptsValidAndSQLEquals(t, opts, "ALTER NETWORK POLICY IF EXISTS %s RENAME TO %s", id.FullyQualifiedName(), newName.FullyQualifiedName())
	})
}

func TestNetworkPolicies_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

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
		opts.name = emptyAccountObjectIdentifier
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
	id := randomAccountObjectIdentifier()

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
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE NETWORK POLICY %s", id.FullyQualifiedName())
	})
}
