package testint

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
todo: add tests for:
  - creating functions with different languages (java, javascript, python, scala, sql) from stages  using [ TARGET_PATH = '<stage_path_and_file_name_to_write>' ]
  - execute and execute-immediate for scripting https://docs.snowflake.com/en/sql-reference/sql/execute-immediate
*/

func TestInt_CreateFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	t.Run("create function for Java", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR)

		definition := testClientHelper().Function.SampleJavaDefinition(t)
		target := fmt.Sprintf("@~/tf-%d.jar", time.Now().Unix())
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVARCHAR)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("x", nil).WithDefaultValue("'abc'").WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForJavaFunctionRequest(id.SchemaObjectId(), *returns, "TestFunc.echoVarchar").
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput)).
			WithTargetPath(target).
			WithFunctionDefinition(definition)
		err := client.Functions.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "JAVA", function.Language)
	})

	t.Run("create function for Javascript", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)

		definition := testClientHelper().Function.SampleJavaScriptDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("d", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request := sdk.NewCreateForJavascriptFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithNullInputBehavior(*sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput))
		err := client.Functions.CreateForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "JAVASCRIPT", function.Language)
	})

	t.Run("create function for Python", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeNumber)

		definition := testClientHelper().Function.SamplePythonDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVariant)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("i", nil).WithArgDataTypeOld(sdk.DataTypeNumber)
		request := sdk.NewCreateForPythonFunctionRequest(id.SchemaObjectId(), *returns, "3.8", "dump").
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinition(definition)
		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "PYTHON", function.Language)
	})

	t.Run("create function for Scala", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeVARCHAR)

		definition := testClientHelper().Function.SampleScalaDefinition(t)
		argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), nil, "Echo.echoVarchar").
			WithResultDataTypeOld(sdk.DataTypeVARCHAR).
			WithOrReplace(true).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithRuntimeVersion("2.12").
			WithFunctionDefinition(definition)
		err := client.Functions.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SCALA", function.Language)
	})

	t.Run("create function for SQL", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithOrReplace(true).
			WithComment("comment")
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})

	t.Run("create function for SQL with no arguments", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true).
			WithComment("comment")
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})
}

func TestInt_OtherFunctions(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertFunction := func(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments, secure bool, withArguments bool) {
		t.Helper()

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.NotEmpty(t, function.CreatedOn)
		assert.Equal(t, id.Name(), function.Name)
		assert.Equal(t, false, function.IsBuiltin)
		assert.Equal(t, false, function.IsAggregate)
		assert.Equal(t, false, function.IsAnsi)
		if withArguments {
			assert.Equal(t, 1, function.MinNumArguments)
			assert.Equal(t, 1, function.MaxNumArguments)
		} else {
			assert.Equal(t, 0, function.MinNumArguments)
			assert.Equal(t, 0, function.MaxNumArguments)
		}
		assert.NotEmpty(t, function.ArgumentsRaw)
		assert.NotEmpty(t, function.ArgumentsOld)
		assert.NotEmpty(t, function.Description)
		assert.NotEmpty(t, function.CatalogName)
		assert.Equal(t, false, function.IsTableFunction)
		assert.Equal(t, false, function.ValidForClustering)
		assert.Equal(t, secure, function.IsSecure)
		assert.Equal(t, false, function.IsExternalFunction)
		assert.Equal(t, "SQL", function.Language)
		assert.Equal(t, false, function.IsMemoizable)
	}

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, cleanup bool, withArguments bool) *sdk.Function {
		t.Helper()
		var id sdk.SchemaObjectIdentifierWithArguments
		if withArguments {
			id = testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		} else {
			id = testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()
		}

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true)
		if withArguments {
			argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
			request = request.WithArguments([]sdk.FunctionArgumentRequest{*argument})
		}
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupFunctionHandle(id))
		}
		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		return function
	}

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false, true)

		id := f.ID()
		nid := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithRenameTo(nid.SchemaObjectId()))
		if err != nil {
			t.Cleanup(cleanupFunctionHandle(id))
		} else {
			t.Cleanup(cleanupFunctionHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Functions.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Functions.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter function: set log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetLogLevel(string(sdk.LogLevelDebug)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetLogLevel(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetTraceLevel(string(sdk.TraceLevelAlways)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetTraceLevel(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetComment("test comment"))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetComment(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, true, true)
	})

	t.Run("alter function: set secure with no arguments", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)
		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithSetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, true, true)
	})

	t.Run("alter function: unset secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id).WithUnsetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("show function for SQL: without like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true, true)
		f2 := createFunctionForSQLHandle(t, true, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)

		require.Contains(t, functions, *f1)
		require.Contains(t, functions, *f2)
	})

	t.Run("show function for SQL: with like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true, true)
		f2 := createFunctionForSQLHandle(t, true, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(sdk.Like{Pattern: &f1.Name}))
		require.NoError(t, err)

		require.Equal(t, 1, len(functions))
		require.Contains(t, functions, *f1)
		require.NotContains(t, functions, *f2)
	})

	t.Run("show function for SQL: no matches", func(t *testing.T) {
		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(sdk.Like{Pattern: sdk.String("non-existing-id-pattern")}))
		require.NoError(t, err)
		require.Equal(t, 0, len(functions))
	})

	t.Run("describe function for SQL", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		details, err := client.Functions.Describe(ctx, f.ID())
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, "SQL", pairs["language"])
		require.Equal(t, "FLOAT", pairs["returns"])
		require.Equal(t, "3.141592654::FLOAT", pairs["body"])
		require.Equal(t, "(X FLOAT)", pairs["signature"])
	})

	t.Run("describe function for SQL: no arguments", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, false)

		details, err := client.Functions.Describe(ctx, f.ID())
		require.NoError(t, err)
		pairs := make(map[string]string)
		for _, detail := range details {
			pairs[detail.Property] = detail.Value
		}
		require.Equal(t, "SQL", pairs["language"])
		require.Equal(t, "FLOAT", pairs["returns"])
		require.Equal(t, "3.141592654::FLOAT", pairs["body"])
		require.Equal(t, "()", pairs["signature"])
	})
}

