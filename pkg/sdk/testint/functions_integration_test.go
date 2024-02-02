package testint

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
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

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, dts))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	t.Run("create function for Java", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		definition := `
		class TestFunc {
			public static String echoVarchar(String x) {
				return x;
			}
		}`
		target := fmt.Sprintf("@~/tf-%d.jar", time.Now().Unix())
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVARCHAR)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeVARCHAR).WithDefaultValue(sdk.String("'abc'"))
		request := sdk.NewCreateForJavaFunctionRequest(id, *returns, "TestFunc.echoVarchar").
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithNullInputBehavior(sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput)).
			WithTargetPath(&target).
			WithFunctionDefinition(&definition)
		err := client.Functions.CreateForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"VARCHAR"}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "JAVA", function.Language)
	})

	t.Run("create function for Javascript", func(t *testing.T) {
		name := "js_factorial"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

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
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		argument := sdk.NewFunctionArgumentRequest("d", sdk.DataTypeFloat)
		request := sdk.NewCreateForJavascriptFunctionRequest(id, *returns, definition).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithNullInputBehavior(sdk.NullInputBehaviorPointer(sdk.NullInputBehaviorCalledOnNullInput))
		err := client.Functions.CreateForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "JAVASCRIPT", function.Language)
	})

	t.Run("create function for Python", func(t *testing.T) {
		name := random.StringN(8)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		definition := `
def dump(i):
	print("Hello World!")`
		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeVariant)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		argument := sdk.NewFunctionArgumentRequest("i", sdk.DataTypeNumber)
		request := sdk.NewCreateForPythonFunctionRequest(id, *returns, "3.8", "dump").
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithFunctionDefinition(&definition)
		err := client.Functions.CreateForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"int"}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "PYTHON", function.Language)
	})

	t.Run("create function for Scala", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		definition := `
		class Echo {
			def echoVarchar(x : String): String = {
				return x
			}
		}`

		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeVARCHAR)
		request := sdk.NewCreateForScalaFunctionRequest(id, sdk.DataTypeVARCHAR, "Echo.echoVarchar").
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithRuntimeVersion(sdk.String("2.12")).
			WithFunctionDefinition(&definition)
		err := client.Functions.CreateForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SCALA", function.Language)
	})

	t.Run("create function for SQL", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeFloat)
		request := sdk.NewCreateForSQLFunctionRequest(id, *returns, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argument}).
			WithOrReplace(sdk.Bool(true)).
			WithComment(sdk.String("comment"))
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})

	t.Run("create function for SQL with no arguments", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		request := sdk.NewCreateForSQLFunctionRequest(id, *returns, definition).
			WithOrReplace(sdk.Bool(true)).
			WithComment(sdk.String("comment"))
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, nil))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})
}

func TestInt_OtherFunctions(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)
	tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	assertFunction := func(t *testing.T, id sdk.SchemaObjectIdentifier, secure bool, withArguments bool) {
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

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifier, dts []sdk.DataType) func() {
		return func() {
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id, dts))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, cleanup bool, withArguments bool) *sdk.Function {
		t.Helper()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))

		definition := "3.141592654::FLOAT"

		dt := sdk.NewFunctionReturnsResultDataTypeRequest(sdk.DataTypeFloat)
		returns := sdk.NewFunctionReturnsRequest().WithResultDataType(dt)
		request := sdk.NewCreateForSQLFunctionRequest(id, *returns, definition).
			WithOrReplace(sdk.Bool(true))
		if withArguments {
			argument := sdk.NewFunctionArgumentRequest("x", sdk.DataTypeFloat)
			request = request.WithArguments([]sdk.FunctionArgumentRequest{*argument})
		}
		err := client.Functions.CreateForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			if withArguments {
				t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))
			} else {
				t.Cleanup(cleanupFunctionHandle(id, nil))
			}
		}
		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		return function
	}

	defaultAlterRequest := func(id sdk.SchemaObjectIdentifier) *sdk.AlterFunctionRequest {
		return sdk.NewAlterFunctionRequest(id, []sdk.DataType{sdk.DataTypeFloat})
	}

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(3))
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))
		} else {
			t.Cleanup(cleanupFunctionHandle(nid, []sdk.DataType{sdk.DataTypeFloat}))
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

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetLogLevel(sdk.String("DEBUG")))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetLogLevel(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetTraceLevel(sdk.String("ALWAYS")))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetTraceLevel(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetComment(sdk.String("test comment")))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetComment(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithSetSecure(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, true, true)
	})

	t.Run("alter function: set secure with no arguments", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, false)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, sdk.NewAlterFunctionRequest(id, nil).WithSetSecure(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, true, false)
	})

	t.Run("alter function: unset secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		err := client.Functions.Alter(ctx, defaultAlterRequest(id).WithUnsetSecure(sdk.Bool(true)))
		require.NoError(t, err)
		assertFunction(t, id, false, true)
	})

	t.Run("alter function: set and unset tags", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
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

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(&sdk.Like{Pattern: &f1.Name}))
		require.NoError(t, err)

		require.Equal(t, 1, len(functions))
		require.Contains(t, functions, *f1)
		require.NotContains(t, functions, *f2)
	})

	t.Run("show function for SQL: no matches", func(t *testing.T) {
		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(&sdk.Like{Pattern: sdk.String(random.String())}))
		require.NoError(t, err)
		require.Equal(t, 0, len(functions))
	})

	t.Run("describe function for SQL", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true, true)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)

		request := sdk.NewDescribeFunctionRequest(id, []sdk.DataType{sdk.DataTypeFloat})
		details, err := client.Functions.Describe(ctx, request)
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
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)

		request := sdk.NewDescribeFunctionRequest(id, nil)
		details, err := client.Functions.Describe(ctx, request)
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
