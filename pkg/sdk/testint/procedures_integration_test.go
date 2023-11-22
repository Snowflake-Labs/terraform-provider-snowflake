package testint

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CreateProcedures(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	_, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)

	cleanupProcedureHandle := func(id sdk.SchemaObjectIdentifier, argumentTypes []sdk.DataType) func() {
		return func() {
			es := []sdk.ProcedureArgumentTypeRequest{}
			for _, item := range argumentTypes {
				es = append(es, *sdk.NewProcedureArgumentTypeRequest().WithArgDataType(item))
			}
			err := client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithArgumentTypes(es))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	t.Run("create procedure for Java: returns result data type", func(t *testing.T) {
		name := "file_reader_java_proc_snowflakefile"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	import java.io.InputStream;
	import java.io.IOException;
	import java.nio.charset.StandardCharsets;
	import com.snowflake.snowpark_java.types.SnowflakeFile;
	import com.snowflake.snowpark_java.Session;
	class FileReader {
		public String execute(Session session, String fileName) throws IOException {
			InputStream input = SnowflakeFile.newInstance(fileName).getInputStream();
			return new String(input.readAllBytes(), StandardCharsets.UTF_8);
		}
	}`
		resultDataType := sdk.NewProcedureReturnsResultDataTypeRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		procedureReturns := sdk.NewProcedureReturnsRequest().WithResultDataType(resultDataType)
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("input").WithArgDataType(sdk.DataTypeVARCHAR)
		procedurePackage := sdk.NewProcedurePackageRequest().WithPackage("com.snowflake:snowpark:latest")
		request := sdk.NewCreateProcedureForJavaProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("11")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithHandler("FileReader.execute").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Java: returns table", func(t *testing.T) {
		name := "filter_by_role"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	import com.snowflake.snowpark_java.*;
	public class Filter {
		public DataFrame filterByRole(Session session, String tableName, String role) {
			DataFrame table = session.table(tableName);
			DataFrame filteredRows = table.filter(Functions.col("role").equal_to(Functions.lit(role)));
			return filteredRows;
		}
	}`
		column1 := sdk.NewProcedureColumnRequest().WithColumnName("id").WithColumnDataType(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest().WithColumnName("name").WithColumnDataType(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest().WithColumnName("role").WithColumnDataType(sdk.DataTypeVARCHAR)
		procedureReturnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		procedureReturns := sdk.NewProcedureReturnsRequest().WithTable(procedureReturnsTable)
		procedureArgument1 := sdk.NewProcedureArgumentRequest().WithArgName("table_name").WithArgDataType(sdk.DataTypeVARCHAR)
		procedureArgument2 := sdk.NewProcedureArgumentRequest().WithArgName("role").WithArgDataType(sdk.DataTypeVARCHAR)
		procedurePackage := sdk.NewProcedurePackageRequest().WithPackage("com.snowflake:snowpark:latest")
		request := sdk.NewCreateProcedureForJavaProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("11")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument1, *procedureArgument2}).
			WithReturns(procedureReturns).
			WithHandler("Filter.filterByRole").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForJava(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Javascript", func(t *testing.T) {
		name := "stproc1"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	var sql_command = "INSERT INTO stproc_test_table1 (num_col1) VALUES (" + FLOAT_PARAM1 + ")";
	try {
		snowflake.execute (
			{sqlText: sql_command}
		);
		return "Succeeded."; // Return a success/error indicator.
	}
	catch (err)  {
		return "Failed: " + err; // Return a success/error indicator.
	}`
		procedureReturns := sdk.NewProcedureReturns2Request().WithResultDataType("STRING")
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("FLOAT_PARAM1").WithArgDataType("FLOAT")
		request := sdk.NewCreateProcedureForJavaScriptProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithNullInputBehavior(sdk.ProcedureNullInputBehaviorStrict).
			WithExecuteAs(sdk.ProcedureExecuteAsCaller).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForJavaScript(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{"FLOAT"}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Scala: returns result data type", func(t *testing.T) {
		name := "file_reader_scala_proc_snowflakefile"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	import java.io.InputStream
	import java.nio.charset.StandardCharsets
	import com.snowflake.snowpark_java.types.SnowflakeFile
	import com.snowflake.snowpark_java.Session
	object FileReader {
		def execute(session: Session, fileName: String): String = {
			var input: InputStream = SnowflakeFile.newInstance(fileName).getInputStream()
			return new String(input.readAllBytes(), StandardCharsets.UTF_8)
		}
	}`
		resultDataType := sdk.NewProcedureReturnsResultDataTypeRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		procedureReturns := sdk.NewProcedureReturnsRequest().WithResultDataType(resultDataType)
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("input").WithArgDataType(sdk.DataTypeVARCHAR)
		procedurePackage := sdk.NewProcedurePackageRequest().WithPackage("com.snowflake:snowpark:latest")
		request := sdk.NewCreateProcedureForScalaProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("2.12")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithHandler("FileReader.execute").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Scala: returns table", func(t *testing.T) {
		name := "filter_by_role"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	import com.snowflake.snowpark.functions._
	import com.snowflake.snowpark._
	object Filter {
		def filterByRole(session: Session, tableName: String, role: String): DataFrame = {
			val table = session.table(tableName)
			val filteredRows = table.filter(col("role") === role)
			return filteredRows
		}
	}`
		column1 := sdk.NewProcedureColumnRequest().WithColumnName("id").WithColumnDataType(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest().WithColumnName("name").WithColumnDataType(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest().WithColumnName("role").WithColumnDataType(sdk.DataTypeVARCHAR)
		procedureReturnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		procedureReturns := sdk.NewProcedureReturnsRequest().WithTable(procedureReturnsTable)
		procedureArgument1 := sdk.NewProcedureArgumentRequest().WithArgName("table_name").WithArgDataType(sdk.DataTypeVARCHAR)
		procedureArgument2 := sdk.NewProcedureArgumentRequest().WithArgName("role").WithArgDataType(sdk.DataTypeVARCHAR)
		procedurePackage := sdk.NewProcedurePackageRequest().WithPackage("com.snowflake:snowpark:latest")
		request := sdk.NewCreateProcedureForScalaProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("2.12")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument1, *procedureArgument2}).
			WithReturns(procedureReturns).
			WithHandler("Filter.filterByRole").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForScala(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Python: returns result data type", func(t *testing.T) {
		name := "joblib_multiprocessing_proc"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
import joblib
from math import sqrt
def joblib_multiprocessing(session, i):
	result = joblib.Parallel(n_jobs=-1)(joblib.delayed(sqrt)(i ** 2) for i in range(10))
	return str(result)`

		resultDataType := sdk.NewProcedureReturnsResultDataTypeRequest().WithResultDataType("STRING")
		procedureReturns := sdk.NewProcedureReturnsRequest().WithResultDataType(resultDataType)
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("i").WithArgDataType("INT")
		procedurePackage1 := sdk.NewProcedurePackageRequest().WithPackage("snowflake-snowpark-python")
		procedurePackage2 := sdk.NewProcedurePackageRequest().WithPackage("joblib")
		request := sdk.NewCreateProcedureForPythonProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("3.8")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithHandler("joblib_multiprocessing").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage1, *procedurePackage2}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{"string"}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for Python: returns table", func(t *testing.T) {
		name := "filterByRole"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
from snowflake.snowpark.functions import col
def filter_by_role(session, table_name, role):
	df = session.table(table_name)
	return df.filter(col("role") == role)`
		column1 := sdk.NewProcedureColumnRequest().WithColumnName("id").WithColumnDataType(sdk.DataTypeNumber)
		column2 := sdk.NewProcedureColumnRequest().WithColumnName("name").WithColumnDataType(sdk.DataTypeVARCHAR)
		column3 := sdk.NewProcedureColumnRequest().WithColumnName("role").WithColumnDataType(sdk.DataTypeVARCHAR)
		procedureReturnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2, *column3})
		procedureReturns := sdk.NewProcedureReturnsRequest().WithTable(procedureReturnsTable)
		procedureArgument1 := sdk.NewProcedureArgumentRequest().WithArgName("table_name").WithArgDataType(sdk.DataTypeVARCHAR)
		procedureArgument2 := sdk.NewProcedureArgumentRequest().WithArgName("role").WithArgDataType(sdk.DataTypeVARCHAR)
		procedurePackage := sdk.NewProcedurePackageRequest().WithPackage("snowflake-snowpark-python")
		request := sdk.NewCreateProcedureForPythonProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithRuntimeVersion(sdk.String("3.8")).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument1, *procedureArgument2}).
			WithReturns(procedureReturns).
			WithHandler("filter_by_role").
			WithPackages([]sdk.ProcedurePackageRequest{*procedurePackage}).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForPython(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR, sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for SQL: returns result data type", func(t *testing.T) {
		name := "output_message"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	BEGIN
		RETURN message;
	END;`

		resultDataType := sdk.NewProcedureReturnsResultDataTypeRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		procedureReturns := sdk.NewProcedureReturns3Request().WithResultDataType(resultDataType).WithNotNull(sdk.Bool(true))
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("message").WithArgDataType(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateProcedureForSQLProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})

	t.Run("create procedure for SQL: returns table", func(t *testing.T) {
		name := "find_invoice_by_id"
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)

		as := `
	DECLARE
		res RESULTSET DEFAULT (SELECT * FROM invoices WHERE id = :id);
	BEGIN
		RETURN TABLE(res);
	END;`
		column1 := sdk.NewProcedureColumnRequest().WithColumnName("id").WithColumnDataType("INTEGER")
		column2 := sdk.NewProcedureColumnRequest().WithColumnName("price").WithColumnDataType("NUMBER(12,2)")
		procedureReturnsTable := sdk.NewProcedureReturnsTableRequest().WithColumns([]sdk.ProcedureColumnRequest{*column1, *column2})
		procedureReturns := sdk.NewProcedureReturns3Request().WithTable(procedureReturnsTable)
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("id").WithArgDataType(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateProcedureForSQLProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForSQL(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(procedures), 1)
	})
}

func TestInt_AlterAndShowProcedures(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	_, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)

	cleanupProcedureHandle := func(id sdk.SchemaObjectIdentifier, argumentTypes []sdk.DataType) func() {
		return func() {
			es := []sdk.ProcedureArgumentTypeRequest{}
			for _, item := range argumentTypes {
				es = append(es, *sdk.NewProcedureArgumentTypeRequest().WithArgDataType(item))
			}
			err := client.Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithArgumentTypes(es))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createProcedureForSQLHandle := func(t *testing.T, cleanup bool) *sdk.Procedure {
		t.Helper()

		as := `
	BEGIN
		RETURN message;
	END;`
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		resultDataType := sdk.NewProcedureReturnsResultDataTypeRequest().WithResultDataType(sdk.DataTypeVARCHAR)
		procedureReturns := sdk.NewProcedureReturns3Request().WithResultDataType(resultDataType).WithNotNull(sdk.Bool(true))
		procedureArgument := sdk.NewProcedureArgumentRequest().WithArgName("message").WithArgDataType(sdk.DataTypeVARCHAR)
		request := sdk.NewCreateProcedureForSQLProcedureRequest(id).
			WithOrReplace(sdk.Bool(true)).
			WithArguments([]sdk.ProcedureArgumentRequest{*procedureArgument}).
			WithReturns(procedureReturns).
			WithExecuteAs(sdk.ProcedureExecuteAsCaller).
			WithProcedureDefinition(sdk.String(as))
		err := client.Procedures.CreateProcedureForSQL(ctx, request)
		require.NoError(t, err)
		if cleanup {
			t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))
		}
		procedure, err := client.Procedures.ShowByID(ctx, id)
		require.NoError(t, err)
		return procedure
	}

	defaultArgumentTypes := []sdk.ProcedureArgumentTypeRequest{
		*sdk.NewProcedureArgumentTypeRequest().WithArgDataType(sdk.DataTypeVARCHAR),
	}

	t.Run("alter procedure: rename", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, false)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.String())
		request := sdk.NewAlterProcedureRequest(id).WithRenameTo(&nid).WithArgumentTypes(defaultArgumentTypes)
		err := client.Procedures.Alter(ctx, request)
		if err != nil {
			t.Cleanup(cleanupProcedureHandle(id, []sdk.DataType{sdk.DataTypeVARCHAR}))
		} else {
			t.Cleanup(cleanupProcedureHandle(nid, []sdk.DataType{sdk.DataTypeVARCHAR}))
		}
		require.NoError(t, err)

		_, err = client.Procedures.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		e, err := client.Procedures.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), e.Name)
	})

	t.Run("alter procedure: set log level", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewProcedureSetRequest().WithLogLevel(sdk.String("DEBUG"))
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter procedure: set trace level", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewProcedureSetRequest().WithTraceLevel(sdk.String("ALWAYS"))
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter procedure: set comment", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		set := sdk.NewProcedureSetRequest().WithComment(sdk.String("comment"))
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSet(set)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter procedure: unset comment", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		unset := sdk.NewProcedureUnsetRequest().WithComment(sdk.Bool(true))
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnset(unset)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter procedure: set execute as", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithExecuteAs(sdk.ProcedureExecuteAsOwner)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("alter procedure: set and unset tags", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)
		setTags := []sdk.TagAssociation{
			{
				Name:  tagTest.ID(),
				Value: "abc",
			},
		}
		request := sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithSetTags(setTags)
		err := client.Procedures.Alter(ctx, request)
		require.NoError(t, err)

		unsetTags := []sdk.ObjectIdentifier{
			tagTest.ID(),
		}
		request = sdk.NewAlterProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes).WithUnsetTags(unsetTags)
		err = client.Procedures.Alter(ctx, request)
		require.NoError(t, err)
	})

	t.Run("show procedure for SQL: without like", func(t *testing.T) {
		f1 := createProcedureForSQLHandle(t, true)
		f2 := createProcedureForSQLHandle(t, true)

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest())
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(procedures), 1)
		require.Contains(t, procedures, *f1)
		require.Contains(t, procedures, *f2)
	})

	t.Run("show procedure for SQL: with like", func(t *testing.T) {
		f1 := createProcedureForSQLHandle(t, true)
		f2 := createProcedureForSQLHandle(t, true)

		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest().WithLike(f1.Name))
		require.NoError(t, err)

		require.Equal(t, 1, len(procedures))
		require.Contains(t, procedures, *f1)
		require.NotContains(t, procedures, *f2)
	})

	t.Run("show procedure for SQL: no matches", func(t *testing.T) {
		procedures, err := client.Procedures.Show(ctx, sdk.NewShowProcedureRequest().WithLike(random.String()))
		require.NoError(t, err)
		require.Equal(t, 0, len(procedures))
	})

	t.Run("describe function for SQL", func(t *testing.T) {
		f := createProcedureForSQLHandle(t, true)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, f.Name)

		request := sdk.NewDescribeProcedureRequest(id).WithArgumentTypes(defaultArgumentTypes)
		details, err := client.Procedures.Describe(ctx, request)
		require.NoError(t, err)
		require.Greater(t, len(details), 0)
	})
}
