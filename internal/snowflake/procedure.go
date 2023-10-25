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

// ProcedureBuilder abstracts the creation of Stored Procedure.
type ProcedureBuilder struct {
	name              string
	schema            string
	db                string
	argumentTypes     []string // (VARCHAR, VARCHAR)
	args              []map[string]string
	returnBehavior    string // VOLATILE, IMMUTABLE
	nullInputBehavior string // "CALLED ON NULL INPUT" or "RETURNS NULL ON NULL INPUT"
	returnType        string
	language          string // SQL, JAVASCRIPT, JAVA, SCALA
	packages          []string
	imports           []string // for Java / Python imports
	handler           string   // for Java / Python handler
	executeAs         string
	comment           string
	statement         string
	runtimeVersion    string // for Python runtime version
}

// QualifiedName prepends the db and schema and appends argument types.
func (pb *ProcedureBuilder) QualifiedName() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("procedures must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"(%v)`, pb.db, pb.schema, pb.name, strings.Join(pb.argumentTypes, ", ")), nil
}

// QualifiedNameWithoutArguments prepends the db and schema if set.
func (pb *ProcedureBuilder) QualifiedNameWithoutArguments() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("procedures must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"`, pb.db, pb.schema, pb.name), nil
}

// Returns the arguments signature of the procedure in a form <PROCEDURE>(<TYPE>, <TYPE>, ..)
func (pb *ProcedureBuilder) ArgumentsSignature() (string, error) {
	return fmt.Sprintf(`%v(%v)`, strings.ToUpper(pb.name), strings.ToUpper(strings.Join(pb.argumentTypes, ", "))), nil
}

// WithArgs sets the args and argumentTypes on the ProcedureBuilder.
func (pb *ProcedureBuilder) WithArgs(args []map[string]string) *ProcedureBuilder {
	pb.args = []map[string]string{}
	for _, arg := range args {
		argName := arg["name"]
		argType := strings.ToUpper(arg["type"])
		pb.args = append(pb.args, map[string]string{"name": argName, "type": argType})
		pb.argumentTypes = append(pb.argumentTypes, argType)
	}
	return pb
}

// WithReturnBehavior.
func (pb *ProcedureBuilder) WithReturnBehavior(s string) *ProcedureBuilder {
	pb.returnBehavior = s
	return pb
}

// WithNullInputBehavior.
func (pb *ProcedureBuilder) WithNullInputBehavior(s string) *ProcedureBuilder {
	pb.nullInputBehavior = s
	return pb
}

// WithReturnType adds the data type of the return type to the ProcedureBuilder.
func (pb *ProcedureBuilder) WithReturnType(s string) *ProcedureBuilder {
	pb.returnType = strings.ToUpper(s)
	return pb
}

// WithExecuteAs sets the execute to OWNER or CALLER.
func (pb *ProcedureBuilder) WithExecuteAs(s string) *ProcedureBuilder {
	pb.executeAs = s
	return pb
}

// WithLanguage sets the language to SQL, JAVA, SCALA or JAVASCRIPT.
func (pb *ProcedureBuilder) WithLanguage(s string) *ProcedureBuilder {
	pb.language = s
	return pb
}

// WithRuntimeVersion.
func (pb *ProcedureBuilder) WithRuntimeVersion(r string) *ProcedureBuilder {
	pb.runtimeVersion = r
	return pb
}

// WithPackages.
func (pb *ProcedureBuilder) WithPackages(s []string) *ProcedureBuilder {
	pb.packages = s
	return pb
}

// WithImports adds jar files to import for Java function or Python file for Python function.
func (pb *ProcedureBuilder) WithImports(s []string) *ProcedureBuilder {
	pb.imports = s
	return pb
}

// WithHandler sets the handler method for Java / Python function.
func (pb *ProcedureBuilder) WithHandler(s string) *ProcedureBuilder {
	pb.handler = s
	return pb
}

// WithComment adds a comment to the ProcedureBuilder.
func (pb *ProcedureBuilder) WithComment(c string) *ProcedureBuilder {
	pb.comment = c
	return pb
}

// WithStatement adds the SQL statement to be used for the procedure.
func (pb *ProcedureBuilder) WithStatement(s string) *ProcedureBuilder {
	pb.statement = s
	return pb
}

// Returns the argument types.
func (pb *ProcedureBuilder) ArgTypes() []string {
	return pb.argumentTypes
}

// Procedure returns a pointer to a Builder that abstracts the DDL operations for a stored procedure.
//
// Supported DDL operations are:
//   - CREATE PROCEDURE
//   - ALTER PROCEDURE
//   - DROP PROCEDURE
//   - SHOW PROCEDURE
//   - DESCRIBE
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/stored-procedures.html)
func NewProcedureBuilder(db, schema, name string, argTypes []string) *ProcedureBuilder {
	return &ProcedureBuilder{
		name:          name,
		db:            db,
		schema:        schema,
		argumentTypes: argTypes,
	}
}

