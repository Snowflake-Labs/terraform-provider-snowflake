package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamOnTableIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name.schema_name.target_table_name"
	streamOnTable, err := streamOnTableIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", streamOnTable.DatabaseName)
	r.Equal("schema_name", streamOnTable.SchemaName)
	r.Equal("target_table_name", streamOnTable.OnTableName)

	// Bad ID -- not enough fields
	id = "database.schema"
	_, err = streamOnTableIDFromString(id)
	r.Equal(fmt.Errorf("invalid format for on_table: database.schema , expected: <database_name.schema_name.target_table_name>"), err)

	// Bad ID
	id = ".."
	_, err = streamOnTableIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = streamOnTableIDFromString(id)
	r.Equal(fmt.Errorf("expecting 1 line"), err)
}
