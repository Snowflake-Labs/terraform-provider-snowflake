package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Streams(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	db, cleanupDb := createDatabase(t, client)
	t.Cleanup(cleanupDb)

	schema, cleanupSchema := createSchema(t, client, db)
	t.Cleanup(cleanupSchema)

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
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)

		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
		assertStream(t, s, id, "Table", "DEFAULT")
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		stageID := sdk.NewAccountObjectIdentifier("EXTERNAL_TABLE_STAGE")
		stageLocation := "@external_table_stage"
		_, _ = createStageWithURL(t, client, stageID, nycWeatherDataURL)

		externalTableId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		err := client.ExternalTables.Create(ctx, sdk.NewCreateExternalTableRequest(externalTableId, stageLocation, sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON)))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(externalTableId))
			require.NoError(t, err)
		})

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnExternalTableRequest(id, externalTableId).WithInsertOnly(sdk.Bool(true)).WithComment(sdk.String("some comment"))
		err = client.Streams.CreateOnExternalTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)

		assert.Equal(t, externalTableId.FullyQualifiedName(), *s.TableName)
		assertStream(t, s, id, "External Table", "INSERT_ONLY")
	})

	t.Run("CreateOnDirectoryTable", func(t *testing.T) {
		stage, cleanupStage := createStageWithDirectory(t, client, db, schema, "test_stage")
		stageId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, stage.Name)
		t.Cleanup(cleanupStage)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnDirectoryTableRequest(id, stageId).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnDirectoryTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)

		assertStream(t, s, id, "Stage", "DEFAULT")
	})

	t.Run("CreateOnView", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		tableId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, table.Name)
		t.Cleanup(cleanupTable)

		viewId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		cleanupView := createView(t, client, viewId, fmt.Sprintf("SELECT id FROM %s", tableId.FullyQualifiedName()))
		t.Cleanup(cleanupView)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnViewRequest(id, viewId).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnView(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)

		assertStream(t, s, id, "View", "DEFAULT")
	})

	t.Run("Clone", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		cloneId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		err = client.Streams.Clone(ctx, sdk.NewCloneStreamRequest(cloneId, id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(cloneId))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(cloneId))
		require.NoError(t, err)

		assertStream(t, s, cloneId, "Table", "DEFAULT")
		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
	})

	t.Run("Alter tags", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID())
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		tag, cleanupTag := createTag(t, client, db, schema)
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

		_, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
	})

	t.Run("Alter comment", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID())
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(sdk.String("some_comment")))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, "some_comment", *s.Comment)

		err = client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(sdk.Bool(true)))
		require.NoError(t, err)

		s, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, "", *s.Comment)

		_, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)

		err = client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
		require.NoError(t, err)

		_, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.Error(t, err)
	})

	t.Run("Show terse", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
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
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
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
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		id2 := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
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
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		idPrefix := "stream_show_"

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, idPrefix+random.AlphanumericN(32))
		req := sdk.NewCreateStreamOnTableRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		id2 := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, idPrefix+random.AlphanumericN(32))
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
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, random.AlphanumericN(32))
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
