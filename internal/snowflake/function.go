// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

// FunctionBuilder abstracts the creation of Function.
type FunctionBuilder struct {
	name              string
	schema            string
	db                string
	argumentTypes     []string // (VARCHAR, VARCHAR)
	args              []map[string]string
	returnBehavior    string // VOLATILE, IMMUTABLE
	nullInputBehavior string // "CALLED ON NULL INPUT" or "RETURNS NULL ON NULL INPUT"
	returnType        string
	language          string
	packages          []string
	imports           []string // for Java / Python imports
	handler           string   // for Java / Python handler
	targetPath        string   // for Java / Python target path
	comment           string
	statement         string
	runtimeVersion    string // for Python runtime version
	secure            bool
}

// QualifiedName prepends the db and schema and appends argument types.
func (pb *FunctionBuilder) QualifiedName() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("functions must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"(%v)`, pb.db, pb.schema, pb.name, strings.Join(pb.argumentTypes, ", ")), nil
}

// QualifiedNameWithoutArguments prepends the db and schema if set.
func (pb *FunctionBuilder) QualifiedNameWithoutArguments() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("functions must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"`, pb.db, pb.schema, pb.name), nil
}

// Returns the arguments signature of the function in a form <function>(<type>, <type>, ..) RETURN <type>.
func (pb *FunctionBuilder) ArgumentsSignature() (string, error) {
	return fmt.Sprintf(`%v(%v) RETURN %v`, pb.name, strings.Join(pb.argumentTypes, ", "), pb.returnType), nil
}

// WithArgs sets the args and argumentTypes on the FunctionBuilder.
func (pb *FunctionBuilder) WithArgs(args []map[string]string) *FunctionBuilder {
	pb.args = []map[string]string{}
	for _, arg := range args {
		argName := arg["name"]
		argType := strings.ToUpper(arg["type"])
		pb.args = append(pb.args, map[string]string{"name": argName, "type": argType})
		pb.argumentTypes = append(pb.argumentTypes, argType)
	}
	return pb
}

// WithRuntimeVersion.
func (pb *FunctionBuilder) WithRuntimeVersion(r string) *FunctionBuilder {
	pb.runtimeVersion = r
	return pb
}

// WithReturnBehavior.
func (pb *FunctionBuilder) WithReturnBehavior(s string) *FunctionBuilder {
	pb.returnBehavior = s
	return pb
}

// WithNullInputBehavior.
func (pb *FunctionBuilder) WithNullInputBehavior(s string) *FunctionBuilder {
	pb.nullInputBehavior = s
	return pb
}

// WithReturnType adds the data type of the return type to the FunctionBuilder.
func (pb *FunctionBuilder) WithReturnType(s string) *FunctionBuilder {
	pb.returnType = strings.ToUpper(s)
	return pb
}

// WithLanguage sets the language to SQL, JAVA or JAVASCRIPT.
func (pb *FunctionBuilder) WithLanguage(s string) *FunctionBuilder {
	pb.language = s
	return pb
}

// WithPackages.
func (pb *FunctionBuilder) WithPackages(s []string) *FunctionBuilder {
	pb.packages = s
	return pb
}

// WithImports adds jar files to import for Java function or Python file for Python function.
func (pb *FunctionBuilder) WithImports(s []string) *FunctionBuilder {
	pb.imports = s
	return pb
}

// WithHandler sets the handler method for Java / Python function.
func (pb *FunctionBuilder) WithHandler(s string) *FunctionBuilder {
	pb.handler = s
	return pb
}

// WithTargetPath sets the target path for compiled jar file for Java function or Python file for Python function.
func (pb *FunctionBuilder) WithTargetPath(s string) *FunctionBuilder {
	pb.targetPath = s
	return pb
}

// WithSecure sets the secure boolean to true
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/sql/create-function)
func (pb *FunctionBuilder) WithSecure() *FunctionBuilder {
	pb.secure = true
	return pb
}

// WithComment adds a comment to the FunctionBuilder.
func (pb *FunctionBuilder) WithComment(c string) *FunctionBuilder {
	pb.comment = c
	return pb
}

// WithStatement adds the SQL/JAVASCRIPT/JAVA statement to be used for the function.
func (pb *FunctionBuilder) WithStatement(s string) *FunctionBuilder {
	pb.statement = s
	return pb
}

// Returns the argument types.
func (pb *FunctionBuilder) ArgTypes() []string {
	return pb.argumentTypes
}

// Function returns a pointer to a Builder that abstracts the DDL operations for a stored function.
//
// Supported DDL operations are:
//   - CREATE FUNCTION
//   - ALTER FUNCTION
//   - DROP FUNCTION
//   - SHOW FUNCTION
//   - DESCRIBE
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/user-defined-functions.html)
func NewFunctionBuilder(db, schema, name string, argTypes []string) *FunctionBuilder {
	return &FunctionBuilder{
		name:          name,
		db:            db,
		schema:        schema,
		argumentTypes: argTypes,
	}
}

