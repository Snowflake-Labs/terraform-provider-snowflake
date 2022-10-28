package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExternalTestTableIDFromString(t *testing.T) {
	t.Helper()
	r := require.New(t)
	// Vanilla
	id := "database_name|schema_name|table"
	table, err := externalTableIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", table.DatabaseName)
	r.Equal("schema_name", table.SchemaName)
	r.Equal("table", table.ExternalTableName)

	// Bad ID -- not enough fields
	id = "database"
	_, err = streamOnObjectIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// Bad ID
	id = "||"
	_, err = streamOnObjectIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = streamOnObjectIDFromString(id)
	r.Equal(fmt.Errorf("1 line at a time"), err)

	// 2 lines
	id = `database_name|schema_name|table
	database_name|schema_name|table`
	_, err = streamOnObjectIDFromString(id)
	r.Equal(fmt.Errorf("1 line at a time"), err)
}

func ExternalTestTableStruct(t *testing.T) {
	t.Helper()
	r := require.New(t)

	// Vanilla
	table := &externalTableID{
		DatabaseName:      "database_name",
		SchemaName:        "schema_name",
		ExternalTableName: "table",
	}
	sID, err := table.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|table", sID)

	// Empty grant
	table = &externalTableID{}
	sID, err = table.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	table = &externalTableID{
		DatabaseName:      "database|name",
		ExternalTableName: "table|name",
	}
	sID, err = table.String()
	r.NoError(err)
	newTable, err := streamOnObjectIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newTable.DatabaseName)
	r.Equal("table|name", newTable.Name)
}
