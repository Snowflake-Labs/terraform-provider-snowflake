package sdk

import (
	"testing"
)

func TestExternalFunctions_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateExternalFunctionOptions {
		return &CreateExternalFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateExternalFunctionOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: must options", func(t *testing.T) {
		opts := defaultOpts()

		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalFunctionOptions", "ApiIntegration"))
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalFunctionOptions", "As"))

		opts = defaultOpts()
		opts.As = "as"
		integration := emptyAccountObjectIdentifier
		opts.ApiIntegration = &integration
		rt := emptySchemaObjectIdentifier
		opts.RequestTranslator = &rt
		st := emptySchemaObjectIdentifier
		opts.ResponseTranslator = &st
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalFunctionOptions", "ApiIntegration"))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateExternalFunctionOptions", "RequestTranslator"))
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateExternalFunctionOptions", "ResponseTranslator"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Secure = Bool(true)
		opts.Arguments = []ExternalFunctionArgument{
			{
				ArgName:     "id",
				ArgDataType: DataTypeNumber,
			},
			{
				ArgName:     "name",
				ArgDataType: DataTypeVARCHAR,
			},
		}
		opts.ResultDataType = DataTypeVARCHAR
		opts.ReturnNullValues = ReturnNullValuesPointer(ReturnNullValuesNotNull)
		opts.NullInputBehavior = NullInputBehaviorPointer(NullInputBehaviorCalledOnNullInput)
		opts.ReturnResultsBehavior = ReturnResultsBehaviorPointer(ReturnResultsBehaviorImmutable)
		opts.Comment = String("comment")
		integration := NewAccountObjectIdentifier("api_integration")
		opts.ApiIntegration = &integration
		opts.Headers = []ExternalFunctionHeader{
			{
				Name:  "header1",
				Value: "value1",
			},
			{
				Name:  "header2",
				Value: "value2",
			},
		}
		opts.ContextHeaders = []ExternalFunctionContextHeader{
			{
				ContextFunction: "CURRENT_ACCOUNT",
			},
			{
				ContextFunction: "CURRENT_USER",
			},
		}
		opts.MaxBatchRows = Int(100)
		opts.Compression = String("GZIP")
		rt := randomSchemaObjectIdentifier()
		opts.RequestTranslator = &rt
		rs := randomSchemaObjectIdentifier()
		opts.ResponseTranslator = &rs
		opts.As = "https://xyz.execute-api.us-west-2.amazonaws.com/prod/remote_echo"
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SECURE EXTERNAL FUNCTION %s (id NUMBER, name VARCHAR) RETURNS VARCHAR NOT NULL CALLED ON NULL INPUT IMMUTABLE COMMENT = 'comment' API_INTEGRATION = "api_integration" HEADERS = ('header1' = 'value1', 'header2' = 'value2') CONTEXT_HEADERS = (CURRENT_ACCOUNT, CURRENT_USER) MAX_BATCH_ROWS = 100 COMPRESSION = GZIP REQUEST_TRANSLATOR = %s RESPONSE_TRANSLATOR = %s AS 'https://xyz.execute-api.us-west-2.amazonaws.com/prod/remote_echo'`, id.FullyQualifiedName(), rt.FullyQualifiedName(), rs.FullyQualifiedName())
	})
}

func TestExternalFunctions_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterExternalFunctionOptions {
		return &AlterExternalFunctionOptions{
			name:              id,
			IfExists:          Bool(true),
			ArgumentDataTypes: []DataType{DataTypeVARCHAR, DataTypeNumber},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterExternalFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment opts.Unset.Headers opts.Unset.ContextHeaders opts.Unset.MaxBatchRows opts.Unset.Compression opts.Unset.Secure opts.Unset.RequestTranslator opts.Unset.ResponseTranslator] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalFunctionUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterExternalFunctionOptions.Unset", "Comment", "Headers", "ContextHeaders", "MaxBatchRows", "Compression", "Secure", "RequestTranslator", "ResponseTranslator"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			MaxBatchRows: Int(100),
		}
		opts.Unset = &ExternalFunctionUnset{
			MaxBatchRows: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalFunctionOptions", "Set", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.Set.ApiIntegration opts.Set.Headers opts.Set.ContextHeaders opts.Set.MaxBatchRows opts.Set.Compression opts.Set.RequestTranslator opts.Set.ResponseTranslator] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			MaxBatchRows: Int(100),
			Headers: []ExternalFunctionHeader{
				{
					Name:  "header1",
					Value: "value1",
				},
				{
					Name:  "header2",
					Value: "value2",
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalFunctionOptions.Set", "ApiIntegration", "Headers", "ContextHeaders", "MaxBatchRows", "Compression", "RequestTranslator", "ResponseTranslator"))
	})

	t.Run("alter: set api integration", func(t *testing.T) {
		opts := defaultOpts()
		integration := NewAccountObjectIdentifier("api_integration")
		opts.Set = &ExternalFunctionSet{
			ApiIntegration: &integration,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET API_INTEGRATION = "api_integration"`, id.FullyQualifiedName())
	})

	t.Run("alter: set headers", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			Headers: []ExternalFunctionHeader{
				{
					Name:  "header1",
					Value: "value1",
				},
				{
					Name:  "header2",
					Value: "value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET HEADERS = ('header1' = 'value1', 'header2' = 'value2')`, id.FullyQualifiedName())
	})

	t.Run("alter: set max batch rows", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			MaxBatchRows: Int(100),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET MAX_BATCH_ROWS = 100`, id.FullyQualifiedName())
	})

	t.Run("alter: set compression", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			Compression: String("GZIP"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET COMPRESSION = GZIP`, id.FullyQualifiedName())
	})

	t.Run("alter: set context headers", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalFunctionSet{
			ContextHeaders: []ExternalFunctionContextHeader{
				{
					ContextFunction: "CURRENT_ACCOUNT",
				},
				{
					ContextFunction: "CURRENT_USER",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET CONTEXT_HEADERS = (CURRENT_ACCOUNT, CURRENT_USER)`, id.FullyQualifiedName())
	})

	t.Run("alter: set request translator", func(t *testing.T) {
		opts := defaultOpts()
		rt := randomSchemaObjectIdentifier()
		opts.Set = &ExternalFunctionSet{
			RequestTranslator: &rt,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET REQUEST_TRANSLATOR = %s`, id.FullyQualifiedName(), rt.FullyQualifiedName())
	})

	t.Run("alter: set response translator", func(t *testing.T) {
		opts := defaultOpts()
		st := randomSchemaObjectIdentifier()
		opts.Set = &ExternalFunctionSet{
			ResponseTranslator: &st,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) SET RESPONSE_TRANSLATOR = %s`, id.FullyQualifiedName(), st.FullyQualifiedName())
	})

	t.Run("alter: unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = []DataType{DataTypeVARCHAR, DataTypeNumber}
		opts.Unset = &ExternalFunctionUnset{
			Comment:            Bool(true),
			Headers:            Bool(true),
			ContextHeaders:     Bool(true),
			MaxBatchRows:       Bool(true),
			Compression:        Bool(true),
			Secure:             Bool(true),
			RequestTranslator:  Bool(true),
			ResponseTranslator: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s (VARCHAR, NUMBER) UNSET COMMENT, HEADERS, CONTEXT_HEADERS, MAX_BATCH_ROWS, COMPRESSION, SECURE, REQUEST_TRANSLATOR, RESPONSE_TRANSLATOR`, id.FullyQualifiedName())
	})

	t.Run("alter: unset with no arguments", func(t *testing.T) {
		opts := defaultOpts()
		opts.ArgumentDataTypes = nil
		opts.Unset = &ExternalFunctionUnset{
			Comment:            Bool(true),
			Headers:            Bool(true),
			ContextHeaders:     Bool(true),
			MaxBatchRows:       Bool(true),
			Compression:        Bool(true),
			Secure:             Bool(true),
			RequestTranslator:  Bool(true),
			ResponseTranslator: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FUNCTION IF EXISTS %s () UNSET COMMENT, HEADERS, CONTEXT_HEADERS, MAX_BATCH_ROWS, COMPRESSION, SECURE, REQUEST_TRANSLATOR, RESPONSE_TRANSLATOR`, id.FullyQualifiedName())
	})
}

func TestExternalFunctions_Show(t *testing.T) {
	defaultOpts := func() *ShowExternalFunctionOptions {
		return &ShowExternalFunctionOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowExternalFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW EXTERNAL FUNCTIONS`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW EXTERNAL FUNCTIONS LIKE 'pattern'`)
	})

	t.Run("show with in", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := defaultOpts()
		opts.In = &In{
			Schema: id,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW EXTERNAL FUNCTIONS IN SCHEMA %s`, id.FullyQualifiedName())
	})
}

func TestExternalFunctions_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeExternalFunctionOptions {
		return &DescribeExternalFunctionOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeExternalFunctionOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
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
