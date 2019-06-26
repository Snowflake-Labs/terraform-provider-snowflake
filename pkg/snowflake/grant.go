package snowflake

import (
	"fmt"
)

type grantType string

const (
	databaseType grantType = "DATABASE"
	schemaType   grantType = "SCHEMA"
	viewType     grantType = "VIEW"
)

// GrantBuilder abstracts the creation of GrantExecutables
type GrantBuilder struct {
	name      string
	grantType grantType
}

// DatabaseGrant returns a pointer to a GrantBuilder for a database
func DatabaseGrant(name string) *GrantBuilder {
	return &GrantBuilder{
		name:      name,
		grantType: databaseType,
	}
}

// SchemaGrant returns a pointer to a GrantBuilder for a schema
func SchemaGrant(name string) *GrantBuilder {
	return &GrantBuilder{
		name:      name,
		grantType: schemaType,
	}
}

// ViewGrant returns a pointer to a GrantBuilder for a view
func ViewGrant(name string) *GrantBuilder {
	return &GrantBuilder{
		name:      name,
		grantType: viewType,
	}
}

// Show returns the SQL that will show all privileges on the grant
func (gb *GrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW GRANTS ON %v "%v"`, gb.grantType, gb.name)
}

type granteeType string

const (
	roleType  granteeType = "ROLE"
	shareType granteeType = "SHARE"
	userType  granteeType = "USER" // user is only supported for RoleGrants.
)

// GrantExecutable abstracts the creation of SQL queries to build grants for
// different resources
type GrantExecutable struct {
	grantName   string
	grantType   grantType
	granteeName string
	granteeType granteeType
}

// Role returns a pointer to a GrantExecutable for a role
func (gb *GrantBuilder) Role(n string) *GrantExecutable {
	return &GrantExecutable{
		grantName:   gb.name,
		grantType:   gb.grantType,
		granteeName: n,
		granteeType: roleType,
	}
}

// Share returns a pointer to a GrantExecutable for a share
func (gb *GrantBuilder) Share(n string) *GrantExecutable {
	return &GrantExecutable{
		grantName:   gb.name,
		grantType:   gb.grantType,
		granteeName: n,
		granteeType: shareType,
	}
}

// Grant returns the SQL that will grant privileges on the grant to the grantee
func (ge *GrantExecutable) Grant(p string) string {
	return fmt.Sprintf(`GRANT %v ON %v "%v" TO %v "%v"`,
		p, ge.grantType, ge.grantName, ge.granteeType, ge.granteeName)
}

// Revoke returns the SQL that will revoke privileges on the grant from the grantee
func (ge *GrantExecutable) Revoke(p string) string {
	return fmt.Sprintf(`REVOKE %v ON %v "%v" FROM %v "%v"`,
		p, ge.grantType, ge.grantName, ge.granteeType, ge.granteeName)
}
