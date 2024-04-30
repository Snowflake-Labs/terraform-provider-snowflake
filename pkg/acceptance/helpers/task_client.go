package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TaskClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewTaskClient(context *TestClientContext, idsGenerator *IdsGenerator) *TaskClient {
	return &TaskClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *TaskClient) client() sdk.Tasks {
	return c.context.client.Tasks
}

func (c *TaskClient) defaultCreateTaskRequest(t *testing.T) *sdk.CreateTaskRequest {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	warehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(c.ids.WarehouseId()))
	return sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").WithWarehouse(warehouseReq)
}

func (c *TaskClient) CreateTask(t *testing.T) (*sdk.Task, func()) {
	t.Helper()
	return c.CreateTaskWithRequest(t, c.defaultCreateTaskRequest(t).WithSchedule(sdk.String("60 minutes")))
}

func (c *TaskClient) CreateTaskWithAfter(t *testing.T, taskId sdk.SchemaObjectIdentifier) (*sdk.Task, func()) {
	t.Helper()
	return c.CreateTaskWithRequest(t, c.defaultCreateTaskRequest(t).WithAfter([]sdk.SchemaObjectIdentifier{taskId}))
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
