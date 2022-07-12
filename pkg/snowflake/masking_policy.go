package snowflake

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// MaskingPolicyBuilder abstracts the creation of SQL queries for a Snowflake Masking Policy
type MaskingPolicyBuilder struct {
	name              string
	db                string
	schema            string
	comment           string
	valueDataType     string
	maskingExpression string
	returnDataType    string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely
func (mpb *MaskingPolicyBuilder) QualifiedName() string {
	var n strings.Builder

	if mpb.db != "" && mpb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, mpb.db, mpb.schema))
	}

	if mpb.db != "" && mpb.schema == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, mpb.db))
	}

	if mpb.db == "" && mpb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, mpb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, mpb.name))

	return n.String()
}

// WithComment adds a comment to the MaskingPolicyBuilder
func (mpb *MaskingPolicyBuilder) WithComment(c string) *MaskingPolicyBuilder {
	mpb.comment = EscapeString(c)
	return mpb
}

// WithValueDataType adds valueDataType to the MaskingPolicyBuilder
func (mpb *MaskingPolicyBuilder) WithValueDataType(dataType string) *MaskingPolicyBuilder {
	mpb.valueDataType = dataType
	return mpb
}

// WithMaskingExpression adds maskingExpression to the MaskingPolicyBuilder
func (mpb *MaskingPolicyBuilder) WithMaskingExpression(maskingExpression string) *MaskingPolicyBuilder {
	mpb.maskingExpression = maskingExpression
	return mpb
}

// WithReturnDataType adds returnDataType to the MaskingPolicyBuilder
func (mpb *MaskingPolicyBuilder) WithReturnDataType(dataType string) *MaskingPolicyBuilder {
	mpb.returnDataType = dataType
	return mpb
}

// MaskingPolicy returns a pointer to a Builder that abstracts the DDL operations for a masking policy.
//
// Supported DDL operations are:
//   - CREATE MASKING POLICY
//	 - ALTER MASKING POLICY
//   - DROP MASKING POLICY
//   - SHOW MASKING POLICIES
//   - DESCRIBE MASKING POLICY
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/security-column-ddm.html)
func MaskingPolicy(name, db, schema string) *MaskingPolicyBuilder {
	return &MaskingPolicyBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL query that will create a masking policy.
func (mpb *MaskingPolicyBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE MASKING POLICY %v AS (VAL %v) RETURNS %v -> `, mpb.QualifiedName(), mpb.valueDataType, mpb.returnDataType))

	q.WriteString(mpb.maskingExpression)

	if mpb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(mpb.comment)))
	}

	return q.String()
}

// Describe returns the SQL query that will describe a masking policy
func (mpb *MaskingPolicyBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE MASKING POLICY %v`, mpb.QualifiedName())
}

// ChangeComment returns the SQL query that will update the comment on the masking policy.
func (mpb *MaskingPolicyBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER MASKING POLICY %v SET COMMENT = '%v'`, mpb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the masking policy.
func (mpb *MaskingPolicyBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER MASKING POLICY %v UNSET COMMENT`, mpb.QualifiedName())
}

// ChangeMaskingExpression returns the SQL query that will update the masking expression on the masking policy.
func (mpb *MaskingPolicyBuilder) ChangeMaskingExpression(maskingExpression string) string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`ALTER MASKING POLICY %v SET BODY -> `, mpb.QualifiedName()))

	q.WriteString(maskingExpression)

	return q.String()
}

// Drop returns the SQL query that will drop a masking policy.
func (mpb *MaskingPolicyBuilder) Drop() string {
	return fmt.Sprintf(`DROP MASKING POLICY %v`, mpb.QualifiedName())
}

// Show returns the SQL query that will show a masking policy.
func (mpb *MaskingPolicyBuilder) Show() string {
	return fmt.Sprintf(`SHOW MASKING POLICIES LIKE '%v' IN SCHEMA "%v"."%v"`, mpb.name, mpb.db, mpb.schema)
}

type MaskingPolicyStruct struct {
	CreatedOn    sql.NullString `db:"created_on"`
	Name         sql.NullString `db:"name"`
	DatabaseName sql.NullString `db:"database_name"`
	SchemaName   sql.NullString `db:"schema_name"`
	Kind         sql.NullString `db:"kind"`
	Owner        sql.NullString `db:"owner"`
	Comment      sql.NullString `db:"comment"`
}

func ScanMaskingPolicies(row *sqlx.Row) (*MaskingPolicyStruct, error) {
	m := &MaskingPolicyStruct{}
	err := row.StructScan(m)
	return m, err
}

func ListMaskingPolicies(databaseName string, schemaName string, db *sql.DB) ([]MaskingPolicyStruct, error) {
	stmt := fmt.Sprintf(`SHOW MASKING POLICIES IN SCHEMA "%s"."%v"`, databaseName, schemaName)
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []MaskingPolicyStruct{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no masking policies found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
