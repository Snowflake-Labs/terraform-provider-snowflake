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
	if vb.secure {
		q.WriteString(" SECURE")
	}
	q.WriteString(fmt.Sprintf(` VIEW "%v"`, vb.name))

	if vb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", vb.comment))
	}

	q.WriteString(fmt.Sprintf(" AS %v", vb.statement))

	return q.String()
}

// Rename returns the SQL query that will rename the view.
func (vb *ViewBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER VIEW "%v" RENAME TO "%v"`, vb.name, newName)
}

// Secure returns the SQL query that will change the view to a secure view.
func (vb *ViewBuilder) Secure() string {
	return fmt.Sprintf(`ALTER VIEW "%v" SET SECURE`, vb.name)
}

// Unsecure returns the SQL query that will change the view to a normal (unsecured) view.
func (vb *ViewBuilder) Unsecure() string {
	return fmt.Sprintf(`ALTER VIEW "%v" UNSET SECURE`, vb.name)
}

// ChangeComment returns the SQL query that will update the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER VIEW "%v" SET COMMENT = '%v'`, vb.name, c)
}

// RemoveComment returns the SQL query that will remove the comment on the view.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (vb *ViewBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER VIEW "%v" UNSET COMMENT`, vb.name)
}

// Show returns the SQL query that will show the row representing this view.
func (vb *ViewBuilder) Show() string {
	return fmt.Sprintf(`SHOW VIEWS LIKE '%v'`, vb.name)
}

// Drop returns the SQL query that will drop the row representing this view.
func (vb *ViewBuilder) Drop() string {
	return fmt.Sprintf(`DROP VIEW "%v"`, vb.name)
}
