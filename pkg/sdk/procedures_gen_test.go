package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func TestProcedures_CreateProcedureForJava(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateProcedureForJavaProcedureOptions {
		return &CreateProcedureForJavaProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateProcedureForJavaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "id",
				ArgDataType: "INTEGER",
			},
			{
				ArgName:     "name",
				ArgDataType: "VARCHAR",
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: "CHAR",
					},
				},
			},
		}
		opts.RuntimeVersion = String("1.8")
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		opts.Imports = []ProcedureImport{
			{
				Import: "test_jar.jar",
			},
		}
		opts.Handler = "TestFunc.echoVarchar"
		opts.ExternalAccessIntegrations = []AccountObjectIdentifier{
			NewAccountObjectIdentifier("ext_integration"),
		}
		opts.Secrets = []ProcedureSecret{
			{
				SecretVariableName: "variable1",
				SecretName:         "name1",
			},
			{
				SecretVariableName: "variable2",
				SecretName:         "name2",
			},
		}
		opts.TargetPath = String("@~/testfunc.jar")
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict: Bool(true),
		}
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile: Bool(true),
		}
		opts.Comment = String("test comment")
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		opts.As = String("return id + name;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (id INTEGER, name VARCHAR) COPY GRANTS RETURNS TABLE (country_code CHAR) LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) TARGET_PATH = '@~/testfunc.jar' STRICT VOLATILE COMMENT = 'test comment' EXECUTE AS CALLER AS 'return id + name;'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateProcedureForJavaScript(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateProcedureForJavaScriptProcedureOptions {
		return &CreateProcedureForJavaScriptProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateProcedureForJavaScriptProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "d",
				ArgDataType: "DOUBLE",
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &ProcedureReturns2{
			ResultDataType: "DOUBLE",
			NotNull:        Bool(true),
		}
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict: Bool(true),
		}
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile: Bool(true),
		}
		opts.Comment = String("test comment")
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		opts.As = String("return 1;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (d DOUBLE) COPY GRANTS RETURNS DOUBLE NOT NULL LANGUAGE JAVASCRIPT STRICT VOLATILE COMMENT = 'test comment' EXECUTE AS CALLER AS 'return 1;'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateProcedureForPython(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateProcedureForPythonProcedureOptions {
		return &CreateProcedureForPythonProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateProcedureForPythonProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "i",
				ArgDataType: "int",
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARIANT",
				Null:           Bool(true),
			},
		}
		opts.RuntimeVersion = String("3.8")
		opts.Packages = []ProcedurePackage{
			{
				Package: "numpy",
			},
			{
				Package: "pandas",
			},
		}
		opts.Imports = []ProcedureImport{
			{
				Import: "numpy",
			},
			{
				Import: "pandas",
			},
		}
		opts.Handler = "udf"
		opts.ExternalAccessIntegrations = []AccountObjectIdentifier{
			NewAccountObjectIdentifier("ext_integration"),
		}
		opts.Secrets = []ProcedureSecret{
			{
				SecretVariableName: "variable1",
				SecretName:         "name1",
			},
			{
				SecretVariableName: "variable2",
				SecretName:         "name2",
			},
		}
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict: Bool(true),
		}
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile: Bool(true),
		}
		opts.Comment = String("test comment")
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		opts.As = String("import numpy as np")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (i int) COPY GRANTS RETURNS VARIANT NULL LANGUAGE PYTHON RUNTIME_VERSION = '3.8' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) STRICT VOLATILE COMMENT = 'test comment' EXECUTE AS CALLER AS 'import numpy as np'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateProcedureForScala(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateProcedureForScalaProcedureOptions {
		return &CreateProcedureForScalaProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateProcedureForScalaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "x",
				ArgDataType: "VARCHAR",
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARCHAR",
				NotNull:        Bool(true),
			},
		}
		opts.RuntimeVersion = String("2.0")
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		opts.Imports = []ProcedureImport{
			{
				Import: "@udf_libs/echohandler.jar",
			},
		}
		opts.Handler = "Echo.echoVarchar"
		opts.TargetPath = String("@~/testfunc.jar")
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict: Bool(true),
		}
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile: Bool(true),
		}
		opts.Comment = String("test comment")
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		opts.As = String("return x")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (x VARCHAR) COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA RUNTIME_VERSION = '2.0' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' TARGET_PATH = '@~/testfunc.jar' STRICT VOLATILE COMMENT = 'test comment' EXECUTE AS CALLER AS 'return x'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateProcedureForSQL(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateProcedureForSQLProcedureOptions {
		return &CreateProcedureForSQLProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateProcedureForSQLProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "message",
				ArgDataType: "VARCHAR",
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &ProcedureReturns3{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARCHAR",
			},
			NotNull: Bool(true),
		}
		opts.StrictOrNot = &ProcedureStrictOrNot{
			Strict: Bool(true),
		}
		opts.VolatileOrNot = &ProcedureVolatileOrNot{
			Volatile: Bool(true),
		}
		opts.Comment = String("test comment")
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		opts.As = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (message VARCHAR) COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SQL STRICT VOLATILE COMMENT = 'test comment' EXECUTE AS CALLER AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
	})
}

func TestProcedures_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DropProcedureOptions {
		return &DropProcedureOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.ArgumentTypes = []ProcedureArgumentType{
			{
				ArgDataType: "VARCHAR",
			},
			{
				ArgDataType: "NUMBER",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP PROCEDURE IF EXISTS %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}

func TestProcedures_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterProcedureOptions {
		return &AlterProcedureOptions{
			name:     id,
			IfExists: Bool(true),
			ArgumentTypes: []ProcedureArgumentType{
				{
					ArgDataType: "VARCHAR",
				},
				{
					ArgDataType: "NUMBER",
				},
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterProcedureOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), random.StringN(12))
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: execute as", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExecuteAs = &ProcedureExecuteAs{
			Caller: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) EXECUTE AS CALLER`, id.FullyQualifiedName())
	})

	t.Run("alter: set log level", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ProcedureSet{
			LogLevel: String("DEBUG"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set trace level", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ProcedureSet{
			TraceLevel: String("DEBUG"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET TRACE_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ProcedureSet{
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ProcedureUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestProcedures_Show(t *testing.T) {
	defaultOpts := func() *ShowProcedureOptions {
		return &ShowProcedureOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowProcedureOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW PROCEDURES`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PROCEDURES LIKE 'pattern'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PROCEDURES IN ACCOUNT`)
	})
}

func TestProcedures_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeProcedureOptions {
		return &DescribeProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeProcedureOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PROCEDURE %s`, id.FullyQualifiedName())
	})
}
