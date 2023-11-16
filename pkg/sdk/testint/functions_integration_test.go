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

func TestInt_Functions(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	_, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	defaultArgumentTypes := []sdk.FunctionArgumentTypeRequest{
		*sdk.NewFunctionArgumentTypeRequest().WithArgDataType("FLOAT"),
	}

	cleanupFunctionHandle := func(id sdk.SchemaObjectIdentifier, argumentTypes []string) func() {
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

		as, at := "3.141592654::FLOAT", "FLOAT"
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String(at))
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		request := sdk.NewCreateFunctionRequest(id).
			WithArguments([]sdk.FunctionArgumentRequest{
				*sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType(at),
			}).
			WithOrReplace(sdk.Bool(true)).
			WithReturns(*returnsRequest).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupFunctionHandle(id, []string{at}))
		}

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		return function
	}

	t.Run("create function for JAVA", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
class TestFunc {
	public static String echoVarchar(String x) {
		return x;
	}
}`
		target := fmt.Sprintf("@~/tf-%d.jar", time.Now().Unix())
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String("VARCHAR"))
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType("VARCHAR")
		strictOrNotRequest := sdk.NewFunctionStrictOrNotRequest().WithCalledOnNullInput(sdk.Bool(true))
		request := sdk.NewCreateFunctionRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithLanguage(sdk.String("JAVA")).
			WithStrictOrNot(strictOrNotRequest).
			WithReturns(*returnsRequest).
			WithHandler(sdk.String("TestFunc.echoVarchar")).
			WithTargetPath(&target).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []string{"VARCHAR"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for PYTHON", func(t *testing.T) {
		name := random.StringN(8)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	def dump(i):
		print("Hello World!")
		`
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String("VARIANT"))
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("i").WithArgDataType("int")
		request := sdk.NewCreateFunctionRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithLanguage(sdk.String("PYTHON")).
			WithReturns(*returnsRequest).
			WithRuntimeVersion(sdk.String("3.8")).
			WithHandler(sdk.String("dump")).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []string{"int"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for SCALA", func(t *testing.T) {
		name := "echo_varchar"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
		class Echo {
			def echoVarchar(x : String): String = {
				return x
			}
		}
		`
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String("VARCHAR"))
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType("VARCHAR")
		request := sdk.NewCreateFunctionRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithLanguage(sdk.String("SCALA")).
			WithReturns(*returnsRequest).
			WithRuntimeVersion(sdk.String("2.12")).
			WithHandler(sdk.String("Echo.echoVarchar")).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []string{"VARCHAR"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for JAVASCRIPT", func(t *testing.T) {
		name := "js_factorial"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
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
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String("double"))
		argumentRequest := sdk.NewFunctionArgumentRequest().WithArgName("d").WithArgDataType("double")
		strictOrNotRequest := sdk.NewFunctionStrictOrNotRequest().WithStrict(sdk.Bool(true))
		request := sdk.NewCreateFunctionRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.FunctionArgumentRequest{*argumentRequest}).
			WithLanguage(sdk.String("JAVASCRIPT")).
			WithStrictOrNot(strictOrNotRequest).
			WithReturns(*returnsRequest).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []string{"double"}))

		functions, err := client.Functions.Show(ctx, sdk.NewShowFunctionRequest())
		require.NoError(t, err)
		require.Equal(t, 1, len(functions))
	})

	t.Run("create function for SQL", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := "3.141592654::FLOAT"
		returnsRequest := sdk.NewFunctionReturnsRequest().WithResultDataType(sdk.String("FLOAT"))
		request := sdk.NewCreateFunctionRequest(id).
			WithArguments([]sdk.FunctionArgumentRequest{
				*sdk.NewFunctionArgumentRequest().WithArgName("x").WithArgDataType("FLOAT"),
			}).
			WithOrReplace(sdk.Bool(true)).
			WithComment(sdk.String("comment")).
			WithReturns(*returnsRequest).
			WithAs(sdk.String(as))
		err := client.Functions.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupFunctionHandle(id, []string{"FLOAT"}))

		function, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id.Name(), function.Name)
		require.Equal(t, "SQL", function.Language)
	})

	t.Run("alter function: rename", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, false)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		request := sdk.NewAlterFunctionRequest(id).WithRenameTo(&nid).WithArgumentTypes(defaultArgumentTypes)
		err := client.Functions.Alter(ctx, request)
		if err != nil {
			t.Cleanup(cleanupFunctionHandle(id, []string{"FLOAT"}))
		} else {
			t.Cleanup(cleanupFunctionHandle(nid, []string{"FLOAT"}))
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

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithLogLevel(sdk.String("DEBUG"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset log level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithLogLevel(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithTraceLevel(sdk.String("ALWAYS"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset trace level", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithTraceLevel(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewFunctionSetRequest().WithComment(sdk.String("comment"))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: unset comment", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithComment(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter function: set secure", func(t *testing.T) {
		f := createFunctionForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
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

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		unset := sdk.NewFunctionUnsetRequest().WithSecure(sdk.Bool(true))
		request := sdk.NewAlterFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Functions.Alter(ctx, request)
		require.NoError(t, err)

		e, err := client.Functions.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, false, e.IsSecure)
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
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)

		request := sdk.NewDescribeFunctionRequest(id).WithArgumentTypes(defaultArgumentTypes)
		details, err := client.Functions.Describe(ctx, request)
		require.NoError(t, err)
		require.Greater(t, len(details), 0)
	})
}
