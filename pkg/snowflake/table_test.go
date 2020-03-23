package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	a := assert.New(t)
	t := Table("test")
	a.NotNil(t)
	a.False(t.transient)
	a.Equal(t.QualifiedName(), `"test"`)

	t.WithComment("great comment")
	a.Equal("great comment", t.comment)

	t.WithColumns(map[string]string{"COL1": "VARCHAR", "COL2": "NUMBER"})
	a.Equal({"COL1": "VARCHAR", "COL2": "NUMBER"}, t.columns)

	q := t.Create()
	a.Equal(`CREATE TABLE "test" ("COL1" VARCHAR, "COL2" NUMBER) COMMENT = 'great comment'`, q)

	q = t.ChangeComment("bad comment")
	a.Equal(`ALTER TABLE "test" SET COMMENT = 'bad comment'`, q)

	q = t.RemoveComment()
	a.Equal(`ALTER TABLE "test" UNSET COMMENT`, q)

	q = t.Drop()
	a.Equal(`DROP TABLE "test"`, q)

	q = t.Show()
	a.Equal(`SHOW TABLES LIKE 'test'`, q)

	q = t.ShowColumns()
	a.Equal(`SHOW COLUMNS IN TABLE "test"`, q)

	v.WithDB("mydb")
	a.Equal(t.QualifiedName(), `"mydb".."test"`)

	q = t.Create()
	a.Equal(`CREATE TABLE "mydb".."test" ("COL1" VARCHAR, "COL2" NUMBER) COMMENT = 'great comment'`, q)

	q = t.Show()
	a.Equal(`SHOW TABLES LIKE 'test' IN DATABASE "mydb"`, q)

	q = t.ShowColumns()
	a.Equal(`SHOW COLUMNS IN TABLE "mydb".."test"`, q)

	q = t.Drop()
	a.Equal(`DROP TABLE "mydb".."test"`, q)
}

func TestQualifiedName(t *testing.T) {
	a := assert.New(t)
	t := Table("table")
	a.Equal(t.QualifiedName(), `"table"`)

	t = View("table").WithDB("db")
	a.Equal(t.QualifiedName(), `"db".."table"`)

	t = View("table").WithSchema("schema")
	a.Equal(t.QualifiedName(), `"schema"."table"`)

	t = View("table").WithDB("db").WithSchema("schema")
	a.Equal(t.QualifiedName(), `"db"."schema"."table"`)
}

func TestColumnsStatement(t *testing.T) {
	a := assert.New(t)
	t := Table("table")
	a.Equal(t.ColumnsStatement(), `()`)


	t := Table("table").WithColumns({"COL1": "VARCHAR", "COL2": "NUMBER"})
	a.Equal(t.ColumnsStatement(), `("COL1" VARCHAR, "COL2" NUMBER)`)

	t := Table("table").WithColumns({"COL1": "STRING", "COL2": "VARCHAR", "COL3": "DATE"})
	a.Equal(t.ColumnsStatement(), `("COL1" STRING, "COL2" VARCHAR, "COL3" DATE)`)
}

func TestRename(t *testing.T) {
	a := assert.New(t)
	t := Table("test")

	q := t.Rename("test2")
	a.Equal(`ALTER TABLE "test" RENAME TO "test2"`, q)

	t.WithDB("testDB")
	q = t.Rename("test3")
	a.Equal(`ALTER TABLE "testDB".."test2" RENAME TO "testDB".."test3"`, q)

	t = Table("test4").WithSchema("testSchema")
	q = v.Rename("test5")
	a.Equal(`ALTER TABÑE "testSchema"."test4" RENAME TO "testSchema"."test5"`, q)
}
