package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	pe "github.com/pkg/errors"
)

// ViewBuilder abstracts the creation of SQL queries for a Snowflake View
type ViewBuilder struct {
	name      string
	db        string
	schema    string
	secure    bool
	replace   bool
	comment   string
	statement string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (vb *ViewBuilder) QualifiedName() (string, error) {
	if vb.db == "" || vb.schema == "" {
		return "", errors.New("Views must specify a database and a schema")
	}

	return fmt.Sprintf(`"%v"."%v"."%v"`, vb.db, vb.schema, vb.name), nil
}

// WithComment adds a comment to the ViewBuilder
func (vb *ViewBuilder) WithComment(c string) *ViewBuilder {
	vb.comment = c
	return vb
}

// WithDB adds the name of the database to the ViewBuilder
func (vb *ViewBuilder) WithDB(db string) *ViewBuilder {
	vb.db = db
	return vb
}

// WithReplace adds the "OR REPLACE" option to the ViewBuilder
func (vb *ViewBuilder) WithReplace() *ViewBuilder {
	vb.replace = true
	return vb
}

// WithSchema adds the name of the schema to the ViewBuilder
func (vb *ViewBuilder) WithSchema(s string) *ViewBuilder {
	vb.schema = s
	return vb
}

// WithSecure sets the secure boolean to true
// [Snowflake Reference](https://docs.snowflake.net/manuals/user-guide/views-secure.html)
func (vb *ViewBuilder) WithSecure() *ViewBuilder {
	vb.secure = true
	return vb
}

// WithStatement adds the SQL statement to be used for the view
func (vb *ViewBuilder) WithStatement(s string) *ViewBuilder {
	vb.statement = s
	return vb
}

// View returns a pointer to a Builder that abstracts the DDL operations for a view.
//
// Supported DDL operations are:
//   - CREATE VIEW
//   - ALTER VIEW
//   - DROP VIEW
//   - SHOW VIEWS
//   - DESCRIBE VIEW
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-table.html#standard-view-management)
func View(name string) *ViewBuilder {
	return &ViewBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new view.
func (vb *ViewBuilder) Create() (string, error) {
	var q strings.Builder

	q.WriteString("CREATE")

	if vb.replace {
		q.WriteString(" OR REPLACE")
	}

	if vb.secure {
		q.WriteString(" SECURE")
	}

	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}

	q.WriteString(fmt.Sprintf(` VIEW %v`, qn))

	if vb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", EscapeString(vb.comment)))
	}

	q.WriteString(fmt.Sprintf(" AS %v", vb.statement))

	return q.String(), nil
}

// Rename returns the SQL query that will rename the view.
func (vb *ViewBuilder) Rename(newName string) (string, error) {
	oldName, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	vb.name = newName

	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER VIEW %v RENAME TO %v`, oldName, qn), nil
}

// Secure returns the SQL query that will change the view to a secure view.
func (vb *ViewBuilder) Secure() (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER VIEW %v SET SECURE`, qn), nil
}

// Unsecure returns the SQL query that will change the view to a normal (unsecured) view.
func (vb *ViewBuilder) Unsecure() (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER VIEW %v UNSET SECURE`, qn), nil
}

// ChangeComment returns the SQL query that will update the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) ChangeComment(c string) (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`ALTER VIEW %v SET COMMENT = '%v'`, qn, EscapeString(c)), nil
}

// RemoveComment returns the SQL query that will remove the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) RemoveComment() (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER VIEW %v UNSET COMMENT`, qn), nil
}

// Show returns the SQL query that will show the row representing this view.
func (vb *ViewBuilder) Show() string {
	return fmt.Sprintf(`SHOW VIEWS LIKE '%v' IN SCHEMA "%v"."%v"`, vb.name, vb.db, vb.schema)
}

// Drop returns the SQL query that will drop the row representing this view.
func (vb *ViewBuilder) Drop() (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`DROP VIEW %v`, qn), nil
}

type view struct {
	Comment      sql.NullString `db:"comment"`
	IsSecure     bool           `db:"is_secure"`
	Name         sql.NullString `db:"name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Text         sql.NullString `db:"text"`
	DatabaseName sql.NullString `db:"database_name"`
}

func ScanView(row *sqlx.Row) (*view, error) {
	r := &view{}
	err := row.StructScan(r)
	return r, err
}

func ListViews(databaseName string, schemaName string, db *sql.DB) ([]view, error) {
	stmt := fmt.Sprintf(`SHOW VIEWS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []view{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no views found")
		return nil, nil
	}
	return dbs, pe.Wrapf(err, "unable to scan row for %s", stmt)
}
