package snowflake

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// FileFormatBuilder abstracts the creation of SQL queries for a Snowflake file format
type FileFormatBuilder struct {
	name                       string
	db                         string
	schema                     string
	formatType                 string
	compression                string
	recordDelimiter            string
	fieldDelimiter             string
	fileExtension              string
	skipHeader                 int
	skipBlankLines             bool
	dateFormat                 string
	timeFormat                 string
	timestampFormat            string
	binaryFormat               string
	escape                     string
	escapeUnenclosedField      string
	trimSpace                  bool
	fieldOptionallyEnclosedBy  string
	nullIf                     []string
	errorOnColumnCountMismatch bool
	replaceInvalidCharacters   bool
	validateUTF8               bool
	emptyFieldAsNull           bool
	skipByteOrderMark          bool
	encoding                   string
	enableOctal                bool
	allowDuplicate             bool
	stripOuterArray            bool
	stripNullValues            bool
	ignoreUTF8Errors           bool
	binaryAsText               bool
	preserveSpace              bool
	stripOuterElement          bool
	disableSnowflakeData       bool
	disableAutoConvert         bool
	comment                    string
}

// QualifiedName prepends the db and schema and escapes everything nicely
func (ffb *FileFormatBuilder) QualifiedName() string {
	var n strings.Builder

	n.WriteString(fmt.Sprintf(`"%v"."%v"."%v"`, ffb.db, ffb.schema, ffb.name))

	return n.String()
}

// WithFormatType adds a comment to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithFormatType(f string) *FileFormatBuilder {
	ffb.formatType = f
	return ffb
}

// WithCompression adds compression to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithCompression(c string) *FileFormatBuilder {
	ffb.compression = c
	return ffb
}

// WithRecordDelimiter adds a record delimiter to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithRecordDelimiter(r string) *FileFormatBuilder {
	ffb.recordDelimiter = r
	return ffb
}

// WithFieldDelimiter adds a field delimiter to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithFieldDelimiter(f string) *FileFormatBuilder {
	ffb.fieldDelimiter = f
	return ffb
}

// WithFileExtension adds a file extension to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithFileExtension(f string) *FileFormatBuilder {
	ffb.fileExtension = f
	return ffb
}

// WithSkipHeader adds skip header to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithSkipHeader(n int) *FileFormatBuilder {
	ffb.skipHeader = n
	return ffb
}

// WithSkipBlankLines adds skip blank lines to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithSkipBlankLines(n bool) *FileFormatBuilder {
	ffb.skipBlankLines = n
	return ffb
}

// WithDateFormat adds date format to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithDateFormat(s string) *FileFormatBuilder {
	ffb.dateFormat = s
	return ffb
}

// WithTimeFormat adds time format to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithTimeFormat(s string) *FileFormatBuilder {
	ffb.timeFormat = s
	return ffb
}

// WithTimestampFormat adds timestamp format to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithTimestampFormat(s string) *FileFormatBuilder {
	ffb.timestampFormat = s
	return ffb
}

// WithBinaryFormat adds binary format to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithBinaryFormat(s string) *FileFormatBuilder {
	ffb.binaryFormat = s
	return ffb
}

// WithEscape adds escape to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithEscape(s string) *FileFormatBuilder {
	ffb.escape = s
	return ffb
}

// WithEscapeUnenclosedField adds escape unenclosed field to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithEscapeUnenclosedField(s string) *FileFormatBuilder {
	ffb.escapeUnenclosedField = s
	return ffb
}

// WithTrimSpace adds trim space to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithTrimSpace(n bool) *FileFormatBuilder {
	ffb.trimSpace = n
	return ffb
}

// WithFieldOptionallyEnclosedBy adds field optionally enclosed by to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithFieldOptionallyEnclosedBy(s string) *FileFormatBuilder {
	ffb.fieldOptionallyEnclosedBy = s
	return ffb
}

// WithNullIf adds null if to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithNullIf(s []string) *FileFormatBuilder {
	ffb.nullIf = s
	return ffb
}

