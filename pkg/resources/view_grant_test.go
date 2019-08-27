package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestViewGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ViewGrant().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestViewGrantCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"view_name":     "test-view",
		"schema_name":   "PUBLIC",
		"database_name": "test-db",
		"privilege":     "SELECT",
		"roles":         []interface{}{"test-role-1", "test-role-2"},
		"shares":        []interface{}{"test-share-1", "test-share-2"},
	}
	d := schema.TestResourceDataRaw(t, resources.ViewGrant().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-view" TO ROLE "test-role-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-view" TO ROLE "test-role-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-view" TO SHARE "test-share-1"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON VIEW "test-db"."PUBLIC"."test-view" TO SHARE "test-share-2"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadViewGrant(mock)
		err := resources.CreateViewGrant(d, db)
		a.NoError(err)
	})
}

func expectReadViewGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "ROLE", "test-role-1", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "ROLE", "test-role-2", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "SHARE", "test-share-1", false, "bob").AddRow(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "VIEW", "test-view", "SHARE", "test-share-2", false, "bob")
	mock.ExpectQuery(`^SHOW GRANTS ON VIEW "test-db"."PUBLIC"."test-view"$`).WillReturnRows(rows)
}
