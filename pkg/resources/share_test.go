package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
)

func TestShare(t *testing.T) {
	r := require.New(t)
	err := resources.Share().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestShareCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":     "test-share",
		"comment":  "great comment",
		"accounts": []interface{}{"bob123", "sue456"},
	}
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SHARE "test-share" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^CREATE DATABASE "TEMP_test-share_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT REFERENCE_USAGE ON DATABASE "TEMP_test-share_\d*" TO SHARE "test-share"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER SHARE "test-share" SET ACCOUNTS=bob123,sue456$`).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectBegin()
		mock.ExpectExec(`^REVOKE REFERENCE_USAGE ON DATABASE "TEMP_test-share_\d*" FROM SHARE "test-share"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec(`^REVOKE USAGE ON DATABASE "TEMP_test-share_\d*" FROM SHARE "test-share"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectExec(`^DROP DATABASE "TEMP_test-share_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadShare(mock)
		err := resources.CreateShare(d, db)
		r.NoError(err)
	})
}

func expectReadShare(mock sqlmock.Sqlmock) {
	// &createdOn, &kind, &name, &databaseName, &to, &owner, &comment
	rows := sqlmock.NewRows([]string{
		"created_on", "kind", "name", "database_name", "to", "owner", "comment",
	}).AddRow("2019-05-19 16:55:36.530 -0700", "SECURE", "test-share", "test_db", "bob123, sue456", "admin", "great comment")
	mock.ExpectQuery(`^SHOW SHARES LIKE 'test-share'$`).WillReturnRows(rows)
}

func TestStripAccountFromName(t *testing.T) {
	r := require.New(t)
	s := "yt12345.my_share"
	r.Equal("my_share", resources.StripAccountFromName(s))

	s = "yt12345.my.share"
	r.Equal("my.share", resources.StripAccountFromName(s))

	s = "no_account_for_some_reason"
	r.Equal("no_account_for_some_reason", resources.StripAccountFromName(s))
}

func TestShareRead(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "test-share",
	}
	d := share(t, "test-share", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Share(d.Id()).Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err := resources.ReadShare(d, db)
		r.Empty(d.State())
		r.Nil(err)
	})
}
