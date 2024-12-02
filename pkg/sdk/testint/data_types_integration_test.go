package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DataTypes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	incorrectBooleanDatatypes := []string{
		"BOOLEAN()",
		"BOOLEAN(1)",
		"BOOL",
	}
	incorrectFloatDatatypes := []string{
		"DOUBLE()",
		"DOUBLE(1)",
		"DOUBLE PRECISION(1)",
	}
	incorrectlyCorrectFloatDatatypes := []string{
		"FLOAT()",
		"FLOAT(20)",
		"FLOAT4(20)",
		"FLOAT8(20)",
		"REAL(20)",
	}
	incorrectNumberDatatypes := []string{
		"NUMBER()",
		"NUMBER(x)",
		"INT()",
		"NUMBER(36, 5, 7)",
	}

	for _, c := range datatypes.ArrayDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of array datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT []::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT []::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
		})
	}

	for _, c := range datatypes.BinaryDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of binary datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TO_BINARY('AB')::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT TO_BINARY('AB')::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT TO_BINARY('AB')::%s(36, 2)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
		})
	}

	for _, c := range datatypes.BooleanDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of boolean datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TRUE::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range incorrectBooleanDatatypes {
		t.Run(fmt.Sprintf("check behavior of boolean datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT TRUE::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	for _, c := range datatypes.DateDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of date datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT '2024-12-02'::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.FloatDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range incorrectFloatDatatypes {
		t.Run(fmt.Sprintf("check behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	// There is no attribute documented for float numbers: https://docs.snowflake.com/en/sql-reference/data-types-numeric#float-float4-float8.
	// However, adding it succeeds for FLOAT, FLOAT4, FLOAT8, and REAL, but ift fails both for DOUBLE and DOUBLE PRECISION.
	for _, c := range incorrectlyCorrectFloatDatatypes {
		t.Run(fmt.Sprintf("document incorrect behavior of float datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1.1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.NoError(t, err)
		})
	}

	for _, c := range datatypes.GeographyDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of geography datatype: %s", c), func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql := fmt.Sprintf("CREATE TABLE %s (i %s)", tableId.FullyQualifiedName(), c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))

			tableId = testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql = fmt.Sprintf("CREATE TABLE %s (i %s())", tableId.FullyQualifiedName(), c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
		})
	}

	for _, c := range datatypes.GeometryDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of geometry datatype: %s", c), func(t *testing.T) {
			tableId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql := fmt.Sprintf("CREATE TABLE %s (i %s)", tableId.FullyQualifiedName(), c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))

			tableId = testClientHelper().Ids.RandomSchemaObjectIdentifier()
			sql = fmt.Sprintf("CREATE TABLE %s (i %s())", tableId.FullyQualifiedName(), c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
			t.Cleanup(testClientHelper().Table.DropFunc(t, tableId))
		})
	}

	for _, c := range datatypes.NumberDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of number datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36, 5)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)
		})
	}

	for _, c := range datatypes.NumberDataTypeSubTypes {
		t.Run(fmt.Sprintf("check behavior of number data type subtype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT 1::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
		})
	}

	for _, c := range incorrectNumberDatatypes {
		t.Run(fmt.Sprintf("check behavior of number datatype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT 1::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			require.Error(t, err)
		})
	}

	for _, c := range datatypes.ObjectDataTypeSynonyms {
		t.Run(fmt.Sprintf("check behavior of object data type subtype: %s", c), func(t *testing.T) {
			sql := fmt.Sprintf("SELECT {}::%s", c)
			_, err := client.QueryUnsafe(ctx, sql)
			assert.NoError(t, err)

			sql = fmt.Sprintf("SELECT {}::%s(36)", c)
			_, err = client.QueryUnsafe(ctx, sql)
			assert.Error(t, err)
		})
	}
}
