package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id, err := idFromString("database_name|schema_name|pipe")
	r.NoError(err)
	r.Equal("database_name", id.Database)
	r.Equal("schema_name", id.Schema)
	r.Equal("pipe", id.Name)

	_, err = idFromString("database")
	r.Equal(fmt.Errorf("wrong number of fields in record"), err)

	// Bad ID
	_, err = idFromString("||")
	r.NoError(err)

	// 0 lines
	_, err = idFromString("")
	r.Equal(fmt.Errorf("EOF"), err)
}

func TestIDStruct(t *testing.T) {
	r := require.New(t)

	// Vanilla
	id := &schemaScopedID{
		Database: "database_name",
		Schema:   "schema_name",
		Name:     "pipe",
	}
	sID, err := id.String()
	r.NoError(err)
	r.Equal("database_name|schema_name|pipe", sID)

	// Empty grant
	id = &schemaScopedID{}
	sID, err = id.String()
	r.NoError(err)
	r.Equal("||", sID)

	// Grant with extra delimiters
	id = &schemaScopedID{
		Database: "database|name",
		Name:     "pipe|name",
	}
	sID, err = id.String()
	r.NoError(err)
	newID, err := idFromString(sID)
	r.NoError(err)
	r.Equal("database|name", newID.Database)
	r.Equal("pipe|name", newID.Name)
}
