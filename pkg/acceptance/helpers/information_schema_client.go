package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

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

type QueryHistory struct {
	QueryId   string
	QueryText string
	QueryTag  string
}

func (c *InformationSchemaClient) GetQueryHistory(t *testing.T, limit int) []QueryHistory {
	t.Helper()
	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT * FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => %d))", limit))
	require.NoError(t, err)
	return collections.Map(result, func(m map[string]*any) QueryHistory {
		require.NotNil(t, m["QUERY_ID"])
		require.NotNil(t, m["QUERY_TEXT"])
		require.NotNil(t, m["QUERY_TAG"])
		return QueryHistory{
			QueryId:   (*m["QUERY_ID"]).(string),
			QueryText: (*m["QUERY_TEXT"]).(string),
			QueryTag:  (*m["QUERY_TAG"]).(string),
		}
	})
}

func (c *InformationSchemaClient) GetQueryHistoryByQueryId(t *testing.T, limit int, queryId string) QueryHistory {
	t.Helper()
	result, err := c.client().QueryUnsafe(context.Background(), fmt.Sprintf("SELECT QUERY_TEXT FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY(RESULT_LIMIT => %d)) WHERE QUERY_ID = '%s'", limit, queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotNil(t, result[0]["QUERY_TEXT"])
	return QueryHistory{
		QueryId:   (*result[0]["QUERY_ID"]).(string),
		QueryText: (*result[0]["QUERY_TEXT"]).(string),
		QueryTag:  (*result[0]["QUERY_TAG"]).(string),
	}
}
