package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DynamicTableCreateAndDrop(t *testing.T) {
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
		name := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, warehouseTest.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(name.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, warehouseTest.ID().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, entity.TargetLag)
	})

	t.Run("test complete with target lag", func(t *testing.T) {
		name := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		targetLag := sdk.TargetLag{
			Downstream: sdk.Bool(true),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, warehouseTest.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(name.Name())}))
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
		_, err := client.DynamicTables.Describe(ctx, sdk.NewDescribeDynamicTableRequest(dynamicTable.ID()))
		require.NoError(t, err)
	})

	t.Run("when dynamic table does not exist", func(t *testing.T) {
		name := sdk.NewSchemaObjectIdentifier("my_db", "my_schema", "does_not_exist")
		_, err := client.DynamicTables.Describe(ctx, sdk.NewDescribeDynamicTableRequest(name))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_DynamicTableAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("alter with suspend or resume", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := createDynamicTable(t, client)
		t.Cleanup(dynamicTableCleanup)

		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateRunning, entities[0].SchedulingState)

		err = client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSuspend(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err = client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateSuspended, entities[0].SchedulingState)

		err = client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithResume(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err = client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateRunning, entities[0].SchedulingState)
	})

	t.Run("alter with refresh", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := createDynamicTable(t, client)
		t.Cleanup(dynamicTableCleanup)

		err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithRefresh(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
	})

	t.Run("alter with suspend and resume", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := createDynamicTable(t, client)
		t.Cleanup(dynamicTableCleanup)

		err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSuspend(sdk.Bool(true)).WithResume(sdk.Bool(true)))
		require.Error(t, err)
		expected := "alter statement needs exactly one action from: set, unset, refresh"
		require.Equal(t, expected, err.Error())
	})

	t.Run("alter with set", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := createDynamicTable(t, client)
		t.Cleanup(dynamicTableCleanup)

		targetLagCases := []string{"10 minutes", "DOWNSTREAM"}
		for _, value := range targetLagCases {
			err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSet(sdk.NewDynamicTableSetRequest().WithTargetLag(sdk.TargetLag{MaximumDuration: sdk.String(value)})))
			require.NoError(t, err)
			entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
			require.NoError(t, err)
			require.Equal(t, 1, len(entities))
			require.Equal(t, value, entities[0].TargetLag)
		}
	})
}
