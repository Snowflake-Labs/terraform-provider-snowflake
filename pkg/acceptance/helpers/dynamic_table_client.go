package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type DynamicTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDynamicTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *DynamicTableClient {
	return &DynamicTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DynamicTableClient) client() sdk.DynamicTables {
	return c.context.client.DynamicTables
}

func (c *DynamicTableClient) CreateDynamicTable(t *testing.T, tableId sdk.SchemaObjectIdentifier) (*sdk.DynamicTable, func()) {
	t.Helper()
	return c.CreateDynamicTableWithOptions(t, c.ids.SchemaId(), random.AlphaN(12), c.ids.WarehouseId(), tableId)
}

func (c *DynamicTableClient) CreateDynamicTableWithOptions(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, name string, warehouseId sdk.AccountObjectIdentifier, tableId sdk.SchemaObjectIdentifier) (*sdk.DynamicTable, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifierInSchema(schemaId, name)
	targetLag := sdk.TargetLag{
		MaximumDuration: sdk.String("2 minutes"),
	}
	query := fmt.Sprintf(`select "ID" from %s`, tableId.FullyQualifiedName())
	comment := random.Comment()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateDynamicTableRequest(id, warehouseId, targetLag, query).WithComment(&comment))
	require.NoError(t, err)

	dynamicTable, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return dynamicTable, c.DropDynamicTableFunc(t, id)
}

func (c *DynamicTableClient) DropDynamicTableFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropDynamicTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
