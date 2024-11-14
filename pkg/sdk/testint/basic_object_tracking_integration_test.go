package testint

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// https://docs.snowflake.com/en/sql-reference/parameters#query-tag
func TestInt_ContextQueryTags(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	sessionId, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)

	queryTag := "some query tag"
	require.NoError(t, client.Parameters.SetSessionParameterOnAccount(ctx, sdk.SessionParameterQueryTag, queryTag))
	t.Cleanup(func() {
		_, err = client.QueryUnsafe(ctx, "ALTER SESSION UNSET QUERY_TAG")
		require.NoError(t, err)
	})

	queryId := executeQueryAndReturnQueryId(t, context.Background(), client)

	result, err := client.QueryUnsafe(ctx, fmt.Sprintf("SELECT QUERY_ID, QUERY_TAG FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY_BY_SESSION(SESSION_ID => %s, RESULT_LIMIT => 2)) WHERE QUERY_ID = '%s'", sessionId, queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, queryId, *result[0]["QUERY_ID"])
	require.Equal(t, queryTag, *result[0]["QUERY_TAG"])

	newQueryTag := "some other query tag"
	ctxWithQueryTag := gosnowflake.WithQueryTag(context.Background(), newQueryTag)
	newQueryId := executeQueryAndReturnQueryId(t, ctxWithQueryTag, client)

	result, err = client.QueryUnsafe(ctx, fmt.Sprintf("SELECT QUERY_ID, QUERY_TAG FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY_BY_SESSION(SESSION_ID => %s, RESULT_LIMIT => 2)) WHERE QUERY_ID = '%s'", sessionId, newQueryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, newQueryId, *result[0]["QUERY_ID"])
	require.Equal(t, newQueryTag, *result[0]["QUERY_TAG"])
}

func executeQueryAndReturnQueryId(t *testing.T, ctx context.Context, client *sdk.Client) string {
	queryIdChan := make(chan string, 1)
	ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)

	_, err := client.QueryUnsafe(ctx, "SELECT 1")
	require.NoError(t, err)

	return <-queryIdChan
}

// https://select.dev/posts/snowflake-query-tags#using-query-comments-instead-of-query-tags
func TestInt_QueryComment(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	sessionId, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)

	queryIdChan := make(chan string, 1)
	metadata := `{"comment": "some comment"}`
	_, err = client.QueryUnsafe(gosnowflake.WithQueryIDChan(ctx, queryIdChan), fmt.Sprintf(`SELECT 1; --%s`, metadata)) // TODO: Check in Snowhouse (check with Sebastian if this approach will be possible to query in snowhouse)
	require.NoError(t, err)
	queryId := <-queryIdChan

	result, err := client.QueryUnsafe(ctx, fmt.Sprintf("SELECT QUERY_ID, QUERY_TEXT FROM TABLE(INFORMATION_SCHEMA.QUERY_HISTORY_BY_SESSION(SESSION_ID => %s, RESULT_LIMIT => 2)) WHERE QUERY_ID = '%s'", sessionId, queryId))
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, queryId, *result[0]["QUERY_ID"])
	require.Equal(t, metadata, strings.Split((*result[0]["QUERY_TEXT"]).(string), "--")[1])
}

func TestInt_QueryComment_ClientSupport(t *testing.T) {

}

// - check comment characters (if they're valid characters etc.) - could be solved by receiving object and turning into e.g. JSON will be done at a lower level

func TestInt_AppName(t *testing.T) {
	// https://community.snowflake.com/s/article/How-to-see-application-name-added-in-the-connection-string-in-Snowsight
	t.Skip("there no way to check client application name by querying Snowflake's")

	version := "v0.99.0"
	config := sdk.DefaultConfig()
	config.Application = fmt.Sprintf("terraform-provider-snowflake:%s", version)
	client, err := sdk.NewClient(config)
	require.NoError(t, err)

	_, err = client.QueryUnsafe(context.Background(), "SELECT 1")
	require.NoError(t, err)
}

// TODO: trying to use connection parameters
