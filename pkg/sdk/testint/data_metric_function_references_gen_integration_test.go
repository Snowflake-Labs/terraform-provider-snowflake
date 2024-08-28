package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_DataMetricFunctionReferences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("view domain", func(t *testing.T) {
		functionId := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "CORE", "AVG")
		statement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
		view, viewCleanup := testClientHelper().View.CreateView(t, statement)
		t.Cleanup(viewCleanup)

		// when we specify schedule by a number of minutes, a cron is returned from Snowflake - see SNOW-1640024
		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(view.ID()).WithSetDataMetricSchedule(*sdk.NewViewSetDataMetricScheduleRequest("5 MINUTE")))
		require.NoError(t, err)
		err = client.Views.Alter(ctx, sdk.NewAlterViewRequest(view.ID()).WithAddDataMetricFunction(*sdk.NewViewAddDataMetricFunctionRequest([]sdk.ViewDataMetricFunction{{
			DataMetricFunction: functionId,
			On:                 []sdk.Column{{Value: "ROLE_NAME"}},
		}})))
		require.NoError(t, err)

		dmfs, err := client.DataMetricFunctionReferences.GetForEntity(ctx, sdk.NewGetForEntityDataMetricFunctionReferenceRequest(view.ID(), sdk.DataMetricFuncionRefEntityDomainView))
		require.NoError(t, err)
		require.Equal(t, 1, len(dmfs))
		dmf := dmfs[0]
		require.Equal(t, string(sdk.DataMetricFuncionRefEntityDomainView), dmf.RefEntityDomain)
		require.Equal(t, functionId.DatabaseName(), dmf.MetricDatabaseName)
		require.Equal(t, functionId.SchemaName(), dmf.MetricSchemaName)
		require.Equal(t, functionId.Name(), dmf.MetricName)
		require.Equal(t, view.ID().DatabaseName(), dmf.RefEntityDatabaseName)
		require.Equal(t, view.ID().SchemaName(), dmf.RefEntitySchemaName)
		require.Equal(t, view.ID().Name(), dmf.RefEntityName)
		require.Equal(t, "*/5 * * * * UTC", dmf.Schedule)
	})
}
