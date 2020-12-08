package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Column struct {
	name  string
	_type string // type is reserved
}

func (c *Column) WithName(name string) *Column {
	c.name = name
	return c
}
func (c *Column) WithType(t string) *Column {
	c._type = t
	return c
}

func (c *Column) getColumnDefinition() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf(`"%v" %v`, EscapeString(c.name), EscapeString(c._type))
}

type Columns []Column

// NewColumns generates columns from a table description
func NewColumns(tds []tableDescription) Columns {
	cs := []Column{}
	for _, td := range tds {
		if td.Kind.String != "COLUMN" {
			continue
		}
		cs = append(cs, Column{
			name:  td.Name.String,
			_type: td.Type.String,
		})
	}
	return Columns(cs)
}

func (c Columns) Flatten() []interface{} {
	flattened := []interface{}{}
	for _, col := range c {
		flat := map[string]interface{}{}
		flat["name"] = col.name
		flat["type"] = col._type

		flattened = append(flattened, flat)
	}
	return flattened
}

func (c Columns) getColumnDefinitions() string {
	// TODO(el): verify Snowflake reflects column order back in desc table calls
	columnDefinitions := []string{}
	for _, column := range c {
		columnDefinitions = append(columnDefinitions, column.getColumnDefinition())
	}

	// NOTE: intentionally blank leading space
	return fmt.Sprintf(" (%s)", strings.Join(columnDefinitions, ", "))
}

// TableBuilder abstracts the creation of SQL queries for a Snowflake schema
type TableBuilder struct {
	name    string
	db      string
	schema  string
	columns Columns
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
func (tb *TableBuilder) WithColumns(c Columns) *TableBuilder {
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
func TableWithColumnDefinitions(name, db, schema string, columns Columns) *TableBuilder {
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
	q.WriteString(tb.columns.getColumnDefinitions())

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
	return fmt.Sprintf(`SHOW TABLES LIKE '%v' IN SCHEMA "%v"."%v"`, tb.name, tb.db, tb.schema)
}

func (tb *TableBuilder) ShowColumns() string {
	return fmt.Sprintf(`DESC TABLE %s`, tb.QualifiedName())
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

type tableDescription struct {
	Name sql.NullString `db:"name"`
	Type sql.NullString `db:"type"`
	Kind sql.NullString `db:"kind"`
}

func ScanTableDescription(rows *sqlx.Rows) ([]tableDescription, error) {
	tds := []tableDescription{}
	for rows.Next() {
		td := tableDescription{}
		err := rows.StructScan(&td)
		if err != nil {
			return nil, err
		}
		tds = append(tds, td)
	}
	return tds, rows.Err()
}
