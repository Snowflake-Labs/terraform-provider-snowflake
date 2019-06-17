package snowflake

import (
	"fmt"
	"strings"
)

// ViewBuilder abstracts the creation of SQL queries for a Snowflake View
type ViewBuilder struct {
	name          string
	secure        bool
	comment       string
	statement     string
	statementArgs []interface{}
}

// WithComment adds a comment to the ViewBuilder
func (vb *ViewBuilder) WithComment(c string) *ViewBuilder {
	vb.comment = c
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

//WithStatmentArgs adds the args to be passed to db.Exec
func (vb *ViewBuilder) WithStatementArgs(args []interface{}) *ViewBuilder {
	vb.statementArgs = args
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

// Create returns the SQL query and args suitable for passing directly to
// db.Exec that will create a new view.
func (vb *ViewBuilder) Create() (string, []interface{}) {
	var q strings.Builder
	var args []interface{}
	q.WriteString("CREATE")
	if vb.secure {
		q.WriteString(" SECURE")
	}
	q.WriteString(" VIEW ?")
	args = append(args, vb.name)

	if vb.comment != "" {
		q.WriteString(" COMMENT = ?")
		args = append(args, vb.comment)
	}

	// The onus is on the user to provide a clean query statement at the moment
	q.WriteString(fmt.Sprintf(" AS %s", vb.statement))

	return q.String(), append(args, vb.statementArgs...)
}

// Rename returns the SQL query and args suitable for passing directly to
// db.Exec that will rename the view.
func (vb *ViewBuilder) Rename(newName string) (string, []interface{}) {
	return "ALTER VIEW ? RENAME TO ?", []interface{}{vb.name, newName}
}

// Secure returns the SQL query and args suitable for passing directly to
// db.Exec that will change the view to a secure view.
func (vb *ViewBuilder) Secure() (string, []interface{}) {
	return "ALTER VIEW ? SET SECURE", []interface{}{vb.name}
}

// Unsecure returns the SQL query and args suitable for passing directly to
// db.Exec that will change the view to a normal (unsecured) view.
func (vb *ViewBuilder) Unsecure() (string, []interface{}) {
	return "ALTER VIEW ? UNSET SECURE", []interface{}{vb.name}
}

// ChangeComment returns the SQL query and args suitable for passing directly to
// db.Exec that will update the comment on the view. Note that comment is the
// only parameter, if more are released this should be abstracted as per the
// generic builder.
func (vb *ViewBuilder) ChangeComment(c string) (string, []interface{}) {
	return "ALTER VIEW ? SET COMMENT = ?", []interface{}{vb.name, c}
}

// RemoveComment returns the SQL query and args suitable for passing directly to
// db.Exec that will remove the comment on the view. Note that comment is the
// only parameter, if more are released this should be abstracted as per the
// generic builder.
func (vb *ViewBuilder) RemoveComment() (string, []interface{}) {
	return "ALTER VIEW ? UNSET COMMENT", []interface{}{vb.name}
}

// Show returns the SQL query and args suitable for passing directly to db.Exec
// that will show the row representing this view.
func (vb *ViewBuilder) Show() (string, []interface{}) {
	return "SHOW VIEWS LIKE ?", []interface{}{vb.name}
}
