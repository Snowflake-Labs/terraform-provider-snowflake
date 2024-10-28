package sdk

import "testing"

func TestConnections_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *CreateConnectionOptions {
		return &CreateConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = invalidAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.ReplicaOf]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		opts.AsReplicaOf = &AsReplicaOf{emptyExternalObjectIdentifier}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		opts.IfNotExists = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION IF NOT EXISTS %s COMMENT = 'comment'", id.FullyQualifiedName())
	})

	t.Run("as replica of", func(t *testing.T) {
		externalId := randomExternalObjectIdentifier()
		opts := defaultOpts()
		opts.name = id
		opts.AsReplicaOf = &AsReplicaOf{externalId}
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION %s AS REPLICA OF %s", id.FullyQualifiedName(), externalId.FullyQualifiedName())
	})

	t.Run("as replica of - all options", func(t *testing.T) {
		externalId := randomExternalObjectIdentifier()
		opts := defaultOpts()
		opts.name = id
		opts.IfNotExists = Bool(true)
		opts.AsReplicaOf = &AsReplicaOf{externalId}
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION IF NOT EXISTS %s AS REPLICA OF %s COMMENT = 'comment'", id.FullyQualifiedName(), externalId.FullyQualifiedName())
	})
}

func TestConnections_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *AlterConnectionOptions {
		return &AlterConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: exactly one field from [opts.EnableConnectionFailover opts.DisableConnectionFailover opts.Primary opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableConnectionFailover = &EnableConnectionFailover{}
		opts.DisableConnectionFailover = &DisableConnectionFailover{}
		opts.Primary = Bool(true)
		opts.Set = &SetConnection{}
		opts.Unset = &UnsetConnection{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterConnectionOptions", "EnableConnectionFailover", "DisableConnectionFailover", "Primary", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SetConnection{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterConnectionOptions.Set", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &UnsetConnection{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterConnectionOptions.Unset", "Comment"))
	})

	t.Run("alter enable failover to accounts", func(t *testing.T) {
		accountIdentifier := randomAccountIdentifier()
		secondAccountIdentifier := randomAccountIdentifier()
		opts := defaultOpts()
		opts.EnableConnectionFailover = &EnableConnectionFailover{
			ToAccounts: []AccountIdentifier{accountIdentifier, secondAccountIdentifier},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s ENABLE FAILOVER TO ACCOUNTS %s, %s", id.FullyQualifiedName(), accountIdentifier.FullyQualifiedName(), secondAccountIdentifier.FullyQualifiedName())
	})

	t.Run("alter disable failover to all accounts", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableConnectionFailover = &DisableConnectionFailover{}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s DISABLE FAILOVER", id.FullyQualifiedName())
	})

	t.Run("alter disable failover to accounts", func(t *testing.T) {
		accountIdentifier := randomAccountIdentifier()
		opts := defaultOpts()
		opts.DisableConnectionFailover = &DisableConnectionFailover{
			ToAccounts: &ToAccounts{[]AccountIdentifier{accountIdentifier}},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s DISABLE FAILOVER TO ACCOUNTS %s", id.FullyQualifiedName(), accountIdentifier.FullyQualifiedName())
	})

	t.Run("primary", func(t *testing.T) {
		opts := defaultOpts()
		opts.Primary = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s PRIMARY", id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SetConnection{Comment: String("test comment")}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s SET COMMENT = 'test comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &UnsetConnection{Comment: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s UNSET COMMENT", id.FullyQualifiedName())
	})
}

func TestConnections_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *DropConnectionOptions {
		return &DropConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP CONNECTION %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP CONNECTION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestConnections_Show(t *testing.T) {
	defaultOpts := func() *ShowConnectionOptions {
		return &ShowConnectionOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW CONNECTIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			String("test_connection_name"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW CONNECTIONS LIKE 'test_connection_name'")
	})
}
