package sdk

import "testing"

func TestApplicationRoles_Show(t *testing.T) {
	appId := RandomAccountObjectIdentifier()

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
		opts.ApplicationName = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("ShowApplicationRoleOptions", "ApplicationRole"))
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
