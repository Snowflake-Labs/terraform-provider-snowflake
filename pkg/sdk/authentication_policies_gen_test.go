package sdk

import "testing"

func TestAuthenticationPolicies_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateAuthenticationPolicyOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.AuthenticationMethods = []AuthenticationMethods{
			{Method: AuthenticationMethodsAll},
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('ALL') COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.AuthenticationMethods = []AuthenticationMethods{
			{Method: AuthenticationMethodsSaml},
			{Method: AuthenticationMethodsPassword},
		}
		opts.MfaAuthenticationMethods = []MfaAuthenticationMethods{
			{Method: MfaAuthenticationMethodsPassword},
		}
		opts.MfaEnrollment = Pointer(MfaEnrollmentOptional)
		opts.ClientTypes = []ClientTypes{
			{ClientType: ClientTypesDrivers},
			{ClientType: ClientTypesSnowSql},
		}
		opts.SecurityIntegrations = []SecurityIntegrationsOption{
			{Name: NewAccountObjectIdentifier("security_integration")},
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('SAML', 'PASSWORD') MFA_AUTHENTICATION_METHODS = ('PASSWORD') MFA_ENROLLMENT = OPTIONAL CLIENT_TYPES = ('DRIVERS', 'SNOWSQL') SECURITY_INTEGRATIONS = (\"security_integration\") COMMENT = 'some comment'", id.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterAuthenticationPolicyOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RenameTo = &emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.AuthenticationMethods opts.Set.MfaAuthenticationMethods opts.Set.MfaEnrollment opts.Set.ClientTypes opts.Set.SecurityIntegrations opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set", "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ClientTypes opts.Unset.AuthenticationMethods opts.Unset.Comment opts.Unset.SecurityIntegrations opts.Unset.MfaAuthenticationMethods opts.Unset.MfaEnrollment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &AuthenticationPolicyUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Unset", "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment"))
	})

	t.Run("alter: set basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods: []AuthenticationMethods{
				{Method: AuthenticationMethodsSaml},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s SET AUTHENTICATION_METHODS = ('SAML')", id.FullyQualifiedName())
	})

	t.Run("alter: set all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods: []AuthenticationMethods{
				{Method: AuthenticationMethodsSaml},
			},
			MfaAuthenticationMethods: []MfaAuthenticationMethods{
				{Method: MfaAuthenticationMethodsPassword},
			},
			MfaEnrollment: Pointer(MfaEnrollmentOptional),
			ClientTypes: []ClientTypes{
				{ClientType: ClientTypesDrivers},
				{ClientType: ClientTypesSnowSql},
			},
			SecurityIntegrations: []SecurityIntegrationsOption{{Name: NewAccountObjectIdentifier("security_integration")}},
			Comment:              String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY IF EXISTS %s SET AUTHENTICATION_METHODS = ('SAML') MFA_AUTHENTICATION_METHODS = ('PASSWORD') MFA_ENROLLMENT = OPTIONAL CLIENT_TYPES = ('DRIVERS', 'SNOWSQL') SECURITY_INTEGRATIONS = (\"security_integration\") COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter: unset basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s UNSET CLIENT_TYPES", id.FullyQualifiedName())
	})

	t.Run("alter: unset all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes:              Bool(true),
			AuthenticationMethods:    Bool(true),
			SecurityIntegrations:     Bool(true),
			MfaAuthenticationMethods: Bool(true),
			MfaEnrollment:            Bool(true),
			Comment:                  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY IF EXISTS %s UNSET CLIENT_TYPES, AUTHENTICATION_METHODS, SECURITY_INTEGRATIONS, MFA_AUTHENTICATION_METHODS, MFA_ENROLLMENT, COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter: renameTo", func(t *testing.T) {
		opts := defaultOpts()
		target := randomSchemaObjectIdentifier()
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s RENAME TO %s", id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP AUTHENTICATION POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Show(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		assertOptsValidAndSQLEquals(t, opts, "SHOW AUTHENTICATION POLICIES")
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("like-pattern"),
		}
		opts.StartsWith = String("starts-with-pattern")
		opts.In = &In{
			Schema: id.SchemaId(),
		}
		opts.Limit = &LimitFrom{
			Rows: Int(10),
			From: String("limit-from"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW AUTHENTICATION POLICIES LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'", id.SchemaId().FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE AUTHENTICATION POLICY %s", id.FullyQualifiedName())
	})
}
