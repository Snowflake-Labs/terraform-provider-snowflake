package snowflake

import (
	"fmt"
)

type futureGrantType string

const (
	futureTableType futureGrantType = "TABLE"
	futureViewType  futureGrantType = "VIEW"
)

// FutureGrantBuilder abstracts the creation of FutureGrantExecutables
type FutureGrantBuilder struct {
	name            string
	qualifiedName   string
	futureGrantType futureGrantType
}

// Name returns the object name for this FutureGrantBuilder
func (fgb *FutureGrantBuilder) Name() string {
	return fgb.name
}

// FutureTableGrant returns a pointer to a FutureGrantBuilder for a table
func FutureTableGrant(db, schema string) GrantBuilder {
	return &FutureGrantBuilder{
		name:            schema,
		qualifiedName:   fmt.Sprintf(`"%v"."%v"`, db, schema),
		futureGrantType: futureTableType,
	}
}

// FutureViewGrant returns a pointer to a FutureGrantBuilder for a view
func FutureViewGrant(db, schema string) GrantBuilder {
	return &FutureGrantBuilder{
		name:            schema,
		qualifiedName:   fmt.Sprintf(`"%v"."%v"`, db, schema),
		futureGrantType: futureViewType,
	}
}

// Show returns the SQL that will show all privileges on the grant
func (fgb *FutureGrantBuilder) Show() string {
	return fmt.Sprintf(`SHOW FUTURE GRANTS IN SCHEMA %v`, fgb.qualifiedName)
}

// FutureGrantExecutable abstracts the creation of SQL queries to build future grants for
// different future grant types.
type FutureGrantExecutable struct {
	grantName       string
	futureGrantType futureGrantType
	granteeName     string
}

// Role returns a pointer to a FutureGrantExecutable for a role
func (fgb *FutureGrantBuilder) Role(n string) GrantExecutable {
	return &FutureGrantExecutable{
		grantName:       fgb.qualifiedName,
		futureGrantType: fgb.futureGrantType,
		granteeName:     n,
	}
}

// Share is not implemented because future objects cannot be granted to shares.
func (gb *FutureGrantBuilder) Share(n string) GrantExecutable {
	return nil
}

// Grant returns the SQL that will grant future privileges on the grant to the grantee
func (fge *FutureGrantExecutable) Grant(p string) string {
	return fmt.Sprintf(`GRANT %v ON FUTURE %vS IN SCHEMA %v TO ROLE "%v"`,
		p, fge.futureGrantType, fge.grantName, fge.granteeName)
}

// Revoke returns the SQL that will revoke future privileges on the grant from the grantee
func (fge *FutureGrantExecutable) Revoke(p string) string {
	return fmt.Sprintf(`REVOKE %v ON FUTURE %vS IN SCHEMA %v FROM ROLE "%v"`,
		p, fge.futureGrantType, fge.grantName, fge.granteeName)
}

// Show returns the SQL that will show all future grants on the schema
func (fge *FutureGrantExecutable) Show() string {
	return fmt.Sprintf(`SHOW FUTURE GRANTS IN SCHEMA %v`, fge.grantName)
}
