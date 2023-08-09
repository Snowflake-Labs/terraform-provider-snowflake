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
