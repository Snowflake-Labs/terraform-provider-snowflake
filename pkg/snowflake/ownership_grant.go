package snowflake

import "fmt"

type OwnershipGrantBuilder struct {
	name string
}

type OwnershipGrantExecutable struct {
	name          string
	grantee       string
	currentGrants string
}

func OwnershipGrant(name string) *OwnershipGrantBuilder {
	return &OwnershipGrantBuilder{name: name}
}

func (gb *OwnershipGrantBuilder) Role(role string, currentGrants string) *OwnershipGrantExecutable {
	return &OwnershipGrantExecutable{
		name:          gb.name,
		grantee:       role,
		currentGrants: currentGrants,
	}
}

func (og *OwnershipGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT OWNERSHIP "%s" TO ROLE "%s" %s CURRENT GRANTS`, og.name, og.grantee, og.currentGrants) // nolint: gosec
}

func (og *OwnershipGrantExecutable) Revoke() string {
	// TODO [ { REVOKE | COPY } CURRENT GRANTS ]
	return fmt.Sprintf(`REVOKE OWNERSHIP "%s" TO ROLE "%s" %s CURRENT GRANTS`, og.name, og.grantee, og.currentGrants) // nolint: gosec
}
