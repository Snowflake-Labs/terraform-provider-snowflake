package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform/helper/schema"
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
		"statement": "SELECT * FROM GREAT_DB.GREAT_SCHEMA.GREAT_TABLE WHERE account_id = 'bobs-account-id'",
		"is_secure": true,
	}
	d := schema.TestResourceDataRaw(t, resources.View().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`USE DATABASE "test_db"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^CREATE SECURE VIEW "good_name" COMMENT = 'great comment' AS SELECT \* FROM GREAT_DB.GREAT_SCHEMA.GREAT_TABLE WHERE account_id = 'bobs-account-id'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		expectReadView(mock)
		err := resources.CreateView(d, db)
		a.NoError(err)
	})
}

func expectReadView(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "reserved", "database_name", "schema_name", "owner", "comment", "text", "is_secure"},
	).AddRow("2019-05-19 16:55:36.530 -0700", "good_name", "", "GREAT_DB", "GREAT_SCHEMA", "admin", "great comment", "SELECT * FROM GREAT_DB.GREAT_SCHEMA.GREAT_TABLE WHERE account_id = 'bobs-account-id'", true)
	mock.ExpectQuery(`^SHOW VIEWS LIKE 'good_name'$`).WillReturnRows(rows)
}
