package sdk

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

		id := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := NewCreateOnTableStreamRequest(id, table.ID()).WithComment(String("some comment"))
		err := client.Streams.CreateOnTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", s.Comment)
		assert.Equal(t, table.ID().FullyQualifiedName(), s.TableName)
		assert.Equal(t, "Table", s.SourceType)
		assert.Equal(t, "DEFAULT", s.Mode)
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		stageID := NewAccountObjectIdentifier("EXTERNAL_TABLE_STAGE")
		stageLocation := "@external_table_stage"
		_, _ = createStageWithURL(t, client, stageID, "s3://snowflake-workshop-lab/weather-nyc")

		externalTableId := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		err := client.ExternalTables.Create(ctx, NewCreateExternalTableRequest(externalTableId, stageLocation, NewExternalTableFileFormatRequest().WithFileFormatType(&ExternalTableFileFormatTypeJSON)))
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.ExternalTables.Drop(ctx, NewDropExternalTableRequest(externalTableId))
			require.NoError(t, err)
		})

		id := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := NewCreateOnExternalTableStreamRequest(id, externalTableId).WithInsertOnly(Bool(true)).WithComment(String("some comment"))
		err = client.Streams.CreateOnExternalTable(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", s.Comment)
		assert.Equal(t, externalTableId.FullyQualifiedName(), s.TableName)
		assert.Equal(t, "External Table", s.SourceType)
		assert.Equal(t, "INSERT_ONLY", s.Mode)
	})

	t.Run("CreateOnStage", func(t *testing.T) {
		stage, cleanupStage := createStageWithDirectory(t, client, db, schema, "test_stage", true)
		stageId := NewSchemaObjectIdentifier(db.Name, schema.Name, stage.Name)
		t.Cleanup(cleanupStage)

		id := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := NewCreateOnStageStreamRequest(id, stageId).WithComment(String("some comment"))
		err := client.Streams.CreateOnStage(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", s.Comment)
		assert.Equal(t, "Stage", s.SourceType)
		assert.Equal(t, "DEFAULT", s.Mode)
	})

	t.Run("CreateOnView", func(t *testing.T) {
		table, cleanupTable := createTable(t, client, db, schema)
		t.Cleanup(cleanupTable)

		viewId := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		cleanupView := createView(t, client, viewId, "")
		t.Cleanup(cleanupView)

		id := NewSchemaObjectIdentifier(db.Name, schema.Name, randomAlphanumericN(t, 32))
		req := NewCreateOnStageStreamRequest(id, stageId).WithComment(String("some comment"))
		err := client.Streams.CreateOnStage(ctx, req)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := client.Streams.Drop(ctx, NewDropStreamRequest(id))
			require.NoError(t, err)
		})

		s, err := client.Streams.ShowByID(ctx, NewShowByIdStreamRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), s.Name)
		assert.Equal(t, db.Name, s.DatabaseName)
		assert.Equal(t, schema.Name, s.SchemaName)
		assert.Equal(t, "some comment", s.Comment)
		assert.Equal(t, "Stage", s.SourceType)
		assert.Equal(t, "DEFAULT", s.Mode)
	})

	t.Run("Clone", func(t *testing.T) {
		// TODO: fill me
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

	t.Run("ShowByID", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe", func(t *testing.T) {
		// TODO: fill me
	})
}
