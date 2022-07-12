package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// PipeBuilder abstracts the creation of SQL queries for a Snowflake schema
type PipeBuilder struct {
	name             string
	db               string
	schema           string
	autoIngest       bool
	awsSnsTopicArn   string
	comment          string
	copyStatement    string
	integration      string
	errorIntegration string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (pb *PipeBuilder) QualifiedName() string {
	var n strings.Builder

	if pb.db != "" && pb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, pb.db, pb.schema))
	}

	if pb.db != "" && pb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, pb.db))
	}

	if pb.db == "" && pb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, pb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, pb.name))

	return n.String()
}

// WithAutoIngest adds the auto_ingest flag to the PipeBuilder
func (pb *PipeBuilder) WithAutoIngest() *PipeBuilder {
	pb.autoIngest = true
	return pb
}

// WithAwsSnsTopicArn adds the aws_sns_topic to the PipeBuilder
func (pb *PipeBuilder) WithAwsSnsTopicArn(s string) *PipeBuilder {
	pb.awsSnsTopicArn = s
	return pb
}

// WithComment adds a comment to the PipeBuilder
func (pb *PipeBuilder) WithComment(c string) *PipeBuilder {
	pb.comment = c
	return pb
}

// WithCopyStatement adds a URL to the PipeBuilder
func (pb *PipeBuilder) WithCopyStatement(s string) *PipeBuilder {
	pb.copyStatement = s
	return pb
}

/// WithIntegration adds Integration specification to the PipeBuilder
func (pb *PipeBuilder) WithIntegration(s string) *PipeBuilder {
	pb.integration = s
	return pb
}

/// WithErrorIntegration adds ErrorIntegration specification to the PipeBuilder
func (pb *PipeBuilder) WithErrorIntegration(s string) *PipeBuilder {
	pb.errorIntegration = s
	return pb
}

// Pipe returns a pointer to a Builder that abstracts the DDL operations for a pipe.
//
// Supported DDL operations are:
//   - CREATE PIPE
//   - ALTER PIPE
//   - DROP PIPE
//   - SHOW PIPE
//
// [Snowflake Reference](https://docs.snowflake.net/manuals/sql-reference/ddl-pipe.html#pipe-management)
func Pipe(name, db, schema string) *PipeBuilder {
	return &PipeBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL statement required to create a pipe
func (pb *PipeBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` PIPE %v`, pb.QualifiedName()))

	if pb.autoIngest {
		q.WriteString(` AUTO_INGEST = TRUE`)
	}

	if pb.integration != "" {
		q.WriteString(fmt.Sprintf(` INTEGRATION = '%v'`, EscapeString(pb.integration)))
	}

	if pb.errorIntegration != "" {
		q.WriteString(fmt.Sprintf(` ERROR_INTEGRATION = '%v'`, EscapeString(pb.errorIntegration)))
	}

	if pb.awsSnsTopicArn != "" {
		q.WriteString(fmt.Sprintf(` AWS_SNS_TOPIC = '%v'`, EscapeString(pb.awsSnsTopicArn)))
	}

	if pb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(pb.comment)))
	}
	if pb.copyStatement != "" {
		q.WriteString(fmt.Sprintf(` AS %v`, pb.copyStatement))
	}

	return q.String()
}

// ChangeComment returns the SQL query that will update the comment on the pipe.
func (pb *PipeBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER PIPE %v SET COMMENT = '%v'`, pb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the pipe.
func (pb *PipeBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER PIPE %v UNSET COMMENT`, pb.QualifiedName())
}

// ChangeErrorIntegration return SQL query that will update the error_integration on the pipe.
func (pb *PipeBuilder) ChangeErrorIntegration(c string) string {
	return fmt.Sprintf(`ALTER PIPE %v SET ERROR_INTEGRATION = %v`, pb.QualifiedName(), EscapeString(c))
}

// RemoveErrorIntegration returns the SQL query that will remove the error_integration on the pipe.
func (pb *PipeBuilder) RemoveErrorIntegration() string {
	return fmt.Sprintf(`ALTER PIPE %v UNSET ERROR_INTEGRATION`, pb.QualifiedName())
}

// Drop returns the SQL query that will drop a pipe.
func (pb *PipeBuilder) Drop() string {
	return fmt.Sprintf(`DROP PIPE %v`, pb.QualifiedName())
}

// Show returns the SQL query that will show a pipe.
func (pb *PipeBuilder) Show() string {
	return fmt.Sprintf(`SHOW PIPES LIKE '%v' IN SCHEMA "%v"."%v"`, pb.name, pb.db, pb.schema)
}

type pipe struct {
	Createdon           string         `db:"created_on"`
	Name                string         `db:"name"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	Definition          string         `db:"definition"`
	Owner               string         `db:"owner"`
	NotificationChannel *string        `db:"notification_channel"`
	Comment             string         `db:"comment"`
	Integration         sql.NullString `db:"integration"`
	ErrorIntegration    sql.NullString `db:"error_integration"`
}

func ScanPipe(row *sqlx.Row) (*pipe, error) {
	p := &pipe{}
	e := row.StructScan(p)
	return p, e
}

func ListPipes(databaseName string, schemaName string, db *sql.DB) ([]pipe, error) {
	stmt := fmt.Sprintf(`SHOW PIPES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []pipe{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no pipes found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
