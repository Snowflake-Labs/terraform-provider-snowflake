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

func TestDatabaseGrant(t *testing.T) {
	r := require.New(t)
	err := resources.DatabaseGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestDatabaseGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"database_name": "test-database",
		"privilege":     "USAGE",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.DatabaseGrant().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "test-database" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadDatabaseGrant(mock)
		err := resources.CreateDatabaseGrant(d, db)
		r.NoError(err)
	})
}

func TestDatabaseGrantRead(t *testing.T) {
	r := require.New(t)

	d := databaseGrant(t, "test-database|||IMPORTED PRIVILIGES", map[string]interface{}{
		"database_name": "test-database",
		"privilege":     "IMPORTED PRIVILIGES",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadDatabaseGrant(mock)
		err := resources.ReadDatabaseGrant(d, db)
		r.NoError(err)
	})
}

func expectReadDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "USAGE", "DATABASE", "test-database", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON DATABASE "test-database"$`).WillReturnRows(rows)
}
