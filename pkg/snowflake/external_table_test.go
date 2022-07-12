package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternalTableCreate(t *testing.T) {
	r := require.New(t)
	s := ExternalTable("test_table", "test_db", "test_schema")
	s.WithColumns([]map[string]string{{"name": "column1", "type": "OBJECT", "as": "expression1"}, {"name": "column2", "type": "VARCHAR", "as": "expression2"}})
	s.WithLocation("location")
	s.WithPattern("pattern")
	s.WithFileFormat("file format")
	r.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_table"`)

	r.Equal(s.Create(), `CREATE EXTERNAL TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT AS expression1, "column2" VARCHAR AS expression2) WITH LOCATION = location REFRESH_ON_CREATE = false AUTO_REFRESH = false PATTERN = 'pattern' FILE_FORMAT = ( file format )`)

	s.WithComment("Test Comment")
	r.Equal(s.Create(), `CREATE EXTERNAL TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT AS expression1, "column2" VARCHAR AS expression2) WITH LOCATION = location REFRESH_ON_CREATE = false AUTO_REFRESH = false PATTERN = 'pattern' FILE_FORMAT = ( file format ) COMMENT = 'Test Comment'`)
}

func TestExternalTableUpdate(t *testing.T) {
	r := require.New(t)
	s := ExternalTable("test_table", "test_db", "test_schema")
	s.WithTags([]TagValue{{Name: "tag1", Value: "value1", Schema: "test_schema", Database: "test_db"}})
	r.Equal(s.Update(), `ALTER EXTERNAL TABLE "test_db"."test_schema"."test_table" TAG "test_db"."test_schema"."tag1" = "value1"`)
}

func TestExternalTableDrop(t *testing.T) {
	r := require.New(t)
	s := ExternalTable("test_table", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP EXTERNAL TABLE "test_db"."test_schema"."test_table"`)
}

func TestExternalTableShow(t *testing.T) {
	r := require.New(t)
	s := ExternalTable("test_table", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW EXTERNAL TABLES LIKE 'test_table' IN SCHEMA "test_db"."test_schema"`)
}
