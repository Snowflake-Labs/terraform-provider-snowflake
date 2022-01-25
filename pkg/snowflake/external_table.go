package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// externalTableBuilder abstracts the creation of SQL queries for a Snowflake schema
type ExternalTableBuilder struct {
	name            string
	db              string
	schema          string
	columns         []map[string]string
	partitionBys    []string
	location        string
	refreshOnCreate bool
	autoRefresh     bool
	pattern         string
	fileFormat      string
	copyGrants      bool
	awsSNSTopic     string
	comment         string
	tags            []TagValue
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (tb *ExternalTableBuilder) QualifiedName() string {
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

// WithComment adds a comment to the ExternalTableBuilder
func (tb *ExternalTableBuilder) WithComment(c string) *ExternalTableBuilder {
	tb.comment = c
	return tb
}

// WithColumns sets the column definitions on the ExternalTableBuilder
func (tb *ExternalTableBuilder) WithColumns(c []map[string]string) *ExternalTableBuilder {
	tb.columns = c
	return tb
}
func (tb *ExternalTableBuilder) WithPartitionBys(c []string) *ExternalTableBuilder {
	tb.partitionBys = c
	return tb
}
func (tb *ExternalTableBuilder) WithLocation(c string) *ExternalTableBuilder {
	tb.location = c
	return tb
}
func (tb *ExternalTableBuilder) WithRefreshOnCreate(c bool) *ExternalTableBuilder {
	tb.refreshOnCreate = c
	return tb
}
func (tb *ExternalTableBuilder) WithAutoRefresh(c bool) *ExternalTableBuilder {
	tb.autoRefresh = c
	return tb
}
func (tb *ExternalTableBuilder) WithPattern(c string) *ExternalTableBuilder {
	tb.pattern = c
	return tb
}
func (tb *ExternalTableBuilder) WithFileFormat(c string) *ExternalTableBuilder {
	tb.fileFormat = c
	return tb
}
func (tb *ExternalTableBuilder) WithCopyGrants(c bool) *ExternalTableBuilder {
	tb.copyGrants = c
	return tb
}
func (tb *ExternalTableBuilder) WithAwsSNSTopic(c string) *ExternalTableBuilder {
	tb.awsSNSTopic = c
	return tb
}

// WithTags sets the tags on the ExternalTableBuilder
func (tb *ExternalTableBuilder) WithTags(tags []TagValue) *ExternalTableBuilder {
	tb.tags = tags
	return tb
}

// ExternalexternalTable returns a pointer to a Builder that abstracts the DDL operations for a externalTable.
//
// Supported DDL operations are:
//   - CREATE externalTable
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/sql/create-external-table.html)

func ExternalTable(name, db, schema string) *ExternalTableBuilder {
	return &ExternalTableBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL statement required to create a externalTable
func (tb *ExternalTableBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE EXTERNAL TABLE %v`, tb.QualifiedName()))

	q.WriteString(fmt.Sprintf(` (`))
	columnDefinitions := []string{}
	for _, columnDefinition := range tb.columns {
		columnDefinitions = append(columnDefinitions, fmt.Sprintf(`"%v" %v AS %v`, EscapeString(columnDefinition["name"]), EscapeString(columnDefinition["type"]), columnDefinition["as"]))
	}
	q.WriteString(strings.Join(columnDefinitions, ", "))
	q.WriteString(fmt.Sprintf(`)`))

	if len(tb.partitionBys) > 0 {
		q.WriteString(` PARTITION BY ( `)
		q.WriteString(EscapeString(strings.Join(tb.partitionBys, ", ")))
		q.WriteString(` )`)
	}

	q.WriteString(` WITH LOCATION = ` + EscapeString(tb.location))
	q.WriteString(fmt.Sprintf(` REFRESH_ON_CREATE = %t`, tb.refreshOnCreate))
	q.WriteString(fmt.Sprintf(` AUTO_REFRESH = %t`, tb.autoRefresh))

	if tb.pattern != "" {
		q.WriteString(fmt.Sprintf(` PATTERN = '%v'`, EscapeString(tb.pattern)))
	}

	q.WriteString(fmt.Sprintf(` FILE_FORMAT = ( %v )`, EscapeString(tb.fileFormat)))

	if tb.awsSNSTopic != "" {
		q.WriteString(fmt.Sprintf(` AWS_SNS_TOPIC = '%v'`, EscapeString(tb.awsSNSTopic)))
	}

	if tb.copyGrants {
		q.WriteString(" COPY GRANTS")
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	if len(tb.tags) > 0 {
		q.WriteString(fmt.Sprintf(` WITH TAG (%s)`, tb.GetTagValueString()))
	}

	return q.String()
}

// Update returns the SQL statement required to update an externalTable
func (tb *ExternalTableBuilder) Update() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER EXTERNAL TABLE %v`, tb.QualifiedName()))

	if len(tb.tags) > 0 {
		q.WriteString(fmt.Sprintf(` TAG %s`, tb.GetTagValueString()))
	}

	return q.String()
}

// Drop returns the SQL query that will drop a externalTable.
func (tb *ExternalTableBuilder) Drop() string {
	return fmt.Sprintf(`DROP EXTERNAL TABLE %v`, tb.QualifiedName())
}

// Show returns the SQL query that will show a externalTable.
func (tb *ExternalTableBuilder) Show() string {
	return fmt.Sprintf(`SHOW EXTERNAL TABLES LIKE '%v' IN SCHEMA "%v"."%v"`, tb.name, tb.db, tb.schema)
}

func (tb *ExternalTableBuilder) GetTagValueString() string {
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

type externalTable struct {
	CreatedOn         sql.NullString `db:"created_on"`
	ExternalTableName sql.NullString `db:"name"`
	DatabaseName      sql.NullString `db:"database_name"`
	SchemaName        sql.NullString `db:"schema_name"`
	Comment           sql.NullString `db:"comment"`
	Owner             sql.NullString `db:"owner"`
}

func ScanExternalTable(row *sqlx.Row) (*externalTable, error) {
	t := &externalTable{}
	e := row.StructScan(t)
	return t, e
}

func ListExternalTables(databaseName string, schemaName string, db *sql.DB) ([]externalTable, error) {
	stmt := fmt.Sprintf(`SHOW EXTERNAL TABLES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []externalTable{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no external tables found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
