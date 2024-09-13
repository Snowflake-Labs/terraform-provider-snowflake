package sdk

import "testing"

func TestSecrets_CreateWithOAuthClientCredentialsFlow(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateWithOAuthClientCredentialsFlowSecretOptions
	defaultOpts := func() *CreateWithOAuthClientCredentialsFlowSecretOptions {
		return &CreateWithOAuthClientCredentialsFlowSecretOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithOAuthClientCredentialsFlowSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithOAuthClientCredentialsFlowSecretOptions", "OrReplace", "IfNotExists"))
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

func TestSecrets_CreateWithOAuthAuthorizationCodeFlow(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateWithOAuthAuthorizationCodeFlowSecretOptions
	defaultOpts := func() *CreateWithOAuthAuthorizationCodeFlowSecretOptions {
		return &CreateWithOAuthAuthorizationCodeFlowSecretOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithOAuthAuthorizationCodeFlowSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithOAuthAuthorizationCodeFlowSecretOptions", "OrReplace", "IfNotExists"))
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

func TestSecrets_CreateWithBasicAuthentication(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateWithBasicAuthenticationSecretOptions
	defaultOpts := func() *CreateWithBasicAuthenticationSecretOptions {
		return &CreateWithBasicAuthenticationSecretOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithBasicAuthenticationSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithBasicAuthenticationSecretOptions", "OrReplace", "IfNotExists"))
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

func TestSecrets_CreateWithGenericString(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateWithGenericStringSecretOptions
	defaultOpts := func() *CreateWithGenericStringSecretOptions {
		return &CreateWithGenericStringSecretOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithGenericStringSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithGenericStringSecretOptions", "OrReplace", "IfNotExists"))
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

func TestSecrets_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterSecretOptions {
		return &AlterSecretOptions{
			name: id,
		}
	}

	setOpts := func() *AlterSecretOptions {
		return &AlterSecretOptions{
			name:     id,
			Set:      &SecretSet{},
			IfExists: Bool(true),
		}
	}

	unsetOpts := func() *AlterSecretOptions {
		return &AlterSecretOptions{
			name:     id,
			Unset:    &SecretUnset{},
			IfExists: Bool(true),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSecretOptions", "Set", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.Set.SetForOAuthClientCredentialsFlow opts.Set.SetForOAuthAuthorizationFlow opts.Set.SetForBasicAuthentication opts.Set.SetForGenericString] should be present", func(t *testing.T) {
		opts := setOpts()
		opts.Set.SetForOAuthAuthorizationFlow = &SetForOAuthAuthorizationFlow{}
		opts.Set.SetForBasicAuthentication = &SetForBasicAuthentication{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSecretOptions.Set", "SetForOAuthClientCredentialsFlow", "SetForOAuthAuthorizationFlow", "SetForBasicAuthentication", "SetForGenericString"))
	})

	t.Run("alter: set options for Oauth Client Credentials Flow", func(t *testing.T) {
		opts := setOpts()
		opts.Set.Comment = String("test")
		opts.Set.SetForOAuthClientCredentialsFlow = &SetForOAuthClientCredentialsFlow{[]SecurityIntegrationScope{{"sample_scope"}}}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET IF EXISTS %s SET COMMENT = 'test' OAUTH_SCOPES = ('sample_scope')", id.FullyQualifiedName())
	})

	t.Run("alter: set options for Oauth Authorization Flow", func(t *testing.T) {
		opts := setOpts()
		opts.Set.Comment = String("test")
		opts.Set.SetForOAuthAuthorizationFlow = &SetForOAuthAuthorizationFlow{
			String("test_token"),
			String("2024-11-11"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET IF EXISTS %s SET COMMENT = 'test' OAUTH_REFRESH_TOKEN = 'test_token' OAUTH_REFRESH_TOKEN_EXPIRY_TIME = '2024-11-11'", id.FullyQualifiedName())
	})

	t.Run("alter: set options for Basic Authentication", func(t *testing.T) {
		opts := setOpts()
		opts.Set.Comment = String("test")
		opts.Set.SetForBasicAuthentication = &SetForBasicAuthentication{
			Username: String("foo"),
			Password: String("bar"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET IF EXISTS %s SET COMMENT = 'test' USERNAME = 'foo' PASSWORD = 'bar'", id.FullyQualifiedName())
	})

	t.Run("alter: set options for Generic string", func(t *testing.T) {
		opts := setOpts()
		opts.Set.Comment = String("test")
		opts.Set.SetForGenericString = &SetForGenericString{String("test")}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET IF EXISTS %s SET COMMENT = 'test' SECRET_STRING = 'test'", id.FullyQualifiedName())
	})

	t.Run("alter: unset options", func(t *testing.T) {
		opts := unsetOpts()
		opts.Unset.Comment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET IF EXISTS %s SET COMMENT = NULL", id.FullyQualifiedName())
	})
}
