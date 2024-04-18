package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type RoleClient struct {
	context *TestClientContext
}

func NewRoleClient(context *TestClientContext) *RoleClient {
	return &RoleClient{
		context: context,
	}
}

func (d *RoleClient) client() sdk.Roles {
	return d.context.client.Roles
}
