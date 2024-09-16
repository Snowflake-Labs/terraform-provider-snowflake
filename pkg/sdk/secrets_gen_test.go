package sdk

import "testing"

func TestSecrets_CreateWithOAuthClientCredentialsFlow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateWithOAuthClientCredentialsFlowSecretOptions {
		return &CreateWithOAuthClientCredentialsFlowSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithOAuthClientCredentialsFlowSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithOAuthClientCredentialsFlowSecretOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		integration := randomAccountObjectIdentifier()

		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.SecurityIntegration = integration
		opts.OauthScopes = []SecurityIntegrationScope{{"test"}}
		opts.Comment = String("foo")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECRET IF NOT EXISTS %s TYPE = OAUTH2 API_AUTHENTICATION = %s OAUTH_SCOPES = ('test') COMMENT = 'foo'", id.FullyQualifiedName(), integration.FullyQualifiedName())
	})
}

func TestSecrets_CreateWithOAuthAuthorizationCodeFlow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateWithOAuthAuthorizationCodeFlowSecretOptions {
		return &CreateWithOAuthAuthorizationCodeFlowSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithOAuthAuthorizationCodeFlowSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithOAuthAuthorizationCodeFlowSecretOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		integration := randomAccountObjectIdentifier()

		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OauthRefreshToken = "foo"
		opts.OauthRefreshTokenExpiryTime = "bar"
		opts.SecurityIntegration = integration
		opts.Comment = String("test")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECRET IF NOT EXISTS %s TYPE = OAUTH2 OAUTH_REFRESH_TOKEN = 'foo' OAUTH_REFRESH_TOKEN_EXPIRY_TIME = 'bar' API_AUTHENTICATION = %s COMMENT = 'test'", id.FullyQualifiedName(), integration.FullyQualifiedName())
	})
}

func TestSecrets_CreateWithBasicAuthentication(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateWithBasicAuthenticationSecretOptions {
		return &CreateWithBasicAuthenticationSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithBasicAuthenticationSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithBasicAuthenticationSecretOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Username = "foo"
		opts.Password = "bar"
		opts.Comment = String("test")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECRET IF NOT EXISTS %s TYPE = PASSWORD USERNAME = 'foo' PASSWORD = 'bar' COMMENT = 'test'", id.FullyQualifiedName())
	})
}

func TestSecrets_CreateWithGenericString(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateWithGenericStringSecretOptions {
		return &CreateWithGenericStringSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateWithGenericStringSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateWithGenericStringSecretOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.SecretString = "test"
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECRET IF NOT EXISTS %s TYPE = GENERIC_STRING SECRET_STRING = 'test'", id.FullyQualifiedName())
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

func TestSecrets_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropSecretOptions {
		return &DropSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP SECRET %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SECRET IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestSecrets_Show(t *testing.T) {
	defaultOpts := func() *ShowSecretOptions {
		return &ShowSecretOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECRETS")
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECRETS LIKE 'pattern'")
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECRETS IN ACCOUNT")
	})
}

func TestSecrets_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeSecretOptions {
		return &DescribeSecretOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSecretOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SECRET %s", id.FullyQualifiedName())
	})
}
