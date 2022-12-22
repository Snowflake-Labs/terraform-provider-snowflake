package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeCreate(t *testing.T) {
	r := require.New(t)
	s := NewPipeBuilder("test_pipe", "test_db", "test_schema")
	r.Equal(`"test_db"."test_schema"."test_pipe"`, s.QualifiedName())

	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe"`, s.Create())

	s.WithAutoIngest()
	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE`, s.Create())

	s.WithComment("Yeehaw")
	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE COMMENT = 'Yeehaw'`, s.Create())

	s.WithCopyStatement("test copy statement ")
	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE COMMENT = 'Yeehaw' AS test copy statement `, s.Create())

	s.WithAwsSnsTopicArn("arn:aws:sns:us-east-1:1234567890123456:mytopic")
	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE AWS_SNS_TOPIC = 'arn:aws:sns:us-east-1:1234567890123456:mytopic' COMMENT = 'Yeehaw' AS test copy statement `, s.Create())

	s.WithIntegration("myintegration")
	r.Equal(`CREATE PIPE "test_db"."test_schema"."test_pipe" AUTO_INGEST = TRUE INTEGRATION = 'myintegration' AWS_SNS_TOPIC = 'arn:aws:sns:us-east-1:1234567890123456:mytopic' COMMENT = 'Yeehaw' AS test copy statement `, s.Create())
}

func TestPipeChangeComment(t *testing.T) {
	r := require.New(t)
	s := NewPipeBuilder("test_pipe", "test_db", "test_schema")
	r.Equal(`ALTER PIPE "test_db"."test_schema"."test_pipe" SET COMMENT = 'worst pipe ever'`, s.ChangeComment("worst pipe ever"))
}

func TestPipeRemoveComment(t *testing.T) {
	r := require.New(t)
	s := NewPipeBuilder("test_pipe", "test_db", "test_schema")
	r.Equal(`ALTER PIPE "test_db"."test_schema"."test_pipe" UNSET COMMENT`, s.RemoveComment())
}

func TestPipeDrop(t *testing.T) {
	r := require.New(t)
	s := NewPipeBuilder("test_pipe", "test_db", "test_schema")
	r.Equal(`DROP PIPE "test_db"."test_schema"."test_pipe"`, s.Drop())
}

func TestPipeShow(t *testing.T) {
	r := require.New(t)
	s := NewPipeBuilder("test_pipe", "test_db", "test_schema")
	r.Equal(`SHOW PIPES LIKE 'test_pipe' IN SCHEMA "test_db"."test_schema"`, s.Show())
}
