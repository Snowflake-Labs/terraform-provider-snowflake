package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TagAssociationTest struct {
	Builder        *TagAssociationBuilder
	ExpectedCreate string
	ExpectedDrop   string
}

func TestTagAssociation(t *testing.T) {
	tests := []TagAssociationTest{
		{Builder: TagAssociation("test_db|test_schema|sensitive").WithObjectIdentifier(`"test_schema"."test_table"`).WithObjectType("TABLE").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "test_schema"."test_table" SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "test_schema"."test_table" UNSET TAG "test_db"."test_schema"."sensitive"`,
		},
		{Builder: TagAssociation("test_db|test_schema|sensitive").WithObjectIdentifier(`"test_schema"."test_table.important"`).WithObjectType("COLUMN").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "test_db"."test_schema"."test_table" ALTER COLUMN important SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "test_db"."test_schema"."test_table" ALTER COLUMN important UNSET TAG "test_db"."test_schema"."sensitive"`},
	}
	for _, testCase := range tests {
		r := require.New(t)
		r.Equal(testCase.ExpectedCreate, testCase.Builder.Create())
		r.Equal(testCase.ExpectedDrop, testCase.Builder.Drop())
	}
}

type TableColumnNameTest struct {
	Builder                               TagAssociationBuilder
	expectedTableName, expectedColumnName string
}

func TestTableColumnName(t *testing.T) {
	tests := []TableColumnNameTest{
		{TagAssociationBuilder{objectIdentifier: `"a"."b"."c"`, objectType: "table"}, "c", ""},
		{TagAssociationBuilder{objectIdentifier: `"db"."schema"."table.column"`, objectType: "column"}, "table", "column"},
	}
	for _, testCase := range tests {
		r := require.New(t)
		tableName, columnName := testCase.Builder.GetTableAndColumnName()
		r.Equal(testCase.expectedTableName, tableName)
		r.Equal(testCase.expectedColumnName, columnName)
	}
}
