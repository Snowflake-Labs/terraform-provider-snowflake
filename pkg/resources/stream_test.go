package resources_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	. "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

var streamColumns = []string{"name", "database_name", "schema_name", "owner", "comment", "table_name", "type", "stale", "mode", "source_type"}

func streamRow(mode string) []driver.Value {
	return []driver.Value{"stream_name", "database_name", "schema_name", "owner_name", "grand comment", "target_table", "DELTA", false, mode, "Table"}
}

func TestStream(t *testing.T) {
	err := resources.Stream().InternalValidate(provider.Provider().Schema, true)
	require.NoError(t, err)
}

func TestStreamCreate(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":              "stream_name",
		"database":          "database_name",
		"schema":            "schema_name",
		"comment":           "great comment",
		"on_table":          "target_db.target_schema.target_table",
		"append_only":       true,
		"insert_only":       false,
		"show_initial_rows": true,
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE STREAM "database_name"."schema_name"."stream_name" ON TABLE "target_db"."target_schema"."target_table" APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		expectOnTableRead(mock)
		err := resources.CreateStream(d, db)
		require.NoError(t, err)
		require.Equal(t, "stream_name", d.Get("name").(string))
	})
}

func TestStreamCreateOnExternalTable(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":              "stream_name",
		"database":          "database_name",
		"schema":            "schema_name",
		"comment":           "great comment",
		"on_table":          "target_db.target_schema.target_table",
		"append_only":       false,
		"insert_only":       true,
		"show_initial_rows": false,
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE STREAM "database_name"."schema_name"."stream_name" ON EXTERNAL TABLE "target_db"."target_schema"."target_table" INSERT_ONLY = true COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		expectOnExternalTableRead(mock)
		err := resources.CreateStream(d, db)
		require.NoError(t, err)
		require.Equal(t, "stream_name", d.Get("name").(string))
	})
}

func TestStreamCreateOnView(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":              "stream_name",
		"database":          "database_name",
		"schema":            "schema_name",
		"comment":           "great comment",
		"on_view":           "target_db.target_schema.target_view",
		"append_only":       true,
		"insert_only":       false,
		"show_initial_rows": true,
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`CREATE STREAM "database_name"."schema_name"."stream_name" ON VIEW "target_db"."target_schema"."target_view" APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'great comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		expectOnViewRead(mock)
		err := resources.CreateStream(d, db)
		require.NoError(t, err)
		require.Equal(t, "stream_name", d.Get("name").(string))
	})
}

func TestStreamOnMultipleSource(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":              "stream_name",
		"database":          "database_name",
		"schema":            "schema_name",
		"comment":           "great comment",
		"on_table":          "target_db.target_schema.target_table",
		"on_view":           "target_db.target_schema.target_view",
		"append_only":       true,
		"insert_only":       false,
		"show_initial_rows": true,
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		err := resources.CreateStream(d, db)
		require.ErrorContains(t, err, "all expectations were already fulfilled,")
	})
}

func expectStreamRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows(streamColumns).AddRow(streamRow("APPEND_ONLY")...)
	mock.ExpectQuery(`SHOW STREAMS LIKE 'stream_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
}

func expectOnTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "database_name", "schema_name", "kind", "comment", "cluster_by", "row", "bytes", "owner", "retention_time", "automatic_clustering", "change_tracking", "is_external"}).
		AddRow("", "target_table", "target_db", "target_schema", "TABLE", "mock comment", "", "", "", "", 1, "OFF", "OFF", "N")
	mock.ExpectQuery(`SHOW TABLES LIKE 'target_table' IN SCHEMA "target_db"."target_schema"`).WillReturnRows(rows)
}

func expectOnExternalTableRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "database_name", "schema_name", "kind", "comment", "cluster_by", "row", "bytes", "owner", "retention_time", "automatic_clustering", "change_tracking", "is_external"}).
		AddRow("", "target_table", "target_db", "target_schema", "TABLE", "mock comment", "", "", "", "", 1, "OFF", "OFF", "Y")
	mock.ExpectQuery(`SHOW TABLES LIKE 'target_table' IN SCHEMA "target_db"."target_schema"`).WillReturnRows(rows)
}

