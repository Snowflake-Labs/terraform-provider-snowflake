package snowflake

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/jmoiron/sqlx"
)

// PrimaryKey structure that represents a tables primary key.
type PrimaryKey struct {
	name string
	keys []string
}

// WithName set the primary key name.
func (pk *PrimaryKey) WithName(name string) *PrimaryKey {
	pk.name = name
	return pk
}

// WithKeys set the primary key keys.
func (pk *PrimaryKey) WithKeys(keys []string) *PrimaryKey {
	pk.keys = keys
	return pk
}

type ColumnDefaultType int

const (
	columnDefaultTypeConstant = iota
	columnDefaultTypeSequence
	columnDefaultTypeExpression
)

type ColumnDefault struct {
	_type      ColumnDefaultType
	expression string
}

type ColumnIdentity struct {
	startNum int
	stepNum  int
}

func (id *ColumnIdentity) WithStartNum(start int) *ColumnIdentity {
	id.startNum = start
	return id
}

func (id *ColumnIdentity) WithStep(step int) *ColumnIdentity {
	id.stepNum = step
	return id
}

func NewColumnDefaultWithConstant(constant string) *ColumnDefault {
	return &ColumnDefault{
		_type:      columnDefaultTypeConstant,
		expression: constant,
	}
}

func NewColumnDefaultWithExpression(expression string) *ColumnDefault {
	return &ColumnDefault{
		_type:      columnDefaultTypeExpression,
		expression: expression,
	}
}

func NewColumnDefaultWithSequence(sequence string) *ColumnDefault {
	return &ColumnDefault{
		_type:      columnDefaultTypeSequence,
		expression: sequence,
	}
}

func (d *ColumnDefault) String(columnType string) string {
	columnType = strings.ToUpper(columnType)

	switch {
	case d._type == columnDefaultTypeExpression:
		return d.expression

	case d._type == columnDefaultTypeSequence:
		return fmt.Sprintf(`%v.NEXTVAL`, d.expression)

	case d._type == columnDefaultTypeConstant && (strings.Contains(columnType, "CHAR") || columnType == "STRING" || columnType == "TEXT"):
		return EscapeSnowflakeString(d.expression)

	default:
		return d.expression
	}
}

func (d *ColumnDefault) UnescapeConstantSnowflakeString(columnType string) string {
	columnType = strings.ToUpper(columnType)

	if d._type == columnDefaultTypeConstant && (strings.Contains(columnType, "CHAR") || columnType == "STRING" || columnType == "TEXT") {
		return UnescapeSnowflakeString(d.expression)
	}

	return d.expression
}

// Column structure that represents a table column.
type Column struct {
	name          string
	_type         string // type is reserved
	nullable      bool
	_default      *ColumnDefault // default is reserved
	identity      *ColumnIdentity
	comment       string // pointer as value is nullable
	maskingPolicy string
}

// WithName set the column name.
func (c *Column) WithName(name string) *Column {
	c.name = name
	return c
}

// WithType set the column type.
func (c *Column) WithType(t string) *Column {
	c._type = t
	return c
}

// WithNullable set if the column is nullable.
func (c *Column) WithNullable(nullable bool) *Column {
	c.nullable = nullable
	return c
}

func (c *Column) WithDefault(cd *ColumnDefault) *Column {
	c._default = cd
	return c
}

// WithComment set the column comment.
func (c *Column) WithComment(comment string) *Column {
	c.comment = comment
	return c
}

func (c *Column) WithMaskingPolicy(maskingPolicy string) *Column {
	c.maskingPolicy = maskingPolicy
	return c
}

func (c *Column) WithIdentity(id *ColumnIdentity) *Column {
	c.identity = id
	return c
}

func (c *Column) getColumnDefinition(withInlineConstraints bool, withComment bool) string {
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

	if c._default != nil {
		colDef.WriteString(fmt.Sprintf(` DEFAULT %v`, c._default.String(c._type)))
	}

	if c.identity != nil {
		colDef.WriteString(fmt.Sprintf(` IDENTITY(%v, %v)`, c.identity.startNum, c.identity.stepNum))
	}

	if strings.TrimSpace(c.maskingPolicy) != "" {
		colDef.WriteString(fmt.Sprintf(` WITH MASKING POLICY %v`, EscapeString(c.maskingPolicy)))
	}

	if withComment {
		colDef.WriteString(fmt.Sprintf(` COMMENT '%v'`, EscapeString(c.comment)))
	}

	return colDef.String()
}

