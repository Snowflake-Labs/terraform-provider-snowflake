package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type RoleClient struct {
	context *TestClientContext
}

func NewRoleClient(context *TestClientContext) *RoleClient {
	return &RoleClient{
		context: context,
	}
}

func (c *RoleClient) client() sdk.Roles {
	return c.context.client.Roles
}

func (c *RoleClient) UseRole(t *testing.T, roleName string) func() {
	t.Helper()
	ctx := context.Background()

	currentRole, err := c.context.client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	err = c.context.client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier(roleName))
	require.NoError(t, err)

	return func() {
		err = c.context.client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier(currentRole))
		require.NoError(t, err)
	}
}

func (c *RoleClient) CreateRole(t *testing.T) (*sdk.Role, func()) {
	t.Helper()
	return c.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(sdk.RandomAccountObjectIdentifier()))
}

func (c *RoleClient) CreateRoleGrantedToCurrentUser(t *testing.T) (*sdk.Role, func()) {
	t.Helper()
	ctx := context.Background()
	role, roleCleanup := c.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(sdk.RandomAccountObjectIdentifier()))

	currentUser, err := c.context.client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = c.client().Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{
		User: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentUser)),
	}))
	require.NoError(t, err)

	return role, roleCleanup
}

func (c *RoleClient) CreateRoleWithRequest(t *testing.T, req *sdk.CreateRoleRequest) (*sdk.Role, func()) {
	t.Helper()
	ctx := context.Background()

	require.True(t, sdk.ValidObjectIdentifier(req.GetName()))
	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	role, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return role, func() {
		err = c.client().Drop(ctx, sdk.NewDropRoleRequest(req.GetName()))
		require.NoError(t, err)
	}
}

// TODO: drop role func
// TODO: clean above methods
