package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MaskingPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewMaskingPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *MaskingPolicyClient {
	return &MaskingPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *MaskingPolicyClient) client() sdk.MaskingPolicies {
	return c.context.client.MaskingPolicies
}

func (c *MaskingPolicyClient) CreateMaskingPolicy(t *testing.T) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	return c.CreateMaskingPolicyInSchema(t, c.ids.SchemaId())
}

func (c *MaskingPolicyClient) CreateMaskingPolicyInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	signature := []sdk.TableColumnSignature{
		{
			Name: c.ids.Alpha(),
			Type: sdk.DataTypeVARCHAR,
		},
		{
			Name: c.ids.Alpha(),
			Type: sdk.DataTypeVARCHAR,
		},
	}
	expression := "REPLACE('X', 1, 2)"
	return c.CreateMaskingPolicyWithOptions(t, schemaId, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{})
}

func (c *MaskingPolicyClient) CreateMaskingPolicyIdentity(t *testing.T, columnType sdk.DataType) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	name := "a"
	signature := []sdk.TableColumnSignature{
		{
			Name: name,
			Type: columnType,
		},
	}
	expression := "a"
	return c.CreateMaskingPolicyWithOptions(t, c.ids.SchemaId(), signature, columnType, expression, &sdk.CreateMaskingPolicyOptions{})
}

func (c *MaskingPolicyClient) CreateMaskingPolicyWithOptions(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, signature []sdk.TableColumnSignature, returns sdk.DataType, expression string, options *sdk.CreateMaskingPolicyOptions) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	err := c.client().Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	maskingPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return maskingPolicy, c.DropMaskingPolicyFunc(t, id)
}

func (c *MaskingPolicyClient) DropMaskingPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		assert.NoError(t, err)
	}
}
