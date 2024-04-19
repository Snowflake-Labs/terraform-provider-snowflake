package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type DatabaseRoleClient struct {
	context *TestClientContext
}

func NewDatabaseRoleClient(context *TestClientContext) *DatabaseRoleClient {
	return &DatabaseRoleClient{
		context: context,
	}
}

func (c *DatabaseRoleClient) client() sdk.DatabaseRoles {
	return c.context.client.DatabaseRoles
}

func (c *DatabaseRoleClient) CreateDatabaseRole(t *testing.T, databaseId sdk.AccountObjectIdentifier) (*sdk.DatabaseRole, func()) {
	t.Helper()
	ctx := context.Background()

	name := random.String()
	id := sdk.NewDatabaseObjectIdentifier(databaseId.Name(), name)

	err := c.client().Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
	require.NoError(t, err)

	databaseRole, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return databaseRole, c.CleanupDatabaseRoleFunc(t, id)
}

func (c *DatabaseRoleClient) CleanupDatabaseRoleFunc(t *testing.T, id sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// to prevent error when db was removed before the role
		_, err := c.context.client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id.DatabaseName()))
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}

		err = c.client().Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
