package sdk

import "testing"

func TestConntections_CreateConnection(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *CreateConnectionConnectionOptions {
		return &CreateConnectionConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateConnectionConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
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
}

func TestConntections_CreateReplicatedConnection(t *testing.T) {
	id := randomAccountObjectIdentifier()
	externalId := randomExternalObjectIdentifier()
	defaultOpts := func() *CreateReplicatedConnectionConnectionOptions {
		return &CreateReplicatedConnectionConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateReplicatedConnectionConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		opts.ReplicaOf = externalId
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.ReplicaOf]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		opts.ReplicaOf = emptyExtenalObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		opts.ReplicaOf = externalId
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION %s AS REPLICA OF %s", id.FullyQualifiedName(), externalId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = id
		opts.IfNotExists = Bool(true)
		opts.ReplicaOf = externalId
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE CONNECTION IF NOT EXISTS %s AS REPLICA OF %s COMMENT = 'comment'", id.FullyQualifiedName(), externalId.FullyQualifiedName())
	})
}

func TestConntections_AlterConnectionFailover(t *testing.T) {
	id := randomAccountObjectIdentifier()
	externalId := randomExternalObjectIdentifier()
	externalIdTwo := randomExternalObjectIdentifier()
	defaultOpts := func() *AlterConnectionFailoverConnectionOptions {
		return &AlterConnectionFailoverConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterConnectionFailoverConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: exactly one field from [opts.EnableConnectionFailover opts.DisableConnectionFailover opts.Primary] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableConnectionFailover = &EnableConnectionFailover{}
		opts.DisableConnectionFailover = &DisableConnectionFailover{}
		opts.Primary = &Primary{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterConnectionFailoverConnectionOptions", "EnableConnectionFailover", "DisableConnectionFailover", "Primary"))
	})

	t.Run("enable connection failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableConnectionFailover = &EnableConnectionFailover{
			Accounts: []ExternalObjectIdentifier{externalId, externalIdTwo},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s ENABLE FAILOVER TO ACCOUNTS %s, %s", id.FullyQualifiedName(), externalId.FullyQualifiedName(), externalIdTwo.FullyQualifiedName())
	})

	t.Run("enable connection failover with ignore edition check", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableConnectionFailover = &EnableConnectionFailover{
			Accounts:           []ExternalObjectIdentifier{externalId, externalIdTwo},
			IgnoreEditionCheck: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s ENABLE FAILOVER TO ACCOUNTS %s, %s IGNORE EDITION CHECK", id.FullyQualifiedName(), externalId.FullyQualifiedName(), externalIdTwo.FullyQualifiedName())
	})

	t.Run("disable connection failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableConnectionFailover = &DisableConnectionFailover{}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s DISABLE FAILOVER", id.FullyQualifiedName())
	})

	t.Run("disable connection failover to accounts", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableConnectionFailover = &DisableConnectionFailover{
			ToAccounts: Bool(true),
			Accounts:   []ExternalObjectIdentifier{externalId, externalIdTwo},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s DISABLE FAILOVER TO ACCOUNTS %s, %s", id.FullyQualifiedName(), externalId.FullyQualifiedName(), externalIdTwo.FullyQualifiedName())
	})
}

func TestConntections_AlterConnection(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *AlterConnectionConnectionOptions {
		return &AlterConnectionConnectionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterConnectionConnectionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Set{}
		opts.Unset = &Unset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterConnectionConnectionOptions", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Set{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterConnectionConnectionOptions.Set", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Unset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterConnectionConnectionOptions.Unset", "Comment"))
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Set{Comment: String("test comment")}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CONNECTION %s SET COMMENT = 'test comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Unset{Comment: Bool(true)}
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
