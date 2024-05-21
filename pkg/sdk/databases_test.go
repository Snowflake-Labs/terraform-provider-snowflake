package sdk

import (
	"errors"
	"testing"
	"time"
)

func TestDatabasesCreate(t *testing.T) {
	defaultOpts := func() *CreateDatabaseOptions {
		return &CreateDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid clone", func(t *testing.T) {
		opts := defaultOpts()
		opts.Clone = &Clone{
			SourceObject: NewAccountObjectIdentifier(""),
			At: &TimeTravel{
				Timestamp: Pointer(time.Now()),
				Offset:    Int(123),
			},
			Before: new(TimeTravel),
		}
		assertOptsInvalidJoinedErrors(t, opts,
			errors.New("only one of AT or BEFORE can be set"),
			errors.New("exactly one of TIMESTAMP, OFFSET or STATEMENT can be set"),
		)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = Pointer(NewAccountObjectIdentifier(""))
		opts.Catalog = Pointer(NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateDatabaseOptions", "Catalog"),
		)
	})

	t.Run("clone", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Clone = &Clone{
			SourceObject: NewAccountObjectIdentifier("db1"),
			At: &TimeTravel{
				Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s CLONE "db1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, opts.name.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Transient = Bool(true)
		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(1)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.DefaultDDLCollation = String("en_US")
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelOnEvent)
		opts.Comment = String("comment")
		opts.Tag = []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT DATABASE IF NOT EXISTS %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s DEFAULT_DDL_COLLATION = 'en_US' LOG_LEVEL = 'INFO' TRACE_LEVEL = 'ON_EVENT' COMMENT = 'comment' TAG ("db1"."schema1"."tag1" = 'v1')`, opts.name.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})
}

func TestDatabasesCreateShared(t *testing.T) {
	defaultOpts := func() *CreateSharedDatabaseOptions {
		return &CreateSharedDatabaseOptions{
			name:      randomAccountObjectIdentifier(),
			fromShare: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid from share name", func(t *testing.T) {
		opts := defaultOpts()
		opts.fromShare = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateSharedDatabaseOptions", "fromShare"))
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomAccountObjectIdentifier()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSharedDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("db")
		opts.ExternalVolume = Pointer(NewAccountObjectIdentifier(""))
		opts.Catalog = Pointer(NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateSharedDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateSharedDatabaseOptions", "Catalog"),
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Transient = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT DATABASE IF NOT EXISTS %s FROM SHARE %s`, opts.name.FullyQualifiedName(), opts.fromShare.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts.OrReplace = Bool(true)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.DefaultDDLCollation = String("en_US")
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelOnEvent)
		opts.Comment = String("comment")
		opts.Tag = []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s FROM SHARE %s EXTERNAL_VOLUME = %s CATALOG = %s DEFAULT_DDL_COLLATION = 'en_US' LOG_LEVEL = 'INFO' TRACE_LEVEL = 'ON_EVENT' COMMENT = 'comment' TAG ("db1"."schema1"."tag1" = 'v1')`, opts.name.FullyQualifiedName(), opts.fromShare.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})
}

func TestDatabasesCreateSecondary(t *testing.T) {
	defaultOpts := func() *CreateSecondaryDatabaseOptions {
		return &CreateSecondaryDatabaseOptions{
			name:            randomAccountObjectIdentifier(),
			primaryDatabase: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid primary database", func(t *testing.T) {
		opts := defaultOpts()
		opts.primaryDatabase = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "primaryDatabase"))
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSecondaryDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = Pointer(NewAccountObjectIdentifier(""))
		opts.Catalog = Pointer(NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateSecondaryDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateSecondaryDatabaseOptions", "Catalog"),
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE IF NOT EXISTS %s AS REPLICA OF %s`, opts.name.FullyQualifiedName(), opts.primaryDatabase.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		primaryDatabaseId := NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier())
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Transient = Bool(true)
		opts.primaryDatabase = primaryDatabaseId
		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(10)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.DefaultDDLCollation = String("en_US")
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelOnEvent)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TRANSIENT DATABASE %s AS REPLICA OF %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 10 EXTERNAL_VOLUME = %s CATALOG = %s DEFAULT_DDL_COLLATION = 'en_US' LOG_LEVEL = 'INFO' TRACE_LEVEL = 'ON_EVENT' COMMENT = 'comment'`, opts.name.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})
}

func TestDatabasesAlter(t *testing.T) {
	defaultOpts := func() *AlterDatabaseOptions {
		return &AlterDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			ExternalVolume: Pointer(NewAccountObjectIdentifier("")),
			Catalog:        Pointer(NewAccountObjectIdentifier("")),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "ExternalVolume"), errInvalidIdentifier("DatabaseSet", "Catalog"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
	})

	t.Run("validation: at least one set option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("DatabaseSet", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ExternalVolume", "Catalog", "DefaultDDLCollation", "LogLevel", "TraceLevel", "Comment"))
	})

	t.Run("validation: at least one unset option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("DatabaseUnset", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ExternalVolume", "Catalog", "DefaultDDLCollation", "LogLevel", "TraceLevel", "Comment"))
	})

	t.Run("validation: invalid external volume identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			ExternalVolume: Pointer(NewAccountObjectIdentifier("")),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "ExternalVolume"))
	})

	t.Run("validation: invalid catalog integration identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			Catalog: Pointer(NewAccountObjectIdentifier("")),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "Catalog"))
	})

	t.Run("validation: invalid NewName identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = Pointer(NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterDatabaseOptions", "NewName"))
	})

	t.Run("validation: invalid SwapWith identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(NewAccountObjectIdentifier(""))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterDatabaseOptions", "SwapWith"))
	})

	t.Run("rename", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.NewName = Pointer(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE IF EXISTS %s RENAME TO %s`, opts.name.FullyQualifiedName(), opts.NewName.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SWAP WITH %s`, opts.name.FullyQualifiedName(), opts.SwapWith.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			DataRetentionTimeInDays:    Int(1),
			MaxDataExtensionTimeInDays: Int(1),
			ExternalVolume:             &externalVolumeId,
			Catalog:                    &catalogId,
			DefaultDDLCollation:        String("en_US"),
			LogLevel:                   Pointer(LogLevelError),
			TraceLevel:                 Pointer(TraceLevelOnEvent),
			Comment:                    String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, EXTERNAL_VOLUME = %s, CATALOG = %s, DEFAULT_DDL_COLLATION = 'en_US', LOG_LEVEL = 'ERROR', TRACE_LEVEL = 'ON_EVENT', COMMENT = 'comment'`, opts.name.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{
			DataRetentionTimeInDays:    Bool(true),
			MaxDataExtensionTimeInDays: Bool(true),
			ExternalVolume:             Bool(true),
			Catalog:                    Bool(true),
			DefaultDDLCollation:        Bool(true),
			LogLevel:                   Bool(true),
			TraceLevel:                 Bool(true),
			Comment:                    Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, DEFAULT_DDL_COLLATION, LOG_LEVEL, TRACE_LEVEL, COMMENT`, opts.name.FullyQualifiedName())
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTag = []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
				Value: "v2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SET TAG "db"."schema"."tag1" = 'v1', "db"."schema"."tag2" = 'v2'`, opts.name.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTag = []ObjectIdentifier{
			NewSchemaObjectIdentifier("db", "schema", "tag1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET TAG "db"."schema"."tag1"`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesAlterReplication(t *testing.T) {
	defaultOpts := func() *AlterDatabaseReplicationOptions {
		return &AlterDatabaseReplicationOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{}
		opts.DisableReplication = &DisableReplication{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("enable replication", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
			IgnoreEditionCheck: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s ENABLE REPLICATION TO ACCOUNTS "account1" IGNORE EDITION CHECK`, opts.name.FullyQualifiedName())
	})

	t.Run("disable replication", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableReplication = &DisableReplication{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s DISABLE REPLICATION TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := defaultOpts()
		opts.Refresh = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s REFRESH`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesAlterFailover(t *testing.T) {
	defaultOpts := func() *AlterDatabaseFailoverOptions {
		return &AlterDatabaseFailoverOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableFailover = &EnableFailover{}
		opts.DisableFailover = &DisableFailover{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
	})

	t.Run("enable failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableFailover = &EnableFailover{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s ENABLE FAILOVER TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("disable failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableFailover = &DisableFailover{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s DISABLE FAILOVER TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("primary", func(t *testing.T) {
		opts := defaultOpts()
		opts.Primary = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s PRIMARY`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesDrop(t *testing.T) {
	defaultOpts := func() *DropDatabaseOptions {
		return &DropDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE %s`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s CASCADE`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - restrict", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Restrict = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s RESTRICT`, opts.name.FullyQualifiedName())
	})

	t.Run("validation: cascade and restrict set together", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		opts.Restrict = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DropDatabaseOptions", "Cascade", "Restrict"))
	})
}

func TestDatabasesUndrop(t *testing.T) {
	defaultOpts := func() *undropDatabaseOptions {
		return &undropDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `UNDROP DATABASE %s`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesShow(t *testing.T) {
	defaultOpts := func() *ShowDatabasesOptions {
		return &ShowDatabasesOptions{}
	}

	t.Run("without show options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES`)
	})

	t.Run("terse", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE DATABASES`)
	})

	t.Run("history", func(t *testing.T) {
		opts := defaultOpts()
		opts.History = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES HISTORY`)
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("db1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES LIKE 'db1'`)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.History = Bool(true)
		opts.Like = &Like{
			Pattern: String("db2"),
		}
		opts.LimitFrom = &LimitFrom{
			Rows: Int(1),
			From: String("db1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE DATABASES HISTORY LIKE 'db2' LIMIT 1 FROM 'db1'`)
	})
}

func TestDatabasesDescribe(t *testing.T) {
	defaultOpts := func() *describeDatabaseOptions {
		return &describeDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE DATABASE %s`, opts.name.FullyQualifiedName())
	})
}
