package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type RowAccessPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewRowAccessPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *RowAccessPolicyClient {
	return &RowAccessPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *RowAccessPolicyClient) client() sdk.RowAccessPolicies {
	return c.context.client.RowAccessPolicies
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicy(t *testing.T) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", sdk.DataTypeNumber)
	body := "true"
	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)

	err := c.client().Create(ctx, createRequest)
	require.NoError(t, err)

	rowAccessPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return rowAccessPolicy, c.DropRowAccessPolicyFunc(t, id)
}

func (c *RowAccessPolicyClient) DropRowAccessPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}

// GetRowAccessPolicyFor is based on https://docs.snowflake.com/en/user-guide/security-row-intro#obtain-database-objects-with-a-row-access-policy.
// TODO: extract getting row access policies as resource (like getting tag in system functions)
func (c *RowAccessPolicyClient) GetOneRowAccessPolicyFor(t *testing.T, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) (*PolicyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := &PolicyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := c.context.client.QueryOneForTests(ctx, s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}
