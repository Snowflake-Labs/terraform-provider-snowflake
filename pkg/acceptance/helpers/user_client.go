package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type UserClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewUserClient(context *TestClientContext, idsGenerator *IdsGenerator) *UserClient {
	return &UserClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *UserClient) client() sdk.Users {
	return c.context.client.Users
}

func (c *UserClient) CreateUser(t *testing.T) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{})
}

// TODO [SNOW-955520]: we have to control the name
func (c *UserClient) CreateUserWithName(t *testing.T, name string) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, sdk.NewAccountObjectIdentifier(name), &sdk.CreateUserOptions{})
}

func (c *UserClient) CreateUserWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateUserOptions) (*sdk.User, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)
	user, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return user, c.DropUserFunc(t, id)
}

func (c *UserClient) DropUserFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropUserOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
