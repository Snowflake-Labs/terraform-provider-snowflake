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

	r.Equal(`CREATE FILE FORMAT "test_db"."test_schema"."test_file_format" TYPE = 'CSV' COMPRESSION = 'AUTO' RECORD_DELIMITER = '\n' FIELD_DELIMITER = ',' FILE_EXTENSION = '.CSV' SKIP_HEADER = 1 DATE_FORMAT = 'AUTO' TIME_FORMAT = 'AUTO' TIMESTAMP_FORMAT = 'AUTO' BINARY_FORMAT = 'HEX' ESCAPE = 'None' ESCAPE_UNENCLOSED_FIELD = 'None' FIELD_OPTIONALLY_ENCLOSED_BY = '"' NULL_IF = () ENCODING = 'UTF8' SKIP_BLANK_LINES = true TRIM_SPACE = true ERROR_ON_COLUMN_COUNT_MISMATCH = false REPLACE_INVALID_CHARACTERS = false VALIDATE_UTF8 = false EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true`, f.Create())
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

	r.Equal(`CREATE FILE FORMAT "test_db"."test_schema"."test_file_format_json" TYPE = 'JSON' COMPRESSION = 'AUTO' DATE_FORMAT = 'AUTO' TIME_FORMAT = 'AUTO' TIMESTAMP_FORMAT = 'AUTO' BINARY_FORMAT = 'HEX' NULL_IF = ('\n', 'NULL') TRIM_SPACE = true ENABLE_OCTAL = false ALLOW_DUPLICATE = false STRIP_OUTER_ARRAY = false STRIP_NULL_VALUES = false REPLACE_INVALID_CHARACTERS = false IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = false`, f.Create())
}

func TestFileFormatChangeComment(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET COMMENT = 'worst format ever'`, f.ChangeComment("worst format ever"))
}

func TestFileFormatChangeCompression(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET COMPRESSION = 'GZIP'`, f.ChangeCompression("GZIP"))
}

func TestFileFormatChangeRecordDelimiter(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET RECORD_DELIMITER = '|'`, f.ChangeRecordDelimiter("|"))
}

func TestFileFormatChangeFieldDelimiter(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FIELD_DELIMITER = '|'`, f.ChangeFieldDelimiter("|"))
}

func TestFileFormatChangeFileExtension(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FILE_EXTENSION = '.csv.gz'`, f.ChangeFileExtension(".csv.gz"))
}

func TestFileFormatChangeDateFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DATE_FORMAT = 'AUTO'`, f.ChangeDateFormat("AUTO"))
}

func TestFileFormatChangeTimeFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TIME_FORMAT = 'AUTO'`, f.ChangeTimeFormat("AUTO"))
}

func TestFileFormatChangeTimestampFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TIMESTAMP_FORMAT = 'AUTO'`, f.ChangeTimestampFormat("AUTO"))
}

func TestFileFormatChangeBinaryFormat(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET BINARY_FORMAT = 'AUTO'`, f.ChangeBinaryFormat("AUTO"))
}

func TestFileFormatChangeValidateUTF8(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET VALIDATE_UTF8 = true`, f.ChangeValidateUTF8(true))
}

func TestFileFormatChangeEmptyFieldAsNull(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET EMPTY_FIELD_AS_NULL = true`, f.ChangeEmptyFieldAsNull(true))
}

func TestFileFormatChangeErrorOnColumnCountMismatch(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ERROR_ON_COLUMN_COUNT_MISMATCH = true`, f.ChangeErrorOnColumnCountMismatch(true))
}

func TestFileFormatChangeEscape(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ESCAPE = 'None'`, f.ChangeEscape("None"))
}

func TestFileFormatChangeFieldOptionallyEnclosedBy(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET FIELD_OPTIONALLY_ENCLOSED_BY = 'None'`, f.ChangeFieldOptionallyEnclosedBy("None"))
}

func TestFileFormatChangeEscapeUnenclosedField(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ESCAPE_UNENCLOSED_FIELD = '\'`, f.ChangeEscapeUnenclosedField("\\"))
}

func TestFileFormatChangeNullIf(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET NULL_IF = ()`, f.ChangeNullIf([]string{}))
}

func TestFileFormatChangeEncoding(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ENCODING = 'UTF8'`, f.ChangeEncoding("UTF8"))
}

func TestFileFormatChangeSkipHeader(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_HEADER = 2`, f.ChangeSkipHeader(2))
}

func TestFileFormatChangeSkipBlankLines(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_BLANK_LINES = true`, f.ChangeSkipBlankLines(true))
}

func TestFileFormatChangeTrimSpace(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET TRIM_SPACE = true`, f.ChangeTrimSpace(true))
}

func TestFileFormatChangeEnableOctal(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ENABLE_OCTAL = true`, f.ChangeEnableOctal(true))
}

func TestFileFormatChangeAllowDuplicate(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET ALLOW_DUPLICATE = true`, f.ChangeAllowDuplicate(true))
}

func TestFileFormatChangeStripOuterArray(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_OUTER_ARRAY = false`, f.ChangeStripOuterArray(false))
}

func TestFileFormatChangeStripNullValues(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_NULL_VALUES = false`, f.ChangeStripNullValues(false))
}

func TestFileFormatChangeReplaceInvalidCharacters(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET REPLACE_INVALID_CHARACTERS = false`, f.ChangeReplaceInvalidCharacters(false))
}

func TestFileFormatChangeIgnoreUTF8Errors(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET IGNORE_UTF8_ERRORS = false`, f.ChangeIgnoreUTF8Errors(false))
}

func TestFileFormatChangeSkipByteOrderMark(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET SKIP_BYTE_ORDER_MARK = false`, f.ChangeSkipByteOrderMark(false))
}

func TestFileFormatChangeBinaryAsText(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET BINARY_AS_TEXT = false`, f.ChangeBinaryAsText(false))
}

func TestFileFormatChangePreserveSpace(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET PRESERVE_SPACE = false`, f.ChangePreserveSpace(false))
}

func TestFileFormatChangeStripOuterElement(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET STRIP_OUTER_ELEMENT = true`, f.ChangeStripOuterElement(true))
}

func TestFileFormatChangeDisableSnowflakeData(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DISABLE_SNOWFLAKE_DATA = true`, f.ChangeDisableSnowflakeData(true))
}

func TestFileFormatChangeDisableAutoConvert(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`ALTER FILE FORMAT "test_db"."test_schema"."test_file_format" SET DISABLE_AUTO_CONVERT = true`, f.ChangeDisableAutoConvert(true))
}

func TestFileFormatDrop(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`DROP FILE FORMAT "test_db"."test_schema"."test_file_format"`, f.Drop())
}

func TestFileFormatDescribe(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`DESCRIBE FILE FORMAT "test_db"."test_schema"."test_file_format"`, f.Describe())
}

func TestFileFormatShow(t *testing.T) {
	r := require.New(t)
	f := FileFormat("test_file_format", "test_db", "test_schema")
	r.Equal(`SHOW FILE FORMATS LIKE 'test_file_format' IN SCHEMA "test_db"."test_schema"`, f.Show())
}
