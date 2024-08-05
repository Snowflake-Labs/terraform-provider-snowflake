package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type GrantClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewGrantClient(context *TestClientContext, idsGenerator *IdsGenerator) *GrantClient {
	return &GrantClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *GrantClient) client() sdk.Grants {
	return c.context.client.Grants
}

func (c *GrantClient) GrantOnSchemaToAccountRole(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, accountRoleId sdk.AccountObjectIdentifier, privileges ...sdk.SchemaPrivilege) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: &schemaId,
			},
		},
		accountRoleId,
		new(sdk.GrantPrivilegesToAccountRoleOptions),
	)
	require.NoError(t, err)
}
