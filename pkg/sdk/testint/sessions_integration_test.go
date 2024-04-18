package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AlterSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	opts := &sdk.AlterSessionOptions{
		Set: &sdk.SessionSet{
			SessionParameters: &sdk.SessionParameters{
				AbortDetachedQuery:    sdk.Bool(true),
				Autocommit:            sdk.Bool(true),
				GeographyOutputFormat: sdk.Pointer(sdk.GeographyOutputFormatGeoJSON),
				WeekOfYearPolicy:      sdk.Int(1),
			},
		},
	}
	err := client.Sessions.AlterSession(ctx, opts)
	require.NoError(t, err)
	cleanup := func() {
		opts = &sdk.AlterSessionOptions{
			Unset: &sdk.SessionUnset{
				SessionParametersUnset: &sdk.SessionParametersUnset{
					AbortDetachedQuery:    sdk.Bool(true),
					Autocommit:            sdk.Bool(true),
					GeographyOutputFormat: sdk.Bool(true),
					WeekOfYearPolicy:      sdk.Bool(true),
				},
			},
		}
		err := client.Sessions.AlterSession(ctx, opts)
		require.NoError(t, err)
	}
	t.Cleanup(cleanup)

	parameter, err := client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameterAbortDetachedQuery)
	require.NoError(t, err)
	assert.Equal(t, "true", parameter.Value)
	parameter, err = client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameterAutocommit)
	require.NoError(t, err)
	assert.Equal(t, "true", parameter.Value)
	parameter, err = client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameterGeographyOutputFormat)
	require.NoError(t, err)
	assert.Equal(t, string(sdk.GeographyOutputFormatGeoJSON), parameter.Value)
	parameter, err = client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameterWeekOfYearPolicy)
	require.NoError(t, err)
	assert.Equal(t, "1", parameter.Value)
}

func TestInt_ShowParameters(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameters, err := client.Parameters.ShowParameters(ctx, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, parameters)
}

func TestInt_ShowAccountParameter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterAutocommit)
	require.NoError(t, err)
	assert.NotEmpty(t, parameter)
}

func TestInt_ShowSessionParameter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameter, err := client.Parameters.ShowSessionParameter(ctx, sdk.SessionParameterAutocommit)
	require.NoError(t, err)
	assert.NotEmpty(t, parameter)
}

func TestInt_ShowObjectParameter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	parameter, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterDataRetentionTimeInDays, sdk.Object{ObjectType: testDb(t).ObjectType(), Name: testDb(t).ID()})
	require.NoError(t, err)
	assert.NotEmpty(t, parameter)
}

func TestInt_ShowUserParameter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	userID := sdk.NewAccountObjectIdentifier(user)
	parameter, err := client.Parameters.ShowUserParameter(ctx, sdk.UserParameterAutocommit, userID)
	require.NoError(t, err)
	assert.NotEmpty(t, parameter)
}

func TestInt_UseWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Cleanup(func() {
		err := client.Sessions.UseWarehouse(ctx, testWarehouse(t).ID())
		require.NoError(t, err)
	})
	// new warehouse created on purpose
	warehouse, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, warehouse.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	expected := warehouse.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseDatabase(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Cleanup(func() {
		err := client.Sessions.UseSchema(ctx, testSchema(t).ID())
		require.NoError(t, err)
	})
	// new database created on purpose
	database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	err := client.Sessions.UseDatabase(ctx, database.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	expected := database.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseSchema(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Cleanup(func() {
		err := client.Sessions.UseSchema(ctx, testSchema(t).ID())
		require.NoError(t, err)
	})
	// new database and schema created on purpose
	database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t, database)
	t.Cleanup(schemaCleanup)
	err := client.Sessions.UseSchema(ctx, schema.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	expected := schema.Name
	assert.Equal(t, expected, actual)
}