func FlattenTablePrimaryKey(pkds []PrimaryKeyDescription) []interface{} {
	flattened := []interface{}{}
	if len(pkds) == 0 {
		return flattened
	}

	sort.SliceStable(pkds, func(i, j int) bool {
		num1, _ := strconv.Atoi(pkds[i].KeySequence.String)
		num2, _ := strconv.Atoi(pkds[j].KeySequence.String)
		return num1 < num2
	})
	// sort our keys on the key sequence

	flat := map[string]interface{}{}
	keys := make([]string, 0, len(pkds))
	var name string
	var nameSet bool

	for _, pk := range pkds {
		// set as empty string, sys_constraint means it was an unnnamed constraint
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

func (c Columns) getColumnDefinitions(withInlineConstraints bool, withComments bool) string {
	// TODO(el): verify Snowflake reflects column order back in desc table calls
	columnDefinitions := []string{}
	for _, column := range c {
		columnDefinitions = append(columnDefinitions, column.getColumnDefinition(withInlineConstraints, withComments))
	}

	// NOTE: intentionally blank leading space
	return fmt.Sprintf(" (%s)", strings.Join(columnDefinitions, ", "))
}

// TableBuilder abstracts the creation of SQL queries for a Snowflake schema.
type TableBuilder struct {
	name                    string
	db                      string
	schema                  string
	columns                 Columns
	comment                 string
	clusterBy               []string
	primaryKey              PrimaryKey
	dataRetentionTimeInDays *int
	changeTracking          bool
	tags                    []TagValue
}

// QualifiedName prepends the db and schema if set and escapes everything nicely.
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

// WithComment adds a comment to the TableBuilder.
func (tb *TableBuilder) WithComment(c string) *TableBuilder {
	tb.comment = c
	return tb
}

// WithColumns sets the column definitions on the TableBuilder.
func (tb *TableBuilder) WithColumns(c Columns) *TableBuilder {
	tb.columns = c
	return tb
}

// WithClustering adds cluster keys/expressions to TableBuilder.
func (tb *TableBuilder) WithClustering(c []string) *TableBuilder {
	tb.clusterBy = c
	return tb
}

// WithPrimaryKey sets the primary key on the TableBuilder.
func (tb *TableBuilder) WithPrimaryKey(pk PrimaryKey) *TableBuilder {
	tb.primaryKey = pk
	return tb
}

// WithDataRetentionTimeInDays sets the data retention time on the TableBuilder.
func (tb *TableBuilder) WithDataRetentionTimeInDays(days int) *TableBuilder {
	tb.dataRetentionTimeInDays = sdk.Int(days)
	return tb
}

// WithChangeTracking sets the change tracking on the TableBuilder.
func (tb *TableBuilder) WithChangeTracking(changeTracking bool) *TableBuilder {
	tb.changeTracking = changeTracking
	return tb
}

// WithTags sets the tags on the TableBuilder.
func (tb *TableBuilder) WithTags(tags []TagValue) *TableBuilder {
	tb.tags = tags
	return tb
}

// AddTag returns the SQL query that will add a new tag to the table.
func (tb *TableBuilder) AddTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER TABLE %s SET TAG "%v"."%v"."%v" = "%v"`, tb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// ChangeTag returns the SQL query that will alter a tag on the table.
func (tb *TableBuilder) ChangeTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER TABLE %s SET TAG "%v"."%v"."%v" = "%v"`, tb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// UnsetTag returns the SQL query that will unset a tag on the table.
func (tb *TableBuilder) UnsetTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER TABLE %s UNSET TAG "%v"."%v"."%v"`, tb.QualifiedName(), tag.Database, tag.Schema, tag.Name)
}

// Function to get clustering definition.
func (tb *TableBuilder) GetClusterKeyString() string {
	return JoinStringList(tb.clusterBy, ", ")
}

func (tb *TableBuilder) GetTagValueString() string {
	var q strings.Builder
	for _, v := range tb.tags {
		fmt.Println(v)
		if v.Schema != "" {
			if v.Database != "" {
				q.WriteString(fmt.Sprintf(`"%v".`, v.Database))
			}
			q.WriteString(fmt.Sprintf(`"%v".`, v.Schema))
		}
		q.WriteString(fmt.Sprintf(`"%v" = "%v", `, v.Name, v.Value))
	}
	return strings.TrimSuffix(q.String(), ", ")
}

func JoinStringList(instrings []string, delimiter string) string {
	return fmt.Sprint(strings.Join(instrings, delimiter))
}

func QuoteStringList(instrings []string) []string {
	clean := make([]string, 0, len(instrings))
	for _, word := range instrings {
		quoted := fmt.Sprintf(`"%s"`, word)
		clean = append(clean, quoted)
	}
	return clean
}

func (tb *TableBuilder) getCreateStatementBody() string {
	var q strings.Builder

	colDef := tb.columns.getColumnDefinitions(true, true)

	if len(tb.primaryKey.keys) > 0 {
		colDef = strings.TrimSuffix(colDef, ")") // strip trailing
		q.WriteString(colDef)
		if tb.primaryKey.name != "" {
			q.WriteString(fmt.Sprintf(` ,CONSTRAINT "%v" PRIMARY KEY(%v)`, tb.primaryKey.name, JoinStringList(QuoteStringList(tb.primaryKey.keys), ",")))
		} else {
			q.WriteString(fmt.Sprintf(` ,PRIMARY KEY(%v)`, JoinStringList(QuoteStringList(tb.primaryKey.keys), ",")))
		}

		q.WriteString(")") // add closing
	} else {
		q.WriteString(colDef)
	}

	return q.String()
}

// Table returns a pointer to a Builder that abstracts the DDL operations for a table.
//
// Supported DDL operations are:
//   - ALTER TABLE
//   - DROP TABLE
//   - SHOW TABLES
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-table.html)
func NewTableBuilder(name, db, schema string) *TableBuilder {
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
func NewTableWithColumnDefinitionsBuilder(name, db, schema string, columns Columns) *TableBuilder {
	return &TableBuilder{
		name:    name,
		db:      db,
		schema:  schema,
		columns: columns,
	}
}

// Create returns the SQL statement required to create a table.
func (tb *TableBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE TABLE %v`, tb.QualifiedName()))
	q.WriteString(tb.getCreateStatementBody())

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	if tb.clusterBy != nil {
		// add optional clustering statement
		q.WriteString(fmt.Sprintf(` CLUSTER BY LINEAR(%v)`, tb.GetClusterKeyString()))
	}

	if tb.dataRetentionTimeInDays != nil {
		q.WriteString(fmt.Sprintf(` DATA_RETENTION_TIME_IN_DAYS = %d`, *tb.dataRetentionTimeInDays))
	}
	q.WriteString(fmt.Sprintf(` CHANGE_TRACKING = %t`, tb.changeTracking))

	if tb.tags != nil {
		q.WriteString(fmt.Sprintf(` WITH TAG (%v)`, tb.GetTagValueString()))
	}

	return q.String()
}

