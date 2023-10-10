package sdk

import (
	"testing"
)

func TestInt_Streams(t *testing.T) {
	//client := testClient(t)
	//ctx := context.Background()
	//
	//db, cleanupDb := createDatabase(t, client)
	//t.Cleanup(cleanupDb)
	//
	//schema, cleanupSchema := createSchema(t, client, db)
	//t.Cleanup(cleanupSchema)
	//
	//t.Run("CreateOnTable", func(t *testing.T) {
	//	table, cleanupTable := createTable(t, client, db, schema)
	//	t.Cleanup(cleanupTable)
	//
	//	id := randomAccountObjectIdentifier(t)
	//	err := client.Streams.CreateOnTable(ctx, NewCreateOnTableStreamRequest(id, table.ID()))
	//	require.NoError(t, err)
	//	t.Cleanup(func() {
	//		err := client.Streams.Drop(ctx, NewDropStreamRequest(id))
	//		require.NoError(t, err)
	//	})
	//
	//})

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
