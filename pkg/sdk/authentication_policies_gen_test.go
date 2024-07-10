package sdk

import "testing"

func TestAuthenticationPolicies_Create(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid CreateAuthenticationPolicyOptions
	defaultOpts := func() *CreateAuthenticationPolicyOptions {
		return &CreateAuthenticationPolicyOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.AuthenticationMethods = []AuthenticationMethods{{Method: "ALL"}}
		opts.MfaAuthenticationMethods = []MfaAuthenticationMethods{{Method: "PASSWORD"}}
		opts.MfaEnrollment = String("OPTIONAL")
		opts.ClientTypes = []ClientTypes{{ClientType: "DRIVERS"}, {ClientType: "SNOWSQL"}}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('ALL') MFA_AUTHENTICATION_METHODS = ('PASSWORD') MFA_ENROLLMENT = OPTIONAL CLIENT_TYPES = ('DRIVERS', 'SNOWSQL') COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestAuthenticationPolicies_Alter(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid AlterAuthenticationPolicyOptions
	defaultOpts := func() *AlterAuthenticationPolicyOptions {
		return &AlterAuthenticationPolicyOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterAuthenticationPolicyOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.AuthenticationMethods opts.Set.MfaAuthenticationMethods opts.Set.MfaEnrollment opts.Set.ClientTypes opts.Set.SecurityIntegrations opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set", "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ClientTypes opts.Unset.AuthenticationMethods opts.Unset.Comment opts.Unset.SecurityIntegrations opts.Unset.MfaAuthenticationMethods opts.Unset.MfaEnrollment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Unset", "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment"))
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

func TestAuthenticationPolicies_Drop(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid DropAuthenticationPolicyOptions
	defaultOpts := func() *DropAuthenticationPolicyOptions {
		return &DropAuthenticationPolicyOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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

func TestAuthenticationPolicies_Show(t *testing.T) {
	// Minimal valid ShowAuthenticationPolicyOptions
	defaultOpts := func() *ShowAuthenticationPolicyOptions {
		return &ShowAuthenticationPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
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

func TestAuthenticationPolicies_Describe(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeAuthenticationPolicyOptions
	defaultOpts := func() *DescribeAuthenticationPolicyOptions {
		return &DescribeAuthenticationPolicyOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
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
