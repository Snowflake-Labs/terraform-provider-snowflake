package snowflake_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("db1")
	r.NotNil(db)

	q := db.Show()
	r.Equal("SHOW DATABASES LIKE 'db1'", q)

	q = db.Drop()
	r.Equal(`DROP DATABASE "db1"`, q)

	q = db.Rename("db2")
	r.Equal(`ALTER DATABASE "db1" RENAME TO "db2"`, q)

	ab := db.Alter()
	r.NotNil(ab)

	ab.SetString(`foo`, `bar`)
	q = ab.Statement()

	r.Equal(`ALTER DATABASE "db1" SET FOO='bar'`, q)

	ab.SetBool(`bam`, false)
	q = ab.Statement()

	r.Equal(`ALTER DATABASE "db1" SET FOO='bar' BAM=false`, q)

	c := db.Create()
	c.SetString("foo", "bar")
	c.SetBool("bam", false)
	q = c.Statement()
	r.Equal(`CREATE DATABASE "db1" FOO='bar' BAM=false`, q)

	// test escaping
	c2 := db.Create()
	c2.SetString("foo", "ba'r")
	q = c2.Statement()
	r.Equal(`CREATE DATABASE "db1" FOO='ba\'r'`, q)
}

func TestDatabaseCreateFromShare(t *testing.T) {
	r := require.New(t)
	db := snowflake.DatabaseFromShare("db1", "abc123", "share1")
	q := db.Create()
	r.Equal(`CREATE DATABASE "db1" FROM SHARE "abc123"."share1"`, q)
}

func TestDatabaseCreateFromDatabase(t *testing.T) {
	r := require.New(t)
	db := snowflake.DatabaseFromDatabase("db1", "abc123")
	q := db.Create()
	r.Equal(`CREATE DATABASE "db1" CLONE "abc123"`, q)
}

func TestDatabaseCreateFromReplica(t *testing.T) {
	r := require.New(t)
	db := snowflake.DatabaseFromReplica("db1", "abc123")
	q := db.Create()
	r.Equal(`CREATE DATABASE "db1" AS REPLICA OF "abc123"`, q)
}

func TestListDatabases(t *testing.T) {
	r := require.New(t)
	mockDB, mock, err := sqlmock.New()
	r.NoError(err)
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{"created_on", "name", "is_default", "is_current", "origin", "owner", "comment", "options", "retention_time"}).AddRow("", "", "", "", "", "", "", "", "")
	mock.ExpectQuery(`SHOW DATABASES`).WillReturnRows(rows)
	_, err = snowflake.ListDatabases(sqlxDB)
	r.NoError(err)
}
