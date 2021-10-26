package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	r := require.New(t)
	err := resources.Tag().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTagCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.Tag().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TAG "test_db"."test_schema"."good_name" COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadTag(mock)
		err := resources.CreateTag(d, db)
		r.NoError(err)
	})
}

func TestTagUpdate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}

	d := tag(t, "test_db|test_schema|good_name", in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^ALTER TAG "test_db"."test_schema"."good_name" SET COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadTag(mock)
		err := resources.UpdateTag(d, db)
		r.NoError(err)
	})
}

func TestTagDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
		"schema":   "test_schema",
		"comment":  "great comment",
	}

	d := tag(t, "test_db|test_schema|good_name", in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^DROP TAG "test_db"."test_schema"."good_name"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := resources.DeleteTag(d, db)
		r.NoError(err)
	})
}

func TestTagRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
		"schema":   "test_schema",
	}

	d := schema.TestResourceDataRaw(t, resources.Tag().Schema, in)
	d.SetId("test_db|test_schema|good_name")

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Tag("good_name").WithDB("test_db").WithSchema("test_schema").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadTag(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}

func expectReadTag(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "owner", "comment"},
	).AddRow("2019-05-19 16:55:36.530 -0700", "good_name", "test_db", "test_schema", "admin", "great comment")
	mock.ExpectQuery(`^SHOW TAGS LIKE 'good_name' IN SCHEMA "test_db"."test_schema"$`).WillReturnRows(rows)
}
