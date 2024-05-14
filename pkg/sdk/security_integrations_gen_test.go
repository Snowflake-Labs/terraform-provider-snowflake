package sdk

import (
	"testing"
)

func TestSecurityIntegrations_CreateSAML2(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateSAML2SecurityIntegrationOptions
	defaultOpts := func() *CreateSAML2SecurityIntegrationOptions {
		return &CreateSAML2SecurityIntegrationOptions{
			name:          id,
			Enabled:       true,
			Saml2Issuer:   "issuer",
			Saml2SsoUrl:   "url",
			Saml2Provider: "provider",
			Saml2X509Cert: "cert",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSAML2SecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSAML2SecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.AllowedEmailPatterns = []EmailPattern{{Pattern: "pattern"}}
		opts.AllowedUserDomains = []UserDomain{{Domain: "domain"}}
		opts.Comment = Pointer("a")
		opts.Saml2EnableSpInitiated = Pointer(true)
		opts.Saml2ForceAuthn = Pointer(true)
		opts.Saml2PostLogoutRedirectUrl = Pointer("redirect")
		opts.Saml2RequestedNameidFormat = Pointer("format")
		opts.Saml2SignRequest = Pointer(true)
		opts.Saml2SnowflakeAcsUrl = Pointer("acs")
		opts.Saml2SnowflakeIssuerUrl = Pointer("issuer")
		opts.Saml2SpInitiatedLoginPageLabel = Pointer("label")

		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'"+
			" ALLOWED_USER_DOMAINS = ('domain') ALLOWED_EMAIL_PATTERNS = ('pattern') SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label' SAML2_ENABLE_SP_INITIATED = true SAML2_SIGN_REQUEST = true"+
			" SAML2_REQUESTED_NAMEID_FORMAT = 'format' SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect' SAML2_FORCE_AUTHN = true SAML2_SNOWFLAKE_ISSUER_URL = 'issuer' SAML2_SNOWFLAKE_ACS_URL = 'acs'"+
			" COMMENT = 'a'", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateSCIM(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateSCIMSecurityIntegrationOptions
	defaultOpts := func() *CreateSCIMSecurityIntegrationOptions {
		return &CreateSCIMSecurityIntegrationOptions{
			name:       id,
			Enabled:    true,
			ScimClient: "GENERIC",
			RunAsRole:  "GENERIC_SCIM_PROVISIONER",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSCIMSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSCIMSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.NetworkPolicy = Pointer(networkPolicyID)
		opts.SyncPassword = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'"+
			" NETWORK_POLICY = %s SYNC_PASSWORD = true", id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterSAML2Integration(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterSAML2IntegrationSecurityIntegrationOptions
	defaultOpts := func() *AlterSAML2IntegrationSecurityIntegrationOptions {
		return &AlterSAML2IntegrationSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSAML2IntegrationSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SAML2IntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SAML2IntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSAML2IntegrationSecurityIntegrationOptions.Set", "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider",
			"Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns", "Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest",
			"Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SAML2IntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSAML2IntegrationSecurityIntegrationOptions.Unset",
			"Enabled", "Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SAML2IntegrationSet{
			Enabled:                        Pointer(true),
			Saml2Issuer:                    Pointer("issuer"),
			Saml2SsoUrl:                    Pointer("url"),
			Saml2Provider:                  Pointer("provider"),
			Saml2X509Cert:                  Pointer("cert"),
			AllowedUserDomains:             []UserDomain{{Domain: "domain"}},
			AllowedEmailPatterns:           []EmailPattern{{Pattern: "pattern"}},
			Saml2SpInitiatedLoginPageLabel: Pointer("label"),
			Saml2EnableSpInitiated:         Pointer(true),
			Saml2SnowflakeX509Cert:         Pointer("cert"),
			Saml2SignRequest:               Pointer(true),
			Saml2RequestedNameidFormat:     Pointer("format"),
			Saml2PostLogoutRedirectUrl:     Pointer("redirect"),
			Saml2ForceAuthn:                Pointer(true),
			Saml2SnowflakeIssuerUrl:        Pointer("issuer"),
			Saml2SnowflakeAcsUrl:           Pointer("acs"),
			Comment:                        Pointer("a"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'"+
			" ALLOWED_USER_DOMAINS = ('domain') ALLOWED_EMAIL_PATTERNS = ('pattern') SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label' SAML2_ENABLE_SP_INITIATED = true SAML2_SNOWFLAKE_X509_CERT = 'cert' SAML2_SIGN_REQUEST = true"+
			" SAML2_REQUESTED_NAMEID_FORMAT = 'format' SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect' SAML2_FORCE_AUTHN = true SAML2_SNOWFLAKE_ISSUER_URL = 'issuer' SAML2_SNOWFLAKE_ACS_URL = 'acs'"+
			" COMMENT = 'a'", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SAML2IntegrationUnset{
			Enabled:                    Pointer(true),
			Saml2ForceAuthn:            Pointer(true),
			Saml2RequestedNameidFormat: Pointer(true),
			Saml2PostLogoutRedirectUrl: Pointer(true),
			Comment:                    Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, SAML2_FORCE_AUTHN, SAML2_REQUESTED_NAMEID_FORMAT, SAML2_POST_LOGOUT_REDIRECT_URL, COMMENT", id.FullyQualifiedName())
	})

	t.Run("refresh SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		opts := defaultOpts()
		opts.RefreshSaml2SnowflakePrivateKey = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterSCIMIntegration(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterSCIMIntegrationSecurityIntegrationOptions
	defaultOpts := func() *AlterSCIMIntegrationSecurityIntegrationOptions {
		return &AlterSCIMIntegrationSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSCIMIntegrationSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SCIMIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SCIMIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SCIMIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSCIMIntegrationSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.Set = &SCIMIntegrationSet{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(networkPolicyID),
			SyncPassword:  Pointer(true),
			Comment:       Pointer("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true NETWORK_POLICY = %s SYNC_PASSWORD = true COMMENT = 'test'",
			id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &SCIMIntegrationUnset{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(true),
			SyncPassword:  Pointer(true),
			Comment:       Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, SYNC_PASSWORD, COMMENT", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DropSecurityIntegrationOptions
	defaultOpts := func() *DropSecurityIntegrationOptions {
		return &DropSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SECURITY INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid DescribeSecurityIntegrationOptions
	defaultOpts := func() *DescribeSecurityIntegrationOptions {
		return &DescribeSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SECURITY INTEGRATION %s", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_Show(t *testing.T) {
	// Minimal valid ShowSecurityIntegrationOptions
	defaultOpts := func() *ShowSecurityIntegrationOptions {
		return &ShowSecurityIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("some pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW SECURITY INTEGRATIONS LIKE 'some pattern'")
	})
}
