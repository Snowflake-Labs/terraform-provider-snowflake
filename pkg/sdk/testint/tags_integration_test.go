package testint

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
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

	assertTagHandle := func(t *testing.T, tag *sdk.Tag, expectedName string, expectedComment string, expectedAllowedValues []string) {
		t.Helper()
		assert.NotEmpty(t, tag.CreatedOn)
		assert.Equal(t, expectedName, tag.Name)
		assert.Equal(t, "ACCOUNTADMIN", tag.Owner)
		assert.Equal(t, expectedComment, tag.Comment)
		assert.Equal(t, expectedAllowedValues, tag.AllowedValues)
	}
	cleanupTagHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}
	createTagHandle := func(t *testing.T) *sdk.Tag {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		return tag
	}

	t.Run("create tag: comment", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		comment := random.Comment()

		request := sdk.NewCreateTagRequest(id).WithComment(&comment)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, comment, nil)
	})

	t.Run("create tag: allowed values", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		values := []string{"value1", "value2"}
		request := sdk.NewCreateTagRequest(id).WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, "", values)
	})

	t.Run("create tag: comment and allowed values", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		comment := random.Comment()
		values := []string{"value1", "value2"}
		err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id).WithOrReplace(true).WithComment(&comment).WithAllowedValues(values))
		expected := "Comment fields: [AllowedValues] are incompatible and cannot be set at the same time"
		require.Equal(t, expected, err.Error())
	})

	t.Run("create tag: no optionals", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTagHandle(id))

		tag, err := client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assertTagHandle(t, tag, name, "", nil)
	})

	t.Run("drop tag: existing", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop tag: non-existing", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("undrop tag: existing", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		err = client.Tags.Undrop(ctx, sdk.NewUndropTagRequest(id))
		require.NoError(t, err)
		_, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("alter tag: set and unset comment", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()

		// alter tag with set comment
		comment := random.Comment()
		set := sdk.NewTagSetRequest().WithComment(comment)
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(set))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, comment, tag.Comment)

		// alter tag with unset comment
		unset := sdk.NewTagUnsetRequest().WithComment(true)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", tag.Comment)
	})

	t.Run("alter tag: set and unset masking policies", func(t *testing.T) {
		policyTest, policyCleanup := createMaskingPolicy(t, client, databaseTest, schemaTest)
		t.Cleanup(policyCleanup)

		tag := createTagHandle(t)
		id := tag.ID()

		policies := []sdk.SchemaObjectIdentifier{policyTest.ID()}
		set := sdk.NewTagSetRequest().WithMaskingPolicies(policies)
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(set))
		require.NoError(t, err)

		unset := sdk.NewTagUnsetRequest().WithMaskingPolicies(policies)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)
	})

	t.Run("alter tag: add and drop allowed values", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()

		values := []string{"value1", "value2"}
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 0, len(tag.AllowedValues))
	})

	t.Run("alter tag: rename", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()

		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithRename(nid))
		if err != nil {
			t.Cleanup(cleanupTagHandle(id))
		} else {
			t.Cleanup(cleanupTagHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		tag, err = client.Tags.ShowByID(ctx, nid)
		require.NoError(t, err)
		assertTagHandle(t, tag, nid.Name(), "", nil)
	})

	t.Run("alter tag: unset allowed values", func(t *testing.T) {
		tag := createTagHandle(t)
		id := tag.ID()
		t.Cleanup(cleanupTagHandle(id))

		values := []string{"value1", "value2"}
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(values))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, values, tag.AllowedValues)

		unset := sdk.NewTagUnsetRequest().WithAllowedValues(true)
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset))
		require.NoError(t, err)

		tag, err = client.Tags.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 0, len(tag.AllowedValues))
	})

	t.Run("show tag: without like", func(t *testing.T) {
		t1 := createTagHandle(t)
		t2 := createTagHandle(t)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest())
		require.NoError(t, err)

		assert.Equal(t, 2, len(tags))
		assert.Contains(t, tags, *t1)
		assert.Contains(t, tags, *t2)
	})

	t.Run("show tag: with like", func(t *testing.T) {
		t1 := createTagHandle(t)
		t2 := createTagHandle(t)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike(t1.Name))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tags))
		assert.Contains(t, tags, *t1)
		assert.NotContains(t, tags, *t2)
	})

	t.Run("show tag: no matches", func(t *testing.T) {
		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tags))
	})
}
