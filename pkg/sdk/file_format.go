package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type FileFormats interface {
	// Create creates a FileFormat.
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptions) error
	// Alter modifies an existing FileFormat
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptions) error
	// Drop removes a FileFormat.
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptions) error
	// Show returns a list of fileFormats.
	Show(ctx context.Context, opts *ShowFileFormatsOptions) ([]*FileFormat, error)
	// ShowByID returns a FileFormat by ID
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error)
	// Describe returns the details of a FileFormat.
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetails, error)
}

var _ FileFormats = (*fileFormats)(nil)

type fileFormats struct {
	client *Client
}

type FileFormat struct {
	Name          SchemaObjectIdentifier
	CreatedOn     time.Time
	Type          FileFormatType
	Owner         string
	Comment       string
	OwnerRoleType string
	Options       FileFormatTypeOptions
}

func (v *FileFormat) ID() SchemaObjectIdentifier {
	return v.Name
}

func (v *FileFormat) ObjectType() ObjectType {
	return ObjectTypeFileFormat
}

type FileFormatRow struct {
	FormatOptions string    `db:"format_options"`
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	DatabaseName  string    `db:"database_name"`
	SchemaName    string    `db:"schema_name"`
	FormatType    string    `db:"type"`
	Owner         string    `db:"owner"`
	Comment       string    `db:"comment"`
	OwnerRoleType string    `db:"owner_role_type"`
}

