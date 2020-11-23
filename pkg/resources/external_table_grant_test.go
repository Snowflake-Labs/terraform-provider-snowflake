package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestExternalTableGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ExternalTableGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestExternalTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"external_table_name": "test-externalTable",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "SELECT",
		"roles":               []interface{}{"test-role-1", "test-role-2"},
		"shares":              []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option":   true,
	}
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-externalTable" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-externalTable" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-externalTable" TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-externalTable" TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadExternalTableGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		r.NoError(err)
	})
}
func TestExternalTableGrantRead(t *testing.T) {
	r := require.New(t)

	d := externalTableGrant(t, "test-db|PUBLIC|test-externalTable|SELECT|false", map[string]interface{}{
		"external_table_name": "test-externalTable",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "SELECT",
		"roles":               []interface{}{},
		"shares":              []interface{}{},
		"with_grant_option":   false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadExternalTableGrant(mock)
		err := resources.ReadExternalTableGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(roles.Len(), 2)

	shares := d.Get("shares").(*schema.Set)
	r.True(shares.Contains("test-share-1"))
	r.True(shares.Contains("test-share-2"))
	r.Equal(shares.Len(), 2)
}

func expectReadExternalTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-externalTable", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-externalTable", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-externalTable", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-externalTable", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON EXTERNAL TABLE "test-db"."PUBLIC"."test-externalTable"$`).WillReturnRows(rows)
}

func TestFutureExternalTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureExternalTableGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		roles := d.Get("roles").(*schema.Set)
		// After the CreateExternalTableGrant has been created a ReadExternalTableGrant reads the current grants
		// and this read should ignore test-role-3 what is returned by SHOW FUTURE GRANTS ON SCHEMA PUBLIC because
		// test-role-3 has been granted to a SELECT on future VIEW and not on future EXTERNAL TABLE
		r.True(roles.Contains("test-role-1"))
		r.True(roles.Contains("test-role-2"))
		r.False(roles.Contains("test-role-3"))
		r.NoError(err)
	})

	b := require.New(t)

	in = map[string]interface{}{
		"on_future":         true,
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	}
	d = schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureExternalTableDatabaseGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		b.NoError(err)
	})
}

func expectReadFutureExternalTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-db.PUBLIC.<EXTERNAL_TABLE>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-db.PUBLIC.<EXTERNAL_TABLE>", "ROLE", "test-role-2", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-db.PUBLIC.<VIEW>", "ROLE", "test-role-3", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureExternalTableDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-db.<EXTERNAL_TABLE>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-db.<EXTERNAL_TABLE>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
