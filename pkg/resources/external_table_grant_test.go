package resources_test

import (
	"database/sql"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExternalTableGrant(t *testing.T) {
	r := require.New(t)
	err := resources.ExternalTableGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestExternalTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"external_table_name": "test-external-table",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "SELECT",
		"roles":               []interface{}{"test-role-1", "test-role-2"},
		"shares":              []interface{}{"test-share-1", "test-share-2"},
		"with_grant_option":   true,
	}
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-external-table" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-external-table" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-external-table" TO SHARE "test-share-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON EXTERNAL TABLE "test-db"."PUBLIC"."test-external-table" TO SHARE "test-share-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadExternalTableGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		r.NoError(err)
	})
}

func TestExternalTableGrantRead(t *testing.T) {
	r := require.New(t)

	d := externalTableGrant(t, "test-db|PUBLIC|test-external-table|SELECT||false", map[string]interface{}{
		"external_table_name": "test-external-table",
		"schema_name":         "PUBLIC",
		"database_name":       "test-db",
		"privilege":           "SELECT",
		"roles":               []interface{}{},
		"shares":              []interface{}{},
		"with_grant_option":   false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadExternalTableGrant(mock)
		err := resources.ReadExternalTableGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())

	shares := d.Get("shares").(*schema.Set)
	r.True(shares.Contains("test-share-1"))
	r.True(shares.Contains("test-share-2"))
	r.Equal(2, shares.Len())
}

func expectReadExternalTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-external-table", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-external-table", "ROLE", "test-role-2", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-external-table", "SHARE", "test-share-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL_TABLE", "test-external-table", "SHARE", "test-share-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON EXTERNAL TABLE "test-db"."PUBLIC"."test-external-table"$`).WillReturnRows(rows)
}

func TestFutureExternalTableGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureExternalTableGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		r.NoError(err)
	})

	b := require.New(t)

	in = map[string]interface{}{
		"on_future":         true,
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": false,
	}
	d = schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE EXTERNAL TABLES IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureExternalTableDatabaseGrant(mock)
		err := resources.CreateExternalTableGrant(d, db)
		b.NoError(err)
	})

	c := require.New(t)

	in = map[string]interface{}{
		"database_name":       "test-db",
		"external_table_name": "test-table",
		"privilege":           "SELECT",
		"roles":               []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option":   false,
	}
	d = schema.TestResourceDataRaw(t, resources.ExternalTableGrant().Resource.Schema, in)
	c.NotNil(d)
	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		err := resources.CreateExternalTableGrant(d, db)
		c.Error(err)
	})
}

func expectReadFutureExternalTableGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureExternalTableDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "EXTERNAL TABLE", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}

func TestParseExternalTableGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseExternalTableGrantID("test-db|PUBLIC|test-external-table|SELECT|false|role1,role2|share1,share2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-external-table", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseExternalTableGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseExternalTableGrantID("test-db❄️PUBLIC❄️test-external-table❄️SELECT❄️true❄️role1,role2❄️share1,share2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-external-table", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
	r.Equal(2, len(grantID.Shares))
	r.Equal("share1", grantID.Shares[0])
	r.Equal("share2", grantID.Shares[1])
}

func TestParseExternalTableGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseExternalTableGrantID("test-db|PUBLIC|test-external-table|SELECT|role1,role2|true")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-external-table", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(true, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
