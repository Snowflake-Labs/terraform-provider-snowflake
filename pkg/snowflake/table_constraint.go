package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

// TableConstraintBuilder abstracts the creation of SQL queries for a Snowflake table constraint.
type TableConstraintBuilder struct {
	name             string
	columns          []string
	constraintType   string
	tableID          string
	enforced         bool
	deferrable       bool
	initially        string
	enable           bool
	validate         bool
	rely             bool
	referenceTableID string
	referenceColumns []string
	match            string
	update           string
	delete           string
	comment          string
}

func NewTableConstraintBuilder(name string, constraintType string, tableID string) *TableConstraintBuilder {
	return &TableConstraintBuilder{
		name:           name,
		constraintType: constraintType,
		tableID:        tableID,
	}
}

// WithComment sets comment.
func (b *TableConstraintBuilder) WithComment(comment string) *TableConstraintBuilder {
	b.comment = comment
	return b
}

// WithColumns sets columns.
func (b *TableConstraintBuilder) WithColumns(columns []string) *TableConstraintBuilder {
	b.columns = columns
	return b
}

// WithEnforced sets enforced.
func (b *TableConstraintBuilder) WithEnforced(enforced bool) *TableConstraintBuilder {
	b.enforced = enforced
	return b
}

// WithDeferrable sets deferrable.
func (b *TableConstraintBuilder) WithDeferrable(deferrable bool) *TableConstraintBuilder {
	b.deferrable = deferrable
	return b
}

// WithInitially sets initially.
func (b *TableConstraintBuilder) WithInitially(initially string) *TableConstraintBuilder {
	b.initially = initially
	return b
}

// WithEnable sets enable.
func (b *TableConstraintBuilder) WithEnable(enable bool) *TableConstraintBuilder {
	b.enable = enable
	return b
}

// WithValidated sets validated.
func (b *TableConstraintBuilder) WithValidate(validate bool) *TableConstraintBuilder {
	b.validate = validate
	return b
}

// WithRely sets rely.
func (b *TableConstraintBuilder) WithRely(rely bool) *TableConstraintBuilder {
	b.rely = rely
	return b
}

// WithReferenceTableID sets referenceTableID.
func (b *TableConstraintBuilder) WithReferenceTableID(referenceTableID string) *TableConstraintBuilder {
	b.referenceTableID = referenceTableID
	return b
}

// WithReferenceColumns sets referenceColumns.
func (b *TableConstraintBuilder) WithReferenceColumns(referenceColumns []string) *TableConstraintBuilder {
	b.referenceColumns = referenceColumns
	return b
}

// WithMatch sets match.
func (b *TableConstraintBuilder) WithMatch(match string) *TableConstraintBuilder {
	b.match = match
	return b
}

// WithUpdate sets update.
func (b *TableConstraintBuilder) WithUpdate(onUpdate string) *TableConstraintBuilder {
	b.update = onUpdate
	return b
}

// WithDelete sets delete.
func (b *TableConstraintBuilder) WithDelete(onDelete string) *TableConstraintBuilder {
	b.delete = onDelete
	return b
}

func (b *TableConstraintBuilder) formattedReferenceColumns() []string {
	formattedColumns := make([]string, len(b.referenceColumns))
	for i, c := range b.referenceColumns {
		formattedColumns[i] = fmt.Sprintf(`"%v"`, EscapeString(c))
	}
	return formattedColumns
}

func (b *TableConstraintBuilder) formattedColumns() []string {
	formattedColumns := make([]string, len(b.columns))
	for i, c := range b.columns {
		formattedColumns[i] = fmt.Sprintf(`"%v"`, EscapeString(c))
	}
	return formattedColumns
}

// Rename returns the SQL query that will rename the table constraint.
func (b *TableConstraintBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER TABLE %v RENAME CONSTRAINT %v TO %v`, b.tableID, b.name, newName)
}

// SetComment returns the SQL query that will update the comment on the table constraint.
func (b *TableConstraintBuilder) SetComment(c string) string {
	return fmt.Sprintf(`COMMENT ON CONSTRAINT %v IS '%v'`, b.name, EscapeString(c))
}

// Drop returns the SQL query that will drop a table constraint.
func (b *TableConstraintBuilder) Drop() string {
	s := fmt.Sprintf(`ALTER TABLE %v DROP CONSTRAINT %v`, b.tableID, b.name)
	/*if b.columns != nil {
		s +=  fmt.Sprintf(` (%v)`, strings.Join(b.formattedColumns(), ", "))
	}*/
	s += " CASCADE"
	return s
}

type TableConstraint struct {
	ConstraintCatalog sql.NullString `db:"CONSTRAINT_CATALOG"`
	ConstraintSchema  sql.NullString `db:"CONSTRAINT_SCHEMA"`
	ConstraintName    sql.NullString `db:"CONSTRAINT_NAME"`
	TableCatalog      sql.NullString `db:"TABLE_CATALOG"`
	TableSchema       sql.NullString `db:"TABLE_SCHEMA"`
	TableName         sql.NullString `db:"TABLE_NAME"`
	ConstraintType    sql.NullString `db:"CONSTRAINT_TYPE"`
	IsDeferrable      sql.NullString `db:"IS_DEFERRABLE"`
	InitiallyDeferred sql.NullString `db:"INITIALLY_DEFERRED"`
	Enforced          sql.NullString `db:"ENFORCED"`
	Comment           sql.NullString `db:"COMMENT"`
}

// Show returns the SQL query that will show a table constraint by ID.
func ShowTableConstraint(name, tableDB, tableSchema, tableName string, db *sql.DB) (*TableConstraint, error) {
	stmt := `SELECT * FROM SNOWFLAKE.INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE TABLE_NAME = '?' AND TABLE_SCHEMA = '?' AND TABLE_CATALOG = '?' AND CONSTRAINT_NAME = '?'`
	rows, err := db.Query(stmt,
		tableName, tableSchema, tableDB, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tableConstraints := []TableConstraint{}
	log.Printf("[DEBUG] tableConstraints is %v", tableConstraints)
	if err := sqlx.StructScan(rows, &tableConstraints); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] no tableConstraints found for constraint %s", name)
			return nil, err
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return &tableConstraints[0], nil
}
