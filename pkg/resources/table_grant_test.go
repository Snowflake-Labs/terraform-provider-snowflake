package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestTableGrant(t *testing.T) {
	r := require.New(t)
	err := resources.TableGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"table_name":        "test-table",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"shares":            []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadTableGrant(mock)
		err := resources.CreateTableGrant(d, db)
		r.NoError(err)
	})
}

func TestTableGrantUpdate(t *testing.T) {
	r := require.New(t)

	// d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, in)
	d := tableGrant(t, "test-db|PUBLIC|test-table|SELECT||false", map[string]interface{}{
		"table_name":    "test-table",
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "SELECT",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	})
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-1"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-2"`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadTableGrant(mock)

		err := resources.UpdateTableGrant(d, db)
		r.NoError(err)
	})
}

func TestTableGrantRead(t *testing.T) {
	r := require.New(t)

	d := tableGrant(t, "test-db|PUBLIC|test-table|SELECT|false", map[string]interface{}{
		"table_name":        "test-table",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadTableGrant(mock)
		err := resources.ReadTableGrant(d, db)
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

func expectReadTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-table", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-table", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-table", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-table", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON TABLE "test-db"."PUBLIC"."test-table"$`).WillReturnRows(rows)
}

func TestFutureTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureTableGrant(mock)
		err := resources.CreateTableGrant(d, db)
		roles := d.Get("roles").(*schema.Set)
		// After the CreateTableGrant has been created a ReadTableGrant reads the current grants
		// and this read should ignore test-role-3 what is returned by SHOW FUTURE GRANTS ON SCHEMA PUBLIC because
		// test-role-3 has been granted to a SELECT on future VIEW and not on future TABLE
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
	d = schema.TestResourceDataRaw(t, resources.TableGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureTableDatabaseGrant(mock)
		err := resources.CreateTableGrant(d, db)
		b.NoError(err)
	})
}

func expectReadFutureTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-db.PUBLIC.<TABLE>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-db.PUBLIC.<TABLE>", "ROLE", "test-role-2", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-db.PUBLIC.<VIEW>", "ROLE", "test-role-3", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureTableDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-db.<TABLE>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "TABLE", "test-db.<TABLE>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
