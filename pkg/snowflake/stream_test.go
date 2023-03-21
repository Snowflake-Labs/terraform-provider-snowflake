package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamCreate(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")

	s.WithOnTable("test_db", "test_schema", "test_target_table")
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" APPEND_ONLY = false INSERT_ONLY = false SHOW_INITIAL_ROWS = false`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = false INSERT_ONLY = false SHOW_INITIAL_ROWS = false`, s.Create())

	s.WithShowInitialRows(true)
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = false INSERT_ONLY = false SHOW_INITIAL_ROWS = true`, s.Create())

	s.WithAppendOnly(true)
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = true INSERT_ONLY = false SHOW_INITIAL_ROWS = true`, s.Create())

	s.WithInsertOnly(true)
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = true INSERT_ONLY = true SHOW_INITIAL_ROWS = true`, s.Create())

	s.WithExternalTable(true)
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON EXTERNAL TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = true INSERT_ONLY = true SHOW_INITIAL_ROWS = true`, s.Create())
}

func TestStreamOnStageCreate(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")

	s.WithOnStage("test_db", "test_schema", "test_target_stage")
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON STAGE "test_db"."test_schema"."test_target_stage"`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(`CREATE STREAM "test_db"."test_schema"."test_stream" ON STAGE "test_db"."test_schema"."test_target_stage" COMMENT = 'Test Comment'`, s.Create())
}

func TestStreamChangeComment(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(`ALTER STREAM "test_db"."test_schema"."test_stream" SET COMMENT = 'new stream comment'`, s.ChangeComment("new stream comment"))
}

func TestStreamRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(`ALTER STREAM "test_db"."test_schema"."test_stream" UNSET COMMENT`, s.RemoveComment())
}

func TestStreamDrop(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(`DROP STREAM "test_db"."test_schema"."test_stream"`, s.Drop())
}

func TestStreamShow(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(`SHOW STREAMS LIKE 'test_stream' IN SCHEMA "test_db"."test_schema"`, s.Show())
}