func TestInt_FunctionsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifierWithArguments) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) {
		t.Helper()

		definition := testClientHelper().Function.SampleSqlDefinition(t)
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).WithOrReplace(true)

		argument := sdk.NewFunctionArgumentRequest("x", nil).WithArgDataTypeOld(sdk.DataTypeFloat)
		request = request.WithArguments([]sdk.FunctionArgumentRequest{*argument})
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeFloat)
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(id1.Name(), schema.ID(), sdk.DataTypeFloat)

		createFunctionForSQLHandle(t, id1)
		createFunctionForSQLHandle(t, id2)

		e1, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)

		e1Id := e1.ID()
		require.NoError(t, err)
		require.Equal(t, id1, e1Id)

		e2, err := client.Functions.ShowByID(ctx, id2)
		require.NoError(t, err)

		e2Id := e2.ID()
		require.NoError(t, err)
		require.Equal(t, id2, e2Id)
	})

	t.Run("show function by id - different name, same arguments", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments(sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		e := testClientHelper().Function.CreateWithIdentifier(t, id1)
		testClientHelper().Function.CreateWithIdentifier(t, id2)

		es, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	t.Run("show function by id - same name, different arguments", func(t *testing.T) {
		name := testClientHelper().Ids.Alpha()
		id1 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeFloat, sdk.DataTypeVARCHAR)
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierWithArgumentsInSchema(name, testClientHelper().Ids.SchemaId(), sdk.DataTypeInt, sdk.DataTypeVARCHAR)
		e := testClientHelper().Function.CreateWithIdentifier(t, id1)
		testClientHelper().Function.CreateWithIdentifier(t, id2)

		es, err := client.Functions.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, *e, *es)
	})

	// TODO [SNOW-1348103]: remove with old function removal for V1
	t.Run("function returns non detailed data types of arguments - old data types", func(t *testing.T) {
		// This test proves that every detailed data types (e.g. VARCHAR(20) and NUMBER(10, 0)) are generalized
		// on Snowflake side (to e.g. VARCHAR and NUMBER) and that sdk.ToDataType mapping function maps detailed types
		// correctly to their generalized counterparts (same as in Snowflake).

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		args := []sdk.FunctionArgumentRequest{
			*sdk.NewFunctionArgumentRequest("A", nil).WithArgDataTypeOld("NUMBER(2, 0)"),
			*sdk.NewFunctionArgumentRequest("B", nil).WithArgDataTypeOld("DECIMAL"),
			*sdk.NewFunctionArgumentRequest("C", nil).WithArgDataTypeOld("INTEGER"),
			*sdk.NewFunctionArgumentRequest("D", nil).WithArgDataTypeOld(sdk.DataTypeFloat),
			*sdk.NewFunctionArgumentRequest("E", nil).WithArgDataTypeOld("DOUBLE"),
			*sdk.NewFunctionArgumentRequest("F", nil).WithArgDataTypeOld("VARCHAR(20)"),
			*sdk.NewFunctionArgumentRequest("G", nil).WithArgDataTypeOld("CHAR"),
			*sdk.NewFunctionArgumentRequest("H", nil).WithArgDataTypeOld(sdk.DataTypeString),
			*sdk.NewFunctionArgumentRequest("I", nil).WithArgDataTypeOld("TEXT"),
			*sdk.NewFunctionArgumentRequest("J", nil).WithArgDataTypeOld(sdk.DataTypeBinary),
			*sdk.NewFunctionArgumentRequest("K", nil).WithArgDataTypeOld("VARBINARY"),
			*sdk.NewFunctionArgumentRequest("L", nil).WithArgDataTypeOld(sdk.DataTypeBoolean),
			*sdk.NewFunctionArgumentRequest("M", nil).WithArgDataTypeOld(sdk.DataTypeDate),
			*sdk.NewFunctionArgumentRequest("N", nil).WithArgDataTypeOld("DATETIME"),
			*sdk.NewFunctionArgumentRequest("O", nil).WithArgDataTypeOld(sdk.DataTypeTime),
			*sdk.NewFunctionArgumentRequest("R", nil).WithArgDataTypeOld(sdk.DataTypeTimestampLTZ),
			*sdk.NewFunctionArgumentRequest("S", nil).WithArgDataTypeOld(sdk.DataTypeTimestampNTZ),
			*sdk.NewFunctionArgumentRequest("T", nil).WithArgDataTypeOld(sdk.DataTypeTimestampTZ),
			*sdk.NewFunctionArgumentRequest("U", nil).WithArgDataTypeOld(sdk.DataTypeVariant),
			*sdk.NewFunctionArgumentRequest("V", nil).WithArgDataTypeOld(sdk.DataTypeObject),
			*sdk.NewFunctionArgumentRequest("W", nil).WithArgDataTypeOld(sdk.DataTypeArray),
			*sdk.NewFunctionArgumentRequest("X", nil).WithArgDataTypeOld(sdk.DataTypeGeography),
			*sdk.NewFunctionArgumentRequest("Y", nil).WithArgDataTypeOld(sdk.DataTypeGeometry),
			*sdk.NewFunctionArgumentRequest("Z", nil).WithArgDataTypeOld("VECTOR(INT, 16)"),
		}
		err := client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
			id,
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(nil).WithResultDataTypeOld(sdk.DataTypeVariant)),
			"3.8",
			"add",
		).
			WithArguments(args).
			WithFunctionDefinition("def add(A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, R, S, T, U, V, W, X, Y, Z): A + A"),
		)
		require.NoError(t, err)

		dataTypes := make([]sdk.DataType, len(args))
		for i, arg := range args {
			dataType, err := datatypes.ParseDataType(string(arg.ArgDataTypeOld))
			require.NoError(t, err)
			dataTypes[i] = sdk.LegacyDataTypeFrom(dataType)
		}
		idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), dataTypes...)

		function, err := client.Functions.ShowByID(ctx, idWithArguments)
		require.NoError(t, err)
		require.Equal(t, dataTypes, function.ArgumentsOld)
	})

	// This test shows behavior of detailed types (e.g. VARCHAR(20) and NUMBER(10, 0)) on Snowflake side for functions.
	// For SHOW, data type is generalized both for argument and return type (to e.g. VARCHAR and NUMBER).
	// FOR DESCRIBE, data type is generalized for argument and works weirdly for the return type: type is generalized to the canonical one, but we also get the attributes.
	for _, tc := range []string{
		"NUMBER(36, 5)",
		"NUMBER(36)",
		"NUMBER",
		"DECIMAL",
		"INTEGER",
		"FLOAT",
		"DOUBLE",
		"VARCHAR",
		"VARCHAR(20)",
		"CHAR",
		"CHAR(10)",
		"TEXT",
		"BINARY",
		"BINARY(1000)",
		"VARBINARY",
		"BOOLEAN",
		"DATE",
		"DATETIME",
		"TIME",
		"TIMESTAMP_LTZ",
		"TIMESTAMP_NTZ",
		"TIMESTAMP_TZ",
		"VARIANT",
		"OBJECT",
		"ARRAY",
		"GEOGRAPHY",
		"GEOMETRY",
		"VECTOR(INT, 16)",
		"VECTOR(FLOAT, 8)",
	} {
		tc := tc
		t.Run(fmt.Sprintf("function returns non detailed data types of arguments for %s", tc), func(t *testing.T) {
			id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			argName := "A"
			dataType, err := datatypes.ParseDataType(tc)
			require.NoError(t, err)
			args := []sdk.FunctionArgumentRequest{
				*sdk.NewFunctionArgumentRequest(argName, dataType),
			}

			err = client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
				id,
				*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(dataType)),
				"3.8",
				"add",
			).
				WithArguments(args).
				WithFunctionDefinition(fmt.Sprintf("def add(%[1]s): %[1]s", argName)),
			)
			require.NoError(t, err)

			oldDataType := sdk.LegacyDataTypeFrom(dataType)
			idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), oldDataType)

			function, err := client.Functions.ShowByID(ctx, idWithArguments)
			require.NoError(t, err)
			assert.Equal(t, []sdk.DataType{oldDataType}, function.ArgumentsOld)
			assert.Equal(t, fmt.Sprintf("%[1]s(%[2]s) RETURN %[2]s", id.Name(), oldDataType), function.ArgumentsRaw)

			details, err := client.Functions.Describe(ctx, idWithArguments)
			require.NoError(t, err)
			pairs := make(map[string]string)
			for _, detail := range details {
				pairs[detail.Property] = detail.Value
			}
			assert.Equal(t, fmt.Sprintf("(%s %s)", argName, oldDataType), pairs["signature"])
			assert.Equal(t, dataType.Canonical(), pairs["returns"])
		})
	}
}
