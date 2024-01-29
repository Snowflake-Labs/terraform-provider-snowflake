package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_RowAccessPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertRowAccessPolicy := func(t *testing.T, rowAccessPolicy *sdk.RowAccessPolicy, id sdk.SchemaObjectIdentifier, comment string) {
		t.Helper()
		assert.NotEmpty(t, rowAccessPolicy.CreatedOn)
		assert.Equal(t, id.Name(), rowAccessPolicy.Name)
		assert.Equal(t, id.DatabaseName(), rowAccessPolicy.DatabaseName)
		assert.Equal(t, id.SchemaName(), rowAccessPolicy.SchemaName)
		assert.Equal(t, "ROW_ACCESS_POLICY", rowAccessPolicy.Kind)
		assert.Equal(t, "ACCOUNTADMIN", rowAccessPolicy.Owner)
		assert.Equal(t, comment, rowAccessPolicy.Comment)
		assert.Empty(t, rowAccessPolicy.Options)
		assert.Equal(t, "ROLE", rowAccessPolicy.OwnerRoleType)
	}

	assertRowAccessPolicyDescription := func(t *testing.T, rowAccessPolicyDescription *sdk.RowAccessPolicyDescription, id sdk.SchemaObjectIdentifier, expectedSignature string, expectedBody string) {
		t.Helper()
		assert.Equal(t, sdk.RowAccessPolicyDescription{
			Name:       id.Name(),
			Signature:  expectedSignature,
			ReturnType: "BOOLEAN",
			Body:       expectedBody,
		}, *rowAccessPolicyDescription)
	}

	cleanupRowAccessPolicyProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
			require.NoError(t, err)
		}
	}

	createRowAccessPolicyRequest := func(t *testing.T, args []sdk.CreateRowAccessPolicyArgsRequest, body string) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		return sdk.NewCreateRowAccessPolicyRequest(id, args, body)
	}

	createRowAccessPolicyBasicRequest := func(t *testing.T) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()

		argName := random.AlphaN(5)
		argType := sdk.DataTypeVARCHAR
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)

		body := "true"

		return createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
	}

	createRowAccessPolicyWithRequest := func(t *testing.T, request *sdk.CreateRowAccessPolicyRequest) *sdk.RowAccessPolicy {
		t.Helper()
		id := request.GetName()

		err := client.RowAccessPolicies.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupRowAccessPolicyProvider(id))

		rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		return rowAccessPolicy
	}

	createRowAccessPolicy := func(t *testing.T) *sdk.RowAccessPolicy {
		t.Helper()
		return createRowAccessPolicyWithRequest(t, createRowAccessPolicyBasicRequest(t))
	}

	t.Run("create row access policy: no optionals", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)

		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		assertRowAccessPolicy(t, rowAccessPolicy, request.GetName(), "")
	})

	t.Run("create row access policy: full", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)
		request.Comment = sdk.String("some comment")

		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		assertRowAccessPolicy(t, rowAccessPolicy, request.GetName(), "some comment")
	})

	t.Run("drop row access policy: existing", func(t *testing.T) {
		request := createRowAccessPolicyBasicRequest(t)
		id := request.GetName()

		err := client.RowAccessPolicies.Create(ctx, request)
		require.NoError(t, err)

		err = client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.RowAccessPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop row access policy: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter row access policy: rename", func(t *testing.T) {
		createRequest := createRowAccessPolicyBasicRequest(t)
		id := createRequest.GetName()

		err := client.RowAccessPolicies.Create(ctx, createRequest)
		require.NoError(t, err)

		newName := random.String()
		newId := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, newName)
		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithRenameTo(&newId)

		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupRowAccessPolicyProvider(id))
		} else {
			t.Cleanup(cleanupRowAccessPolicyProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.RowAccessPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		rowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertRowAccessPolicy(t, rowAccessPolicy, newId, "")
	})

	t.Run("alter row access policy: set and unset comment", func(t *testing.T) {
		rowAccessPolicy := createRowAccessPolicy(t)
		id := rowAccessPolicy.ID()

		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithSetComment(sdk.String("new comment"))
		err := client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredRowAccessPolicy.Comment)

		alterRequest = sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetComment(sdk.Bool(true))
		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicy, err = client.RowAccessPolicies.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredRowAccessPolicy.Comment)
	})

	t.Run("alter row access policy: set body", func(t *testing.T) {
		rowAccessPolicy := createRowAccessPolicy(t)
		id := rowAccessPolicy.ID()

		alterRequest := sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String("false"))
		err := client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "false", alteredRowAccessPolicyDescription.Body)

		alterRequest = sdk.NewAlterRowAccessPolicyRequest(id).WithSetBody(sdk.String("true"))
		err = client.RowAccessPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredRowAccessPolicyDescription, err = client.RowAccessPolicies.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "true", alteredRowAccessPolicyDescription.Body)
	})

	t.Run("alter row access policy: set and unset tags", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		rowAccessPolicy := createRowAccessPolicy(t)
		id := rowAccessPolicy.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterRowAccessPolicyRequest(id).WithSetTags(tags)

		err := client.RowAccessPolicies.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeRowAccessPolicy)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterRowAccessPolicyRequest(id).WithUnsetTags(unsetTags)

		err = client.RowAccessPolicies.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeRowAccessPolicy)
		require.Error(t, err)
	})

	t.Run("show row access policy: default", func(t *testing.T) {
		rowAccessPolicy1 := createRowAccessPolicy(t)
		rowAccessPolicy2 := createRowAccessPolicy(t)

		showRequest := sdk.NewShowRowAccessPolicyRequest()
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("show row access policy: with options", func(t *testing.T) {
		rowAccessPolicy1 := createRowAccessPolicy(t)
		rowAccessPolicy2 := createRowAccessPolicy(t)

		showRequest := sdk.NewShowRowAccessPolicyRequest().
			WithLike(&sdk.Like{Pattern: &rowAccessPolicy1.Name}).
			WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(testDb(t).Name, testSchema(t).Name)})
		returnedRowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedRowAccessPolicies))
		assert.Contains(t, returnedRowAccessPolicies, *rowAccessPolicy1)
		assert.NotContains(t, returnedRowAccessPolicies, *rowAccessPolicy2)
	})

	t.Run("describe row access policy: existing", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataTypeVARCHAR
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), argType), body)
	})

	t.Run("describe row access policy: with data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataTypeTimestamp
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), sdk.DataTypeTimestampNTZ), body)
	})

	t.Run("describe row access policy: with data type normalization", func(t *testing.T) {
		argName := random.AlphaN(5)
		argType := sdk.DataType("VARCHAR(200)")
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), sdk.DataTypeVARCHAR), body)
	})

	t.Run("describe row access policy: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		_, err := client.RowAccessPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
