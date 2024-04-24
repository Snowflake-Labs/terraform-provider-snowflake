package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type AlertClient struct {
	context *TestClientContext
}

func NewAlertClient(context *TestClientContext) *AlertClient {
	return &AlertClient{
		context: context,
	}
}

func (c *AlertClient) client() sdk.Alerts {
	return c.context.client.Alerts
}

func (c *AlertClient) CreateAlert(t *testing.T) (*sdk.Alert, func()) {
	t.Helper()
	schedule := "USING CRON * * * * * UTC"
	condition := "SELECT 1"
	action := "SELECT 1"
	return c.CreateAlertWithOptions(t, schedule, condition, action, &sdk.CreateAlertOptions{})
}

func (c *AlertClient) CreateAlertWithOptions(t *testing.T, schedule string, condition string, action string, opts *sdk.CreateAlertOptions) (*sdk.Alert, func()) {
	t.Helper()
	ctx := context.Background()

	name := random.String()
	id := sdk.NewSchemaObjectIdentifier(c.context.database, c.context.schema, name)

	err := c.client().Create(ctx, id, c.context.warehouseId(), schedule, condition, action, opts)
	require.NoError(t, err)

	alert, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return alert, c.DropAlertFunc(t, id)
}

func (c *AlertClient) DropAlertFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// TODO: add if exists
		err := c.client().Drop(ctx, id)
		require.NoError(t, err)
	}
}
