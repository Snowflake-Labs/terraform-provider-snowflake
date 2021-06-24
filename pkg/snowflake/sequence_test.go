package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSequenceCreate(t *testing.T) {
	r := require.New(t)
	s := Sequence("test_sequence", "test_db", "test_schema")

	r.Equal(`"test_db"."test_schema"."test_sequence"`, s.QualifiedName())

	r.Equal(`CREATE SEQUENCE "test_db"."test_schema"."test_sequence"`, s.Create())

	s.WithComment("Test Comment")
	r.Equal(`CREATE SEQUENCE "test_db"."test_schema"."test_sequence" COMMENT = 'Test Comment'`, s.Create())
	s.WithIncrement(5)
	r.Equal(`CREATE SEQUENCE "test_db"."test_schema"."test_sequence" INCREMENT = 5 COMMENT = 'Test Comment'`, s.Create())
	s.WithStart(26)
	r.Equal(`CREATE SEQUENCE "test_db"."test_schema"."test_sequence" START = 26 INCREMENT = 5 COMMENT = 'Test Comment'`, s.Create())
}

func TestSequenceDrop(t *testing.T) {
	r := require.New(t)
	s := Sequence("test_sequence", "test_db", "test_schema")
	r.Equal(`DROP SEQUENCE "test_db"."test_schema"."test_sequence"`, s.Drop())
}

func TestSequenceShow(t *testing.T) {
	r := require.New(t)
	s := Sequence("test_sequence", "test_db", "test_schema")
	r.Equal(`SHOW SEQUENCES LIKE 'test_sequence' IN SCHEMA "test_db"."test_schema"`, s.Show())
}
