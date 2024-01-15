package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"strings"
)

// we need to keep track of literally everything to construct a unique identifier that can be imported
type GrantPrivilegesToAccountRoleID struct {
	RoleName         string
	Privileges       []string
	AllPrivileges    bool
	WithGrantOption  bool
	OnAccount        bool
	OnAccountObject  bool
	OnSchema         bool
	OnSchemaObject   bool
	All              bool
	Future           bool
	ObjectType       string
	ObjectName       string
	ObjectTypePlural string
	InSchema         bool
	SchemaName       string
	InDatabase       bool
	DatabaseName     string
}

func NewGrantPrivilegesToAccountRoleID(id string) GrantPrivilegesToAccountRoleID {
	parts := strings.Split(id, "|")
	privileges := strings.Split(parts[1], ",")
	if len(privileges) == 1 && privileges[0] == "" {
		privileges = []string{}
	}
	return GrantPrivilegesToAccountRoleID{
		RoleName:         parts[0],
		Privileges:       privileges,
		AllPrivileges:    parts[2] == "true",
		WithGrantOption:  parts[3] == "true",
		OnAccount:        parts[4] == "true",
		OnAccountObject:  parts[5] == "true",
		OnSchema:         parts[6] == "true",
		OnSchemaObject:   parts[7] == "true",
		All:              parts[8] == "true",
		Future:           parts[9] == "true",
		ObjectType:       parts[10],
		ObjectName:       parts[11],
		ObjectTypePlural: parts[12],
		InSchema:         parts[13] == "true",
		SchemaName:       parts[14],
		InDatabase:       parts[15] == "true",
		DatabaseName:     parts[16],
	}
}

func (v GrantPrivilegesToAccountRoleID) String() string {
	return helpers.EncodeSnowflakeID(v.RoleName, v.Privileges, v.AllPrivileges, v.WithGrantOption, v.OnAccount, v.OnAccountObject, v.OnSchema, v.OnSchemaObject, v.All, v.Future, v.ObjectType, v.ObjectName, v.ObjectTypePlural, v.InSchema, v.SchemaName, v.InDatabase, v.DatabaseName)
}
