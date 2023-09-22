package sdk

import (
	"context"
	"testing"

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
		id := randomAccountObjectIdentifier(t)
		opts := &CreateDynamicTableOptions{
			OrReplace: Bool(true),
			TargetLag: TargetLag("2 minutes"),
			Query:     "select id from " + tableTest.ID().FullyQualifiedName(),
			Comment:   String("comment"),
		}
		err := client.DynamicTables.Create(ctx, id, warehouseTest.ID(), opts)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, id, &DropDynamicTableOptions{})
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, &ShowDynamicTableOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, id.Name(), entity.Name)
		require.Equal(t, warehouseTest.ID().Name(), entity.Warehouse)
		require.Equal(t, opts.TargetLag.String(), entity.TargetLag)
	})
}
