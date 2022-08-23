package snowflake_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
)

func TestQualifiedNameDatabase(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`"test"`, db.QualifiedName())
}

func TestCreateDatabase(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")

	r.Equal(`CREATE DATABASE "test"`, db.Create())

	db.Transient()
	r.Equal(`CREATE TRANSIENT DATABASE "test"`, db.Create())

	db.Clone("other")
	r.Equal(`CREATE TRANSIENT DATABASE "test" CLONE "other"`, db.Create())

	db.WithDataRetentionDays(7)
	r.Equal(`CREATE TRANSIENT DATABASE "test" CLONE "other" DATA_RETENTION_TIME_IN_DAYS = 7`, db.Create())

	db.WithComment("Yee'haw")
	r.Equal(`CREATE TRANSIENT DATABASE "test" CLONE "other" DATA_RETENTION_TIME_IN_DAYS = 7 COMMENT = 'Yee\'haw'`, db.Create())
}

func TestDatabaseCreateFromShare(t *testing.T) {
	r := require.New(t)
	db := snowflake.DatabaseFromShare("db1", "abc123", "share1")
	q := db.Create()
	r.Equal(`CREATE DATABASE "db1" FROM SHARE "abc123"."share1"`, q)

	db = snowflake.DatabaseFromShare("db1", "org1\".\"account1", "share1")
	db.WithComment("This is comment")
	q = db.Create()
	r.Equal(`CREATE DATABASE "db1" FROM SHARE "org1"."account1"."share1" COMMENT = 'This is comment'`, q)
}

func TestDatabaseCreateFromReplica(t *testing.T) {
	r := require.New(t)
	db := snowflake.DatabaseFromReplica("db1", "abc123")
	q := db.Create()
	r.Equal(`CREATE DATABASE "db1" AS REPLICA OF "abc123"`, q)
}

func TestDatabaseRename(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("db1")

	r.Equal(`ALTER DATABASE "db1" RENAME TO "db2"`, db.Rename("db2"))
}

func TestDatabaseSwap(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`ALTER DATABASE "test" SWAP WITH "target"`, db.Swap("target"))
}

func TestDatabaseChangeComment(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`ALTER DATABASE "test" SET COMMENT = 'test\' db'`, db.ChangeComment("test' db"))
}

func TestDatabaseRemoveComment(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`ALTER DATABASE "test" UNSET COMMENT`, db.RemoveComment())
}

func TestDatabaseChangeDataRetentionDays(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`ALTER DATABASE "test" SET DATA_RETENTION_TIME_IN_DAYS = 22`, db.ChangeDataRetentionDays(22))
}

func TestDatabaseRemoveDataRetentionDays(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`ALTER DATABASE "test" UNSET DATA_RETENTION_TIME_IN_DAYS`, db.RemoveDataRetentionDays())
}

func TestDatabaseDrop(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("db1")

	r.Equal(`DROP DATABASE "db1"`, db.Drop())
}

func TestDatabaseUndrop(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`UNDROP DATABASE "test"`, db.Undrop())
}

func TestDatabaseUse(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("test")
	r.Equal(`USE DATABASE "test"`, db.Use())
}

func TestDatabaseShow(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("db1")

	r.Equal("SHOW DATABASES LIKE 'db1'", db.Show())
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

func TestEnableReplicationAccounts(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("good_name")
	r.Equal(db.EnableReplicationAccounts("good_name", "account1"), `ALTER DATABASE "good_name" ENABLE REPLICATION TO ACCOUNTS account1`)
}

func TestDisableReplicationAccounts(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("good_name")
	r.Equal(db.DisableReplicationAccounts("good_name", "account1"), `ALTER DATABASE "good_name" DISABLE REPLICATION TO ACCOUNTS account1`)
}

func TestGetRemovedAccountsFromReplicationConfiguration(t *testing.T) {
	r := require.New(t)
	db := snowflake.Database("good_name")

	oldAccounts := []interface{}{"acc1", "acc2", "acc3"}
	newAccounts := []interface{}{"acc1", "acc2"}

	r.Equal(db.GetRemovedAccountsFromReplicationConfiguration(oldAccounts, newAccounts), []interface{}{"acc3"})
}
