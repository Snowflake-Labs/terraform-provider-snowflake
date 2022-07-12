package snowflake

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestView(t *testing.T) {
	r := require.New(t)
	db := "some_database"
	schema := "some_schema"
	view := "test"

	v := View(view).WithDB(db).WithSchema(schema)
	r.NotNil(v)
	r.False(v.secure)
	qn, err := v.QualifiedName()
	r.NoError(err)
	r.Equal(qn, fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, view))

	v.WithSecure()
	r.True(v.secure)

	v.WithComment("great' comment")
	v.WithStatement("SELECT * FROM DUMMY LIMIT 1")
	r.Equal("SELECT * FROM DUMMY LIMIT 1", v.statement)

	v.WithStatement("SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1")

	q, err := v.Create()
	r.NoError(err)
	r.Equal(`CREATE SECURE VIEW "some_database"."some_schema"."test" COMMENT = 'great\' comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q, err = v.Secure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" SET SECURE`, q)

	q, err = v.Unsecure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" UNSET SECURE`, q)

	q, err = v.ChangeComment("bad' comment")
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" SET COMMENT = 'bad\' comment'`, q)

	q, err = v.RemoveComment()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" UNSET COMMENT`, q)

	q, err = v.Drop()
	r.NoError(err)
	r.Equal(`DROP VIEW "some_database"."some_schema"."test"`, q)

	q = v.Show()
	r.Equal(`SHOW VIEWS LIKE 'test' IN SCHEMA "some_database"."some_schema"`, q)

	v.WithDB("mydb")
	qn, err = v.QualifiedName()
	r.NoError(err)
	r.Equal(qn, `"mydb"."some_schema"."test"`)

	q, err = v.Create()
	r.NoError(err)
	r.Equal(`CREATE SECURE VIEW "mydb"."some_schema"."test" COMMENT = 'great\' comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q, err = v.Secure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "mydb"."some_schema"."test" SET SECURE`, q)

	q = v.Show()
	r.Equal(`SHOW VIEWS LIKE 'test' IN SCHEMA "mydb"."some_schema"`, q)

	q, err = v.Drop()
	r.NoError(err)
	r.Equal(`DROP VIEW "mydb"."some_schema"."test"`, q)
}

func TestQualifiedName(t *testing.T) {
	r := require.New(t)
	v := View("view").WithDB("db").WithSchema("schema")
	qn, err := v.QualifiedName()
	r.NoError(err)
	r.Equal(qn, `"db"."schema"."view"`)
}

func TestRename(t *testing.T) {
	r := require.New(t)
	v := View("test").WithDB("db").WithSchema("schema")

	q, err := v.Rename("test2")
	r.NoError(err)
	r.Equal(`ALTER VIEW "db"."schema"."test" RENAME TO "db"."schema"."test2"`, q)

	v.WithDB("testDB")
	q, err = v.Rename("test3")
	r.NoError(err)
	r.Equal(`ALTER VIEW "testDB"."schema"."test2" RENAME TO "testDB"."schema"."test3"`, q)

	v = View("test4").WithDB("db").WithSchema("testSchema")
	q, err = v.Rename("test5")
	r.NoError(err)
	r.Equal(`ALTER VIEW "db"."testSchema"."test4" RENAME TO "db"."testSchema"."test5"`, q)
}
