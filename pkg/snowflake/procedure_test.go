package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getProcedure(withArgs bool) *ProcedureBuilder {
	s := Procedure("test_db", "test_schema", "test_proc", []string{})
	s.WithStatement(`var message = "Hi"` + "\n" + `return message`)
	s.WithReturnType("varchar")
	s.WithExecuteAs("CALLER")
	if withArgs {
		s.WithArgs([]map[string]string{
			{"name": "user", "type": "varchar"},
			{"name": "eventdt", "type": "date"}})
	}
	return s
}

func TestProcedureQualifiedName(t *testing.T) {
	r := require.New(t)
	s := getProcedure(true)
	qn, _ := s.QualifiedName()
	r.Equal(`"test_db"."test_schema"."test_proc"(VARCHAR, DATE)`, qn)
	qna, _ := s.QualifiedNameWithoutArguments()
	r.Equal(`"test_db"."test_schema"."test_proc"`, qna)
}

func TestProcedureCreate(t *testing.T) {
	r := require.New(t)
	s := getProcedure(true)

	r.Equal([]string{"VARCHAR", "DATE"}, s.ArgTypes())
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE PROCEDURE "test_db"."test_schema"."test_proc"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR LANGUAGE javascript EXECUTE AS CALLER AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestProcedureCreateWithOptionalParams(t *testing.T) {
	r := require.New(t)
	s := getProcedure(true)
	s.WithNullInputBehavior("RETURNS NULL ON NULL INPUT")
	s.WithReturnBehavior("IMMUTABLE")
	s.WithComment("this is cool proc!")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE PROCEDURE "test_db"."test_schema"."test_proc"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR LANGUAGE javascript RETURNS NULL ON NULL INPUT` +
		` IMMUTABLE COMMENT = 'this is cool proc!' EXECUTE AS CALLER AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestProcedureDrop(t *testing.T) {
	r := require.New(t)

	// Without arg
	s := getProcedure(false)
	stmnt, _ := s.Drop()
	r.Equal(stmnt, `DROP PROCEDURE "test_db"."test_schema"."test_proc"()`)

	// With arg
	ss := getProcedure(true)
	stmnt, _ = ss.Drop()
	r.Equal(`DROP PROCEDURE "test_db"."test_schema"."test_proc"(VARCHAR, DATE)`, stmnt)
}

func TestProcedureShow(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)
	stmnt := s.Show()
	r.Equal(stmnt, `SHOW PROCEDURES LIKE 'test_proc' IN SCHEMA "test_db"."test_schema"`)
}

func TestProcedureRename(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)

	stmnt, _ := s.Rename("new_proc")
	expected := `ALTER PROCEDURE "test_db"."test_schema"."test_proc"() RENAME TO "test_db"."test_schema"."new_proc"`
	r.Equal(expected, stmnt)
}

func TestProcedureChangeComment(t *testing.T) {
	r := require.New(t)
	s := getProcedure(true)

	stmnt, _ := s.ChangeComment("not used")
	expected := `ALTER PROCEDURE "test_db"."test_schema"."test_proc"(VARCHAR, DATE) SET COMMENT = 'not used'`
	r.Equal(expected, stmnt)
}

func TestProcedureRemoveComment(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)

	stmnt, _ := s.RemoveComment()
	expected := `ALTER PROCEDURE "test_db"."test_schema"."test_proc"() UNSET COMMENT`
	r.Equal(expected, stmnt)
}

func TestProcedureChangeExecuteAs(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)

	stmnt, _ := s.ChangeExecuteAs("OWNER")
	expected := `ALTER PROCEDURE "test_db"."test_schema"."test_proc"() EXECUTE AS OWNER`
	r.Equal(expected, stmnt)
}

func TestProcedureDescribe(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)

	stmnt, _ := s.Describe()
	expected := `DESCRIBE PROCEDURE "test_db"."test_schema"."test_proc"()`
	r.Equal(expected, stmnt)
}

func TestProcedureArgumentsSignature(t *testing.T) {
	r := require.New(t)
	s := getProcedure(false)
	sign, _ := s.ArgumentsSignature()
	r.Equal("TEST_PROC()", sign)
	s = getProcedure(true)
	sign, _ = s.ArgumentsSignature()
	r.Equal("TEST_PROC(VARCHAR, DATE)", sign)
}
