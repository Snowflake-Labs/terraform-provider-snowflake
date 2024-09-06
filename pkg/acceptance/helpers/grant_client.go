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

func (c *GrantClient) RevokePrivilegesOnSchemaObjectFromAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	schemaObjectIdentifier sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: objectType,
					Name:       schemaObjectIdentifier,
				},
			},
		},
		accountRoleId,
		new(sdk.RevokePrivilegesFromAccountRoleOptions),
	)

	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnSchemaObjectToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	schemaObjectIdentifier sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: objectType,
					Name:       schemaObjectIdentifier,
				},
			},
		},
		accountRoleId,
		&sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	require.NoError(t, err)
}

func (c *GrantClient) RevokePrivilegesOnDatabaseFromDatabaseRole(
	t *testing.T,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleId,
		new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnDatabaseToDatabaseRole(
	t *testing.T,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleId,
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	require.NoError(t, err)
}