// WithErrorOnColumnCountMismatch adds error on column count mistmatch to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithErrorOnColumnCountMismatch(n bool) *FileFormatBuilder {
	ffb.errorOnColumnCountMismatch = n
	return ffb
}

// WithReplaceInvalidCharacters adds replace invalid characters to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithReplaceInvalidCharacters(n bool) *FileFormatBuilder {
	ffb.replaceInvalidCharacters = n
	return ffb
}

// WithValidateUTF8 adds validate utf8 to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithValidateUTF8(n bool) *FileFormatBuilder {
	ffb.validateUTF8 = n
	return ffb
}

// WithEmptyFieldAsNull adds empty field as null to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithEmptyFieldAsNull(n bool) *FileFormatBuilder {
	ffb.emptyFieldAsNull = n
	return ffb
}

// WithSkipByteOrderMark adds skip byte order mark to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithSkipByteOrderMark(n bool) *FileFormatBuilder {
	ffb.skipByteOrderMark = n
	return ffb
}

// WithEnableOctal adds enable octal to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithEnableOctal(n bool) *FileFormatBuilder {
	ffb.enableOctal = n
	return ffb
}

// WithAllowDuplicate adds allow duplicate to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithAllowDuplicate(n bool) *FileFormatBuilder {
	ffb.allowDuplicate = n
	return ffb
}

// WithStripOuterArray adds strip outer array to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithStripOuterArray(n bool) *FileFormatBuilder {
	ffb.stripOuterArray = n
	return ffb
}

// WithStripNullValues adds strip null values to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithStripNullValues(n bool) *FileFormatBuilder {
	ffb.stripNullValues = n
	return ffb
}

// WithIgnoreUTF8Errors adds ignore UTF8 errors to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithIgnoreUTF8Errors(n bool) *FileFormatBuilder {
	ffb.ignoreUTF8Errors = n
	return ffb
}

// WithBinaryAsText adds binary as text to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithBinaryAsText(n bool) *FileFormatBuilder {
	ffb.binaryAsText = n
	return ffb
}

// WithPreserveSpace adds preserve space to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithPreserveSpace(n bool) *FileFormatBuilder {
	ffb.preserveSpace = n
	return ffb
}

// WithStripOuterElement adds strip outer element to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithStripOuterElement(n bool) *FileFormatBuilder {
	ffb.stripOuterElement = n
	return ffb
}

// WithDisableSnowflakeData adds disable Snowflake data to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithDisableSnowflakeData(n bool) *FileFormatBuilder {
	ffb.disableSnowflakeData = n
	return ffb
}

// WithDisableAutoConvert adds disbale auto convert to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithDisableAutoConvert(n bool) *FileFormatBuilder {
	ffb.disableAutoConvert = n
	return ffb
}

// WithEncoding adds encoding to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithEncoding(e string) *FileFormatBuilder {
	ffb.encoding = e
	return ffb
}

// WithComment adds a comment to the FileFormatBuilder
func (ffb *FileFormatBuilder) WithComment(c string) *FileFormatBuilder {
	ffb.comment = c
	return ffb
}

// FileFormat returns a pointer to a Builder that abstracts the DDL operations for a file format.
//
// Supported DDL operations are:
//   - CREATE FILE FORMAT
//   - ALTER FILE FORMAT
//   - DROP FILE FORMAT
//   - SHOW FILE FORMATS
//   - DESCRIBE FILE FORMAT
//
// [Snowflake Reference](https://docs.snowflake.com/en/sql-reference/sql/create-file-format.html)
func FileFormat(name, db, schema string) *FileFormatBuilder {
	return &FileFormatBuilder{
		name:   name,
		db:     db,
		schema: schema,
	}
}

