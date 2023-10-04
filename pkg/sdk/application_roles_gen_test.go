package sdk

import "testing"

func TestApplicationRoles_Create(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid CreateApplicationRoleOptions
	defaultOpts := func() *CreateApplicationRoleOptions {
		return &CreateApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateApplicationRoleOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE APPLICATION ROLE %s", id.Name())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE APPLICATION ROLE IF NOT EXISTS %s ", id.Name())
	})
}

func TestApplicationRoles_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid AlterApplicationRoleOptions
	defaultOpts := func() *AlterApplicationRoleOptions {
		return &AlterApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.SetComment opts.UnsetComment] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("some comment")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RenameTo", "SetComment", "UnsetComment"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		newName := NewAccountObjectIdentifier("")
		opts.RenameTo = &newName
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("rename to", func(t *testing.T) {
		opts := defaultOpts()
		newID := NewAccountObjectIdentifier("new_name")
		opts.IfExists = Bool(true)
		opts.RenameTo = &newID
		assertOptsValidAndSQLEquals(t, opts, "ALTER APPLICATION ROLE %s IF EXISTS RENAME TO %s", id.Name(), newID.Name())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "ALTER APPLICATION ROLE %s SET COMMENT = 'some comment'", id.Name())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER APPLICATION ROLE %s UNSET COMMENT", id.Name())
	})
}

func TestApplicationRoles_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid DropApplicationRoleOptions
	defaultOpts := func() *DropApplicationRoleOptions {
		return &DropApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP APPLICATION ROLE %s IF EXISTS", id.Name())
	})
}

func TestApplicationRoles_Show(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid ShowApplicationRoleOptions
	defaultOpts := func() *ShowApplicationRoleOptions {
		return &ShowApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestApplicationRoles_Grant(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid GrantApplicationRoleOptions
	defaultOpts := func() *GrantApplicationRoleOptions {
		return &GrantApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *GrantApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.GrantTo.ParentRole opts.GrantTo.ApplicationRole opts.GrantTo.Application] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.GrantTo.ParentRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.GrantTo.ApplicationRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.GrantTo.Application] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestApplicationRoles_Revoke(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	// Minimal valid RevokeApplicationRoleOptions
	defaultOpts := func() *RevokeApplicationRoleOptions {
		return &RevokeApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RevokeApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RevokeFrom.ParentRole opts.RevokeFrom.ApplicationRole opts.RevokeFrom.Application] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.ParentRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.ApplicationRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.Application] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
