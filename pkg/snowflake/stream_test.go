package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamCreate(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")

	s.WithOnTable("test_db", "test_schema", "test_target_table")
	r.Equal(s.Create(), `CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" APPEND_ONLY = false`)

	s.WithComment("Test Comment")
	r.Equal(s.Create(), `CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = false`)

	s.WithAppendOnly(true)
	r.Equal(s.Create(), `CREATE STREAM "test_db"."test_schema"."test_stream" ON TABLE "test_db"."test_schema"."test_target_table" COMMENT = 'Test Comment' APPEND_ONLY = true`)
}

func TestStreamChangeComment(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(s.ChangeComment("new stream comment"), `ALTER STREAM "test_db"."test_schema"."test_stream" SET COMMENT = 'new stream comment'`)
}

func TestStreamRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(s.RemoveComment(), `ALTER STREAM "test_db"."test_schema"."test_stream" UNSET COMMENT`)
}

func TestStreamDrop(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP STREAM "test_db"."test_schema"."test_stream"`)
}

func TestStreamShow(t *testing.T) {
	r := require.New(t)
	s := Stream("test_stream", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW STREAMS LIKE 'test_stream' IN DATABASE "test_db"`)
}
