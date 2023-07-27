package sdk

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalTablesCreate1(t *testing.T) {
}

func TestExternalTablesCreate2(t *testing.T) {
}

func TestExternalTablesCreate3(t *testing.T) {
}

func TestExternalTablesAlter(t *testing.T) {
	t.Run("refresh", func(t *testing.T) {
		opts := AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			Refresh: &ExternalTableRefresh{
				RelativePath: String("some/path"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE IF EXISTS "external_table" REFRESH 'some/path'`
		assert.Equal(t, expected, actual)
	})

	t.Run("add files", func(t *testing.T) {
		opts := AlterExternalTableOptions{
			name:     NewAccountObjectIdentifier("external_table"),
			AddFiles: []string{"one/file.txt", "scond/file.txt"},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" ADD FILES ('one/file.txt', 'second/file.txt')`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove files", func(t *testing.T) {
		opts := AlterExternalTableOptions{
			name:        NewAccountObjectIdentifier("external_table"),
			RemoveFiles: []string{"one/file.txt", "scond/file.txt"},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" REMOVE FILES ('one/file.txt', 'second/file.txt')`
		assert.Equal(t, expected, actual)
	})

	t.Run("set", func(t *testing.T) {
		opts := AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			Set: &ExternalTableSet{
				AutoRefresh: Bool(true),
				Tag: []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tag1"),
						Value: "tag_value1",
					},
					{
						Name:  NewAccountObjectIdentifier("tag2"),
						Value: "tag_value2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" SET AUTO_REFRESH = TRUE TAG ("tag1" = 'tag_value1', "tag2" = 'tag_value2')`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset", func(t *testing.T) {
		opts := AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			Unset: &ExternalTableUnset{
				Tag: []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" UNSET TAG "tag1", "tag2"`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesAlterPartitions(t *testing.T) {
	t.Run("add partition", func(t *testing.T) {
		opts := AlterExternalTablePartitionOptions{
			name:     NewAccountObjectIdentifier("external_table"),
			IfExists: Bool(true),
			AddPartitions: []Partition{
				{
					ColumnName: "one",
					Value:      "one_value",
				},
				{
					ColumnName: "two",
					Value:      "two_value",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" IF EXISTS ADD PARTITION ("one" = 'one_value', "two" = 'two_value')`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove partition", func(t *testing.T) {
		opts := AlterExternalTablePartitionOptions{
			name:          NewAccountObjectIdentifier("external_table"),
			IfExists:      Bool(true),
			DropPartition: String("partition_location"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" IF EXISTS DROP PARTITION LOCATION 'partition_location'`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesDrop(t *testing.T) {
	t.Run("restrict", func(t *testing.T) {
		opts := DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP EXTERNAL TABLE IF EXISTS "external_table" RESTRICT`
		assert.Equal(t, expected, actual)
	})

	t.Run("cascade", func(t *testing.T) {
		opts := DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			DropOption: &ExternalTableDropOption{
				Cascade: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP EXTERNAL TABLE IF EXISTS "external_table" CASCADE`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesShow(t *testing.T) {
	t.Run("all options", func(t *testing.T) {
		opts := ShowExternalTableOptions{
			Terse: Bool(true),
			Like: &Like{
				Pattern: String("some_pattern"),
			},
			In: &In{
				Account: Bool(true),
			},
			StartsWith: String("some_external_table"),
			LimitFrom: &LimitFrom{
				Rows: Int(123),
				From: String("some_string"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES LIKE 'some_pattern' IN ACCOUNT STARTS WITH 'some_external_table' LIMIT 123 FROM 'some_string'`
		assert.Equal(t, expected, actual)
	})

	t.Run("in database", func(t *testing.T) {
		opts := ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Database: NewAccountObjectIdentifier("database_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES IN DATABASE "database_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("in schema", func(t *testing.T) {
		opts := ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Schema: NewSchemaIdentifier("database_name", "schema_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES IN SCHEMA "database_name"."schema_name"`
		assert.Equal(t, expected, actual)
	})
}

type A struct {
	a string
	b *string
}

func TestExternalTablesDescribe(t *testing.T) {
	a := A{}
	_ = a

	t.Run("type columns", func(t *testing.T) {
		opts := DescribeExternalTableOptions{
			name:        NewAccountObjectIdentifier("external_table"),
			ColumnsType: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE EXTERNAL TABLE "external_table" TYPE = COLUMNS`
		assert.Equal(t, expected, actual)
	})

	t.Run("type stage", func(t *testing.T) {
		opts := DescribeExternalTableOptions{
			name:      NewAccountObjectIdentifier("external_table"),
			StageType: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE EXTERNAL TABLE "external_table" TYPE = STAGE`
		assert.Equal(t, expected, actual)
	})
}

func copyWithFields[In, Opts interface{}](in In, opts Opts) *In {
	inType := reflect.ValueOf(&in).Elem()
	optsType := reflect.ValueOf(&opts).Elem()

	if inType.Kind() != reflect.Struct || optsType.Kind() != reflect.Struct {
		panic("...")
	}

	for i := 0; i < optsType.NumField(); i++ {
		optsFieldMeta := optsType.Type().Field(i)
		optsFieldValue := optsType.Field(i)
		for j := 0; j < inType.NumField(); j++ {
			inFieldMeta := inType.Type().Field(j)
			inFieldValue := inType.Field(j)
			if optsFieldMeta.Name == inFieldMeta.Name {
				fmt.Printf("Setting a: '%s', b: '%s'\n", inFieldMeta.Name, optsFieldMeta.Name)
				inFieldValue.Set(optsFieldValue)
			}
		}
	}

	return &in
}
