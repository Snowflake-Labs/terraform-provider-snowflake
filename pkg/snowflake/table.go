package snowflake

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

//CREATE TABLE "FULFIL_TEST"."PUBLIC"."TABLENAME" ("C1" STRING);

// TableBuilder abstracts the creation of SQL queries for a Snowflake schema
type TableBuilder struct {
	name    string
	db      string
	schema  string
	columns map[string]string
	comment string
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

// WithComment adds a comment to the TableBuilder
func (tb *TableBuilder) WithComment(c string) *TableBuilder {
	tb.comment = c
	return tb
}

// Table returns a pointer to a Builder that abstracts the DDL operations for a table.
//
// Supported DDL operations are:
//   - ALTER TABLE
//   - DROP TABLE
//   - SHOW TABLES
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-table.html)
func Table(name, db, schema string) *TableBuilder {
	return &TableBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Table returns a pointer to a Builder that abstracts the DDL operations for a table.
//
// Supported DDL operations are:
//   - CREATE TABLE
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-table.html)
func TableWithColumnDefinition(name, db, schema string, columns map[string]string) *TableBuilder {
	return &TableBuilder{
		name:    name,
		db:      db,
		schema:  schema,
		columns: columns,
	}
}

// Create returns the SQL statement required to create a table
func (tb *TableBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TABLE %v`, tb.QualifiedName()))
	q.WriteString(fmt.Sprintf(` (`))
	columnDefinitions := []string{}
	for columnName, columnType := range tb.columns {
		columnDefinitions = append(columnDefinitions, fmt.Sprintf(`"%v" %v`, columnName, columnType))
	}
	q.WriteString(strings.Join(columnDefinitions, ", "))
	q.WriteString(fmt.Sprintf(`)`))

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	return q.String()
}

// ChangeComment returns the SQL query that will update the comment on the table.
func (tb *TableBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TABLE %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the table.
func (tb *TableBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TABLE %v UNSET COMMENT`, tb.QualifiedName())
}

// Drop returns the SQL query that will drop a table.
func (tb *TableBuilder) Drop() string {
	return fmt.Sprintf(`DROP TABLE %v`, tb.QualifiedName())
}

// Show returns the SQL query that will show a table.
func (tb *TableBuilder) Show() string {
	return fmt.Sprintf(`SHOW TABLES LIKE '%v' IN DATABASE "%v"`, tb.name, tb.db)
}

type table struct {
	Name       string `db:"name"`
	Type       string `db:"type"`
	Kind       string `db:"kind"`
	Null       string `db:"null?"`
	PrimaryKey string `db:"primary key"`
	UniqueKey  string `db:"unique key"`
	Check      string `db:"check"`
	Expression string `db:"expression"`
	Comment    string `db:"comment"`
}

func ScanTable(row *sqlx.Row) (*table, error) {
	t := &table{}
	e := row.StructScan(t)
	return t, e
}
