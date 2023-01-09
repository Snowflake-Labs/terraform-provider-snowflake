package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserOwnershipGrantBuilder struct {
	user          string
	currentGrants string
}

type UserOwnershipGrantExecutable struct {
	role          string
	granteeType   granteeType
	grantee       string
	currentGrants string
}

func NewUserOwnershipGrantBuilder(user string, currentGrants string) *UserOwnershipGrantBuilder {
	return &UserOwnershipGrantBuilder{user: user, currentGrants: currentGrants}
}

func (gb *UserOwnershipGrantBuilder) Role(role string) *UserOwnershipGrantExecutable {
	return &UserOwnershipGrantExecutable{
		role:          role,
		granteeType:   "USER",
		grantee:       gb.user,
		currentGrants: gb.currentGrants,
	}
}

func (gr *UserOwnershipGrantExecutable) Grant() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}

func (gr *UserOwnershipGrantExecutable) Revoke() string {
	return fmt.Sprintf(`GRANT OWNERSHIP ON %s "%s" TO ROLE "%s" %s CURRENT GRANTS`, gr.granteeType, gr.grantee, gr.role, gr.currentGrants) // nolint: gosec
}

type UserOwnershipGrant struct {
	Name                  sql.NullString `db:"name"`
	CreatedOn             sql.NullString `db:"created_on"`
	LoginName             sql.NullString `db:"login_name"`
	DisplayName           sql.NullString `db:"display_name"`
	FirstName             sql.NullString `db:"first_name"`
	LastName              sql.NullString `db:"last_name"`
	Email                 sql.NullString `db:"email"`
	MinsToUnlock          sql.NullString `db:"mins_to_unlock"`
	DaysToExpiry          sql.NullString `db:"days_to_expiry"`
	Comment               sql.NullString `db:"comment"`
	Disabled              sql.NullString `db:"disabled"`
	MustChangePassword    sql.NullString `db:"must_change_password"`
	SnowflakeLock         sql.NullString `db:"snowflake_lock"`
	DefaultWarehouse      sql.NullString `db:"default_warehouse"`
	DefaultNamespace      sql.NullString `db:"default_namespace"`
	DefaultRole           sql.NullString `db:"default_role"`
	DefaultSecondaryRoles sql.NullString `db:"default_secondary_roles"`
	ExtAuthnDuo           sql.NullString `db:"ext_authn_duo"`
	ExtAuthnUID           sql.NullString `db:"ext_authn_uid"`
	MinsToBypassMFA       sql.NullString `db:"mins_to_bypass_mfa"`
	Owner                 sql.NullString `db:"owner"`
	LastSuccessLogin      sql.NullString `db:"last_success_login"`
	ExpiresAtTime         sql.NullString `db:"expires_at_time"`
	LockedUntilTime       sql.NullString `db:"locked_until_time"`
	HasPassword           sql.NullString `db:"has_password"`
	HasRsaPublicKey       sql.NullString `db:"has_rsa_public_key"`
}

func ScanUserOwnershipGrant(row *sqlx.Row) (*UserOwnershipGrant, error) {
	uog := &UserOwnershipGrant{}
	err := row.StructScan(uog)
	return uog, err
}