// Create returns the SQL query that will create a new file format.
func (ffb *FileFormatBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(`CREATE`)

	q.WriteString(fmt.Sprintf(` FILE FORMAT %v`, ffb.QualifiedName()))
	q.WriteString(fmt.Sprintf(` TYPE = '%v'`, ffb.formatType))

	if ffb.compression != "" {
		q.WriteString(fmt.Sprintf(` COMPRESSION = '%v'`, ffb.compression))
	}

	if ffb.recordDelimiter != "" {
		q.WriteString(fmt.Sprintf(` RECORD_DELIMITER = '%v'`, ffb.recordDelimiter))
	}

	if ffb.fieldDelimiter != "" {
		q.WriteString(fmt.Sprintf(` FIELD_DELIMITER = '%v'`, ffb.fieldDelimiter))
	}

	if ffb.fileExtension != "" {
		q.WriteString(fmt.Sprintf(` FILE_EXTENSION = '%v'`, ffb.fileExtension))
	}

	if ffb.skipHeader > 0 {
		q.WriteString(fmt.Sprintf(` SKIP_HEADER = %v`, ffb.skipHeader))
	}

	if ffb.dateFormat != "" {
		q.WriteString(fmt.Sprintf(` DATE_FORMAT = '%v'`, ffb.dateFormat))
	}

	if ffb.timeFormat != "" {
		q.WriteString(fmt.Sprintf(` TIME_FORMAT = '%v'`, ffb.timeFormat))
	}

	if ffb.timestampFormat != "" {
		q.WriteString(fmt.Sprintf(` TIMESTAMP_FORMAT = '%v'`, ffb.timestampFormat))
	}

	if ffb.binaryFormat != "" {
		q.WriteString(fmt.Sprintf(` BINARY_FORMAT = '%v'`, ffb.binaryFormat))
	}

	if ffb.escape != "" {
		q.WriteString(fmt.Sprintf(` ESCAPE = '%v'`, EscapeString(ffb.escape)))
	}

	if ffb.escapeUnenclosedField != "" {
		q.WriteString(fmt.Sprintf(` ESCAPE_UNENCLOSED_FIELD = '%v'`, ffb.escapeUnenclosedField))
	}

	if ffb.fieldOptionallyEnclosedBy != "" {
		q.WriteString(fmt.Sprintf(` FIELD_OPTIONALLY_ENCLOSED_BY = '%v'`, EscapeString(ffb.fieldOptionallyEnclosedBy)))
	}

	if len(ffb.nullIf) > 0 {
		nullIfStr := "'" + strings.Join(ffb.nullIf, "', '") + "'"
		q.WriteString(fmt.Sprintf(` NULL_IF = (%v)`, nullIfStr))
	} else if strings.ToUpper(ffb.formatType) != "XML" {
		q.WriteString(` NULL_IF = ()`)
	}

	if ffb.encoding != "" {
		q.WriteString(fmt.Sprintf(` ENCODING = '%v'`, ffb.encoding))
	}

	// set boolean values
	if ffb.formatType == "CSV" {
		q.WriteString(fmt.Sprintf(` SKIP_BLANK_LINES = %v`, ffb.skipBlankLines))
		q.WriteString(fmt.Sprintf(` TRIM_SPACE = %v`, ffb.trimSpace))
		q.WriteString(fmt.Sprintf(` ERROR_ON_COLUMN_COUNT_MISMATCH = %v`, ffb.errorOnColumnCountMismatch))
		q.WriteString(fmt.Sprintf(` REPLACE_INVALID_CHARACTERS = %v`, ffb.replaceInvalidCharacters))
		q.WriteString(fmt.Sprintf(` VALIDATE_UTF8 = %v`, ffb.validateUTF8))
		q.WriteString(fmt.Sprintf(` EMPTY_FIELD_AS_NULL = %v`, ffb.emptyFieldAsNull))
		q.WriteString(fmt.Sprintf(` SKIP_BYTE_ORDER_MARK = %v`, ffb.skipByteOrderMark))
	} else if ffb.formatType == "JSON" {
		q.WriteString(fmt.Sprintf(` TRIM_SPACE = %v`, ffb.trimSpace))
		q.WriteString(fmt.Sprintf(` ENABLE_OCTAL = %v`, ffb.enableOctal))
		q.WriteString(fmt.Sprintf(` ALLOW_DUPLICATE = %v`, ffb.allowDuplicate))
		q.WriteString(fmt.Sprintf(` STRIP_OUTER_ARRAY = %v`, ffb.stripOuterArray))
		q.WriteString(fmt.Sprintf(` STRIP_NULL_VALUES = %v`, ffb.stripNullValues))
		q.WriteString(fmt.Sprintf(` REPLACE_INVALID_CHARACTERS = %v`, ffb.replaceInvalidCharacters))
		q.WriteString(fmt.Sprintf(` IGNORE_UTF8_ERRORS = %v`, ffb.ignoreUTF8Errors))
		q.WriteString(fmt.Sprintf(` SKIP_BYTE_ORDER_MARK = %v`, ffb.skipByteOrderMark))
	} else if ffb.formatType == "AVRO" || ffb.formatType == "ORC" {
		q.WriteString(fmt.Sprintf(` TRIM_SPACE = %v`, ffb.trimSpace))
	} else if ffb.formatType == "PARQUET" {
		q.WriteString(fmt.Sprintf(` BINARY_AS_TEXT = %v`, ffb.binaryAsText))
		q.WriteString(fmt.Sprintf(` TRIM_SPACE = %v`, ffb.trimSpace))
	} else if ffb.formatType == "XML" {
		q.WriteString(fmt.Sprintf(` IGNORE_UTF8_ERRORS = %v`, ffb.ignoreUTF8Errors))
		q.WriteString(fmt.Sprintf(` PRESERVE_SPACE = %v`, ffb.preserveSpace))
		q.WriteString(fmt.Sprintf(` STRIP_OUTER_ELEMENT = %v`, ffb.stripOuterElement))
		q.WriteString(fmt.Sprintf(` DISABLE_SNOWFLAKE_DATA = %v`, ffb.disableSnowflakeData))
		q.WriteString(fmt.Sprintf(` DISABLE_AUTO_CONVERT = %v`, ffb.disableAutoConvert))
		q.WriteString(fmt.Sprintf(` SKIP_BYTE_ORDER_MARK = %v`, ffb.skipByteOrderMark))
	}

	if ffb.comment != "" {
		q.WriteString(fmt.Sprintf(` COMMENT = '%v'`, EscapeString(ffb.comment)))
	}

	return q.String()
}

