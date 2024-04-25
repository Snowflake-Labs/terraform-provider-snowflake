package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PasswordPolicyClient struct {
	context *TestClientContext
}

func NewPasswordPolicyClient(context *TestClientContext) *PasswordPolicyClient {
	return &PasswordPolicyClient{
		context: context,
	}
}

func (c *PasswordPolicyClient) client() sdk.PasswordPolicies {
	return c.context.client.PasswordPolicies
}

func (c *PasswordPolicyClient) CreatePasswordPolicy(t *testing.T) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return c.CreatePasswordPolicyInSchema(t, c.context.schemaId())
}

func (c *PasswordPolicyClient) CreatePasswordPolicyInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return c.CreatePasswordPolicyInSchemaWithOptions(t, schemaId, nil)
}

func (c *PasswordPolicyClient) CreatePasswordPolicyWithOptions(t *testing.T, options *sdk.CreatePasswordPolicyOptions) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return c.CreatePasswordPolicyInSchemaWithOptions(t, c.context.schemaId(), options)
}

func (c *PasswordPolicyClient) CreatePasswordPolicyInSchemaWithOptions(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, options *sdk.CreatePasswordPolicyOptions) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	name := random.AlphanumericN(12)
	id := sdk.NewSchemaObjectIdentifier(schemaId.DatabaseName(), schemaId.Name(), name)

	err := c.client().Create(ctx, id, options)
	require.NoError(t, err)

	passwordPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return passwordPolicy, c.DropPasswordPolicyFunc(t, id)
}

func (c *PasswordPolicyClient) DropPasswordPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropPasswordPolicyOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
