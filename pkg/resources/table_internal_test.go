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
	r.Equal("database_name", table.Database)
	r.Equal("schema_name", table.Schema)
	r.Equal("table", table.Name)

	// Bad ID -- not enough fields
	id = "database"
	_, err = tableIDFromString(id)
	r.Equal(fmt.Errorf("wrong number of fields in record"), err)

	// Bad ID
	id = "||"
	_, err = tableIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = tableIDFromString(id)
	r.Equal(fmt.Errorf("EOF"), err)
}

func TestTableStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	table := &tableID{
		Database: "database_name",
		Schema:   "schema_name",
		Name:     "table",
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
		Database: "database|name",
		Name:     "table|name",
	}
	sID, err = table.String()
	r.NoError(err)
	newTable, err := tableIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newTable.Database)
	r.Equal("table|name", newTable.Name)
}
