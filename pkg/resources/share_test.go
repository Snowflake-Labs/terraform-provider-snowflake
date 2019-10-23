package resources_test

import (
	"database/sql"
	"log"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		"accounts": []interface{}{"bob123", "sue456"},
	}
	d := schema.TestResourceDataRaw(t, resources.Share().Schema, in)
	a.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SHARE "test-share" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^CREATE DATABASE "TEMP_test-share_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "TEMP_test-share_\d*" TO SHARE "test-share"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER SHARE "test-share" SET ACCOUNTS=bob123,sue456$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^REVOKE USAGE ON DATABASE "TEMP_test-share_\d*" FROM SHARE "test-share"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^DROP DATABASE "TEMP_test-share_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadShare(mock)
		err := resources.CreateShare(d, db)
		a.NoError(err)
	})
}

func TestCaseInsensitiveShares(t *testing.T) {
	a := assert.New(t)

	in_a := map[string]interface{}{
		"name":     "test-share-case-insensitive",
		"comment":  "great comment",
		"accounts": []interface{}{"accountname"},
	}

	in_b := map[string]interface{}{
		"name":     "test-share-case-insensitive",
		"comment":  "great comment",
		"accounts": []interface{}{"ACCOUNTNAME"},
	}

	log.Printf("CREATING NEW SHARE INA")

	d_a := schema.TestResourceDataRaw(t, resources.Share().Schema, in_a)

	log.Printf("CREATING NEW SHARE INB")
	d_b := schema.TestResourceDataRaw(t, resources.Share().Schema, in_b)

	a.NotNil(d_a)
	a.NotNil(d_b)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^CREATE SHARE "test-share-case-insensitive" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^CREATE DATABASE "TEMP_test-share-case-insensitive_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT USAGE ON DATABASE "TEMP_test-share-case-insensitive_\d*" TO SHARE "test-share-case-insensitive"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^ALTER SHARE "test-share-case-insensitive" SET ACCOUNTS=accountname$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^REVOKE USAGE ON DATABASE "TEMP_test-share-case-insensitive_\d*" FROM SHARE "test-share-case-insensitive"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^DROP DATABASE "TEMP_test-share-case-insensitive_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadShareCaseInsensitive(mock)
		err := resources.CreateShare(d_a, db)
		a.NoError(err)
	})

	log.Printf("CREATING NEW SHARE")

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		// mock.ExpectExec(`^CREATE SHARE "test-share-case-insensitive" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec(`^CREATE DATABASE "TEMP_test-share-case-insensitive_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec(`^GRANT USAGE ON DATABASE "TEMP_test-share-case-insensitive_\d*" TO SHARE "test-share-case-insensitive"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec(`^ALTER SHARE "test-share-case-insensitive" SET ACCOUNTS=ACCOUNTNAME$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec(`^REVOKE USAGE ON DATABASE "TEMP_test-share-case-insensitive_\d*" FROM SHARE "test-share-case-insensitive"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec(`^DROP DATABASE "TEMP_test-share-case-insensitive_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		// expectReadShareCaseInsensitive(mock)
		// err := resources.CreateShare(d_b, db)
		mock.ExpectExec(`^SHOW SHARES LIKE "test-share-case-insensitive" COMMENT='great comment'$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^SHOW GRANTS ON DATABASE "TEMP_test-share-case-insensitive_\d*"$`).WillReturnResult(sqlmock.NewResult(1, 1))
		//expectReadShareCaseInsensitive(mock)
		err := resources.ReadShare(d_b, db)
		log.Printf("error %s", err)
	})
}

func expectReadShareCaseInsensitive(mock sqlmock.Sqlmock) {
	// &createdOn, &kind, &name, &databaseName, &to, &owner, &comment
	rows := sqlmock.NewRows([]string{
		"created_on", "kind", "name", "database_name", "to", "owner", "comment",
	}).AddRow("2019-10-08 16:55:36.530 -0700", "SECURE", "test-share-case-insensitive", "test_db", "ACCOUNTNAME", "admin", "great comment")
	mock.ExpectQuery(`^SHOW SHARES LIKE 'test-share-case-insensitive'$`).WillReturnRows(rows)
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