type ShowFileFormatsOptionsResult struct {
	// CSV + shared fields
	Type                       string   `json:"TYPE"`
	RecordDelimiter            string   `json:"RECORD_DELIMITER"`
	FieldDelimiter             string   `json:"FIELD_DELIMITER"`
	FileExtension              string   `json:"FILE_EXTENSION"`
	SkipHeader                 int      `json:"SKIP_HEADER"`
	ParseHeader                bool     `json:"PARSE_HEADER"`
	DateFormat                 string   `json:"DATE_FORMAT"`
	TimeFormat                 string   `json:"TIME_FORMAT"`
	TimestampFormat            string   `json:"TIMESTAMP_FORMAT"`
	BinaryFormat               string   `json:"BINARY_FORMAT"`
	Escape                     string   `json:"ESCAPE"`
	EscapeUnenclosedField      string   `json:"ESCAPE_UNENCLOSED_FIELD"`
	TrimSpace                  bool     `json:"TRIM_SPACE"`
	FieldOptionallyEnclosedBy  string   `json:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	NullIf                     []string `json:"NULL_IF"`
	Compression                string   `json:"COMPRESSION"`
	ErrorOnColumnCountMismatch bool     `json:"ERROR_ON_COLUMN_COUNT_MISMATCH"`
	ValidateUtf8               bool     `json:"VALIDATE_UTF8"`
	SkipBlankLines             bool     `json:"SKIP_BLANK_LINES"`
	ReplaceInvalidCharacters   bool     `json:"REPLACE_INVALID_CHARACTERS"`
	EmptyFieldAsNull           bool     `json:"EMPTY_FIELD_AS_NULL"`
	SkipByteOrderMark          bool     `json:"SKIP_BYTE_ORDER_MARK"`
	Encoding                   string   `json:"ENCODING"`

	// JSON fields
	EnableOctal      bool `json:"ENABLE_OCTAL"`
	AllowDuplicate   bool `json:"ALLOW_DUPLICATE"`
	StripOuterArray  bool `json:"STRIP_OUTER_ARRAY"`
	StripNullValues  bool `json:"STRIP_NULL_VALUES"`
	IgnoreUtf8Errors bool `json:"IGNORE_UTF8_ERRORS"`

	// Parquet fields
	BinaryAsText bool `json:"BINARY_AS_TEXT"`

	// XML fields
	PreserveSpace        bool `json:"PRESERVE_SPACE"`
	StripOuterElement    bool `json:"STRIP_OUTER_ELEMENT"`
	DisableSnowflakeData bool `json:"DISABLE_SNOWFLAKE_DATA"`
	DisableAutoConvert   bool `json:"DISABLE_AUTO_CONVERT"`
}

func (row *FileFormatRow) toFileFormat() *FileFormat {
	inputOptions := ShowFileFormatsOptionsResult{}
	err := json.Unmarshal([]byte(row.FormatOptions), &inputOptions)
	if err != nil {
		fmt.Printf("%s", err)
		panic("cannot parse options json")
	}

	ff := &FileFormat{
		Name:          NewSchemaObjectIdentifier(row.DatabaseName, row.SchemaName, row.Name),
		CreatedOn:     row.CreatedOn,
		Type:          FileFormatType(row.FormatType),
		Owner:         row.Owner,
		Comment:       row.Comment,
		OwnerRoleType: row.OwnerRoleType,
		Options:       FileFormatTypeOptions{},
	}

	newNullIf := []NullString{}
	for _, s := range inputOptions.NullIf {
		newNullIf = append(newNullIf, NullString{s})
	}

	switch ff.Type {
	case FileFormatTypeCsv:
		ff.Options.CsvCompression = (*CsvCompression)(&inputOptions.Compression)
		ff.Options.CsvRecordDelimiter = &inputOptions.RecordDelimiter
		ff.Options.CsvFieldDelimiter = &inputOptions.FieldDelimiter
		ff.Options.CsvFileExtension = &inputOptions.FileExtension
		ff.Options.CsvParseHeader = &inputOptions.ParseHeader
		ff.Options.CsvSkipHeader = &inputOptions.SkipHeader
		ff.Options.CsvSkipBlankLines = &inputOptions.SkipBlankLines
		ff.Options.CsvDateFormat = &inputOptions.DateFormat
		ff.Options.CsvTimeFormat = &inputOptions.TimeFormat
		ff.Options.CsvTimestampFormat = &inputOptions.TimestampFormat
		ff.Options.CsvBinaryFormat = (*BinaryFormat)(&inputOptions.BinaryFormat)
		ff.Options.CsvEscape = &inputOptions.Escape
		ff.Options.CsvEscapeUnenclosedField = &inputOptions.EscapeUnenclosedField
		ff.Options.CsvTrimSpace = &inputOptions.TrimSpace
		ff.Options.CsvFieldOptionallyEnclosedBy = &inputOptions.FieldOptionallyEnclosedBy
		ff.Options.CsvNullIf = &newNullIf
		ff.Options.CsvErrorOnColumnCountMismatch = &inputOptions.ErrorOnColumnCountMismatch
		ff.Options.CsvReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.CsvEmptyFieldAsNull = &inputOptions.EmptyFieldAsNull
		ff.Options.CsvSkipByteOrderMark = &inputOptions.SkipByteOrderMark
		ff.Options.CsvEncoding = (*CsvEncoding)(&inputOptions.Encoding)
	case FileFormatTypeJson:
		ff.Options.JsonCompression = (*JsonCompression)(&inputOptions.Compression)
		ff.Options.JsonDateFormat = &inputOptions.DateFormat
		ff.Options.JsonTimeFormat = &inputOptions.TimeFormat
		ff.Options.JsonTimestampFormat = &inputOptions.TimestampFormat
		ff.Options.JsonBinaryFormat = (*BinaryFormat)(&inputOptions.BinaryFormat)
		ff.Options.JsonTrimSpace = &inputOptions.TrimSpace
		ff.Options.JsonNullIf = &newNullIf
		ff.Options.JsonFileExtension = &inputOptions.FileExtension
		ff.Options.JsonEnableOctal = &inputOptions.EnableOctal
		ff.Options.JsonAllowDuplicate = &inputOptions.AllowDuplicate
		ff.Options.JsonStripOuterArray = &inputOptions.StripOuterArray
		ff.Options.JsonStripNullValues = &inputOptions.StripNullValues
		ff.Options.JsonReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.JsonIgnoreUtf8Errors = &inputOptions.IgnoreUtf8Errors
		ff.Options.JsonSkipByteOrderMark = &inputOptions.SkipByteOrderMark
	case FileFormatTypeAvro:
		ff.Options.AvroTrimSpace = &inputOptions.TrimSpace
		ff.Options.AvroNullIf = &newNullIf
		ff.Options.AvroCompression = (*AvroCompression)(&inputOptions.Compression)
		ff.Options.AvroReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
	case FileFormatTypeOrc:
		ff.Options.OrcTrimSpace = &inputOptions.TrimSpace
		ff.Options.OrcReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.OrcNullIf = &newNullIf
	case FileFormatTypeParquet:
		ff.Options.ParquetTrimSpace = &inputOptions.TrimSpace
		ff.Options.ParquetNullIf = &newNullIf
		ff.Options.ParquetCompression = (*ParquetCompression)(&inputOptions.Compression)
		ff.Options.ParquetBinaryAsText = &inputOptions.BinaryAsText
		ff.Options.ParquetReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
	case FileFormatTypeXml:
		ff.Options.XmlCompression = (*XmlCompression)(&inputOptions.Compression)
		ff.Options.XmlIgnoreUtf8Errors = &inputOptions.IgnoreUtf8Errors
		ff.Options.XmlPreserveSpace = &inputOptions.PreserveSpace
		ff.Options.XmlStripOuterElement = &inputOptions.StripOuterElement
		ff.Options.XmlDisableSnowflakeData = &inputOptions.DisableSnowflakeData
		ff.Options.XmlDisableAutoConvert = &inputOptions.DisableAutoConvert
		ff.Options.XmlReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.XmlSkipByteOrderMark = &inputOptions.SkipByteOrderMark
	}

	return ff
}

type FileFormatType string

const (
	FileFormatTypeCsv     FileFormatType = "CSV"
	FileFormatTypeJson    FileFormatType = "JSON"
	FileFormatTypeAvro    FileFormatType = "AVRO"
	FileFormatTypeOrc     FileFormatType = "ORC"
	FileFormatTypeParquet FileFormatType = "PARQUET"
	FileFormatTypeXml     FileFormatType = "XML"
)

type BinaryFormat string

var (
	BinaryFormatHex    BinaryFormat = "HEX"
	BinaryFormatBase64 BinaryFormat = "BASE64"
	BinaryFormatUtf8   BinaryFormat = "UTF8"
)

type CsvCompression string

var (
	CsvCompressionAuto       CsvCompression = "AUTO"
	CsvCompressionGzip       CsvCompression = "GZIP"
	CsvCompressionBz2        CsvCompression = "BZ2"
	CsvCompressionBrotli     CsvCompression = "BROTLI"
	CsvCompressionZstd       CsvCompression = "ZSTD"
	CsvCompressionDeflate    CsvCompression = "DEFLATE"
	CsvCompressionRawDeflate CsvCompression = "RAW_DEFLATE"
	CsvCompressionNone       CsvCompression = "NONE"
)

type CsvEncoding string

var (
	CsvEncodingBIG5        CsvEncoding = "BIG5"
	CsvEncodingEUCJP       CsvEncoding = "EUCJP"
	CsvEncodingEUCKR       CsvEncoding = "EUCKR"
	CsvEncodingGB18030     CsvEncoding = "GB18030"
	CsvEncodingIBM420      CsvEncoding = "IBM420"
	CsvEncodingIBM424      CsvEncoding = "IBM424"
	CsvEncodingISO2022CN   CsvEncoding = "ISO2022CN"
	CsvEncodingISO2022JP   CsvEncoding = "ISO2022JP"
	CsvEncodingISO2022KR   CsvEncoding = "ISO2022KR"
	CsvEncodingISO88591    CsvEncoding = "ISO88591"
	CsvEncodingISO88592    CsvEncoding = "ISO88592"
	CsvEncodingISO88595    CsvEncoding = "ISO88595"
	CsvEncodingISO88596    CsvEncoding = "ISO88596"
	CsvEncodingISO88597    CsvEncoding = "ISO88597"
	CsvEncodingISO88598    CsvEncoding = "ISO88598"
	CsvEncodingISO88599    CsvEncoding = "ISO88599"
	CsvEncodingISO885915   CsvEncoding = "ISO885915"
	CsvEncodingKOI8R       CsvEncoding = "KOI8R"
	CsvEncodingSHIFTJIS    CsvEncoding = "SHIFTJIS"
	CsvEncodingUTF8        CsvEncoding = "UTF8"
	CsvEncodingUTF16       CsvEncoding = "UTF16"
	CsvEncodingUTF16BE     CsvEncoding = "UTF16BE"
	CsvEncodingUTF16LE     CsvEncoding = "UTF16LE"
	CsvEncodingUTF32       CsvEncoding = "UTF32"
	CsvEncodingUTF32BE     CsvEncoding = "UTF32BE"
	CsvEncodingUTF32LE     CsvEncoding = "UTF32LE"
	CsvEncodingWINDOWS1250 CsvEncoding = "WINDOWS1250"
	CsvEncodingWINDOWS1251 CsvEncoding = "WINDOWS1251"
	CsvEncodingWINDOWS1252 CsvEncoding = "WINDOWS1252"
	CsvEncodingWINDOWS1253 CsvEncoding = "WINDOWS1253"
	CsvEncodingWINDOWS1254 CsvEncoding = "WINDOWS1254"
	CsvEncodingWINDOWS1255 CsvEncoding = "WINDOWS1255"
	CsvEncodingWINDOWS1256 CsvEncoding = "WINDOWS1256"
)

type JsonCompression string

var (
	JsonCompressionAuto       JsonCompression = "AUTO"
	JsonCompressionGzip       JsonCompression = "GZIP"
	JsonCompressionBz2        JsonCompression = "BZ2"
	JsonCompressionBrotli     JsonCompression = "BROTLI"
	JsonCompressionZstd       JsonCompression = "ZSTD"
	JsonCompressionDeflate    JsonCompression = "DEFLATE"
	JsonCompressionRawDeflate JsonCompression = "RAW_DEFLATE"
	JsonCompressionNone       JsonCompression = "NONE"
)

type AvroCompression string

var (
	AvroCompressionAuto       AvroCompression = "AUTO"
	AvroCompressionGzip       AvroCompression = "GZIP"
	AvroCompressionBrotli     AvroCompression = "BROTLI"
	AvroCompressionZstd       AvroCompression = "ZSTD"
	AvroCompressionDeflate    AvroCompression = "DEFLATE"
	AvroCompressionRawDeflate AvroCompression = "RAW_DEFLATE"
	AvroCompressionNone       AvroCompression = "NONE"
)

type ParquetCompression string

var (
	ParquetCompressionAuto   ParquetCompression = "AUTO"
	ParquetCompressionLzo    ParquetCompression = "LZO"
	ParquetCompressionSnappy ParquetCompression = "SNAPPY"
	ParquetCompressionNone   ParquetCompression = "NONE"
)

type XmlCompression string

var (
	XmlCompressionAuto       XmlCompression = "AUTO"
	XmlCompressionGzip       XmlCompression = "GZIP"
	XmlCompressionBz2        XmlCompression = "BZ2"
	XmlCompressionBrotli     XmlCompression = "BROTLI"
	XmlCompressionZstd       XmlCompression = "ZSTD"
	XmlCompressionDeflate    XmlCompression = "DEFLATE"
	XmlCompressionRawDeflate XmlCompression = "RAW_DEFLATE"
	XmlCompressionNone       XmlCompression = "NONE"
)

type NullString struct {
	S string `ddl:"parameter,no_equals,single_quotes"`
}

type CreateFileFormatOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"` //lint:ignore U1000 This is used in the ddl tag
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary   *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	fileFormat  bool                   `ddl:"static" sql:"FILE FORMAT"` //lint:ignore U1000 This is used in the ddl tag
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Type        FileFormatType         `ddl:"parameter" sql:"TYPE"`

	FileFormatTypeOptions
}

