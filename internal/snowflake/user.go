// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/helpers"
	"github.com/jmoiron/sqlx"
)

func NewUserBuilder(name string) *Builder {
	return &Builder{
		entityType: UserType,
		name:       name,
	}
}

type User struct {
	Comment               sql.NullString `db:"comment"`
	DefaultNamespace      sql.NullString `db:"default_namespace"`
	DefaultRole           sql.NullString `db:"default_role"`
	DefaultSecondaryRoles sql.NullString `db:"default_secondary_roles"`
	DefaultWarehouse      sql.NullString `db:"default_warehouse"`
	Disabled              bool           `db:"disabled"`
	DisplayName           sql.NullString `db:"display_name"`
	Email                 sql.NullString `db:"email"`
	FirstName             sql.NullString `db:"first_name"`
	HasRsaPublicKey       bool           `db:"has_rsa_public_key"`
	LastName              sql.NullString `db:"last_name"`
	LoginName             sql.NullString `db:"login_name"`
	Name                  sql.NullString `db:"name"`
}

func ScanUser(row *sqlx.Row) (*User, error) {
	r := &User{}
	err := row.StructScan(r)
	return r, err
}

func ScanUserDescription(rows *sqlx.Rows) (*User, error) {
	r := &User{}
	var err error

	for rows.Next() {
		userProp := &DescribeUserProp{}
		if err := rows.StructScan(userProp); err != nil {
			return nil, err
		}

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
		case "DEFAULT_SECONDARY_ROLES":
			if len(userProp.Value.String) > 0 {
				defaultSecondaryRoles := helpers.ListContentToString(userProp.Value.String)
				r.DefaultSecondaryRoles = sql.NullString{String: defaultSecondaryRoles, Valid: true}
			} else {
				r.DefaultSecondaryRoles = sql.NullString{Valid: false}
			}
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

func ListUsers(pattern string, db *sql.DB) ([]User, error) {
	stmt := fmt.Sprintf(`SHOW USERS like '%s'`, pattern)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []User{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no users found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
