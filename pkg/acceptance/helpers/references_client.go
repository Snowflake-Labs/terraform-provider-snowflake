package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type PolicyReferencesClient struct {
	context *TestClientContext
}

func NewPolicyReferencesClient(context *TestClientContext) *PolicyReferencesClient {
	return &PolicyReferencesClient{
		context: context,
	}
}

func (c *PolicyReferencesClient) client() sdk.RowAccessPolicies {
	return c.context.client.RowAccessPolicies
}

// GetPolicyReferences is based on https://docs.snowflake.com/en/sql-reference/functions/policy_references.
func (c *PolicyReferencesClient) GetPolicyReferences(t *testing.T, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) ([]PolicyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := []PolicyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := c.context.client.QueryForTests(ctx, &s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}

// GetPolicyReference is based on https://docs.snowflake.com/en/sql-reference/functions/policy_references.
func (c *PolicyReferencesClient) GetPolicyReference(t *testing.T, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) (*PolicyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := &PolicyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := c.context.client.QueryOneForTests(ctx, s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}
