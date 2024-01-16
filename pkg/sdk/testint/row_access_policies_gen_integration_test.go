package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_RowAccessPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

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

		// TODO: args should be required
		return sdk.NewCreateRowAccessPolicyRequest(id, body).WithArgs(args)
	}

	createRowAccessPolicyBasicRequest := func(t *testing.T) *sdk.CreateRowAccessPolicyRequest {
		t.Helper()

		argName := random.AlphaN(5)
		argType := string(sdk.DataTypeVARCHAR)
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
		// TODO: fill me
	})

	t.Run("create row access policy: full", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop row access policy: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop row access policy: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: rename", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset comment", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset body", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter row access policy: set and unset tags", func(t *testing.T) {
		// TODO: fill me
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
		argType := string(sdk.DataTypeVARCHAR)
		args := sdk.NewCreateRowAccessPolicyArgsRequest(argName, argType)
		body := "true"

		request := createRowAccessPolicyRequest(t, []sdk.CreateRowAccessPolicyArgsRequest{*args}, body)
		rowAccessPolicy := createRowAccessPolicyWithRequest(t, request)

		returnedRowAccessPolicyDescription, err := client.RowAccessPolicies.Describe(ctx, rowAccessPolicy.ID())
		require.NoError(t, err)

		assertRowAccessPolicyDescription(t, returnedRowAccessPolicyDescription, rowAccessPolicy.ID(), fmt.Sprintf("(%s %s)", strings.ToUpper(argName), argType), body)
	})

	t.Run("describe row access policy: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, "does_not_exist")

		_, err := client.RowAccessPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
