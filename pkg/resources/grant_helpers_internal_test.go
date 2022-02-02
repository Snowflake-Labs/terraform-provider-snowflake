package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrantIDFromString(t *testing.T) {
	r := require.New(t)
	// Vanilla without GrantOption
	id := "database_name|schema|view_name|privilege|test1,test2"
	grant, err := grantIDFromString(id)
	r.NoError(err)

	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("privilege", grant.Privilege)
	r.Equal(false, grant.GrantOption)

	// Vanilla with GrantOption
	id = "database_name|schema|view_name|privilege|test1,test2|true"
	grant, err = grantIDFromString(id)
	r.NoError(err)

	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("privilege", grant.Privilege)
	r.Equal(true, grant.GrantOption)

	// No view
	id = "database_name|||privilege|"
	grant, err = grantIDFromString(id)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ObjectName)
	r.Equal("privilege", grant.Privilege)
	r.Equal(false, grant.GrantOption)

	// Bad ID -- not enough fields
	id = "database|name-privilege"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("5 or 6 fields allowed"), err)

	// Bad ID -- privilege in wrong area
	id = "database||name-privilege"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("5 or 6 fields allowed"), err)

	// too many fields
	id = "database_name|schema|view_name|privilege|false|2|too-many"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("5 or 6 fields allowed"), err)

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
		Roles:        []string{"test1", "test2"},
		GrantOption:  true,
	}
	gID, err := grant.String()
	r.NoError(err)
	r.Equal("database_name|schema|view_name|priv|test1,test2|true", gID)

	// Empty grant
	grant = &grantID{}
	gID, err = grant.String()
	r.NoError(err)
	r.Equal("|||||false", gID)

	// Grant with extra delimiters
	grant = &grantID{
		ResourceName: "database|name",
		SchemaName:   "schema|name",
		ObjectName:   "view|name",
		Privilege:    "priv",
		Roles:        []string{"test3", "test4"},
		GrantOption:  false,
	}
	gID, err = grant.String()
	r.NoError(err)
	newGrant, err := grantIDFromString(gID)
	r.NoError(err)
	r.Equal("database|name", newGrant.ResourceName)
	r.Equal("schema|name", newGrant.SchemaName)
	r.Equal("view|name", newGrant.ObjectName)
	r.Equal("priv", newGrant.Privilege)
	r.Equal([]string{"test3", "test4"}, newGrant.Roles)
	r.Equal(false, newGrant.GrantOption)
}
