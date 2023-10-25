package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

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

func NewRoleOwnershipGrantBuilder(role string, currentGrants string) *RoleOwnershipGrantBuilder {
	return &RoleOwnershipGrantBuilder{role: role, currentGrants: currentGrants}
}

func (gb *RoleOwnershipGrantBuilder) Role(role string) *RoleOwnershipGrantExecutable {
	return &RoleOwnershipGrantExecutable{
		role:          role,
		granteeType:   "ROLE",
		grantee:       gb.role,
		currentGrants: gb.currentGrants,
	}
}

func (gr *RoleOwnershipGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}

func (gr *RoleOwnershipGrantExecutable) Revoke() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}

type RoleOwnershipGrant struct {
	CreatedOn       sql.NullString `db:"created_on"`
	Name            sql.NullString `db:"name"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	IsInherited     sql.NullString `db:"is_inherited"`
	AssignedToUsers sql.NullString `db:"assigned_to_users"`
	GrantedToRoles  sql.NullString `db:"granted_to_roles"`
	GrantedRoles    sql.NullString `db:"granted_roles"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
}

func ScanRoleOwnershipGrant(row *sqlx.Row) (*RoleOwnershipGrant, error) {
	rog := &RoleOwnershipGrant{}
	err := row.StructScan(rog)
	return rog, err
}
