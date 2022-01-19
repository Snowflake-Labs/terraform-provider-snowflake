package resources_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestStageGrant(t *testing.T) {
	r := require.New(t)
	err := resources.StageGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStageGrantCreate(t *testing.T) {
	r := require.New(t)

	for _, test_priv := range []string{"USAGE", "READ"} {
		in := map[string]interface{}{
			"stage_name":        "test-stage",
			"schema_name":       "test-schema",
			"database_name":     "test-db",
			"privilege":         test_priv,
			"roles":             []interface{}{"test-role-1", "test-role-2"},
			"shares":            []interface{}{"test-share-1", "test-share-2"},
			"with_grant_option": true,
		}
		d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, in)
		r.NotNil(d)

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO ROLE "test-role-1" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO ROLE "test-role-2" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO SHARE "test-share-1" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(fmt.Sprintf(`^GRANT %s ON STAGE "test-db"."test-schema"."test-stage" TO SHARE "test-share-2" WITH GRANT OPTION$`, test_priv)).WillReturnResult(sqlmock.NewResult(1, 1))
			expectReadStageGrant(mock, test_priv)
			err := resources.CreateStageGrant(d, db)
			r.NoError(err)
		})
	}
}

func TestStageGrantRead(t *testing.T) {
	r := require.New(t)

	d := stageGrant(t, "test-db|test-schema|test-stage|USAGE||false", map[string]interface{}{
		"stage_name":        "test-stage",
		"schema_name":       "test-schema",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{},
		"shares":            []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStageGrant(mock, "USAGE")
		err := resources.ReadStageGrant(d, db)
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

func expectReadStageGrant(mock sqlmock.Sqlmock, test_priv string) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "STAGE", "test-stage", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON STAGE "test-db"."test-schema"."test-stage"$`).WillReturnRows(rows)
}

func TestFutureStageGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "USAGE",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE STAGES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE STAGES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureStageGrant(mock)
		err := resources.CreateStageGrant(d, db)
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
	d = schema.TestResourceDataRaw(t, resources.StageGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE STAGES IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE STAGES IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureStageDatabaseGrant(mock)
		err := resources.CreateStageGrant(d, db)
		b.NoError(err)
	})
}

func expectReadFutureStageGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureStageDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "STAGE", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
