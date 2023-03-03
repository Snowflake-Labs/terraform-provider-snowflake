package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/stretchr/testify/require"
)

func TestParseProcedureGrantID(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseProcedureGrantID("MY_DATABASE❄️MY_SCHEMA❄️MY_PROCEDURE❄️❄️privilege_name❄️false❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_PROCEDURE", grantID.ObjectName)
	r.Equal(0, len(grantID.ArgumentDataTypes))
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(2, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseProcedureGrantIDWithArgs(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseProcedureGrantID("MY_DATABASE❄️MY_SCHEMA❄️MY_PROCEDURE❄️string,string❄️privilege_name❄️false❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_PROCEDURE", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(2, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseProcedureGrantOldID(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseProcedureGrantID("MY_DATABASE|MY_SCHEMA|MY_PROCEDURE()|privilege_name|role1,role2|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_PROCEDURE", grantID.ObjectName)
	r.Equal(0, len(grantID.ArgumentDataTypes))
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseProcedureGrantOldIDWithArgs(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseProcedureGrantID("MY_DATABASE|MY_SCHEMA|MY_PROCEDURE( string , string)|privilege_name|role1,role2|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_PROCEDURE", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}

func TestParseProcedureGrantOldIDWithArgsAndNames(t *testing.T) {
	r := require.New(t)
	grantID, err := resources.ParseFunctionGrantID("MY_DATABASE|MY_SCHEMA|MY_PROCEDURE(A string, B string)|privilege_name|role1,role2|false")
	r.NoError(err)
	r.Equal("MY_DATABASE", grantID.DatabaseName)
	r.Equal("MY_SCHEMA", grantID.SchemaName)
	r.Equal("MY_PROCEDURE", grantID.ObjectName)
	r.Equal(2, len(grantID.ArgumentDataTypes))
	r.Equal("string", grantID.ArgumentDataTypes[0])
	r.Equal("string", grantID.ArgumentDataTypes[1])
	r.Equal("privilege_name", grantID.Privilege)
	r.Equal(2, len(grantID.Roles))
	r.Equal(0, len(grantID.Shares))
	r.Equal(false, grantID.WithGrantOption)
}
