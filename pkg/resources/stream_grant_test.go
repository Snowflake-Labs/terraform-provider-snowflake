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

func TestStreamGrant(t *testing.T) {
	r := require.New(t)
	err := resources.StreamGrant().Resource.InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStreamGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"stream_name":       "test-stream",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`^GRANT SELECT ON STREAM "test-db"."PUBLIC"."test-stream" TO ROLE "test-role-1" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`^GRANT SELECT ON STREAM "test-db"."PUBLIC"."test-stream" TO ROLE "test-role-2" WITH GRANT OPTION$`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadStreamGrant(mock)
		err := resources.CreateStreamGrant(d, db)
		r.NoError(err)
	})
}

func TestStreamGrantRead(t *testing.T) {
	r := require.New(t)

	d := streamGrant(t, "test-db|PUBLIC|test-stream|SELECT||false", map[string]interface{}{
		"stream_name":       "test-stream",
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{},
		"with_grant_option": false,
	})

	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectReadStreamGrant(mock)
		err := resources.ReadStreamGrant(d, db)
		r.NoError(err)
	})

	roles := d.Get("roles").(*schema.Set)
	r.True(roles.Contains("test-role-1"))
	r.True(roles.Contains("test-role-2"))
	r.Equal(2, roles.Len())
}

func expectReadStreamGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "granted_on", "name", "granted_to", "grantee_name", "grant_option", "granted_by",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-stream", "ROLE", "test-role-1", false, "bob",
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-stream", "ROLE", "test-role-2", false, "bob",
	)
	mock.ExpectQuery(`^SHOW GRANTS ON STREAM "test-db"."PUBLIC"."test-stream"$`).WillReturnRows(rows)
}

func TestFutureStreamGrantCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"on_future":         true,
		"schema_name":       "PUBLIC",
		"database_name":     "test-db",
		"privilege":         "SELECT",
		"roles":             []interface{}{"test-role-1", "test-role-2"},
		"with_grant_option": true,
	}
	d := schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE STREAMS IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-1" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE STREAMS IN SCHEMA "test-db"."PUBLIC" TO ROLE "test-role-2" WITH GRANT OPTION$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureStreamGrant(mock)
		err := resources.CreateStreamGrant(d, db)
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
	d = schema.TestResourceDataRaw(t, resources.StreamGrant().Resource.Schema, in)
	b.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE STREAMS IN DATABASE "test-db" TO ROLE "test-role-1"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(
			`^GRANT SELECT ON FUTURE STREAMS IN DATABASE "test-db" TO ROLE "test-role-2"$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFutureStreamDatabaseGrant(mock)
		err := resources.CreateStreamGrant(d, db)
		b.NoError(err)
	})
}

func expectReadFutureStreamGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-db.PUBLIC.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN SCHEMA "test-db"."PUBLIC"$`).WillReturnRows(rows)
}

func expectReadFutureStreamDatabaseGrant(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "privilege", "grant_on", "name", "grant_to", "grantee_name", "grant_option",
	}).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-db.<SCHEMA>", "ROLE", "test-role-1", false,
	).AddRow(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "SELECT", "STREAM", "test-db.<SCHEMA>", "ROLE", "test-role-2", false,
	)
	mock.ExpectQuery(`^SHOW FUTURE GRANTS IN DATABASE "test-db"$`).WillReturnRows(rows)
}

func TestParseStreamGrantID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseStreamGrantID("test-db|PUBLIC|test-stream|SELECT|false|role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-stream", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseStreamGrantEmojiID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseStreamGrantID("test-db❄️PUBLIC❄️test-stream❄️SELECT❄️false❄️role1,role2")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-stream", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}

func TestParseStreamGrantOldID(t *testing.T) {
	r := require.New(t)

	grantID, err := resources.ParseStreamGrantID("test-db|PUBLIC|test-stream|SELECT|role1,role2|false")
	r.NoError(err)
	r.Equal("test-db", grantID.DatabaseName)
	r.Equal("PUBLIC", grantID.SchemaName)
	r.Equal("test-stream", grantID.ObjectName)
	r.Equal("SELECT", grantID.Privilege)
	r.Equal(false, grantID.WithGrantOption)
	r.Equal(2, len(grantID.Roles))
	r.Equal("role1", grantID.Roles[0])
	r.Equal("role2", grantID.Roles[1])
}
