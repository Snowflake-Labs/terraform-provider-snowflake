package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SchemaBuilder abstracts the creation of SQL queries for a Snowflake schema
type SchemaBuilder struct {
	name                 string
	db                   string
	comment              string
	managedAccess        bool
	transient            bool
	setDataRetentionDays bool
	dataRetentionDays    int
}

// QualifiedName prepends the db if set and escapes everything nicely
func (sb *SchemaBuilder) QualifiedName() string {
	var n strings.Builder

	if sb.db != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, sb.db))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, sb.name))

	return n.String()
}

// Managed adds the WITH MANAGED ACCESS flag to the SchemaBuilder
func (sb *SchemaBuilder) Managed() *SchemaBuilder {
	sb.managedAccess = true
	return sb
}

// Transient adds the TRANSIENT flag to the SchemaBuilder
func (sb *SchemaBuilder) Transient() *SchemaBuilder {
	sb.transient = true
	return sb
}

// WithComment adds a comment to the SchemaBuilder
func (sb *SchemaBuilder) WithComment(c string) *SchemaBuilder {
	sb.comment = c
	return sb
}

// WithDataRetentionDays adds the days to retain data to the SchemaBuilder (must
// be 0-1 for standard edition, 0-90 for enterprise edition)
func (sb *SchemaBuilder) WithDataRetentionDays(d int) *SchemaBuilder {
	sb.setDataRetentionDays = true
	sb.dataRetentionDays = d
	return sb
}

// WithDB adds the name of the database to the SchemaBuilder
func (sb *SchemaBuilder) WithDB(db string) *SchemaBuilder {
	sb.db = db
	return sb
}

// Schema returns a pointer to a Builder that abstracts the DDL operations for a schema.
//
// Supported DDL operations are:
//   - CREATE SCHEMA
//   - ALTER SCHEMA
//   - DROP SCHEMA
//   - UNDROP SCHEMA
//   - USE SCHEMA
//   - SHOW SCHEMAS
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-database.html#schema-management)
func Schema(name string) *SchemaBuilder {
	return &SchemaBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new schema.
func (sb *SchemaBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	if sb.transient {
		q.WriteString(` TRANSIENT`)
	}

	q.WriteString(fmt.Sprintf(` SCHEMA %v`, sb.QualifiedName()))

	if sb.managedAccess {
		q.WriteString(` WITH MANAGED ACCESS`)
	}

	if sb.setDataRetentionDays {
		q.WriteString(fmt.Sprintf(` DATA_RETENTION_TIME_IN_DAYS = %d`, sb.dataRetentionDays))
	}

	if sb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, sb.comment))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the schema.
func (sb *SchemaBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER SCHEMA %v RENAME TO "%v"`, sb.QualifiedName(), newName)
}

// Swap returns the SQL query that Swaps all objects (tables, views, etc.) and
// metadata, including identifiers, between the two specified schemas.
func (sb *SchemaBuilder) Swap(targetSchema string) string {
	return fmt.Sprintf(`ALTER SCHEMA %v SWAP WITH "%v"`, sb.QualifiedName(), targetSchema)
}

// ChangeComment returns the SQL query that will update the comment on the schema.
func (sb *SchemaBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER SCHEMA %v SET COMMENT = '%v'`, sb.QualifiedName(), c)
}

// RemoveComment returns the SQL query that will remove the comment on the schema.
func (sb *SchemaBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER SCHEMA %v UNSET COMMENT`, sb.QualifiedName())
}

// ChangeDataRetentionDays returns the SQL query that will update the data retention days on the schema.
func (sb *SchemaBuilder) ChangeDataRetentionDays(d int) string {
	return fmt.Sprintf(`ALTER SCHEMA %v SET DATA_RETENTION_TIME_IN_DAYS = %d`, sb.QualifiedName(), d)
}

// RemoveDataRetentionDays returns the SQL query that will remove the data retention days on the schema.
func (sb *SchemaBuilder) RemoveDataRetentionDays() string {
	return fmt.Sprintf(`ALTER SCHEMA %v UNSET DATA_RETENTION_TIME_IN_DAYS`, sb.QualifiedName())
}

// Manage returns the SQL query that will enable managed access for a schema.
func (sb *SchemaBuilder) Manage() string {
	return fmt.Sprintf(`ALTER SCHEMA %v ENABLE MANAGED ACCESS`, sb.QualifiedName())
}

// Unmanage returns the SQL query that will disble managed access for a schema.
func (sb *SchemaBuilder) Unmanage() string {
	return fmt.Sprintf(`ALTER SCHEMA %v DISABLE MANAGED ACCESS`, sb.QualifiedName())
}

// Drop returns the SQL query that will drop a schema.
func (sb *SchemaBuilder) Drop() string {
	return fmt.Sprintf(`DROP SCHEMA %v`, sb.QualifiedName())
}

// Undrop returns the SQL query that will undrop a schema.
func (sb *SchemaBuilder) Undrop() string {
	return fmt.Sprintf(`UNDROP SCHEMA %v`, sb.QualifiedName())
}

// Use returns the SQL query that will use a schema.
func (sb *SchemaBuilder) Use() string {
	return fmt.Sprintf(`USE SCHEMA %v`, sb.QualifiedName())
}

// Show returns the SQL query that will show a schema.
func (sb *SchemaBuilder) Show() string {
	q := strings.Builder{}

	q.WriteString(fmt.Sprintf(`SHOW SCHEMAS LIKE '%v'`, sb.name))

	if sb.db != "" {
		q.WriteString(fmt.Sprintf(` IN DATABASE "%v"`, sb.db))
	}

	return q.String()
}

type schema struct {
	Name          sql.NullString `db:"name"`
	DatabaseName  sql.NullString `db:"database_name"`
	Comment       sql.NullString `db:"comment"`
	Options       sql.NullString `db:"options"`
	RetentionTime sql.NullInt64  `db:"retention_time"`
}

func ScanSchema(row *sqlx.Row) (*schema, error) {
	r := &schema{}
	err := row.StructScan(r)
	return r, err
}
