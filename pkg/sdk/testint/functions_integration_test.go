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

func TestInt_CreateFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifier, argumentTypes []sdk.DataType) func() {
		return func() {
			es := []sdk.FunctionArgumentTypeRequest{}
			for _, item := range argumentTypes {
				es = append(es, *sdk.NewFunctionArgumentTypeRequest().WithArgDataType(item))
			}
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id).WithArgumentTypes(es))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	t.Run("create function for Java", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		definition := `
	class TestFunc {
		public static String echoVarchar(String x) {
			return x;
		}
	}`
		target := fmt.Sprintf("@~/tf-%d.jar", time.Now().Unix())
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType(sdk.DataTypeVARCHAR).WithDefault(sdk.String("abc"))
		request := sdk.NewCreateFunctionForJavaFunctionRequest(id, returnsRequest, "TestFunc.echoVarchar", definition).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithNullInputBehavior(sdk.FunctionNullInputBehaviorCalledOnNullInput).
			WithTargetPath(&target)
		err := client.Functions.CreateFunctionForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"VARCHAR"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for Javascript", func(t *testing.T) {
		name := "js_factorial"
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		definition := `
		if (D <= 0) {
			return 1;
		} else {
			var result = 1;
			for (var i = 2; i <= D; i++) {
				result = result * i;
			}
			return result;
		}
		`
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeFloat)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("d").WithArgDataType(sdk.DataTypeFloat)
		request := sdk.NewCreateFunctionForJavascriptFunctionRequest(id, returnsRequest, definition).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithNullInputBehavior(sdk.FunctionNullInputBehaviorCalledOnNullInput)
		err := client.Functions.CreateFunctionForJavascript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for Python", func(t *testing.T) {
		name := random.StringN(8)
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		definition := `
def dump(i):
	print("Hello World!")
		`
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeVariant)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("i").WithArgDataType(sdk.DataTypeNumber)
		request := sdk.NewCreateFunctionForPythonFunctionRequest(id, returnsRequest, "3.8", "dump", definition).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest})
		err := client.Functions.CreateFunctionForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"int"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for Scala", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		definition := `
			class Echo {
				def echoVarchar(x : String): String = {
					return x
				}
			}
			`
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateFunctionForScalaFunctionRequest(id, returnsRequest, "Echo.echoVarchar", definition).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithRuntimeVersion(sdk.String("2.12"))
		err := client.Functions.CreateFunctionForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"VARCHAR"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for SQL", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		definition := "3.141592654::FLOAT"
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeFloat)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType(sdk.DataTypeFloat)
		request := sdk.NewCreateFunctionForSQLFunctionRequest(id, returnsRequest, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithOrReplace(sdk.Bool(true)).
			WithComment(sdk.String("comment"))
		err := client.Functions.CreateFunctionForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"FLOAT"}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})
}

func TestInt_AlterAndShowFunctions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	tagTest, tagCleanup := createTag(t, client, testDb(t), testSchema(t))
	t.Cleanup(tagCleanup)

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifier, argumentTypes []sdk.DataType) func() {
		return func() {
			es := []sdk.FunctionArgumentTypeRequest{}
			for _, item := range argumentTypes {
				es = append(es, *sdk.NewFunctionArgumentTypeRequest().WithArgDataType(item))
			}
			err := client.Functions.Drop(ctx, sdk.NewDropFunctionRequest(id).WithArgumentTypes(es))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createFunctionForSQLHandle := func(t *testing.T, cleanup bool) *sdk.Function {
		t.Helper()

		definition := "3.141592654::FLOAT"
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.DataTypeFloat)
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType(sdk.DataTypeFloat)
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.String())
		request := sdk.NewCreateFunctionForSQLFunctionRequest(id, returnsRequest, definition).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithOrReplace(sdk.Bool(true))
		err := client.Functions.CreateFunctionForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{sdk.DataTypeFloat}))
		}
		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		return function
	}

	defaultArgumentTypes := []sdk.FunctionArgumentTypeRequest{
		*sdk.NewFunctionArgumentTypeRequest().WithArgDataType(sdk.DataTypeFloat),
	}

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		nid := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, random.String())
		request := sdk.NewAlterFunctionRequest(id).WithRenameTo(&nid).WithArgumentTypes(defaultArgumentTypes)
		err := client.Functions.Alter(ctx, request)
		if err != nil {
			t.Cleanup(cleanupFunctionHandle(id, []sdk.DataType{"FLOAT"}))
		} else {
			t.Cleanup(cleanupFunctionHandle(nid, []sdk.DataType{"FLOAT"}))
		}
		require.NoError(t, err)

		_, err = client.Functions.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Functions.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter function: set log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithLogLevel(sdk.String("DEBUG"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithLogLevel(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithTraceLevel(sdk.String("ALWAYS"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithTraceLevel(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithComment(sdk.String("comment"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithComment(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithSecure(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		e, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, true, e.IsSecure)
	})

	t.Run("alter function: unset secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithSecure(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		e, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, false, e.IsSecure)
	})

	t.Run("alter function: set and unset tags", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)
		setTags := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "abc",
			},
		}
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSetTags(setTags)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		unsetTags := []sdk.ObjectIdentifier{
			tagTest.ID(),
		}
		request = sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnsetTags(unsetTags)
		err = client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("show function for SQL: without like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true)
		f2 := createFunctionForSQLHandle(t, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)

		require.Equal(t, 2, len(functions))
		require.Contains(t, functions, *f1)
		require.Contains(t, functions, *f2)
	})

	t.Run("show function for SQL: with like", func(t *testing.T) {
		f1 := createFunctionForSQLHandle(t, true)
		f2 := createFunctionForSQLHandle(t, true)

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(f1.Name))
		require.NoError(t, err)

		require.Equal(t, 1, len(functions))
		require.Contains(t, functions, *f1)
		require.NotContains(t, functions, *f2)
	})

	t.Run("show function for SQL: no matches", func(t *testing.T) {
		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest().WithLike(random.String()))
		require.NoError(t, err)
		require.Equal(t, 0, len(functions))
	})

	t.Run("describe function for SQL", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, f.Name)

		request := sdk.NewDescribeFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes)
		details, err := client.Functions.Describe(ctx, request)
		require.NoError(t, err)
		require.Greater(t, len(details), 0)
	})
}
