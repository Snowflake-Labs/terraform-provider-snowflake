package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseRoleCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *createDatabaseRoleOptions {
		return &createDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("OrReplace", "IfNotExists"))
	})

	t.Run("validation: multiple errors", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errOneOf("OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(false)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE IF NOT EXISTS %s COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleAlter(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *alterDatabaseRoleOptions {
		return &alterDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := setUpOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: NewDatabaseObjectIdentifier("", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := NewDatabaseObjectIdentifier(id.DatabaseName()+randomStringN(t, 1), randomStringN(t, 12))

		opts := setUpOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: newId,
		}
		assertOptsInvalidJoinedErrors(t, opts, errDifferentDatabase)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &DatabaseRoleUnset{
			Comment: false,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsAtLeastOneProperty)
	})

	t.Run("rename", func(t *testing.T) {
		newId := NewDatabaseObjectIdentifier(id.DatabaseName(), randomStringN(t, 12))

		opts := setUpOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: newId,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("set comment to empty", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = ''`, id.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleDrop(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *dropDatabaseRoleOptions {
		return &dropDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestDatabaseRolesShow(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	setUpOpts := func() *showDatabaseRoleOptions {
		return &showDatabaseRoleOptions{
			Database: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeShowOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.Database = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES IN DATABASE %s`, id.FullyQualifiedName())
	})

	t.Run("show with like", func(t *testing.T) {
		opts := setUpOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s`, id.Name(), id.FullyQualifiedName())
	})
}

func TestDatabaseRoleBuilder(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	t.Run("should remove other optionals", func(t *testing.T) {
		request := NewAlterDatabaseRoleRequest(id).
			WithSet(NewDatabaseRoleSetRequest("comment")).
			WithRename(NewDatabaseRoleRenameRequest(randomDatabaseObjectIdentifier(t))).
			WithUnsetComment()

		assert.Nil(t, request.set)
		assert.Nil(t, request.rename)
		assert.NotNil(t, request.unset)
	})
}
