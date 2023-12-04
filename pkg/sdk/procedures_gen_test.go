package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
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
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterProcedureOptions", "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "UnsetComment", "SetTags", "UnsetTags", "ExecuteAs"))
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

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PROCEDURE %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}
