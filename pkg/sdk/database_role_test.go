package sdk

import (
	"testing"
)

func TestDatabaseRoleCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	defaultOpts := func() *createDatabaseRoleOptions {
		return &createDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("OrReplace", "IfNotExists"))
	})

	t.Run("validation: multiple errors", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errOneOf("OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(false)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE IF NOT EXISTS %s COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleAlter(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	defaultOpts := func() *alterDatabaseRoleOptions {
		return &alterDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: NewDatabaseObjectIdentifier("", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := NewDatabaseObjectIdentifier(id.DatabaseName()+randomStringN(t, 1), randomStringN(t, 12))

		opts := defaultOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: newId,
		}
		assertOptsInvalidJoinedErrors(t, opts, errDifferentDatabase)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseRoleUnset{
			Comment: false,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("rename", func(t *testing.T) {
		newId := NewDatabaseObjectIdentifier(id.DatabaseName(), randomStringN(t, 12))

		opts := defaultOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: newId,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("set comment to empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = ''`, id.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleDrop(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	defaultOpts := func() *dropDatabaseRoleOptions {
		return &dropDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestDatabaseRolesShow(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *showDatabaseRoleOptions {
		return &showDatabaseRoleOptions{
			Database: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeShowOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Database = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES IN DATABASE %s`, id.FullyQualifiedName())
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s`, id.Name(), id.FullyQualifiedName())
	})
}

func TestDatabaseRoles_Grant(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *grantDatabaseRoleOptions {
		return &grantDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("grant to database role", func(t *testing.T) {
		databaseRoleId := randomDatabaseObjectIdentifier(t)
		opts := setUpOpts()
		opts.Role.DatabaseRoleName = &databaseRoleId

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO ROLE %s`, id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("grant to account role", func(t *testing.T) {
		accountRoleId := randomAccountObjectIdentifier(t)
		opts := setUpOpts()
		opts.Role.AccountRoleName = &accountRoleId

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO ROLE %s`, id.FullyQualifiedName(), accountRoleId.FullyQualifiedName())
	})
}

func TestDatabaseRoles_Revoke(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *revokeDatabaseRoleOptions {
		return &revokeDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("revoke from database role", func(t *testing.T) {
		databaseRoleId := randomDatabaseObjectIdentifier(t)
		opts := setUpOpts()
		opts.Role.DatabaseRoleName = &databaseRoleId

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("revoke from account role", func(t *testing.T) {
		accountRoleId := randomAccountObjectIdentifier(t)
		opts := setUpOpts()
		opts.Role.AccountRoleName = &accountRoleId

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), accountRoleId.FullyQualifiedName())
	})
}

func TestDatabaseRoles_GrantToShare(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)
	share := randomAccountObjectIdentifier(t)

	setUpOpts := func() *grantDatabaseRoleToShareOptions {
		return &grantDatabaseRoleToShareOptions{
			name:  id,
			Share: share,
		}
	}

	t.Run("grant to share", func(t *testing.T) {
		opts := setUpOpts()

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO SHARE %s`, id.FullyQualifiedName(), share.FullyQualifiedName())
	})
}

func TestDatabaseRoles_RevokeFromShare(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)
	share := randomAccountObjectIdentifier(t)

	setUpOpts := func() *revokeDatabaseRoleFromShareOptions {
		return &revokeDatabaseRoleFromShareOptions{
			name:  id,
			Share: share,
		}
	}

	t.Run("revoke from share", func(t *testing.T) {
		opts := setUpOpts()

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM SHARE %s`, id.FullyQualifiedName(), share.FullyQualifiedName())
	})
}