// ChangeComment returns the SQL query that will update the comment on the file format.
func (ffb *FileFormatBuilder) ChangeComment(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET COMMENT = '%v'`, ffb.QualifiedName(), c)
}

// RemoveComment returns the SQL query that will remove the comment on the file format.
func (ffb *FileFormatBuilder) RemoveComment() string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v UNSET COMMENT`, ffb.QualifiedName())
}

// ChangeCompression returns the SQL query that will update the compression on the file format.
func (ffb *FileFormatBuilder) ChangeCompression(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET COMPRESSION = '%v'`, ffb.QualifiedName(), c)
}

// ChangeRecordDelimiter returns the SQL query that will update the record delimiter on the file format.
func (ffb *FileFormatBuilder) ChangeRecordDelimiter(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET RECORD_DELIMITER = '%v'`, ffb.QualifiedName(), c)
}

// ChangeDateFormat returns the SQL query that will update the date format on the file format.
func (ffb *FileFormatBuilder) ChangeDateFormat(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET DATE_FORMAT = '%v'`, ffb.QualifiedName(), c)
}

// ChangeTimeFormat returns the SQL query that will update the time format on the file format.
func (ffb *FileFormatBuilder) ChangeTimeFormat(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET TIME_FORMAT = '%v'`, ffb.QualifiedName(), c)
}

