package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestView(t *testing.T) {
	r := require.New(t)
	err := resources.View().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestViewCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":      "good_name",
		"database":  "test_db",
		"comment":   "great comment",
		"statement": "SELECT * FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'",
		"is_secure": true,
	}
	d := schema.TestResourceDataRaw(t, resources.View().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE SECURE VIEW "test_db"."PUBLIC"."good_name" COMMENT = 'great comment' AS SELECT \* FROM test_db.PUBLIC.GREAT_TABLE WHERE account_id = 'bobs-account-id'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadView(mock)
		err := resources.CreateView(d, db)
		a.NoError(err)
	})
}

func expectReadView(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "reserved", "database_name", "schema_name", "owner", "comment", "text", "is_secure", "is_materialized"},
	).AddRow("2019-05-19 16:55:36.530 -0700", "good_name", "", "test_db", "GREAT_SCHEMA", "admin", "great comment", "SELECT * FROM test_db.GREAT_SCHEMA.GREAT_TABLE WHERE account_id = 'bobs-account-id'", true, false)
	mock.ExpectQuery(`^SHOW VIEWS LIKE 'good_name' IN DATABASE "test_db"$`).WillReturnRows(rows)
}
