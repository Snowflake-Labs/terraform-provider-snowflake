package snowflake

import (
	"fmt"
	"strings"
)

// MaskingPolicyBuilder abstracts the creation of SQL queries for a Snowflake Masking Policy.
type MaskingPolicyBuilder struct {
	name   string
	db     string
	schema string
}

// QualifiedName prepends the db and schema if set and escapes everything nicely.
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

// MaskingPolicy returns a pointer to a Builder that abstracts the DDL operations for a masking policy.
//
// Supported DDL operations are:
//   - CREATE MASKING POLICY
//   - ALTER MASKING POLICY
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
