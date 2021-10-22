package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getFunction(withArgs bool) *FunctionBuilder {
	s := Function("test_db", "test_schema", "test_func", []string{})
	s.WithStatement(`var message = "Hi"` + "\n" + `return message`)
	s.WithReturnType("varchar")
	if withArgs {
		s.WithArgs([]map[string]string{
			{"name": "user", "type": "varchar"},
			{"name": "eventdt", "type": "date"}})
	}
	return s
}

func TestFunctionQualifiedName(t *testing.T) {
	r := require.New(t)
	s := getFunction(true)
	qn, _ := s.QualifiedName()
	r.Equal(`"test_db"."test_schema"."test_func"(VARCHAR, DATE)`, qn)
	qna, _ := s.QualifiedNameWithoutArguments()
	r.Equal(`"test_db"."test_schema"."test_func"`, qna)
}

func TestFunctionCreate(t *testing.T) {
	r := require.New(t)
	s := getFunction(true)

	r.Equal([]string{"VARCHAR", "DATE"}, s.ArgTypes())
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithOptionalParams(t *testing.T) {
	r := require.New(t)
	s := getFunction(true)
	s.WithNullInputBehavior("RETURNS NULL ON NULL INPUT")
	s.WithReturnBehavior("IMMUTABLE")
	s.WithComment("this is cool func!")
	s.WithLanguage("JAVASCRIPT")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR LANGUAGE JAVASCRIPT RETURNS NULL ON NULL INPUT` +
		` IMMUTABLE COMMENT = 'this is cool func!' AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestFunctionDrop(t *testing.T) {
	r := require.New(t)

	// Without arg
	s := getFunction(false)
	stmnt, _ := s.Drop()
	r.Equal(stmnt, `DROP FUNCTION "test_db"."test_schema"."test_func"()`)

	// With arg
	ss := getFunction(true)
	stmnt, _ = ss.Drop()
	r.Equal(`DROP FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE)`, stmnt)
}

func TestFunctionShow(t *testing.T) {
	r := require.New(t)
	s := getFunction(false)
	stmnt := s.Show()
	r.Equal(stmnt, `SHOW USER FUNCTIONS LIKE 'test_func' IN SCHEMA "test_db"."test_schema"`)
}

func TestFunctionRename(t *testing.T) {
	r := require.New(t)
	s := getFunction(false)

	stmnt, _ := s.Rename("new_func")
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"() RENAME TO "test_db"."test_schema"."new_func"`
	r.Equal(expected, stmnt)
}

func TestFunctionChangeComment(t *testing.T) {
	r := require.New(t)
	s := getFunction(true)

	stmnt, _ := s.ChangeComment("not used")
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE) SET COMMENT = 'not used'`
	r.Equal(expected, stmnt)
}

func TestFunctionRemoveComment(t *testing.T) {
	r := require.New(t)
	s := getFunction(false)

	stmnt, _ := s.RemoveComment()
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"() UNSET COMMENT`
	r.Equal(expected, stmnt)
}

func TestFunctionDescribe(t *testing.T) {
	r := require.New(t)
	s := getFunction(false)

	stmnt, _ := s.Describe()
	expected := `DESCRIBE FUNCTION "test_db"."test_schema"."test_func"()`
	r.Equal(expected, stmnt)
}

func TestFunctionArgumentsSignature(t *testing.T) {
	r := require.New(t)
	s := getFunction(false)
	sign, _ := s.ArgumentsSignature()
	r.Equal("test_func() RETURN VARCHAR", sign)
	s = getFunction(true)
	sign, _ = s.ArgumentsSignature()
	r.Equal("test_func(VARCHAR, DATE) RETURN VARCHAR", sign)
}
