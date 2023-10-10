package sdk_integration_tests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SessionPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	database, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := sdk.createSchema(t, client, database)
	t.Cleanup(schemaCleanup)

	assertSessionPolicy := func(t *testing.T, sessionPolicy *sdk.SessionPolicy, id sdk.SchemaObjectIdentifier, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, sessionPolicy.CreatedOn)
		assert.Equal(t, id.Name(), sessionPolicy.Name)
		assert.Equal(t, id.SchemaName(), sessionPolicy.SchemaName)
		assert.Equal(t, id.DatabaseName(), sessionPolicy.DatabaseName)
		assert.Equal(t, "ACCOUNTADMIN", sessionPolicy.Owner)
		assert.Equal(t, expectedComment, sessionPolicy.Comment)
		assert.Equal(t, "SESSION_POLICY", sessionPolicy.Kind)
		assert.Equal(t, "", sessionPolicy.Options)
	}

	assertSessionPolicyDescription := func(
		t *testing.T,
		sessionPolicyDescription *sdk.SessionPolicyDescription,
		id sdk.SchemaObjectIdentifier,
	) {
		t.Helper()
		assert.NotEmpty(t, sessionPolicyDescription.CreatedOn)
		assert.Equal(t, id.Name(), sessionPolicyDescription.Name)
		assert.Equal(t, 240, sessionPolicyDescription.SessionIdleTimeoutMins)
		assert.Equal(t, 240, sessionPolicyDescription.SessionUIIdleTimeoutMins)
		assert.Equal(t, "", sessionPolicyDescription.Comment)
	}

	cleanupSessionPolicyProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id))
			require.NoError(t, err)
		}
	}

	createSessionPolicy := func(t *testing.T) *sdk.SessionPolicy {
		t.Helper()
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return sessionPolicy
	}

	t.Run("create session_policy: complete case", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := randomComment(t)

		request := sdk.NewCreateSessionPolicyRequest(id).
			WithSessionIdleTimeoutMins(sdk.Int(5)).
			WithSessionUiIdleTimeoutMins(sdk.Int(34)).
			WithComment(&comment).
			WithIfNotExists(sdk.Bool(true))

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, sessionPolicy, id, comment)
	})

	t.Run("create session_policy: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		request := sdk.NewCreateSessionPolicyRequest(id)

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, sessionPolicy, id, "")
	})

	t.Run("drop session_policy: existing", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)

		err = client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("drop session_policy: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, "does_not_exist")

		err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id))
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("alter session_policy: set value and unset value", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		alterRequest := sdk.NewAlterSessionPolicyRequest(id).WithSet(sdk.NewSessionPolicySetRequest().WithComment(sdk.String("new comment")))
		err = client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredSessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredSessionPolicy.Comment)

		alterRequest = sdk.NewAlterSessionPolicyRequest(id).WithUnset(sdk.NewSessionPolicyUnsetRequest().WithComment(sdk.Bool(true)))
		err = client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredSessionPolicy, err = client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredSessionPolicy.Comment)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		tag, tagCleanup := sdk.createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterSessionPolicyRequest(id).WithSetTags(tags)

		err = client.SessionPolicies.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeSessionPolicy)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterSessionPolicyRequest(id).WithUnsetTags(unsetTags)

		err = client.SessionPolicies.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeSessionPolicy)
		require.Error(t, err)
	})

	t.Run("alter session_policy: rename", func(t *testing.T) {
		name := randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, sdk.NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)

		newName := randomString(t)
		newId := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, newName)
		alterRequest := sdk.NewAlterSessionPolicyRequest(id).WithRenameTo(&newId)

		err = client.SessionPolicies.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupSessionPolicyProvider(id))
		} else {
			t.Cleanup(cleanupSessionPolicyProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertSessionPolicy(t, sessionPolicy, newId, "")
	})

	t.Run("show session_policy: default", func(t *testing.T) {
		sessionPolicy1 := createSessionPolicy(t)
		sessionPolicy2 := createSessionPolicy(t)

		showRequest := sdk.NewShowSessionPolicyRequest()
		returnedSessionPolicies, err := client.SessionPolicies.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedSessionPolicies))
		assert.Contains(t, returnedSessionPolicies, *sessionPolicy1)
		assert.Contains(t, returnedSessionPolicies, *sessionPolicy2)
	})

	t.Run("describe session_policy", func(t *testing.T) {
		sessionPolicy := createSessionPolicy(t)

		returnedSessionPolicy, err := client.SessionPolicies.Describe(ctx, sessionPolicy.ID())
		require.NoError(t, err)

		assertSessionPolicyDescription(t, returnedSessionPolicy, sessionPolicy.ID())
	})
}
