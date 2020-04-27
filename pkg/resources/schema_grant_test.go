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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaGrant(t *testing.T) {
	r := require.New(t)
	err := resources.SchemaGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestSchemaGrantCreate(t *testing.T) {
	r := require.New(t)

	for _, test_priv := range []string{"USAGE", "MODIFY"} {
		in := map[string]interface{}{
			"schema_name":   "test-schema",
			"database_name": "test-db",
			"privilege":     test_priv,
			"roles":         []interface{}{"test-role-1", "test-role-2"},
			"shares":        []interface{}{"test-share-1", "test-share-2"},
		}
		d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Schema, in)
		r.NotNil(d)

		WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-1"$`, test_priv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO ROLE "test-role-2"$`, test_priv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-1"$`, test_priv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(
				fmt.Sprintf(`^GRANT %s ON SCHEMA "test-db"."test-schema" TO SHARE "test-share-2"$`, test_priv),
			).WillReturnResult(sqlmock.NewResult(1, 1))
			expectReadSchemaGrant(mock, test_priv)
			err := resources.CreateSchemaGrant(d, db)
			r.NoError(err)
		})
	}
}

func expectReadSchemaGrant(mock sqlmock.Sqlmock, test_priv string) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "SCHEMA", "test-schema", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "SCHEMA", "test-schema", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "SCHEMA", "test-schema", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), test_priv, "SCHEMA", "test-schema", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON SCHEMA "test-db"."test-schema"$`).WillReturnRows(rows)
}

func TestFutureSchemaGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":     true,
		"database_name": "test-db",
		"privilege":     "USAGE",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.SchemaGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT USAGE ON FUTURE SCHEMAS IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureSchemaGrant(mock)
		err := resources.CreateSchemaGrant(d, db)
		r.NoError(err)
	})
}

func expectReadFutureSchemaGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "SCHEMA", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "SCHEMA", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}
