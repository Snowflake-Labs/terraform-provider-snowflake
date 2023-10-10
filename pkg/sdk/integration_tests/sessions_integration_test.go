package sdk_integration_tests

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AlterSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	opts := &sdk.AlterSessionOptions{
		Set: &sdk.SessionSet{
			&sdk.SessionParameters{
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
				&sdk.SessionParametersUnset{
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
	databaseTest, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	parameter, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterDataRetentionTimeInDays, sdk.Object{ObjectType: databaseTest.ObjectType(), Name: databaseTest.ID()})
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
	originalWH, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalWHIdentifier := sdk.NewAccountObjectIdentifier(originalWH)
		if !sdk.validObjectidentifier(originalWHIdentifier) {
			return
		}
		err := client.Sessions.UseWarehouse(ctx, originalWHIdentifier)
		require.NoError(t, err)
	})
	warehouseTest, warehouseCleanup := sdk.createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err = client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	expected := warehouseTest.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseDatabase(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	originalDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalDBIdentifier := sdk.NewAccountObjectIdentifier(originalDB)
		if !sdk.validObjectidentifier(originalDBIdentifier) {
			return
		}
		err := client.Sessions.UseDatabase(ctx, originalDBIdentifier)
		require.NoError(t, err)
	})
	databaseTest, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err = client.Sessions.UseDatabase(ctx, databaseTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	expected := databaseTest.Name
	assert.Equal(t, expected, actual)
}

func TestInt_UseSchema(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := sdk.createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	originalSchema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	originalDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		originalSchemaIdentifier := sdk.NewDatabaseObjectIdentifier(originalDB, originalSchema)
		if !sdk.validObjectidentifier(originalSchemaIdentifier) {
			return
		}
		err := client.Sessions.UseSchema(ctx, originalSchemaIdentifier)
		require.NoError(t, err)
	})
	err = client.Sessions.UseSchema(ctx, schemaTest.ID())
	require.NoError(t, err)
	actual, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	expected := schemaTest.Name
	assert.Equal(t, expected, actual)
}
