package testint

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SecurityIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	acsURL := fmt.Sprintf("https://%s.snowflakecomputing.com/fed/login", testClientHelper().Context.CurrentAccount(t))
	issuerURL := fmt.Sprintf("https://%s.snowflakecomputing.com", testClientHelper().Context.CurrentAccount(t))

	cleanupSecurityIntegration := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(sdk.Pointer(true)))
			assert.NoError(t, err)
		})
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
		},
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	require.NoError(t, err)

	cert := strings.TrimPrefix(certPEM.String(), "-----BEGIN CERTIFICATE-----\n")
	cert = strings.TrimSuffix(cert, "-----END CERTIFICATE-----\n")

	createSAML2Integration := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateSaml2SecurityIntegrationRequest)) {
		t.Helper()
		_, err := client.ExecForTests(ctx, "ALTER ACCOUNT SET ENABLE_IDENTIFIER_FIRST_LOGIN = true")
		require.NoError(t, err)

		saml2Req := sdk.NewCreateSaml2SecurityIntegrationRequest(siID, false, "test", "https://example.com", "Custom", cert)
		if with != nil {
			with(saml2Req)
		}
		err = client.SecurityIntegrations.CreateSaml2(ctx, saml2Req)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
	}

	createSCIMIntegration := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateScimSecurityIntegrationRequest)) {
		t.Helper()
		role, roleCleanup := testClientHelper().Role.CreateRoleWithName(t, "GENERIC_SCIM_PROVISIONER")
		t.Cleanup(roleCleanup)
		testClientHelper().Role.GrantRoleToCurrentRole(t, role.ID())

		scimReq := sdk.NewCreateScimSecurityIntegrationRequest(siID, false, sdk.ScimSecurityIntegrationScimClientGeneric, sdk.ScimSecurityIntegrationRunAsRoleGenericScimProvisioner)
		if with != nil {
			with(scimReq)
		}
		err = client.SecurityIntegrations.CreateScim(ctx, scimReq)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
	}

	assertSecurityIntegration := func(t *testing.T, si *sdk.SecurityIntegration, id sdk.AccountObjectIdentifier, siType string, enabled bool, comment string) {
		t.Helper()
		assert.Equal(t, id.Name(), si.Name)
		assert.Equal(t, siType, si.IntegrationType)
		assert.Equal(t, enabled, si.Enabled)
		assert.Equal(t, comment, si.Comment)
		assert.Equal(t, "SECURITY", si.Category)
	}

	assertSCIMDescribe := func(details []sdk.SecurityIntegrationProperty, enabled, networkPolicy, runAsRole, syncPassword, comment string) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "RUN_AS_ROLE", Type: "String", Value: runAsRole, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SYNC_PASSWORD", Type: "Boolean", Value: syncPassword, Default: "true"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: comment, Default: ""})
	}

	type saml2details struct {
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

	assertSAML2Describe := func(details []sdk.SecurityIntegrationProperty, d saml2details) {
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

	t.Run("CreateSaml2", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSAML2Integration(t, id, func(r *sdk.CreateSaml2SecurityIntegrationRequest) {
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
			// WithSaml2SnowflakeX509Cert(sdk.Pointer(x509))
		})
		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2details{
			provider:                  "Custom",
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "https://example.com",
			issuer:                    "test",
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

	t.Run("AlterSAML2Integration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSAML2Integration(t, id, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).
			WithSet(
				sdk.NewSaml2IntegrationSetRequest().
					WithEnabled(sdk.Pointer(true)).
					WithSaml2Issuer(sdk.Pointer("issuer")).
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
					// WithSaml2SnowflakeX509Cert(sdk.Pointer(cert)).
					WithAllowedEmailPatterns([]sdk.EmailPattern{{Pattern: "^(.+dev)@example.com$"}}).
					WithAllowedUserDomains([]sdk.UserDomain{{Domain: "example.com"}}),
			)
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSAML2Describe(details, saml2details{
			provider:                  "OKTA",
			enableSPInitiated:         "true",
			spInitiatedLoginPageLabel: "label",
			ssoURL:                    "http://example.com",
			issuer:                    "issuer",
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
					WithSaml2PostLogoutRedirectUrl(sdk.Pointer(true)),
			)
		err = client.SecurityIntegrations.AlterSaml2(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_FORCE_AUTHN", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_REQUESTED_NAMEID_FORMAT", Type: "String", Value: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", Default: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SAML2_POST_LOGOUT_REDIRECT_URL", Type: "String", Value: "", Default: ""})
	})

	t.Run("AlterSAML2Integration - REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSAML2Integration(t, id, nil)

		setRequest := sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithRefreshSaml2SnowflakePrivateKey(sdk.Pointer(true))
		err := client.SecurityIntegrations.AlterSaml2(ctx, setRequest)
		require.NoError(t, err)
	})

	t.Run("AlterSAML2Integration - set and unset tags", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSAML2Integration(t, id, nil)

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
}
