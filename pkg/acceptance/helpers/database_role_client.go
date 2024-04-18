package helpers

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type DatabaseRoleClient struct {
	context *TestClientContext
}

func NewDatabaseRoleClient(context *TestClientContext) *DatabaseRoleClient {
	return &DatabaseRoleClient{
		context: context,
	}
}

func (d *DatabaseRoleClient) client() sdk.DatabaseRoles {
	return d.context.client.DatabaseRoles
}
