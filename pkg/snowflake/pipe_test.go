package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeCreate(t *testing.T) {
	r := require.New(t)
	s := Pipe("test_pipe", "test_db", "test_schema")
	r.Equal(s.QualifiedName(), `"test_db"."test_schema"."test_pipe"`)

	r.Equal(s.Create(), `CREATE PIPE "test_db"."test_schema"."test_pipe"`)

	s.WithAutoIngest()
	r.Equal(s.Create(), `CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE`)

	s.WithComment("Yeehaw")
	r.Equal(s.Create(), `CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE COMMENT = 'Yeehaw'`)

	s.WithCopyStatement("test copy statement ")
	r.Equal(s.Create(), `CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE COMMENT = 'Yeehaw' AS test copy statement `)
}

func TestPipeChangeComment(t *testing.T) {
	r := require.New(t)
	s := Pipe("test_pipe", "test_db", "test_schema")
	r.Equal(s.ChangeComment("worst pipe ever"), `ALTER PIPE "test_db"."test_schema"."test_pipe" SET COMMENT = 'worst pipe ever'`)
}

func TestPipeRemoveComment(t *testing.T) {
	r := require.New(t)
	s := Pipe("test_pipe", "test_db", "test_schema")
	r.Equal(s.RemoveComment(), `ALTER PIPE "test_db"."test_schema"."test_pipe" UNSET COMMENT`)
}

func TestPipeDrop(t *testing.T) {
	r := require.New(t)
	s := Pipe("test_pipe", "test_db", "test_schema")
	r.Equal(s.Drop(), `DROP PIPE "test_db"."test_schema"."test_pipe"`)
}

func TestPipeShow(t *testing.T) {
	r := require.New(t)
	s := Pipe("test_pipe", "test_db", "test_schema")
	r.Equal(s.Show(), `SHOW PIPES LIKE 'test_pipe' IN DATABASE "test_db"`)
}
