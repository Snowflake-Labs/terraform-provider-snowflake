package testint

import (
	"context"
	"fmt"
	"testing"

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

	t.Run("CreateOnTable", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
		assert.Equal(t, "Table", *s.SourceType)
		assert.Equal(t, "DEFAULT", *s.Mode)
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		stageID := sdk.NewAccountObjectIdentifier("EXTERNAL_TABLE_STAGE")
		stageLocation := "@external_table_stage"
		_, _ = createStageWithURL(t, client, stageID, "s3://snowflake-workshop-lab/weather-nyc")

		externalTableId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		err := client.ExternalTables.Create(ctx, sdk.NewCreateExternalTableRequest(externalTableId, stageLocation, sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON)))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(externalTableId))
			require.NoError(t, err)
		})

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnExternalTableStreamRequest(id, externalTableId).WithInsertOnly(sdk.Bool(true)).WithComment(sdk.String("some comment"))
		err = client.Streams.CreateOnExternalTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, externalTableId.FullyQualifiedName(), *s.TableName)
		assert.Equal(t, "External Table", *s.SourceType)
		assert.Equal(t, "INSERT_ONLY", *s.Mode)
	})

	t.Run("CreateOnStage", func(t *testing.T) {
		stage, cleanupStage := createStageWithDirectory(t, client, db, schema, "test_stage")
		stageId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, stage.Name)
		t.Cleanup(cleanupStage)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnStageStreamRequest(id, stageId).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnStage(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, "Stage", *s.SourceType)
		assert.Equal(t, "DEFAULT", *s.Mode)
	})

	t.Run("CreateOnView", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		tableId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, table.Name)
		t.Cleanup(cleanupTable)

		viewId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		cleanupView := createView(t, client, viewId, fmt.Sprintf("SELECT id FROM %s", tableId.FullyQualifiedName()))
		t.Cleanup(cleanupView)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnViewStreamRequest(id, viewId).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnView(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, "View", *s.SourceType)
		assert.Equal(t, "DEFAULT", *s.Mode)
	})

	t.Run("Clone", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		cloneId := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		err = client.Streams.Clone(ctx, sdk.NewCloneStreamRequest(cloneId, id))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(cloneId))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(cloneId))
		require.NoError(t, err)
		assert.Equal(t, cloneId.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", *s.Comment)
		assert.Equal(t, table.ID().FullyQualifiedName(), *s.TableName)
		assert.Equal(t, "Table", *s.SourceType)
		assert.Equal(t, "DEFAULT", *s.Mode)
	})

	t.Run("Alter tags", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID())
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

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID())
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

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
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

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, sdk.NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.Show(ctx, sdk.NewShowStreamRequest().WithTerse(sdk.Bool(true)))
		require.NoError(t, err)

		stream, err := sdk.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assert.Equal(t, id.Name(), stream.Name)
		assert.Equal(t, db.Name, stream.DatabaseName)
		assert.Equal(t, schema.Name, stream.SchemaName)
		assert.Equal(t, table.Name, *stream.TableOn)
	})

	t.Run("Show", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
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

		stream, err := sdk.FindOne[sdk.Stream](s, func(stream sdk.Stream) bool { return id.Name() == stream.Name })
		require.NoError(t, err)

		assert.Equal(t, id.Name(), stream.Name)
		assert.Equal(t, db.Name, stream.DatabaseName)
		assert.Equal(t, schema.Name, stream.SchemaName)
		assert.Nil(t, stream.TableOn)
		assert.Equal(t, "some comment", *stream.Comment)
		assert.Equal(t, table.ID().FullyQualifiedName(), *stream.TableName)
		assert.Equal(t, "Table", *stream.SourceType)
		assert.Equal(t, "DEFAULT", *stream.Mode)
	})

	t.Run("Describe", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		id := sdk.NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := sdk.NewCreateOnTableStreamRequest(id, table.ID()).WithComment(sdk.String("some comment"))
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
