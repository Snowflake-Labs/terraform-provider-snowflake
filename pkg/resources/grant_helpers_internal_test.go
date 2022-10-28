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

	// too many fields
	id = "database_name|schema|view_name|privilege|false|2|too-many"
	_, err = grantIDFromString(id)
	r.Equal(fmt.Errorf("1 to 6 fields allowed in ID"), err)

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

func TestGrantLegacyID(t *testing.T) {
	// Testing that grants with legacy ID structure resolves to expected output
	r := require.New(t)
	gID := "database_name|schema|view_name|priv|true"
	grant, err := grantIDFromString(gID)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("priv", grant.Privilege)
	r.Equal([]string{}, grant.Roles)
	r.Equal(true, grant.GrantOption)

	gID = "database_name|schema|view_name|priv|false"
	grant, err = grantIDFromString(gID)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("priv", grant.Privilege)
	r.Equal([]string{}, grant.Roles)
	r.Equal(false, grant.GrantOption)

	gID = "database_name|schema|view_name|priv"
	grant, err = grantIDFromString(gID)
	r.NoError(err)
	r.Equal("database_name", grant.ResourceName)
	r.Equal("schema", grant.SchemaName)
	r.Equal("view_name", grant.ObjectName)
	r.Equal("priv", grant.Privilege)
	r.Equal([]string{}, grant.Roles)
	r.Equal(false, grant.GrantOption)

}

func TestGrantIDFromStringRoleGrant(t *testing.T) {
	r := require.New(t)
	gID := "role_a||||role1,role2|"
	grant, err := grantIDFromString(gID)
	r.NoError(err)
	r.Equal("role_a", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ObjectName)
	r.Equal("", grant.Privilege)
	r.Equal([]string{"role1", "role2"}, grant.Roles)
	r.Equal(false, grant.GrantOption)

	// Testing the legacy ID structure passes as expected
	gID = "role_a"
	grant, err = grantIDFromString(gID)
	r.NoError(err)
	r.Equal("role_a", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ObjectName)
	r.Equal("", grant.Privilege)
	r.Equal([]string{}, grant.Roles)
	r.Equal(false, grant.GrantOption)

	gID = "role_b||||role3,role4|false"
	grant, err = grantIDFromString(gID)
	r.NoError(err)
	r.Equal("role_b", grant.ResourceName)
	r.Equal("", grant.SchemaName)
	r.Equal("", grant.ObjectName)
	r.Equal("", grant.Privilege)
	r.Equal([]string{"role3", "role4"}, grant.Roles)
	r.Equal(false, grant.GrantOption)
}
