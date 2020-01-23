package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema_name|pipe"
	pipe, err := pipeIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", pipe.DatabaseName)
	r.Equal("schema_name", pipe.SchemaName)
	r.Equal("pipe", pipe.PipeName)

	// Bad ID -- not enough fields
	id = "database"
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("3 fields allowed"), err)

	// Bad ID
	id = "||"
	_, err = pipeIDFromString(id)
	r.NoError(err)

	// 0 lines
	id = ""
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("1 line per pipe"), err)

	// 2 lines
	id = `database_name|schema_name|pipe
	database_name|schema_name|pipe`
	_, err = pipeIDFromString(id)
	r.Equal(fmt.Errorf("1 line per pipe"), err)
}

func TestPipeStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	pipe := &pipeID{
		DatabaseName: "database_name",
		SchemaName:   "schema_name",
		PipeName:     "pipe",
	}
	sID, err := pipe.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|pipe", sID)

	// Empty grant
	pipe = &pipeID{}
	sID, err = pipe.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	pipe = &pipeID{
		DatabaseName: "database|name",
		PipeName:     "pipe|name",
	}
	sID, err = pipe.String()
	r.NoError(err)
	newPipe, err := pipeIDFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newPipe.DatabaseName)
	r.Equal("pipe|name", newPipe.PipeName)
}