// ChangeTimestampFormat returns the SQL query that will update the timestamp format on the file format.
func (ffb *FileFormatBuilder) ChangeTimestampFormat(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET TIMESTAMP_FORMAT = '%v'`, ffb.QualifiedName(), c)
}

// ChangeBinaryFormat returns the SQL query that will update the binary format on the file format.
func (ffb *FileFormatBuilder) ChangeBinaryFormat(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET BINARY_FORMAT = '%v'`, ffb.QualifiedName(), c)
}

// ChangeErrorOnColumnCountMismatch returns the SQL query that will update the error_on_column_count_mismatch on the file format.
func (ffb *FileFormatBuilder) ChangeErrorOnColumnCountMismatch(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ERROR_ON_COLUMN_COUNT_MISMATCH = %v`, ffb.QualifiedName(), c)
}

// ChangeValidateUTF8 returns the SQL query that will update the error_on_column_count_mismatch on the file format.
func (ffb *FileFormatBuilder) ChangeValidateUTF8(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET VALIDATE_UTF8 = %v`, ffb.QualifiedName(), c)
}

// ChangeEmptyFieldAsNull returns the SQL query that will update the error_on_column_count_mismatch on the file format.
func (ffb *FileFormatBuilder) ChangeEmptyFieldAsNull(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET EMPTY_FIELD_AS_NULL = %v`, ffb.QualifiedName(), c)
}

// ChangeEscape returns the SQL query that will update the escape on the file format.
func (ffb *FileFormatBuilder) ChangeEscape(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ESCAPE = '%v'`, ffb.QualifiedName(), c)
}

// ChangeEscapeUnenclosedField returns the SQL query that will update the escape unenclosed field on the file format.
func (ffb *FileFormatBuilder) ChangeEscapeUnenclosedField(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ESCAPE_UNENCLOSED_FIELD = '%v'`, ffb.QualifiedName(), c)
}

// ChangeFileExtension returns the SQL query that will update the FILE_EXTENSION on the file format.
func (ffb *FileFormatBuilder) ChangeFileExtension(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET FILE_EXTENSION = '%v'`, ffb.QualifiedName(), c)
}

// ChangeFieldDelimiter returns the SQL query that will update the FIELD_DELIMITER on the file format.
func (ffb *FileFormatBuilder) ChangeFieldDelimiter(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET FIELD_DELIMITER = '%v'`, ffb.QualifiedName(), c)
}

// ChangeFieldOptionallyEnclosedBy returns the SQL query that will update the field optionally enclosed by on the file format.
func (ffb *FileFormatBuilder) ChangeFieldOptionallyEnclosedBy(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET FIELD_OPTIONALLY_ENCLOSED_BY = '%v'`, ffb.QualifiedName(), c)
}

// ChangeNullIf returns the SQL query that will update the null if on the file format.
func (ffb *FileFormatBuilder) ChangeNullIf(c []string) string {
	nullIfStr := ""
	if len(c) > 0 {
		nullIfStr = "'" + strings.Join(c, "', '") + "'"
	}
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET NULL_IF = (%v)`, ffb.QualifiedName(), nullIfStr)
}

// ChangeEncoding returns the SQL query that will update the encoding on the file format.
func (ffb *FileFormatBuilder) ChangeEncoding(c string) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ENCODING = '%v'`, ffb.QualifiedName(), c)
}

// ChangeSkipHeader returns the SQL query that will update the skip header on the file format.
func (ffb *FileFormatBuilder) ChangeSkipHeader(c int) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET SKIP_HEADER = %v`, ffb.QualifiedName(), c)
}

// ChangeSkipBlankLines returns the SQL query that will update SKIP_BLANK_LINES on the file format.
func (ffb *FileFormatBuilder) ChangeSkipBlankLines(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET SKIP_BLANK_LINES = %v`, ffb.QualifiedName(), c)
}

// ChangeTrimSpace returns the SQL query that will update TRIM_SPACE on the file format.
func (ffb *FileFormatBuilder) ChangeTrimSpace(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET TRIM_SPACE = %v`, ffb.QualifiedName(), c)
}

// ChangeEnableOctal returns the SQL query that will update ENABLE_OCTAL on the file format.
func (ffb *FileFormatBuilder) ChangeEnableOctal(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ENABLE_OCTAL = %v`, ffb.QualifiedName(), c)
}

