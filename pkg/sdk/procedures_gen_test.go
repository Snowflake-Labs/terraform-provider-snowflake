package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func TestProcedures_CreateForJava(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForJavaProcedureOptions {
		return &CreateForJavaProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateForJavaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: function definition", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetPath = String("@~/testfunc.jar")
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("TARGET_PATH must be nil when AS is nil"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavaProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavaProcedureOptions", "RuntimeVersion"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "id",
				ArgDataType: DataTypeNumber,
			},
			{
				ArgName:      "name",
				ArgDataType:  DataTypeVARCHAR,
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
		}
		opts.RuntimeVersion = "1.8"
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
		opts.Secrets = []Secret{
			{
				VariableName: "variable1",
				Name:         "name1",
			},
			{
				VariableName: "variable2",
				Name:         "name2",
			},
		}
		opts.TargetPath = String("@~/testfunc.jar")
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.Comment = String("test comment")
		opts.ExecuteAs = ExecuteAsPointer(ExecuteAsCaller)
		opts.ProcedureDefinition = String("return id + name;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (id NUMBER, name VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE (country_code VARCHAR) LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) TARGET_PATH = '@~/testfunc.jar' STRICT COMMENT = 'test comment' EXECUTE AS CALLER AS 'return id + name;'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateForJavaScript(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForJavaScriptProcedureOptions {
		return &CreateForJavaScriptProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateForJavaScriptProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavaScriptProcedureOptions", "ProcedureDefinition"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "d",
				ArgDataType:  "DOUBLE",
				DefaultValue: String("1.0"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.ResultDataType = "DOUBLE"
		opts.NotNull = Bool(true)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.Comment = String("test comment")
		opts.ExecuteAs = ExecuteAsPointer(ExecuteAsCaller)
		opts.ProcedureDefinition = "return 1;"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (d DOUBLE DEFAULT 1.0) COPY GRANTS RETURNS DOUBLE NOT NULL LANGUAGE JAVASCRIPT STRICT COMMENT = 'test comment' EXECUTE AS CALLER AS 'return 1;'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateForPython(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForPythonProcedureOptions {
		return &CreateForPythonProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateForPythonProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonProcedureOptions", "RuntimeVersion"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "i",
				ArgDataType:  "int",
				DefaultValue: String("1"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARIANT",
				Null:           Bool(true),
			},
		}
		opts.RuntimeVersion = "3.8"
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
		opts.Secrets = []Secret{
			{
				VariableName: "variable1",
				Name:         "name1",
			},
			{
				VariableName: "variable2",
				Name:         "name2",
			},
		}
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.Comment = String("test comment")
		opts.ExecuteAs = ExecuteAsPointer(ExecuteAsCaller)
		opts.ProcedureDefinition = String("import numpy as np")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (i int DEFAULT 1) COPY GRANTS RETURNS VARIANT NULL LANGUAGE PYTHON RUNTIME_VERSION = '3.8' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) STRICT COMMENT = 'test comment' EXECUTE AS CALLER AS 'import numpy as np'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateForScala(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForScalaProcedureOptions {
		return &CreateForScalaProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateForScalaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: function definition", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetPath = String("@~/testfunc.jar")
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("TARGET_PATH must be nil when AS is nil"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForScalaProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForScalaProcedureOptions", "RuntimeVersion"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "x",
				ArgDataType:  "VARCHAR",
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARCHAR",
				NotNull:        Bool(true),
			},
		}
		opts.RuntimeVersion = "2.0"
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
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.Comment = String("test comment")
		opts.ExecuteAs = ExecuteAsPointer(ExecuteAsCaller)
		opts.ProcedureDefinition = String("return x")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (x VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA RUNTIME_VERSION = '2.0' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' TARGET_PATH = '@~/testfunc.jar' STRICT COMMENT = 'test comment' EXECUTE AS CALLER AS 'return x'`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateForSQL(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForSQLProcedureOptions {
		return &CreateForSQLProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateForSQLProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForSQLProcedureOptions", "ProcedureDefinition"))
	})

	t.Run("create with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureSQLReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		opts.ProcedureDefinition = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE PROCEDURE %s () RETURNS FLOAT LANGUAGE SQL AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "message",
				ArgDataType:  "VARCHAR",
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = ProcedureSQLReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARCHAR",
			},
			NotNull: Bool(true),
		}
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.Comment = String("test comment")
		opts.ExecuteAs = ExecuteAsPointer(ExecuteAsCaller)
		opts.ProcedureDefinition = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE PROCEDURE %s (message VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SQL STRICT COMMENT = 'test comment' EXECUTE AS CALLER AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
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

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP PROCEDURE %s ()`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		assertOptsValidAndSQLEquals(t, opts, `DROP PROCEDURE IF EXISTS %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}

func TestProcedures_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterProcedureOptions {
		return &AlterProcedureOptions{
			name:              id,
			IfExists:          Bool(true),
			ArgumentDataTypes: []DataType{DataTypeVARCHAR, DataTypeNumber},
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

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetLogLevel = String("DEBUG")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterProcedureOptions", "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "UnsetComment", "SetTags", "UnsetTags", "ExecuteAs"))
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), random.StringN(12))
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: execute as", func(t *testing.T) {
		opts := defaultOpts()
		executeAs := ExecuteAsCaller
		opts.ExecuteAs = &executeAs
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) EXECUTE AS CALLER`, id.FullyQualifiedName())
	})

	t.Run("alter: set log level", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetLogLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set log level with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = nil
		opts.SetLogLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s () SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set trace level", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTraceLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET TRACE_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER PROCEDURE IF EXISTS %s (VARCHAR, NUMBER) SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
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

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PROCEDURE %s ()`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PROCEDURE %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}

