package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
)

func TestFunctions_CreateFunctionForJava(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateFunctionForJavaFunctionOptions {
		return &CreateFunctionForJavaFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateFunctionForJavaFunctionOptions)(nil)
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
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:     "id",
				ArgDataType: DataTypeNumber,
			},
			{
				ArgName:     "name",
				ArgDataType: DataTypeVARCHAR,
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = &FunctionReturns{
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
		returnNullValues := FunctionReturnNullValuesNotNull
		opts.ReturnNullValues = &returnNullValues
		nullInputBehavior := FunctionNullInputBehaviorCalledOnNullInput
		opts.NullInputBehavior = &nullInputBehavior
		returnResultsBehavior := FunctionReturnResultsBehaviorImmutable
		opts.ReturnResultsBehavior = &returnResultsBehavior
		opts.RuntimeVersion = String("2.0")
		opts.Comment = String("comment")
		opts.Imports = []FunctionImports{
			{
				Import: "@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar",
			},
		}
		opts.Packages = []FunctionPackages{
			{
				Package: "com.snowflake:snowpark:1.2.0",
			},
		}
		opts.Handler = "TestFunc.echoVarchar"
		opts.ExternalAccessIntegrations = []AccountObjectIdentifier{
			NewAccountObjectIdentifier("ext_integration"),
		}
		opts.Secrets = []FunctionSecret{
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
		opts.FunctionDefinition = "return id + name;"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s (id NUMBER, name VARCHAR) COPY GRANTS RETURNS TABLE (country_code VARCHAR, country_name VARCHAR) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) TARGET_PATH = '@~/testfunc.jar' AS 'return id + name;'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateFunctionForJavascript(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateFunctionForJavascriptFunctionOptions {
		return &CreateFunctionForJavascriptFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateFunctionForJavascriptFunctionOptions)(nil)
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
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:     "d",
				ArgDataType: DataTypeFloat,
			},
		}
		opts.CopyGrants = Bool(true)
		float := DataTypeFloat
		opts.Returns = &FunctionReturns{
			ResultDataType: &float,
		}
		returnNullValues := FunctionReturnNullValuesNotNull
		opts.ReturnNullValues = &returnNullValues
		nullInputBehavior := FunctionNullInputBehaviorCalledOnNullInput
		opts.NullInputBehavior = &nullInputBehavior
		returnResultsBehavior := FunctionReturnResultsBehaviorImmutable
		opts.ReturnResultsBehavior = &returnResultsBehavior
		opts.Comment = String("comment")
		opts.FunctionDefinition = "return 1;"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s (d FLOAT) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS 'return 1;'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateFunctionForPython(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateFunctionForPythonFunctionOptions {
		return &CreateFunctionForPythonFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateFunctionForPythonFunctionOptions)(nil)
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
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:     "i",
				ArgDataType: DataTypeNumber,
			},
		}
		opts.CopyGrants = Bool(true)
		varint := DataTypeVariant
		opts.Returns = &FunctionReturns{
			ResultDataType: &varint,
		}
		returnNullValues := FunctionReturnNullValuesNotNull
		opts.ReturnNullValues = &returnNullValues
		nullInputBehavior := FunctionNullInputBehaviorCalledOnNullInput
		opts.NullInputBehavior = &nullInputBehavior
		returnResultsBehavior := FunctionReturnResultsBehaviorImmutable
		opts.ReturnResultsBehavior = &returnResultsBehavior
		opts.RuntimeVersion = "3.8"
		opts.Comment = String("comment")
		opts.Imports = []FunctionImports{
			{
				Import: "numpy",
			},
			{
				Import: "pandas",
			},
		}
		opts.Packages = []FunctionPackages{
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
		opts.Secrets = []FunctionSecret{
			{
				SecretVariableName: "variable1",
				SecretName:         "name1",
			},
			{
				SecretVariableName: "variable2",
				SecretName:         "name2",
			},
		}
		opts.FunctionDefinition = "import numpy as np"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s (i NUMBER) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.8' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = name1, 'variable2' = name2) AS 'import numpy as np'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateFunctionForScala(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateFunctionForScalaFunctionOptions {
		return &CreateFunctionForScalaFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateFunctionForScalaFunctionOptions)(nil)
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
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:     "x",
				ArgDataType: DataTypeVARCHAR,
			},
		}
		opts.CopyGrants = Bool(true)
		varchar := DataTypeVARCHAR
		opts.Returns = &FunctionReturns{
			ResultDataType: &varchar,
		}
		returnNullValues := FunctionReturnNullValuesNotNull
		opts.ReturnNullValues = &returnNullValues
		nullInputBehavior := FunctionNullInputBehaviorCalledOnNullInput
		opts.NullInputBehavior = &nullInputBehavior
		returnResultsBehavior := FunctionReturnResultsBehaviorImmutable
		opts.ReturnResultsBehavior = &returnResultsBehavior
		opts.RuntimeVersion = String("2.0")
		opts.Comment = String("comment")
		opts.Imports = []FunctionImports{
			{
				Import: "@udf_libs/echohandler.jar",
			},
		}
		opts.Handler ="Echo.echoVarchar"
		opts.FunctionDefinition = "return x"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s (x VARCHAR) COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS 'return x'`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateFunctionForSQL(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateFunctionForSQLFunctionOptions {
		return &CreateFunctionForSQLFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateFunctionForSQLFunctionOptions)(nil)
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
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.IfNotExists = Bool(true)
		opts.CopyGrants = Bool(true)
		dt := DataTypeFloat
		opts.Returns = &FunctionReturns{
			ResultDataType: &dt,
		}
		returnNullValues := FunctionReturnNullValuesNotNull
		opts.ReturnNullValues = &returnNullValues
		returnResultsBehavior := FunctionReturnResultsBehaviorImmutable
		opts.ReturnResultsBehavior = &returnResultsBehavior
		opts.Memoizable = Bool(true)
		opts.Comment = String("comment")
		opts.FunctionDefinition = "3.141592654::FLOAT"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION IF NOT EXISTS %s COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS '3.141592654::FLOAT'`, id.FullyQualifiedName())
	})
}

func TestFunctions_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

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

	t.Run("all options", func(t *testing.T) {
		opts := &DropFunctionOptions{
			name: id,
		}
		opts.IfExists = Bool(true)
		opts.ArgumentTypes = []FunctionArgumentType{
			{
				ArgDataType: DataTypeVARCHAR,
			},
			{
				ArgDataType: DataTypeNumber,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP FUNCTION IF EXISTS %s (VARCHAR, NUMBER)`, id.FullyQualifiedName())
	})
}

func TestFunctions_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterFunctionOptions {
		return &AlterFunctionOptions{
			name:     id,
			IfExists: Bool(true),
			ArgumentTypes: []FunctionArgumentType{
				{
					ArgDataType: DataTypeVARCHAR,
				},
				{
					ArgDataType: DataTypeNumber,
				},
			},
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

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), random.StringN(12))
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: set log level", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			LogLevel: String("DEBUG"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET LOG_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set trace level", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			TraceLevel: String("DEBUG"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET TRACE_LEVEL = 'DEBUG'`, id.FullyQualifiedName())
	})

	t.Run("alter: set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: set secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			Secure: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: unset secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &FunctionUnset{
			Secure: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &FunctionUnset{
			Comment: Bool(true),
		}
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

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
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
	id := RandomSchemaObjectIdentifier()

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

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE FUNCTION %s`, id.FullyQualifiedName())
	})
}
