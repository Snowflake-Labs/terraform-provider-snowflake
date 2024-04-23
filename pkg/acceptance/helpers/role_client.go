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

func (c *RoleClient) CreateRoleWithName(t *testing.T, name string) (*sdk.Role, func()) {
	t.Helper()
	return c.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(sdk.NewAccountObjectIdentifier(name)))
}

func (c *RoleClient) CreateRoleGrantedToCurrentUser(t *testing.T) (*sdk.Role, func()) {
	t.Helper()

	role, roleCleanup := c.CreateRole(t)
	c.GrantRoleToCurrentUser(t, role.ID())
	return role, roleCleanup
}

func (c *RoleClient) CreateRoleWithRequest(t *testing.T, req *sdk.CreateRoleRequest) (*sdk.Role, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	role, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return role, c.DropRoleFunc(t, req.GetName())
}

func (c *RoleClient) DropRoleFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropRoleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *RoleClient) GrantRoleToCurrentUser(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	currentUser, err := c.context.client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = c.client().Grant(ctx, sdk.NewGrantRoleRequest(id, sdk.GrantRole{
		User: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentUser)),
	}))
	require.NoError(t, err)
}

// TODO: move later to grants client
func (c *RoleClient) GrantOwnershipOnAccountObject(t *testing.T, roleId sdk.AccountObjectIdentifier, objectId sdk.AccountObjectIdentifier, objectType sdk.ObjectType) {
	t.Helper()
	ctx := context.Background()

	err := c.context.client.Grants.GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectId,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: sdk.Pointer(roleId),
		},
		new(sdk.GrantOwnershipOptions),
	)
	require.NoError(t, err)
}
