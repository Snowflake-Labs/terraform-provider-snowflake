package sdk

import "testing"

func TestNetworkRules_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid CreateNetworkRuleOptions
	defaultOpts := func() *CreateNetworkRuleOptions {
		return &CreateNetworkRuleOptions{
			name: id,
			Type: NetworkRuleTypeIpv4,
			ValueList: []NetworkRuleValue{
				{Value: "0.0.0.0"},
				{Value: "1.1.1.1"},
			},
			Mode: NetworkRuleModeIngress,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNetworkRuleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE NETWORK RULE %s TYPE = IPV4 VALUE_LIST = ('0.0.0.0', '1.1.1.1') MODE = INGRESS`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE NETWORK RULE %s TYPE = IPV4 VALUE_LIST = ('0.0.0.0', '1.1.1.1') MODE = INGRESS COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestNetworkRules_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid AlterNetworkRuleOptions
	defaultOpts := func() *AlterNetworkRuleOptions {
		return &AlterNetworkRuleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterNetworkRuleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set opts.Unset] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkRuleOptions", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.ValueList opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkRuleSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkRuleOptions.Set", "ValueList", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ValueList opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &NetworkRuleUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterNetworkRuleOptions.Unset", "ValueList", "Comment"))
	})

	t.Run("all options set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NetworkRuleSet{
			ValueList: []NetworkRuleValue{
				{Value: "0.0.0.0"},
				{Value: "1.1.1.1"},
			},
			Comment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER NETWORK RULE %s SET VALUE_LIST = ('0.0.0.0', '1.1.1.1'), COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("all options unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &NetworkRuleUnset{
			ValueList: Bool(true),
			Comment:   Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER NETWORK RULE %s UNSET VALUE_LIST, COMMENT`, id.FullyQualifiedName())
	})
}

func TestNetworkRules_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DropNetworkRuleOptions
	defaultOpts := func() *DropNetworkRuleOptions {
		return &DropNetworkRuleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNetworkRuleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP NETWORK RULE %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP NETWORK RULE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestNetworkRules_Show(t *testing.T) {
	// Minimal valid ShowNetworkRuleOptions
	defaultOpts := func() *ShowNetworkRuleOptions {
		return &ShowNetworkRuleOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNetworkRuleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW NETWORK RULES`)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("name"),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database-name"),
		}
		opts.StartsWith = String("abc")
		opts.Limit = &LimitFrom{
			Rows: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW NETWORK RULES LIKE 'name' IN DATABASE "database-name" STARTS WITH 'abc' LIMIT 10`)
	})
}

func TestNetworkRules_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	// Minimal valid DescribeNetworkRuleOptions
	defaultOpts := func() *DescribeNetworkRuleOptions {
		return &DescribeNetworkRuleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNetworkRuleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE NETWORK RULE %s`, id.FullyQualifiedName())
	})
}
