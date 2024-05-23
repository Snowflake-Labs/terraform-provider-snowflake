package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Streams(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db := testDb(t)
	schema := testSchema(t)

	assertStream := func(t *testing.T, s *sdk.Stream, id sdk.SchemaObjectIdentifier, sourceType string, mode string) {
		t.Helper()
		assert.NotNil(t, s)
		assert.Nil(t, s.TableOn)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, sourceType, *s.SourceType)
		assert.Equal(t, mode, *s.Mode)
	}

	t.Run("CreateOnTable", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
		assertStream(t, s, id, "Table", "DEFAULT")
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		stageID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
		_, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t, stageID)
		t.Cleanup(stageCleanup)

		externalTableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.ExternalTables.Create(ctx, sdk.NewCreateExternalTableRequest(externalTableId, stageLocation).WithFileFormat(*sdk.NewExternalTableFileFormatRequest().WithFileFormatType(sdk.ExternalTableFileFormatTypeJSON)))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(externalTableId))
			require.NoError(t, err)
		})

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnExternalTableRequest(id, externalTableId).WithInsertOnly(sdk.Bool(true)).WithComment(sdk.String("some comment"))
		err = client.Streams.CreateOnExternalTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, externalTableId.FullyQualifiedName(), *s.TableName)
		assertStream(t, s, id, "External Table", "INSERT_ONLY")
	})

	t.Run("CreateOnDirectoryTable", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStageWithDirectory(t)
		stageId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, stage.Name)
		t.Cleanup(cleanupStage)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnDirectoryTableRequest(id, stageId).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnDirectoryTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStream(t, s, id, "Stage", "DEFAULT")
	})

	t.Run("CreateOnView", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		tableId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, table.Name)
		t.Cleanup(cleanupTable)

		view, cleanupView := testClientHelper().View.CreateView(t, fmt.Sprintf("SELECT id FROM %s", tableId.FullyQualifiedName()))
		t.Cleanup(cleanupView)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnViewRequest(id, view.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnView(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		assertStream(t, s, id, "View", "DEFAULT")
	})

	t.Run("Clone", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		cloneId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err = client.Streams.Clone(ctx, sdk.NewCloneStreamRequest(cloneId, id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(cloneId))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, cloneId)
		require.NoError(t, err)

		assertStream(t, s, cloneId, "Table", "DEFAULT")
		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
	})

	t.Run("Alter tags", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID())
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		tag, cleanupTag := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(cleanupTag)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStream)
		require.Error(t, err)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: "tag_value",
			},
		}))
		require.NoError(t, err)

		tagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStream)
		require.NoError(t, err)
		assert.Equal(t, "tag_value", tagValue)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetTags([]sdk.ObjectIdentifier{tag.ID()}))
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeStream)
		require.Error(t, err)

		_, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Alter comment", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID())
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(sdk.String("some_comment")))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "some_comment", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(sdk.Bool(true)))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		_, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, id)
		require.NoError(t, err)

		err = client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("Show terse", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().WithTerse(sdk.Bool(true)))
		require.NoError(t, err)

		stream, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assert.NotNil(t, stream)
		assert.Equal(t, id.Name(), stream.Name)
		assert.Equal(t, db.Name, stream.DatabaseName)
		assert.Equal(t, schema.Name, stream.SchemaName)
		assert.Equal(t, table.Name, *stream.TableOn)
		assert.Nil(t, stream.Comment)
		assert.Nil(t, stream.SourceType)
		assert.Nil(t, stream.Mode)
	})

	t.Run("Show single with options", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().
			WithTerse(sdk.Bool(false)).
			WithIn(&sdk.In{
				Schema: sdk.NewDatabaseObjectIdentifier(db.Name, schema.Name),
			}).
			WithLike(&sdk.Like{
				Pattern: sdk.String(id.Name()),
			}).
			WithStartsWith(sdk.String(id.Name())).
			WithLimit(&sdk.LimitFrom{
				Rows: sdk.Int(1),
			}))
		require.NoError(t, err)
		assert.Equal(t, 1, len(s))

		stream, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assertStream(t, stream, id, "Table", "DEFAULT")
	})

	t.Run("Show multiple", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req2 := sdk.NewCreateStreamOnTableRequest(id2, table.ID()).WithComment(sdk.String("some comment"))
		err = client.Streams.CreateOnTable(ctx, req2)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id2))
			require.NoError(t, err)
		})

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest())
		require.NoError(t, err)
		assert.Equal(t, 2, len(s))

		stream, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		stream2, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id2.Name() == stream.Name })
		require.NoError(t, err)

		assertStream(t, stream, id, "Table", "DEFAULT")
		assertStream(t, stream2, id2, "Table", "DEFAULT")
	})

	t.Run("Show multiple with options", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		idPrefix := "stream_show_"

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix(idPrefix)
		req2 := sdk.NewCreateStreamOnTableRequest(id2, table.ID()).WithComment(sdk.String("some comment"))
		err = client.Streams.CreateOnTable(ctx, req2)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id2))
			require.NoError(t, err)
		})

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().
			WithTerse(sdk.Bool(false)).
			WithIn(&sdk.In{
				Schema: sdk.NewDatabaseObjectIdentifier(db.Name, schema.Name),
			}).
			WithStartsWith(sdk.String(idPrefix)).
			WithLimit(&sdk.LimitFrom{
				Rows: sdk.Int(2),
			}))
		require.NoError(t, err)
		assert.Equal(t, 2, len(s))

		stream, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		stream2, err := collections.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id2.Name() == stream.Name })
		require.NoError(t, err)

		assertStream(t, stream, id, "Table", "DEFAULT")
		assertStream(t, stream2, id2, "Table", "DEFAULT")
	})

	t.Run("Describe", func(t *testing.T) {
		table, cleanupTable := testClientHelper().Table.CreateTableInSchema(t, schema.ID())
		t.Cleanup(cleanupTable)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.Describe(ctx, sdk.NewDescribeStreamRequest(id))
		require.NoError(t, err)
		assert.NotNil(t, s)

		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Nil(t, s.TableOn)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
		assert.Equal(t, "Table", *s.SourceType)
		assert.Equal(t, "DEFAULT", *s.Mode)
	})
}

func TestInt_StreamsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, cleanupTable := testClientHelper().Table.CreateTable(t)
	t.Cleanup(cleanupTable)

	cleanupStreamHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createStreamHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.Streams.CreateOnTable(ctx, sdk.NewCreateStreamOnTableRequest(id, table.ID()))
		require.NoError(t, err)
		t.Cleanup(cleanupStreamHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createStreamHandle(t, id1)
		createStreamHandle(t, id2)

		e1, err := client.Streams.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Streams.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
