package sdk

import "testing"

func TestApplicationRoles_Grant(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	// Minimal valid GrantApplicationRoleOptions
	defaultOpts := func() *GrantApplicationRoleOptions {
		return &GrantApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *GrantApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [GrantApplicationRoleOptions.To.RoleName GrantApplicationRoleOptions.To.ApplicationRoleName GrantApplicationRoleOptions.To.ApplicationName] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantApplicationRoleOptions.To", "RoleName", "ApplicationRoleName", "ApplicationName"))
	})

	t.Run("grant to role", func(t *testing.T) {
		roleId := randomAccountObjectIdentifier()

		opts := defaultOpts()
		opts.To = KindOfRole{
			RoleName: &roleId,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO ROLE %s`, id.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("grant to application role", func(t *testing.T) {
		appRoleId := randomDatabaseObjectIdentifier()

		opts := defaultOpts()
		opts.To = KindOfRole{
			ApplicationRoleName: &appRoleId,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO APPLICATION ROLE %s`, id.FullyQualifiedName(), appRoleId.FullyQualifiedName())
	})

	t.Run("grant to application", func(t *testing.T) {
		appId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.To = KindOfRole{
			ApplicationName: &appId,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLICATION ROLE %s TO APPLICATION %s`, id.FullyQualifiedName(), appId.FullyQualifiedName())
	})
}

func TestApplicationRoles_Revoke(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	// Minimal valid RevokeApplicationRoleOptions
	defaultOpts := func() *RevokeApplicationRoleOptions {
		return &RevokeApplicationRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RevokeApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("revoke from role", func(t *testing.T) {
		roleId := randomAccountObjectIdentifier()

		opts := defaultOpts()
		opts.From = KindOfRole{
			RoleName: &roleId,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("revoke from application role", func(t *testing.T) {
		appRoleId := randomDatabaseObjectIdentifier()

		opts := defaultOpts()
		opts.From = KindOfRole{
			ApplicationRoleName: &appRoleId,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM APPLICATION ROLE %s`, id.FullyQualifiedName(), appRoleId.FullyQualifiedName())
	})

	t.Run("revoke from application", func(t *testing.T) {
		appId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.From = KindOfRole{
			ApplicationName: &appId,
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE APPLICATION ROLE %s FROM APPLICATION %s`, id.FullyQualifiedName(), appId.FullyQualifiedName())
	})
}

func TestApplicationRoles_Show(t *testing.T) {
	appId := randomAccountObjectIdentifier()

	// Minimal valid ShowApplicationRoleOptions
	defaultOpts := func() *ShowApplicationRoleOptions {
		return &ShowApplicationRoleOptions{
			ApplicationName: appId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowApplicationRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.ApplicationName]", func(t *testing.T) {
		opts := defaultOpts()
		opts.ApplicationName = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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
