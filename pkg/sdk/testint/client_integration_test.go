package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func TestInt_Client_AdditionalMetadata(t *testing.T) {
	client := testClient(t)
	metadata := tracking.NewMetadata("v1.13.1002-rc-test", resources.Database, tracking.CreateOperation)

	assertQueryMetadata := func(t *testing.T, queryId string) {
		t.Helper()
		queryText := testClientHelper().InformationSchema.GetQueryHistoryByQueryId(t, 20, queryId).QueryText
		parsedMetadata, err := tracking.ParseMetadata(queryText)
		require.NoError(t, err)
		require.Equal(t, metadata, parsedMetadata)
	}

	t.Run("query one", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		row := struct {
			One int `db:"ONE"`
		}{}
		err := client.QueryOneForTests(ctx, &row, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("query", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		var rows []struct {
			One int `db:"ONE"`
		}
		err := client.QueryForTests(ctx, &rows, "SELECT 1 AS ONE")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})

	t.Run("exec", func(t *testing.T) {
		queryIdChan := make(chan string, 1)
		ctx := context.Background()
		ctx = tracking.NewContext(ctx, metadata)
		ctx = gosnowflake.WithQueryIDChan(ctx, queryIdChan)
		_, err := client.ExecForTests(ctx, "SELECT 1")
		require.NoError(t, err)

		assertQueryMetadata(t, <-queryIdChan)
	})
}
