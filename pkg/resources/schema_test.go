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

func TestSchema(t *testing.T) {
	r := require.New(t)
	err := resources.Schema().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSchemaCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":         "good_name",
		"database":     "test_db",
		"comment":      "great comment",
		"is_transient": true,
		"is_managed":   true,
	}
	d := schema.TestResourceDataRaw(t, resources.Schema().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE TRANSIENT SCHEMA "test_db"."good_name" WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadSchema(mock)
		err := resources.CreateSchema(d, db)
		r.NoError(err)
	})
}

func TestSchemaRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "good_name",
		"database": "test_db",
	}

	d := schema.TestResourceDataRaw(t, resources.Schema().Schema, in)
	d.SetId("test_db|good_name")

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Schema("good_name").WithDB("test_db").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadSchema(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}

func expectReadSchema(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "is_default", "is_current", "database_name", "owner", "comment", "options", "retention_time"},
	).AddRow("2019-05-19 16:55:36.530 -0700", "good_name", "N", "Y", "test_db", "admin", "great comment", "TRANSIENT, MANAGED ACCESS", 1)
	mock.ExpectQuery(`^SHOW SCHEMAS LIKE 'good_name' IN DATABASE "test_db"$`).WillReturnRows(rows)
}