// ChangeClusterBy returns the SQL query to change cluastering on table.
func (tb *TableBuilder) ChangeClusterBy(cb string) string {
	return fmt.Sprintf(`ALTER TABLE %v CLUSTER BY LINEAR(%v)`, tb.QualifiedName(), cb)
}

// ChangeComment returns the SQL query that will update the comment on the table.
func (tb *TableBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TABLE %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(c))
}

// ChangeDataRetention returns the SQL query that will update the DATA_RETENTION_TIME_IN_DAYS on the table.
func (tb *TableBuilder) ChangeDataRetention(days int) string {
	return fmt.Sprintf(`ALTER TABLE %v SET DATA_RETENTION_TIME_IN_DAYS = %d`, tb.QualifiedName(), days)
}

// ChangeChangeTracking returns the SQL query that will update the CHANGE_TRACKING on the table.
func (tb *TableBuilder) ChangeChangeTracking(changeTracking bool) string {
	return fmt.Sprintf(`ALTER TABLE %v SET CHANGE_TRACKING = %t`, tb.QualifiedName(), changeTracking)
}

// AddColumn returns the SQL query that will add a new column to the table.
func (tb *TableBuilder) AddColumn(name string, dataType string, nullable bool, _default *ColumnDefault, identity *ColumnIdentity, comment string, maskingPolicy string) string {
	col := Column{
		name:          name,
		_type:         dataType,
		nullable:      nullable,
		_default:      _default,
		identity:      identity,
		comment:       comment,
		maskingPolicy: maskingPolicy,
	}
	return fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition(true, true))
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

	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN %s`, tb.QualifiedName(), col.getColumnDefinition(false, false))
}

func (tb *TableBuilder) ChangeColumnComment(name string, comment string) string {
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%v" COMMENT '%v'`, tb.QualifiedName(), EscapeString(name), EscapeString(comment))
}

