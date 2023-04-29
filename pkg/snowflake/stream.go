package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// StreamBuilder abstracts the creation of SQL queries for a Snowflake stream.
type StreamBuilder struct {
	name            string
	db              string
	schema          string
	externalTable   bool
	onTable         string
	onView          string
	onStage         string
	appendOnly      bool
	insertOnly      bool
	showInitialRows bool
	comment         string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely.
func (sb *StreamBuilder) QualifiedName() string {
	var n strings.Builder

	if sb.db != "" && sb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, sb.db, sb.schema))
	}

	if sb.db != "" && sb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, sb.db))
	}

	if sb.db == "" && sb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, sb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, sb.name))

	return n.String()
}

func (sb *StreamBuilder) WithComment(c string) *StreamBuilder {
	sb.comment = c
	return sb
}

func (sb *StreamBuilder) WithOnTable(d string, s string, t string) *StreamBuilder {
	sb.onTable = fmt.Sprintf(`"%v"."%v"."%v"`, d, s, t)
	return sb
}

func (sb *StreamBuilder) WithExternalTable(b bool) *StreamBuilder {
	sb.externalTable = b
	return sb
}

func (sb *StreamBuilder) WithOnView(d string, s string, t string) *StreamBuilder {
	sb.onView = fmt.Sprintf(`"%v"."%v"."%v"`, d, s, t)
	return sb
}

func (sb *StreamBuilder) WithOnStage(d string, s string, t string) *StreamBuilder {
	sb.onStage = fmt.Sprintf(`"%v"."%v"."%v"`, d, s, t)
	return sb
}

func (sb *StreamBuilder) WithAppendOnly(b bool) *StreamBuilder {
	sb.appendOnly = b
	return sb
}

func (sb *StreamBuilder) WithInsertOnly(b bool) *StreamBuilder {
	sb.insertOnly = b
	return sb
}

func (sb *StreamBuilder) WithShowInitialRows(b bool) *StreamBuilder {
	sb.showInitialRows = b
	return sb
}

// Stream returns a pointer to a Builder that abstracts the DDL operations for a stream.
//
// Supported DDL operations are:
//   - CREATE Stream
//   - ALTER Stream
//   - DROP Stream
//   - SHOW Stream
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/sql/create-stream.html)
func Stream(name, db, schema string) *StreamBuilder {
	return &StreamBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL statement required to create a stream.
func (sb *StreamBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE STREAM %v`, sb.QualifiedName()))

	q.WriteString(` ON`)

	switch {
	case sb.onTable != "":
		if sb.externalTable {
			q.WriteString(` EXTERNAL`)
		}
		q.WriteString(fmt.Sprintf(` TABLE %v`, sb.onTable))
	case sb.onView != "":
		q.WriteString(fmt.Sprintf(` VIEW %v`, sb.onView))
	case sb.onStage != "":
		q.WriteString(fmt.Sprintf(` STAGE %v`, sb.onStage))
	}

	if sb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(sb.comment)))
	}

	if sb.onStage == "" {
		q.WriteString(fmt.Sprintf(` APPEND_ONLY = %v`, sb.appendOnly))

		q.WriteString(fmt.Sprintf(` INSERT_ONLY = %v`, sb.insertOnly))

		q.WriteString(fmt.Sprintf(` SHOW_INITIAL_ROWS = %v`, sb.showInitialRows))
	}

	return q.String()
}

// ChangeComment returns the SQL query that will update the comment on the stream.
func (sb *StreamBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER STREAM %v SET COMMENT = '%v'`, sb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the stream.
func (sb *StreamBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER STREAM %v UNSET COMMENT`, sb.QualifiedName())
}

// Drop returns the SQL query that will drop a stream.
func (sb *StreamBuilder) Drop() string {
	return fmt.Sprintf(`DROP STREAM %v`, sb.QualifiedName())
}

// Show returns the SQL query that will show a stream.
func (sb *StreamBuilder) Show() string {
	return fmt.Sprintf(`SHOW STREAMS LIKE '%v' IN SCHEMA "%v"."%v"`, sb.name, sb.db, sb.schema)
}

type DescStreamRow struct {
	CreatedOn       sql.NullString `db:"created_on"`
	StreamName      sql.NullString `db:"name"`
	DatabaseName    sql.NullString `db:"database_name"`
	SchemaName      sql.NullString `db:"schema_name"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
	ShowInitialRows bool           `db:"show_initial_rows"`
	TableName       sql.NullString `db:"table_name"`
	Type            sql.NullString `db:"type"`
	Stale           sql.NullString `db:"stale"`
	Mode            sql.NullString `db:"mode"`
	SourceType      sql.NullString `db:"source_type"`
}

func ScanStream(row *sqlx.Row) (*DescStreamRow, error) {
	t := &DescStreamRow{}
	e := row.StructScan(t)
	return t, e
}

func ListStreams(databaseName string, schemaName string, db *sql.DB) ([]DescStreamRow, error) {
	stmt := fmt.Sprintf(`SHOW STREAMS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []DescStreamRow{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no stages found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
