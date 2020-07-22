package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
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
func (vb *ViewBuilder) QualifiedName() string {
	var n strings.Builder

	if vb.db != "" && vb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, vb.db, vb.schema))
	}

	if vb.db != "" && vb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, vb.db))
	}

	if vb.db == "" && vb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, vb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, vb.name))

	return n.String()
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
func (vb *ViewBuilder) Create() string {
	var q strings.Builder

	q.WriteString("CREATE")

	if vb.replace {
		q.WriteString(" OR REPLACE")
	}

	if vb.secure {
		q.WriteString(" SECURE")
	}

	q.WriteString(fmt.Sprintf(` VIEW %v`, vb.QualifiedName()))

	if vb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", vb.comment))
	}

	q.WriteString(fmt.Sprintf(" AS %v", vb.statement))

	return q.String()
}

// Rename returns the SQL query that will rename the view.
func (vb *ViewBuilder) Rename(newName string) string {
	oldName := vb.QualifiedName()
	vb.name = newName
	return fmt.Sprintf(`ALTER VIEW %v RENAME TO %v`, oldName, vb.QualifiedName())
}

// Secure returns the SQL query that will change the view to a secure view.
func (vb *ViewBuilder) Secure() string {
	return fmt.Sprintf(`ALTER VIEW %v SET SECURE`, vb.QualifiedName())
}

// Unsecure returns the SQL query that will change the view to a normal (unsecured) view.
func (vb *ViewBuilder) Unsecure() string {
	return fmt.Sprintf(`ALTER VIEW %v UNSET SECURE`, vb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER VIEW %v SET COMMENT = '%v'`, vb.QualifiedName(), c)
}

// RemoveComment returns the SQL query that will remove the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER VIEW %v UNSET COMMENT`, vb.QualifiedName())
}

// Show returns the SQL query that will show the row representing this view.
func (vb *ViewBuilder) Show() string {
	if vb.db == "" {
		return fmt.Sprintf(`SHOW VIEWS LIKE '%v'`, vb.name)
	}
	return fmt.Sprintf(`SHOW VIEWS LIKE '%v' IN DATABASE "%v"`, vb.name, vb.db)
}

// Drop returns the SQL query that will drop the row representing this view.
func (vb *ViewBuilder) Drop() string {
	return fmt.Sprintf(`DROP VIEW %v`, vb.QualifiedName())
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
