package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SessionPolicies(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchema(t, client, database)
	t.Cleanup(schemaCleanup)

	assertSessionPolicy := func(t *testing.T, sessionPolicy *SessionPolicy, id SchemaObjectIdentifier, expectedComment string) {
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
		sessionPolicyDescription *SessionPolicyDescription,
		id SchemaObjectIdentifier,
	) {
		t.Helper()
		assert.NotEmpty(t, sessionPolicyDescription.CreatedOn)
		assert.Equal(t, id.Name(), sessionPolicyDescription.Name)
		assert.Equal(t, 240, sessionPolicyDescription.SessionIdleTimeoutMins)
		assert.Equal(t, 240, sessionPolicyDescription.SessionUIIdleTimeoutMins)
		assert.Equal(t, "", sessionPolicyDescription.Comment)
	}

	cleanupSessionPolicyProvider := func(id SchemaObjectIdentifier) func() {
		return func() {
			err := client.SessionPolicies.Drop(ctx, NewDropSessionPolicyRequest(id))
			require.NoError(t, err)
		}
	}

	createSessionPolicy := func(t *testing.T) *SessionPolicy {
		t.Helper()
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return sessionPolicy
	}

	t.Run("create session_policy: complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := randomComment(t)

		request := NewCreateSessionPolicyRequest(id).
			WithSessionIdleTimeoutMins(Int(5)).
			WithSessionUiIdleTimeoutMins(Int(34)).
			WithComment(&comment).
			WithIfNotExists(Bool(true))

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		databaseRole, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, databaseRole, id, comment)
	})

	t.Run("create session_policy: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		request := NewCreateSessionPolicyRequest(id)

		err := client.SessionPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		databaseRole, err := client.SessionPolicies.ShowByID(ctx, id)

		require.NoError(t, err)
		assertSessionPolicy(t, databaseRole, id, "")
	})

	t.Run("drop session_policy: existing", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)

		err = client.SessionPolicies.Drop(ctx, NewDropSessionPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("drop session_policy: non-existing", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, "does_not_exist")

		err := client.SessionPolicies.Drop(ctx, NewDropSessionPolicyRequest(id))
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("alter session_policy: set value and unset value", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		alterRequest := NewAlterSessionPolicyRequest(id).WithSet(NewSessionPolicySetRequest().WithComment(String("new comment")))
		err = client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err := client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredDatabaseRole.Comment)

		alterRequest = NewAlterSessionPolicyRequest(id).WithUnset(NewSessionPolicyUnsetRequest().WithComment(Bool(true)))
		err = client.SessionPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err = client.SessionPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredDatabaseRole.Comment)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupSessionPolicyProvider(id))

		tagValue := "abc"
		tags := []TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := NewAlterSessionPolicyRequest(id).WithSetTags(tags)

		err = client.SessionPolicies.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, ObjectTypeSessionPolicy)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := NewAlterSessionPolicyRequest(id).WithUnsetTags(unsetTags)

		err = client.SessionPolicies.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, ObjectTypeSessionPolicy)
		require.Error(t, err)
	})

	t.Run("alter database_role: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.SessionPolicies.Create(ctx, NewCreateSessionPolicyRequest(id))
		require.NoError(t, err)

		newName := randomString(t)
		newId := NewSchemaObjectIdentifier(database.Name, schema.Name, newName)
		alterRequest := NewAlterSessionPolicyRequest(id).WithRenameTo(&newId)

		err = client.SessionPolicies.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupSessionPolicyProvider(id))
		} else {
			t.Cleanup(cleanupSessionPolicyProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.SessionPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)

		databaseRole, err := client.SessionPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertSessionPolicy(t, databaseRole, newId, "")
	})

	t.Run("show session_policy: default", func(t *testing.T) {
		sessionPolicy1 := createSessionPolicy(t)
		sessionPolicy2 := createSessionPolicy(t)

		showRequest := NewShowSessionPolicyRequest()
		returnedDatabaseRoles, err := client.SessionPolicies.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedDatabaseRoles))
		assert.Contains(t, returnedDatabaseRoles, *sessionPolicy1)
		assert.Contains(t, returnedDatabaseRoles, *sessionPolicy2)
	})

	t.Run("describe session_policy", func(t *testing.T) {
		sessionPolicy := createSessionPolicy(t)

		returnedSessionPolicy, err := client.SessionPolicies.Describe(ctx, sessionPolicy.ID())
		require.NoError(t, err)

		assertSessionPolicyDescription(t, returnedSessionPolicy, sessionPolicy.ID())
	})
}
