package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestManagedAccount(t *testing.T) {
	r := require.New(t)
	err := resources.ManagedAccount().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestManagedAccountCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":           "test-account",
		"admin_name":     "bob",
		"admin_password": "abc123ABC",
		"comment":        "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.ManagedAccount().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE MANAGED ACCOUNT "test-account" ADMIN_NAME='bob' ADMIN_PASSWORD='abc123ABC' COMMENT='great comment' TYPE='READER'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadManagedAccount(mock)
		err := resources.CreateManagedAccount(d, db)
		r.NoError(err)
	})
}

func expectReadManagedAccount(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"name", "cloud", "region", "locator", "created_on", "url", "is_reader", "comment",
	}).AddRow("test-account", "aws", "ap-southeast-2", "locatorstring", "2019-01-01", "www.test.com", true, "great comment")
	mock.ExpectQuery(`^SHOW MANAGED ACCOUNTS LIKE 'test-account'$`).WillReturnRows(rows)
}
