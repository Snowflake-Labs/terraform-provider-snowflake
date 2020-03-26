package snowflake

import (
	"fmt"
)

type onAllGrantType string

const (
	onAllTableType onAllGrantType = "TABLE"
	onAllViewType  onAllGrantType = "VIEW"
)

// AllGrantBuilder abstracts the creation of AllGrantBuilder
type OnAllGrantBuilder struct {
	name           string
	qualifiedName  string
	onAllGrantType onAllGrantType
}

// Name returns the object name for this AllGrantBuilder
func (agb *OnAllGrantBuilder) Name() string {
	return agb.name
}

// OnAllTableGrant returns a pointer to a AllGrantBuilder for a table
func OnAllTableGrant(db, schema string) GrantBuilder {
	return &OnAllGrantBuilder{
		name:           schema,
		qualifiedName:  fmt.Sprintf(`"%v"."%v"`, db, schema),
		onAllGrantType: onAllTableType,
	}
}

// AllViewGrant returns a pointer to a AllGrantBuilder for a view
func OnAllViewGrant(db, schema string) GrantBuilder {
	return &OnAllGrantBuilder{
		name:           schema,
		qualifiedName:  fmt.Sprintf(`"%v"."%v"`, db, schema),
		onAllGrantType: onAllViewType,
	}
}

// Show for OnAllGrantBuilder is not implemented because
// It is necessary to provide role to show privileges for it
func (agb *OnAllGrantBuilder) Show() string {
	return ""
}

// AllGrantExecutable abstracts the creation of SQL queries to build all grants for
// different all grant types.
type OnAllGrantExecutable struct {
	grantName      string
	onAllGrantType onAllGrantType
	granteeName    string
}

// Role returns a pointer to a AllGrantExecutable for a role
func (agb *OnAllGrantBuilder) Role(n string) GrantExecutable {
	return &OnAllGrantExecutable{
		grantName:      agb.qualifiedName,
		onAllGrantType: agb.onAllGrantType,
		granteeName:    n,
	}
}

// Share is not implemented because all objects cannot be granted to shares.
func (agb *OnAllGrantBuilder) Share(n string) GrantExecutable {
	return nil
}

// Grant returns the SQL that will grant all privileges on the grant to the grantee
func (age *OnAllGrantExecutable) Grant(p string) string {
	return fmt.Sprintf(`GRANT %v ON ALL %vS IN SCHEMA %v TO ROLE "%v"`,
		p, age.onAllGrantType, age.grantName, age.granteeName)
}

// Revoke returns the SQL that will revoke all privileges on the grant from the grantee
func (age *OnAllGrantExecutable) Revoke(p string) string {
	return fmt.Sprintf(`REVOKE %v ON ALL %vS IN SCHEMA %v FROM ROLE "%v"`,
		p, age.onAllGrantType, age.grantName, age.granteeName)
}

// Show returns the SQL that will show all all grants on the schema
func (age *OnAllGrantExecutable) Show() string {
	return fmt.Sprintf(`SHOW GRANTS TO ROLE %v`, age.granteeName)
}
