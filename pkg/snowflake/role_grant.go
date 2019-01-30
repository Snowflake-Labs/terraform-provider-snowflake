package snowflake

import "fmt"

type RoleGrantBuilder struct {
	name string
}

type granteeType string

const (
	userType granteeType = "USER"
	roleType granteeType = "ROLE"
)

type RoleGrantExecutable struct {
	name        string
	granteeType granteeType
	grantee     string
}

func RoleGrant(name string) *RoleGrantBuilder {
	return &RoleGrantBuilder{name: name}
}

func (gb *RoleGrantBuilder) User(user string) *RoleGrantExecutable {
	return &RoleGrantExecutable{
		name:        gb.name,
		granteeType: userType,
		grantee:     user,
	}
}

func (gb *RoleGrantBuilder) Role(role string) *RoleGrantExecutable {
	return &RoleGrantExecutable{
		name:        gb.name,
		granteeType: roleType,
		grantee:     role,
	}
}

func (gr *RoleGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT ROLE %s TO %s %s`, gr.name, gr.granteeType, gr.grantee)
}

func (gr *RoleGrantExecutable) Revoke() string {
	return fmt.Sprintf(`REVOKE ROLE %s FROM %s %s`, gr.name, gr.granteeType, gr.grantee)
}
