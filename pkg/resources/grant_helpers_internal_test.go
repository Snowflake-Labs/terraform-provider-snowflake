package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrantIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla
	id := "database_name|schema|view_name|privilege|on_future|on_all"
	grant, err := grantIDFromString(id)
	r.NoError(err)

	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("privilege", grant.Privilege)
	r.Equal("on_future", grant.OnFuture)
	r.Equal("on_all", grant.OnAll)

	// No view
	id = "database_name|||privilege|on_future|on_all"
	grant, err = grantIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ObjectName)
	r.Equal("privilege", grant.Privilege)

	// Bad ID -- not enough fields
	id = "database|name-privilege"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("6 fields allowed"), err)

	// Bad ID -- on_all in wrong area
	id = "database||||||on_all"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("6 fields allowed"), err)

	// too many fields
	id = "database_name|schema|view_name|privilege|on_future|on_all|extra"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("6 fields allowed"), err)

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

	// Vanilla
	grant := &grantID{
		ResourceName: "database_name",
		SchemaName:   "schema",
		ObjectName:   "view_name",
		Privilege:    "priv",
		OnFuture:     "on_future",
		OnAll:        "on_all",
	}
	gID, err := grant.String()
	r.NoError(err)
	r.Equal("database_name|schema|view_name|priv|on_future|on_all", gID)

	// Empty grant
	grant = &grantID{}
	gID, err = grant.String()
	r.NoError(err)
	r.Equal("|||||", gID)

	// Grant with extra delimiters
	grant = &grantID{
		ResourceName: "database|name",
		SchemaName:   "schema|name",
		ObjectName:   "view|name",
		Privilege:    "priv",
		OnFuture:     "on|future",
		OnAll:        "on|all",
	}
	gID, err = grant.String()
	r.NoError(err)
	newGrant, err := grantIDFromString(gID)
	r.NoError(err)
	r.Equal("database|name", newGrant.ResourceName)
	r.Equal("schema|name", newGrant.SchemaName)
	r.Equal("view|name", newGrant.ObjectName)
	r.Equal("priv", newGrant.Privilege)
	r.Equal("on|future", newGrant.OnFuture)
	r.Equal("on|all", newGrant.OnAll)
}
