package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AuthenticationPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	cert := random.GenerateX509(t)

	assertAuthenticationPolicy := func(t *testing.T, authenticationPolicy *sdk.AuthenticationPolicy, id sdk.SchemaObjectIdentifier, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, authenticationPolicy.CreatedOn)
		assert.Equal(t, id.Name(), authenticationPolicy.Name)
		assert.Equal(t, id.SchemaName(), authenticationPolicy.SchemaName)
		assert.Equal(t, id.DatabaseName(), authenticationPolicy.DatabaseName)
		assert.Equal(t, "", authenticationPolicy.Options)
		assert.Equal(t, "ACCOUNTADMIN", authenticationPolicy.Owner)
		assert.Equal(t, expectedComment, authenticationPolicy.Comment)
		assert.Equal(t, "ROLE", authenticationPolicy.OwnerRoleType)
	}

	cleanupAuthenticationPolicyProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	cleanupSecurityIntegration := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(true))
			assert.NoError(t, err)
		})
	}

	createAuthenticationPolicy := func(t *testing.T) *sdk.AuthenticationPolicy {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupAuthenticationPolicyProvider(id))

		authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return authenticationPolicy
	}

	defaultCreateRequest := func() *sdk.CreateAuthenticationPolicyRequest {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := "some_comment"
		return sdk.NewCreateAuthenticationPolicyRequest(id).
			WithOrReplace(true).
			WithComment(comment)
	}

	createSAML2Integration := func(t *testing.T, with func(*sdk.CreateSaml2SecurityIntegrationRequest)) sdk.AccountObjectIdentifier {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		issuer := testClientHelper().Ids.Alpha()
		saml2Req := sdk.NewCreateSaml2SecurityIntegrationRequest(id, issuer, "https://example.com", sdk.Saml2SecurityIntegrationSaml2ProviderCustom, cert)
		if with != nil {
			with(saml2Req)
		}
		err := client.SecurityIntegrations.CreateSaml2(ctx, saml2Req)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, id)
		_, showErr := client.SecurityIntegrations.ShowByID(ctx, id)
		require.NoError(t, showErr)

		return id
	}

	t.Run("Create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreateAuthenticationPolicyRequest(id).
			WithAuthenticationMethods([]sdk.AuthenticationMethods{{Method: "Password"}}).
			WithComment(comment)

		err := client.AuthenticationPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupAuthenticationPolicyProvider(id))

		authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertAuthenticationPolicy(t, authenticationPolicy, id, comment)
	})

	t.Run("Alter - set authentication methods", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithAuthenticationMethods([]sdk.AuthenticationMethods{{Method: "PASSWORD"}})))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "AUTHENTICATION_METHODS", Value: "[PASSWORD]"})
	})

	t.Run("Alter - set client types", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithClientTypes([]sdk.ClientTypes{{ClientType: "DRIVERS"}, {ClientType: "SNOWSQL"}})))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "CLIENT_TYPES", Value: "[DRIVERS, SNOWSQL]"})
	})

	t.Run("Alter - set security integrations", func(t *testing.T) {
		secId := createSAML2Integration(t, func(r *sdk.CreateSaml2SecurityIntegrationRequest) {
			r.WithEnabled(true)
		})
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithSecurityIntegrations([]sdk.SecurityIntegrationsOption{{Name: secId.Name()}})))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "SECURITY_INTEGRATIONS", Value: fmt.Sprintf("[%s]", secId.Name())})
	})

	t.Run("Alter - set mfa authentication methods", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithMfaAuthenticationMethods([]sdk.MfaAuthenticationMethods{{Method: "PASSWORD"}})))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "MFA_AUTHENTICATION_METHODS", Value: "[PASSWORD]"})
	})

	t.Run("Alter - set mfa enrollment", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithMfaEnrollment("REQUIRED")))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "MFA_ENROLLMENT", Value: "REQUIRED"})
	})

	t.Run("Alter - set comment", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithComment("new comment")))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "COMMENT", Value: "new comment"})
	})

	t.Run("Alter - unset authentication methods", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithAuthenticationMethods(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "AUTHENTICATION_METHODS", Value: "[ALL]"})
	})

	t.Run("Alter - unset client types", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithClientTypes(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "CLIENT_TYPES", Value: "[ALL]"})
	})

	t.Run("Alter - unset security integrations", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithSecurityIntegrations(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "SECURITY_INTEGRATIONS", Value: "[ALL]"})
	})

	t.Run("Alter - unset mfa authentication methods", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithMfaAuthenticationMethods(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "MFA_AUTHENTICATION_METHODS", Value: "[PASSWORD, SAML]"})
	})

	t.Run("Alter - unset mfa enrollment", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithMfaEnrollment(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "MFA_ENROLLMENT", Value: "OPTIONAL"})
	})

	t.Run("Alter - unset comment", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().WithComment(true)))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "COMMENT", Value: "null"})
	})

	t.Run("Alter - rename", func(t *testing.T) {
		req := defaultCreateRequest()
		client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		t.Cleanup(cleanupAuthenticationPolicyProvider(newId))
		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithRenameTo(newId))
		require.NoError(t, alterErr)

		_, descErr := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		assert.ErrorIs(t, descErr, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id))
		require.NoError(t, err)

		err = client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.AuthenticationPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("Drop: non-existing", func(t *testing.T) {
		err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		authenticationPolicy1 := createAuthenticationPolicy(t)
		authenticationPolicy2 := createAuthenticationPolicy(t)

		showRequest := sdk.NewShowAuthenticationPolicyRequest()
		returnedAuthenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.LessOrEqual(t, 2, len(returnedAuthenticationPolicies))
		assert.Contains(t, returnedAuthenticationPolicies, *authenticationPolicy1)
		assert.Contains(t, returnedAuthenticationPolicies, *authenticationPolicy2)
	})

	t.Run("Describe", func(t *testing.T) {
		request := defaultCreateRequest()
		client.AuthenticationPolicies.Create(ctx, request)

		desc, err := client.AuthenticationPolicies.Describe(ctx, request.GetName())
		require.NoError(t, err)

		assert.Equal(t, 8, len(desc))
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Property: "COMMENT", Value: "some_comment"})
	})
}
