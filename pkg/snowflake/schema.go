package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SchemaBuilder abstracts the creation of SQL queries for a Snowflake schema.
type SchemaBuilder struct {
	name                 string
	db                   string
	comment              string
	managedAccess        bool
	transient            bool
	setDataRetentionDays bool
	dataRetentionDays    int
	tags                 []TagValue
}

// QualifiedName prepends the db if set and escapes everything nicely.
func (sb *SchemaBuilder) QualifiedName() string {
	var n strings.Builder

	if sb.db != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, sb.db))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, sb.name))

	return n.String()
}

// Managed adds the WITH MANAGED ACCESS flag to the SchemaBuilder.
func (sb *SchemaBuilder) Managed() *SchemaBuilder {
	sb.managedAccess = true
	return sb
}

// Transient adds the TRANSIENT flag to the SchemaBuilder.
func (sb *SchemaBuilder) Transient() *SchemaBuilder {
	sb.transient = true
	return sb
}

// WithComment adds a comment to the SchemaBuilder.
func (sb *SchemaBuilder) WithComment(c string) *SchemaBuilder {
	sb.comment = c
	return sb
}

// WithDataRetentionDays adds the days to retain data to the SchemaBuilder (must
// be 0-1 for standard edition, 0-90 for enterprise edition).
func (sb *SchemaBuilder) WithDataRetentionDays(d int) *SchemaBuilder {
	sb.setDataRetentionDays = true
	sb.dataRetentionDays = d
	return sb
}

// WithDB adds the name of the database to the SchemaBuilder.
func (sb *SchemaBuilder) WithDB(db string) *SchemaBuilder {
	sb.db = db
	return sb
}

// WithTags sets the tags on the SchemaBuilder.
func (sb *SchemaBuilder) WithTags(tags []TagValue) *SchemaBuilder {
	sb.tags = tags
	return sb
}

// AddTag returns the SQL query that will add a new tag to the schema.
func (sb *SchemaBuilder) AddTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER SCHEMA %s SET TAG "%v"."%v"."%v" = "%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// ChangeTag returns the SQL query that will alter a tag on the schema.
func (sb *SchemaBuilder) ChangeTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER SCHEMA %s SET TAG "%v"."%v"."%v" = "%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name, tag.Value)
}

// UnsetTag returns the SQL query that will unset a tag on the schema.
func (sb *SchemaBuilder) UnsetTag(tag TagValue) string {
	return fmt.Sprintf(`ALTER SCHEMA %s UNSET TAG "%v"."%v"."%v"`, sb.QualifiedName(), tag.Database, tag.Schema, tag.Name)
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
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(sb.comment)))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the schema.
func (sb *SchemaBuilder) Rename(newName string) string {
	oldName := sb.QualifiedName()
	sb.name = newName
	return fmt.Sprintf(`ALTER SCHEMA %v RENAME TO %v`, oldName, sb.QualifiedName())
}

// Swap returns the SQL query that Swaps all objects (tables, views, etc.) and
// metadata, including identifiers, between the two specified schemas.
func (sb *SchemaBuilder) Swap(targetSchema string) string {
	sourceSchema := sb.QualifiedName()
	sb.name = targetSchema
	return fmt.Sprintf(`ALTER SCHEMA %v SWAP WITH %v`, sourceSchema, sb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the schema.
func (sb *SchemaBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER SCHEMA %v SET COMMENT = '%v'`, sb.QualifiedName(), EscapeString(c))
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
	RetentionTime sql.NullString `db:"retention_time"`
}

func ScanSchema(row *sqlx.Row) (*schema, error) {
	r := &schema{}
	err := row.StructScan(r)
	return r, err
}

func ListSchemas(databaseName string, db *sql.DB) ([]schema, error) {
	stmt := fmt.Sprintf(`SHOW SCHEMAS IN DATABASE "%v"`, databaseName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []schema{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Println("[DEBUG] no schemas found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
