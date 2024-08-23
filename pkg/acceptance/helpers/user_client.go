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

func (c *UserClient) CreateUserWithPrefix(t *testing.T, prefix string) (*sdk.User, func()) {
	t.Helper()
	return c.CreateUserWithOptions(t, c.ids.RandomAccountObjectIdentifierWithPrefix(prefix), &sdk.CreateUserOptions{})
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

func (c *UserClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.User, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *UserClient) Disable(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterUserOptions{Set: &sdk.UserSet{ObjectProperties: &sdk.UserObjectProperties{Disable: sdk.Bool(true)}}})
	require.NoError(t, err)
}
