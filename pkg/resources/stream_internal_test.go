package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema_name|stream"
	stream, err := streamIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", stream.Database)
	r.Equal("schema_name", stream.Schema)
	r.Equal("stream", stream.Name)

	// Bad ID -- not enough fields
	id = "database"
	_, err = streamIDFromString(id)
	r.Equal(fmt.Errorf("wrong number of fields in record"), err)

	// Bad ID
	id = "||"
	_, err = streamIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = streamIDFromString(id)
	r.Equal(fmt.Errorf("EOF"), err)
}

func TestStreamStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	stream := &streamID{
		Database: "database_name",
		Schema:   "schema_name",
		Name:     "stream_name",
	}
	sID, err := stream.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|stream_name", sID)

	// Empty grant
	stream = &streamID{}
	sID, err = stream.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	stream = &streamID{
		Database: "database|name",
		Name:     "stream|name",
	}
	sID, err = stream.String()
	r.NoError(err)
	newStream, err := streamIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newStream.Database)
	r.Equal("stream|name", newStream.Name)
}

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
