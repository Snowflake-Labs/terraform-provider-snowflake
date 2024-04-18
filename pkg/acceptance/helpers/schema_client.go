package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SchemaClient struct {
	context *TestClientContext
}

func NewSchemaClient(context *TestClientContext) *SchemaClient {
	return &SchemaClient{
		context: context,
	}
}

func (c *SchemaClient) client() sdk.Schemas {
	return c.context.client.Schemas
}

func (c *SchemaClient) CreateSchema(t *testing.T, database *sdk.Database) (*sdk.Schema, func()) {
	t.Helper()
	return c.CreateSchemaWithIdentifier(t, database, random.StringRange(8, 28))
}

func (c *SchemaClient) CreateSchemaWithIdentifier(t *testing.T, database *sdk.Database, name string) (*sdk.Schema, func()) {
	t.Helper()
	ctx := context.Background()
	schemaID := sdk.NewDatabaseObjectIdentifier(database.Name, name)
	err := c.client().Create(ctx, schemaID, nil)
	require.NoError(t, err)
	schema, err := c.client().ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(database.Name, name))
	require.NoError(t, err)
	return schema, func() {
		err := c.client().Drop(ctx, schemaID, nil)
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}
		require.NoError(t, err)
		err = c.context.client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(c.context.database, c.context.schema))
		require.NoError(t, err)
	}
}
