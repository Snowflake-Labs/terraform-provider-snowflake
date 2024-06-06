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

func TestInt_CortexSearchServiceCreateAndDrop(t *testing.T) {
	client := testClient(t)

	tableTest, tableCleanup := testClientHelper().Table.CreateTable(t)
	t.Cleanup(tableCleanup)

	ctx := context.Background()
	t.Run("test complete", func(t *testing.T) {
		name := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		targetLag := "2 minutes"
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(name, "id", testWarehouse(t).ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.CortexSearchServices.Drop(ctx, sdk.NewDropCortexSearchServiceRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.CortexSearchServices.Show(ctx, sdk.NewShowCortexSearchServiceRequest().WithLike(&sdk.Like{Pattern: sdk.String(name.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)

		cortexSearchServiceById, err := client.CortexSearchServices.ShowByID(ctx, name)
		require.NoError(t, err)
		require.NotNil(t, cortexSearchServiceById)
		require.Equal(t, name.Name(), cortexSearchServiceById.Name)
		require.Equal(t, name.DatabaseName(), cortexSearchServiceById.DatabaseName)
		require.Equal(t, name.SchemaName(), cortexSearchServiceById.SchemaName)
	})
}

func TestInt_CortexSearchServiceDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	table, tableCleanup := testClientHelper().Table.CreateTable(t)
	t.Cleanup(tableCleanup)

	cortexSearchService, cortexSearchServiceCleanup := testClientHelper().CortexSearchService.CreateCortexSearchService(t, table.ID())
	t.Cleanup(cortexSearchServiceCleanup)

	t.Run("when cortex search service exists", func(t *testing.T) {
		_, err := client.CortexSearchServices.Describe(ctx, sdk.NewDescribeCortexSearchServiceRequest(cortexSearchService.ID()))
		require.NoError(t, err)
	})

	t.Run("when cortex search service does not exist", func(t *testing.T) {
		_, err := client.CortexSearchServices.Describe(ctx, sdk.NewDescribeCortexSearchServiceRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_CortexSearchServiceAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	table, tableCleanup := testClientHelper().Table.CreateTable(t)
	t.Cleanup(tableCleanup)

	t.Run("alter with set", func(t *testing.T) {
		cortexSearchService, cortexSearchServiceCleanup := testClientHelper().CortexSearchService.CreateCortexSearchService(t, table.ID())
		t.Cleanup(cortexSearchServiceCleanup)

		err := client.CortexSearchServices.Alter(ctx, sdk.NewAlterCortexSearchServiceRequest(cortexSearchService.ID()).WithSet(sdk.NewCortexSearchServiceSetRequest().WithTargetLag("10 minutes")))
		require.NoError(t, err)
		entities, err := client.CortexSearchServices.Show(ctx, sdk.NewShowCortexSearchServiceRequest().WithLike(&sdk.Like{Pattern: sdk.String(cortexSearchService.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
	})
}

func TestInt_CortexSearchServicesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouseTest := testWarehouse(t)

	cleanupCortexSearchServiceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.CortexSearchServices.Drop(ctx, sdk.NewDropCortexSearchServiceRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createCortexSearchServiceHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		tableTest, tableCleanup := testClientHelper().Table.CreateTable(t)
		t.Cleanup(tableCleanup)
		on := "ID"
		targetLag := "2 minutes"
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		err := client.CortexSearchServices.Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, on, warehouseTest.ID(), targetLag, query).WithOrReplace(true))
		require.NoError(t, err)
		t.Cleanup(cleanupCortexSearchServiceHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createCortexSearchServiceHandle(t, id1)
		createCortexSearchServiceHandle(t, id2)

		e1, err := client.CortexSearchServices.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.CortexSearchServices.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
