package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Views(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := createTable(t, client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	// TODO: fill
	assertView := func(t *testing.T, task *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
	}

	// TODO: fill
	assertViewWithOptions := func(t *testing.T, task *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
	}

	// TODO: fill
	assertViewTerse := func(t *testing.T, task *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
	}

	cleanupViewProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
			require.NoError(t, err)
		}
	}

	createViewBasicRequest := func(t *testing.T) *sdk.CreateViewRequest {
		t.Helper()
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		return sdk.NewCreateViewRequest(id, sql)
	}

	createViewWithRequest := func(t *testing.T, request *sdk.CreateViewRequest) *sdk.View {
		t.Helper()
		id := request.GetName()

		err := client.Views.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupViewProvider(id))

		view, err := client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		return view
	}

	createView := func(t *testing.T) *sdk.View {
		t.Helper()
		return createViewWithRequest(t, createViewBasicRequest(t))
	}

	_, _ = assertViewTerse, createView

	t.Run("create view: no optionals", func(t *testing.T) {
		request := createViewBasicRequest(t)

		view := createViewWithRequest(t, request)

		assertView(t, view, request.GetName())
	})

	t.Run("create view: almost complete case", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		// row access policy is not added
		// masking policy is not added
		// recursive is not used
		request := createViewBasicRequest(t).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithTemporary(sdk.Bool(true)).
			WithColumns([]sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("column_with_comment").WithComment(sdk.String("column comment")),
			}).
			WithCopyGrants(sdk.Bool(true)).
			WithComment(sdk.String("comment")).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v2",
			}})

		id := request.GetName()

		view := createViewWithRequest(t, request)

		assertViewWithOptions(t, view, id)
	})

	t.Run("drop view: existing", func(t *testing.T) {
		request := createViewBasicRequest(t)
		id := request.GetName()

		err := client.Views.Create(ctx, request)
		require.NoError(t, err)

		err = client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
		require.NoError(t, err)

		_, err = client.Views.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop view: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
