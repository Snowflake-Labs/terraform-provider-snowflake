package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SessionPolicyClient struct {
	context *TestClientContext
}

func NewSessionPolicyClient(context *TestClientContext) *SessionPolicyClient {
	return &SessionPolicyClient{
		context: context,
	}
}

func (c *SessionPolicyClient) client() sdk.SessionPolicies {
	return c.context.client.SessionPolicies
}

func (c *SessionPolicyClient) CreateSessionPolicy(t *testing.T) (*sdk.SessionPolicy, func()) {
	t.Helper()
	id := c.context.newSchemaObjectIdentifier(random.StringN(12))
	return c.CreateSessionPolicyWithOptions(t, id, sdk.NewCreateSessionPolicyRequest(id))
}

func (c *SessionPolicyClient) CreateSessionPolicyWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateSessionPolicyRequest) (*sdk.SessionPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	sessionPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return sessionPolicy, c.DropSessionPolicyFunc(t, id)
}

func (c *SessionPolicyClient) DropSessionPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSessionPolicyRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
