package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CurrentSession(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	session, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session)
}

func TestInt_CurrentDatabase(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Sessions.UseDatabase(ctx, databaseTest.ID())
	require.NoError(t, err)
	db, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, db)
}

func TestInt_CurrentSchema(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	err := client.Sessions.UseSchema(ctx, schemaTest.ID())
	require.NoError(t, err)
	schema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
}

func TestInt_CurrentWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, warehouse)
}
