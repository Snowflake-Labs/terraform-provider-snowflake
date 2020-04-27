package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestView(t *testing.T) {
	a := require.New(t)
	v := View("test")
	a.NotNil(v)
	a.False(v.secure)
	a.Equal(v.QualifiedName(), `"test"`)

	v.WithSecure()
	a.True(v.secure)

	v.WithComment("great comment")
	a.Equal("great comment", v.comment)

	v.WithStatement("SELECT * FROM DUMMY LIMIT 1")
	a.Equal("SELECT * FROM DUMMY LIMIT 1", v.statement)

	v.WithStatement("SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1")

	q := v.Create()
	a.Equal(`CREATE SECURE VIEW "test" COMMENT = 'great comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q = v.Secure()
	a.Equal(`ALTER VIEW "test" SET SECURE`, q)

	q = v.Unsecure()
	a.Equal(`ALTER VIEW "test" UNSET SECURE`, q)

	q = v.ChangeComment("bad comment")
	a.Equal(`ALTER VIEW "test" SET COMMENT = 'bad comment'`, q)

	q = v.RemoveComment()
	a.Equal(`ALTER VIEW "test" UNSET COMMENT`, q)

	q = v.Drop()
	a.Equal(`DROP VIEW "test"`, q)

	q = v.Show()
	a.Equal(`SHOW VIEWS LIKE 'test'`, q)

	v.WithDB("mydb")
	a.Equal(v.QualifiedName(), `"mydb".."test"`)

	q = v.Create()
	a.Equal(`CREATE SECURE VIEW "mydb".."test" COMMENT = 'great comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q = v.Secure()
	a.Equal(`ALTER VIEW "mydb".."test" SET SECURE`, q)

	q = v.Show()
	a.Equal(`SHOW VIEWS LIKE 'test' IN DATABASE "mydb"`, q)

	q = v.Drop()
	a.Equal(`DROP VIEW "mydb".."test"`, q)
}

func TestQualifiedName(t *testing.T) {
	a := require.New(t)
	v := View("view")
	a.Equal(v.QualifiedName(), `"view"`)

	v = View("view").WithDB("db")
	a.Equal(v.QualifiedName(), `"db".."view"`)

	v = View("view").WithSchema("schema")
	a.Equal(v.QualifiedName(), `"schema"."view"`)

	v = View("view").WithDB("db").WithSchema("schema")
	a.Equal(v.QualifiedName(), `"db"."schema"."view"`)
}

func TestRename(t *testing.T) {
	a := require.New(t)
	v := View("test")

	q := v.Rename("test2")
	a.Equal(`ALTER VIEW "test" RENAME TO "test2"`, q)

	v.WithDB("testDB")
	q = v.Rename("test3")
	a.Equal(`ALTER VIEW "testDB".."test2" RENAME TO "testDB".."test3"`, q)

	v = View("test4").WithSchema("testSchema")
	q = v.Rename("test5")
	a.Equal(`ALTER VIEW "testSchema"."test4" RENAME TO "testSchema"."test5"`, q)
}
