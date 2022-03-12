package snowflake

import "fmt"

type RoleOwnershipGrantBuilder struct {
	name string
}

type RoleOwnershipGrantExecutable struct {
	name        string
	granteeType granteeType
	grantee     string
}

func RoleOwnershipGrant(name string) *RoleOwnershipGrantBuilder {
	return &RoleOwnershipGrantBuilder{name: name}
}

func (gb *RoleOwnershipGrantBuilder) User(user string) *RoleOwnershipGrantExecutable {
	return &RoleOwnershipGrantExecutable{
		name:        gb.name,
		granteeType: userType,
		grantee:     user,
	}
}

func (gb *RoleOwnershipGrantBuilder) Role(role string) *RoleOwnershipGrantExecutable {
	return &RoleOwnershipGrantExecutable{
		name:        gb.name,
		granteeType: roleType,
		grantee:     role,
	}
}

func (gr *RoleOwnershipGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" COPY CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.name) // nolint: gosec
}

func (gr *RoleOwnershipGrantExecutable) Revoke() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "ACCOUNTADMIN" COPY CURRENT GRANTS`, gr.granteeType, gr.grantee) // nolint: gosec
}
