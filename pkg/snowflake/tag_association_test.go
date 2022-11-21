package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TagAssociationTest struct {
	Builder        TagAssociationBuilder
	ExpectedCreate string
	ExpectedDrop   string
}

func TestTagAssociation(t *testing.T) {
	tests := []TagAssociationTest{
		{Builder: TagAssociationBuilder{
			"test_db", "test_schema.test_table", "TABLE", "test_schema", "sensitive", "true",
		}, ExpectedCreate: `ALTER TABLE test_schema.test_table SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop: `ALTER TABLE test_schema.test_table UNSET TAG "test_db"."test_schema"."sensitive"`,
		},
		{Builder: TagAssociationBuilder{
			"test_db", "test_schema.test_table.important", "COLUMN", "test_schema", "sensitive", "true",
		}, ExpectedCreate: `ALTER TABLE test_schema.test_table ALTER COLUMN important SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop: `ALTER TABLE test_schema.test_table ALTER COLUMN important UNSET TAG "test_db"."test_schema"."sensitive"`},
	}
	for _, testCase := range tests {
		r := require.New(t)
		r.Equal(testCase.Builder.Create(), testCase.ExpectedCreate)
		r.Equal(testCase.Builder.Drop(), testCase.ExpectedDrop)
	}
}
