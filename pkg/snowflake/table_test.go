package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	r := require.New(t)
	s := TableWithColumnDefinition("test_table", "test_db", "test_schema", map[string]string{"column1": "OBJECT", "column2": "VARCHAR(255)"})
	r.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_table"`)

	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT, "column2" VARCHAR(255))`)

	s.WithComment("Test Comment")
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT, "column2" VARCHAR(255)) COMMENT = 'Test Comment'`)
}

func TestTableChangeComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeComment("new table comment"), `ALTER TABLE "test_db"."test_schema"."test_table" SET COMMENT = 'new table comment'`)
}

func TestTableRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.RemoveComment(), `ALTER TABLE "test_db"."test_schema"."test_table" UNSET COMMENT`)
}

func TestTableDrop(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP TABLE "test_db"."test_schema"."test_table"`)
}

func TestTableShow(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW TABLES LIKE 'test_table' IN DATABASE "test_db"`)
}
