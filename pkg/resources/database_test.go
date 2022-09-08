package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
)

func TestDatabase(t *testing.T) {
	r := require.New(t)
	err := resources.Database().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestDatabaseCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "tst-terraform-good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "tst-terraform-good_name" COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabase_Create_WithValidReplicationConfiguration(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "tst-terraform-good_name",
		"comment": "great comment",
		"replication_configuration": []interface{}{map[string]interface{}{
			"accounts":             []interface{}{"account1", "account2"},
			"ignore_edition_check": "true",
		}},
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "tst-terraform-good_name" COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`ALTER DATABASE "tst-terraform-good_name" ENABLE REPLICATION TO ACCOUNTS account1, account2`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)

		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabase_Create_WithReplicationConfig_AndFalseIgnoreEditionCheck(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "tst-terraform-good_name",
		"comment": "great comment",
		"replication_configuration": []interface{}{map[string]interface{}{
			"accounts":             []interface{}{"acc_to_replicate"},
			"ignore_edition_check": "false",
		}},
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		err := resources.CreateDatabase(d, db)
		r.EqualError(err, "error enabling replication - ignore edition check was set to false")
	})
}

func expectRead(mock sqlmock.Sqlmock) {
	dbRows := sqlmock.NewRows([]string{"created_on", "name", "is_default", "is_current", "origin", "owner", "comment", "options", "retention_time"}).AddRow("created_on", "tst-terraform-good_name", "is_default", "is_current", "origin", "owner", "mock comment", "options", "1")
	mock.ExpectQuery("SHOW DATABASES LIKE 'tst-terraform-good_name'").WillReturnRows(dbRows)
}

func TestDatabaseRead(t *testing.T) {
	r := require.New(t)

	d := database(t, "tst-terraform-good_name", map[string]interface{}{
		"name": "tst-terraform-good_name",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectRead(mock)
		err := resources.ReadDatabase(d, db)
		r.NoError(err)
		r.Equal("tst-terraform-good_name", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal(1, d.Get("data_retention_time_in_days").(int))
	})
}

func TestDatabaseDelete(t *testing.T) {
	r := require.New(t)

	d := database(t, "tst-terraform-good_name-drop_it", map[string]interface{}{"name": "tst-terraform-good_name-drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP DATABASE "tst-terraform-good_name-drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromShare(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "tst-terraform-good_name",
		"from_share": map[string]interface{}{
			"provider": "abc123",
			"share":    "my_share",
		},
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "tst-terraform-good_name" FROM SHARE "abc123"."my_share"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromDatabase(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":          "tst-terraform-good_name",
		"from_database": "abc123",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "tst-terraform-good_name" CLONE "abc123"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromReplica(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "tst-terraform-good_name",
		"from_replica": "abc123",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "tst-terraform-good_name" AS REPLICA OF "abc123"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateTransient(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "tst-terraform-good_name",
		"comment":      "great comment",
		"is_transient": true,
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE TRANSIENT DATABASE "tst-terraform-good_name" COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateTransientFromDatabase(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":          "tst-terraform-good_name",
		"from_database": "abc123",
		"is_transient":  true,
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE TRANSIENT DATABASE "tst-terraform-good_name" CLONE "abc123"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}
