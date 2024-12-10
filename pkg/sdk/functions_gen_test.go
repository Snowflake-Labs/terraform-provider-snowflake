package sdk

import (
	"fmt"
	"testing"
)

func wrapFunctionDefinition(def string) string {
	return fmt.Sprintf(`$$%s$$`, def)
}

func TestFunctions_CreateForJava(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForJavaFunctionOptions {
		return &CreateForJavaFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForJavaFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.Handler] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVarchar,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavaFunctionOptions", "Handler"))
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateForJavaFunctionOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat},
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavaFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
				ResultDataType:    dataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat, ColumnDataType: dataTypeFloat},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat},
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
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
		assertOptsInvalidJoinedErrors(t, opts, NewError("IMPORTS must not be empty when AS is nil"))
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("all options - old data types", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:        "id",
				ArgDataTypeOld: DataTypeNumber,
			},
			{
				ArgName:        "name",
				ArgDataTypeOld: DataTypeVARCHAR,
				DefaultValue:   String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{
						ColumnName:        "country_code",
						ColumnDataTypeOld: DataTypeVARCHAR,
					},
					{
						ColumnName:        "country_name",
						ColumnDataTypeOld: DataTypeVARCHAR,
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
		opts.Secrets = []SecretReference{
			{
				VariableName: "variable1",
				Name:         secretId,
			},
			{
				VariableName: "variable2",
				Name:         secretId2,
			},
		}
		opts.TargetPath = String("@~/testfunc.jar")
		opts.FunctionDefinition = String(wrapFunctionDefinition("return id + name;"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("id" NUMBER, "name" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR, "country_name" VARCHAR) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`, id.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:     "id",
				ArgDataType: dataTypeNumber,
			},
			{
				ArgName:      "name",
				ArgDataType:  dataTypeVarchar,
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{
						ColumnName:     "country_code",
						ColumnDataType: dataTypeVarchar,
					},
					{
						ColumnName:     "country_name",
						ColumnDataType: dataTypeVarchar,
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
		opts.Secrets = []SecretReference{
			{
				VariableName: "variable1",
				Name:         secretId,
			},
			{
				VariableName: "variable2",
				Name:         secretId2,
			},
		}
		opts.TargetPath = String("@~/testfunc.jar")
		opts.FunctionDefinition = String(wrapFunctionDefinition("return id + name;"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("id" NUMBER(36, 2), "name" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS TABLE ("country_code" VARCHAR(100), "country_name" VARCHAR(100)) NOT NULL LANGUAGE JAVA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@~/my_decrement_udf_package_dir/my_decrement_udf_jar.jar') PACKAGES = ('com.snowflake:snowpark:1.2.0') HANDLER = 'TestFunc.echoVarchar' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) TARGET_PATH = '@~/testfunc.jar' AS $$return id + name;$$`, id.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName())
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

	t.Run("validation: [opts.FunctionDefinition] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVarchar,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForJavascriptFunctionOptions", "FunctionDefinition"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat},
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
			Table:          &FunctionReturnsTable{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForJavascriptFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
				ResultDataType:    dataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat, ColumnDataType: dataTypeFloat},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat},
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("all options - old data types", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:        "d",
				ArgDataTypeOld: DataTypeFloat,
				DefaultValue:   String("1.0"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Comment = String("comment")
		opts.FunctionDefinition = wrapFunctionDefinition("return 1;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "d",
				ArgDataType:  dataTypeFloat,
				DefaultValue: String("1.0"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Comment = String("comment")
		opts.FunctionDefinition = wrapFunctionDefinition("return 1;")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("d" FLOAT DEFAULT 1.0) COPY GRANTS RETURNS FLOAT NOT NULL LANGUAGE JAVASCRIPT CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' AS $$return 1;$$`, id.FullyQualifiedName())
	})
}

func TestFunctions_CreateForPython(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateForPythonFunctionOptions {
		return &CreateForPythonFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateForPythonFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.RuntimeVersion] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVarchar,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonFunctionOptions", "RuntimeVersion"))
	})

	t.Run("validation: [opts.Handler] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVarchar,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForPythonFunctionOptions", "Handler"))
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateForPythonFunctionOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - one correct, one incorrect", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat},
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
			Table:          &FunctionReturnsTable{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForPythonFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
				ResultDataType:    dataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat},
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat, ColumnDataType: dataTypeFloat},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
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

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("all options - old data types", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:        "i",
				ArgDataTypeOld: DataTypeNumber,
				DefaultValue:   String("1"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeVariant,
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
		opts.Secrets = []SecretReference{
			{
				VariableName: "variable1",
				Name:         secretId,
			},
			{
				VariableName: "variable2",
				Name:         secretId2,
			},
		}
		opts.FunctionDefinition = String(wrapFunctionDefinition("import numpy as np"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("i" NUMBER DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.8' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`, id.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "i",
				ArgDataType:  dataTypeNumber,
				DefaultValue: String("1"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVariant,
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
		opts.Secrets = []SecretReference{
			{
				VariableName: "variable1",
				Name:         secretId,
			},
			{
				VariableName: "variable2",
				Name:         secretId2,
			},
		}
		opts.FunctionDefinition = String(wrapFunctionDefinition("import numpy as np"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("i" NUMBER(36, 2) DEFAULT 1) COPY GRANTS RETURNS VARIANT NOT NULL LANGUAGE PYTHON CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '3.8' COMMENT = 'comment' IMPORTS = ('numpy', 'pandas') PACKAGES = ('numpy', 'pandas') HANDLER = 'udf' EXTERNAL_ACCESS_INTEGRATIONS = ("ext_integration") SECRETS = ('variable1' = %s, 'variable2' = %s) AS $$import numpy as np$$`, id.FullyQualifiedName(), secretId.FullyQualifiedName(), secretId2.FullyQualifiedName())
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

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.Handler] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResultDataType = dataTypeVarchar
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForScalaFunctionOptions", "Handler"))
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateForScalaFunctionOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: exactly one field from [opts.ResultDataTypeOld opts.ResultDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaFunctionOptions", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.ResultDataTypeOld opts.ResultDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResultDataTypeOld = DataTypeFloat
		opts.ResultDataType = dataTypeFloat
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaFunctionOptions", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat},
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForScalaFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
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
		assertOptsInvalidJoinedErrors(t, opts, NewError("IMPORTS must not be empty when AS is nil"))
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("all options - old data types", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:        "x",
				ArgDataTypeOld: DataTypeVARCHAR,
				DefaultValue:   String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.ResultDataTypeOld = DataTypeVARCHAR
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.RuntimeVersion = "2.0"
		opts.Comment = String("comment")
		opts.Imports = []FunctionImport{
			{
				Import: "@udf_libs/echohandler.jar",
			},
		}
		opts.Handler = "Echo.echoVarchar"
		opts.FunctionDefinition = String(wrapFunctionDefinition("return x"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("x" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS VARCHAR NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "x",
				ArgDataType:  dataTypeVarchar,
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.ResultDataType = dataTypeVarchar
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.RuntimeVersion = "2.0"
		opts.Comment = String("comment")
		opts.Imports = []FunctionImport{
			{
				Import: "@udf_libs/echohandler.jar",
			},
		}
		opts.Handler = "Echo.echoVarchar"
		opts.FunctionDefinition = String(wrapFunctionDefinition("return x"))
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("x" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS VARCHAR(100) NOT NULL LANGUAGE SCALA CALLED ON NULL INPUT IMMUTABLE RUNTIME_VERSION = '2.0' COMMENT = 'comment' IMPORTS = ('@udf_libs/echohandler.jar') HANDLER = 'Echo.echoVarchar' AS $$return x$$`, id.FullyQualifiedName())
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

	t.Run("validation: [opts.FunctionDefinition] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeVarchar,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateForSQLFunctionOptions", "FunctionDefinition"))
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Arguments.ArgDataTypeOld opts.Arguments.ArgDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Arguments = []FunctionArgument{
			{ArgName: "arg", ArgDataTypeOld: DataTypeFloat, ArgDataType: dataTypeFloat},
			{ArgName: "arg"},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Arguments", "ArgDataTypeOld", "ArgDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType opts.Returns.Table] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
			Table:          &FunctionReturnsTable{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns", "ResultDataType", "Table"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.ResultDataType.ResultDataTypeOld opts.Returns.ResultDataType.ResultDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
				ResultDataType:    dataTypeFloat,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.ResultDataType", "ResultDataTypeOld", "ResultDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat},
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg", ColumnDataTypeOld: DataTypeFloat, ColumnDataType: dataTypeFloat},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("validation: exactly one field from [opts.Returns.Table.Columns.ColumnDataTypeOld opts.Returns.Table.Columns.ColumnDataType] should be present - one valid, one invalid", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			Table: &FunctionReturnsTable{
				Columns: []FunctionColumn{
					{ColumnName: "arg"},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateForSQLFunctionOptions.Returns.Table.Columns", "ColumnDataTypeOld", "ColumnDataType"))
	})

	t.Run("create with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeFloat,
			},
		}
		opts.FunctionDefinition = wrapFunctionDefinition("3.141592654::FLOAT")
		assertOptsValidAndSQLEquals(t, opts, `CREATE FUNCTION %s () RETURNS FLOAT AS $$3.141592654::FLOAT$$`, id.FullyQualifiedName())
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("all options - old data types", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:        "message",
				ArgDataTypeOld: "VARCHAR",
				DefaultValue:   String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataTypeOld: DataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Memoizable = Bool(true)
		opts.Comment = String("comment")
		opts.FunctionDefinition = wrapFunctionDefinition("3.141592654::FLOAT")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("message" VARCHAR DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Temporary = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []FunctionArgument{
			{
				ArgName:      "message",
				ArgDataType:  dataTypeVarchar,
				DefaultValue: String("'test'"),
			},
		}
		opts.CopyGrants = Bool(true)
		opts.Returns = FunctionReturns{
			ResultDataType: &FunctionReturnsResultDataType{
				ResultDataType: dataTypeFloat,
			},
		}
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Memoizable = Bool(true)
		opts.Comment = String("comment")
		opts.FunctionDefinition = wrapFunctionDefinition("3.141592654::FLOAT")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY SECURE FUNCTION %s ("message" VARCHAR(100) DEFAULT 'test') COPY GRANTS RETURNS FLOAT NOT NULL IMMUTABLE MEMOIZABLE COMMENT = 'comment' AS $$3.141592654::FLOAT$$`, id.FullyQualifiedName())
	})
}

func TestFunctions_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR, DataTypeNumber)
	secretId := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterFunctionOptions {
		return &AlterFunctionOptions{
			name:     id,
			IfExists: Bool(true),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifierWithArguments
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		target := emptySchemaObjectIdentifier
		opts.RenameTo = &target
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.Set opts.Unset opts.SetSecure opts.UnsetSecure opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFunctionOptions", "RenameTo", "Set", "Unset", "SetSecure", "UnsetSecure", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.RenameTo opts.Set opts.Unset opts.SetSecure opts.UnsetSecure opts.SetTags opts.UnsetTags] should be present - two present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{}
		opts.Unset = &FunctionUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFunctionOptions", "RenameTo", "Set", "Unset", "SetSecure", "UnsetSecure", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Comment opts.Set.ExternalAccessIntegrations opts.Set.SecretsList opts.Set.EnableConsoleOutput opts.Set.LogLevel opts.Set.MetricLevel opts.Set.TraceLevel] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterFunctionOptions.Set", "Comment", "ExternalAccessIntegrations", "SecretsList", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment opts.Unset.ExternalAccessIntegrations opts.Unset.EnableConsoleOutput opts.Unset.LogLevel opts.Unset.MetricLevel opts.Unset.TraceLevel] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &FunctionUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterFunctionOptions.Unset", "Comment", "ExternalAccessIntegrations", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"))
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			Comment:    String("comment"),
			TraceLevel: Pointer(TraceLevelOff),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s SET COMMENT = 'comment', TRACE_LEVEL = 'OFF'`, id.FullyQualifiedName())
	})

	t.Run("alter: set empty secrets", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			SecretsList: &SecretsList{},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s SET SECRETS = ()`, id.FullyQualifiedName())
	})

	t.Run("alter: set non-empty secrets", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &FunctionSet{
			SecretsList: &SecretsList{
				[]SecretReference{
					{VariableName: "abc", Name: secretId},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s SET SECRETS = ('abc' = %s)`, id.FullyQualifiedName(), secretId.FullyQualifiedName())
	})

	t.Run("alter: unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &FunctionUnset{
			Comment:    Bool(true),
			TraceLevel: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s UNSET COMMENT, TRACE_LEVEL`, id.FullyQualifiedName())
	})

	t.Run("alter: set secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s SET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: unset secure", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetSecure = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s UNSET SECURE`, id.FullyQualifiedName())
	})

	t.Run("alter: set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestFunctions_Drop(t *testing.T) {
	noArgsId := randomSchemaObjectIdentifierWithArguments()
	id := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR, DataTypeNumber)

	defaultOpts := func() *DropFunctionOptions {
		return &DropFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifierWithArguments
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = noArgsId
		assertOptsValidAndSQLEquals(t, opts, `DROP FUNCTION %s`, noArgsId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &DropFunctionOptions{
			name: id,
		}
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP FUNCTION IF EXISTS %s`, id.FullyQualifiedName())
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
		opts.In = &ExtendedIn{
			In: In{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER FUNCTIONS IN ACCOUNT`)
	})
}

func TestFunctions_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR, DataTypeNumber)

	defaultOpts := func() *DescribeFunctionOptions {
		return &DescribeFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifierWithArguments
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE FUNCTION %s`, id.FullyQualifiedName())
	})
}
