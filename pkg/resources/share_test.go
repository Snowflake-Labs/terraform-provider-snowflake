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

func TestShare(t *testing.T) {
	r := require.New(t)
	err := resources.Share().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestShareCreate(t *testing.T) {
	a := assert.New(t)

	in := map[string]interface{}{
		"name":     "test-share",
		"comment":  "great comment",
		"accounts": []string{"bob123", "sue456"},
	}
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SHARE "test-share" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER SHARE "test-share" SET ACCOUNTS=bob123,sue456$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadShare(mock)
		err := resources.CreateShare(d, db)
		a.NoError(err)
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
	a := assert.New(t)
	s := "yt12345.my_share"
	a.Equal("my_share", resources.StripAccountFromName(s))

	s = "yt12345.my.share"
	a.Equal("my.share", resources.StripAccountFromName(s))

	s = "no_account_for_some_reason"
	a.Equal("no_account_for_some_reason", resources.StripAccountFromName(s))
}
