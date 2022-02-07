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

func TestMaterializedViewGrant(t *testing.T) {
	r := require.New(t)
	err := resources.MaterializedViewGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestMaterializedViewGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"materialized_view_name": "test-materialized-view",
		"schema_name":            "PUBLIC",
		"database_name":          "test-db",
		"privilege":              "SELECT",
		"roles":                  []interface{}{"test-role-1", "test-role-2"},
		"shares":                 []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option":      true,
	}
	d := schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-materialized-view" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-materialized-view" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-materialized-view" TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-materialized-view" TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadMaterializedViewGrant(mock)
		err := resources.CreateMaterializedViewGrant(d, db)
		r.NoError(err)
	})
}

func TestMaterializedViewGrantRead(t *testing.T) {
	r := require.New(t)

	d := materializedViewGrant(t, "test-db|PUBLIC|test-materialized-view|SELECT||false", map[string]interface{}{
		"materialized_view_name": "test-materialized-view",
		"schema_name":            "PUBLIC",
		"database_name":          "test-db",
		"privilege":              "SELECT",
		"roles":                  []interface{}{},
		"shares":                 []interface{}{},
		"with_grant_option":      false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadMaterializedViewGrant(mock)
		err := resources.ReadMaterializedViewGrant(d, db)
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

func expectReadMaterializedViewGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-materialized-view", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-materialized-view", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-materialized-view", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-materialized-view", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON MATERIALIZED VIEW "test-db"."PUBLIC"."test-materialized-view"$`).WillReturnRows(rows)
}

func TestFutureMaterializedViewGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureMaterializedViewGrant(mock)
		err := resources.CreateMaterializedViewGrant(d, db)
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
	d = schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE MATERIALIZED VIEWS IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureMaterializedViewDatabaseGrant(mock)
		err := resources.CreateMaterializedViewGrant(d, db)
		b.NoError(err)
	})

	// Validate specifying on_future=false and schema_name="" generates an error
	m := require.New(t)

	in = map[string]interface{}{
		"on_future":         false,
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	}
	d = schema.TestResourceDataRaw(t, resources.MaterializedViewGrant().Resource.Schema, in)
	m.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		err := resources.CreateMaterializedViewGrant(d, db)
		m.Error(err)
	})
}

func expectReadFutureMaterializedViewGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-db.PUBLIC.<VIEW>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-db.PUBLIC.<VIEW>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureMaterializedViewDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-db.<VIEW>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "MATERIALIZED_VIEW", "test-db.<VIEW>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
