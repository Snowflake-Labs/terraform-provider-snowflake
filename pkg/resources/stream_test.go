package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	r := require.New(t)
	err := resources.Stream().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestStreamCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":              "stream_name",
		"database":          "database_name",
		"schema":            "schema_name",
		"comment":           "great comment",
		"on_table":          "target_db.target_schema.target_table",
		"append_only":       true,
		"show_initial_rows": true,
	}
	d := stream(t, "database_name|schema_name|stream_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE STREAM "database_name"."schema_name"."stream_name" ON TABLE "target_db"."target_schema"."target_table" COMMENT = 'great comment' APPEND_ONLY = true SHOW_INITIAL_ROWS = true`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		err := resources.CreateStream(d, db)
		r.NoError(err)
		r.Equal("stream_name", d.Get("name").(string))
	})
}

func expectStreamRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"name", "database_name", "schema_name", "owner", "comment", "table_name", "type", "stale", "mode"}).AddRow("stream_name", "database_name", "schema_name", "owner_name", "grand comment", "target_table", "DELTA", false, "APPEND_ONLY")
	mock.ExpectQuery(`SHOW STREAMS LIKE 'stream_name' IN DATABASE "database_name"`).WillReturnRows(rows)
}

func TestStreamRead(t *testing.T) {
	r := require.New(t)

	d := stream(t, "database_name|schema_name|stream_name", map[string]interface{}{"name": "stream_name", "comment": "grand comment"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectStreamRead(mock)
		err := resources.ReadStream(d, db)
		r.NoError(err)
		r.Equal("stream_name", d.Get("name").(string))
		r.Equal("grand comment", d.Get("comment").(string))

		// Test when resource is not found, checking if state will be empty
		r.NotEmpty(d.State())
		q := snowflake.Stream("stream_name", "database_name", "schema_name").Show()
		mock.ExpectQuery(q).WillReturnError(sql.ErrNoRows)
		err2 := resources.ReadStream(d, db)
		r.Empty(d.State())
		r.Nil(err2)
	})
}

func TestStreamDelete(t *testing.T) {
	r := require.New(t)

	d := stream(t, "database_name|schema_name|drop_it", map[string]interface{}{"name": "drop_it"})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP STREAM "database_name"."schema_name"."drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteStream(d, db)
		r.NoError(err)
	})
}

func TestStreamUpdate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":        "stream_name",
		"database":    "database_name",
		"schema":      "schema_name",
		"comment":     "new stream comment",
		"on_table":    "target_table",
		"append_only": true,
	}

	d := stream(t, "database_name|schema_name|stream_name", in)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER STREAM "database_name"."schema_name"."stream_name" SET COMMENT = 'new stream comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		err := resources.UpdateStream(d, db)
		r.NoError(err)
	})
}