// ChangeAllowDuplicate returns the SQL query that will update ALLOW_DUPLICATE on the file format.
func (ffb *FileFormatBuilder) ChangeAllowDuplicate(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET ALLOW_DUPLICATE = %v`, ffb.QualifiedName(), c)
}

// ChangeStripOuterArray returns the SQL query that will update STRIP_OUTER_ARRAY on the file format.
func (ffb *FileFormatBuilder) ChangeStripOuterArray(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET STRIP_OUTER_ARRAY = %v`, ffb.QualifiedName(), c)
}

// ChangeStripNullValues returns the SQL query that will update STRIP_NULL_VALUES on the file format.
func (ffb *FileFormatBuilder) ChangeStripNullValues(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET STRIP_NULL_VALUES = %v`, ffb.QualifiedName(), c)
}

// ChangeReplaceInvalidCharacters returns the SQL query that will update REPLACE_INVALID_CHARACTERS on the file format.
func (ffb *FileFormatBuilder) ChangeReplaceInvalidCharacters(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET REPLACE_INVALID_CHARACTERS = %v`, ffb.QualifiedName(), c)
}

// ChangeIgnoreUTF8Errors returns the SQL query that will update IGNORE_UTF8_ERRORS on the file format.
func (ffb *FileFormatBuilder) ChangeIgnoreUTF8Errors(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET IGNORE_UTF8_ERRORS = %v`, ffb.QualifiedName(), c)
}

// ChangeSkipByteOrderMark returns the SQL query that will update SKIP_BYTE_ORDER_MARK on the file format.
func (ffb *FileFormatBuilder) ChangeSkipByteOrderMark(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET SKIP_BYTE_ORDER_MARK = %v`, ffb.QualifiedName(), c)
}

// ChangeBinaryAsText returns the SQL query that will update BINARY_AS_TEXT on the file format.
func (ffb *FileFormatBuilder) ChangeBinaryAsText(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET BINARY_AS_TEXT = %v`, ffb.QualifiedName(), c)
}

// ChangePreserveSpace returns the SQL query that will update PRESERVE_SPACE on the file format.
func (ffb *FileFormatBuilder) ChangePreserveSpace(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET PRESERVE_SPACE = %v`, ffb.QualifiedName(), c)
}

// ChangeStripOuterElement returns the SQL query that will update STRIP_OUTER_ELEMENT on the file format.
func (ffb *FileFormatBuilder) ChangeStripOuterElement(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET STRIP_OUTER_ELEMENT = %v`, ffb.QualifiedName(), c)
}

// ChangeDisableSnowflakeData returns the SQL query that will update DISABLE_SNOWFLAKE_DATA on the file format.
func (ffb *FileFormatBuilder) ChangeDisableSnowflakeData(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET DISABLE_SNOWFLAKE_DATA = %v`, ffb.QualifiedName(), c)
}

// ChangeDisableAutoConvert returns the SQL query that will update DISABLE_AUTO_CONVERT on the file format.
func (ffb *FileFormatBuilder) ChangeDisableAutoConvert(c bool) string {
	return fmt.Sprintf(`ALTER FILE FORMAT %v SET DISABLE_AUTO_CONVERT = %v`, ffb.QualifiedName(), c)
}

// Drop returns the SQL query that will drop a file format.
func (ffb *FileFormatBuilder) Drop() string {
	return fmt.Sprintf(`DROP FILE FORMAT %v`, ffb.QualifiedName())
}

// Describe returns the SQL query that will describe a file format..
func (ffb *FileFormatBuilder) Describe() string {
	return fmt.Sprintf(`DESCRIBE FILE FORMAT %v`, ffb.QualifiedName())
}

// Show returns the SQL query that will show a file format.
func (ffb *FileFormatBuilder) Show() string {
	return fmt.Sprintf(`SHOW FILE FORMATS LIKE '%v' IN DATABASE "%v"`, ffb.name, ffb.db)
}

