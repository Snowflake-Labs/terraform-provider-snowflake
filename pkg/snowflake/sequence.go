package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Sequence returns a pointer to a Builder for a sequence
func Sequence(name, db, schema string) *SequenceBuilder {
	return &SequenceBuilder{
		name:      name,
		db:        db,
		schema:    schema,
		increment: 1,
		start:     1,
	}
}

type sequence struct {
	Name       sql.NullString `db:"name"`
	DBName     sql.NullString `db:"database_name"`
	SchemaName sql.NullString `db:"schema_name"`
	NextValue  sql.NullString `db:"next_value"`
	Increment  sql.NullString `db:"interval"`
	CreatedOn  sql.NullString `db:"created_on"`
	Owner      sql.NullString `db:"owner"`
	Comment    sql.NullString `db:"comment"`
}

type SequenceBuilder struct {
	name      string
	db        string
	schema    string
	increment int
	comment   string
	start     int
}

// Drop returns the SQL query that will drop a sequence.
func (sb *SequenceBuilder) Drop() string {
	return fmt.Sprintf(`DROP SEQUENCE %v`, sb.QualifiedName())
}

// Show returns the SQL query that will show a sequence.
func (sb *SequenceBuilder) Show() string {
	return fmt.Sprintf(`SHOW SEQUENCES LIKE '%v' IN SCHEMA "%v"."%v"`, sb.name, sb.db, sb.schema)
}

func (sb *SequenceBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE SEQUENCE %v`, sb.QualifiedName()))
	if sb.start != 1 {
		q.WriteString(fmt.Sprintf(` START = %d`, sb.start))
	}
	if sb.increment != 1 {
		q.WriteString(fmt.Sprintf(` INCREMENT = %d`, sb.increment))
	}
	if sb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(sb.comment)))
	}
	return q.String()
}

func (sb *SequenceBuilder) WithComment(comment string) *SequenceBuilder {
	sb.comment = comment
	return sb
}

func (sb *SequenceBuilder) WithIncrement(increment int) *SequenceBuilder {
	sb.increment = increment
	return sb
}

func (sb *SequenceBuilder) WithStart(start int) *SequenceBuilder {
	sb.start = start
	return sb
}

func (sb *SequenceBuilder) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"`, sb.db, sb.schema, sb.name)
}

func (sb *SequenceBuilder) Address() string {
	return AddressEscape(sb.db, sb.schema, sb.name)
}

func ScanSequence(row *sqlx.Row) (*sequence, error) {
	d := &sequence{}
	e := row.StructScan(d)
	return d, e
}

func ListSequences(databaseName string, schemaName string, db *sql.DB) ([]sequence, error) {
	stmt := fmt.Sprintf(`SHOW SEQUENCES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []sequence{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no sequences found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
