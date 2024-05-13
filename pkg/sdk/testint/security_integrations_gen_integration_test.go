package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SecurityIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupSecurityIntegration := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()
		t.Cleanup(func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(id).WithIfExists(sdk.Pointer(true)))
			assert.NoError(t, err)
		})
	}

	createSCIMIntegration := func(t *testing.T, siID sdk.AccountObjectIdentifier, with func(*sdk.CreateSCIMSecurityIntegrationRequest)) {
		t.Helper()
		roleID := sdk.NewAccountObjectIdentifier("GENERIC_SCIM_PROVISIONER")
		err := client.Roles.Create(ctx, sdk.NewCreateRoleRequest(roleID).WithIfNotExists(true))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(roleID))
			assert.NoError(t, err)
		})
		currentRole := testClientHelper().Context.CurrentRole(t)
		err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(roleID, sdk.GrantRole{Role: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentRole))}))
		require.NoError(t, err)

		scimReq := sdk.NewCreateSCIMSecurityIntegrationRequest(siID, false, "GENERIC", roleID.Name())
		if with != nil {
			with(scimReq)
		}
		err = client.SecurityIntegrations.CreateSCIM(ctx, scimReq)
		require.NoError(t, err)
		cleanupSecurityIntegration(t, siID)
	}

	assertSecurityIntegration := func(t *testing.T, si *sdk.SecurityIntegration, id sdk.AccountObjectIdentifier, siType string, enabled bool, comment string) {
		t.Helper()
		assert.Equal(t, id.Name(), si.Name)
		assert.Equal(t, siType, si.IntegrationType)
		assert.Equal(t, enabled, si.Enabled)
		assert.Equal(t, comment, si.Comment)
	}

	assertSCIMDescribe := func(details []sdk.SecurityIntegrationProperty, enabled, networkPolicy, runAsRole, syncPassword, comment string) {
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: enabled, Default: "false"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "NETWORK_POLICY", Type: "String", Value: networkPolicy, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "RUN_AS_ROLE", Type: "String", Value: runAsRole, Default: ""})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "SYNC_PASSWORD", Type: "Boolean", Value: syncPassword, Default: "true"})
		assert.Contains(t, details, sdk.SecurityIntegrationProperty{Name: "COMMENT", Type: "String", Value: comment, Default: ""})
	}

	t.Run("CreateSCIM", func(t *testing.T) {
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
		t.Cleanup(networkPolicyCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, func(r *sdk.CreateSCIMSecurityIntegrationRequest) {
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

	t.Run("AlterSCIMIntegration", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		createSCIMIntegration(t, id, nil)

		setRequest := sdk.NewAlterSCIMIntegrationSecurityIntegrationRequest(id).
			WithSet(
				sdk.NewSCIMIntegrationSetRequest().
					WithEnabled(sdk.Bool(true)).
					WithSyncPassword(sdk.Bool(true)).
					WithComment(sdk.String("altered")),
			)
		err := client.SecurityIntegrations.AlterSCIMIntegration(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.SecurityIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSCIMDescribe(details, "true", "", "GENERIC_SCIM_PROVISIONER", "true", "altered")
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
