package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	pe "github.com/pkg/errors"
)

// RowAccessPolicyBuilder abstracts the creation of SQL queries for a Snowflake Row Access Policy
type RowAccessPolicyBuilder struct {
	name                string
	db                  string
	schema              string
	comment             string
	signature           map[string]interface{}
	rowAccessExpression string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (rapb *RowAccessPolicyBuilder) QualifiedName() string {
	var n strings.Builder

	if rapb.db != "" && rapb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, rapb.db, rapb.schema))
	}

	if rapb.db != "" && rapb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, rapb.db))
	}

	if rapb.db == "" && rapb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, rapb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, rapb.name))

	return n.String()
}

// WithComment adds a comment to the RowAccessPolicyBuilder
func (rapb *RowAccessPolicyBuilder) WithComment(c string) *RowAccessPolicyBuilder {
	rapb.comment = EscapeString(c)
	return rapb
}

// WithSignature adds signature to the RowAccessPolicyBuilder
func (rapb *RowAccessPolicyBuilder) WithSignature(signature map[string]interface{}) *RowAccessPolicyBuilder {
	rapb.signature = signature
	return rapb
}

// WithRowAccessExpression adds rowAccessExpression to the RowAccessPolicyBuilder
func (rapb *RowAccessPolicyBuilder) WithRowAccessExpression(rowAccessExpression string) *RowAccessPolicyBuilder {
	rapb.rowAccessExpression = rowAccessExpression
	return rapb
}

// RowAccessPolicy returns a pointer to a Builder that abstracts the DDL operations for a row access policy.
//
// Supported DDL operations are:
//   - CREATE ROW ACCESS POLICY
//	 - ALTER ROW ACCESS POLICY
//   - DROP ROW ACCESS POLICY
//   - SHOW ROW ACCESS POLICIES
//   - DESCRIBE ROW ACCESS POLICY
//
func RowAccessPolicy(name, db, schema string) *RowAccessPolicyBuilder {
	return &RowAccessPolicyBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL query that will create a row access policy.
func (rapb *RowAccessPolicyBuilder) Create() string {
	q := strings.Builder{}

	sortedSignature := []string{}
	for _, k := range sortInterfaceStrings(rapb.signature) {
		sortedSignature = append(sortedSignature, EscapeString(fmt.Sprintf(`%v %v`, k, rapb.signature[k])))
	}
	parsedSignature := fmt.Sprintf(`%v`, strings.Join(sortedSignature, ", "))

	q.WriteString(fmt.Sprintf(`CREATE ROW ACCESS POLICY %v AS (%v) RETURNS BOOLEAN -> `, rapb.QualifiedName(), parsedSignature))

	q.WriteString(rapb.rowAccessExpression)

	if rapb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(rapb.comment)))
	}

	return q.String()
}

// Describe returns the SQL query that will describe a row access policy
func (rapb *RowAccessPolicyBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE ROW ACCESS POLICY %v`, rapb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the row access policy.
func (rapb *RowAccessPolicyBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER ROW ACCESS POLICY %v SET COMMENT = '%v'`, rapb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the row access policy.
func (rapb *RowAccessPolicyBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER ROW ACCESS POLICY %v UNSET COMMENT`, rapb.QualifiedName())
}

// ChangeRowAccessExpression returns the SQL query that will update the row access expression on the row access policy.
func (rapb *RowAccessPolicyBuilder) ChangeRowAccessExpression(rowAccessExpression string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER ROW ACCESS POLICY %v SET BODY -> `, rapb.QualifiedName()))

	q.WriteString(rowAccessExpression)

	return q.String()
}

// Drop returns the SQL query that will drop a row access policy.
func (rapb *RowAccessPolicyBuilder) Drop() string {
	return fmt.Sprintf(`DROP ROW ACCESS POLICY %v`, rapb.QualifiedName())
}

// Show returns the SQL query that will show a row access policy.
func (rapb *RowAccessPolicyBuilder) Show() string {
	return fmt.Sprintf(`SHOW ROW ACCESS POLICIES LIKE '%v' IN SCHEMA "%v"."%v"`, rapb.name, rapb.db, rapb.schema)
}

type RowAccessPolicyStruct struct {
	CreatedOn    sql.NullString `db:"created_on"`
	Name         sql.NullString `db:"name"`
	DatabaseName sql.NullString `db:"database_name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Kind         sql.NullString `db:"kind"`
	Owner        sql.NullString `db:"owner"`
	Comment      sql.NullString `db:"comment"`
}

func ScanRowAccessPolicies(row *sqlx.Row) (*RowAccessPolicyStruct, error) {
	m := &RowAccessPolicyStruct{}
	err := row.StructScan(m)
	return m, err
}

func ListRowAccessPolicies(databaseName string, schemaName string, db *sql.DB) ([]RowAccessPolicyStruct, error) {
	stmt := fmt.Sprintf(`SHOW ROW ACCESS POLICIES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []RowAccessPolicyStruct{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no row access policies found")
		return nil, nil
	}
	return dbs, pe.Wrapf(err, "unable to scan row for %s", stmt)
}