func (opts *CreateFileFormatOptions) validate() error {
	fields := opts.FileFormatTypeOptions.fieldsByType()

	for formatType := range fields {
		if opts.Type == formatType {
			continue
		}
		if anyValueSet(fields[formatType]...) {
			return fmt.Errorf("Cannot set %s fields when TYPE = %s", formatType, opts.Type)
		}
	}

	err := opts.FileFormatTypeOptions.validate()
	if err != nil {
		return err
	}

	return nil
}

func (v *fileFormats) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptions) error {
	if opts == nil {
		opts = &CreateFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type AlterFileFormatOptions struct {
	alter      bool                   `ddl:"static" sql:"ALTER"`       //lint:ignore U1000 This is used in the ddl tag
	fileFormat bool                   `ddl:"static" sql:"FILE FORMAT"` //lint:ignore U1000 This is used in the ddl tag
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`

	Rename *AlterFileFormatRenameOptions
	Set    *FileFormatTypeOptions `ddl:"list,no_comma" sql:"SET"`
}

func (opts *AlterFileFormatOptions) validate() error {
	if !exactlyOneValueSet(opts.Rename, opts.Set) {
		return fmt.Errorf("Only one of Rename or Set can be set at once.")
	}
	if valueSet(opts.Set) {
		err := opts.Set.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type AlterFileFormatRenameOptions struct {
	NewName SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
}

type FileFormatTypeOptions struct {
	// CSV type options
	CsvCompression                *CsvCompression `ddl:"parameter" sql:"COMPRESSION"`
	CsvRecordDelimiter            *string         `ddl:"parameter,single_quotes" sql:"RECORD_DELIMITER"`
	CsvFieldDelimiter             *string         `ddl:"parameter,single_quotes" sql:"FIELD_DELIMITER"`
	CsvFileExtension              *string         `ddl:"parameter,single_quotes" sql:"FILE_EXTENSION"`
	CsvParseHeader                *bool           `ddl:"parameter" sql:"PARSE_HEADER"`
	CsvSkipHeader                 *int            `ddl:"parameter" sql:"SKIP_HEADER"`
	CsvSkipBlankLines             *bool           `ddl:"parameter" sql:"SKIP_BLANK_LINES"`
	CsvDateFormat                 *string         `ddl:"parameter,single_quotes" sql:"DATE_FORMAT"`
	CsvTimeFormat                 *string         `ddl:"parameter,single_quotes" sql:"TIME_FORMAT"`
	CsvTimestampFormat            *string         `ddl:"parameter,single_quotes" sql:"TIMESTAMP_FORMAT"`
	CsvBinaryFormat               *BinaryFormat   `ddl:"parameter" sql:"BINARY_FORMAT"`
	CsvEscape                     *string         `ddl:"parameter,single_quotes" sql:"ESCAPE"`
	CsvEscapeUnenclosedField      *string         `ddl:"parameter,single_quotes" sql:"ESCAPE_UNENCLOSED_FIELD"`
	CsvTrimSpace                  *bool           `ddl:"parameter" sql:"TRIM_SPACE"`
	CsvFieldOptionallyEnclosedBy  *string         `ddl:"parameter,single_quotes" sql:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	CsvNullIf                     *[]NullString   `ddl:"parameter,parentheses" sql:"NULL_IF"`
	CsvErrorOnColumnCountMismatch *bool           `ddl:"parameter" sql:"ERROR_ON_COLUMN_COUNT_MISMATCH"`
	CsvReplaceInvalidCharacters   *bool           `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	CsvEmptyFieldAsNull           *bool           `ddl:"parameter" sql:"EMPTY_FIELD_AS_NULL"`
	CsvSkipByteOrderMark          *bool           `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`
	CsvEncoding                   *CsvEncoding    `ddl:"parameter,single_quotes" sql:"ENCODING"`

	// JSON type options
	JsonCompression              *JsonCompression `ddl:"parameter" sql:"COMPRESSION"`
	JsonDateFormat               *string          `ddl:"parameter,single_quotes" sql:"DATE_FORMAT"`
	JsonTimeFormat               *string          `ddl:"parameter,single_quotes" sql:"TIME_FORMAT"`
	JsonTimestampFormat          *string          `ddl:"parameter,single_quotes" sql:"TIMESTAMP_FORMAT"`
	JsonBinaryFormat             *BinaryFormat    `ddl:"parameter" sql:"BINARY_FORMAT"`
	JsonTrimSpace                *bool            `ddl:"parameter" sql:"TRIM_SPACE"`
	JsonNullIf                   *[]NullString    `ddl:"parameter,parentheses" sql:"NULL_IF"`
	JsonFileExtension            *string          `ddl:"parameter,single_quotes" sql:"FILE_EXTENSION"`
	JsonEnableOctal              *bool            `ddl:"parameter" sql:"ENABLE_OCTAL"`
	JsonAllowDuplicate           *bool            `ddl:"parameter" sql:"ALLOW_DUPLICATE"`
	JsonStripOuterArray          *bool            `ddl:"parameter" sql:"STRIP_OUTER_ARRAY"`
	JsonStripNullValues          *bool            `ddl:"parameter" sql:"STRIP_NULL_VALUES"`
	JsonReplaceInvalidCharacters *bool            `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	JsonIgnoreUtf8Errors         *bool            `ddl:"parameter" sql:"IGNORE_UTF8_ERRORS"`
	JsonSkipByteOrderMark        *bool            `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`

	// AVRO type options
	AvroCompression              *AvroCompression `ddl:"parameter" sql:"COMPRESSION"`
	AvroTrimSpace                *bool            `ddl:"parameter" sql:"TRIM_SPACE"`
	AvroReplaceInvalidCharacters *bool            `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	AvroNullIf                   *[]NullString    `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// ORC type options
	OrcTrimSpace                *bool         `ddl:"parameter" sql:"TRIM_SPACE"`
	OrcReplaceInvalidCharacters *bool         `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	OrcNullIf                   *[]NullString `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// PARQUET type options
	ParquetCompression              *ParquetCompression `ddl:"parameter" sql:"COMPRESSION"`
	ParquetSnappyCompression        *bool               `ddl:"parameter" sql:"SNAPPY_COMPRESSION"`
	ParquetBinaryAsText             *bool               `ddl:"parameter" sql:"BINARY_AS_TEXT"`
	ParquetTrimSpace                *bool               `ddl:"parameter" sql:"TRIM_SPACE"`
	ParquetReplaceInvalidCharacters *bool               `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	ParquetNullIf                   *[]NullString       `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// XML type options
	XmlCompression              *XmlCompression `ddl:"parameter" sql:"COMPRESSION"`
	XmlIgnoreUtf8Errors         *bool           `ddl:"parameter" sql:"IGNORE_UTF8_ERRORS"`
	XmlPreserveSpace            *bool           `ddl:"parameter" sql:"PRESERVE_SPACE"`
	XmlStripOuterElement        *bool           `ddl:"parameter" sql:"STRIP_OUTER_ELEMENT"`
	XmlDisableSnowflakeData     *bool           `ddl:"parameter" sql:"DISABLE_SNOWFLAKE_DATA"`
	XmlDisableAutoConvert       *bool           `ddl:"parameter" sql:"DISABLE_AUTO_CONVERT"`
	XmlReplaceInvalidCharacters *bool           `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	XmlSkipByteOrderMark        *bool           `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`

	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *FileFormatTypeOptions) fieldsByType() map[FileFormatType][]any {
	return map[FileFormatType][]any{
		FileFormatTypeCsv: {
			opts.CsvCompression,
			opts.CsvRecordDelimiter,
			opts.CsvFieldDelimiter,
			opts.CsvFileExtension,
			opts.CsvParseHeader,
			opts.CsvSkipHeader,
			opts.CsvSkipBlankLines,
			opts.CsvDateFormat,
			opts.CsvTimeFormat,
			opts.CsvTimestampFormat,
			opts.CsvBinaryFormat,
			opts.CsvEscape,
			opts.CsvEscapeUnenclosedField,
			opts.CsvTrimSpace,
			opts.CsvFieldOptionallyEnclosedBy,
			opts.CsvNullIf,
			opts.CsvErrorOnColumnCountMismatch,
			opts.CsvReplaceInvalidCharacters,
			opts.CsvEmptyFieldAsNull,
			opts.CsvSkipByteOrderMark,
			opts.CsvEncoding,
		},
		FileFormatTypeJson: {
			opts.JsonCompression,
			opts.JsonDateFormat,
			opts.JsonTimeFormat,
			opts.JsonTimestampFormat,
			opts.JsonBinaryFormat,
			opts.JsonTrimSpace,
			opts.JsonNullIf,
			opts.JsonFileExtension,
			opts.JsonEnableOctal,
			opts.JsonAllowDuplicate,
			opts.JsonStripOuterArray,
			opts.JsonStripNullValues,
			opts.JsonReplaceInvalidCharacters,
			opts.JsonIgnoreUtf8Errors,
			opts.JsonSkipByteOrderMark,
		},
		FileFormatTypeAvro: {
			opts.AvroCompression,
			opts.AvroTrimSpace,
			opts.AvroReplaceInvalidCharacters,
			opts.AvroNullIf,
		},
		FileFormatTypeOrc: {
			opts.OrcTrimSpace,
			opts.OrcReplaceInvalidCharacters,
			opts.OrcNullIf,
		},
		FileFormatTypeParquet: {
			opts.ParquetCompression,
			opts.ParquetSnappyCompression,
			opts.ParquetBinaryAsText,
			opts.ParquetTrimSpace,
			opts.ParquetReplaceInvalidCharacters,
			opts.ParquetNullIf,
		},
		FileFormatTypeXml: {
			opts.XmlCompression,
			opts.XmlIgnoreUtf8Errors,
			opts.XmlPreserveSpace,
			opts.XmlStripOuterElement,
			opts.XmlDisableSnowflakeData,
			opts.XmlDisableAutoConvert,
			opts.XmlReplaceInvalidCharacters,
			opts.XmlSkipByteOrderMark,
		},
	}
}

func (opts *FileFormatTypeOptions) validate() error {
	fields := opts.fieldsByType()
	count := 0

	for formatType := range fields {
		if anyValueSet(fields[formatType]...) {
			count += 1
			if count > 1 {
				return fmt.Errorf("Cannot set options for different format types")
			}
		}
	}

	if everyValueSet(opts.CsvParseHeader, opts.CsvSkipHeader) && *opts.CsvParseHeader {
		return fmt.Errorf("ParseHeader and SkipHeader cannot be set simultaneously")
	}

	if everyValueSet(opts.JsonIgnoreUtf8Errors, opts.JsonReplaceInvalidCharacters) && *opts.JsonIgnoreUtf8Errors && *opts.JsonReplaceInvalidCharacters {
		return fmt.Errorf("IgnoreUtf8Errors and ReplaceInvalidCharacters cannot be set simultaneously")
	}

	if everyValueSet(opts.ParquetCompression, opts.ParquetSnappyCompression) && *opts.ParquetSnappyCompression {
		return fmt.Errorf("Compression and SnappyCompression cannot be set simultaneously")
	}

	if everyValueSet(opts.XmlIgnoreUtf8Errors, opts.XmlReplaceInvalidCharacters) && *opts.XmlIgnoreUtf8Errors && *opts.XmlReplaceInvalidCharacters {
		return fmt.Errorf("IgnoreUtf8Errors and ReplaceInvalidCharacters cannot be set simultaneously")
	}

	validEnclosedBy := []string{"null", "'", `"`}
	if valueSet(opts.CsvFieldOptionallyEnclosedBy) && !slices.Contains(validEnclosedBy, *opts.CsvFieldOptionallyEnclosedBy) {
		return fmt.Errorf("CsvFieldOptionallyEnclosedBy must be one of %v", validEnclosedBy)
	}
	return nil
}

func (v *fileFormats) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptions) error {
	if opts == nil {
		opts = &AlterFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type DropFileFormatOptions struct {
	drop       bool                   `ddl:"static" sql:"DROP"`        //lint:ignore U1000 This is used in the ddl tag
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"` //lint:ignore U1000 This is used in the ddl tag
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropFileFormatOptions) validate() error {
	return nil
}

func (v *fileFormats) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptions) error {
	if opts == nil {
		opts = &DropFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type ShowFileFormatsOptions struct {
	show        bool  `ddl:"static" sql:"SHOW"`         //lint:ignore U1000 This is used in the ddl tag
	fileFormats bool  `ddl:"static" sql:"FILE FORMATS"` //lint:ignore U1000 This is used in the ddl tag
	Like        *Like `ddl:"keyword" sql:"LIKE"`
	In          *In   `ddl:"keyword" sql:"IN"`
}

func (opts *ShowFileFormatsOptions) validate() error {
	return nil
}

func (v *fileFormats) Show(ctx context.Context, opts *ShowFileFormatsOptions) ([]*FileFormat, error) {
	if opts == nil {
		opts = &ShowFileFormatsOptions{}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []FileFormatRow
	err = v.client.query(ctx, &rows, sql)
	fileFormats := make([]*FileFormat, len(rows))
	for i, row := range rows {
		fileFormats[i] = row.toFileFormat()
	}
	return fileFormats, err
}

func (v *fileFormats) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error) {
	fileFormats, err := v.client.FileFormats.Show(ctx, &ShowFileFormatsOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	for _, FileFormat := range fileFormats {
		if FileFormat.ID() == id {
			return FileFormat, nil
		}
	}
	return nil, ErrObjectNotExistOrAuthorized
}

type FileFormatDetails struct {
	Type    FileFormatType
	Options FileFormatTypeOptions
}

type FileFormatDetailsRow struct {
	Property         string
	Property_Type    string
	Property_Value   string
	Property_Default string
}

type describeFileFormatOptions struct {
	describe   bool                   `ddl:"static" sql:"DESCRIBE"`    //lint:ignore U1000 This is used in the ddl tag
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"` //lint:ignore U1000 This is used in the ddl tag
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describeFileFormatOptions) validate() error {
	return nil
}

func (v *fileFormats) Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetails, error) {
	opts := &describeFileFormatOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []FileFormatDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	details := FileFormatDetails{}
	for _, row := range rows {
		if row.Property == "TYPE" {
			details.Type = FileFormatType(row.Property_Value)
			break
		}
	}

	switch details.Type {
	case FileFormatTypeCsv:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "RECORD_DELIMITER":
				details.Options.CsvRecordDelimiter = &v
			case "FIELD_DELIMITER":
				details.Options.CsvFieldDelimiter = &v
			case "FILE_EXTENSION":
				details.Options.CsvFileExtension = &v
			case "SKIP_HEADER":
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, v, err)
				}
				i0 := int(i)
				details.Options.CsvSkipHeader = &i0
			case "PARSE_HEADER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvParseHeader = &b
			case "DATE_FORMAT":
				details.Options.CsvDateFormat = &v
			case "TIME_FORMAT":
				details.Options.CsvTimeFormat = &v
			case "TIMESTAMP_FORMAT":
				details.Options.CsvTimestampFormat = &v
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.CsvBinaryFormat = &bf
			case "ESCAPE":
				details.Options.CsvEscape = &v
			case "ESCAPE_UNENCLOSED_FIELD":
				details.Options.CsvEscapeUnenclosedField = &v
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvTrimSpace = &b
			case "FIELD_OPTIONALLY_ENCLOSED_BY":
				details.Options.CsvFieldOptionallyEnclosedBy = &v
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.CsvNullIf = &newNullIf
			case "COMPRESSION":
				comp := CsvCompression(v)
				details.Options.CsvCompression = &comp
			case "ERROR_ON_COLUMN_COUNT_MISMATCH":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvErrorOnColumnCountMismatch = &b
			// case "VALIDATE_UTF8":
			// 	details.Options.C = &v
			case "SKIP_BLANK_LINES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvSkipBlankLines = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvReplaceInvalidCharacters = &b
			case "EMPTY_FIELD_AS_NULL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvEmptyFieldAsNull = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvSkipByteOrderMark = &b
			case "ENCODING":
				enc := CsvEncoding(v)
				details.Options.CsvEncoding = &enc
			}
		}
	case FileFormatTypeJson:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "FILE_EXTENSION":
				details.Options.JsonFileExtension = &v
			case "DATE_FORMAT":
				details.Options.JsonDateFormat = &v
			case "TIME_FORMAT":
				details.Options.JsonTimeFormat = &v
			case "TIMESTAMP_FORMAT":
				details.Options.JsonTimestampFormat = &v
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.JsonBinaryFormat = &bf
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.JsonNullIf = &newNullIf
			case "COMPRESSION":
				comp := JsonCompression(v)
				details.Options.JsonCompression = &comp
			case "ENABLE_OCTAL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonEnableOctal = &b
			case "ALLOW_DUPLICATE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonAllowDuplicate = &b
			case "STRIP_OUTER_ARRAY":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonStripOuterArray = &b
			case "STRIP_NULL_VALUES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonStripNullValues = &b
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonIgnoreUtf8Errors = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonSkipByteOrderMark = &b
			}
		}
	case FileFormatTypeAvro:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.AvroNullIf = &newNullIf
			case "COMPRESSION":
				comp := AvroCompression(v)
				details.Options.AvroCompression = &comp
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeOrc:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.OrcTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.OrcNullIf = &newNullIf
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.OrcReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeParquet:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.ParquetNullIf = &newNullIf
			case "COMPRESSION":
				comp := ParquetCompression(v)
				details.Options.ParquetCompression = &comp
			case "BINARY_AS_TEXT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetBinaryAsText = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeXml:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "COMPRESSION":
				comp := XmlCompression(v)
				details.Options.XmlCompression = &comp
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlIgnoreUtf8Errors = &b
			case "PRESERVE_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlPreserveSpace = &b
			case "STRIP_OUTER_ELEMENT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlStripOuterElement = &b
			case "DISABLE_SNOWFLAKE_DATA":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast DISABLE_SNOWFLAKE_DATA value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlDisableSnowflakeData = &b
			case "DISABLE_AUTO_CONVERT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlDisableAutoConvert = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlSkipByteOrderMark = &b
			}
		}
	default:
		return nil, fmt.Errorf("Describe did not return format type")
	}

	return &details, nil
}
