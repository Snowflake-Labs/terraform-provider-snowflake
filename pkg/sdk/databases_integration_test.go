package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabasesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	databaseTest2, databaseCleanup2 := createDatabase(t, client)

	t.Run("without show options", func(t *testing.T) {
		databases, err := client.Databases.Show(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(databases), 2)
		databaseIDs := make([]AccountObjectIdentifier, len(databases))
		for i, database := range databases {
			databaseIDs[i] = database.ID()
		}
		assert.Contains(t, databaseIDs, databaseTest.ID())
		assert.Contains(t, databaseIDs, databaseTest2.ID())
	})

	t.Run("with terse", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Terse: Bool(true),
			Like: &Like{
				Pattern: String(databaseTest.Name),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))
		database := databases[0]
		assert.Equal(t, databaseTest.Name, database.Name)
		assert.NotEmpty(t, database.CreatedOn)
		assert.Empty(t, database.DroppedOn)
		assert.Empty(t, database.Owner)
	})

	t.Run("with history", func(t *testing.T) {
		// need to drop a database to test if the "dropped_on" column is populated
		databaseCleanup2()
		showOptions := &ShowDatabasesOptions{
			History: Bool(true),
			Like: &Like{
				Pattern: String(databaseTest2.Name),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))
		database := databases[0]
		assert.Equal(t, databaseTest2.Name, database.Name)
		assert.NotEmpty(t, database.DroppedOn)
	})

	t.Run("with like starts with", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			StartsWith: String(databaseTest.Name),
			LimitFrom: &LimitFrom{
				Rows: Int(1),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(databases))
		database := databases[0]
		assert.Equal(t, databaseTest.Name, database.Name)
	})

	t.Run("when searching a non-existent database", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		databases, err := client.Databases.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(databases))
	})

	/*
		// there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
		t.Run("when limiting the number of results", func(t *testing.T) {
			showOptions := &MaskingPolicyShowOptions{
				In: &In{
					Schema: schemaTest.ID(),
				},
				Limit: Int(1),
			}
			maskingPolicies, err := client.MaskingPolicies.Show(ctx, showOptions)
			require.NoError(t, err)
			assert.Equal(t, 1, len(maskingPolicies))
		})
	*/
}
