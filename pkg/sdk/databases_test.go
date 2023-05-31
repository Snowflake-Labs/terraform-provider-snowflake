package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabasesShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowDatabasesOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("terse", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Terse: Bool(true),
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW TERSE DATABASES`
		assert.Equal(t, expected, actual)
	})

	t.Run("history", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			History: Bool(true),
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW DATABASES HISTORY`
		assert.Equal(t, expected, actual)
	})

	t.Run("like", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Like: &Like{
				Pattern: String("db1"),
			},
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW DATABASES LIKE 'db1'`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete", func(t *testing.T) {
		showOptions := &ShowDatabasesOptions{
			Terse:   Bool(true),
			History: Bool(true),
			Like: &Like{
				Pattern: String("db2"),
			},
			LimitFrom: &LimitFrom{
				Rows: Int(1),
				From: String("db1"),
			},
		}
		actual, err := structToSQL(showOptions)
		require.NoError(t, err)
		expected := `SHOW TERSE DATABASES HISTORY LIKE 'db2' LIMIT 1 FROM 'db1'`
		assert.Equal(t, expected, actual)
	})
}
