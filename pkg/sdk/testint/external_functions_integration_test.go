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

	defaultDataTypes := []sdk.DataType{sdk.DataTypeVARCHAR}

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupExternalFuncionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, dts).WithIfExists(sdk.Bool(true)))
			require.NoError(t, err)
		}
	}

	assertExternalFunction := func(t *testing.T, id sdk.SchemaObjectIdentifier, secure bool, dts []sdk.DataType) {
		t.Helper()

		e, err := client.ExternalFunctions.ShowByID(ctx, id, dts)
		require.NoError(t, err)

		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, fmt.Sprintf(`"%v"`, id.SchemaName()), e.SchemaName)
		require.Equal(t, false, e.IsBuiltin)
		require.Equal(t, false, e.IsAggregate)
		require.Equal(t, false, e.IsAnsi)
		if len(dts) > 0 {
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

	createExternalFunction := func(t *testing.T, dt sdk.DataType) *sdk.ExternalFunction {
		t.Helper()

		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(4))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		argument := sdk.NewExternalFunctionArgumentRequest("x", dt)
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument})
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFuncionHandle(id, []sdk.DataType{sdk.DataTypeVariant}))

		e, err := client.ExternalFunctions.ShowByID(ctx, id, defaultDataTypes)
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

		assertExternalFunction(t, id, true, defaultDataTypes)
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

		assertExternalFunction(t, id, false, nil)
	})

	t.Run("alter external function: set api integration", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(5))
		createApiIntegrationHandle(t, integration)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().
			WithApiIntegration(&integration)
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("alter external function: set headers", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		headers := []sdk.ExternalFunctionHeaderRequest{
			{
				Name:  "measure",
				Value: "kilometers",
			},
		}
		set := sdk.NewExternalFunctionSetRequest().WithHeaders(headers)
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("alter external function: set context headers", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

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
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("alter external function: set compression", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().WithCompression(sdk.String("AUTO"))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("alter external function: set max batch rows", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		set := sdk.NewExternalFunctionSetRequest().WithMaxBatchRows(sdk.Int(20))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("alter external function: unset", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		unset := sdk.NewExternalFunctionUnsetRequest().
			WithComment(sdk.Bool(true)).
			WithHeaders(sdk.Bool(true))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithUnset(unset)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true, defaultDataTypes)
	})

	t.Run("show external function: with like", func(t *testing.T) {
		e1 := createExternalFunction(t, sdk.DataTypeVARCHAR)
		e2 := createExternalFunction(t, sdk.DataTypeVARCHAR)

		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(&sdk.Like{Pattern: sdk.String(e1.Name)}))
		require.NoError(t, err)

		require.Equal(t, 1, len(es))
		require.Contains(t, es, *e1)
		require.NotContains(t, es, *e2)
	})

	t.Run("show external function: with in", func(t *testing.T) {
		otherDb, otherDbCleanup := createDatabase(t, testClient(t))
		t.Cleanup(otherDbCleanup)

		e1 := createExternalFunction(t, sdk.DataTypeVARCHAR)

		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(databaseTest.Name, schemaTest.Name)}))
		require.NoError(t, err)

		require.Contains(t, es, *e1)

		es, err = client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(&sdk.In{Database: sdk.NewAccountObjectIdentifier(databaseTest.Name)}))
		require.NoError(t, err)

		require.Contains(t, es, *e1)

		es, err = client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithIn(&sdk.In{Database: otherDb.ID()}))
		require.NoError(t, err)

		require.Empty(t, es)
	})

	t.Run("show external function: no matches", func(t *testing.T) {
		es, err := client.ExternalFunctions.Show(ctx, sdk.NewShowExternalFunctionRequest().WithLike(&sdk.Like{Pattern: sdk.String(random.String())}))
		require.NoError(t, err)
		require.Equal(t, 0, len(es))
	})

	t.Run("show external function by id", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		es, err := client.ExternalFunctions.ShowByID(ctx, id, []sdk.DataType{sdk.DataTypeVARCHAR})
		require.NoError(t, err)
		require.Equal(t, *e, *es)

		_, err = client.ExternalFunctions.ShowByID(ctx, id, nil)
		require.Error(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("describe external function", func(t *testing.T) {
		e := createExternalFunction(t, sdk.DataTypeVARCHAR)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		request := sdk.NewDescribeExternalFunctionRequest(id, []sdk.DataType{sdk.DataTypeVARCHAR})
		details, err := client.ExternalFunctions.Describe(ctx, request)
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, "EXTERNAL", pairs["language"])
		require.Equal(t, "VARIANT", pairs["returns"])
		require.Equal(t, "VOLATILE", pairs["volatility"])
		require.Equal(t, "AUTO", pairs["compression"])
		require.Equal(t, "(X VARCHAR)", pairs["signature"])
	})
}

func TestInt_ExternalFunctionsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupExternalFuncionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, dts).WithIfExists(sdk.Bool(true)))
			require.NoError(t, err)
		}
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

	createExternalFunction := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		integration := sdk.NewAccountObjectIdentifier(random.AlphaN(4))
		createApiIntegrationHandle(t, integration)

		argument := sdk.NewExternalFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument})
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFuncionHandle(id, []sdk.DataType{sdk.DataTypeVariant}))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := createSchemaWithIdentifier(t, client, databaseTest, random.AlphaN(8))
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createExternalFunction(t, id1)
		createExternalFunction(t, id2)

		e1, err := client.ExternalFunctions.ShowByID(ctx, id1, []sdk.DataType{sdk.DataTypeVARCHAR})
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.ExternalFunctions.ShowByID(ctx, id2, []sdk.DataType{sdk.DataTypeVARCHAR})
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
