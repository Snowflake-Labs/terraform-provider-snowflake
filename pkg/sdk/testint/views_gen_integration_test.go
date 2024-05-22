package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1016430]: add tests for setting masking policy on creation
// TODO [SNOW-1016430]: add tests for setting recursive on creation
func TestInt_Views(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := testClientHelper().Table.CreateTable(t)
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
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

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

	// source https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2085
	t.Run("create view: no table reference", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateViewRequest(id, "SELECT NULL AS TYPE")

		view := createViewWithRequest(t, request)

		assertView(t, view, request.GetName())
	})

	t.Run("create view: almost complete case", func(t *testing.T) {
		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		request := createViewBasicRequest(t).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithTemporary(sdk.Bool(true)).
			WithColumns([]sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("COLUMN_WITH_COMMENT").WithComment(sdk.String("column comment")),
			}).
			WithCopyGrants(sdk.Bool(true)).
			WithComment(sdk.String("comment")).
			WithRowAccessPolicy(sdk.NewViewRowAccessPolicyRequest(rowAccessPolicy.ID(), []string{"column_with_comment"})).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v2",
			}})

		id := request.GetName()

		view := createViewWithRequest(t, request)

		assertViewWithOptions(t, view, id, true, "comment")
		rowAccessPolicyReference, err := testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy.Name, rowAccessPolicyReference.PolicyName)
		assert.Equal(t, "ROW_ACCESS_POLICY", rowAccessPolicyReference.PolicyKind)
		assert.Equal(t, view.ID().Name(), rowAccessPolicyReference.RefEntityName)
		assert.Equal(t, "VIEW", rowAccessPolicyReference.RefEntityDomain)
		assert.Equal(t, "ACTIVE", rowAccessPolicyReference.PolicyStatus)
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

	t.Run("alter view: rename", func(t *testing.T) {
		createRequest := createViewBasicRequest(t)
		id := createRequest.GetName()

		err := client.Views.Create(ctx, createRequest)
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterRequest := sdk.NewAlterViewRequest(id).WithRenameTo(&newId)

		err = client.Views.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupViewProvider(id))
		} else {
			t.Cleanup(cleanupViewProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.Views.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		view, err := client.Views.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertView(t, view, newId)
	})

	t.Run("alter view: set and unset values", func(t *testing.T) {
		view := createView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterViewRequest(id).WithSetComment(sdk.String("new comment"))
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredView.Comment)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetSecure(sdk.Bool(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, true, alteredView.IsSecure)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetChangeTracking(sdk.Bool(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "ON", alteredView.ChangeTracking)

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetComment(sdk.Bool(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredView.Comment)

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetSecure(sdk.Bool(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, false, alteredView.IsSecure)

		alterRequest = sdk.NewAlterViewRequest(id).WithSetChangeTracking(sdk.Bool(false))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.Views.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "OFF", alteredView.ChangeTracking)
	})

	t.Run("alter view: set and unset tag", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		view := createView(t)
		id := view.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterViewRequest(id).WithSetTags(tags)

		err := client.Views.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		// setting object type to view results in:
		// SQL compilation error: Invalid value VIEW for argument OBJECT_TYPE. Please use object type TABLE for all kinds of table-like objects.
		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTable)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterViewRequest(id).WithUnsetTags(unsetTags)

		err = client.Views.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTable)
		require.Error(t, err)
	})

	t.Run("alter view: set and unset masking policy", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicyIdentity(t, sdk.DataTypeNumber)
		t.Cleanup(maskingPolicyCleanup)

		view := createView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterViewRequest(id).WithSetMaskingPolicyOnColumn(
			sdk.NewViewSetColumnMaskingPolicyRequest("id", maskingPolicy.ID()),
		)
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredViewDetails, err := client.Views.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, 1, len(alteredViewDetails))
		assert.Equal(t, maskingPolicy.ID().FullyQualifiedName(), sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(*alteredViewDetails[0].PolicyName).FullyQualifiedName())

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetMaskingPolicyOnColumn(
			sdk.NewViewUnsetColumnMaskingPolicyRequest("id"),
		)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredViewDetails, err = client.Views.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, 1, len(alteredViewDetails))
		assert.Empty(t, alteredViewDetails[0].PolicyName)
	})

	t.Run("alter view: set and unset tags on column", func(t *testing.T) {
		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		view := createView(t)
		id := view.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}

		alterRequest := sdk.NewAlterViewRequest(id).WithSetTagsOnColumn(
			sdk.NewViewSetColumnTagsRequest("id", tags),
		)
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		columnId := sdk.NewTableColumnIdentifier(id.DatabaseName(), id.SchemaName(), id.Name(), "ID")
		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), columnId, sdk.ObjectTypeColumn)
		require.NoError(t, err)
		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}

		alterRequest = sdk.NewAlterViewRequest(id).WithUnsetTagsOnColumn(
			sdk.NewViewUnsetColumnTagsRequest("id", unsetTags),
		)
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), columnId, sdk.ObjectTypeColumn)
		require.Error(t, err)
	})

	t.Run("alter view: add and drop row access policies", func(t *testing.T) {
		rowAccessPolicy, rowAccessPolicyCleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicyCleanup)
		rowAccessPolicy2, rowAccessPolicy2Cleanup := testClientHelper().RowAccessPolicy.CreateRowAccessPolicy(t)
		t.Cleanup(rowAccessPolicy2Cleanup)

		view := createView(t)
		id := view.ID()

		// add policy
		alterRequest := sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []string{"ID"}))
		err := client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err := testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy.ID().Name(), rowAccessPolicyReference.PolicyName)
		assert.Equal(t, "ROW_ACCESS_POLICY", rowAccessPolicyReference.PolicyKind)
		assert.Equal(t, view.ID().Name(), rowAccessPolicyReference.RefEntityName)
		assert.Equal(t, "VIEW", rowAccessPolicyReference.RefEntityDomain)
		assert.Equal(t, "ACTIVE", rowAccessPolicyReference.PolicyStatus)

		// remove policy
		alterRequest = sdk.NewAlterViewRequest(id).WithDropRowAccessPolicy(sdk.NewViewDropRowAccessPolicyRequest(rowAccessPolicy.ID()))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		_, err = testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.Error(t, err, "no rows in result set")

		// add policy again
		alterRequest = sdk.NewAlterViewRequest(id).WithAddRowAccessPolicy(sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy.ID(), []string{"ID"}))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err = testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy.ID().Name(), rowAccessPolicyReference.PolicyName)

		// drop and add other policy simultaneously
		alterRequest = sdk.NewAlterViewRequest(id).WithDropAndAddRowAccessPolicy(sdk.NewViewDropAndAddRowAccessPolicyRequest(
			*sdk.NewViewDropRowAccessPolicyRequest(rowAccessPolicy.ID()),
			*sdk.NewViewAddRowAccessPolicyRequest(rowAccessPolicy2.ID(), []string{"ID"}),
		))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		rowAccessPolicyReference, err = testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicy2.ID().Name(), rowAccessPolicyReference.PolicyName)

		// drop all policies
		alterRequest = sdk.NewAlterViewRequest(id).WithDropAllRowAccessPolicies(sdk.Bool(true))
		err = client.Views.Alter(ctx, alterRequest)
		require.NoError(t, err)

		_, err = testClientHelper().RowAccessPolicy.GetRowAccessPolicyFor(t, view.ID(), sdk.ObjectTypeView)
		require.Error(t, err, "no rows in result set")
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

	// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2506
	t.Run("show view by id: same name in different schemas", func(t *testing.T) {
		// we assume that SF returns views alphabetically
		schemaName := "aaaa" + random.StringRange(8, 28)
		schema, schemaCleanup := testClientHelper().Schema.CreateSchemaWithName(t, schemaName)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		request1 := sdk.NewCreateViewRequest(id1, sql)
		request2 := sdk.NewCreateViewRequest(id2, sql)

		createViewWithRequest(t, request1)
		createViewWithRequest(t, request2)

		returnedView1, err := client.Views.ShowByID(ctx, id1)
		require.NoError(t, err)
		returnedView2, err := client.Views.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id1, returnedView1.ID())
		require.Equal(t, id2, returnedView2.ID())
	})

	t.Run("describe view", func(t *testing.T) {
		view := createView(t)

		returnedViewDetails, err := client.Views.Describe(ctx, view.ID())
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedViewDetails))
		assertViewDetailsRow(t, &returnedViewDetails[0])
	})

	t.Run("describe view: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		_, err := client.Views.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_ViewsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)
	table, tableCleanup := testClientHelper().Table.CreateTable(t)
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	cleanupViewHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Views.Drop(ctx, sdk.NewDropViewRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createViewHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.Views.Create(ctx, sdk.NewCreateViewRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupViewHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createViewHandle(t, id1)
		createViewHandle(t, id2)

		e1, err := client.Views.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Views.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
