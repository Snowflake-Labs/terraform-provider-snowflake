package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestRowAccessPolicyCreate(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	m.WithSignature(map[string]interface{}{"n": "string", "v": "string"})
	m.WithRowAccessExpression(`case
	when current_role() in ('ANALYST') then true
	else false
end`)
	m.WithComment("This is a test comment")

	q := m.Create()
	r.Equal(`CREATE ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy" AS (n string, v string) RETURNS BOOLEAN -> case
	when current_role() in ('ANALYST') then true
	else false
end COMMENT = 'This is a test comment'`, q)
}

func TestRowAccessPolicyDescribe(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.Describe()
	r.Equal(`DESCRIBE ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy"`, q)
}

func TestRowAccessPolicyDrop(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.Drop()
	r.Equal(`DROP ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy"`, q)
}

func TestRowAccessPolicyChangeComment(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.ChangeComment("test comment!")
	r.Equal(`ALTER ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy" SET COMMENT = 'test comment!'`, q)
}

func TestRowAccessPolicyRemoveComment(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.RemoveComment()
	r.Equal(`ALTER ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy" UNSET COMMENT`, q)
}

func TestRowAccessChangeRowAccessExpression(t *testing.T) {
	r := require.New(t)
	m := snowflake.RowAccessPolicy("test_row_access_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.ChangeRowAccessExpression(`case
	when current_role() in ('ANALYST') then true
    else false
end`)

	r.Equal(`ALTER ROW ACCESS POLICY "test_db"."test_schema"."test_row_access_policy" SET BODY -> case
	when current_role() in ('ANALYST') then true
    else false
end`, q)
}
