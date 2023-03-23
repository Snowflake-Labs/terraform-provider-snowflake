package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/stretchr/testify/require"
)

func TestParseFunctionGrantID(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION||privilege_name|false|role1,role2|share1,share2")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(0, len(grantID.ArgumentDataTypes))
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(2, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseFunctionGrantIDWithArgs(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION|string,string|privilege_name|false|role1,role2|share1,share2")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(2, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseFunctionEmojiGrantIDWithArgs(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE❄️MY_SCHEMA❄️MY_FUNCTION❄️string,string❄️privilege_name❄️true❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(2, len(grantID.Shares))
	r.Equal(true, grantID.WithGrantOption)
}

func TestParseFunctionGrantOldID(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION()|privilege_name|role1,role2|true")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(0, len(grantID.ArgumentDataTypes))
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(true, grantID.WithGrantOption)
}

func TestParseFunctionGrantOldIDWithArgs(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION(string, string)|privilege_name|role1,role2|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseFunctionGrantOldIDWithArgsAndNames(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION(A string, B string)|privilege_name|role1,role2|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseFunctionGrantReallyOldIDWithArgsAndNames(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_FUNCTION(A string, B string)|privilege_name|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_FUNCTION", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(0, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}
