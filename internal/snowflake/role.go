// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewRoleBuilder(db *sql.DB, name string) *RoleBuilder {
	return &RoleBuilder{
		db:   db,
		name: name,
	}
}

type RoleBuilder struct {
	name    string
	comment string
	tags    []TagValue
	db      *sql.DB
}

func (b *RoleBuilder) WithName(name string) *RoleBuilder {
	b.name = name
	return b
}

func (b *RoleBuilder) WithComment(comment string) *RoleBuilder {
	b.comment = comment
	return b
}

func (b *RoleBuilder) WithTags(tags []TagValue) *RoleBuilder {
	b.tags = tags
	return b
}

func (b *RoleBuilder) Create() error {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE ROLE "%v"`, b.name))
	if b.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", b.comment))
	}
	if len(b.tags) > 0 {
		q.WriteString(" TAG (")
		for i, tag := range b.tags {
			q.WriteString(fmt.Sprintf(`"%v"."%v"."%v" = "%v"`, tag.Database, tag.Schema, tag.Name, tag.Value))
			if i < len(b.tags)-1 {
				q.WriteString(", ")
			}
		}
		q.WriteString(")")
	}
	_, err := b.db.Exec(q.String())
	return err
}

func (b *RoleBuilder) SetComment(comment string) error {
	q := fmt.Sprintf(`ALTER ROLE "%s" SET COMMENT = '%v'`, b.name, comment)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) UnsetComment() error {
	q := fmt.Sprintf(`ALTER ROLE "%v" UNSET COMMENT`, b.name)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) UnsetTag(tag TagValue) error {
	q := fmt.Sprintf(`ALTER ROLE %s UNSET TAG "%v"."%v"."%v"`, b.name, tag.Database, tag.Schema, tag.Name)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) SetTag(tag TagValue) error {
	q := fmt.Sprintf(`ALTER ROLE %s SET TAG  "%v"."%v"."%v" = "%v"`, b.name, tag.Database, tag.Schema, tag.Name, tag.Value)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) ChangeTag(tag TagValue) error {
	q := fmt.Sprintf(`ALTER ROLE "%s" SET TAG "%v"."%v"."%v" = "%v"`, b.name, tag.Database, tag.Schema, tag.Name, tag.Value)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) Drop() error {
	q := fmt.Sprintf(`DROP ROLE "%s"`, b.name)
	_, err := b.db.Exec(q)
	return err
}

func (b *RoleBuilder) Show() (*Role, error) {
	stmt := fmt.Sprintf(`SHOW ROLES LIKE '%s'`, b.name)
	row := QueryRow(b.db, stmt)
	r := &Role{}
	err := row.StructScan(r)
	return r, err
}

func (b *RoleBuilder) Rename(newName string) error {
	stmt := fmt.Sprintf(`ALTER ROLE "%s" RENAME TO "%s"`, b.name, newName)
	_, err := b.db.Exec(stmt)
	return err
}

type Role struct {
	Name    sql.NullString `db:"name"`
	Comment sql.NullString `db:"comment"`
	Owner   sql.NullString `db:"owner"`
}

func ListRoles(db *sql.DB, rolePattern string) ([]*Role, error) {
	stmt := strings.Builder{}
	stmt.WriteString("SHOW ROLES")
	if rolePattern != "" {
		stmt.WriteString(fmt.Sprintf(` LIKE '%v'`, rolePattern))
	}
	rows, err := Query(db, stmt.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*Role{}
	if err := sqlx.StructScan(rows, &roles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no roles found")
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan stmt = %v err = %w", stmt, err)
	}
	return roles, nil
}
