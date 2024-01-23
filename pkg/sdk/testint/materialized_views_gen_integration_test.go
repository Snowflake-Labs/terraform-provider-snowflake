package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_MaterializedViews(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := createTable(t, client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	cleanupMaterializedViewProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.MaterializedViews.Drop(ctx, sdk.NewDropMaterializedViewRequest(id))
			require.NoError(t, err)
		}
	}

	createMaterializedViewBasicRequest := func(t *testing.T) *sdk.CreateMaterializedViewRequest {
		t.Helper()
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		return sdk.NewCreateMaterializedViewRequest(id, sql)
	}

	createMaterializedViewWithRequest := func(t *testing.T, request *sdk.CreateMaterializedViewRequest) *sdk.MaterializedView {
		t.Helper()
		id := request.GetName()

		err := client.MaterializedViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupMaterializedViewProvider(id))

		materializedView, err := client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		return materializedView
	}

	createMaterializedView := func(t *testing.T) *sdk.MaterializedView {
		t.Helper()
		return createMaterializedViewWithRequest(t, createMaterializedViewBasicRequest(t))
	}

	t.Run("create materialized view: no optionals", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("create materialized view: complete case", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop materialized view: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop view: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: rename", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: set cluster by", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: recluster suspend and resume", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: suspend and resume", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter materialized view: set and unset values", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show materialized view: default", func(t *testing.T) {
		view1 := createMaterializedView(t)
		view2 := createMaterializedView(t)

		showRequest := sdk.NewShowMaterializedViewRequest()
		returnedViews, err := client.MaterializedViews.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedViews))
		assert.Contains(t, returnedViews, *view1)
		assert.Contains(t, returnedViews, *view2)
	})

	t.Run("show materialized view: with options", func(t *testing.T) {
		view1 := createMaterializedView(t)
		view2 := createMaterializedView(t)

		showRequest := sdk.NewShowMaterializedViewRequest().
			WithLike(&sdk.Like{Pattern: &view1.Name}).
			WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(testDb(t).Name, testSchema(t).Name)})
		returnedViews, err := client.MaterializedViews.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedViews))
		assert.Contains(t, returnedViews, *view1)
		assert.NotContains(t, returnedViews, *view2)
	})

	t.Run("describe materialized view", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe materialized view: non-existing", func(t *testing.T) {
		// TODO: fill me
	})
}
