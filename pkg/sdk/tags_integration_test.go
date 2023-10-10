package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tags(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	assertTagHandle := func(t *testing.T, tag *Tag, expectedName string, expectedComment string, expectedAllowedValues []string) {
		t.Helper()
		assert.NotEmpty(t, tag.CreatedOn)
		assert.Equal(t, expectedName, tag.Name)
		assert.Equal(t, "ACCOUNTADMIN", tag.Owner)
		assert.Equal(t, expectedComment, tag.Comment)
		assert.Equal(t, expectedAllowedValues, tag.AllowedValues)
	}
	cleanupTagHandle := func(id SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tags.Drop(ctx, NewDropTagRequest(id))
			require.NoError(t, err)
		}
	}
	createTagHandle := func(t *testing.T) *Tag {
		t.Helper()

		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, randomString(t))
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		return tag
	}

	t.Run("create tag: comment", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		comment := randomComment(t)

		request := NewCreateTagRequest(id).WithComment(&comment)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, comment, nil)
	})

	t.Run("create tag: allowed values", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		values := []string{"value1", "value2"}
		request := NewCreateTagRequest(id).WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, "", values)
	})

	t.Run("create tag: comment and allowed values", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		comment := randomComment(t)
		values := []string{"value1", "value2"}
		err := client.Tags.Create(ctx, NewCreateTagRequest(id).WithOrReplace(true).WithComment(&comment).WithAllowedValues(values))
		expected := "Comment fields: [AllowedValues] are incompatible and cannot be set at the same time"
		require.Equal(t, expected, err.Error())
	})

	t.Run("create tag: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, "", nil)
	})

	t.Run("drop tag: existing", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)

		err = client.Tags.Drop(ctx, NewDropTagRequest(id))
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("drop tag: non-existing", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.Tags.Drop(ctx, NewDropTagRequest(id))
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("undrop tag: existing", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		err = client.Tags.Drop(ctx, NewDropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)

		err = client.Tags.Undrop(ctx, NewUndropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("alter tag: set and unset comment", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		// alter tag with set comment
		comment := randomComment(t)
		set := NewTagSetRequest().WithComment(comment)
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithSet(set))
		require.NoError(t, err)

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, comment, tag.Comment)

		// alter tag with unset comment
		unset := NewTagUnsetRequest().WithComment(true)
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", tag.Comment)
	})

	t.Run("alter tag: set and unset masking policies", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		policyTest, policyCleanup := createMaskingPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(policyCleanup)

		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		policies := []string{policyTest.Name}
		set := NewTagSetRequest().WithMaskingPolicies(policies)
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithSet(set))
		require.NoError(t, err)

		unset := NewTagUnsetRequest().WithMaskingPolicies(policies)
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)
	})

	t.Run("alter tag: add and drop allowed values", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		values := []string{"value1", "value2"}
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithDrop(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 0, len(tag.AllowedValues))
	})

	t.Run("alter tag: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)

		nid := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, randomString(t))
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithRename(nid))
		if err != nil {
			t.Cleanup(cleanupTagHandle(id))
		} else {
			t.Cleanup(cleanupTagHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)

		tag, err := client.Tags.ShowByID(ctx, nid)
		require.NoError(t, err)
		assertTagHandle(t, tag, nid.Name(), "", nil)
	})

	t.Run("alter tag: unset allowed values", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		values := []string{"value1", "value2"}
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		unset := NewTagUnsetRequest().WithAllowedValues(true)
		err = client.Tags.Alter(ctx, NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 0, len(tag.AllowedValues))
	})

	t.Run("show tag: without like", func(t *testing.T) {
		t1 := createTagHandle(t)
		t2 := createTagHandle(t)

		tags, err := client.Tags.Show(ctx, NewShowTagRequest())
		require.NoError(t, err)

		assert.Equal(t, 2, len(tags))
		assert.Contains(t, tags, *t1)
		assert.Contains(t, tags, *t2)
	})

	t.Run("show tag: with like", func(t *testing.T) {
		t1 := createTagHandle(t)
		t2 := createTagHandle(t)

		tags, err := client.Tags.Show(ctx, NewShowTagRequest().WithLike(t1.Name))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tags))
		assert.Contains(t, tags, *t1)
		assert.NotContains(t, tags, *t2)
	})

	t.Run("show tag: no matches", func(t *testing.T) {
		tags, err := client.Tags.Show(ctx, NewShowTagRequest().WithLike("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tags))
	})
}
