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

	assertViewWithOptions := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier, isSecure bool, comment string) {
		t.Helper()
		assert.NotEmpty(t, view.CreatedOn)
		assert.Equal(t, id.Name(), view.Name)
		// Kind is filled out only in TERSE response.
		assert.Empty(t, view.Kind)
		assert.Empty(t, view.Reserved)
		assert.Equal(t, testDb(t).Name, view.DatabaseName)
		assert.Equal(t, testSchema(t).Name, view.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", view.Owner)
		assert.Equal(t, comment, view.Comment)
		assert.NotEmpty(t, view.Text)
		assert.Equal(t, isSecure, view.IsSecure)
		assert.Equal(t, false, view.IsMaterialized)
		assert.Equal(t, "ROLE", view.OwnerRoleType)
		assert.Equal(t, "OFF", view.ChangeTracking)
	}

	assertView := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assertViewWithOptions(t, view, id, false, "")
	}

	assertViewTerse := func(t *testing.T, view *sdk.View, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.NotEmpty(t, view.CreatedOn)
		assert.Equal(t, id.Name(), view.Name)
		assert.Equal(t, "VIEW", view.Kind)
		assert.Equal(t, testDb(t).Name, view.DatabaseName)
		assert.Equal(t, testSchema(t).Name, view.SchemaName)

		// all below are not contained in the terse response, that's why all of them we expect to be empty
		assert.Empty(t, view.Reserved)
		assert.Empty(t, view.Owner)
		assert.Empty(t, view.Comment)
		assert.Empty(t, view.Text)
		assert.Empty(t, view.IsSecure)
		assert.Empty(t, view.IsMaterialized)
		assert.Empty(t, view.OwnerRoleType)
		assert.Empty(t, view.ChangeTracking)
	}

	assertViewDetailsRow := func(t *testing.T, viewDetails *sdk.ViewDetails) {
		t.Helper()
		assert.Equal(t, sdk.ViewDetails{
			Name:       "ID",
			Type:       "NUMBER(38,0)",
			Kind:       "COLUMN",
			IsNullable: true,
			Default:    nil,
			IsPrimary:  false,
			IsUnique:   false,
			Check:      nil,
			Expression: nil,
			Comment:    nil,
			PolicyName: nil,
		}, *viewDetails)
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

		assertViewWithOptions(t, view, id, true, "comment")
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

	t.Run("show view: default", func(t *testing.T) {
		view1 := createView(t)
		view2 := createView(t)

		showRequest := sdk.NewShowViewRequest()
		returnedViews, err := client.Views.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedViews))
		assert.Contains(t, returnedViews, *view1)
		assert.Contains(t, returnedViews, *view2)
	})

	t.Run("show view: terse", func(t *testing.T) {
		view := createView(t)

		showRequest := sdk.NewShowViewRequest().WithTerse(sdk.Bool(true))
		returnedViews, err := client.Views.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedViews))
		assertViewTerse(t, &returnedViews[0], view.ID())
	})

	t.Run("show view: with options", func(t *testing.T) {
		view1 := createView(t)
		view2 := createView(t)

		showRequest := sdk.NewShowViewRequest().
			WithLike(&sdk.Like{Pattern: &view1.Name}).
			WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(testDb(t).Name, testSchema(t).Name)}).
			WithLimit(&sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedViews, err := client.Views.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedViews))
		assert.Contains(t, returnedViews, *view1)
		assert.NotContains(t, returnedViews, *view2)
	})

	t.Run("describe view", func(t *testing.T) {
		view := createView(t)

		returnedViewDetails, err := client.Views.Describe(ctx, view.ID())
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedViewDetails))
		assertViewDetailsRow(t, &returnedViewDetails[0])
	})
}
