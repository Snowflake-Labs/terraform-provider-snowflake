package sdk

import "testing"

func TestDatabaseRoleCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *CreateDatabaseRoleOptions {
		return &CreateDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateDatabaseRoleOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := setUpOpts()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE IF NOT EXISTS %s COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleAlter(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

	setUpOpts := func() *AlterDatabaseRoleOptions {
		return &AlterDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterDatabaseRoleOptions = nil
		assertOptsInvalid(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := setUpOpts()
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsInvalid(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := setUpOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: NewDatabaseObjectIdentifier("", ""),
		}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := NewDatabaseObjectIdentifier(id.DatabaseName()+randomStringN(t, 1), randomStringN(t, 12))

		opts := setUpOpts()
		opts.Rename = &DatabaseRoleRename{
			Name: newId,
		}
		assertOptsInvalid(t, opts, errDifferentDatabase)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := setUpOpts()
		opts.Unset = &DatabaseRoleUnset{
			Comment: false,
		}
		assertOptsInvalid(t, opts, errAlterNeedsAtLeastOneProperty)
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
