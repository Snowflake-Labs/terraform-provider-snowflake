package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TableClient struct {
	context *TestClientContext
}

func NewTableClient(context *TestClientContext) *TableClient {
	return &TableClient{
		context: context,
	}
}

func (c *TableClient) client() sdk.Tables {
	return c.context.client.Tables
}

func (c *TableClient) CreateTable(t *testing.T) (*sdk.Table, func()) {
	t.Helper()
	return c.CreateTableInSchema(t, c.context.schemaId())
}

func (c *TableClient) CreateTableInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.Table, func()) {
	t.Helper()

	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
	}
	name := random.StringRange(8, 28)
	return c.CreateTableWithColumns(t, schemaId, name, columns)
}

func (c *TableClient) CreateTableWithColumns(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, name string, columns []sdk.TableColumnRequest) (*sdk.Table, func()) {
	t.Helper()

	id := sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), name)
	ctx := context.Background()

	dbCreateRequest := sdk.NewCreateTableRequest(id, columns)
	err := c.client().Create(ctx, dbCreateRequest)
	require.NoError(t, err)

	table, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return table, c.DropTableFunc(t, id)
}

func (c *TableClient) DropTableFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// to prevent error when schema was removed before the table
		_, err := c.context.client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()))
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}

		dropErr := c.client().Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, dropErr)
	}
}
