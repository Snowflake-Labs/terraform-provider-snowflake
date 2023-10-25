package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TagAssociationTest struct {
	Builder        *TagAssociationBuilder
	ExpectedCreate string
	ExpectedDrop   string
	ExpectedShow   string
}

func TestTagAssociation(t *testing.T) {
	tests := []TagAssociationTest{
		{
			Builder:        NewTagAssociationBuilder("test_db|test_schema|sensitive").WithObjectIdentifier(`"test_db"."test_schema"."test_table"`).WithObjectType("TABLE").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "test_db"."test_schema"."test_table" SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "test_db"."test_schema"."test_table" UNSET TAG "test_db"."test_schema"."sensitive"`,
			ExpectedShow:   `SELECT SYSTEM$GET_TAG('"test_db"."test_schema"."sensitive"', '"test_db"."test_schema"."test_table"', 'TABLE') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`,
		},
		{
			Builder:        NewTagAssociationBuilder("test_db|test_schema|sensitive").WithObjectIdentifier(`"test_db"."test_schema"."test_table.important"`).WithObjectType("COLUMN").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "important" SET TAG "test_db"."test_schema"."sensitive" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "important" UNSET TAG "test_db"."test_schema"."sensitive"`,
			ExpectedShow:   `SELECT SYSTEM$GET_TAG('"test_db"."test_schema"."sensitive"', '"test_db"."test_schema"."test_table"."important"', 'COLUMN') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`,
		},
		{
			Builder:        NewTagAssociationBuilder("tag_db|tag_schema|sensitive").WithObjectIdentifier(`"table_db"."table_schema"."test_table.important"`).WithObjectType("COLUMN").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "table_db"."table_schema"."test_table" MODIFY COLUMN "important" SET TAG "tag_db"."tag_schema"."sensitive" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "table_db"."table_schema"."test_table" MODIFY COLUMN "important" UNSET TAG "tag_db"."tag_schema"."sensitive"`,
			ExpectedShow:   `SELECT SYSTEM$GET_TAG('"tag_db"."tag_schema"."sensitive"', '"table_db"."table_schema"."test_table"."important"', 'COLUMN') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`,
		},
		{
			Builder:        NewTagAssociationBuilder("OPERATION_DB|SECURITY|PII_2").WithObjectIdentifier(`"OPERATION_DB"."SECURITY"."test_table.important"`).WithObjectType("COLUMN").WithTagValue("true"),
			ExpectedCreate: `ALTER TABLE "OPERATION_DB"."SECURITY"."test_table" MODIFY COLUMN "important" SET TAG "OPERATION_DB"."SECURITY"."PII_2" = 'true'`,
			ExpectedDrop:   `ALTER TABLE "OPERATION_DB"."SECURITY"."test_table" MODIFY COLUMN "important" UNSET TAG "OPERATION_DB"."SECURITY"."PII_2"`,
			ExpectedShow:   `SELECT SYSTEM$GET_TAG('"OPERATION_DB"."SECURITY"."PII_2"', '"OPERATION_DB"."SECURITY"."test_table"."important"', 'COLUMN') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`,
		},
	}
	for _, testCase := range tests {
		r := require.New(t)
		r.Equal(testCase.ExpectedCreate, testCase.Builder.Create())
		r.Equal(testCase.ExpectedDrop, testCase.Builder.Drop())
		r.Equal(testCase.ExpectedShow, testCase.Builder.Show())
	}
}

type TableColumnNameTest struct {
	Builder                               *TagAssociationBuilder
	expectedTableName, expectedColumnName string
}

func TestTableColumnName(t *testing.T) {
	tests := []TableColumnNameTest{
		{NewTagAssociationBuilder("a|b|sensitive").WithObjectIdentifier(`"a"."b"."c"`).WithObjectType("TABLE"), `"a"."b"."c"`, ""},
		{NewTagAssociationBuilder("db|schema|sensitive").WithObjectIdentifier(`"db"."schema"."table.column"`).WithObjectType("COLUMN"), `"db"."schema"."table"`, "column"},
	}
	for _, testCase := range tests {
		r := require.New(t)
		tableName, columnName := testCase.Builder.GetTableAndColumnName()
		r.Equal(testCase.expectedTableName, tableName)
		r.Equal(testCase.expectedColumnName, columnName)
	}
}
