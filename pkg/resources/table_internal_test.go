package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema_name|table"
	table, err := tableIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", table.DatabaseName)
	r.Equal("schema_name", table.SchemaName)
	r.Equal("table", table.TableName)

	// Bad ID -- not enough fields
	id = "database"
	_, err = tableIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// Bad ID
	id = "||"
	_, err = tableIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = tableIDFromString(id)
	r.Equal(fmt.Errorf("1 line at a time"), err)

	// 2 lines
	id = `database_name|schema_name|table
	database_name|schema_name|table`
	_, err = tableIDFromString(id)
	r.Equal(fmt.Errorf("1 line at a time"), err)
}

func TestTableStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	table := &tableID{
		DatabaseName: "database_name",
		SchemaName:   "schema_name",
		TableName:    "table",
	}
	sID, err := table.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|table", sID)

	// Empty grant
	table = &tableID{}
	sID, err = table.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	table = &tableID{
		DatabaseName: "database|name",
		TableName:    "table|name",
	}
	sID, err = table.String()
	r.NoError(err)
	newTable, err := tableIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newTable.DatabaseName)
	r.Equal("table|name", newTable.TableName)
}
