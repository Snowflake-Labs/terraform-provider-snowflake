package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type RowAccessPolicyClient struct {
	context *TestClientContext
}

func NewRowAccessPolicyClient(context *TestClientContext) *RowAccessPolicyClient {
	return &RowAccessPolicyClient{
		context: context,
	}
}

func (c *RowAccessPolicyClient) client() sdk.RowAccessPolicies {
	return c.context.client.RowAccessPolicies
}

func (c *RowAccessPolicyClient) CreateRowAccessPolicy(t *testing.T) (*sdk.RowAccessPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.context.newSchemaObjectIdentifier(random.AlphanumericN(12))
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
func (c *RowAccessPolicyClient) GetRowAccessPolicyFor(t *testing.T, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) (*PolicyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := &PolicyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := c.context.client.QueryOneForTests(ctx, s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}

type PolicyReference struct {
	PolicyDb          string         `db:"POLICY_DB"`
	PolicySchema      string         `db:"POLICY_SCHEMA"`
	PolicyName        string         `db:"POLICY_NAME"`
	PolicyKind        string         `db:"POLICY_KIND"`
	RefDatabaseName   string         `db:"REF_DATABASE_NAME"`
	RefSchemaName     string         `db:"REF_SCHEMA_NAME"`
	RefEntityName     string         `db:"REF_ENTITY_NAME"`
	RefEntityDomain   string         `db:"REF_ENTITY_DOMAIN"`
	RefColumnName     sql.NullString `db:"REF_COLUMN_NAME"`
	RefArgColumnNames string         `db:"REF_ARG_COLUMN_NAMES"`
	TagDatabase       sql.NullString `db:"TAG_DATABASE"`
	TagSchema         sql.NullString `db:"TAG_SCHEMA"`
	TagName           sql.NullString `db:"TAG_NAME"`
	PolicyStatus      string         `db:"POLICY_STATUS"`
}
