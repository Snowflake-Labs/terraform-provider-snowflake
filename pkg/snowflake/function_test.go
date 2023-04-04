package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getJavaScriptFuction(withArgs bool) *FunctionBuilder {
	s := NewFunctionBuilder("test_db", "test_schema", "test_func", []string{})
	s.WithStatement(`var message = "Hi"` + "\n" + `return message`)
	s.WithReturnType("varchar")
	if withArgs {
		s.WithArgs([]map[string]string{
			{"name": "user", "type": "varchar"},
			{"name": "eventdt", "type": "date"},
		})
	}
	return s
}

const javafunc = `class CoolFunc {` + "\n" +
	`  public static String test(String u, int c) {` + "\n" +
	`    return u;` + "\n" +
	`  }` + "\n" +
	`}`

func getJavaFuction(withArgs bool) *FunctionBuilder {
	s := NewFunctionBuilder("test_db", "test_schema", "test_func", []string{})
	s.WithReturnType("varchar")
	s.WithStatement(javafunc)
	if withArgs {
		s.WithArgs([]map[string]string{
			{"name": "user", "type": "varchar"},
			{"name": "count", "type": "number"},
		})
	}
	return s
}

const pythonfunc = `def add_py(i):` + "\n " +
	` return i+1`

func getPythonFunction(withArgs bool) *FunctionBuilder {
	s := NewFunctionBuilder("test_db", "test_schema", "test_func", []string{})
	s.WithReturnType("int")
	s.WithStatement(pythonfunc)
	if withArgs {
		s.WithArgs([]map[string]string{
			{"name": "arg", "type": "int"},
		})
	}
	return s
}

func TestFunctionQualifiedName(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)
	qn, _ := s.QualifiedName()
	r.Equal(`"test_db"."test_schema"."test_func"(VARCHAR, DATE)`, qn)
	qna, _ := s.QualifiedNameWithoutArguments()
	r.Equal(`"test_db"."test_schema"."test_func"`, qna)
}

