package snowflake

import (
	"fmt"
	"strings"
)

// TableBuilder abstracts the creation of SQL queries for a Snowflake Table
type TableBuilder struct {
	name      	string
	db        	string
	schema    	string
	columns		map[string]string // {column_name: column_type}
	comment   	string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (tb *TableBuilder) QualifiedName() string {
	var n strings.Builder

	if tb.db != "" && tb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, tb.db, tb.schema))
	}

	if tb.db != "" && tb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, tb.db))
	}

	if tb.db == "" && tb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, tb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, tb.name))

	return n.String()
}

// ColumnsStatement creates the statement to create the different columns
func (tb *TableBuilder) ColumnsStatement() string {
	var s strings.Builder
	var c []string

	for cn, ct := range tb.columns {
		append(c, fmt.Sprintf(`"%v" %v`, cn, ct))
	}

	s.writeString(fmt.Sprintf("(%v)", strings.Join(", ", c)))

	return s.String()
}

// WithComment adds a comment to the TableBuilder
func (tb *TableBuilder) WithComment(c string) *TableBuilder {
	tb.comment = c
	return tb
}

// WithDB adds the name of the database to the TableBuilder
func (tb *TableBuilder) WithDB(db string) *TableBuilder {
	tb.db = db
	return tb
}

// WithSchema adds the name of the schema to the TableBuilder
func (tb *TableBuilder) WithSchema(s string) *TableBuilder {
	tb.schema = s
	return tb
}

// WithColumns adds the columns map to the table
func (tb *TableBuilder) WithColumns(c map[string]string) *TableBuilder {
	tb.columns = c
	return tb
}

// Table returns a pointer to a Builder that abstracts the DDL operations for a table.
//
// Supported DDL operations are:
//   - CREATE TABLE
//   - ALTER TABLE
//   - SHOW TABLE
//   - DROP TABLE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-table.html#table-management)
func Table(name string) *TableBuilder {
	return &TableBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new table.
func (tb *TableBuilder) Create() string {
	var q strings.Builder

	q.WriteString("CREATE")

	q.WriteString(fmt.Sprintf(` TABLE %v`, tb.QualifiedName()))

	q.WriteString(fmt.Sprintf(" %v", tb.ColumnsStatement()))

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", tb.comment))
	}	

	return q.String()
}

// Rename returns the SQL query that will rename the table.
func (tb *TableBuilder) Rename(newName string) string {
	oldName := tb.QualifiedName()
	tb.name = newName
	return fmt.Sprintf(`ALTER TABLE %v RENAME TO %v`, oldName, tb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the table.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (tb *TableBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TABLE %v SET COMMENT = '%v'`, tb.QualifiedName(), c)
}

// RemoveComment returns the SQL query that will remove the comment on the table.
// Note that comment is the only parameter, if more are released this should be
// abstracted as per the generic builder.
func (tb *TableBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TABLE %v UNSET COMMENT`, tb.QualifiedName())
}

// Show returns the SQL query that will show the row representing this table.
func (tb *TableBuilder) Show() string {
	if tb.db == "" {
		return fmt.Sprintf(`SHOW TABLES LIKE '%v'`, tb.name)
	}
	return fmt.Sprintf(`SHOW TABLES LIKE '%v' IN DATABASE "%v"`, tb.name, tb.db)
}

// ShowColumns returns the SQL query that will show the columns of this table.
func (tb *TableBuilder) ShowColumns() string {	
	return fmt.Sprintf(`SHOW COLUMNS IN TABLE %v`, tb.QualifiedName())
}

// Drop returns the SQL query that will drop the row representing this table.
func (tb *TableBuilder) Drop() string {
	return fmt.Sprintf(`DROP TABLE %v`, tb.QualifiedName())
}
