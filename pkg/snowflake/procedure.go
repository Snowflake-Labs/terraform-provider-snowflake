package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	pe "github.com/pkg/errors"
)

// ProcedureBuilder abstracts the creation of Stored Procedure
type ProcedureBuilder struct {
	name              string
	schema            string
	db                string
	argumentTypes     []string // (VARCHAR, VARCHAR)
	args              []map[string]string
	returnBehavior    string // VOLATILE, IMMUTABLE
	nullInputBehavior string // "CALLED ON NULL INPUT" or "RETURNS NULL ON NULL INPUT"
	returnType        string
	executeAs         string
	comment           string
	statement         string
}

// QualifiedName prepends the db and schema and appends argument types
func (pb *ProcedureBuilder) QualifiedName() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("Procedures must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"(%v)`, pb.db, pb.schema, pb.name, strings.Join(pb.argumentTypes, ", ")), nil
}

// QualifiedNameWithoutArguments prepends the db and schema if set
func (pb *ProcedureBuilder) QualifiedNameWithoutArguments() (string, error) {
	if pb.db == "" || pb.schema == "" || pb.name == "" {
		return "", errors.New("Procedures must specify a database a schema and a name")
	}
	return fmt.Sprintf(`"%v"."%v"."%v"`, pb.db, pb.schema, pb.name), nil
}

// Returns the arguments signature of the procedure in a form <PROCEDURE>(<TYPE>, <TYPE>, ..)
func (pb *ProcedureBuilder) ArgumentsSignature() (string, error) {
	return fmt.Sprintf(`%v(%v)`, strings.ToUpper(pb.name), strings.ToUpper(strings.Join(pb.argumentTypes, ", "))), nil
}

// WithArgs sets the args and argumentTypes on the ProcedureBuilder
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

// WithReturnBehavior
func (pb *ProcedureBuilder) WithReturnBehavior(s string) *ProcedureBuilder {
	pb.returnBehavior = s
	return pb
}

// WithNullInputBehavior
func (pb *ProcedureBuilder) WithNullInputBehavior(s string) *ProcedureBuilder {
	pb.nullInputBehavior = s
	return pb
}

// WithReturnType adds the data type of the return type to the ProcedureBuilder
func (pb *ProcedureBuilder) WithReturnType(s string) *ProcedureBuilder {
	pb.returnType = strings.ToUpper(s)
	return pb
}

// WithExecuteAs sets the execute to OWNER or CALLER
func (pb *ProcedureBuilder) WithExecuteAs(s string) *ProcedureBuilder {
	pb.executeAs = s
	return pb
}

// WithComment adds a comment to the ProcedureBuilder
func (pb *ProcedureBuilder) WithComment(c string) *ProcedureBuilder {
	pb.comment = c
	return pb
}

// WithStatement adds the SQL statement to be used for the procedure
func (pb *ProcedureBuilder) WithStatement(s string) *ProcedureBuilder {
	pb.statement = s
	return pb
}

// Returns the argument types
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
func Procedure(db, schema, name string, argTypes []string) *ProcedureBuilder {
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
	q.WriteString(" LANGUAGE javascript")
	if pb.nullInputBehavior != "" {
		q.WriteString(fmt.Sprintf(` %v`, EscapeString(pb.nullInputBehavior)))
	}
	if pb.returnBehavior != "" {
		q.WriteString(fmt.Sprintf(` %v`, EscapeString(pb.returnBehavior)))
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
func (vb *ProcedureBuilder) ChangeComment(c string) (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`ALTER PROCEDURE %v SET COMMENT = '%v'`, qn, EscapeString(c)), nil
}

// RemoveComment returns the SQL query that will remove the comment on the procedure.
func (vb *ProcedureBuilder) RemoveComment() (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER PROCEDURE %v UNSET COMMENT`, qn), nil
}

// ChangeExecuteAs returns the SQL query that will update the call mode on the procedure.
func (vb *ProcedureBuilder) ChangeExecuteAs(c string) (string, error) {
	qn, err := vb.QualifiedName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`ALTER PROCEDURE %v EXECUTE AS %v`, qn, c), nil
}

// Show returns the SQL query that will show the row representing this procedure.
// This show statement returns all procedures with the given name (overloaded ones)
func (pb *ProcedureBuilder) Show() string {
	return fmt.Sprintf(`SHOW PROCEDURES LIKE '%v' IN SCHEMA "%v"."%v"`, pb.name, pb.db, pb.schema)
}

// To describe the procedure the name must be specified as fully qualified name
// including argument types
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

type procedure struct {
	Comment sql.NullString `db:"description"`
	// Snowflake returns is_secure in the show procedure output, but it is irrelevant
	Name         sql.NullString `db:"name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Text         sql.NullString `db:"text"`
	DatabaseName sql.NullString `db:"catalog_name"`
	Arguments    sql.NullString `db:"arguments"`
}

type procedureDescription struct {
	Property sql.NullString `db:"property"`
	Value    sql.NullString `db:"value"`
}

// ScanProcedureDescription reads through the rows with property and value columns
// and returns a slice of procedureDescription structs
func ScanProcedureDescription(rows *sqlx.Rows) ([]procedureDescription, error) {
	pdsl := []procedureDescription{}
	for rows.Next() {
		pd := procedureDescription{}
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
func ScanProcedures(rows *sqlx.Rows) ([]*procedure, error) {
	var pcs []*procedure
	for rows.Next() {
		r := &procedure{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, r)
	}
	return pcs, rows.Err()
}

func ListProcedures(databaseName string, schemaName string, db *sql.DB) ([]procedure, error) {
	stmt := fmt.Sprintf(`SHOW PROCEDURES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []procedure{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no procedures found")
		return nil, nil
	}
	return dbs, pe.Wrapf(err, "unable to scan row for %s", stmt)
}
