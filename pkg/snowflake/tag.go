package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TagBuilder abstracts the creation of SQL queries for a Snowflake tag
type TagBuilder struct {
	name    string
	db      string
	schema  string
	comment string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (tb *TagBuilder) QualifiedName() string {
	var n strings.Builder

	if tb.db != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, tb.db))
	}

	if tb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, tb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, tb.name))

	return n.String()
}

// WithComment adds a comment to the TagBuilder
func (tb *TagBuilder) WithComment(c string) *TagBuilder {
	tb.comment = c
	return tb
}

// WithDB adds the name of the database to the TagBuilder
func (tb *TagBuilder) WithDB(db string) *TagBuilder {
	tb.db = db
	return tb
}

// WithSchema adds the name of the schema to the TagBuilder
func (tb *TagBuilder) WithSchema(schema string) *TagBuilder {
	tb.schema = schema
	return tb
}

// Tag returns a pointer to a Builder that abstracts the DDL operations for a tag.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - ALTER TAG
//   - DROP TAG
//   - UNDROP TAG
//   - SHOW TAGS
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func Tag(name string) *TagBuilder {
	return &TagBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new tag.
func (tb *TagBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TAG %v`, tb.QualifiedName()))

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the tag.
func (tb *TagBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER TAG %v RENAME TO "%v"`, tb.QualifiedName(), newName)
}

// ChangeComment returns the SQL query that will update the comment on the tag.
func (tb *TagBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TAG %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the tag.
func (tb *TagBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TAG %v UNSET COMMENT`, tb.QualifiedName())
}

// Drop returns the SQL query that will drop a tag.
func (tb *TagBuilder) Drop() string {
	return fmt.Sprintf(`DROP TAG %v`, tb.QualifiedName())
}

// Undrop returns the SQL query that will undrop a tag.
func (tb *TagBuilder) Undrop() string {
	return fmt.Sprintf(`UNDROP TAG %v`, tb.QualifiedName())
}

// Show returns the SQL query that will show a tag.
func (tb *TagBuilder) Show() string {
	q := strings.Builder{}

	q.WriteString(fmt.Sprintf(`SHOW TAGS LIKE '%v'`, tb.name))

	if tb.schema != "" && tb.db != "" {
		q.WriteString(fmt.Sprintf(` IN SCHEMA "%v"."%v"`, tb.db, tb.schema))
	} else if tb.db != "" {
		q.WriteString(fmt.Sprintf(` IN DATABASE "%v"`, tb.db))
	}

	return q.String()
}

type tag struct {
	Name         sql.NullString `db:"name"`
	DatabaseName sql.NullString `db:"database_name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Comment      sql.NullString `db:"comment"`
}

type TagValue struct {
	Name     string
	Database string
	Schema   string
	Value    string
}

func ScanTag(row *sqlx.Row) (*tag, error) {
	r := &tag{}
	err := row.StructScan(r)
	return r, err
}

// ListTags returns a list of tags in a database or schema
func ListTags(databaseName, schemaName string, db *sql.DB) ([]tag, error) {
	stmt := fmt.Sprintf(`SHOW TAGS IN SCHEMA "%v"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []tag{}
	err = sqlx.StructScan(rows, &tags)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no tags found")
		return nil, nil
	}
	return tags, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
