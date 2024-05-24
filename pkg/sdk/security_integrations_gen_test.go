package sdk

import (
	"testing"
)

func TestSecurityIntegrations_CreateOauthCustom(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateOauthForCustomClientsSecurityIntegrationOptions
	defaultOpts := func() *CreateOauthForCustomClientsSecurityIntegrationOptions {
		return &CreateOauthForCustomClientsSecurityIntegrationOptions{
			name:             id,
			OauthClientType:  OauthSecurityIntegrationClientTypePublic,
			OauthRedirectUri: "uri",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOauthForCustomClientsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOauthForCustomClientsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		roleID, role2ID, npID := randomAccountObjectIdentifier(), randomAccountObjectIdentifier(), randomAccountObjectIdentifier()
		opts.IfNotExists = Bool(true)
		opts.OauthClientType = OauthSecurityIntegrationClientTypePublic
		opts.OauthRedirectUri = "uri"
		opts.Enabled = Pointer(true)
		opts.OauthAllowNonTlsRedirectUri = Pointer(true)
		opts.OauthEnforcePkce = Pointer(true)
		opts.OauthUseSecondaryRoles = Pointer(OauthSecurityIntegrationUseSecondaryRolesNone)
		opts.PreAuthorizedRolesList = &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{roleID}}
		opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{role2ID}}
		opts.OauthIssueRefreshTokens = Pointer(true)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.NetworkPolicy = Pointer(npID)
		opts.OauthClientRsaPublicKey = Pointer("key")
		opts.OauthClientRsaPublicKey2 = Pointer("key2")
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = CUSTOM OAUTH_CLIENT_TYPE = 'PUBLIC' OAUTH_REDIRECT_URI = 'uri' ENABLED = true"+
			" OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true OAUTH_ENFORCE_PKCE = true OAUTH_USE_SECONDARY_ROLES = NONE PRE_AUTHORIZED_ROLES_LIST = (%s) BLOCKED_ROLES_LIST = (%s)"+
			" OAUTH_ISSUE_REFRESH_TOKENS = true OAUTH_REFRESH_TOKEN_VALIDITY = 42 NETWORK_POLICY = %s OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key' OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2' COMMENT = 'a'",
			id.FullyQualifiedName(), roleID.FullyQualifiedName(), role2ID.FullyQualifiedName(), npID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateOauthPartner(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateOauthForPartnerApplicationsSecurityIntegrationOptions
	defaultOpts := func() *CreateOauthForPartnerApplicationsSecurityIntegrationOptions {
		return &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
			name:        id,
			OauthClient: OauthSecurityIntegrationClientTableauDesktop,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateOauthForPartnerApplicationsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOauthForPartnerApplicationsSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: OAUTH_REDIRECT_URI is required when OAUTH_CLIENT=LOOKER", func(t *testing.T) {
		opts := &CreateOauthForPartnerApplicationsSecurityIntegrationOptions{
			name:        id,
			OauthClient: OauthSecurityIntegrationClientLooker,
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("OauthRedirectUri is required when OauthClient is LOOKER"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = OAUTH OAUTH_CLIENT = TABLEAU_DESKTOP", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		blockedRoleID := randomAccountObjectIdentifier()
		opts.IfNotExists = Bool(true)
		opts.OauthClient = OauthSecurityIntegrationClientLooker
		opts.OauthRedirectUri = Pointer("uri")
		opts.Enabled = Pointer(true)
		opts.OauthIssueRefreshTokens = Pointer(true)
		opts.OauthRefreshTokenValidity = Pointer(42)
		opts.OauthUseSecondaryRoles = Pointer(OauthSecurityIntegrationUseSecondaryRolesNone)
		opts.BlockedRolesList = &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{blockedRoleID}}
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = OAUTH OAUTH_CLIENT = LOOKER OAUTH_REDIRECT_URI = 'uri' ENABLED = true OAUTH_ISSUE_REFRESH_TOKENS = true"+
			" OAUTH_REFRESH_TOKEN_VALIDITY = 42 OAUTH_USE_SECONDARY_ROLES = NONE BLOCKED_ROLES_LIST = (%s) COMMENT = 'a'", id.FullyQualifiedName(), blockedRoleID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateSaml2(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateSaml2SecurityIntegrationOptions
	defaultOpts := func() *CreateSaml2SecurityIntegrationOptions {
		return &CreateSaml2SecurityIntegrationOptions{
			name:          id,
			Enabled:       true,
			Saml2Issuer:   "issuer",
			Saml2SsoUrl:   "url",
			Saml2Provider: "provider",
			Saml2X509Cert: "cert",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSaml2SecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSaml2SecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
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
		opts.Saml2SnowflakeX509Cert = Pointer("cert")

		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SAML2 ENABLED = true SAML2_ISSUER = 'issuer' SAML2_SSO_URL = 'url' SAML2_PROVIDER = 'provider' SAML2_X509_CERT = 'cert'"+
			" ALLOWED_USER_DOMAINS = ('domain') ALLOWED_EMAIL_PATTERNS = ('pattern') SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label' SAML2_ENABLE_SP_INITIATED = true SAML2_SNOWFLAKE_X509_CERT = 'cert' SAML2_SIGN_REQUEST = true"+
			" SAML2_REQUESTED_NAMEID_FORMAT = 'format' SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect' SAML2_FORCE_AUTHN = true SAML2_SNOWFLAKE_ISSUER_URL = 'issuer' SAML2_SNOWFLAKE_ACS_URL = 'acs'"+
			" COMMENT = 'a'", id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_CreateScim(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid CreateScimSecurityIntegrationOptions
	defaultOpts := func() *CreateScimSecurityIntegrationOptions {
		return &CreateScimSecurityIntegrationOptions{
			name:       id,
			Enabled:    true,
			ScimClient: "GENERIC",
			RunAsRole:  "GENERIC_SCIM_PROVISIONER",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateScimSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateScimSecurityIntegrationOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SECURITY INTEGRATION %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.IfNotExists = Pointer(true)
		opts.NetworkPolicy = Pointer(networkPolicyID)
		opts.SyncPassword = Pointer(true)
		opts.Comment = Pointer("a")
		assertOptsValidAndSQLEquals(t, opts, "CREATE SECURITY INTEGRATION IF NOT EXISTS %s TYPE = SCIM ENABLED = true SCIM_CLIENT = 'GENERIC' RUN_AS_ROLE = 'GENERIC_SCIM_PROVISIONER'"+
			" NETWORK_POLICY = %s SYNC_PASSWORD = true COMMENT = 'a'", id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterOauthPartner(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterOauthForPartnerApplicationsSecurityIntegrationOptions
	defaultOpts := func() *AlterOauthForPartnerApplicationsSecurityIntegrationOptions {
		return &AlterOauthForPartnerApplicationsSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterOauthForPartnerApplicationsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Set", "Enabled", "OauthIssueRefreshTokens",
			"OauthRedirectUri", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles", "BlockedRolesList", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions.Unset",
			"Enabled", "OauthUseSecondaryRoles"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{}
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForPartnerApplicationsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("empty roles lists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			BlockedRolesList: &BlockedRolesList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET BLOCKED_ROLES_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		roleID := randomAccountObjectIdentifier()
		opts.Set = &OauthForPartnerApplicationsIntegrationSet{
			Enabled:                   Pointer(true),
			OauthRedirectUri:          Pointer("uri"),
			OauthIssueRefreshTokens:   Pointer(true),
			OauthRefreshTokenValidity: Pointer(42),
			OauthUseSecondaryRoles:    Pointer(OauthSecurityIntegrationUseSecondaryRolesNone),
			BlockedRolesList:          &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{roleID}},
			Comment:                   Pointer("a"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_REFRESH_TOKEN_VALIDITY = 42,"+
			" OAUTH_USE_SECONDARY_ROLES = NONE, BLOCKED_ROLES_LIST = (%s), COMMENT = 'a'", id.FullyQualifiedName(), roleID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForPartnerApplicationsIntegrationUnset{
			Enabled:                Pointer(true),
			OauthUseSecondaryRoles: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, OAUTH_USE_SECONDARY_ROLES", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterOauthCustom(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterOauthForCustomClientsSecurityIntegrationOptions
	defaultOpts := func() *AlterOauthForCustomClientsSecurityIntegrationOptions {
		return &AlterOauthForCustomClientsSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterOauthForCustomClientsSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Set", "Enabled", "OauthRedirectUri", "OauthAllowNonTlsRedirectUri",
			"OauthEnforcePkce", "PreAuthorizedRolesList", "BlockedRolesList", "OauthIssueRefreshTokens", "OauthRefreshTokenValidity", "OauthUseSecondaryRoles",
			"NetworkPolicy", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForCustomClientsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions.Unset",
			"Enabled", "NetworkPolicy", "OauthUseSecondaryRoles", "OauthClientRsaPublicKey", "OauthClientRsaPublicKey2"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{}
		opts.Unset = &OauthForCustomClientsIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOauthForCustomClientsSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("empty roles lists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			PreAuthorizedRolesList: &PreAuthorizedRolesList{},
			BlockedRolesList:       &BlockedRolesList{},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET PRE_AUTHORIZED_ROLES_LIST = (), BLOCKED_ROLES_LIST = ()", id.FullyQualifiedName())
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		roleID, role2ID, npID := randomAccountObjectIdentifier(), randomAccountObjectIdentifier(), randomAccountObjectIdentifier()
		opts.Set = &OauthForCustomClientsIntegrationSet{
			Enabled:                     Pointer(true),
			OauthRedirectUri:            Pointer("uri"),
			OauthAllowNonTlsRedirectUri: Pointer(true),
			OauthEnforcePkce:            Pointer(true),
			OauthUseSecondaryRoles:      Pointer(OauthSecurityIntegrationUseSecondaryRolesNone),
			PreAuthorizedRolesList:      &PreAuthorizedRolesList{PreAuthorizedRolesList: []AccountObjectIdentifier{roleID}},
			BlockedRolesList:            &BlockedRolesList{BlockedRolesList: []AccountObjectIdentifier{role2ID}},
			OauthIssueRefreshTokens:     Pointer(true),
			OauthRefreshTokenValidity:   Pointer(42),
			NetworkPolicy:               Pointer(npID),
			OauthClientRsaPublicKey:     Pointer("key"),
			OauthClientRsaPublicKey2:    Pointer("key2"),
			Comment:                     Pointer("a"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, OAUTH_REDIRECT_URI = 'uri', OAUTH_ALLOW_NON_TLS_REDIRECT_URI = true, OAUTH_ENFORCE_PKCE = true,"+
			" PRE_AUTHORIZED_ROLES_LIST = (%s), BLOCKED_ROLES_LIST = (%s), OAUTH_ISSUE_REFRESH_TOKENS = true, OAUTH_REFRESH_TOKEN_VALIDITY = 42, OAUTH_USE_SECONDARY_ROLES = NONE,"+
			" NETWORK_POLICY = %s, OAUTH_CLIENT_RSA_PUBLIC_KEY = 'key', OAUTH_CLIENT_RSA_PUBLIC_KEY_2 = 'key2', COMMENT = 'a'", id.FullyQualifiedName(), roleID.FullyQualifiedName(), role2ID.FullyQualifiedName(), npID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OauthForCustomClientsIntegrationUnset{
			Enabled:                  Pointer(true),
			OauthUseSecondaryRoles:   Pointer(true),
			NetworkPolicy:            Pointer(true),
			OauthClientRsaPublicKey:  Pointer(true),
			OauthClientRsaPublicKey2: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, OAUTH_CLIENT_RSA_PUBLIC_KEY, OAUTH_CLIENT_RSA_PUBLIC_KEY_2, OAUTH_USE_SECONDARY_ROLES", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterSaml2(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterSaml2IntegrationSecurityIntegrationOptions
	defaultOpts := func() *AlterSaml2SecurityIntegrationOptions {
		return &AlterSaml2SecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSaml2SecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSaml2SecurityIntegrationOptions", "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Set", "Enabled", "Saml2Issuer", "Saml2SsoUrl", "Saml2Provider",
			"Saml2X509Cert", "AllowedUserDomains", "AllowedEmailPatterns", "Saml2SpInitiatedLoginPageLabel", "Saml2EnableSpInitiated", "Saml2SnowflakeX509Cert", "Saml2SignRequest",
			"Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Saml2ForceAuthn", "Saml2SnowflakeIssuerUrl", "Saml2SnowflakeAcsUrl", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Saml2IntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterSaml2SecurityIntegrationOptions.Unset",
			"Saml2ForceAuthn", "Saml2RequestedNameidFormat", "Saml2PostLogoutRedirectUrl", "Comment"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{}
		opts.Unset = &Saml2IntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSaml2SecurityIntegrationOptions", "Set", "Unset", "RefreshSaml2SnowflakePrivateKey", "SetTags", "UnsetTags"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &Saml2IntegrationSet{
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
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, SAML2_ISSUER = 'issuer', SAML2_SSO_URL = 'url', SAML2_PROVIDER = 'provider', SAML2_X509_CERT = 'cert',"+
			" ALLOWED_USER_DOMAINS = ('domain'), ALLOWED_EMAIL_PATTERNS = ('pattern'), SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'label', SAML2_ENABLE_SP_INITIATED = true, SAML2_SNOWFLAKE_X509_CERT = 'cert', SAML2_SIGN_REQUEST = true,"+
			" SAML2_REQUESTED_NAMEID_FORMAT = 'format', SAML2_POST_LOGOUT_REDIRECT_URL = 'redirect', SAML2_FORCE_AUTHN = true, SAML2_SNOWFLAKE_ISSUER_URL = 'issuer', SAML2_SNOWFLAKE_ACS_URL = 'acs',"+
			" COMMENT = 'a'", id.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &Saml2IntegrationUnset{
			Saml2ForceAuthn:            Pointer(true),
			Saml2RequestedNameidFormat: Pointer(true),
			Saml2PostLogoutRedirectUrl: Pointer(true),
			Comment:                    Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET SAML2_FORCE_AUTHN, SAML2_REQUESTED_NAMEID_FORMAT, SAML2_POST_LOGOUT_REDIRECT_URL, COMMENT", id.FullyQualifiedName())
	})

	t.Run("refresh SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		opts := defaultOpts()
		opts.RefreshSaml2SnowflakePrivateKey = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
	})
}

func TestSecurityIntegrations_AlterScim(t *testing.T) {
	id := randomAccountObjectIdentifier()

	// Minimal valid AlterScimSecurityIntegrationOptions
	defaultOpts := func() *AlterScimSecurityIntegrationOptions {
		return &AlterScimSecurityIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterScimSecurityIntegrationOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{
			Enabled: Pointer(true),
		}
		opts.name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterScimSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one of the fields [opts.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{}
		opts.Unset = &ScimIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterScimSecurityIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ScimIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Set", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.*] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ScimIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterScimSecurityIntegrationOptions.Unset", "Enabled", "NetworkPolicy", "SyncPassword", "Comment"))
	})

	t.Run("all options - set", func(t *testing.T) {
		opts := defaultOpts()
		networkPolicyID := randomAccountObjectIdentifier()
		opts.Set = &ScimIntegrationSet{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(networkPolicyID),
			SyncPassword:  Pointer(true),
			Comment:       Pointer("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s SET ENABLED = true, NETWORK_POLICY = %s, SYNC_PASSWORD = true, COMMENT = 'test'",
			id.FullyQualifiedName(), networkPolicyID.FullyQualifiedName())
	})

	t.Run("all options - unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ScimIntegrationUnset{
			Enabled:       Pointer(true),
			NetworkPolicy: Pointer(true),
			SyncPassword:  Pointer(true),
			Comment:       Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER SECURITY INTEGRATION %s UNSET ENABLED, NETWORK_POLICY, SYNC_PASSWORD, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
			{
				Name:  NewAccountObjectIdentifier("second-name"),
				Value: "second-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s SET TAG "name" = 'value', "second-name" = 'second-value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
			NewAccountObjectIdentifier("second-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SECURITY INTEGRATION %s UNSET TAG "name", "second-name"`, id.FullyQualifiedName())
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
