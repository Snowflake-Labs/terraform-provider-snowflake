package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupExternalFuncionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, dts).WithIfExists(sdk.Bool(true)))
			require.NoError(t, err)
		}
	}

	assertExternalFunction := func(t *testing.T, id sdk.SchemaObjectIdentifier, secure bool, withArguments bool) {
		t.Helper()

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)

		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, fmt.Sprintf(`"%v"`, id.SchemaName()), e.SchemaName)
		require.Equal(t, false, e.IsBuiltin)
		require.Equal(t, false, e.IsAggregate)
		require.Equal(t, false, e.IsAnsi)
		if withArguments {
			require.Equal(t, 1, e.MinNumArguments)
			require.Equal(t, 1, e.MaxNumArguments)
		} else {
			require.Equal(t, 0, e.MinNumArguments)
			require.Equal(t, 0, e.MaxNumArguments)
		}
		require.NotEmpty(t, e.Arguments)
		require.NotEmpty(t, e.Description)
		require.NotEmpty(t, e.CatalogName)
		require.Equal(t, false, e.IsTableFunction)
		require.Equal(t, false, e.ValidForClustering)
		require.Equal(t, secure, e.IsSecure)
		require.Equal(t, true, e.IsExternalFunction)
		require.Equal(t, "EXTERNAL", e.Language)
		require.Equal(t, false, e.IsMemoizable)
		require.Equal(t, false, e.IsDataMetric)
	}

	createApiIntegrationHandle := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE API INTEGRATION %s API_PROVIDER = aws_api_gateway API_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/hello_cloud_account_role' API_ALLOWED_PREFIXES = ('https://xyz.execute-api.us-west-2.amazonaws.com/production') ENABLED = true`, id.FullyQualifiedName()))
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err = client.ExecForTests(ctx, fmt.Sprintf(`DROP API INTEGRATION %s`, id.FullyQualifiedName()))
			require.NoError(t, err)
		})
	}

	createExternalFunction := func(t *testing.T) *sdk.ExternalFunction {
		t.Helper()

		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(4))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		argument := sdk.NewExternalFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument})
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFuncionHandle(id, []sdk.DataType{sdk.DataTypeVariant}))

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	t.Run("create external function", func(t *testing.T) {
		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(4))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		argument := sdk.NewExternalFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		headers := []sdk.ExternalFunctionHeaderRequest{
			{
				Name:  "measure",
				Value: "kilometers",
			},
		}
		ch := []sdk.ExternalFunctionContextHeaderRequest{
			{
				ContextFunction: "CURRENT_DATE",
			},
			{
				ContextFunction: "CURRENT_TIMESTAMP",
			},
		}
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument}).
			WithNullInputBehavior(sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput)).
			WithHeaders(headers).
			WithContextHeaders(ch).
			WithMaxBatchRows(sdk.Int(10)).
			WithCompression(sdk.String("GZIP"))
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFuncionHandle(id, []sdk.DataType{sdk.DataTypeVariant}))

		assertExternalFunction(t, id, true, true)
	})

	t.Run("create external function without arguments", func(t *testing.T) {
		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(4))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as)
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFuncionHandle(id, nil))

		assertExternalFunction(t, id, false, false)
	})

	t.Run("alter external function: set api integration", func(t *testing.T) {
		e := createExternalFunction(t)

		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(5))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().
			WithApiIntegration(&integration)
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true, true)
	})

	t.Run("alter external function: set headers", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		headers := []sdk.ExternalFunctionHeaderRequest{
			{
				Name:  "measure",
				Value: "kilometers",
			},
		}
		set := sdk.NewExternalFunctionSetRequest().WithHeaders(headers)
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, true)
	})

	t.Run("alter external function: set context headers", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		ch := []sdk.ExternalFunctionContextHeaderRequest{
			{
				ContextFunction: "CURRENT_DATE",
			},
			{
				ContextFunction: "CURRENT_TIMESTAMP",
			},
		}
		set := sdk.NewExternalFunctionSetRequest().WithContextHeaders(ch)
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, true)
	})

	t.Run("alter external function: set compression", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().WithCompression(sdk.String("AUTO"))
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, true)
	})

	t.Run("alter external function: set max batch rows", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().WithMaxBatchRows(sdk.Int(20))
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, true)
	})

	t.Run("alter external function: unset", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		unset := sdk.NewExternalFunctionUnsetRequest().
			WithComment(sdk.Bool(true)).
			WithHeaders(sdk.Bool(true))
		request := sdk.NewAlterExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR}).WithUnset(unset)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true, true)
	})

	t.Run("show external function: with like", func(t *testing.T) {
		e1 := createExternalFunction(t)
		e2 := createExternalFunction(t)

		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(&sdk.Like{Pattern: sdk.String(e1.Name)}))
		require.NoError(t, err)

		require.Equal(t, 1, len(es))
		require.Contains(t, es, *e1)
		require.NotContains(t, es, *e2)
	})

	t.Run("show external function: no matches", func(t *testing.T) {
		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(&sdk.Like{Pattern: sdk.String(random.String())}))
		require.NoError(t, err)
		require.Equal(t, 0, len(es))
	})
}
