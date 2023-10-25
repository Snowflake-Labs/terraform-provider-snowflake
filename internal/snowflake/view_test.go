// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

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

	vb := NewViewBuilder(view).WithDB(db).WithSchema(schema)
	r.NotNil(vb)
	r.False(vb.secure)
	qn, err := vb.QualifiedName()
	r.NoError(err)
	r.Equal(fmt.Sprintf(`"%v"."%v"."%v"`, db, schema, view), qn)

	vb.WithSecure()
	r.True(vb.secure)

	vb.WithComment("great' comment")
	vb.WithStatement("SELECT * FROM DUMMY LIMIT 1")
	r.Equal("SELECT * FROM DUMMY LIMIT 1", vb.statement)

	vb.WithCopyGrants()
	r.True(vb.copyGrants)

	vb.WithStatement("SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1")

	q, err := vb.Create()
	r.NoError(err)
	r.Equal(`CREATE SECURE VIEW "some_database"."some_schema"."test" COPY GRANTS COMMENT = 'great\' comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q, err = vb.Secure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" SET SECURE`, q)

	q, err = vb.Unsecure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" UNSET SECURE`, q)

	q, err = vb.ChangeComment("bad' comment")
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" SET COMMENT = 'bad\' comment'`, q)

	q, err = vb.RemoveComment()
	r.NoError(err)
	r.Equal(`ALTER VIEW "some_database"."some_schema"."test" UNSET COMMENT`, q)

	q, err = vb.Drop()
	r.NoError(err)
	r.Equal(`DROP VIEW "some_database"."some_schema"."test"`, q)

	q = vb.Show()
	r.Equal(`SHOW VIEWS LIKE 'test' IN SCHEMA "some_database"."some_schema"`, q)

	vb.WithDB("mydb")
	qn, err = vb.QualifiedName()
	r.NoError(err)
	r.Equal(`"mydb"."some_schema"."test"`, qn)

	q, err = vb.Create()
	r.NoError(err)
	r.Equal(`CREATE SECURE VIEW "mydb"."some_schema"."test" COPY GRANTS COMMENT = 'great\' comment' AS SELECT * FROM DUMMY WHERE blah = 'blahblah' LIMIT 1`, q)

	q, err = vb.Secure()
	r.NoError(err)
	r.Equal(`ALTER VIEW "mydb"."some_schema"."test" SET SECURE`, q)

	q = vb.Show()
	r.Equal(`SHOW VIEWS LIKE 'test' IN SCHEMA "mydb"."some_schema"`, q)

	q, err = vb.Drop()
	r.NoError(err)
	r.Equal(`DROP VIEW "mydb"."some_schema"."test"`, q)
}

func TestQualifiedName(t *testing.T) {
	r := require.New(t)
	v := NewViewBuilder("view").WithDB("db").WithSchema("schema")
	qn, err := v.QualifiedName()
	r.NoError(err)
	r.Equal(`"db"."schema"."view"`, qn)
}

func TestRename(t *testing.T) {
	r := require.New(t)
	v := NewViewBuilder("test").WithDB("db").WithSchema("schema")

	q, err := v.Rename("test2")
	r.NoError(err)
	r.Equal(`ALTER VIEW "db"."schema"."test" RENAME TO "db"."schema"."test2"`, q)

	v.WithDB("testDB")
	q, err = v.Rename("test3")
	r.NoError(err)
	r.Equal(`ALTER VIEW "testDB"."schema"."test2" RENAME TO "testDB"."schema"."test3"`, q)

	v = NewViewBuilder("test4").WithDB("db").WithSchema("testSchema")
	q, err = v.Rename("test5")
	r.NoError(err)
	r.Equal(`ALTER VIEW "db"."testSchema"."test4" RENAME TO "db"."testSchema"."test5"`, q)
}
