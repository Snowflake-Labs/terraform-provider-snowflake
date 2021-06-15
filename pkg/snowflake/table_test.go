package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	cols := []Column{
		{
			name:  "column1",
			_type: "OBJECT",
		},
		{
			name:  "column2",
			_type: "VARCHAR",
		},
	}

	s.WithColumns(Columns(cols))
	r.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_table"`)

	r.Equal(`CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT, "column2" VARCHAR)`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT, "column2" VARCHAR) COMMENT = 'Test Comment'`)

	s.WithClustering([]string{"column1"})
	r.Equal(s.Create(), `CREATE TABLE "test_db"."test_schema"."test_table" ("column1" OBJECT, "column2" VARCHAR) COMMENT = 'Test Comment' CLUSTER BY LINEAR(column1)`)

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

func TestTableAddColumn(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.AddColumn("new_column", "VARIANT"), `ALTER TABLE "test_db"."test_schema"."test_table" ADD COLUMN "new_column" VARIANT`)
}

func TestTableDropColumn(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropColumn("old_column"), `ALTER TABLE "test_db"."test_schema"."test_table" DROP COLUMN "old_column"`)
}

func TestTableChangeColumnType(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeColumnType("old_column", "BIGINT"), `ALTER TABLE "test_db"."test_schema"."test_table" MODIFY COLUMN "old_column" BIGINT`)
}

func TestTableChangeClusterBy(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.ChangeClusterBy("column2, column3"), `ALTER TABLE "test_db"."test_schema"."test_table" CLUSTER BY LINEAR(column2, column3)`)
}

func TestTableDropClusterBy(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.DropClustering(), `ALTER TABLE "test_db"."test_schema"."test_table" DROP CLUSTERING KEY`)
}

func TestTableDrop(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP TABLE "test_db"."test_schema"."test_table"`)
}

func TestTableShow(t *testing.T) {
	r := require.New(t)
	s := Table("test_table", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW TABLES LIKE 'test_table' IN SCHEMA "test_db"."test_schema"`)
}
