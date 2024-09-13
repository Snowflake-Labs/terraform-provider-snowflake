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
		securityIntegration := NewAccountObjectIdentifier("security_integration")
		oauthScopes := []SecurityIntegrationScope{{Scope: "sample_scope"}}

		opts := defaultOpts()
		opts.SecurityIntegration = securityIntegration
		opts.OauthScopes = oauthScopes
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECRET %s TYPE = OAUTH2 API_AUTHENTICATION = %s OAUTH_SCOPES = ('sample_scope')", id.FullyQualifiedName(), securityIntegration.FullyQualifiedName())
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
	// Minimal valid AlterSecretOptions
	defaultOpts := func() *AlterSecretOptions {
		return &AlterSecretOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSecretOptions", "Set", "Unset"))
	})

	t.Run("alter: set options", func(t *testing.T) {
		oauthScopes := &OAuthScopes{[]SecurityIntegrationScope{{Scope: "different_scope_name"}}}

		opts := defaultOpts()
		opts.Set = &SecretSet{
			OAuthScopes:                 oauthScopes,
			OauthRefreshToken:           String("refresh_token"),
			OauthRefreshTokenExpiryTime: String("2024-10-10"),
		}
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECRET %s SET OAUTH_SCOPES = ('different_scope_name')", id.FullyQualifiedName())
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
