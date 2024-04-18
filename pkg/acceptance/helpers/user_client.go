package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type UserClient struct {
	context *TestClientContext
}

func NewUserClient(context *TestClientContext) *UserClient {
	return &UserClient{
		context: context,
	}
}

func (d *UserClient) client() sdk.Users {
	return d.context.client.Users
}
