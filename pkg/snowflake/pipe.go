package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// PipeBuilder abstracts the creation of SQL queries for a Snowflake schema.
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

// QualifiedName prepends the db and schema if set and escapes everything nicely.
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

type Pipe struct {
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

func ListPipes(databaseName string, schemaName string, db *sql.DB) ([]Pipe, error) {
	stmt := fmt.Sprintf(`SHOW PIPES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Pipe{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no pipes found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
