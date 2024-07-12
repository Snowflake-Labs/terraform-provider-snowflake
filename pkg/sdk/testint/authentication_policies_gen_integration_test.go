package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe", func(t *testing.T) {
		// TODO: fill me
	})
}