func (tb *TableBuilder) ChangeColumnMaskingPolicy(name string, maskingPolicy string) string {
	if strings.TrimSpace(maskingPolicy) == "" {
		return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%v" UNSET MASKING POLICY`, tb.QualifiedName(), EscapeString(name))
	}
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%v" SET MASKING POLICY %v`, tb.QualifiedName(), EscapeString(name), EscapeString(maskingPolicy))
}

func (tb *TableBuilder) DropColumnDefault(name string) string {
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%v" DROP DEFAULT`, tb.QualifiedName(), EscapeString(name))
}

// RemoveComment returns the SQL query that will remove the comment on the table.
func (tb *TableBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TABLE %v UNSET COMMENT`, tb.QualifiedName())
}

// Return sql to set/unset null constraint on column.
func (tb *TableBuilder) ChangeNullConstraint(name string, nullable bool) string {
	if nullable {
		return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%s" DROP NOT NULL`, tb.QualifiedName(), name)
	}
	return fmt.Sprintf(`ALTER TABLE %s MODIFY COLUMN "%s" SET NOT NULL`, tb.QualifiedName(), name)
}

func (tb *TableBuilder) ChangePrimaryKey(newPk PrimaryKey) string {
	tb.WithPrimaryKey(newPk)
	pks := JoinStringList(QuoteStringList(newPk.keys), ", ")
	if tb.primaryKey.name != "" {
		return fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT "%v" PRIMARY KEY(%v)`, tb.QualifiedName(), tb.primaryKey.name, pks)
	}
	return fmt.Sprintf(`ALTER TABLE %s ADD PRIMARY KEY(%v)`, tb.QualifiedName(), pks)
}

func (tb *TableBuilder) DropPrimaryKey() string {
	return fmt.Sprintf(`ALTER TABLE %s DROP PRIMARY KEY`, tb.QualifiedName())
}

// RemoveClustering returns the SQL query that will remove data clustering from the table.
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

func (tb *TableBuilder) Rename(newName string) string {
	oldName := tb.QualifiedName()
	tb.name = newName
	return fmt.Sprintf(`ALTER TABLE %s RENAME TO %s`, oldName, tb.QualifiedName())
}

type Table struct {
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
	RetentionTime       sql.NullInt32  `db:"retention_time"`
	AutomaticClustering sql.NullString `db:"automatic_clustering"`
	ChangeTracking      sql.NullString `db:"change_tracking"`
	IsExternal          sql.NullString `db:"is_external"`
}

func ScanTable(row *sqlx.Row) (*Table, error) {
	t := &Table{}
	e := row.StructScan(t)
	return t, e
}

type TableDescription struct {
	Name          sql.NullString `db:"name"`
	Type          sql.NullString `db:"type"`
	Kind          sql.NullString `db:"kind"`
	Nullable      sql.NullString `db:"null?"`
	Default       sql.NullString `db:"default"`
	Comment       sql.NullString `db:"comment"`
	MaskingPolicy sql.NullString `db:"policy name"`
}

type PrimaryKeyDescription struct {
	ColumnName     sql.NullString `db:"column_name"`
	KeySequence    sql.NullString `db:"key_sequence"`
	ConstraintName sql.NullString `db:"constraint_name"`
}
