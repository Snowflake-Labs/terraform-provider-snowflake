package testint

import (
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
			err := client.AuthenticationPolicies.Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id))
			require.NoError(t, err)
		}
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
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "AUTHENTICATION_METHODS", Value: "[PASSWORD]"})
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
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "CLIENT_TYPES", Value: "[DRIVERS, SNOWSQL]"})
	})

	t.Run("Alter - set security integrations", func(t *testing.T) {
		req := defaultCreateRequest()
		err := client.AuthenticationPolicies.Create(ctx, req)
		t.Cleanup(cleanupAuthenticationPolicyProvider(req.GetName()))

		alterErr := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(req.GetName()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().WithSecurityIntegrations([]sdk.SecurityIntegrationsOption{{Name: "sec-integration"}})))
		require.NoError(t, alterErr)

		desc, err := client.AuthenticationPolicies.Describe(ctx, req.GetName())
		require.NoError(t, err)
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "SECURITY_INTEGRATIONS", Value: "[sec-integration]"})
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
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "MFA_AUTHENTICATION_METHODS", Value: "[PASSWORD]"})
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
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "MFA_ENROLLMENT", Value: "REQUIRED"})
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

		assert.Equal(t, 2, len(desc))
		assert.Contains(t, desc, sdk.AuthenticationPolicyDescription{Name: "COMMENT", Value: "some comment"})
	})
}
