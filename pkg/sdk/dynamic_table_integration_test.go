package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DynamicTableCreate(t *testing.T) {
	client := testClient(t)

	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	tableTest, tableCleanup := createTable(t, client, databaseTest, schemaTest)
	t.Cleanup(tableCleanup)

	ctx := context.Background()
	t.Run("test complete", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		targetLag := TargetLag{
			Lagtime: String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := randomComment(t)
		err := client.DynamicTables.Create(ctx, NewCreateDynamicTableRequest(name, warehouseTest.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, NewDropDynamicTableRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, &ShowDynamicTableOptions{
			Like: &Like{
				Pattern: String(name.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, warehouseTest.ID().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.Lagtime, entity.TargetLag)
	})

	t.Run("test complete with target lag", func(t *testing.T) {
		name := randomAccountObjectIdentifier(t)
		targetLag := TargetLag{
			Downstream: Bool(true),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := randomComment(t)
		err := client.DynamicTables.Create(ctx, NewCreateDynamicTableRequest(name, warehouseTest.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, NewDropDynamicTableRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, &ShowDynamicTableOptions{
			Like: &Like{
				Pattern: String(name.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, warehouseTest.ID().Name(), entity.Warehouse)
		require.Equal(t, "DOWNSTREAM", entity.TargetLag)
	})
}

func TestInt_DynamicTableDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	dynamicTable, dynamicTableCleanup := createDynamicTable(t, client)
	t.Cleanup(dynamicTableCleanup)

	t.Run("when dynamic table exists", func(t *testing.T) {
		_, err := client.DynamicTables.Describe(ctx, NewDescribeDynamicTableRequest(dynamicTable.ID()))
		require.NoError(t, err)
	})

	t.Run("when dynamic table does not exist", func(t *testing.T) {
		name := NewAccountObjectIdentifier("does_not_exist")
		_, err := client.DynamicTables.Describe(ctx, NewDescribeDynamicTableRequest(name))
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})
}
