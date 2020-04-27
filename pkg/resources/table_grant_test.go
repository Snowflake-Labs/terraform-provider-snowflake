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

func TestTableGrant(t *testing.T) {
	r := require.New(t)
	err := resources.TableGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"table_name":    "test-table",
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "SELECT",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON TABLE "test-db"."PUBLIC"."test-table" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadTableGrant(mock)
		err := resources.CreateTableGrant(d, db)
		r.NoError(err)
	})
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
		"on_future":     true,
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "SELECT",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.TableGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureTableGrant(mock)
		err := resources.CreateTableGrant(d, db)
		r.NoError(err)
	})

	b := require.New(t)

	in = map[string]interface{}{
		"on_future":     true,
		"database_name": "test-db",
		"privilege":     "SELECT",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
	}
	d = schema.TestResourceDataRaw(t, resources.TableGrant().Schema, in)
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
