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

func ScanUserDescription(rows *sqlx.Rows) (*user, error) {
	r := &user{}
	var err error

	for rows.Next() {
		userProp := &DescribeUserProp{}
		err = rows.StructScan(userProp)

		// The "DESCRIBE USER ..." command returns the string "null" for null values
		if userProp.Value.String == "null" {
			userProp.Value.Valid = false
			userProp.Value.String = ""
		}

		switch propery := userProp.Property; propery {
		case "COMMENT":
			r.Comment = userProp.Value
		case "DEFAULT_NAMESPACE":
			r.DefaultNamespace = userProp.Value
		case "DEFAULT_ROLE":
			r.DefaultRole = userProp.Value
		case "DEFAULT_WAREHOUSE":
			r.DefaultWarehouse = userProp.Value
		case "DISABLED":
			r.Disabled = userProp.Value.String == "true"
		case "DISPLAY_NAME":
			r.DisplayName = userProp.Value
		case "EMAIL":
			r.Email = userProp.Value
		case "FIRST_NAME":
			r.FirstName = userProp.Value
		case "RSA_PUBLIC_KEY_FP":
			r.HasRsaPublicKey = userProp.Value.Valid
		case "LAST_NAME":
			r.LastName = userProp.Value
		case "LOGIN_NAME":
			r.LoginName = userProp.Value
		case "NAME":
			r.Name = userProp.Value
		}
	}

	err = rows.Err()

	return r, err
}

type DescribeUserProp struct {
	Property string         `db:"property"`
	Value    sql.NullString `db:"value"`
}
