package testint

import (
	"context"
	"testing"

	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tags(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	assertTagHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier, expectedComment string, expectedAllowedValues []string) {
		t.Helper()
		assertions.AssertThatObject(t, objectassert.Tag(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(expectedComment).
			HasAllowedValues(expectedAllowedValues...).
			HasOwnerRoleType("ROLE"),
		)
	}

	t.Run("create tag: comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreateTagRequest(id).WithComment(&comment)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, comment, nil)
	})

	t.Run("create tag: allowed values", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		values := []string{"value1", "value2"}
		request := sdk.NewCreateTagRequest(id).WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, "", values)
	})

	t.Run("create tag: comment and allowed values", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		comment := random.Comment()
		values := []string{"value1", "value2"}
		request := sdk.NewCreateTagRequest(id).
			WithOrReplace(true).
			WithComment(&comment).
			WithAllowedValues(values)
		err := client.Tags.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, comment, values)
	})

	t.Run("create tag: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

		assertTagHandle(t, id, "", nil)
	})

	t.Run("drop tag: existing", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop tag: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("undrop tag: existing", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
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
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		// alter tag with set comment
		comment := random.Comment()
		set := sdk.NewTagSetRequest().WithComment(comment)
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithIfExists(true).WithSet(set))
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
		policyTest, policyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(policyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
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
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
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
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()

		nid := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithRename(nid))
		if err != nil {
			t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))
		} else {
			t.Cleanup(testClientHelper().Tag.DropTagFunc(t, nid))
		}
		require.NoError(t, err)

		_, err = client.Tags.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertTagHandle(t, nid, "", nil)
	})

	t.Run("alter tag: unset allowed values", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)
		id := tag.ID()
		t.Cleanup(testClientHelper().Tag.DropTagFunc(t, id))

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
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest())
		require.NoError(t, err)

		assert.Equal(t, 2, len(tags))
		assert.Contains(t, tags, *tag1)
		assert.Contains(t, tags, *tag2)
	})

	t.Run("show tag: with like", func(t *testing.T) {
		tag1, tag1Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag1Cleanup)
		tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tag2Cleanup)

		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike(tag1.Name))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tags))
		assert.Contains(t, tags, *tag1)
		assert.NotContains(t, tags, *tag2)
	})

	t.Run("show tag: no matches", func(t *testing.T) {
		tags, err := client.Tags.Show(ctx, sdk.NewShowTagRequest().WithLike("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tags))
	})
}

func TestInt_TagsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, tag1Cleanup := testClientHelper().Tag.CreateTagWithIdentifier(t, id1)
		t.Cleanup(tag1Cleanup)

		_, tag2Cleanup := testClientHelper().Tag.CreateTagWithIdentifier(t, id2)
		t.Cleanup(tag2Cleanup)

		e1, err := client.Tags.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Tags.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