func TestFunctionCreate(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)

	r.Equal([]string{"VARCHAR", "DATE"}, s.ArgTypes())
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithSecure(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)
	s.WithSecure()

	r.Equal([]string{"VARCHAR", "DATE"}, s.ArgTypes())
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE SECURE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, eventdt DATE) RETURNS VARCHAR AS $$` +
		`var message = "Hi"` + "\nreturn message$$"
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithJavaScriptFunction(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)
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

func TestFunctionCreateWithJavaFunction(t *testing.T) {
	r := require.New(t)
	s := getJavaFuction(true)
	s.WithNullInputBehavior("RETURNS NULL ON NULL INPUT")
	s.WithReturnBehavior("IMMUTABLE")
	s.WithComment("this is cool func!")
	s.WithLanguage("JAVA")
	s.WithHandler("CoolFunc.test")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, count NUMBER) RETURNS VARCHAR` +
		` LANGUAGE JAVA RETURNS NULL ON NULL INPUT IMMUTABLE COMMENT = 'this is cool func!'` +
		` HANDLER = 'CoolFunc.test' AS $$` + javafunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithJavaFunctionWithImports(t *testing.T) {
	r := require.New(t)
	s := getJavaFuction(true)
	s.WithNullInputBehavior("RETURNS NULL ON NULL INPUT")
	s.WithReturnBehavior("IMMUTABLE")
	s.WithComment("this is cool func!")
	s.WithLanguage("JAVA")
	s.WithImports([]string{"@~/stage/myudf1.jar", "@~/stage/myudf2.jar"})
	s.WithHandler("CoolFunc.test")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, count NUMBER) RETURNS VARCHAR` +
		` LANGUAGE JAVA RETURNS NULL ON NULL INPUT IMMUTABLE COMMENT = 'this is cool func!'` +
		` IMPORTS = ('@~/stage/myudf1.jar', '@~/stage/myudf2.jar') HANDLER = 'CoolFunc.test'` +
		` AS $$` + javafunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithJavaFunctionWithTargetPath(t *testing.T) {
	r := require.New(t)
	s := getJavaFuction(true)
	s.WithNullInputBehavior("RETURNS NULL ON NULL INPUT")
	s.WithReturnBehavior("IMMUTABLE")
	s.WithComment("this is cool func!")
	s.WithLanguage("JAVA")
	s.WithTargetPath("@~/stage/myudf1.jar")
	s.WithHandler("CoolFunc.test")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(user VARCHAR, count NUMBER) RETURNS VARCHAR` +
		` LANGUAGE JAVA RETURNS NULL ON NULL INPUT IMMUTABLE COMMENT = 'this is cool func!'` +
		` HANDLER = 'CoolFunc.test' TARGET_PATH = '@~/stage/myudf1.jar'` +
		` AS $$` + javafunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithPythonFunction(t *testing.T) {
	r := require.New(t)
	s := getPythonFunction(true)
	s.WithComment("this is cool func!")
	s.WithLanguage("PYTHON")
	s.WithRuntimeVersion("3.8")
	s.WithHandler("CoolFunc.test")
	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(arg INT) RETURNS INT` +
		` LANGUAGE PYTHON RUNTIME_VERSION = '3.8' COMMENT = 'this is cool func!'` +
		` HANDLER = 'CoolFunc.test' AS $$` + pythonfunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithPythonFunctionWithPackages(t *testing.T) {
	r := require.New(t)
	s := getPythonFunction(true)

	pkgs := []string{"numpy", "pandas"}

	s.WithComment("this is cool func!")
	s.WithLanguage("PYTHON")
	s.WithRuntimeVersion("3.8")
	s.WithPackages(pkgs)
	s.WithHandler("CoolFunc.test")

	createStmnt, _ := s.Create()

	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(arg INT) RETURNS INT` +
		` LANGUAGE PYTHON RUNTIME_VERSION = '3.8' PACKAGES = ('numpy', 'pandas') COMMENT = 'this is cool func!'` +
		` HANDLER = 'CoolFunc.test' AS $$` + pythonfunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithPythonFunctionWithImports(t *testing.T) {
	r := require.New(t)
	s := getPythonFunction(true)
	s.WithComment("this is cool func!")
	s.WithLanguage("PYTHON")
	s.WithRuntimeVersion("3.8")
	s.WithHandler("CoolFunc.test")
	s.WithImports([]string{"@~/stage/myudf1.py", "@~/stage/myudf2.py"})

	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(arg INT) RETURNS INT` +
		` LANGUAGE PYTHON RUNTIME_VERSION = '3.8' COMMENT = 'this is cool func!'` +
		` IMPORTS = ('@~/stage/myudf1.py', '@~/stage/myudf2.py') HANDLER = 'CoolFunc.test'` +
		` AS $$` + pythonfunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionCreateWithPythonFunctionWithTargetPath(t *testing.T) {
	r := require.New(t)
	s := getPythonFunction(true)
	s.WithComment("this is cool func!")
	s.WithLanguage("PYTHON")
	s.WithRuntimeVersion("3.8")
	s.WithTargetPath("@~/stage/myudf1.py")
	s.WithHandler("CoolFunc.test")

	createStmnt, _ := s.Create()
	expected := `CREATE OR REPLACE FUNCTION "test_db"."test_schema"."test_func"` +
		`(arg INT) RETURNS INT` +
		` LANGUAGE PYTHON RUNTIME_VERSION = '3.8' COMMENT = 'this is cool func!'` +
		` HANDLER = 'CoolFunc.test' TARGET_PATH = '@~/stage/myudf1.py'` +
		` AS $$` + pythonfunc + `$$`
	r.Equal(expected, createStmnt)
}

func TestFunctionDrop(t *testing.T) {
	r := require.New(t)

	// Without arg
	s := getJavaScriptFuction(false)
	stmnt, _ := s.Drop()
	r.Equal(`DROP FUNCTION "test_db"."test_schema"."test_func"()`, stmnt)

	// With arg
	ss := getJavaScriptFuction(true)
	stmnt, _ = ss.Drop()
	r.Equal(`DROP FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE)`, stmnt)
}

func TestFunctionShow(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(false)
	stmnt := s.Show()
	r.Equal(`SHOW USER FUNCTIONS LIKE 'test_func' IN SCHEMA "test_db"."test_schema"`, stmnt)
}

func TestFunctionRename(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(false)

	stmnt, _ := s.Rename("new_func")
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"() RENAME TO "test_db"."test_schema"."new_func"`
	r.Equal(expected, stmnt)
}

func TestFunctionSecure(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)

	stmnt, _ := s.Secure()
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE) SET SECURE`
	r.Equal(expected, stmnt)
}

func TestFunctionUnsecure(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)

	stmnt, _ := s.Unsecure()
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE) UNSET SECURE`
	r.Equal(expected, stmnt)
}

func TestFunctionChangeComment(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(true)

	stmnt, _ := s.ChangeComment("not used")
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"(VARCHAR, DATE) SET COMMENT = 'not used'`
	r.Equal(expected, stmnt)
}

func TestFunctionRemoveComment(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(false)

	stmnt, _ := s.RemoveComment()
	expected := `ALTER FUNCTION "test_db"."test_schema"."test_func"() UNSET COMMENT`
	r.Equal(expected, stmnt)
}

func TestFunctionDescribe(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(false)

	stmnt, _ := s.Describe()
	expected := `DESCRIBE FUNCTION "test_db"."test_schema"."test_func"()`
	r.Equal(expected, stmnt)
}

func TestFunctionArgumentsSignature(t *testing.T) {
	r := require.New(t)
	s := getJavaScriptFuction(false)
	sign, _ := s.ArgumentsSignature()
	r.Equal("test_func() RETURN VARCHAR", sign)
	s = getJavaScriptFuction(true)
	sign, _ = s.ArgumentsSignature()
	r.Equal("test_func(VARCHAR, DATE) RETURN VARCHAR", sign)
}
