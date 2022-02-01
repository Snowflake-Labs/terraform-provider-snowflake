package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	r := require.New(t)
	err := resources.Database().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestDatabaseCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":    "good_name",
		"comment": "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "good_name" COMMENT='great comment`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func expectRead(mock sqlmock.Sqlmock) {
	dbRows := sqlmock.NewRows([]string{"created_on", "name", "is_default", "is_current", "origin", "owner", "comment", "options", "retention_time"}).AddRow("created_on", "good_name", "is_default", "is_current", "origin", "owner", "mock comment", "options", "1")
	mock.ExpectQuery("SHOW DATABASES LIKE 'good_name'").WillReturnRows(dbRows)
}

func TestDatabaseRead(t *testing.T) {
	r := require.New(t)

	d := database(t, "good_name", map[string]interface{}{
		"name": "good_name",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectRead(mock)
		err := resources.ReadDatabase(d, db)
		r.NoError(err)
		r.Equal("good_name", d.Get("name").(string))
		r.Equal("mock comment", d.Get("comment").(string))
		r.Equal(1, d.Get("data_retention_time_in_days").(int))
	})
}

func TestDatabaseDelete(t *testing.T) {
	r := require.New(t)

	d := database(t, "drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP DATABASE "drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromShare(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "good_name",
		"from_share": map[string]interface{}{
			"provider": "abc123",
			"share":    "my_share",
		},
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "good_name" FROM SHARE "abc123"."my_share"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromDatabase(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":          "good_name",
		"from_database": "abc123",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "good_name" CLONE "abc123"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}

func TestDatabaseCreateFromReplica(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "good_name",
		"from_replica": "abc123",
	}
	d := schema.TestResourceDataRaw(t, resources.Database().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE DATABASE "good_name" AS REPLICA OF "abc123"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectRead(mock)
		err := resources.CreateDatabase(d, db)
		r.NoError(err)
	})
}
