package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TableBuilder abstracts the creation of SQL queries for a Snowflake schema
type TableBuilder struct {
	name    string
	db      string
	schema  string
	columns []map[string]string
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

// WithColumns sets the column definitions on the TableBuilder
func (tb *TableBuilder) WithColumns(c []map[string]string) *TableBuilder {
	tb.columns = c
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
func TableWithColumnDefinitions(name, db, schema string, columns []map[string]string) *TableBuilder {
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
	q.WriteString(fmt.Sprintf(`CREATE TABLE %v`, tb.QualifiedName()))

	q.WriteString(fmt.Sprintf(` (`))
	columnDefinitions := []string{}
	for _, columnDefinition := range tb.columns {
		columnDefinitions = append(columnDefinitions, fmt.Sprintf(`"%v" %v`, EscapeString(columnDefinition["name"]), EscapeString(columnDefinition["type"])))
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
	CreatedOn           sql.NullString `db:"created_on"`
	TableName           sql.NullString `db:"name"`
	DatabaseName        sql.NullString `db:"database_name"`
	SchemaName          sql.NullString `db:"schema_name"`
	Kind                sql.NullString `db:"kind"`
	Comment             sql.NullString `db:"comment"`
	ClusterBy           sql.NullString `db:"cluster_by"`
	Rows                sql.NullString `db:"row"`
	Bytes               sql.NullString `db:"bytes"`
	Owner               sql.NullString `db:"owner"`
	RetentionTime       sql.NullString `db:"retention_time"`
	AutomaticClustering sql.NullString `db:"automatic_clustering"`
	ChangeTracking      sql.NullString `db:"change_tracking"`
}

func ScanTable(row *sqlx.Row) (*table, error) {
	t := &table{}
	e := row.StructScan(t)
	return t, e
}
