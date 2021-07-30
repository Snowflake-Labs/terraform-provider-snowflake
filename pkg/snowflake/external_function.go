package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ExternalFunctionBuilder abstracts the creation of SQL queries for a Snowflake schema
type ExternalFunctionBuilder struct {
	name                  string
	db                    string
	schema                string
	args                  []map[string]string
	argtypes              string // only used for 'DESC FUNCTION' & 'DROP FUNCTION' commands as of today (list of args types is required)
	nullInputBehavior     string
	returnType            string
	returnNullAllowed     bool
	returnBehavior        string
	apiIntegration        string
	headers               []map[string]string
	contextHeaders        []string
	maxBatchRows          int
	compression           string
	urlOfProxyAndResource string
	comment               string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (fb *ExternalFunctionBuilder) QualifiedName() string {
	var n strings.Builder

	if fb.db != "" && fb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, fb.db, fb.schema))
	}

	if fb.db != "" && fb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, fb.db))
	}

	if fb.db == "" && fb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, fb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, fb.name))

	return n.String()
}

// QualifiedNameWithArgTypes appends all args' types to the qualified name. This is required to invoke 'DESC FUNCTION' and 'DROP FUNCTION' commands.
func (fb *ExternalFunctionBuilder) QualifiedNameWithArgTypes() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`%v (%s)`, fb.QualifiedName(), fb.argtypes))
	return q.String()
}

// WithArgs sets the args on the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithArgs(args []map[string]string) *ExternalFunctionBuilder {
	fb.args = args
	return fb
}

// WithArgTypes sets the args on the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithArgTypes(argtypes string) *ExternalFunctionBuilder {
	argtypeslist := strings.ReplaceAll(argtypes, "-", ", ")
	fb.argtypes = argtypeslist
	return fb
}

// WithNullInputBehavior adds a nullInputBehavior to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithNullInputBehavior(nullInputBehavior string) *ExternalFunctionBuilder {
	fb.nullInputBehavior = nullInputBehavior
	return fb
}

// WithReturnType adds a returnType to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithReturnType(returnType string) *ExternalFunctionBuilder {
	fb.returnType = returnType
	return fb
}

// WithReturnNullAllowed adds a returnNullAllowed to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithReturnNullAllowed(returnNullAllowed bool) *ExternalFunctionBuilder {
	fb.returnNullAllowed = returnNullAllowed
	return fb
}

// WithReturnBehavior adds a returnBehavior to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithReturnBehavior(returnBehavior string) *ExternalFunctionBuilder {
	fb.returnBehavior = returnBehavior
	return fb
}

// WithAPIIntegration adds a apiIntegration to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithAPIIntegration(apiIntegration string) *ExternalFunctionBuilder {
	fb.apiIntegration = apiIntegration
	return fb
}

// WithHeaders sets the headers on the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithHeaders(headers []map[string]string) *ExternalFunctionBuilder {
	fb.headers = headers
	return fb
}

// WithContextHeaders sets the context headers on the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithContextHeaders(contextHeaders []string) *ExternalFunctionBuilder {
	fb.contextHeaders = contextHeaders
	return fb
}

// WithMaxBatchRows adds a maxBatchRows to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithMaxBatchRows(maxBatchRows int) *ExternalFunctionBuilder {
	fb.maxBatchRows = maxBatchRows
	return fb
}

// WithCompression adds a compression to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithCompression(compression string) *ExternalFunctionBuilder {
	fb.compression = compression
	return fb
}

// WithURLOfProxyAndResource adds a urlOfProxyAndResource to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithURLOfProxyAndResource(urlOfProxyAndResource string) *ExternalFunctionBuilder {
	fb.urlOfProxyAndResource = urlOfProxyAndResource
	return fb
}

// WithComment adds a comment to the ExternalFunctionBuilder
func (fb *ExternalFunctionBuilder) WithComment(c string) *ExternalFunctionBuilder {
	fb.comment = c
	return fb
}

// ExternalFunction returns a pointer to a Builder that abstracts the DDL operations for an external function.
//
// Supported DDL operations are:
//   - CREATE EXTERNAL FUNCTION
//   - ALTER EXTERNAL FUNCTION
//   - DROP FUNCTION
//   - SHOW EXTERNAL FUNCTIONS
//   - DESCRIBE FUNCTION
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/ddl-udf.html#external-function-management)
func ExternalFunction(name, db, schema string) *ExternalFunctionBuilder {
	return &ExternalFunctionBuilder{
		name:              name,
		db:                db,
		schema:            schema,
		returnNullAllowed: true,
	}
}

