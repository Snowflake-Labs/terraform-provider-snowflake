package testint

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt_Client_UnsafeQuery(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test show databases", func(t *testing.T) {
		sql := fmt.Sprintf("SHOW DATABASES LIKE '%%%s%%'", testDb(t).Name)
		results, err := client.QueryUnsafe(ctx, sql)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		row := results[0]
		assert.Equal(t, testDb(t).Name, *row["name"])
		assert.NotEmpty(t, *row["created_on"])
		assert.Equal(t, "STANDARD", *row["kind"])
		assert.Equal(t, "ACCOUNTADMIN", *row["owner"])
		assert.Equal(t, "", *row["options"])
		assert.Equal(t, "", *row["comment"])
		assert.Equal(t, "N", *row["is_default"])
		assert.Equal(t, "Y", *row["is_current"])
		assert.Nil(t, *row["budget"])
	})

	t.Run("test more results", func(t *testing.T) {
		db1, db1Cleanup := createDatabase(t, client)
		t.Cleanup(db1Cleanup)
		db2, db2Cleanup := createDatabase(t, client)
		t.Cleanup(db2Cleanup)
		db3, db3Cleanup := createDatabase(t, client)
		t.Cleanup(db3Cleanup)

		sql := "SHOW DATABASES"
		results, err := client.QueryUnsafe(ctx, sql)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(results), 4)
		names := make([]any, len(results))
		for i, r := range results {
			names[i] = *r["name"]
		}
		assert.Contains(t, names, testDb(t).Name)
		assert.Contains(t, names, db1.Name)
		assert.Contains(t, names, db2.Name)
		assert.Contains(t, names, db3.Name)
	})
}
