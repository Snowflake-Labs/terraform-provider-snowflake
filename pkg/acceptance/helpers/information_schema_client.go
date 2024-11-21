package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type InformationSchemaClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewInformationSchemaClient(context *TestClientContext, idsGenerator *IdsGenerator) *InformationSchemaClient {
	return &InformationSchemaClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *InformationSchemaClient) client() *sdk.Client {
	return c.context.client
}

func (c *InformationSchemaClient) GetQueryTextByQueryId(t *testing.T, queryId string) string {
	t.Helper()
	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT QUERY_TEXT FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => 20)) WHERE QUERY_ID = '%s'", queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0]["QUERY_TEXT"])
	return (*result[0]["QUERY_TEXT"]).(string)
}

func (c *InformationSchemaClient) GetQueryTagByQueryId(t *testing.T, queryId string) string {
	t.Helper()
	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT QUERY_TAG FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => 20)) WHERE QUERY_ID = '%s'", queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0]["QUERY_TAG"])
	return (*result[0]["QUERY_TAG"]).(string)
}