// Create returns the SQL statement required to create an  external function
func (fb *ExternalFunctionBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE EXTERNAL FUNCTION %v`, fb.QualifiedName()))

	q.WriteString(` (`)
	args := []string{}
	for _, arg := range fb.args {
		args = append(args, fmt.Sprintf(`%v %v`, EscapeString(arg["name"]), EscapeString(arg["type"])))
	}
	q.WriteString(strings.Join(args, ", "))
	q.WriteString(`)`)

	q.WriteString(` RETURNS ` + EscapeString(fb.returnType))

	if !fb.returnNullAllowed {
		q.WriteString(` NOT`)
	}
	q.WriteString(` NULL`)

	if fb.nullInputBehavior != "" {
		q.WriteString(fmt.Sprintf(` %v`, EscapeString(fb.nullInputBehavior)))
	}

	q.WriteString(fmt.Sprintf(` %v`, EscapeString(fb.returnBehavior)))

	if fb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(fb.comment)))
	}

	q.WriteString(fmt.Sprintf(` API_INTEGRATION = '%v'`, EscapeString(fb.apiIntegration)))

	if len(fb.headers) > 0 {
		q.WriteString(` HEADERS = (`)
		headers := []string{}
		for _, header := range fb.headers {
			headers = append(headers, fmt.Sprintf(`'%v' = '%v'`, EscapeString(header["name"]), EscapeString(header["value"])))
		}
		q.WriteString(strings.Join(headers, ", "))
		q.WriteString(`)`)
	}

	if len(fb.contextHeaders) > 0 {
		q.WriteString(` CONTEXT_HEADERS = (`)
		q.WriteString(EscapeString(strings.Join(fb.contextHeaders, ", ")))
		q.WriteString(`)`)
	}

	if fb.maxBatchRows > 0 {
		q.WriteString(fmt.Sprintf(` MAX_BATCH_ROWS = %d`, fb.maxBatchRows))
	}

	if fb.compression != "" {
		q.WriteString(fmt.Sprintf(` COMPRESSION = '%v'`, EscapeString(fb.compression)))
	}

	q.WriteString(fmt.Sprintf(` AS '%v'`, EscapeString(fb.urlOfProxyAndResource)))

	return q.String()
}

// Drop returns the SQL query that will drop an external function.
func (fb *ExternalFunctionBuilder) Drop() string {
	return fmt.Sprintf(`DROP FUNCTION %v`, fb.QualifiedNameWithArgTypes())
}

// Show returns the SQL query that will show an external function.
func (fb *ExternalFunctionBuilder) Show() string {
	return fmt.Sprintf(`SHOW EXTERNAL FUNCTIONS LIKE '%v' IN SCHEMA "%v"."%v"`, fb.name, fb.db, fb.schema)
}

// Describe returns the SQL query that will describe an external function.
func (fb *ExternalFunctionBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE FUNCTION %s`, fb.QualifiedNameWithArgTypes())
}

type externalFunction struct {
	CreatedOn            sql.NullString `db:"created_on"`
	ExternalFunctionName sql.NullString `db:"name"`
	DatabaseName         sql.NullString `db:"catalog_name"`
	SchemaName           sql.NullString `db:"schema_name"`
	Comment              sql.NullString `db:"description"`
	IsExternalFunction   sql.NullString `db:"is_external_function"`
	Language             sql.NullString `db:"language"`
}

// ScanExternalFunction
func ScanExternalFunction(row *sqlx.Row) (*externalFunction, error) {
	f := &externalFunction{}
	e := row.StructScan(f)
	return f, e
}

type externalFunctionDescription struct {
	Property sql.NullString `db:"property"`
	Value    sql.NullString `db:"value"`
}

// ScanExternalFunctionDescription
func ScanExternalFunctionDescription(rows *sqlx.Rows) ([]externalFunctionDescription, error) {
	efds := []externalFunctionDescription{}
	for rows.Next() {
		efd := externalFunctionDescription{}
		err := rows.StructScan(&efd)
		if err != nil {
			return nil, err
		}
		efds = append(efds, efd)
	}
	return efds, rows.Err()
}

func ListExternalFunctions(databaseName string, schemaName string, db *sql.DB) ([]externalFunction, error) {
	stmt := fmt.Sprintf(`SHOW EXTERNAL FUNCTIONS IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []externalFunction{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no external functions found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
