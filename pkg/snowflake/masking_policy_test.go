package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestMaskingPolicyCreate(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	m.WithValueDataType("string")
	m.WithMaskingExpression(`case
	when current_role() in ('ANALYST') then val
	else '*********'
end`)
	m.WithReturnDataType("string")
	m.WithComment("This is a test comment")

	q := m.Create()
	r.Equal(`CREATE MASKING POLICY "test_db"."test_schema"."test_masking_policy" AS (VAL string) RETURNS string -> case
	when current_role() in ('ANALYST') then val
	else '*********'
end COMMENT = 'This is a test comment'`, q)
}

func TestMaskingPolicyDescribe(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.Describe()
	r.Equal(`DESCRIBE MASKING POLICY "test_db"."test_schema"."test_masking_policy"`, q)
}

func TestMaskingPolicyDrop(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.Drop()
	r.Equal(`DROP MASKING POLICY "test_db"."test_schema"."test_masking_policy"`, q)
}

func TestMaskingPolicyChangeComment(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.ChangeComment("test comment!")
	r.Equal(`ALTER MASKING POLICY "test_db"."test_schema"."test_masking_policy" SET COMMENT = 'test comment!'`, q)
}

func TestMaskingPolicyRemoveComment(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.RemoveComment()
	r.Equal(`ALTER MASKING POLICY "test_db"."test_schema"."test_masking_policy" UNSET COMMENT`, q)
}

func TestMaskingChangeMaskingExpression(t *testing.T) {
	r := require.New(t)
	m := snowflake.MaskingPolicy("test_masking_policy", "test_db", "test_schema")
	r.NotNil(m)

	q := m.ChangeMaskingExpression(`case
	when current_role() in ('ANALYST') then val
    else sha2(val, 512)
end`)

	r.Equal(`ALTER MASKING POLICY "test_db"."test_schema"."test_masking_policy" SET BODY -> case
	when current_role() in ('ANALYST') then val
    else sha2(val, 512)
end`, q)
}
