package snowflake

import "fmt"

type RoleOwnershipGrantBuilder struct {
	role          string
	currentGrants string
}

type RoleOwnershipGrantExecutable struct {
	role          string
	granteeType   granteeType
	grantee       string
	currentGrants string
}

func RoleOwnershipGrant(role string, currentGrants string) *RoleOwnershipGrantBuilder {
	return &RoleOwnershipGrantBuilder{role: role, currentGrants: currentGrants}
}

func (gb *RoleOwnershipGrantBuilder) Role(role string) *RoleOwnershipGrantExecutable {
	return &RoleOwnershipGrantExecutable{
		role:          gb.role,
		granteeType:   "Role",
		grantee:       role,
		currentGrants: gb.currentGrants,
	}
}

func (gr *RoleOwnershipGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}

func (gr *RoleOwnershipGrantExecutable) Revoke() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}
