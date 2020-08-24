package snowflake

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func User(name string) *Builder {
	return &Builder{
		entityType: UserType,
		name:       name,
	}
}

type user struct {
	Comment          sql.NullString `db:"comment"`
	DefaultNamespace sql.NullString `db:"default_namespace"`
	DefaultRole      sql.NullString `db:"default_role"`
	DefaultWarehouse sql.NullString `db:"default_warehouse"`
	Disabled         bool           `db:"disabled"`
	DisplayName      sql.NullString `db:"display_name"`
	Email            sql.NullString `db:"email"`
	FirstName        sql.NullString `db:"first_name"`
	HasRsaPublicKey  bool           `db:"has_rsa_public_key"`
	LastName         sql.NullString `db:"last_name"`
	LoginName        sql.NullString `db:"login_name"`
	Name             sql.NullString `db:"name"`
}

func ScanUser(row *sqlx.Row) (*user, error) {
	r := &user{}
	err := row.StructScan(r)
	return r, err
}
