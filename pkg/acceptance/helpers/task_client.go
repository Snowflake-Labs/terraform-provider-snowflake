package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TaskClient struct {
	context *TestClientContext
}

func NewTaskClient(context *TestClientContext) *TaskClient {
	return &TaskClient{
		context: context,
	}
}

func (c *TaskClient) client() sdk.Tasks {
	return c.context.client.Tasks
}

func (c *TaskClient) CreateTask(t *testing.T) (*sdk.Task, func()) {
	t.Helper()
	id := c.context.newSchemaObjectIdentifier(random.AlphanumericN(12))
	warehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(c.context.warehouseId()))
	return c.CreateTaskWithRequest(t, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").WithSchedule(sdk.String("60 minutes")).WithWarehouse(warehouseReq))
}

func (c *TaskClient) CreateTaskWithRequest(t *testing.T, request *sdk.CreateTaskRequest) (*sdk.Task, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.GetName()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	task, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return task, c.DropTaskFunc(t, id)
}

func (c *TaskClient) DropTaskFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropTaskRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
