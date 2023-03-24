package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// DatabaseRoleBuilder abstracts the creation of sql queries for a snowflake database role.
type DatabaseRoleBuilder struct {
	name    string
	db      string
	comment string
}

// GetFullName prepends db and schema to in parameter.
func (builder *DatabaseRoleBuilder) GetFullName(name string) string {
	var n strings.Builder
	n.WriteString(fmt.Sprintf(`"%v"."%v"`, builder.db, name))
	return n.String()
}

// QualifiedName prepends the db and schema and escapes everything nicely.
func (builder *DatabaseRoleBuilder) QualifiedName() string {
	return builder.GetFullName(builder.name)
}

// Name returns the name of the database role.
func (builder *DatabaseRoleBuilder) Name() string {
	return builder.name
}

// WithComment adds a comment to the DatabaseRoleBuilder.
func (builder *DatabaseRoleBuilder) WithComment(c string) *DatabaseRoleBuilder {
	builder.comment = c
	return builder
}

// Database role returns a pointer to a Builder that abstracts the DDL operations for a database role.
//
// Supported DDL operations are:
//   - CREATE DATABASE ROLE
//   - ALTER DATABASE ROLE
//   - DROP DATABASE ROLE
//   - DESCRIBE DATABASE ROLE
//   - SHOW DATABASE ROLES
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/security-access-control-considerations#label-access-control-considerations-database-roles)

func NewDatabaseRoleBuilder(name, db string) *DatabaseRoleBuilder {
	return &DatabaseRoleBuilder{
		name: name,
		db:   db,
	}
}

// Create returns the SQL that will create a new database role.
func (builder *DatabaseRoleBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)
	q.WriteString(fmt.Sprintf(` DATABASE ROLE %v`, builder.QualifiedName()))

	if builder.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(builder.comment)))
	}
	return q.String()
}

// ChangeComment returns the sql that will change the comment for the database role.
func (builder *DatabaseRoleBuilder) ChangeComment(newComment string) string {
	return fmt.Sprintf(`ALTER DATABASE ROLE %v SET COMMENT = '%v'`, builder.QualifiedName(), EscapeString(newComment))
}

// Drop returns the sql that will remove the Database Role.
func (builder *DatabaseRoleBuilder) Drop() string {
	return fmt.Sprintf(`DROP DATABASE ROLE %v`, builder.QualifiedName())
}

// Describe returns the sql that will describe a Database Role.
func (builder *DatabaseRoleBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE DATABASE ROLE %v`, builder.QualifiedName())
}

// Show returns the sql that will show a Database Role.
func (builder *DatabaseRoleBuilder) Show() string {
	return fmt.Sprintf(`SHOW DATABASE ROLES IN DATABASE "%v"`, EscapeString(builder.db))
}

type DatabaseRole struct {
	CreatedOn    string  `db:"created_on"`
	Name         string  `db:"name"`
	DatabaseName string  `db:"database_name"`
	Owner        string  `db:"owner"`
	Comment      *string `db:"comment"`
}

func (dr *DatabaseRole) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"`, EscapeString(dr.DatabaseName), EscapeString(dr.Name))
}

// ScanDatabaseRole turns a sql row into a database role object.
func ScanDatabaseRole(row *sqlx.Row) (*DatabaseRole, error) {
	dr := &DatabaseRole{}
	e := row.StructScan(dr)
	return dr, e
}

func ListDatabaseRoles(databaseName string, db *sql.DB) ([]DatabaseRole, error) {
	stmt := fmt.Sprintf(`SHOW DATABASE ROLES IN DATABASE "%s"`, databaseName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	databaseRoles := []DatabaseRole{}
	if err := sqlx.StructScan(rows, &databaseRoles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no Database Roles found")
			return nil, nil
		}
		return databaseRoles, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return databaseRoles, nil
}
