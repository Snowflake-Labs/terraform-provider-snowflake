package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SecurityIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	acsURL := testClientHelper().Context.ACSURL(t)
	issuerURL := testClientHelper().Context.IssuerURL(t)
	cert := random.GenerateX509(t)
	rsaKey := random.GenerateRSAPublicKey(t)

	revertParameter := testClientHelper().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterEnableIdentifierFirstLogin, "true")
	t.Cleanup(revertParameter)

	cleanupSecurityIntegration := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(sdk.Pointer(true)))
			assert.NoError(t, err)
		})
	}
	createOauthCustom := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateOauthCustomSecurityIntegrationRequest)) *sdk.SecurityIntegration {
		t.Helper()

		req := sdk.NewCreateOauthCustomSecurityIntegrationRequest(siID, sdk.OauthSecurityIntegrationClientTypePublic, "https://example.com")
		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateOauthCustom(ctx, req)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
		integration, err := client.SecurityIntegrations.ShowByID(ctx, siID)
		require.NoError(t, err)

		return integration
	}
	createOauthPartner := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateOauthPartnerSecurityIntegrationRequest)) *sdk.SecurityIntegration {
		t.Helper()

		req := sdk.NewCreateOauthPartnerSecurityIntegrationRequest(siID, sdk.OauthSecurityIntegrationClientLooker).
			WithOauthRedirectUri(sdk.Pointer("http://example.com"))

		if with != nil {
			with(req)
		}
		err := client.SecurityIntegrations.CreateOauthPartner(ctx, req)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
		integration, err := client.SecurityIntegrations.ShowByID(ctx, siID)
		require.NoError(t, err)

		return integration
	}
	createSAML2Integration := func(t *testing.T, siID sdk.AccountObjectIdentifier, issuer string, with func(*sdk.CreateSaml2SecurityIntegrationRequest)) *sdk.SecurityIntegration {
		t.Helper()

		saml2Req := sdk.NewCreateSaml2SecurityIntegrationRequest(siID, false, issuer, "https://example.com", "Custom", cert)
		if with != nil {
			with(saml2Req)
		}
		err := client.SecurityIntegrations.CreateSaml2(ctx, saml2Req)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
		integration, err := client.SecurityIntegrations.ShowByID(ctx, siID)
		require.NoError(t, err)

		return integration
	}

	createSCIMIntegration := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateScimSecurityIntegrationRequest)) *sdk.SecurityIntegration {
		t.Helper()
		role, roleCleanup := testClientHelper().Role.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(snowflakeroles.GenericScimProvisioner).WithOrReplace(true))
		t.Cleanup(roleCleanup)
		testClientHelper().Role.GrantRoleToCurrentRole(t, role.ID())

		scimReq := sdk.NewCreateScimSecurityIntegrationRequest(siID, false, sdk.ScimSecurityIntegrationScimClientGeneric, sdk.ScimSecurityIntegrationRunAsRoleGenericScimProvisioner)
		if with != nil {
			with(scimReq)
		}
		err := client.SecurityIntegrations.CreateScim(ctx, scimReq)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
		integration, err := client.SecurityIntegrations.ShowByID(ctx, siID)
		require.NoError(t, err)

		return integration
	}

	assertSecurityIntegration := func(t *testing.T, si *sdk.SecurityIntegration, id sdk.AccountObjectIdentifier, siType string, enabled bool, comment string) {
		t.Helper()
		assert.Equal(t, id.Name(), si.Name)
		assert.Equal(t, siType, si.IntegrationType)
		assert.Equal(t, enabled, si.Enabled)
		assert.Equal(t, comment, si.Comment)
		assert.Equal(t, "SECURITY", si.Category)
	}

	type snowflakeOauthPartnerDetails struct {
		enabled                 string
		oauthIssueRefreshTokens string
		refreshTokenValidity    string
		useSecondaryRoles       string
		preAuthorizedRolesList  string
		blockedRolesList        string
		networkPolicy           string
		comment                 string
	}

	assertOauthPartner := func(details []sdk.SecurityIntegrationProperty, d snowflakeOauthPartnerDetails) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: d.enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ISSUE_REFRESH_TOKENS", Type: "Boolean", Value: d.oauthIssueRefreshTokens, Default: "true"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_REFRESH_TOKEN_VALIDITY", Type: "Integer", Value: d.refreshTokenValidity, Default: "7776000"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: d.useSecondaryRoles, Default: "NONE"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "PRE_AUTHORIZED_ROLES_LIST", Type: "List", Value: d.preAuthorizedRolesList, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: d.networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
		// Chech one-by-one because snowflake returns a few extra roles
		found, err := collections.FindOne(details, func(d sdk.SecurityIntegrationProperty) bool { return d.Name == "BLOCKED_ROLES_LIST" })
		assert.NoError(t, err)
		roles := strings.Split(found.Value, ",")
		for _, exp := range strings.Split(d.blockedRolesList, ",") {
			assert.Contains(t, roles, exp)
		}
	}

	assertOauthCustom := func(details []sdk.SecurityIntegrationProperty, d snowflakeOauthPartnerDetails, allowNonTlsRedirectUri, clientType, enforcePkce string) {
		assertOauthPartner(details, d)
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOW_NON_TLS_REDIRECT_URI", Type: "Boolean", Value: allowNonTlsRedirectUri, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_TYPE", Type: "String", Value: clientType, Default: "CONFIDENTIAL"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_ENFORCE_PKCE", Type: "Boolean", Value: enforcePkce, Default: "false"})
		// Keys are hashed in snowflake, so we check only if these fields are present
		keys := make(map[string]struct{})
		for _, detail := range details {
			keys[detail.Name] = struct{}{}
		}
		assert.Contains(t, keys, "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP")
		assert.Contains(t, keys, "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP")
	}

	assertSCIMDescribe := func(details []sdk.SecurityIntegrationProperty, enabled, networkPolicy, runAsRole, syncPassword, comment string) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "RUN_AS_ROLE", Type: "String", Value: runAsRole, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SYNC_PASSWORD", Type: "Boolean", Value: syncPassword, Default: "true"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: comment, Default: ""})
	}

	type saml2Details struct {
		provider                  string
		enableSPInitiated         string
		spInitiatedLoginPageLabel string
		ssoURL                    string
		issuer                    string
		requestedNameIDFormat     string
		forceAuthn                string
		postLogoutRedirectUrl     string
		signrequest               string
		comment                   string
		snowflakeIssuerURL        string
		snowflakeAcsURL           string
		allowedUserDomains        string
		allowedEmailPatterns      string
	}

	assertSAML2Describe := func(details []sdk.SecurityIntegrationProperty, d saml2Details) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_X509_CERT", Type: "String", Value: cert, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_PROVIDER", Type: "String", Value: d.provider, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_ENABLE_SP_INITIATED", Type: "Boolean", Value: d.enableSPInitiated, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL", Type: "String", Value: d.spInitiatedLoginPageLabel, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SSO_URL", Type: "String", Value: d.ssoURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_ISSUER", Type: "String", Value: d.issuer, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_REQUESTED_NAMEID_FORMAT", Type: "String", Value: d.requestedNameIDFormat, Default: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_FORCE_AUTHN", Type: "Boolean", Value: d.forceAuthn, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_POST_LOGOUT_REDIRECT_URL", Type: "String", Value: d.postLogoutRedirectUrl, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SIGN_REQUEST", Type: "Boolean", Value: d.signrequest, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_DIGEST_METHODS_USED", Type: "String", Value: "http://www.w3.org/2001/04/xmlenc#sha256", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SIGNATURE_METHODS_USED", Type: "String", Value: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: d.comment, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SNOWFLAKE_ISSUER_URL", Type: "String", Value: d.snowflakeIssuerURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_SNOWFLAKE_ACS_URL", Type: "String", Value: d.snowflakeAcsURL, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ALLOWED_USER_DOMAINS", Type: "List", Value: d.allowedUserDomains, Default: "[]"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ALLOWED_EMAIL_PATTERNS", Type: "List", Value: d.allowedEmailPatterns, Default: "[]"})
	}

	t.Run("CreateOauthPartner", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)

		integration := createOauthPartner(t, id, func(r *sdk.CreateOauthPartnerSecurityIntegrationRequest) {
			r.WithBlockedRolesList(&sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithComment(sdk.Pointer("a")).
				WithEnabled(sdk.Pointer(true)).
				WithOauthIssueRefreshTokens(sdk.Pointer(true)).
				WithOauthRefreshTokenValidity(sdk.Pointer(12345)).
				WithOauthUseSecondaryRoles(sdk.Pointer(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit))
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthPartner(details, snowflakeOauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "12345",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			blockedRolesList:        role1.Name,
			comment:                 "a",
		})

		assertSecurityIntegration(t, integration, id, "OAUTH - LOOKER", true, "a")
	})

	t.Run("CreateOauthCustom", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)
		role2, role2Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role2Cleanup)

		integration := createOauthCustom(t, id, func(r *sdk.CreateOauthCustomSecurityIntegrationRequest) {
			r.WithBlockedRolesList(&sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
				WithComment(sdk.Pointer("a")).
				WithEnabled(sdk.Pointer(true)).
				WithNetworkPolicy(sdk.Pointer(sdk.NewAccountObjectIdentifier(networkPolicy.Name))).
				WithOauthAllowNonTlsRedirectUri(sdk.Pointer(true)).
				WithOauthClientRsaPublicKey(sdk.Pointer(rsaKey)).
				WithOauthClientRsaPublicKey2(sdk.Pointer(rsaKey)).
				WithOauthEnforcePkce(sdk.Pointer(true)).
				WithOauthIssueRefreshTokens(sdk.Pointer(true)).
				WithOauthRefreshTokenValidity(sdk.Pointer(12345)).
				WithOauthUseSecondaryRoles(sdk.Pointer(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
				WithPreAuthorizedRolesList(&sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: []sdk.AccountObjectIdentifier{role2.ID()}})
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthCustom(details, snowflakeOauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "12345",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  role2.Name,
			blockedRolesList:        role1.Name,
			networkPolicy:           networkPolicy.Name,
			comment:                 "a",
		}, "true", string(sdk.OauthSecurityIntegrationClientTypePublic), "true")

		assertSecurityIntegration(t, integration, id, "OAUTH - CUSTOM", true, "a")
	})

	t.Run("CreateSaml2", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()

		createSAML2Integration(t, id, issuer, func(r *sdk.CreateSaml2SecurityIntegrationRequest) {
			r.WithAllowedEmailPatterns([]sdk.EmailPattern{{Pattern: "^(.+dev)@example.com$"}}).
				WithAllowedUserDomains([]sdk.UserDomain{{Domain: "example.com"}}).
				WithComment(sdk.Pointer("a")).
				WithSaml2EnableSpInitiated(sdk.Pointer(true)).
				WithSaml2ForceAuthn(sdk.Pointer(true)).
				WithSaml2PostLogoutRedirectUrl(sdk.Pointer("http://example.com/logout")).
				WithSaml2RequestedNameidFormat(sdk.Pointer("urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified")).
				WithSaml2SignRequest(sdk.Pointer(true)).
				WithSaml2SnowflakeAcsUrl(&acsURL).
				WithSaml2SnowflakeIssuerUrl(&issuerURL).
				WithSaml2SpInitiatedLoginPageLabel(sdk.Pointer("label"))
			// TODO: fix after format clarification
			// WithSaml2SnowflakeX509Cert(sdk.Pointer(x509))
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2Details{
			provider:                  "Custom",
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "https://example.com",
			issuer:                    issuer,
			requestedNameIDFormat:     "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			forceAuthn:                "true",
			postLogoutRedirectUrl:     "http://example.com/logout",
			signrequest:               "true",
			comment:                   "a",
			snowflakeIssuerURL:        issuerURL,
			snowflakeAcsURL:           acsURL,
			allowedUserDomains:        "[example.com]",
			allowedEmailPatterns:      "[^(.+dev)@example.com$]",
		})

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SAML2", false, "a")
	})

	t.Run("CreateScim", func(t *testing.T) {
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, func(r *sdk.CreateScimSecurityIntegrationRequest) {
			r.WithComment(sdk.Pointer("a")).
				WithNetworkPolicy(sdk.Pointer(sdk.NewAccountObjectIdentifier(networkPolicy.Name))).
				WithSyncPassword(sdk.Pointer(false))
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "false", networkPolicy.Name, "GENERIC_SCIM_PROVISIONER", "false", "a")

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SCIM - GENERIC", false, "a")
	})

	t.Run("AlterOauthPartner", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createOauthPartner(t, id, func(r *sdk.CreateOauthPartnerSecurityIntegrationRequest) {
			r.WithOauthRedirectUri(sdk.Pointer("http://example.com"))
		})

		setRequest := sdk.NewAlterOauthPartnerSecurityIntegrationRequest(id).
			WithSet(
				sdk.NewOauthPartnerIntegrationSetRequest().
					WithBlockedRolesList(sdk.NewBlockedRolesListRequest()).
					WithComment(sdk.Pointer("a")).
					WithEnabled(sdk.Pointer(true)).
					WithOauthIssueRefreshTokens(sdk.Pointer(true)).
					WithOauthRedirectUri(sdk.Pointer("http://example2.com")).
					WithOauthRefreshTokenValidity(sdk.Pointer(22222)).
					WithOauthUseSecondaryRoles(sdk.Pointer(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)),
			)
		err := client.SecurityIntegrations.AlterOauthPartner(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthPartner(details, snowflakeOauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "22222",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  "",
			blockedRolesList:        "ACCOUNTADMIN,SECURITYADMIN",
			networkPolicy:           "",
			comment:                 "a",
		})

		unsetRequest := sdk.NewAlterOauthPartnerSecurityIntegrationRequest(id).
			WithUnset(
				sdk.NewOauthPartnerIntegrationUnsetRequest().
					WithEnabled(sdk.Pointer(true)).
					WithOauthUseSecondaryRoles(sdk.Pointer(true)),
			)
		err = client.SecurityIntegrations.AlterOauthPartner(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: "NONE", Default: "NONE"})
	})

	t.Run("AlterOauthPartner - set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createOauthPartner(t, id, nil)

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterOauthPartnerSecurityIntegrationRequest(id).WithSetTags(tags)

		err := client.SecurityIntegrations.AlterOauthPartner(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterOauthPartnerSecurityIntegrationRequest(id).WithUnsetTags(unsetTags)

		err = client.SecurityIntegrations.AlterOauthPartner(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.Error(t, err)
	})

	t.Run("AlterOauthCustom", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createOauthCustom(t, id, nil)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)
		role1, role1Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role1Cleanup)
		role2, role2Cleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(role2Cleanup)

		setRequest := sdk.NewAlterOauthCustomSecurityIntegrationRequest(id).
			WithSet(
				sdk.NewOauthCustomIntegrationSetRequest().
					WithEnabled(sdk.Pointer(true)).
					WithBlockedRolesList(&sdk.BlockedRolesListRequest{BlockedRolesList: []sdk.AccountObjectIdentifier{role1.ID()}}).
					WithComment(sdk.Pointer("a")).
					WithNetworkPolicy(sdk.Pointer(sdk.NewAccountObjectIdentifier(networkPolicy.Name))).
					WithOauthAllowNonTlsRedirectUri(sdk.Pointer(true)).
					WithOauthClientRsaPublicKey(sdk.Pointer(rsaKey)).
					WithOauthClientRsaPublicKey2(sdk.Pointer(rsaKey)).
					WithOauthEnforcePkce(sdk.Pointer(true)).
					WithOauthIssueRefreshTokens(sdk.Pointer(true)).
					WithOauthRedirectUri(sdk.Pointer("http://example2.com")).
					WithOauthRefreshTokenValidity(sdk.Pointer(22222)).
					WithOauthUseSecondaryRoles(sdk.Pointer(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit)).
					WithPreAuthorizedRolesList(&sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: []sdk.AccountObjectIdentifier{role2.ID()}}),
			)
		err := client.SecurityIntegrations.AlterOauthCustom(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertOauthCustom(details, snowflakeOauthPartnerDetails{
			enabled:                 "true",
			oauthIssueRefreshTokens: "true",
			refreshTokenValidity:    "22222",
			useSecondaryRoles:       string(sdk.OauthSecurityIntegrationUseSecondaryRolesImplicit),
			preAuthorizedRolesList:  role2.Name,
			blockedRolesList:        role1.Name,
			networkPolicy:           networkPolicy.Name,
			comment:                 "a",
		}, "true", string(sdk.OauthSecurityIntegrationClientTypePublic), "true")

		unsetRequest := sdk.NewAlterOauthCustomSecurityIntegrationRequest(id).
			WithUnset(
				sdk.NewOauthCustomIntegrationUnsetRequest().
					WithEnabled(sdk.Bool(true)).
					WithOauthUseSecondaryRoles(sdk.Bool(true)).
					WithNetworkPolicy(sdk.Bool(true)).
					WithOauthClientRsaPublicKey(sdk.Bool(true)).
					WithOauthClientRsaPublicKey2(sdk.Bool(true)),
			)
		err = client.SecurityIntegrations.AlterOauthCustom(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_USE_SECONDARY_ROLES", Type: "String", Value: "NONE", Default: "NONE"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_FP", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterOauthCustom - set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createOauthCustom(t, id, nil)

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterOauthCustomSecurityIntegrationRequest(id).WithSetTags(tags)

		err := client.SecurityIntegrations.AlterOauthCustom(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterOauthCustomSecurityIntegrationRequest(id).WithUnsetTags(unsetTags)

		err = client.SecurityIntegrations.AlterOauthCustom(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.Error(t, err)
	})
	t.Run("AlterSAML2Integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()
		createSAML2Integration(t, id, issuer, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).
			WithSet(
				sdk.NewSaml2IntegrationSetRequest().
					WithEnabled(sdk.Pointer(true)).
					WithSaml2Issuer(sdk.Pointer(issuer)).
					WithSaml2SsoUrl(sdk.Pointer("http://example.com")).
					WithSaml2Provider(sdk.Pointer("OKTA")).
					WithSaml2X509Cert(sdk.Pointer(cert)).
					WithComment(sdk.Pointer("a")).
					WithSaml2EnableSpInitiated(sdk.Pointer(true)).
					WithSaml2ForceAuthn(sdk.Pointer(true)).
					WithSaml2PostLogoutRedirectUrl(sdk.Pointer("http://example.com/logout")).
					WithSaml2RequestedNameidFormat(sdk.Pointer("urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified")).
					WithSaml2SignRequest(sdk.Pointer(true)).
					WithSaml2SnowflakeAcsUrl(&acsURL).
					WithSaml2SnowflakeIssuerUrl(&issuerURL).
					WithSaml2SpInitiatedLoginPageLabel(sdk.Pointer("label")).
					WithAllowedEmailPatterns([]sdk.EmailPattern{{Pattern: "^(.+dev)@example.com$"}}).
					WithAllowedUserDomains([]sdk.UserDomain{{Domain: "example.com"}}),
				// TODO: fix after format clarification
				// WithSaml2SnowflakeX509Cert(sdk.Pointer(cert)).
			)
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2Details{
			provider:                  "OKTA",
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "http://example.com",
			issuer:                    issuer,
			requestedNameIDFormat:     "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			forceAuthn:                "true",
			postLogoutRedirectUrl:     "http://example.com/logout",
			signrequest:               "true",
			comment:                   "a",
			snowflakeIssuerURL:        issuerURL,
			snowflakeAcsURL:           acsURL,
			allowedUserDomains:        "[example.com]",
			allowedEmailPatterns:      "[^(.+dev)@example.com$]",
		})

		unsetRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).
			WithUnset(
				sdk.NewSaml2IntegrationUnsetRequest().
					WithSaml2ForceAuthn(sdk.Pointer(true)).
					WithSaml2RequestedNameidFormat(sdk.Pointer(true)).
					WithSaml2PostLogoutRedirectUrl(sdk.Pointer(true)).
					WithComment(sdk.Pointer(true)),
			)
		err = client.SecurityIntegrations.AlterSaml2(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_FORCE_AUTHN", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_REQUESTED_NAMEID_FORMAT", Type: "String", Value: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", Default: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_POST_LOGOUT_REDIRECT_URL", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterSAML2Integration - REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()
		createSAML2Integration(t, id, issuer, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithRefreshSaml2SnowflakePrivateKey(sdk.Pointer(true))
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)
	})

	t.Run("AlterSAML2Integration - set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()
		createSAML2Integration(t, id, issuer, nil)

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithSetTags(tags)

		err := client.SecurityIntegrations.AlterSaml2(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithUnsetTags(unsetTags)

		err = client.SecurityIntegrations.AlterSaml2(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.Error(t, err)
	})

	t.Run("AlterSCIMIntegration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)

		setRequest := sdk.NewAlterScimSecurityIntegrationRequest(id).
			WithSet(
				sdk.NewScimIntegrationSetRequest().
					WithNetworkPolicy(sdk.Pointer(sdk.NewAccountObjectIdentifier(networkPolicy.Name))).
					WithEnabled(sdk.Bool(true)).
					WithSyncPassword(sdk.Bool(false)).
					WithComment(sdk.String("altered")),
			)
		err := client.SecurityIntegrations.AlterScim(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "true", networkPolicy.Name, "GENERIC_SCIM_PROVISIONER", "false", "altered")

		unsetRequest := sdk.NewAlterScimSecurityIntegrationRequest(id).
			WithUnset(
				sdk.NewScimIntegrationUnsetRequest().
					WithNetworkPolicy(sdk.Bool(true)).
					WithSyncPassword(sdk.Bool(true)),
			)
		err = client.SecurityIntegrations.AlterScim(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "true", "", "GENERIC_SCIM_PROVISIONER", "true", "altered")
	})

	t.Run("AlterSCIMIntegration - set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterScimSecurityIntegrationRequest(id).WithSetTags(tags)

		err := client.SecurityIntegrations.AlterScim(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnsetTags(unsetTags)

		err = client.SecurityIntegrations.AlterScim(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeIntegration)
		require.Error(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NotNil(t, si)
		require.NoError(t, err)

		err = client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id))
		require.NoError(t, err)

		si, err = client.SecurityIntegrations.ShowByID(ctx, id)
		require.Nil(t, si)
		require.Error(t, err)
	})

	t.Run("Drop non-existing", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier("does_not_exist")

		err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "false", "", "GENERIC_SCIM_PROVISIONER", "true", "")
	})

	t.Run("ShowByID", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		si, err := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)
		assertSecurityIntegration(t, si, id, "SCIM - GENERIC", false, "")
	})

	t.Run("Show OauthPartner", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si1 := createOauthPartner(t, id, nil)
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// more than one oauth partner integration is not allowed, create a custom one
		si2 := createOauthCustom(t, id2, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(&sdk.Like{
			Pattern: sdk.Pointer(id.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show OauthCustom", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si1 := createOauthCustom(t, id, nil)
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si2 := createOauthCustom(t, id2, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(&sdk.Like{
			Pattern: sdk.Pointer(id.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show SAML2", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si1 := createSAML2Integration(t, id, testClientHelper().Ids.Alpha(), nil)
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si2 := createSAML2Integration(t, id2, testClientHelper().Ids.Alpha(), nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(&sdk.Like{
			Pattern: sdk.Pointer(id.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})

	t.Run("Show SCIM", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si1 := createSCIMIntegration(t, id, nil)
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifier()
		si2 := createSCIMIntegration(t, id2, nil)

		returnedIntegrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest().WithLike(&sdk.Like{
			Pattern: sdk.Pointer(id.Name()),
		}))
		require.NoError(t, err)
		assert.Contains(t, returnedIntegrations, *si1)
		assert.NotContains(t, returnedIntegrations, *si2)
	})
}
