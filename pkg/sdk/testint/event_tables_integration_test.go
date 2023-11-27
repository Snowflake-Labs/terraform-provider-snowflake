package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_EventTables(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	tagTest, tagCleaup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleaup)

	assertEventTableHandle := func(t *testing.T, et *sdk.EventTable, expectedName string, expectedComment string, expectedAllowedValues []string) {
		t.Helper()
		assert.NotEmpty(t, et.CreatedOn)
		assert.Equal(t, expectedName, et.Name)
		assert.Equal(t, "ACCOUNTADMIN", et.Owner)
		assert.Equal(t, expectedComment, et.Comment)
	}

	cleanupTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", id.DatabaseName(), id.SchemaName(), id.Name()))
			require.NoError(t, err)
		}
	}

	createEventTableHandle := func(t *testing.T) *sdk.EventTable {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		et, err := client.EventTables.ShowByID(ctx, id)
		require.NoError(t, err)
		return et
	}

	t.Run("create event tables: all options", func(t *testing.T) {
		name := random.StringN(4)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		request := sdk.NewCreateEventTableRequest(id).
			WithChangeTracking(sdk.Bool(true)).
			WithDefaultDdlCollation(sdk.String("en_US")).
			WithDataRetentionTimeInDays(sdk.Int(1)).
			WithMaxDataExtensionTimeInDays(sdk.Int(2)).
			WithComment(sdk.String("test")).
			WithIfNotExists(sdk.Bool(true)).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
			})
		err := client.EventTables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))
	})

	t.Run("show event table: without like", func(t *testing.T) {
		et1 := createEventTableHandle(t)
		et2 := createEventTableHandle(t)

		tables, err := client.EventTables.Show(ctx, sdk.NewShowEventTableRequest())
		require.NoError(t, err)

		assert.Equal(t, 2, len(tables))
		assert.Contains(t, tables, *et1)
		assert.Contains(t, tables, *et2)
	})

	t.Run("show event table: with like", func(t *testing.T) {
		et1 := createEventTableHandle(t)
		et2 := createEventTableHandle(t)

		tables, err := client.EventTables.Show(ctx, sdk.NewShowEventTableRequest().WithLike(et1.Name))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tables))
		assert.Contains(t, tables, *et1)
		assert.NotContains(t, tables, *et2)
	})

	t.Run("show event table: no matches", func(t *testing.T) {
		tables, err := client.EventTables.Show(ctx, sdk.NewShowEventTableRequest().WithLike("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tables))
	})

	t.Run("describe event table", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		details, err := client.EventTables.Describe(ctx, sdk.NewDescribeEventTableRequest(id))
		require.NoError(t, err)
		assert.Equal(t, "TIMESTAMP", details.Name)
	})

	t.Run("alter event table: set and unset comment", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		comment := random.Comment()
		set := sdk.NewEventTableSetRequest().WithComment(&comment)
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithSet(set))
		require.NoError(t, err)

		et, err := client.EventTables.ShowByID(ctx, id)
		require.NoError(t, err)
		assertEventTableHandle(t, et, name, comment, nil)

		unset := sdk.NewEventTableUnsetRequest().WithComment(sdk.Bool(true))
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithUnset(unset))
		require.NoError(t, err)

		et, err = client.EventTables.ShowByID(ctx, id)
		require.NoError(t, err)
		assertEventTableHandle(t, et, name, "", nil)
	})

	t.Run("alter event table: set and unset change tacking", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		set := sdk.NewEventTableSetRequest().WithChangeTracking(sdk.Bool(true))
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithSet(set))
		require.NoError(t, err)

		et, err := client.EventTables.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, true, et.ChangeTracking)

		unset := sdk.NewEventTableUnsetRequest().WithChangeTracking(sdk.Bool(true))
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithUnset(unset))
		require.NoError(t, err)

		et, err = client.EventTables.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, false, et.ChangeTracking)
	})

	t.Run("alter event table: set and unset tag", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		set := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "v1",
			},
		}
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithSetTags(set))
		require.NoError(t, err)

		unset := []sdk.ObjectIdentifier{tagTest.ID()}
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithUnsetTags(unset))
		require.NoError(t, err)
	})

	t.Run("alter event table: rename", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)

		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupTableHandle(t, id))
		} else {
			t.Cleanup(cleanupTableHandle(t, nid))
		}
		require.NoError(t, err)

		_, err = client.EventTables.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		_, err = client.EventTables.ShowByID(ctx, nid)
		require.NoError(t, err)
	})

	t.Run("alter event table: clustering action with drop", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		err := client.EventTables.Create(ctx, sdk.NewCreateEventTableRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupTableHandle(t, id))

		action := sdk.NewEventTableClusteringActionRequest().WithDropClusteringKey(sdk.Bool(true))
		err = client.EventTables.Alter(ctx, sdk.NewAlterEventTableRequest(id).WithClusteringAction(action))
		require.NoError(t, err)
	})
}
