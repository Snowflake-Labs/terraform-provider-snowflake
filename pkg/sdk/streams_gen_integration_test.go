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
		table, cleanupTable := createTableWithTimeTravel(t, client, db, schema, 1)
		t.Cleanup(cleanupTable)

		id := randomAccountObjectIdentifier(t)
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
		assert.Equal(t, table.Name, s.TableName)
		assert.Equal(t, "TABLE", s.SourceType)
		assert.Equal(t, "DEFAULT", s.Mode)
	})

	t.Run("CreateOnExternalTable", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateOnStage", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("CreateOnView", func(t *testing.T) {
		// TODO: fill me
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