type fileFormatShow struct {
	CreatedOn      sql.NullString `db:"created_on"`
	FileFormatName sql.NullString `db:"name"`
	DatabaseName   sql.NullString `db:"database_name"`
	SchemaName     sql.NullString `db:"schema_name"`
	FormatType     sql.NullString `db:"type"`
	Owner          sql.NullString `db:"owner"`
	Comment        sql.NullString `db:"comment"`
	FormatOptions  sql.NullString `db:"format_options"`
}

type fileFormatOptions struct {
	Type                       string   `json:"TYPE"`
	Compression                string   `json:"COMPRESSION,omitempty"`
	RecordDelimiter            string   `json:"RECORD_DELIMITER,omitempty"`
	FieldDelimiter             string   `json:"FIELD_DELIMITER,omitempty"`
	FileExtension              string   `json:"FILE_EXTENSION,omitempty"`
	SkipHeader                 int      `json:"SKIP_HEADER,omitempty"`
	DateFormat                 string   `json:"DATE_FORMAT,omitempty"`
	TimeFormat                 string   `json:"TIME_FORMAT,omitempty"`
	TimestampFormat            string   `json:"TIMESTAMP_FORMAT,omitempty"`
	BinaryFormat               string   `json:"BINARY_FORMAT,omitempty"`
	Escape                     string   `json:"ESCAPE,omitempty"`
	EscapeUnenclosedField      string   `json:"ESCAPE_UNENCLOSED_FIELD,omitempty"`
	TrimSpace                  bool     `json:"TRIM_SPACE,omitempty"`
	FieldOptionallyEnclosedBy  string   `json:"FIELD_OPTIONALLY_ENCLOSED_BY,omitempty"`
	NullIf                     []string `json:"NULL_IF,omitempty"`
	ErrorOnColumnCountMismatch bool     `json:"ERROR_ON_COLUMN_COUNT_MISMATCH,omitempty"`
	ValidateUTF8               bool     `json:"VALIDATE_UTF8,omitempty"`
	SkipBlankLines             bool     `json:"SKIP_BLANK_LINES,omitempty"`
	ReplaceInvalidCharacters   bool     `json:"REPLACE_INVALID_CHARACTERS,omitempty"`
	EmptyFieldAsNull           bool     `json:"EMPTY_FIELD_AS_NULL,omitempty"`
	SkipByteOrderMark          bool     `json:"SKIP_BYTE_ORDER_MARK,omitempty"`
	Encoding                   string   `json:"ENCODING,omitempty"`
	EnabelOctal                bool     `json:"ENABLE_OCTAL,omitempty"`
	AllowDuplicate             bool     `json:"ALLOW_DUPLICATE,omitempty"`
	StripOuterArray            bool     `json:"STRIP_OUTER_ARRAY,omitempty"`
	StripNullValues            bool     `json:"STRIP_NULL_VALUES,omitempty"`
	IgnoreUTF8Errors           bool     `json:"IGNORE_UTF8_ERRORS,omitempty"`
	BinaryAsText               bool     `json:"BINARY_AS_TEXT,omitempty"`
	PreserveSpace              bool     `json:"PRESERVE_SPACE,omitempty"`
	StripOuterElement          bool     `json:"STRIP_OUTER_ELEMENT,omitempty"`
	DisableSnowflakeData       bool     `json:"DISABLE_SNOWFLAKE_DATA,omitempty"`
	DisableAutoConvert         bool     `json:"DISABLE_AUTO_CONVERT,omitempty"`
}

func ScanFileFormatShow(row *sqlx.Row) (*fileFormatShow, error) {
	r := &fileFormatShow{}
	err := row.StructScan(r)
	return r, err
}

func ParseFormatOptions(fileOptions string) (*fileFormatOptions, error) {
	ff := &fileFormatOptions{}
	err := json.Unmarshal([]byte(fileOptions), ff)
	return ff, err
}
