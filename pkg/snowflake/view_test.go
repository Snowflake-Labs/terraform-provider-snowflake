package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestView(t *testing.T) {
	r := require.New(t)
	v := View("test")
	r.NotNil(v)
	r.False(v.secure)
	r.Equal(v.QualifiedName(), `"test"`)

	v.WithSecure()
	r.True(v.secure)

	v.WithComment("great comment")
	r.Equal("great comment", v.comment)

	v.WithStatement("SELECT * FROM DUMMY LIMIT 1")
	r.Equal("SELECT * FROM DUMMY LIMIT 1", v.statement)

	v.WithStatement("SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1")

	q := v.Create()
	r.Equal(`CREATE SECURE VIEW "test" COMMENT = 'great comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q = v.Secure()
	r.Equal(`ALTER VIEW "test" SET SECURE`, q)

	q = v.Unsecure()
	r.Equal(`ALTER VIEW "test" UNSET SECURE`, q)

	q = v.ChangeComment("bad comment")
	r.Equal(`ALTER VIEW "test" SET COMMENT = 'bad comment'`, q)

	q = v.RemoveComment()
	r.Equal(`ALTER VIEW "test" UNSET COMMENT`, q)

	q = v.Drop()
	r.Equal(`DROP VIEW "test"`, q)

	q = v.Show()
	r.Equal(`SHOW VIEWS LIKE 'test'`, q)

	v.WithDB("mydb")
	r.Equal(v.QualifiedName(), `"mydb".."test"`)

	q = v.Create()
	r.Equal(`CREATE SECURE VIEW "mydb".."test" COMMENT = 'great comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q = v.Secure()
	r.Equal(`ALTER VIEW "mydb".."test" SET SECURE`, q)

	q = v.Show()
	r.Equal(`SHOW VIEWS LIKE 'test' IN DATABASE "mydb"`, q)

	q = v.Drop()
	r.Equal(`DROP VIEW "mydb".."test"`, q)
}

func TestQualifiedName(t *testing.T) {
	r := require.New(t)
	v := View("view")
	r.Equal(v.QualifiedName(), `"view"`)

	v = View("view").WithDB("db")
	r.Equal(v.QualifiedName(), `"db".."view"`)

	v = View("view").WithSchema("schema")
	r.Equal(v.QualifiedName(), `"schema"."view"`)

	v = View("view").WithDB("db").WithSchema("schema")
	r.Equal(v.QualifiedName(), `"db"."schema"."view"`)
}

func TestRename(t *testing.T) {
	r := require.New(t)
	v := View("test")

	q := v.Rename("test2")
	r.Equal(`ALTER VIEW "test" RENAME TO "test2"`, q)

	v.WithDB("testDB")
	q = v.Rename("test3")
	r.Equal(`ALTER VIEW "testDB".."test2" RENAME TO "testDB".."test3"`, q)

	v = View("test4").WithSchema("testSchema")
	q = v.Rename("test5")
	r.Equal(`ALTER VIEW "testSchema"."test4" RENAME TO "testSchema"."test5"`, q)
}
