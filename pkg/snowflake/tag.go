package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/jmoiron/sqlx"
)

// TagBuilder abstracts the creation of SQL queries for a Snowflake tag.
type TagBuilder struct {
	name                 string
	db                   string
	schema               string
	comment              string
	allowedValues        string
	maskingPolicyBuilder *MaskingPolicyBuilder
}

// QualifiedName prepends the db and schema if set and escapes everything nicely.
func (tb *TagBuilder) QualifiedName() string {
	var n strings.Builder

	if tb.db != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, tb.db))
	}

	if tb.schema != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, tb.schema))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, tb.name))

	return n.String()
}

// WithComment adds a comment to the TagBuilder.
func (tb *TagBuilder) WithComment(c string) *TagBuilder {
	tb.comment = c
	return tb
}

// WithDB adds the name of the database to the TagBuilder.
func (tb *TagBuilder) WithDB(db string) *TagBuilder {
	tb.db = db
	return tb
}

// WithSchema adds the name of the schema to the TagBuilder.
func (tb *TagBuilder) WithSchema(schema string) *TagBuilder {
	tb.schema = schema
	return tb
}

// WithAllowedValues adds the allowed values to the query.
func (tb *TagBuilder) WithAllowedValues(av []string) *TagBuilder {
	tb.allowedValues = helpers.ListToSnowflakeString(av)
	return tb
}

// WithMaskingPolicy adds a pointer to a MaskingPolicyBuilder to the TagBuilder.
func (tb *TagBuilder) WithMaskingPolicy(mpb *MaskingPolicyBuilder) *TagBuilder {
	tb.maskingPolicyBuilder = mpb
	return tb
}

// Tag returns a pointer to a Builder that abstracts the DDL operations for a tag.
//
// Supported DDL operations are:
//   - CREATE TAG
//   - ALTER TAG
//   - DROP TAG
//   - UNDROP TAG
//   - SHOW TAGS
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func NewTagBuilder(name string) *TagBuilder {
	return &TagBuilder{
		name: name,
	}
}

// Create returns the SQL query that will create a new tag.
func (tb *TagBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` TAG %v`, tb.QualifiedName()))

	if tb.allowedValues != "" {
		q.WriteString(fmt.Sprintf(` ALLOWED_VALUES %v`, tb.allowedValues))
	}

	if tb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(tb.comment)))
	}

	return q.String()
}

// Rename returns the SQL query that will rename the tag.
func (tb *TagBuilder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER TAG %v RENAME TO "%v"`, tb.QualifiedName(), newName)
}

// ChangeComment returns the SQL query that will update the comment on the tag.
func (tb *TagBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER TAG %v SET COMMENT = '%v'`, tb.QualifiedName(), EscapeString(c))
}

// RemoveComment returns the SQL query that will remove the comment on the tag.
func (tb *TagBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER TAG %v UNSET COMMENT`, tb.QualifiedName())
}

// AddAllowedValues returns the SQL query that will add the allowed_values.
func (tb *TagBuilder) AddAllowedValues(avs []string) string {
	return fmt.Sprintf(`ALTER TAG %v ADD ALLOWED_VALUES %v`, tb.QualifiedName(), helpers.ListToSnowflakeString(avs))
}

// DropAllowedValues returns the SQL query that will drop the unwanted allowed_values.
func (tb *TagBuilder) DropAllowedValues(davs []string) string {
	return fmt.Sprintf(`ALTER TAG %v DROP ALLOWED_VALUES %v`, tb.QualifiedName(), helpers.ListToSnowflakeString(davs))
}

// RemoveAllowedValues returns the SQL query that will remove the allowed_values from the tag.
func (tb *TagBuilder) RemoveAllowedValues() string {
	return fmt.Sprintf(`ALTER TAG %v UNSET ALLOWED_VALUES`, tb.QualifiedName())
}

// Drop returns the SQL query that will drop a tag.
func (tb *TagBuilder) Drop() string {
	return fmt.Sprintf(`DROP TAG %v`, tb.QualifiedName())
}

// Undrop returns the SQL query that will undrop a tag.
func (tb *TagBuilder) Undrop() string {
	return fmt.Sprintf(`UNDROP TAG %v`, tb.QualifiedName())
}

// AddMaskingPolicy returns the SQL query that will add a masking policy to a tag.
func (tb *TagBuilder) AddMaskingPolicy() string {
	return fmt.Sprintf(`ALTER TAG %v SET MASKING POLICY %v`, tb.QualifiedName(), tb.maskingPolicyBuilder.QualifiedName())
}

// RemoveMaskingPolicy returns the SQL query that will remove a masking policy from a tag.
func (tb *TagBuilder) RemoveMaskingPolicy() string {
	return fmt.Sprintf(`ALTER TAG %v UNSET MASKING POLICY %v`, tb.QualifiedName(), tb.maskingPolicyBuilder.QualifiedName())
}

// Show returns the SQL query that will show a tag.
func (tb *TagBuilder) Show() string {
	q := strings.Builder{}

	q.WriteString(fmt.Sprintf(`SHOW TAGS LIKE '%v'`, tb.name))

	if tb.schema != "" && tb.db != "" {
		q.WriteString(fmt.Sprintf(` IN SCHEMA "%v"."%v"`, tb.db, tb.schema))
	} else if tb.db != "" {
		q.WriteString(fmt.Sprintf(` IN DATABASE "%v"`, tb.db))
	}

	return q.String()
}

// Returns sql to show a tag with a specific policy attached to it.
func (tb *TagBuilder) ShowAttachedPolicy() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`SELECT * from table ("%v".information_schema.policy_references(ref_entity_name => '%v', ref_entity_domain => 'TAG')) where policy_db='%v' and policy_schema='%v' and policy_name='%v'`, tb.db, tb.QualifiedName(), tb.maskingPolicyBuilder.db, tb.maskingPolicyBuilder.schema, tb.maskingPolicyBuilder.name))

	return q.String()
}

type Tag struct {
	Name          sql.NullString `db:"name"`
	DatabaseName  sql.NullString `db:"database_name"`
	SchemaName    sql.NullString `db:"schema_name"`
	Comment       sql.NullString `db:"comment"`
	AllowedValues sql.NullString `db:"allowed_values"`
}

type TagPolicyAttachment struct {
	PolicyDB        sql.NullString `db:"POLICY_DB"`
	PolicySchema    sql.NullString `db:"POLICY_SCHEMA"`
	PolicyName      sql.NullString `db:"POLICY_NAME"`
	PolicyKind      sql.NullString `db:"POLICY_KIND"`
	RefDB           sql.NullString `db:"REF_DATABASE_NAME"`
	RefSchema       sql.NullString `db:"REF_SCHEMA_NAME"`
	RefEntity       sql.NullString `db:"REF_ENTITY_NAME"`
	RefEntityDomain sql.NullString `db:"REF_ENTITY_DOMAIN"`
}

type TagValue struct {
	Name     string
	Database string
	Schema   string
	Value    string
}

func ScanTag(row *sqlx.Row) (*Tag, error) {
	r := &Tag{}
	err := row.StructScan(r)
	return r, err
}

func ScanTagPolicy(row *sqlx.Row) (*TagPolicyAttachment, error) {
	r := &TagPolicyAttachment{}
	err := row.StructScan(r)
	return r, err
}
