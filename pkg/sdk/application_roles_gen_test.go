package sdk

import "testing"

func TestApplicationRoles_Create(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

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
		opts.name = NewDatabaseObjectIdentifier("", "")
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE APPLICATION ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE APPLICATION ROLE IF NOT EXISTS %s COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestApplicationRoles_Alter(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

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
		opts.name = NewDatabaseObjectIdentifier("", "")
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
		newName := NewDatabaseObjectIdentifier("", "")
		opts.RenameTo = &newName
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RenameTo", "SetComment", "UnsetComment"))
	})

	t.Run("rename to", func(t *testing.T) {
		opts := defaultOpts()
		newID := NewDatabaseObjectIdentifier("db_name", "new_name")
		opts.IfExists = Bool(true)
		opts.RenameTo = &newID
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION ROLE IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), newID.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION ROLE %s SET COMMENT = 'some comment'`, id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER APPLICATION ROLE %s UNSET COMMENT`, id.FullyQualifiedName())
	})
}

func TestApplicationRoles_Drop(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

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
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP APPLICATION ROLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestApplicationRoles_Show(t *testing.T) {
	appId := randomAccountObjectIdentifier(t)

	// Minimal valid ShowApplicationRoleOptions
	defaultOpts := func() *ShowApplicationRoleOptions {
		return &ShowApplicationRoleOptions{
			ApplicationName: appId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: valid identifier for [opts.ApplicationName]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ApplicationName = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: Int(123),
			From: String("some limit"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW APPLICATION ROLES IN APPLICATION %s LIMIT 123 FROM 'some limit'`, appId.FullyQualifiedName())
	})
}

func TestApplicationRoles_Grant(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

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
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.GrantTo.ParentRole opts.GrantTo.ApplicationRole opts.GrantTo.Application] should be present", func(t *testing.T) {
		opts := defaultOpts()
		parentRole := randomAccountObjectIdentifier(t)
		appRole := randomDatabaseObjectIdentifier(t)
		opts.GrantTo = ApplicationGrantOptions{
			ParentRole:      &parentRole,
			ApplicationRole: &appRole,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.GrantTo.ParentRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewAccountObjectIdentifier("")
		opts.GrantTo = ApplicationGrantOptions{
			ParentRole: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.GrantTo.ApplicationRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewDatabaseObjectIdentifier("", "")
		opts.GrantTo = ApplicationGrantOptions{
			ApplicationRole: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.GrantTo.Application] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewAccountObjectIdentifier("")
		opts.GrantTo = ApplicationGrantOptions{
			Application: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("parent role", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier(t)
		opts.GrantTo = ApplicationGrantOptions{
			ParentRole: &roleID,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO ROLE %s`, id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("application role", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomDatabaseObjectIdentifier(t)
		opts.GrantTo = ApplicationGrantOptions{
			ApplicationRole: &roleID,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO APPLICATION ROLE %s`, id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("application", func(t *testing.T) {
		opts := defaultOpts()
		appID := randomAccountObjectIdentifier(t)
		opts.GrantTo = ApplicationGrantOptions{
			Application: &appID,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO APPLICATION %s`, id.FullyQualifiedName(), appID.FullyQualifiedName())
	})
}

func TestApplicationRoles_Revoke(t *testing.T) {
	id := randomDatabaseObjectIdentifier(t)

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
		opts.name = NewDatabaseObjectIdentifier("", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RevokeFrom.ParentRole opts.RevokeFrom.ApplicationRole opts.RevokeFrom.Application] should be present", func(t *testing.T) {
		opts := defaultOpts()
		parentRole := randomAccountObjectIdentifier(t)
		appRole := randomDatabaseObjectIdentifier(t)
		opts.RevokeFrom = ApplicationGrantOptions{
			ParentRole:      &parentRole,
			ApplicationRole: &appRole,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.ParentRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewAccountObjectIdentifier("")
		opts.RevokeFrom = ApplicationGrantOptions{
			ParentRole: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.ApplicationRole] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewDatabaseObjectIdentifier("", "")
		opts.RevokeFrom = ApplicationGrantOptions{
			ApplicationRole: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("validation: valid identifier for [opts.RevokeFrom.Application] if set", func(t *testing.T) {
		opts := defaultOpts()
		invalidID := NewAccountObjectIdentifier("")
		opts.RevokeFrom = ApplicationGrantOptions{
			Application: &invalidID,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ParentRole", "ApplicationRole", "Application"))
	})

	t.Run("parent role", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier(t)
		opts.RevokeFrom = ApplicationGrantOptions{
			ParentRole: &roleID,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("application role", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomDatabaseObjectIdentifier(t)
		opts.RevokeFrom = ApplicationGrantOptions{
			ApplicationRole: &roleID,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM APPLICATION ROLE %s`, id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("application", func(t *testing.T) {
		opts := defaultOpts()
		appID := randomAccountObjectIdentifier(t)
		opts.RevokeFrom = ApplicationGrantOptions{
			Application: &appID,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM APPLICATION %s`, id.FullyQualifiedName(), appID.FullyQualifiedName())
	})
}
