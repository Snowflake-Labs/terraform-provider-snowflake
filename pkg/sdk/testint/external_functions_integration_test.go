package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	defaultDataTypes := []sdk.DataType{sdk.DataTypeVARCHAR}

	databaseTest, schemaTest := testDb(t), testSchema(t)

	integration, integrationCleanup := createApiIntegration(t, client)
	t.Cleanup(integrationCleanup)

	cleanupExternalFunctionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id.WithoutArguments(), dts).WithIfExists(sdk.Bool(true)))
			require.NoError(t, err)
		}
	}

	createExternalFunction := func(t *testing.T) *sdk.ExternalFunction {
		t.Helper()
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, random.StringN(4), defaultDataTypes)
		argument := sdk.NewExternalFunctionArgumentRequest("x", defaultDataTypes[0])
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as).
			WithOrReplace(sdk.Bool(true)).
			WithSecure(sdk.Bool(true)).
			WithArguments([]sdk.ExternalFunctionArgumentRequest{*argument})
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFunctionHandle(id.WithoutArguments(), []sdk.DataType{sdk.DataTypeVariant}))

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	assertExternalFunction := func(t *testing.T, id sdk.SchemaObjectIdentifier, secure bool) {
		t.Helper()
		dts := id.Arguments()

		e, err := client.ExternalFunctions.ShowByID(ctx, id)
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

	t.Run("create external function", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, random.StringN(4), defaultDataTypes)
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
		t.Cleanup(cleanupExternalFunctionHandle(id.WithoutArguments(), []sdk.DataType{sdk.DataTypeVariant}))

		assertExternalFunction(t, id, true)
	})

	t.Run("create external function without arguments", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, random.StringN(4), nil)
		as := "https://xyz.execute-api.us-west-2.amazonaws.com/production/remote_echo"
		request := sdk.NewCreateExternalFunctionRequest(id, sdk.DataTypeVariant, &integration, as)
		err := client.ExternalFunctions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalFunctionHandle(id, nil))

		assertExternalFunction(t, id, false)
	})

	t.Run("alter external function: set api integration", func(t *testing.T) {
		e := createExternalFunction(t)
		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
		set := sdk.NewExternalFunctionSetRequest().
			WithApiIntegration(&integration)
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true)
	})

	t.Run("alter external function: set headers", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
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
		assertExternalFunction(t, id, true)
	})

	t.Run("alter external function: set context headers", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
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
		assertExternalFunction(t, id, true)
	})

	t.Run("alter external function: set compression", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
		set := sdk.NewExternalFunctionSetRequest().WithCompression(sdk.String("AUTO"))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true)
	})

	t.Run("alter external function: set max batch rows", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
		set := sdk.NewExternalFunctionSetRequest().WithMaxBatchRows(sdk.Int(20))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithSet(set)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)
		assertExternalFunction(t, id, true)
	})

	t.Run("alter external function: unset", func(t *testing.T) {
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
		unset := sdk.NewExternalFunctionUnsetRequest().
			WithComment(sdk.Bool(true)).
			WithHeaders(sdk.Bool(true))
		request := sdk.NewAlterExternalFunctionRequest(id, defaultDataTypes).WithUnset(unset)
		err := client.ExternalFunctions.Alter(ctx, request)
		require.NoError(t, err)

		assertExternalFunction(t, id, true)
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

	t.Run("show external function: with in", func(t *testing.T) {
		otherDb, otherDbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(otherDbCleanup)

		e1 := createExternalFunction(t)

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
		e := createExternalFunction(t)

		id := sdk.NewSchemaObjectIdentifierWithArguments(databaseTest.Name, schemaTest.Name, e.Name, defaultDataTypes)
		es, err := client.ExternalFunctions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, *e, *es)

		_, err = client.ExternalFunctions.ShowByID(ctx, id.WithoutArguments())
		require.Error(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("describe external function", func(t *testing.T) {
		e := createExternalFunction(t)
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
