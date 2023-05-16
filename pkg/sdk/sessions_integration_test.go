package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_UseWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	originalWH, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalWHIdentifier := NewAccountObjectIdentifier(originalWH)
		if !validObjectidentifier(originalWHIdentifier) {
			return
		}
		err := client.Sessions.UseWarehouse(ctx, originalWHIdentifier)
		require.NoError(t, err)
	})
	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	expected := warehouseTest.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseDatabase(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	originalDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalDBIdentifier := NewAccountObjectIdentifier(originalDB)
		if !validObjectidentifier(originalDBIdentifier) {
			return
		}
		err := client.Sessions.UseDatabase(ctx, originalDBIdentifier)
		require.NoError(t, err)
	})
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err = client.Sessions.UseDatabase(ctx, databaseTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	expected := databaseTest.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseSchema(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	originalSchema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	originalDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalSchemaIdentifier := NewSchemaIdentifier(originalDB, originalSchema)
		if !validObjectidentifier(originalSchemaIdentifier) {
			return
		}
		err := client.Sessions.UseSchema(ctx, originalSchemaIdentifier)
		require.NoError(t, err)
	})
	err = client.Sessions.UseSchema(ctx, schemaTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	expected := schemaTest.Name
	assert.Equal(t, expected, actual)
}
