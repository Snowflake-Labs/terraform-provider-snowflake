package sdk

import (
	"testing"
)

func TestFunctions_CreateForJava(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForJavaFunctionOptions {
		return &CreateForJavaFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForJavaFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: function definition", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetPath = String("@~/testfunc.jar")
		opts.Packages = []FunctionPackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("TARGET_PATH must be nil when AS is nil"))
		assertOptsInvalidJoinedErrors(t, opts, NewError("PACKAGES must be empty when AS is nil"))
		assertOptsInvalidJoinedErrors(t, opts, NewError("IMPORTS must not be empty when AS is nil"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavaFunctionOptions", "Handler"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
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
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: DataTypeVARCHAR,
					},
					{
						ColumnName:     "country_name",
						ColumnDataType: DataTypeVARCHAR,
					},
				},
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.RuntimeVersion = String("2.0")
		opts.Comment = String("comment")
		opts.Imports = []FunctionImport{
			{
				Import: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar",
			},
		}
		opts.Packages = []FunctionPackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
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
		opts.FunctionDefinition = String("return id + name;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (id NUMBER, name VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE (country_code VARCHAR, country_name VARCHAR) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) TARGET_PATH = '@~/testfunc.jar' AS 'return id + name;'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateForJavascript(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForJavascriptFunctionOptions {
		return &CreateForJavascriptFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForJavascriptFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavascriptFunctionOptions", "FunctionDefinition"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "d",
				ArgDataType:  DataTypeFloat,
				DefaultValue: String("1.0"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Comment = String("comment")
		opts.FunctionDefinition = "return 1;"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (d FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS 'return 1;'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateForPython(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForPythonFunctionOptions {
		return &CreateForPythonFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForPythonFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonFunctionOptions", "RuntimeVersion"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonFunctionOptions", "Handler"))
	})

	t.Run("validation: function definition", func(t *testing.T) {
		opts := defaultOpts()
		opts.Packages = []FunctionPackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("IMPORTS must not be empty when AS is nil"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "i",
				ArgDataType:  DataTypeNumber,
				DefaultValue: String("1"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeVariant,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.RuntimeVersion = "3.8"
		opts.Comment = String("comment")
		opts.Imports = []FunctionImport{
			{
				Import: "numpy",
			},
			{
				Import: "pandas",
			},
		}
		opts.Packages = []FunctionPackage{
			{
				Package: "numpy",
			},
			{
				Package: "pandas",
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
		opts.FunctionDefinition = String("import numpy as np")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (i NUMBER DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.8' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) AS 'import numpy as np'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateForScala(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForScalaFunctionOptions {
		return &CreateForScalaFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForScalaFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: function definition", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetPath = String("@~/testfunc.jar")
		opts.Packages = []FunctionPackage{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, NewError("TARGET_PATH must be nil when AS is nil"))
		assertOptsInvalidJoinedErrors(t, opts, NewError("PACKAGES must be empty when AS is nil"))
		assertOptsInvalidJoinedErrors(t, opts, NewError("IMPORTS must not be empty when AS is nil"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResultDataType = DataTypeVARCHAR
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForScalaFunctionOptions", "Handler"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "x",
				ArgDataType:  DataTypeVARCHAR,
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.ResultDataType = DataTypeVARCHAR
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.RuntimeVersion = String("2.0")
		opts.Comment = String("comment")
		opts.Imports = []FunctionImport{
			{
				Import: "@udf_libs/echohandler.jar",
			},
		}
		opts.Handler = "Echo.echoVarchar"
		opts.FunctionDefinition = String("return x")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (x VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS 'return x'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateForSQL(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForSQLFunctionOptions {
		return &CreateForSQLFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForSQLFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: returns", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: options are missing", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeVARCHAR,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForSQLFunctionOptions", "FunctionDefinition"))
	})

	t.Run("create with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		opts.FunctionDefinition = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE FUNCTION %s () RETURNS FLOAT AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "message",
				ArgDataType:  "VARCHAR",
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: DataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Memoizable = Bool(true)
		opts.Comment = String("comment")
		opts.FunctionDefinition = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (message VARCHAR DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
	})
}

func TestFunctions_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropFunctionOptions {
		return &DropFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP FUNCTION %s ()`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &DropFunctionOptions{
			name: id,
		}
		opts.IfExists = Bool(true)
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		assertOptsValidAndSQLEquals(t, opts, `DROP FUNCTION IF EXISTS %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}

func TestFunctions_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterFunctionOptions {
		return &AlterFunctionOptions{
			name:              id,
			IfExists:          Bool(true),
			ArgumentDataTypes: []DataType{DataTypeVARCHAR, DataTypeNumber},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFunctionOptions", "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "SetSecure", "UnsetLogLevel", "UnsetTraceLevel", "UnsetSecure", "UnsetComment", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetLogLevel = String("DEBUG")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFunctionOptions", "RenameTo", "SetComment", "SetLogLevel", "SetTraceLevel", "SetSecure", "UnsetLogLevel", "UnsetTraceLevel", "UnsetSecure", "UnsetComment", "SetTags", "UnsetTags"))
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: set log level with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = nil
		opts.SetLogLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s () SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set log level", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetLogLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set trace level", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTraceLevel = String("DEBUG")
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET TRACE_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: set secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: unset log level", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetLogLevel = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET LOG_LEVEL`, id.FullyQualifiedName())
	})

	t.Run("alter: unset trace level", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTraceLevel = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET TRACE_LEVEL`, id.FullyQualifiedName())
	})

	t.Run("alter: unset secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestFunctions_Show(t *testing.T) {
	defaultOpts := func() *ShowFunctionOptions {
		return &ShowFunctionOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER FUNCTIONS`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER FUNCTIONS LIKE 'pattern'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER FUNCTIONS IN ACCOUNT`)
	})
}

func TestFunctions_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeFunctionOptions {
		return &DescribeFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("no arguments", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE FUNCTION %s ()`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE FUNCTION %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}
