package testint

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DynamicTableCreateAndDrop(t *testing.T) {
	client := testClient(t)

	tableTest, tableCleanup := testClientHelper().Table.CreateTable(t, testSchema(t).ID())
	t.Cleanup(tableCleanup)

	ctx := context.Background()
	t.Run("test complete", func(t *testing.T) {
		name := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.String())
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, testWarehouse(t).ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
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
		require.Equal(t, testWarehouse(t).ID().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, entity.TargetLag)

		dynamicTableById, err := client.DynamicTables.ShowByID(ctx, name)
		require.NoError(t, err)
		require.NotNil(t, dynamicTableById)
		require.Equal(t, name.Name(), dynamicTableById.Name)
		require.Equal(t, testWarehouse(t).ID().Name(), dynamicTableById.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, dynamicTableById.TargetLag)
	})

	t.Run("test complete with target lag", func(t *testing.T) {
		name := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.String())
		targetLag := sdk.TargetLag{
			Downstream: sdk.Bool(true),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, testWarehouse(t).ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
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
		require.Equal(t, testWarehouse(t).ID().Name(), entity.Warehouse)
		require.Equal(t, "DOWNSTREAM", entity.TargetLag)
		require.Equal(t, sdk.DynamicTableRefreshModeIncremental, entity.RefreshMode)
		require.Contains(t, entity.Text, "initialize = 'ON_CREATE'")
		require.Contains(t, entity.Text, "refresh_mode = 'AUTO'")
	})

	t.Run("test complete with refresh mode and initialize", func(t *testing.T) {
		name := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.String())
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		refreshMode := sdk.DynamicTableRefreshModeFull
		initialize := sdk.DynamicTableInitializeOnSchedule
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, testWarehouse(t).ID(), targetLag, query).WithOrReplace(true).WithInitialize(initialize).WithRefreshMode(refreshMode).WithComment(&comment))
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
		require.Equal(t, testWarehouse(t).ID().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, entity.TargetLag)
		require.Equal(t, sdk.DynamicTableRefreshModeFull, entity.RefreshMode)
		require.Contains(t, entity.Text, "initialize = 'ON_SCHEDULE'")
		require.Contains(t, entity.Text, "refresh_mode = 'FULL'")
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
		require.Equal(t, sdk.DynamicTableSchedulingStateActive, entities[0].SchedulingState)

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
		require.Equal(t, sdk.DynamicTableSchedulingStateActive, entities[0].SchedulingState)
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
		sdk.ErrorsEqual(t, sdk.JoinErrors(sdk.ErrExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set")), err)
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

func TestInt_DynamicTablesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouseTest, databaseTest, schemaTest := testWarehouse(t), testDb(t), testSchema(t)

	cleanupDynamicTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createDynamicTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		tableTest, tableCleanup := testClientHelper().Table.CreateTable(t, schemaTest.ID())
		t.Cleanup(tableCleanup)
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(id, warehouseTest.ID(), targetLag, query).WithOrReplace(true))
		require.NoError(t, err)
		t.Cleanup(cleanupDynamicTableHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithIdentifier(t, databaseTest, random.AlphaN(8))
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createDynamicTableHandle(t, id1)
		createDynamicTableHandle(t, id2)

		e1, err := client.DynamicTables.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.DynamicTables.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
