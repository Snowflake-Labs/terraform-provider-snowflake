package snowflake_test

import (
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	a := assert.New(t)
	db := snowflake.Database("db1")
	a.NotNil(db)

	q := db.Show()
	a.Equal("SHOW DATABASES LIKE 'db1'", q)

	q = db.Drop()
	a.Equal(`DROP DATABASE "db1"`, q)

	q = db.Rename("db2")
	a.Equal(`ALTER DATABASE "db1" RENAME TO "db2"`, q)

	ab := db.Alter()
	a.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	a.Equal(`ALTER DATABASE "db1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	a.Equal(`ALTER DATABASE "db1" SET FOO='bar' BAM=false`, q)

	c := db.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	a.Equal(`CREATE DATABASE "db1" FOO='bar' BAM=false`, q)
}

func TestDatabaseCreateFromShare(t *testing.T) {
	a := assert.New(t)
	db := snowflake.DatabaseFromShare("db1", "abc123", "share1")
	q := db.Create()
	a.Equal(`CREATE DATABASE "db1" FROM SHARE "abc123"."share1"`, q)
}

func TestDatabaseCreateFromDatabase(t *testing.T) {
	a := assert.New(t)
	db := snowflake.DatabaseFromDatabase("db1", "abc123")
	q := db.Create()
	a.Equal(`CREATE DATABASE "db1" CLONE "abc123"`, q)
}
