package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileFormatCreateCSV(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.QualifiedName(), `"test_db"."test_schema"."test_file_format"`)

	f.WithFormatType("CSV")
	f.WithCompression("AUTO")
	f.WithRecordDelimiter("\\n")
	f.WithFieldDelimiter(",")
	f.WithFileExtension(".CSV")
	f.WithSkipHeader(1)
	f.WithSkipBlankLines(true)
	f.WithDateFormat("AUTO")
	f.WithTimeFormat("AUTO")
	f.WithTimestampFormat("AUTO")
	f.WithBinaryFormat("HEX")
	f.WithEscape("None")
	f.WithEscapeUnenclosedField("None")
	f.WithTrimSpace(true)
	f.WithFieldOptionallyEnclosedBy("\"")
	f.WithNullIf([]string{})
	f.WithErrorOnColumnCountMismatch(false)
	f.WithReplaceInvalidCharacters(false)
	f.WithValidateUTF8(false)
	f.WithEmptyFieldAsNull(true)
	f.WithSkipByteOrderMark(true)
	f.WithEncoding("UTF8")

	r.Equal(f.Create(), `CREATE FILE FORMAT "test_db"."test_schema"."test_file_format" TYPE = 'CSV' COMPRESSION = 'AUTO' RECORD_DELIMITER = '\n' FIELD_DELIMITER = ',' FILE_EXTENSION = '.CSV' SKIP_HEADER = 1 DATE_FORMAT = 'AUTO' TIME_FORMAT = 'AUTO' TIMESTAMP_FORMAT = 'AUTO' BINARY_FORMAT = 'HEX' ESCAPE = 'None' ESCAPE_UNENCLOSED_FIELD = 'None' FIELD_OPTIONALLY_ENCLOSED_BY = '"' NULL_IF = () ENCODING = 'UTF8' SKIP_BLANK_LINES = true TRIM_SPACE = true ERROR_ON_COLUMN_COUNT_MISMATCH = false REPLACE_INVALID_CHARACTERS = false VALIDATE_UTF8 = false EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true`)
}

func TestFileFormatCreateJSON(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format_json", "test_db", "test_schema")
	r.Equal(f.QualifiedName(), `"test_db"."test_schema"."test_file_format_json"`)

	f.WithFormatType("JSON")
	f.WithCompression("AUTO")
	f.WithDateFormat("AUTO")
	f.WithTimeFormat("AUTO")
	f.WithTimestampFormat("AUTO")
	f.WithBinaryFormat("HEX")
	f.WithTrimSpace(true)
	f.WithNullIf([]string{"\\n", "NULL"})
	f.WithAllowDuplicate(false)
	f.WithStripOuterArray(false)
	f.WithStripNullValues(false)
	f.WithIgnoreUTF8Errors(true)

	r.Equal(f.Create(), `CREATE FILE FORMAT "test_db"."test_schema"."test_file_format_json" TYPE = 'JSON' COMPRESSION = 'AUTO' DATE_FORMAT = 'AUTO' TIME_FORMAT = 'AUTO' TIMESTAMP_FORMAT = 'AUTO' BINARY_FORMAT = 'HEX' NULL_IF = ('\n', 'NULL') TRIM_SPACE = true ENABLE_OCTAL = false ALLOW_DUPLICATE = false STRIP_OUTER_ARRAY = false STRIP_NULL_VALUES = false REPLACE_INVALID_CHARACTERS = false IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = false`)
}

func TestFileFormatChangeComment(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeComment("worst format ever"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET COMMENT = 'worst format ever'`)
}

func TestFileFormatChangeCompression(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeCompression("GZIP"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET COMPRESSION = 'GZIP'`)
}

func TestFileFormatChangeRecordDelimiter(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeRecordDelimiter("|"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET RECORD_DELIMITER = '|'`)
}

func TestFileFormatChangeFieldDelimiter(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeFieldDelimiter("|"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FIELD_DELIMITER = '|'`)
}

func TestFileFormatChangeFileExtension(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeFileExtension(".csv.gz"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FILE_EXTENSION = '.csv.gz'`)
}

func TestFileFormatChangeDateFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeDateFormat("AUTO"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DATE_FORMAT = 'AUTO'`)
}

func TestFileFormatChangeTimeFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeTimeFormat("AUTO"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TIME_FORMAT = 'AUTO'`)
}

func TestFileFormatChangeTimestampFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeTimestampFormat("AUTO"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TIMESTAMP_FORMAT = 'AUTO'`)
}

func TestFileFormatChangeBinaryFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeBinaryFormat("AUTO"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET BINARY_FORMAT = 'AUTO'`)
}

func TestFileFormatChangeValidateUTF8(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeValidateUTF8(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET VALIDATE_UTF8 = true`)
}

func TestFileFormatChangeEmptyFieldAsNull(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeEmptyFieldAsNull(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET EMPTY_FIELD_AS_NULL = true`)
}

func TestFileFormatChangeErrorOnColumnCountMismatch(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeErrorOnColumnCountMismatch(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ERROR_ON_COLUMN_COUNT_MISMATCH = true`)
}

func TestFileFormatChangeEscape(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeEscape("None"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ESCAPE = 'None'`)
}

func TestFileFormatChangeFieldOptionallyEnclosedBy(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeFieldOptionallyEnclosedBy("None"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FIELD_OPTIONALLY_ENCLOSED_BY = 'None'`)
}

func TestFileFormatChangeEscapeUnenclosedField(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeEscapeUnenclosedField("\\"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ESCAPE_UNENCLOSED_FIELD = '\'`)
}

func TestFileFormatChangeNullIf(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeNullIf([]string{}), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET NULL_IF = ()`)
}

func TestFileFormatChangeEncoding(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeEncoding("UTF8"), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ENCODING = 'UTF8'`)
}

func TestFileFormatChangeSkipHeader(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeSkipHeader(2), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_HEADER = 2`)
}

func TestFileFormatChangeSkipBlankLines(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeSkipBlankLines(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_BLANK_LINES = true`)
}

func TestFileFormatChangeTrimSpace(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeTrimSpace(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TRIM_SPACE = true`)
}

func TestFileFormatChangeEnableOctal(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeEnableOctal(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ENABLE_OCTAL = true`)
}

func TestFileFormatChangeAllowDuplicate(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeAllowDuplicate(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ALLOW_DUPLICATE = true`)
}

func TestFileFormatChangeStripOuterArray(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeStripOuterArray(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_OUTER_ARRAY = false`)
}

func TestFileFormatChangeStripNullValues(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeStripNullValues(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_NULL_VALUES = false`)
}

func TestFileFormatChangeReplaceInvalidCharacters(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeReplaceInvalidCharacters(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET REPLACE_INVALID_CHARACTERS = false`)
}

func TestFileFormatChangeIgnoreUTF8Errors(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeIgnoreUTF8Errors(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET IGNORE_UTF8_ERRORS = false`)
}

func TestFileFormatChangeSkipByteOrderMark(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeSkipByteOrderMark(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_BYTE_ORDER_MARK = false`)
}

func TestFileFormatChangeSnappyCompression(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeSnappyCompression(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SNAPPY_COMPRESSION = false`)
}

func TestFileFormatChangeBinaryAsText(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeBinaryAsText(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET BINARY_AS_TEXT = false`)
}

func TestFileFormatChangePreserveSpace(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangePreserveSpace(false), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET PRESERVE_SPACE = false`)
}

func TestFileFormatChangeStripOuterElement(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeStripOuterElement(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_OUTER_ELEMENT = true`)
}

func TestFileFormatChangeDisableSnowflakeData(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeDisableSnowflakeData(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DISABLE_SNOWFLAKE_DATA = true`)
}

func TestFileFormatChangeDisableAutoConvert(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.ChangeDisableAutoConvert(true), `ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DISABLE_AUTO_CONVERT = true`)
}

func TestFileFormatDrop(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.Drop(), `DROP FILE FORMAT "test_db"."test_schema"."test_file_format"`)
}

func TestFileFormatDescribe(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.Describe(), `DESCRIBE FILE FORMAT "test_db"."test_schema"."test_file_format"`)
}

func TestFileFormatShow(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(f.Show(), `SHOW FILE FORMATS LIKE 'test_file_format' IN DATABASE "test_db"`)
}
