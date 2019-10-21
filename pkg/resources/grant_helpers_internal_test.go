package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrantIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema|view_name|privilege"
	grant, err := grantIDFromString(id)
	r.NoError(err)

	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ViewName)
	r.Equal("privilege", grant.Privilege)

	// No view
	id = "database_name|||privilege"
	grant, err = grantIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ViewName)
	r.Equal("privilege", grant.Privilege)

	// Bad ID -- not enough fields
	id = "database|name-privilege"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("4 fields allowed"), err)

	// Bad ID
	id = "database||||name-privilege"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("4 fields allowed"), err)

	// 0 lines
	id = ""
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("1 line per grant"), err)

	// 2 lines
	id = `database_name|schema|view_name|privilege
	database_name|schema|view_name|privilege`
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("1 line per grant"), err)
}

func TestGrantStruct(t *testing.T) {
	r := require.New(t)
	grant := &grantID{
		ResourceName: "database_name",
		SchemaName:   "schema",
		ViewName:     "view_name",
		Privilege:    "priv",
	}
	gID, err := grant.String()
	r.NoError(err)
	r.Equal("database_name|schema|view_name|priv", gID)

	grant = &grantID{}
	gID, err = grant.String()
	r.NoError(err)
	r.Equal("|||", gID)

}
