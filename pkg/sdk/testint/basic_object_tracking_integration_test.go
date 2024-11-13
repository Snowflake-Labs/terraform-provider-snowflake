package testint

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

// https://docs.snowflake.com/en/sql-reference/parameters#query-tag
func TestInt_ContextQueryTags(t *testing.T) {
	client := testClient(t)
	queryId := make(chan string, 1)

	ctx := context.Background()
	ctx = gosnowflake.WithQueryTag(ctx, "TERRAFORM_PROVIDE_TEST_QUERY_TAG")
	ctx = gosnowflake.WithQueryIDChan(ctx, queryId)
	_, err := client.QueryUnsafe(ctx, "SELECT 1")
	require.NoError(t, err)

	log.Println("Query id: " + <-queryId)
}

// https://select.dev/posts/snowflake-query-tags#using-query-comments-instead-of-query-tags
func TestInt_QueryComment(t *testing.T) {
	client := testClient(t)
	queryId := make(chan string, 1)

	ctx := context.Background()
	ctx = gosnowflake.WithQueryIDChan(ctx, queryId)
	_, err := client.QueryUnsafe(ctx, `SELECT 1; --{"comment": "some comment"}`)
	require.NoError(t, err)

	log.Println("Query id: " + <-queryId)
}

func TestInt_AppName(t *testing.T) {
	version := "v0.99.0"
	config := &gosnowflake.Config{
		// TODO: We filter by appname + it's somehow transferred to Terraform
		Application: fmt.Sprintf("terraform-provider-snowflake:%s", version),
	}
	client, _ := sdk.NewClient(config)
	_ = client
}

// TODO: trying to use connection parameters