// Create returns the SQL query that will create a new function.
func (pb *FunctionBuilder) Create() (string, error) {
	var q strings.Builder

	q.WriteString("CREATE OR REPLACE")

	if pb.secure {
		q.WriteString(" SECURE")
	}

	qn, err := pb.QualifiedNameWithoutArguments()
	if err != nil {
		return "", err
	}

	q.WriteString(fmt.Sprintf(" FUNCTION %v", qn))

	q.WriteString(`(`)
	args := []string{}
	for _, arg := range pb.args {
		args = append(args, fmt.Sprintf(`%v %v`, EscapeString(arg["name"]), EscapeString(arg["type"])))
	}
	q.WriteString(strings.Join(args, ", "))
	q.WriteString(`)`)

	q.WriteString(fmt.Sprintf(" RETURNS %v", pb.returnType))
	if pb.language != "" {
		q.WriteString(fmt.Sprintf(" LANGUAGE %v", pb.language))
	}

	if pb.nullInputBehavior != "" {
		q.WriteString(fmt.Sprintf(` %v`, EscapeString(pb.nullInputBehavior)))
	}
	if pb.returnBehavior != "" {
		q.WriteString(fmt.Sprintf(` %v`, EscapeString(pb.returnBehavior)))
	}

	if pb.runtimeVersion != "" {
		q.WriteString(fmt.Sprintf(" RUNTIME_VERSION = '%v'", EscapeString(pb.runtimeVersion)))
	}

	if len(pb.packages) > 0 {
		q.WriteString(` PACKAGES = (`)
		packages := []string{}
		for _, pack := range pb.packages {
			packages = append(packages, fmt.Sprintf(`'%v'`, pack))
		}
		q.WriteString(strings.Join(packages, ", "))
		q.WriteString(`)`)
	}

	if pb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", EscapeString(pb.comment)))
	}

	if len(pb.imports) > 0 {
		q.WriteString(` IMPORTS = (`)
		imports := []string{}
		for _, imp := range pb.imports {
			imports = append(imports, fmt.Sprintf(`'%v'`, imp))
		}
		q.WriteString(strings.Join(imports, ", "))
		q.WriteString(`)`)
	}

	if pb.handler != "" {
		q.WriteString(fmt.Sprintf(" HANDLER = '%v'", pb.handler))
	}

	if pb.targetPath != "" {
		q.WriteString(fmt.Sprintf(" TARGET_PATH = '%v'", pb.targetPath))
	}

	q.WriteString(fmt.Sprintf(" AS $$%v$$", pb.statement))
	return q.String(), nil
}

// Rename returns the SQL query that will rename the function.
func (pb *FunctionBuilder) Rename(newName string) (string, error) {
	oldName, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	pb.name = newName

	qn, err := pb.QualifiedNameWithoutArguments()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER FUNCTION %v RENAME TO %v`, oldName, qn), nil
}

// Secure returns the SQL query that will change the function to a secure function.
func (pb *FunctionBuilder) Secure() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER FUNCTION %v SET SECURE`, qn), nil
}

// Unsecure returns the SQL query that will change the function to a normal (unsecured) function.
func (pb *FunctionBuilder) Unsecure() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER FUNCTION %v UNSET SECURE`, qn), nil
}

// ChangeComment returns the SQL query that will update the comment on the function.
func (pb *FunctionBuilder) ChangeComment(c string) (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`ALTER FUNCTION %v SET COMMENT = '%v'`, qn, EscapeString(c)), nil
}

// RemoveComment returns the SQL query that will remove the comment on the function.
func (pb *FunctionBuilder) RemoveComment() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER FUNCTION %v UNSET COMMENT`, qn), nil
}

// Show returns the SQL query that will show the row representing this function.
// This show statement returns all functions with the given name (overloaded ones).
func (pb *FunctionBuilder) Show() string {
	return fmt.Sprintf(`SHOW USER FUNCTIONS LIKE '%v' IN SCHEMA "%v"."%v"`, pb.name, pb.db, pb.schema)
}

// To describe the function the name must be specified as fully qualified name
// including argument types.
func (pb *FunctionBuilder) Describe() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`DESCRIBE FUNCTION %v`, qn), nil
}

// Drop returns the SQL query that will drop the function.
func (pb *FunctionBuilder) Drop() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`DROP FUNCTION %v`, qn), nil
}

type Function struct {
	Comment      sql.NullString `db:"description"`
	IsSecure     sql.NullString `db:"is_secure"`
	Name         sql.NullString `db:"name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Text         sql.NullString `db:"text"`
	DatabaseName sql.NullString `db:"database_name"`
	Arguments    sql.NullString `db:"arguments"`
}

type FunctionDescription struct {
	Property sql.NullString `db:"property"`
	Value    sql.NullString `db:"value"`
}

// ScanFunctionDescription reads through the rows with property and value columns
// and returns a slice of functionDescription structs.
func ScanFunctionDescription(rows *sqlx.Rows) ([]FunctionDescription, error) {
	pdsl := []FunctionDescription{}
	for rows.Next() {
		pd := FunctionDescription{}
		err := rows.StructScan(&pd)
		if err != nil {
			return nil, err
		}
		pdsl = append(pdsl, pd)
	}
	return pdsl, rows.Err()
}

// SHOW FUNCTION can return more than one item because of function names overloading
// https://docs.snowflake.com/en/sql-reference/sql/show-functions.html
func ScanFunctions(rows *sqlx.Rows) ([]*Function, error) {
	var pcs []*Function
	for rows.Next() {
		r := &Function{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, r)
	}
	return pcs, rows.Err()
}

type UserFunctions struct {
	Name        sql.NullString `db:"name"`
	Arguments   sql.NullString `db:"arguments"`
	Description sql.NullString `db:"description"`
	Language    sql.NullString `db:"language"`
}

func ListUserFunctions(databaseName string, schemaName string, db *sql.DB) ([]UserFunctions, error) {
	stmt := fmt.Sprintf(`SHOW USER FUNCTIONS IN SCHEMA "%v"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []UserFunctions{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no functions found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
}
