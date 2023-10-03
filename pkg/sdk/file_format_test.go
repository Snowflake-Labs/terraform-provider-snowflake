package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileFormatsCreate(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: NewSchemaObjectIdentifier("db1", "schema2", "format3"),
			Type: FileFormatTypeCSV,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE FILE FORMAT "db1"."schema2"."format3" TYPE = CSV`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete CSV", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeCSV,

			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVCompression:               &CSVCompressionBz2,
				CSVRecordDelimiter:           String("-"),
				CSVFieldDelimiter:            String(":"),
				CSVFileExtension:             String("csv"),
				CSVParseHeader:               Bool(true),
				CSVSkipHeader:                Int(5),
				CSVSkipBlankLines:            Bool(true),
				CSVDateFormat:                String("YYYY-MM-DD"),
				CSVTimeFormat:                String("HH:mm:SS"),
				CSVTimestampFormat:           String("time"),
				CSVBinaryFormat:              &BinaryFormatUTF8,
				CSVEscape:                    String("\\"),
				CSVEscapeUnenclosedField:     String("§"),
				CSVTrimSpace:                 Bool(true),
				CSVFieldOptionallyEnclosedBy: String("&"),
				CSVNullIf: &[]NullString{
					{"nul"},
					{"nulll"},
				},
				CSVErrorOnColumnCountMismatch: Bool(true),
				CSVReplaceInvalidCharacters:   Bool(true),
				CSVEmptyFieldAsNull:           Bool(true),
				CSVSkipByteOrderMark:          Bool(true),
				CSVEncoding:                   &CSVEncodingISO2022KR,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = CSV COMPRESSION = BZ2 RECORD_DELIMITER = '-' FIELD_DELIMITER = ':' FILE_EXTENSION = 'csv' PARSE_HEADER = true SKIP_HEADER = 5 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'time' BINARY_FORMAT = UTF8 ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '§' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '&' NULL_IF = ('nul', 'nulll') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = 'ISO2022KR'`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete JSON", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeJSON,

			FileFormatTypeOptions: FileFormatTypeOptions{
				JSONCompression:     &JSONCompressionBrotli,
				JSONDateFormat:      String("YYYY-MM-DD"),
				JSONTimeFormat:      String("HH:mm:SS"),
				JSONTimestampFormat: String("aze"),
				JSONBinaryFormat:    &BinaryFormatHex,
				JSONTrimSpace:       Bool(true),
				JSONNullIf: &[]NullString{
					{"c1"},
					{"c2"},
				},
				JSONFileExtension:            String("json"),
				JSONEnableOctal:              Bool(true),
				JSONAllowDuplicate:           Bool(true),
				JSONStripOuterArray:          Bool(true),
				JSONStripNullValues:          Bool(true),
				JSONReplaceInvalidCharacters: Bool(true),
				JSONIgnoreUTF8Errors:         Bool(true),
				JSONSkipByteOrderMark:        Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = JSON COMPRESSION = BROTLI DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'aze' BINARY_FORMAT = HEX TRIM_SPACE = true NULL_IF = ('c1', 'c2') FILE_EXTENSION = 'json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true REPLACE_INVALID_CHARACTERS = true IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = true`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete Avro", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeAvro,

			FileFormatTypeOptions: FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionDeflate,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"nul"}},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = AVRO COMPRESSION = DEFLATE TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete ORC", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeORC,

			FileFormatTypeOptions: FileFormatTypeOptions{
				ORCTrimSpace:                Bool(true),
				ORCReplaceInvalidCharacters: Bool(true),
				ORCNullIf:                   &[]NullString{{"nul"}},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = ORC TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete Parquet", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeParquet,

			FileFormatTypeOptions: FileFormatTypeOptions{
				ParquetCompression:              &ParquetCompressionLzo,
				ParquetSnappyCompression:        Bool(true),
				ParquetBinaryAsText:             Bool(true),
				ParquetTrimSpace:                Bool(true),
				ParquetReplaceInvalidCharacters: Bool(true),
				ParquetNullIf:                   &[]NullString{{"nil"}},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = PARQUET COMPRESSION = LZO SNAPPY_COMPRESSION = true BINARY_AS_TEXT = true TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete XML", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeXML,

			FileFormatTypeOptions: FileFormatTypeOptions{
				XMLCompression:              &XMLCompressionZstd,
				XMLIgnoreUTF8Errors:         Bool(true),
				XMLPreserveSpace:            Bool(true),
				XMLStripOuterElement:        Bool(true),
				XMLDisableSnowflakeData:     Bool(true),
				XMLDisableAutoConvert:       Bool(true),
				XMLReplaceInvalidCharacters: Bool(true),
				XMLSkipByteOrderMark:        Bool(true),

				Comment: String("test file format"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS "db4"."schema5"."format6" TYPE = XML COMPRESSION = ZSTD IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_SNOWFLAKE_DATA = true DISABLE_AUTO_CONVERT = true REPLACE_INVALID_CHARACTERS = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'test file format'`
		assert.Equal(t, expected, actual)
	})

	t.Run("previous test", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: NewSchemaObjectIdentifier("test_db", "test_schema", "test_file_format"),
			Type: FileFormatTypeCSV,

			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVNullIf:                     &[]NullString{{"NULL"}},
				CSVSkipBlankLines:             Bool(false),
				CSVTrimSpace:                  Bool(false),
				CSVErrorOnColumnCountMismatch: Bool(true),
				CSVReplaceInvalidCharacters:   Bool(false),
				CSVEmptyFieldAsNull:           Bool(false),
				CSVSkipByteOrderMark:          Bool(false),

				Comment: String("great comment"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE FILE FORMAT "test_db"."test_schema"."test_file_format" TYPE = CSV SKIP_BLANK_LINES = false TRIM_SPACE = false NULL_IF = ('NULL') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = false EMPTY_FIELD_AS_NULL = false SKIP_BYTE_ORDER_MARK = false COMMENT = 'great comment'`
		assert.Equal(t, expected, actual)
	})
}

func TestFileFormatsAlter(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "fileformat"),
			Rename: &AlterFileFormatRenameOptions{
				NewName: NewSchemaObjectIdentifier("new_db", "new_schema", "new_fileformat"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FILE FORMAT IF EXISTS "db"."schema"."fileformat" RENAME TO "new_db"."new_schema"."new_fileformat"`
		assert.Equal(t, expected, actual)
	})

	t.Run("set", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "fileformat"),
			Set: &FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionBrotli,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"nil"}},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER FILE FORMAT IF EXISTS "db"."schema"."fileformat" SET COMPRESSION = BROTLI TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`
		assert.Equal(t, expected, actual)
	})
}

func TestFileFormatsDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropFileFormatOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "ff"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP FILE FORMAT "db"."schema"."ff"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with IfExists", func(t *testing.T) {
		opts := &DropFileFormatOptions{
			name:     NewSchemaObjectIdentifier("db", "schema", "ff"),
			IfExists: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP FILE FORMAT IF EXISTS "db"."schema"."ff"`
		assert.Equal(t, expected, actual)
	})
}

func TestFileFormatsShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowFileFormatsOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW FILE FORMATS`
		assert.Equal(t, expected, actual)
	})

	t.Run("with show options", func(t *testing.T) {
		opts := &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String("test"),
			},
			In: &In{
				Schema: NewDatabaseObjectIdentifier("db", "schema"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW FILE FORMATS LIKE 'test' IN SCHEMA "db"."schema"`
		assert.Equal(t, expected, actual)
	})
}

func TestFileFormatsShowById(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		id := NewSchemaObjectIdentifier("db", "schema", "ff")
		opts := &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: NewDatabaseObjectIdentifier(id.databaseName, id.schemaName),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW FILE FORMATS LIKE 'ff' IN SCHEMA "db"."schema"`
		assert.Equal(t, expected, actual)
	})
}

func TestFileFormatsDescribe(t *testing.T) {
	opts := &describeFileFormatOptions{
		name: NewSchemaObjectIdentifier("db", "schema", "ff"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `DESCRIBE FILE FORMAT "db"."schema"."ff"`
	assert.Equal(t, expected, actual)
}