func expectOnViewRead(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"created_on", "name", "database_name", "schema_name", "kind", "comment", "cluster_by", "row", "bytes", "owner", "retention_time", "automatic_clustering", "change_tracking", "is_external"}).
		AddRow(time.Now(), "target_view", "target_db", "target_schema", "VIEW", "mock comment", "", "", "", "", 1, "OFF", "OFF", "Y")
	mock.ExpectQuery(`SHOW VIEWS LIKE 'target_view' IN SCHEMA "target_db"."target_schema"`).WillReturnRows(rows)
}

func TestStreamRead(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":    "stream_name",
		"comment": "grand comment",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectStreamRead(mock)
		err := resources.ReadStream(d, db)
		require.NoError(t, err)
		require.Equal(t, "stream_name", d.Get("name").(string))
		require.Equal(t, "database_name", d.Get("database").(string))
		require.Equal(t, "schema_name", d.Get("schema").(string))
		require.Equal(t, "grand comment", d.Get("comment").(string))

		// Test when resource is not found, checking if state will be empty
		require.NotEmpty(t, d.State())
		client := sdk.NewClientFromDB(db)
		ctx := context.Background()
		_, err = client.Streams.ShowByID(ctx, sdk.NewShowByIdStreamRequest(sdk.NewSchemaObjectIdentifier("database", "schema", "name")))
		require.Error(t, err)
		err2 := resources.ReadStream(d, db)
		require.Empty(t, d.State())
		require.Nil(t, err2)
	})
}

func TestStreamReadAppendOnlyMode(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":    "stream_name",
		"comment": "grand comment",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows(streamColumns).AddRow(streamRow("APPEND_ONLY")...)
		mock.ExpectQuery(`SHOW STREAMS LIKE 'stream_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
		err := resources.ReadStream(d, db)
		require.NoError(t, err)
		require.Equal(t, true, d.Get("append_only").(bool))
	})
}

func TestStreamReadInsertOnlyMode(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":    "stream_name",
		"comment": "grand comment",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows(streamColumns).AddRow(streamRow("INSERT_ONLY")...)
		mock.ExpectQuery(`SHOW STREAMS LIKE 'stream_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
		err := resources.ReadStream(d, db)
		require.NoError(t, err)
		require.Equal(t, true, d.Get("insert_only").(bool))
	})
}

func TestStreamReadDefaultMode(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":    "stream_name",
		"comment": "grand comment",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows(streamColumns).AddRow(streamRow("DEFAULT")...)
		mock.ExpectQuery(`SHOW STREAMS LIKE 'stream_name' IN SCHEMA "database_name"."schema_name"`).WillReturnRows(rows)
		err := resources.ReadStream(d, db)
		require.NoError(t, err)
		require.Equal(t, false, d.Get("append_only").(bool))
	})
}

func TestStreamDelete(t *testing.T) {
	d := stream(t, "database_name|schema_name|drop_it", map[string]any{
		"name": "drop_it",
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP STREAM "database_name"."schema_name"."drop_it"`).WillReturnResult(sqlmock.NewResult(1, 1))
		err := resources.DeleteStream(d, db)
		require.NoError(t, err)
	})
}

func TestStreamUpdate(t *testing.T) {
	d := stream(t, "database_name|schema_name|stream_name", map[string]any{
		"name":        "stream_name",
		"database":    "database_name",
		"schema":      "schema_name",
		"comment":     "new stream comment",
		"on_table":    "target_table",
		"append_only": true,
		"insert_only": false,
	})

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`ALTER STREAM "database_name"."schema_name"."stream_name" SET COMMENT = 'new stream comment'`).WillReturnResult(sqlmock.NewResult(1, 1))
		expectStreamRead(mock)
		err := resources.UpdateStream(d, db)
		require.NoError(t, err)
	})
}
