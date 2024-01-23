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

// TODO [SNOW-1016430]: add tests for setting masking policy on creation
func TestInt_MaterializedViews(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	table, tableCleanup := createTable(t, client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	sql := fmt.Sprintf("SELECT id FROM %s", table.ID().FullyQualifiedName())

	assertMaterializedViewWithOptions := func(t *testing.T, view *sdk.MaterializedView, id sdk.SchemaObjectIdentifier, isSecure bool, comment string, clusterBy string) {
		t.Helper()
		assert.NotEmpty(t, view.CreatedOn)
		assert.Equal(t, id.Name(), view.Name)
		assert.Empty(t, view.Reserved)
		assert.Equal(t, testDb(t).Name, view.DatabaseName)
		assert.Equal(t, testSchema(t).Name, view.SchemaName)
		assert.Equal(t, clusterBy, view.ClusterBy)
		assert.Equal(t, 0, view.Rows)
		assert.Equal(t, 0, view.Bytes)
		assert.Equal(t, testDb(t).Name, view.SourceDatabaseName)
		assert.Equal(t, testSchema(t).Name, view.SourceSchemaName)
		assert.Equal(t, table.Name, view.SourceTableName)
		assert.NotEmpty(t, view.RefreshedOn)
		assert.NotEmpty(t, view.CompactedOn)
		assert.Equal(t, "ACCOUNTADMIN", view.Owner)
		assert.Equal(t, false, view.Invalid)
		assert.Equal(t, "", view.InvalidReason)
		assert.NotEmpty(t, view.BehindBy)
		assert.Equal(t, comment, view.Comment)
		assert.NotEmpty(t, view.Text)
		assert.Equal(t, isSecure, view.IsSecure)
		assert.Equal(t, clusterBy != "", view.AutomaticClustering)
		assert.Equal(t, "ROLE", view.OwnerRoleType)
		assert.Equal(t, "", view.Budget)
	}

	assertMaterializedView := func(t *testing.T, view *sdk.MaterializedView, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assertMaterializedViewWithOptions(t, view, id, false, "", "")
	}

	assertViewDetailsRow := func(t *testing.T, materializedViewDetails *sdk.MaterializedViewDetails) {
		t.Helper()
		assert.Equal(t, sdk.MaterializedViewDetails{
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
		}, *materializedViewDetails)
	}

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
		request := createMaterializedViewBasicRequest(t)

		view := createMaterializedViewWithRequest(t, request)

		assertMaterializedView(t, view, request.GetName())
	})

	t.Run("create materialized view: almost complete case", func(t *testing.T) {
		rowAccessPolicyId, rowAccessPolicyCleanup := createRowAccessPolicy(t, client, testSchema(t))
		t.Cleanup(rowAccessPolicyCleanup)

		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		request := createMaterializedViewBasicRequest(t).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithColumns([]sdk.MaterializedViewColumnRequest{
				*sdk.NewMaterializedViewColumnRequest("COLUMN_WITH_COMMENT").WithComment(sdk.String("column comment")),
			}).
			WithCopyGrants(sdk.Bool(true)).
			WithComment(sdk.String("comment")).
			WithRowAccessPolicy(sdk.NewMaterializedViewRowAccessPolicyRequest(rowAccessPolicyId, []string{"column_with_comment"})).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v2",
			}}).
			WithClusterBy(sdk.NewMaterializedViewClusterByRequest().WithExpressions([]sdk.MaterializedViewClusterByExpressionRequest{{"COLUMN_WITH_COMMENT"}}))

		id := request.GetName()

		view := createMaterializedViewWithRequest(t, request)

		assertMaterializedViewWithOptions(t, view, id, true, "comment", fmt.Sprintf(`LINEAR("%s")`, "COLUMN_WITH_COMMENT"))
		rowAccessPolicyReference, err := getRowAccessPolicyFor(t, client, view.ID(), sdk.ObjectTypeView)
		require.NoError(t, err)
		assert.Equal(t, rowAccessPolicyId.Name(), rowAccessPolicyReference.PolicyName)
		assert.Equal(t, "ROW_ACCESS_POLICY", rowAccessPolicyReference.PolicyKind)
		assert.Equal(t, view.ID().Name(), rowAccessPolicyReference.RefEntityName)
		assert.Equal(t, "MATERIALIZED_VIEW", rowAccessPolicyReference.RefEntityDomain)
		assert.Equal(t, "ACTIVE", rowAccessPolicyReference.PolicyStatus)
	})

	t.Run("drop materialized view: existing", func(t *testing.T) {
		request := createMaterializedViewBasicRequest(t)
		id := request.GetName()

		err := client.MaterializedViews.Create(ctx, request)
		require.NoError(t, err)

		err = client.MaterializedViews.Drop(ctx, sdk.NewDropMaterializedViewRequest(id))
		require.NoError(t, err)

		_, err = client.MaterializedViews.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop view: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		err := client.MaterializedViews.Drop(ctx, sdk.NewDropMaterializedViewRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter materialized view: rename", func(t *testing.T) {
		createRequest := createMaterializedViewBasicRequest(t)
		id := createRequest.GetName()

		err := client.MaterializedViews.Create(ctx, createRequest)
		require.NoError(t, err)

		newName := random.String()
		newId := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, newName)
		alterRequest := sdk.NewAlterMaterializedViewRequest(id).WithRenameTo(&newId)

		err = client.MaterializedViews.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupMaterializedViewProvider(id))
		} else {
			t.Cleanup(cleanupMaterializedViewProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.MaterializedViews.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		view, err := client.MaterializedViews.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertMaterializedView(t, view, newId)
	})

	t.Run("alter materialized view: set cluster by", func(t *testing.T) {
		view := createMaterializedView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterMaterializedViewRequest(id).WithClusterBy(sdk.NewMaterializedViewClusterByRequest().WithExpressions([]sdk.MaterializedViewClusterByExpressionRequest{{"ID"}}))
		err := client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf(`LINEAR("%s")`, "ID"), alteredView.ClusterBy)
	})

	t.Run("alter materialized view: recluster suspend and resume", func(t *testing.T) {
		request := createMaterializedViewBasicRequest(t).WithClusterBy(sdk.NewMaterializedViewClusterByRequest().WithExpressions([]sdk.MaterializedViewClusterByExpressionRequest{{"ID"}}))
		view := createMaterializedViewWithRequest(t, request)
		id := view.ID()

		assert.Equal(t, true, view.AutomaticClustering)

		alterRequest := sdk.NewAlterMaterializedViewRequest(id).WithSuspendRecluster(sdk.Bool(true))
		err := client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, false, alteredView.AutomaticClustering)

		alterRequest = sdk.NewAlterMaterializedViewRequest(id).WithResumeRecluster(sdk.Bool(true))
		err = client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, true, alteredView.AutomaticClustering)
	})

	t.Run("alter materialized view: suspend and resume", func(t *testing.T) {
		view := createMaterializedView(t)
		id := view.ID()

		assert.Equal(t, false, view.Invalid)

		alterRequest := sdk.NewAlterMaterializedViewRequest(id).WithSuspend(sdk.Bool(true))
		err := client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, true, alteredView.Invalid)

		alterRequest = sdk.NewAlterMaterializedViewRequest(id).WithResume(sdk.Bool(true))
		err = client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, false, alteredView.Invalid)
	})

	t.Run("alter materialized view: set and unset values", func(t *testing.T) {
		view := createMaterializedView(t)
		id := view.ID()

		alterRequest := sdk.NewAlterMaterializedViewRequest(id).WithSet(
			sdk.NewMaterializedViewSetRequest().WithSecure(sdk.Bool(true)),
		)
		err := client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err := client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, true, alteredView.IsSecure)

		alterRequest = sdk.NewAlterMaterializedViewRequest(id).WithSet(
			sdk.NewMaterializedViewSetRequest().WithComment(sdk.String("comment")),
		)
		err = client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "comment", alteredView.Comment)

		alterRequest = sdk.NewAlterMaterializedViewRequest(id).WithUnset(
			sdk.NewMaterializedViewUnsetRequest().WithComment(sdk.Bool(true)),
		)
		err = client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredView.Comment)

		alterRequest = sdk.NewAlterMaterializedViewRequest(id).WithUnset(
			sdk.NewMaterializedViewUnsetRequest().WithSecure(sdk.Bool(true)),
		)
		err = client.MaterializedViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredView, err = client.MaterializedViews.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, false, alteredView.IsSecure)
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
		view := createMaterializedView(t)

		returnedViewDetails, err := client.MaterializedViews.Describe(ctx, view.ID())
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedViewDetails))
		assertViewDetailsRow(t, &returnedViewDetails[0])
	})

	t.Run("describe materialized view: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		_, err := client.MaterializedViews.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