func TestProcedures_Call(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CallProcedureOptions {
		return &CallProcedureOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CallProcedureOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CALL %s ()`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"province => 'Manitoba'", "amount => 127.4"}
		assertOptsValidAndSQLEquals(t, opts, `CALL %s (province => 'Manitoba', amount => 127.4) INTO :ret`, id.FullyQualifiedName())

		opts = defaultOpts()
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"'Manitoba'", "127.4"}
		assertOptsValidAndSQLEquals(t, opts, `CALL %s ('Manitoba', 127.4) INTO :ret`, id.FullyQualifiedName())
	})
}

func TestProcedures_CreateAndCallForJava(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateAndCallForJavaProcedureOptions {
		return &CreateAndCallForJavaProcedureOptions{
			Name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAndCallForJavaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForJavaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "name",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForJavaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForJavaProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForJavaProcedureOptions", "RuntimeVersion"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{},
		}
		opts.RuntimeVersion = "1.8"
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:latest",
			},
		}
		opts.Handler = "TestFunc.echoVarchar"
		opts.ProcedureName = id
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:latest') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`, id.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "id",
				ArgDataType: DataTypeNumber,
			},
			{
				ArgName:     "name",
				ArgDataType: DataTypeVARCHAR,
			},
		}
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
		}
		opts.RuntimeVersion = "1.8"
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
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.ProcedureDefinition = String("return id + name;")
		cte := NewAccountObjectIdentifier("album_info_1976")
		opts.WithClause = &ProcedureWithClause{
			CteName:    cte,
			CteColumns: []string{"x", "y"},
			Statement:  "(select m.album_ID, m.album_name, b.band_name from music_albums)",
		}
		opts.ProcedureName = id
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"1", "rnd"}
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE (id NUMBER, name VARCHAR) RETURNS TABLE (country_code VARCHAR) LANGUAGE JAVA RUNTIME_VERSION = '1.8' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' STRICT AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`, id.FullyQualifiedName(), cte.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestProcedures_CreateAndCallForScala(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateAndCallForScalaProcedureOptions {
		return &CreateAndCallForScalaProcedureOptions{
			Name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAndCallForScalaProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForScalaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "name",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForScalaProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForScalaProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForScalaProcedureOptions", "RuntimeVersion"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{},
		}
		opts.RuntimeVersion = "2.12"
		opts.Packages = []ProcedurePackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		opts.Handler = "TestFunc.echoVarchar"
		opts.ProcedureName = id
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE SCALA RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' CALL %s ()`, id.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:     "id",
				ArgDataType: DataTypeNumber,
			},
			{
				ArgName:     "name",
				ArgDataType: DataTypeVARCHAR,
			},
		}
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
		}
		opts.RuntimeVersion = "2.12"
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
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.ProcedureDefinition = String("return id + name;")
		cte := NewAccountObjectIdentifier("album_info_1976")
		opts.WithClauses = []ProcedureWithClause{
			{
				CteName:    cte,
				CteColumns: []string{"x", "y"},
				Statement:  "(select m.album_ID, m.album_name, b.band_name from music_albums)",
			},
		}
		opts.ProcedureName = id
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"1", "rnd"}
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE (id NUMBER, name VARCHAR) RETURNS TABLE (country_code VARCHAR) LANGUAGE SCALA RUNTIME_VERSION = '2.12' PACKAGES = ('com.snowflake:snowpark:1.2.0') IMPORTS = ('test_jar.jar') HANDLER = 'TestFunc.echoVarchar' STRICT AS 'return id + name;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1, rnd) INTO :ret`, id.FullyQualifiedName(), cte.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestProcedures_CreateAndCallForPython(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateAndCallForPythonProcedureOptions {
		return &CreateAndCallForPythonProcedureOptions{
			Name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAndCallForPythonProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForPythonProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "name",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForPythonProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForPythonProcedureOptions", "Handler"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForPythonProcedureOptions", "RuntimeVersion"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{},
		}
		opts.RuntimeVersion = "3.8"
		opts.Packages = []ProcedurePackage{
			{
				Package: "snowflake-snowpark-python",
			},
		}
		opts.Handler = "udf"
		opts.ProcedureName = id
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE PYTHON RUNTIME_VERSION = '3.8' PACKAGES = ('snowflake-snowpark-python') HANDLER = 'udf' CALL %s ()`, id.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "i",
				ArgDataType:  "int",
				DefaultValue: String("1"),
			},
		}
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: "VARIANT",
				Null:           Bool(true),
			},
		}
		opts.RuntimeVersion = "3.8"
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
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.ProcedureDefinition = String("import numpy as np")
		cte := NewAccountObjectIdentifier("album_info_1976")
		opts.WithClauses = []ProcedureWithClause{
			{
				CteName:    cte,
				CteColumns: []string{"x", "y"},
				Statement:  "(select m.album_ID, m.album_name, b.band_name from music_albums)",
			},
		}
		opts.ProcedureName = id
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"1"}
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE (i int DEFAULT 1) RETURNS VARIANT NULL LANGUAGE PYTHON RUNTIME_VERSION = '3.8' PACKAGES = ('numpy', 'pandas') IMPORTS = ('numpy', 'pandas') HANDLER = 'udf' STRICT AS 'import numpy as np' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`, id.FullyQualifiedName(), cte.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestProcedures_CreateAndCallForJavaScript(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateAndCallForJavaScriptProcedureOptions {
		return &CreateAndCallForJavaScriptProcedureOptions{
			Name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAndCallForJavaScriptProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForJavaScriptProcedureOptions", "ProcedureDefinition"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResultDataType = "DOUBLE"
		opts.ProcedureDefinition = "return 1;"
		opts.ProcedureName = id
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE () RETURNS DOUBLE LANGUAGE JAVASCRIPT AS 'return 1;' CALL %s ()`, id.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "d",
				ArgDataType:  "DOUBLE",
				DefaultValue: String("1.0"),
			},
		}
		opts.ResultDataType = "DOUBLE"
		opts.NotNull = Bool(true)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.ProcedureDefinition = "return 1;"
		cte := NewAccountObjectIdentifier("album_info_1976")
		opts.WithClauses = []ProcedureWithClause{
			{
				CteName:    cte,
				CteColumns: []string{"x", "y"},
				Statement:  "(select m.album_ID, m.album_name, b.band_name from music_albums)",
			},
		}
		opts.ProcedureName = id
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"1"}
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE (d DOUBLE DEFAULT 1.0) RETURNS DOUBLE NOT NULL LANGUAGE JAVASCRIPT STRICT AS 'return 1;' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`, id.FullyQualifiedName(), cte.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestProcedures_CreateAndCallForSQL(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	defaultOpts := func() *CreateAndCallForSQLProcedureOptions {
		return &CreateAndCallForSQLProcedureOptions{
			Name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAndCallForSQLProcedureOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Name = NewAccountObjectIdentifier("")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForSQLProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{
				Columns: []ProcedureColumn{
					{
						ColumnName:     "name",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAndCallForSQLProcedureOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateAndCallForSQLProcedureOptions", "ProcedureDefinition"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = ProcedureReturns{
			Table: &ProcedureReturnsTable{},
		}
		opts.ProcedureDefinition = "3.141592654::FLOAT"
		opts.ProcedureName = id
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE () RETURNS TABLE () LANGUAGE SQL AS '3.141592654::FLOAT' CALL %s ()`, id.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []ProcedureArgument{
			{
				ArgName:      "message",
				ArgDataType:  "VARCHAR",
				DefaultValue: String("'test'"),
			},
		}
		opts.Returns = ProcedureReturns{
			ResultDataType: &ProcedureReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorStrict)
		opts.ProcedureDefinition = "3.141592654::FLOAT"
		cte := NewAccountObjectIdentifier("album_info_1976")
		opts.WithClauses = []ProcedureWithClause{
			{
				CteName:    cte,
				CteColumns: []string{"x", "y"},
				Statement:  "(select m.album_ID, m.album_name, b.band_name from music_albums)",
			},
		}
		opts.ProcedureName = id
		opts.ScriptingVariable = String(":ret")
		opts.CallArguments = []string{"1"}
		assertOptsValidAndSQLEquals(t, opts, `WITH %s AS PROCEDURE (message VARCHAR DEFAULT 'test') RETURNS FLOAT LANGUAGE SQL STRICT AS '3.141592654::FLOAT' , %s (x, y) AS (select m.album_ID, m.album_name, b.band_name from music_albums) CALL %s (1) INTO :ret`, id.FullyQualifiedName(), cte.FullyQualifiedName(), id.FullyQualifiedName())
	})
}
