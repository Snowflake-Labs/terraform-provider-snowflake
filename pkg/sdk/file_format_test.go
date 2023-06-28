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
			Type: FileFormatTypeCsv,
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
			Type:        FileFormatTypeCsv,

			FileFormatTypeOptions: FileFormatTypeOptions{
				CsvCompression:               &CsvCompressionBz2,
				CsvRecordDelimiter:           String("-"),
				CsvFieldDelimiter:            String(":"),
				CsvFileExtension:             String("csv"),
				CsvParseHeader:               Bool(true),
				CsvSkipHeader:                Int(5),
				CsvSkipBlankLines:            Bool(true),
				CsvDateFormat:                String("YYYY-MM-DD"),
				CsvTimeFormat:                String("HH:mm:SS"),
				CsvTimestampFormat:           String("time"),
				CsvBinaryFormat:              &BinaryFormatUtf8,
				CsvEscape:                    String("\\"),
				CsvEscapeUnenclosedField:     String("§"),
				CsvTrimSpace:                 Bool(true),
				CsvFieldOptionallyEnclosedBy: String("&"),
				CsvNullIf: &[]NullString{
					{"nul"},
					{"nulll"},
				},
				CsvErrorOnColumnCountMismatch: Bool(true),
				CsvReplaceInvalidCharacters:   Bool(true),
				CsvEmptyFieldAsNull:           Bool(true),
				CsvSkipByteOrderMark:          Bool(true),
				CsvEncoding:                   &CsvEncodingISO2022KR,
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
			Type:        FileFormatTypeJson,

			FileFormatTypeOptions: FileFormatTypeOptions{
				JsonCompression:     &JsonCompressionBrotli,
				JsonDateFormat:      String("YYYY-MM-DD"),
				JsonTimeFormat:      String("HH:mm:SS"),
				JsonTimestampFormat: String("aze"),
				JsonBinaryFormat:    &BinaryFormatHex,
				JsonTrimSpace:       Bool(true),
				JsonNullIf: &[]NullString{
					{"c1"},
					{"c2"},
				},
				JsonFileExtension:            String("json"),
				JsonEnableOctal:              Bool(true),
				JsonAllowDuplicate:           Bool(true),
				JsonStripOuterArray:          Bool(true),
				JsonStripNullValues:          Bool(true),
				JsonReplaceInvalidCharacters: Bool(true),
				JsonIgnoreUtf8Errors:         Bool(true),
				JsonSkipByteOrderMark:        Bool(true),
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

	t.Run("complete Orc", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        NewSchemaObjectIdentifier("db4", "schema5", "format6"),
			IfNotExists: Bool(true),
			Type:        FileFormatTypeOrc,

			FileFormatTypeOptions: FileFormatTypeOptions{
				OrcTrimSpace:                Bool(true),
				OrcReplaceInvalidCharacters: Bool(true),
				OrcNullIf:                   &[]NullString{{"nul"}},
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
			Type:        FileFormatTypeXml,

			FileFormatTypeOptions: FileFormatTypeOptions{
				XmlCompression:              &XmlCompressionZstd,
				XmlIgnoreUtf8Errors:         Bool(true),
				XmlPreserveSpace:            Bool(true),
				XmlStripOuterElement:        Bool(true),
				XmlDisableSnowflakeData:     Bool(true),
				XmlDisableAutoConvert:       Bool(true),
				XmlReplaceInvalidCharacters: Bool(true),
				XmlSkipByteOrderMark:        Bool(true),

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
			Type: FileFormatTypeCsv,

			FileFormatTypeOptions: FileFormatTypeOptions{
				CsvNullIf:                     &[]NullString{{"NULL"}},
				CsvSkipBlankLines:             Bool(false),
				CsvTrimSpace:                  Bool(false),
				CsvErrorOnColumnCountMismatch: Bool(true),
				CsvReplaceInvalidCharacters:   Bool(false),
				CsvEmptyFieldAsNull:           Bool(false),
				CsvSkipByteOrderMark:          Bool(false),

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
				Schema: NewSchemaIdentifier("db", "schema"),
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
		opts := &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String(NewSchemaObjectIdentifier("db", "schema", "ff").Name()),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW FILE FORMATS LIKE 'ff'`
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