// Create returns the SQL query that will create a new procedure.
func (pb *ProcedureBuilder) Create() (string, error) {
	var q strings.Builder

	q.WriteString("CREATE OR REPLACE")

	qn, err := pb.QualifiedNameWithoutArguments()
	if err != nil {
		return "", err
	}

	q.WriteString(fmt.Sprintf(" PROCEDURE %v", qn))

	q.WriteString(`(`)
	args := []string{}
	for _, arg := range pb.args {
		args = append(args, fmt.Sprintf(`%v %v`, EscapeString(arg["name"]), EscapeString(arg["type"])))
	}
	q.WriteString(strings.Join(args, ", "))
	q.WriteString(`)`)

	q.WriteString(fmt.Sprintf(" RETURNS %v", pb.returnType))
	if pb.language != "" {
		q.WriteString(fmt.Sprintf(" LANGUAGE %v", EscapeString(pb.language)))
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
	if pb.comment != "" {
		q.WriteString(fmt.Sprintf(" COMMENT = '%v'", EscapeString(pb.comment)))
	}
	q.WriteString(fmt.Sprintf(" EXECUTE AS %v", pb.executeAs))
	q.WriteString(fmt.Sprintf(" AS $$%v$$", pb.statement))
	return q.String(), nil
}

// Rename returns the SQL query that will rename the procedure.
func (pb *ProcedureBuilder) Rename(newName string) (string, error) {
	oldName, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	pb.name = newName

	qn, err := pb.QualifiedNameWithoutArguments()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER PROCEDURE %v RENAME TO %v`, oldName, qn), nil
}

// ChangeComment returns the SQL query that will update the comment on the procedure.
func (pb *ProcedureBuilder) ChangeComment(c string) (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`ALTER PROCEDURE %v SET COMMENT = '%v'`, qn, EscapeString(c)), nil
}

// RemoveComment returns the SQL query that will remove the comment on the procedure.
func (pb *ProcedureBuilder) RemoveComment() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER PROCEDURE %v UNSET COMMENT`, qn), nil
}

// ChangeExecuteAs returns the SQL query that will update the call mode on the procedure.
func (pb *ProcedureBuilder) ChangeExecuteAs(c string) (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER PROCEDURE %v EXECUTE AS %v`, qn, c), nil
}

// Show returns the SQL query that will show the row representing this procedure.
// This show statement returns all procedures with the given name (overloaded ones).
func (pb *ProcedureBuilder) Show() string {
	return fmt.Sprintf(`SHOW PROCEDURES LIKE '%v' IN SCHEMA "%v"."%v"`, pb.name, pb.db, pb.schema)
}

// To describe the procedure the name must be specified as fully qualified name
// including argument types.
func (pb *ProcedureBuilder) Describe() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`DESCRIBE PROCEDURE %v`, qn), nil
}

// Drop returns the SQL query that will drop the procedure.
func (pb *ProcedureBuilder) Drop() (string, error) {
	qn, err := pb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`DROP PROCEDURE %v`, qn), nil
}

type Procedure struct {
	Comment sql.NullString `db:"description"`
	// Snowflake returns is_secure in the show procedure output, but it is irrelevant
	Name         sql.NullString `db:"name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Text         sql.NullString `db:"text"`
	DatabaseName sql.NullString `db:"catalog_name"`
	Arguments    sql.NullString `db:"arguments"`
}

type ProcedureDescription struct {
	Property sql.NullString `db:"property"`
	Value    sql.NullString `db:"value"`
}

// ScanProcedureDescription reads through the rows with property and value columns
// and returns a slice of procedureDescription structs.
func ScanProcedureDescription(rows *sqlx.Rows) ([]ProcedureDescription, error) {
	pdsl := []ProcedureDescription{}
	for rows.Next() {
		pd := ProcedureDescription{}
		err := rows.StructScan(&pd)
		if err != nil {
			return nil, err
		}
		pdsl = append(pdsl, pd)
	}
	return pdsl, rows.Err()
}

// SHOW PROCEDURE can return more than one item because of procedure names overloading
// https://docs.snowflake.com/en/sql-reference/sql/show-procedures.html
func ScanProcedures(rows *sqlx.Rows) ([]*Procedure, error) {
	var pcs []*Procedure
	for rows.Next() {
		r := &Procedure{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, r)
	}
	return pcs, rows.Err()
}

func ListProcedures(databaseName string, schemaName string, db *sql.DB) ([]Procedure, error) {
	stmt := fmt.Sprintf(`SHOW PROCEDURES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []Procedure{}
	if err := sqlx.StructScan(rows, &dbs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[DEBUG] no procedures found")
			return nil, nil
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}
	return dbs, nil
}
