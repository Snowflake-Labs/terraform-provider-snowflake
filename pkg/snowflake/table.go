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
	name      string
	db        string
	schema    string
	columns   Columns
	comment   string
	clusterBy []string
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

// WithClustering adds cluster keys/expressions to TableBuilder
func (tb *TableBuilder) WithClustering(c []string) *TableBuilder {
	tb.clusterBy = c
	return tb
}

//Function to get clustering definition
func (tb *TableBuilder) GetClusterKeyString() string {

	return fmt.Sprint(strings.Join(tb.clusterBy[:], ", "))
}

//function to take the literal snowflake cluster statement returned from SHOW TABLES and convert it to a list of keys.
func ClusterStatementToList(clusterStatement string) []string {
	if clusterStatement == "" {
		return nil
	}

	cleanStatement := strings.TrimSuffix(strings.Replace(clusterStatement, "LINEAR(", "", 1), ")")
	// remove cluster statement and trailing parenthesis

	var clean []string

	for _, s := range strings.Split(cleanStatement, ",") {
		clean = append(clean, strings.TrimSpace(s))
	}

	return clean

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

	if tb.clusterBy != nil {
		//add optional clustering statement
		q.WriteString(fmt.Sprintf(` CLUSTER BY LINEAR(%v)`, tb.GetClusterKeyString()))

	}

	return q.String()
}

// ChangeClusterBy returns the SQL query to change cluastering on table
func (tb *TableBuilder) ChangeClusterBy(cb string) string {
	return fmt.Sprintf(`ALTER TABLE %v CLUSTER BY LINEAR(%v)`, tb.QualifiedName(), cb)
}

// ChangeComment returns the SQL query that will update the comment on the table.
func (tb *TableBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TABLE %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(c))
}

// AddColumn returns the SQL query that will add a new column to the table.
func (tb *TableBuilder) AddColumn(name string, dataType string) string {
	col := Column{
		name:  name,
		_type: dataType,
	}
	return fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition())
}

// DropColumn returns the SQL query that will add a new column to the table.
func (tb *TableBuilder) DropColumn(name string) string {
	return fmt.Sprintf(`ALTER TABLE %s DROP COLUMN "%s"`, tb.QualifiedName(), name)
}

// ChangeColumnType returns the SQL query that will change the type of the named column to the given type.
func (tb *TableBuilder) ChangeColumnType(name string, dataType string) string {
	col := Column{
		name:  name,
		_type: dataType,
	}
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition())
}

// RemoveComment returns the SQL query that will remove the comment on the table.
func (tb *TableBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TABLE %v UNSET COMMENT`, tb.QualifiedName())
}

// RemoveClustering returns the SQL query that will remove data clustering from the table
func (tb *TableBuilder) DropClustering() string {
	return fmt.Sprintf(`ALTER TABLE %v DROP CLUSTERING KEY`, tb.QualifiedName())
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
