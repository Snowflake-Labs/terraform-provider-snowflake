package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SystemFunctionsClient struct {
	context *TestClientContext
}

func NewSystemFunctionsClient(context *TestClientContext) *SystemFunctionsClient {
	return &SystemFunctionsClient{
		context: context,
	}
}

func (c *SystemFunctionsClient) client() sdk.SystemFunctions {
	return c.context.client.SystemFunctions
}

func (c *SystemFunctionsClient) GetTag(t *testing.T, tagId sdk.SchemaObjectIdentifier, objectId sdk.ObjectIdentifier, objectType sdk.ObjectType) (*string, error) {
	t.Helper()
	ctx := context.Background()
	client := c.context.client.SystemFunctions

	return client.GetTag(ctx, tagId, objectId, objectType)
}
