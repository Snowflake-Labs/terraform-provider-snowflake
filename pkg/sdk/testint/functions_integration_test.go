package testint

import (
	"context"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

		definition := `
		class TestFunc {
			public static String echoVarchar(String x) {
				return x;
			}
		}`
		target := fmt.Sprintf("@~/tf-%d.jar", time.Now().Unix())
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVARCHAR)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeVARCHAR).WithDefaultValue("'abc'")
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

		definition := `
		if (D <= 0) {
			return 1;
		} else {
			var result = 1;
			for (var i = 2; i <= D; i++) {
				result = result * i;
			}
			return result;
		}`

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("d", sdk.DataTypeFloat)
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

		definition := `
def dump(i):
	print("Hello World!")`
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVariant)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("i", sdk.DataTypeNumber)
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

		definition := `
		class Echo {
			def echoVarchar(x : String): String = {
				return x
			}
		}`

		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForScalaFunctionRequest(id.SchemaObjectId(), sdk.DataTypeVARCHAR, "Echo.echoVarchar").
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

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeFloat)
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

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
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

	tagTest, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

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
		assert.NotEmpty(t, function.Arguments)
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

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).
			WithOrReplace(true)
		if withArguments {
			argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeFloat)
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

	defaultAlterRequest := func(id sdk.SchemaObjectIdentifierWithArguments) *sdk.AlterFunctionRequest {
		return sdk.NewAlterFunctionRequest(id)
	}

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false, true)

		id := f.ID()
		nid := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArguments()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithRenameTo(nid.SchemaObjectId()))
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
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetLogLevel(string(sdk.LogLevelDebug)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetLogLevel(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetTraceLevel(string(sdk.TraceLevelAlways)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetTraceLevel(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetComment("test comment"))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetComment(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetSecure(true))
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
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetSecure(true))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set and unset tags", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := f.ID()
		setTags := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "v1",
			},
		}
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetTags(setTags))
		require.NoError(t, err)
		assertFunction(t, id, false, true)

		unsetTags := []sdk.ObjectIdentifier{
			tagTest.ID(),
		}
		err = client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetTags(unsetTags))
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

		definition := "3.141592654::FLOAT"
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(*dt)
		request := sdk.NewCreateForSQLFunctionRequest(id.SchemaObjectId(), *returns, definition).WithOrReplace(true)

		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeFloat)
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

	t.Run("function returns non detailed data types of arguments", func(t *testing.T) {
		// This test proves that every detailed data type (e.g. VARCHAR(20) and NUMBER(10, 0)) is generalized
		// (to e.g. VARCHAR and NUMBER) and that sdk.ToDataType mapping function maps detailed types correctly to
		// their generalized counterparts.

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		args := []sdk.FunctionArgumentRequest{
			*sdk.NewFunctionArgumentRequest("A", "NUMBER(2, 0)"),
			*sdk.NewFunctionArgumentRequest("B", "DECIMAL"),
			*sdk.NewFunctionArgumentRequest("C", "INTEGER"),
			*sdk.NewFunctionArgumentRequest("D", sdk.DataTypeFloat),
			*sdk.NewFunctionArgumentRequest("E", "DOUBLE"),
			*sdk.NewFunctionArgumentRequest("F", "VARCHAR(20)"),
			*sdk.NewFunctionArgumentRequest("G", "CHAR"),
			*sdk.NewFunctionArgumentRequest("H", sdk.DataTypeString),
			*sdk.NewFunctionArgumentRequest("I", "TEXT"),
			*sdk.NewFunctionArgumentRequest("J", sdk.DataTypeBinary),
			*sdk.NewFunctionArgumentRequest("K", "VARBINARY"),
			*sdk.NewFunctionArgumentRequest("L", sdk.DataTypeBoolean),
			*sdk.NewFunctionArgumentRequest("M", sdk.DataTypeDate),
			*sdk.NewFunctionArgumentRequest("N", "DATETIME"),
			*sdk.NewFunctionArgumentRequest("O", sdk.DataTypeTime),
			*sdk.NewFunctionArgumentRequest("P", sdk.DataTypeTimestamp),
			*sdk.NewFunctionArgumentRequest("R", sdk.DataTypeTimestampLTZ),
			*sdk.NewFunctionArgumentRequest("S", sdk.DataTypeTimestampNTZ),
			*sdk.NewFunctionArgumentRequest("T", sdk.DataTypeTimestampTZ),
			*sdk.NewFunctionArgumentRequest("U", sdk.DataTypeVariant),
			*sdk.NewFunctionArgumentRequest("V", sdk.DataTypeObject),
			*sdk.NewFunctionArgumentRequest("W", sdk.DataTypeArray),
			*sdk.NewFunctionArgumentRequest("X", sdk.DataTypeGeography),
			*sdk.NewFunctionArgumentRequest("Y", sdk.DataTypeGeometry),
			*sdk.NewFunctionArgumentRequest("Z", "VECTOR(INT, 16)"),
		}
		err := client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
			id,
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVariant)),
			"3.8",
			"add",
		).
			WithArguments(args).
			WithFunctionDefinition("def add(A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, R, S, T, U, V, W, X, Y, Z): A + A"),
		)
		require.NoError(t, err)

		dataTypes := make([]sdk.DataType, len(args))
		for i, arg := range args {
			dataTypes[i], err = sdk.ToDataType(string(arg.ArgDataType))
			require.NoError(t, err)
		}
		idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), dataTypes...)

		_, err = client.Functions.ShowByID(ctx, idWithArguments)
		require.NoError(t, err)
	})

	t.Run("function returns non detailed data types of arguments", func(t *testing.T) {
		// This test proves that every detailed data type (e.g. VARCHAR(20) and NUMBER(10, 0)) is generalized
		// (to e.g. VARCHAR and NUMBER) and that sdk.ToDataType mapping function maps detailed types correctly to
		// their generalized counterparts.

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		args := []sdk.FunctionArgumentRequest{
			*sdk.NewFunctionArgumentRequest("A", "NUMBER(2, 0)"),
			*sdk.NewFunctionArgumentRequest("B", "DECIMAL"),
			*sdk.NewFunctionArgumentRequest("C", "INTEGER"),
			*sdk.NewFunctionArgumentRequest("D", sdk.DataTypeFloat),
			*sdk.NewFunctionArgumentRequest("E", "DOUBLE"),
			*sdk.NewFunctionArgumentRequest("F", "VARCHAR(20)"),
			*sdk.NewFunctionArgumentRequest("G", "CHAR"),
			*sdk.NewFunctionArgumentRequest("H", sdk.DataTypeString),
			*sdk.NewFunctionArgumentRequest("I", "TEXT"),
			*sdk.NewFunctionArgumentRequest("J", sdk.DataTypeBinary),
			*sdk.NewFunctionArgumentRequest("K", "VARBINARY"),
			*sdk.NewFunctionArgumentRequest("L", sdk.DataTypeBoolean),
			*sdk.NewFunctionArgumentRequest("M", sdk.DataTypeDate),
			*sdk.NewFunctionArgumentRequest("N", "DATETIME"),
			*sdk.NewFunctionArgumentRequest("O", sdk.DataTypeTime),
			*sdk.NewFunctionArgumentRequest("P", sdk.DataTypeTimestamp),
			*sdk.NewFunctionArgumentRequest("R", sdk.DataTypeTimestampLTZ),
			*sdk.NewFunctionArgumentRequest("S", sdk.DataTypeTimestampNTZ),
			*sdk.NewFunctionArgumentRequest("T", sdk.DataTypeTimestampTZ),
			*sdk.NewFunctionArgumentRequest("U", sdk.DataTypeVariant),
			*sdk.NewFunctionArgumentRequest("V", sdk.DataTypeObject),
			*sdk.NewFunctionArgumentRequest("W", sdk.DataTypeArray),
			*sdk.NewFunctionArgumentRequest("X", sdk.DataTypeGeography),
			*sdk.NewFunctionArgumentRequest("Y", sdk.DataTypeGeometry),
			*sdk.NewFunctionArgumentRequest("Z", "VECTOR(INT, 16)"),
		}
		err := client.Functions.CreateForPython(ctx, sdk.NewCreateForPythonFunctionRequest(
			id,
			*sdk.NewFunctionReturnsRequest().WithResultDataType(*sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVariant)),
			"3.8",
			"add",
		).
			WithArguments(args).
			WithFunctionDefinition("def add(A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, R, S, T, U, V, W, X, Y, Z): A + A"),
		)
		require.NoError(t, err)

		dataTypes := make([]sdk.DataType, len(args))
		for i, arg := range args {
			dataTypes[i], err = sdk.ToDataType(string(arg.ArgDataType))
			require.NoError(t, err)
		}
		idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(id.DatabaseName(), id.SchemaName(), id.Name(), dataTypes...)

		function, err := client.Functions.ShowByID(ctx, idWithArguments)
		require.NoError(t, err)
		require.Equal(t, dataTypes, function.Arguments)
	})
}
