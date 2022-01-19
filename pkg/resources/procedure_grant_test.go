package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestProcedureGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ProcedureGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestProcedureGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"procedure_name": "test-procedure",
		"arguments": []interface{}{map[string]interface{}{
			"name": "a",
			"type": "array",
		}, map[string]interface{}{
			"name": "b",
			"type": "string",
		}},
		"return_type":       "string",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"shares":            []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.ProcedureGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON PROCEDURE "test-db"."PUBLIC"."test-procedure"\(ARRAY, STRING\) TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON PROCEDURE "test-db"."PUBLIC"."test-procedure"\(ARRAY, STRING\) TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON PROCEDURE "test-db"."PUBLIC"."test-procedure"\(ARRAY, STRING\) TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON PROCEDURE "test-db"."PUBLIC"."test-procedure"\(ARRAY, STRING\) TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadProcedureGrant(mock)
		err := resources.CreateProcedureGrant(d, db)
		r.NoError(err)
	})
}

func TestProcedureGrantRead(t *testing.T) {
	r := require.New(t)

	d := procedureGrant(t, "test-db|PUBLIC|test-procedure(A ARRAY, B STRING):STRING|USAGE||false", map[string]interface{}{
		"procedure_name": "test-procedure",
		"arguments": []interface{}{map[string]interface{}{
			"name": "a",
			"type": "array",
		}, map[string]interface{}{
			"name": "b",
			"type": "string",
		}},
		"return_type":       "string",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadProcedureGrant(mock)
		err := resources.ReadProcedureGrant(d, db)
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

func expectReadProcedureGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.test-schema.\"test-procedure(A ARRAY, B STRING):STRING\"", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.test-schema.\"test-procedure(A ARRAY, B STRING):STRING\"", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.test-schema.\"test-procedure(A ARRAY, B STRING):STRING\"", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.test-schema.\"test-procedure(A ARRAY, B STRING):STRING\"", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON PROCEDURE "test-db"."PUBLIC"."test-procedure"\(ARRAY, STRING\)$`).WillReturnRows(rows)
}

func TestFutureProcedureGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.ProcedureGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE PROCEDURES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE PROCEDURES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureProcedureGrant(mock)
		err := resources.CreateProcedureGrant(d, db)
		r.NoError(err)
	})

	b := require.New(t)

	in = map[string]interface{}{
		"on_future":         true,
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	}
	d = schema.TestResourceDataRaw(t, resources.ProcedureGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE PROCEDURES IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE PROCEDURES IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureProcedureDatabaseGrant(mock)
		err := resources.CreateProcedureGrant(d, db)
		b.NoError(err)
	})
}

func expectReadFutureProcedureGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureProcedureDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "PROCEDURE", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
