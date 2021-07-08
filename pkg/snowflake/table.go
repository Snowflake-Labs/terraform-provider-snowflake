package snowflake

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PrimaryKey struct {
	name string
	keys []string
}

func (pk *PrimaryKey) WithName(name string) *PrimaryKey {
	pk.name = name
	return pk
}

func (pk *PrimaryKey) WithKeys(keys []string) *PrimaryKey {
	pk.keys = keys
	return pk
}

type Column struct {
	name     string
	_type    string // type is reserved
	nullable bool
}

func (c *Column) WithName(name string) *Column {
	c.name = name
	return c
}
func (c *Column) WithType(t string) *Column {
	c._type = t
	return c
}

func (c *Column) WithNullable(nullable bool) *Column {
	c.nullable = nullable
	return c
}

func (c *Column) getColumnDefinition(withInlineConstraints bool) string {

	if c == nil {
		return ""
	}
	var colDef strings.Builder
	colDef.WriteString(fmt.Sprintf(`"%v" %v`, EscapeString(c.name), EscapeString(c._type)))
	if withInlineConstraints {
		if !c.nullable {
			colDef.WriteString(` NOT NULL`)
		}
	}

	return colDef.String()

}

func FlattenTablePrimaryKey(pkds []primaryKeyDescription) []interface{} {
	flattened := []interface{}{}
	if len(pkds) == 0 {
		return flattened
	}

	sort.SliceStable(pkds, func(i, j int) bool {
		num1, _ := strconv.Atoi(pkds[i].KeySequence.String)
		num2, _ := strconv.Atoi(pkds[j].KeySequence.String)
		return num1 < num2
	})
	//sort our keys on the key sequence

	flat := map[string]interface{}{}
	var keys []string
	var name string
	var nameSet bool

	for _, pk := range pkds {
		//set as empty string, sys_constraint means it was an unnnamed constraint
		if strings.Contains(pk.ConstraintName.String, "SYS_CONSTRAINT") && !nameSet {
			name = ""
			nameSet = true
		}
		if !nameSet {
			name = pk.ConstraintName.String
			nameSet = true
		}

		keys = append(keys, pk.ColumnName.String)

	}

	flat["name"] = name
	flat["keys"] = keys
	flattened = append(flattened, flat)
	return flattened

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
			name:     td.Name.String,
			_type:    td.Type.String,
			nullable: td.IsNullable(),
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
		flat["nullable"] = col.nullable

		flattened = append(flattened, flat)
	}
	return flattened
}

func (c Columns) getColumnDefinitions(withInlineConstraints bool) string {
	// TODO(el): verify Snowflake reflects column order back in desc table calls
	columnDefinitions := []string{}
	for _, column := range c {
		columnDefinitions = append(columnDefinitions, column.getColumnDefinition(withInlineConstraints))
	}

	// NOTE: intentionally blank leading space
	return fmt.Sprintf(" (%s)", strings.Join(columnDefinitions, ", "))
}

// TableBuilder abstracts the creation of SQL queries for a Snowflake schema
type TableBuilder struct {
	name       string
	db         string
	schema     string
	columns    Columns
	comment    string
	clusterBy  []string
	primaryKey PrimaryKey
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

func (tb *TableBuilder) WithPrimaryKey(pk PrimaryKey) *TableBuilder {
	tb.primaryKey = pk
	return tb
}

//Function to get clustering definition
func (tb *TableBuilder) GetClusterKeyString() string {

	return JoinStringList(tb.clusterBy[:], ", ")
}

func JoinStringList(instrings []string, delimiter string) string {

	return fmt.Sprint(strings.Join(instrings[:], delimiter))

}

func quoteStringList(instrings []string) []string {
	var clean []string
	for _, word := range instrings {
		quoted := fmt.Sprintf(`"%s"`, word)
		clean = append(clean, quoted)

	}
	return clean

}

func (tb *TableBuilder) getCreateStatementBody() string {
	var q strings.Builder

	colDef := tb.columns.getColumnDefinitions(true)

	if len(tb.primaryKey.keys) > 0 {
		colDef = strings.TrimSuffix(colDef, ")") //strip trailing
		q.WriteString(colDef)
		if tb.primaryKey.name != "" {
			q.WriteString(fmt.Sprintf(` ,CONSTRAINT "%v" PRIMARY KEY(%v)`, tb.primaryKey.name, JoinStringList(quoteStringList(tb.primaryKey.keys), ",")))

		} else {
			q.WriteString(fmt.Sprintf(` ,PRIMARY KEY(%v)`, JoinStringList(quoteStringList(tb.primaryKey.keys), ",")))
		}

		q.WriteString(")") // add closing
	} else {
		q.WriteString(colDef)
	}

	return q.String()
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
	q.WriteString(tb.getCreateStatementBody())

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
func (tb *TableBuilder) AddColumn(name string, dataType string, nullable bool) string {
	col := Column{
		name:     name,
		_type:    dataType,
		nullable: nullable,
	}
	return fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition(true))
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
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition(false))
}

// RemoveComment returns the SQL query that will remove the comment on the table.
func (tb *TableBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TABLE %v UNSET COMMENT`, tb.QualifiedName())
}

// Return sql to set/unset null constraint on column
func (tb *TableBuilder) ChangeNullConstraint(name string, nullable bool) string {
	if nullable {
		return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%s" DROP NOT NULL`, tb.QualifiedName(), name)
	} else {
		return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%s" SET NOT NULL`, tb.QualifiedName(), name)
	}
}

func (tb *TableBuilder) ChangePrimaryKey(newPk PrimaryKey) string {
	tb.WithPrimaryKey(newPk)
	pks := JoinStringList(quoteStringList(newPk.keys), ", ")
	if tb.primaryKey.name != "" {
		return fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT "%v" PRIMARY KEY(%v)`, tb.QualifiedName(), tb.primaryKey.name, pks)
	}
	return fmt.Sprintf(`ALTER TABLE %s ADD PRIMARY KEY(%v)`, tb.QualifiedName(), pks)
}

func (tb *TableBuilder) DropPrimaryKey() string {
	return fmt.Sprintf(`ALTER TABLE %s DROP PRIMARY KEY`, tb.QualifiedName())
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

func (tb *TableBuilder) ShowPrimaryKeys() string {
	return fmt.Sprintf(`SHOW PRIMARY KEYS IN TABLE %s`, tb.QualifiedName())
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
	Name     sql.NullString `db:"name"`
	Type     sql.NullString `db:"type"`
	Kind     sql.NullString `db:"kind"`
	Nullable sql.NullString `db:"null?"`
}

func (td *tableDescription) IsNullable() bool {
	if td.Nullable.String == "Y" {
		return true
	} else {
		return false
	}
}

type primaryKeyDescription struct {
	ColumnName     sql.NullString `db:"column_name"`
	KeySequence    sql.NullString `db:"key_sequence"`
	ConstraintName sql.NullString `db:"constraint_name"`
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

func ScanPrimaryKeyDescription(rows *sqlx.Rows) ([]primaryKeyDescription, error) {
	pkds := []primaryKeyDescription{}
	for rows.Next() {
		pk := primaryKeyDescription{}
		err := rows.StructScan(&pk)
		if err != nil {
			return nil, err
		}
		pkds = append(pkds, pk)
	}
	return pkds, rows.Err()
}
