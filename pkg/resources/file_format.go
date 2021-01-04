package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"
)

const (
	fileFormatIDDelimiter = '|'
)

var fileFormatSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the file format; must be unique for the database and schema in which the file format is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the file format.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the file format.",
		ForceNew:    true,
	},
	"format_type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the format of the input files (for data loading) or output files (for data unloading).",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"CSV", "JSON", "AVRO", "ORC", "PARQUET", "XML"}, true),
	},
	"compression": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the current compression algorithm for the data file.",
	},
	"record_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies one or more singlebyte or multibyte characters that separate records in an input file (data loading) or unloaded file (data unloading).",
	},
	"field_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies one or more singlebyte or multibyte characters that separate fields in an input file (data loading) or unloaded file (data unloading).",
	},
	"file_extension": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the extension for files unloaded to a stage.",
	},
	"skip_header": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Number of lines at the start of the file to skip.",
	},
	"skip_blank_lines": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies to skip any blank lines encountered in the data files.",
	},
	"date_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of date values in the data files (data loading) or table (data unloading).",
	},
	"time_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of time values in the data files (data loading) or table (data unloading).",
	},
	"timestamp_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of timestamp values in the data files (data loading) or table (data unloading).",
	},
	"binary_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the encoding format for binary input or output.",
	},
	"escape": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for field values.",
	},
	"escape_unenclosed_field": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for unenclosed field values only.",
	},
	"trim_space": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to remove white space from fields.",
	},
	"field_optionally_enclosed_by": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Character used to enclose strings.",
	},
	"null_if": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "String used to convert to and from SQL NULL.",
	},
	"error_on_column_count_mismatch": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to generate a parsing error if the number of delimited columns (i.e. fields) in an input file does not match the number of columns in the corresponding table.",
	},
	"replace_invalid_characters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�).",
	},
	"validate_utf8": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to validate UTF-8 character encoding in string column data.",
	},
	"empty_field_as_null": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to insert SQL NULL for empty fields in an input file, which are represented by two successive delimiters.",
	},
	"skip_byte_order_mark": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to skip the BOM (byte order mark), if present in a data file.",
	},
	"encoding": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "String (constant) that specifies the character set of the source data when loading data into a table.",
	},
	"enable_octal": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that enables parsing of octal numbers.",
	},
	"allow_duplicate": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies to allow duplicate object field names (only the last one will be preserved).",
	},
	"strip_outer_array": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that instructs the JSON parser to remove outer brackets.",
	},
	"strip_null_values": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that instructs the JSON parser to remove object fields or array elements containing null values.",
	},
	"ignore_utf8_errors": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether UTF-8 encoding errors produce error conditions.",
	},
	"snappy_compression": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether unloaded file(s) are compressed using the SNAPPY algorithm.",
	},
	"binary_as_text": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to interpret columns with no defined logical data type as UTF-8 text.",
	},
	"preserve_space": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser preserves leading and trailing spaces in element content.",
	},
	"strip_outer_element": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser strips out the outer XML element, exposing 2nd level elements as separate documents.",
	},
	"disable_snowflake_data": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser disables recognition of Snowflake semi-structured data tags.",
	},
	"disable_auto_convert": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser disables automatic conversion of numeric and Boolean values from text to native representation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the file format.",
	},
}

type fileFormatID struct {
	DatabaseName   string
	SchemaName     string
	FileFormatName string
}

// String() takes in a fileFormatID object and returns a pipe-delimited string:
// DatabaseName|SchemaName|FileFormatName
func (fi *fileFormatID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = fileFormatIDDelimiter
	dataIdentifiers := [][]string{{fi.DatabaseName, fi.SchemaName, fi.FileFormatName}}
	err := csvWriter.WriteAll(dataIdentifiers)
	if err != nil {
		return "", err
	}
	strFileFormatID := strings.TrimSpace(buf.String())
	return strFileFormatID, nil
}

// fileFormatIDFromString() takes in a pipe-delimited string: DatabaseName|SchemaName|FileFormatName
// and returns a fileFormatID object
func fileFormatIDFromString(stringID string) (*fileFormatID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = fileFormatIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line per file format")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("3 fields allowed")
	}

	fileFormatResult := &fileFormatID{
		DatabaseName:   lines[0][0],
		SchemaName:     lines[0][1],
		FileFormatName: lines[0][2],
	}
	return fileFormatResult, nil
}

// FileFormat returns a pointer to the resource representing a file format
func FileFormat() *schema.Resource {
	return &schema.Resource{
		Create: CreateFileFormat,
		Read:   ReadFileFormat,
		Update: UpdateFileFormat,
		Delete: DeleteFileFormat,
		Exists: FileFormatExists,

		Schema: fileFormatSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateFileFormat implements schema.CreateFunc
func CreateFileFormat(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	database := data.Get("database").(string)
	schema := data.Get("schema").(string)

	builder := snowflake.FileFormat(name, database, schema)

	builder.WithFormatType(data.Get("format_type").(string))

	// Set optionals
	if v, ok := data.GetOk("compression"); ok {
		builder.WithCompression(v.(string))
	}

	if v, ok := data.GetOk("record_delimiter"); ok {
		builder.WithRecordDelimiter(v.(string))
	}

	if v, ok := data.GetOk("field_delimiter"); ok {
		builder.WithFieldDelimiter(v.(string))
	}

	if v, ok := data.GetOk("file_extension"); ok {
		builder.WithFileExtension(v.(string))
	}

	if v, ok := data.GetOk("skip_header"); ok {
		builder.WithSkipHeader(v.(int))
	}

	if v, ok := data.GetOk("skip_blank_lines"); ok {
		builder.WithSkipBlankLines(v.(bool))
	}

	if v, ok := data.GetOk("date_format"); ok {
		builder.WithDateFormat(v.(string))
	}

	if v, ok := data.GetOk("time_format"); ok {
		builder.WithTimeFormat(v.(string))
	}

	if v, ok := data.GetOk("timestamp_format"); ok {
		builder.WithTimestampFormat(v.(string))
	}

	if v, ok := data.GetOk("binary_format"); ok {
		builder.WithBinaryFormat(v.(string))
	}

	if v, ok := data.GetOk("escape"); ok {
		builder.WithEscape(v.(string))
	}

	if v, ok := data.GetOk("escape_unenclosed_field"); ok {
		builder.WithEscapeUnenclosedField(v.(string))
	}

	if v, ok := data.GetOk("trim_space"); ok {
		builder.WithTrimSpace(v.(bool))
	}

	if v, ok := data.GetOk("field_optionally_enclosed_by"); ok {
		builder.WithFieldOptionallyEnclosedBy(v.(string))
	}

	if v, ok := data.GetOk("null_if"); ok {
		builder.WithNullIf(expandStringList(v.([]interface{})))
	}

	if v, ok := data.GetOk("error_on_column_count_mismatch"); ok {
		builder.WithErrorOnColumnCountMismatch(v.(bool))
	}

	if v, ok := data.GetOk("replace_invalid_characters"); ok {
		builder.WithReplaceInvalidCharacters(v.(bool))
	}

	if v, ok := data.GetOk("validate_utf8"); ok {
		builder.WithValidateUTF8(v.(bool))
	}

	if v, ok := data.GetOk("empty_field_as_null"); ok {
		builder.WithEmptyFieldAsNull(v.(bool))
	}

	if v, ok := data.GetOk("skip_byte_order_mark"); ok {
		builder.WithSkipByteOrderMark(v.(bool))
	}

	if v, ok := data.GetOk("encoding"); ok {
		builder.WithEncoding(v.(string))
	}

	if v, ok := data.GetOk("enable_octal"); ok {
		builder.WithEnableOctal(v.(bool))
	}

	if v, ok := data.GetOk("allow_duplicate"); ok {
		builder.WithAllowDuplicate(v.(bool))
	}

	if v, ok := data.GetOk("strip_outer_array"); ok {
		builder.WithStripOuterArray(v.(bool))
	}

	if v, ok := data.GetOk("strip_null_values"); ok {
		builder.WithStripNullValues(v.(bool))
	}

	if v, ok := data.GetOk("ignore_utf8_errors"); ok {
		builder.WithIgnoreUTF8Errors(v.(bool))
	}

	if v, ok := data.GetOk("snappy_compression"); ok {
		builder.WithSnappyCompression(v.(bool))
	}

	if v, ok := data.GetOk("binary_as_text"); ok {
		builder.WithBinaryAsText(v.(bool))
	}

	if v, ok := data.GetOk("preserve_space"); ok {
		builder.WithPreserveSpace(v.(bool))
	}

	if v, ok := data.GetOk("strip_outer_element"); ok {
		builder.WithStripOuterElement(v.(bool))
	}

	if v, ok := data.GetOk("disable_snowflake_data"); ok {
		builder.WithDisableSnowflakeData(v.(bool))
	}

	if v, ok := data.GetOk("disable_auto_convert"); ok {
		builder.WithDisableAutoConvert(v.(bool))
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating file format %v", name)
	}

	fileFormatID := &fileFormatID{
		DatabaseName:   database,
		SchemaName:     schema,
		FileFormatName: name,
	}
	dataIDInput, err := fileFormatID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadFileFormat(data, meta)
}

// ReadFileFormat implements schema.ReadFunc
func ReadFileFormat(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	fileFormatID, err := fileFormatIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schema := fileFormatID.SchemaName
	fileFormat := fileFormatID.FileFormatName

	ff := snowflake.FileFormat(fileFormat, dbName, schema).Show()
	row := snowflake.QueryRow(db, ff)

	f, err := snowflake.ScanFileFormatShow(row)
	if err != nil {
		return err
	}

	opts, err := snowflake.ParseFormatOptions(f.FormatOptions.String)
	if err != nil {
		return err
	}

	err = data.Set("name", f.FileFormatName.String)
	if err != nil {
		return err
	}

	err = data.Set("database", f.DatabaseName.String)
	if err != nil {
		return err
	}

	err = data.Set("schema", f.SchemaName.String)
	if err != nil {
		return err
	}

	err = data.Set("format_type", opts.Type)
	if err != nil {
		return err
	}

	err = data.Set("compression", opts.Type)
	if err != nil {
		return err
	}

	err = data.Set("record_delimiter", opts.RecordDelimiter)
	if err != nil {
		return err
	}

	err = data.Set("field_delimiter", opts.FieldDelimiter)
	if err != nil {
		return err
	}

	err = data.Set("file_extension", opts.FileExtension)
	if err != nil {
		return err
	}

	err = data.Set("skip_header", opts.SkipHeader)
	if err != nil {
		return err
	}

	err = data.Set("skip_blank_lines", opts.SkipBlankLines)
	if err != nil {
		return err
	}

	err = data.Set("date_format", opts.DateFormat)
	if err != nil {
		return err
	}

	err = data.Set("time_format", opts.TimeFormat)
	if err != nil {
		return err
	}

	err = data.Set("timestamp_format", opts.TimestampFormat)
	if err != nil {
		return err
	}

	err = data.Set("binary_format", opts.BinaryFormat)
	if err != nil {
		return err
	}

	err = data.Set("escape", opts.Escape)
	if err != nil {
		return err
	}

	err = data.Set("escape_unenclosed_field", opts.EscapeUnenclosedField)
	if err != nil {
		return err
	}

	err = data.Set("trim_space", opts.TrimSpace)
	if err != nil {
		return err
	}

	err = data.Set("field_optionally_enclosed_by", opts.FieldOptionallyEnclosedBy)
	if err != nil {
		return err
	}

	err = data.Set("null_if", opts.NullIf)
	if err != nil {
		return err
	}

	err = data.Set("error_on_column_count_mismatch", opts.ErrorOnColumnCountMismatch)
	if err != nil {
		return err
	}

	err = data.Set("replace_invalid_characters", opts.ReplaceInvalidCharacters)
	if err != nil {
		return err
	}

	err = data.Set("validate_utf8", opts.ValidateUTF8)
	if err != nil {
		return err
	}

	err = data.Set("empty_field_as_null", opts.EmptyFieldAsNull)
	if err != nil {
		return err
	}

	err = data.Set("skip_byte_order_mark", opts.SkipByteOrderMark)
	if err != nil {
		return err
	}

	err = data.Set("encoding", opts.Encoding)
	if err != nil {
		return err
	}

	err = data.Set("enable_octal", opts.EnabelOctal)
	if err != nil {
		return err
	}

	err = data.Set("allow_duplicate", opts.AllowDuplicate)
	if err != nil {
		return err
	}

	err = data.Set("strip_outer_array", opts.StripOuterArray)
	if err != nil {
		return err
	}

	err = data.Set("strip_null_values", opts.StripNullValues)
	if err != nil {
		return err
	}

	err = data.Set("ignore_utf8_errors", opts.IgnoreUTF8Errors)
	if err != nil {
		return err
	}

	err = data.Set("snappy_compression", opts.SnappyCompression)
	if err != nil {
		return err
	}

	err = data.Set("binary_as_text", opts.BinaryAsText)
	if err != nil {
		return err
	}

	err = data.Set("preserve_space", opts.PreserveSpace)
	if err != nil {
		return err
	}

	err = data.Set("strip_outer_element", opts.StripOuterElement)
	if err != nil {
		return err
	}

	err = data.Set("disable_snowflake_data", opts.DisableSnowflakeData)
	if err != nil {
		return err
	}

	err = data.Set("disable_auto_convert", opts.DisableAutoConvert)
	if err != nil {
		return err
	}

	err = data.Set("comment", f.Comment.String)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFileFormat implements schema.UpdateFunc
func UpdateFileFormat(data *schema.ResourceData, meta interface{}) error {
	// https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	data.Partial(true)

	fileFormatID, err := fileFormatIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schema := fileFormatID.SchemaName
	fileFormat := fileFormatID.FileFormatName

	builder := snowflake.FileFormat(fileFormat, dbName, schema)
	fmt.Println(builder)

	db := meta.(*sql.DB)
	if data.HasChange("compression") {
		_, change := data.GetChange("compression")
		q := builder.ChangeCompression(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format compression on %v", data.Id())
		}

		data.SetPartial("compression")
	}

	if data.HasChange("record_delimiter") {
		_, change := data.GetChange("record_delimiter")
		q := builder.ChangeRecordDelimiter(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format record delimiter on %v", data.Id())
		}

		data.SetPartial("record_delimiter")
	}

	if data.HasChange("field_delimiter") {
		_, change := data.GetChange("field_delimiter")
		q := builder.ChangeFieldDelimiter(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format field delimiter on %v", data.Id())
		}

		data.SetPartial("field_delimiter")
	}

	if data.HasChange("file_extension") {
		_, change := data.GetChange("file_extension")
		q := builder.ChangeFileExtension(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format file extension on %v", data.Id())
		}

		data.SetPartial("file_extension")
	}

	if data.HasChange("skip_header") {
		_, change := data.GetChange("skip_header")
		q := builder.ChangeSkipHeader(change.(int))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip header on %v", data.Id())
		}

		data.SetPartial("skip_header")
	}

	if data.HasChange("skip_blank_lines") {
		_, change := data.GetChange("skip_blank_lines")
		q := builder.ChangeSkipBlankLines(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip blank lines on %v", data.Id())
		}

		data.SetPartial("skip_blank_lines")
	}

	if data.HasChange("date_format") {
		_, change := data.GetChange("date_format")
		q := builder.ChangeDateFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format date format on %v", data.Id())
		}

		data.SetPartial("date_format")
	}

	if data.HasChange("time_format") {
		_, change := data.GetChange("time_format")
		q := builder.ChangeTimeFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format time format on %v", data.Id())
		}

		data.SetPartial("time_format")
	}

	if data.HasChange("timestamp_format") {
		_, change := data.GetChange("timestamp_format")
		q := builder.ChangeTimestampFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format timstamp format on %v", data.Id())
		}

		data.SetPartial("timestamp_format")
	}

	if data.HasChange("binary_format") {
		_, change := data.GetChange("binary_format")
		q := builder.ChangeBinaryFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format binary format on %v", data.Id())
		}

		data.SetPartial("binary_format")
	}

	if data.HasChange("escape") {
		_, change := data.GetChange("escape")
		q := builder.ChangeEscape(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format escape on %v", data.Id())
		}

		data.SetPartial("escape")
	}

	if data.HasChange("escape_unenclosed_field") {
		_, change := data.GetChange("escape_unenclosed_field")
		q := builder.ChangeEscapeUnenclosedField(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format escape_unenclosed_field on %v", data.Id())
		}

		data.SetPartial("escape_unenclosed_field")
	}

	if data.HasChange("field_optionally_enclosed_by") {
		_, change := data.GetChange("field_optionally_enclosed_by")
		q := builder.ChangeFieldOptionallyEnclosedBy(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format field_optionally_enclosed_by on %v", data.Id())
		}

		data.SetPartial("field_optionally_enclosed_by")
	}

	if data.HasChange("encoding") {
		_, change := data.GetChange("encoding")
		q := builder.ChangeEncoding(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format encoding on %v", data.Id())
		}

		data.SetPartial("encoding")
	}

	if data.HasChange("comment") {
		_, change := data.GetChange("comment")
		q := builder.ChangeComment(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format comment on %v", data.Id())
		}

		data.SetPartial("comment")
	}

	if data.HasChange("trim_space") {
		_, change := data.GetChange("trim_space")
		q := builder.ChangeTrimSpace(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format trim_space on %v", data.Id())
		}

		data.SetPartial("trim_space")
	}

	if data.HasChange("error_on_column_count_mismatch") {
		_, change := data.GetChange("error_on_column_count_mismatch")
		q := builder.ChangeErrorOnColumnCountMismatch(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format error_on_column_count_mismatch on %v", data.Id())
		}

		data.SetPartial("error_on_column_count_mismatch")
	}

	if data.HasChange("replace_invalid_characters") {
		_, change := data.GetChange("replace_invalid_characters")
		q := builder.ChangeReplaceInvalidCharacters(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format replace_invalid_characters on %v", data.Id())
		}

		data.SetPartial("replace_invalid_characters")
	}

	if data.HasChange("validate_utf8") {
		_, change := data.GetChange("validate_utf8")
		q := builder.ChangeValidateUTF8(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format validate_utf8 on %v", data.Id())
		}

		data.SetPartial("validate_utf8")
	}

	if data.HasChange("empty_field_as_null") {
		_, change := data.GetChange("empty_field_as_null")
		q := builder.ChangeEmptyFieldAsNull(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format empty_field_as_null on %v", data.Id())
		}

		data.SetPartial("empty_field_as_null")
	}

	if data.HasChange("skip_byte_order_mark") {
		_, change := data.GetChange("skip_byte_order_mark")
		q := builder.ChangeSkipByteOrderMark(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip_byte_order_mark on %v", data.Id())
		}

		data.SetPartial("skip_byte_order_mark")
	}

	if data.HasChange("enable_octal") {
		_, change := data.GetChange("enable_octal")
		q := builder.ChangeEnableOctal(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format enable_octal on %v", data.Id())
		}

		data.SetPartial("enable_octal")
	}

	if data.HasChange("allow_duplicate") {
		_, change := data.GetChange("allow_duplicate")
		q := builder.ChangeAllowDuplicate(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format allow_duplicate on %v", data.Id())
		}

		data.SetPartial("allow_duplicate")
	}

	if data.HasChange("strip_outer_array") {
		_, change := data.GetChange("strip_outer_array")
		q := builder.ChangeStripOuterArray(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_outer_array on %v", data.Id())
		}

		data.SetPartial("strip_outer_array")
	}

	if data.HasChange("strip_null_values") {
		_, change := data.GetChange("strip_null_values")
		q := builder.ChangeStripNullValues(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_null_values on %v", data.Id())
		}

		data.SetPartial("strip_null_values")
	}

	if data.HasChange("ignore_utf8_errors") {
		_, change := data.GetChange("ignore_utf8_errors")
		q := builder.ChangeIgnoreUTF8Errors(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format ignore_utf8_errors on %v", data.Id())
		}

		data.SetPartial("ignore_utf8_errors")
	}

	if data.HasChange("snappy_compression") {
		_, change := data.GetChange("snappy_compression")
		q := builder.ChangeSnappyCompression(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format snappy_compression on %v", data.Id())
		}

		data.SetPartial("snappy_compression")
	}

	if data.HasChange("binary_as_text") {
		_, change := data.GetChange("binary_as_text")
		q := builder.ChangeBinaryAsText(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format binary_as_text on %v", data.Id())
		}

		data.SetPartial("binary_as_text")
	}

	if data.HasChange("preserve_space") {
		_, change := data.GetChange("preserve_space")
		q := builder.ChangePreserveSpace(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format preserve_space on %v", data.Id())
		}

		data.SetPartial("preserve_space")
	}

	if data.HasChange("strip_outer_element") {
		_, change := data.GetChange("strip_outer_element")
		q := builder.ChangeStripOuterElement(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_outer_element on %v", data.Id())
		}

		data.SetPartial("strip_outer_element")
	}

	if data.HasChange("disable_snowflake_data") {
		_, change := data.GetChange("disable_snowflake_data")
		q := builder.ChangeDisableSnowflakeData(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format disable_snowflake_data on %v", data.Id())
		}

		data.SetPartial("disable_snowflake_data")
	}

	if data.HasChange("disable_auto_convert") {
		_, change := data.GetChange("disable_auto_convert")
		q := builder.ChangeDisableAutoConvert(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format disable_auto_convert on %v", data.Id())
		}

		data.SetPartial("disable_auto_convert")
	}

	if data.HasChange("null_if") {
		_, change := data.GetChange("null_if")
		q := builder.ChangeNullIf(expandStringList(change.([]interface{})))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format null_if on %v", data.Id())
		}

		data.SetPartial("null_if")
	}

	return ReadFileFormat(data, meta)
}

// DeleteFileFormat implements schema.DeleteFunc
func DeleteFileFormat(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	fileFormatID, err := fileFormatIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schema := fileFormatID.SchemaName
	fileFormat := fileFormatID.FileFormatName

	q := snowflake.FileFormat(fileFormat, dbName, schema).Drop()

	err = snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error deleting file format %v", data.Id())
	}

	data.SetId("")

	return nil
}

// FileFormatExists implements schema.ExistsFunc
func FileFormatExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	fileFormatID, err := fileFormatIDFromString(data.Id())
	if err != nil {
		return false, err
	}

	dbName := fileFormatID.DatabaseName
	schema := fileFormatID.SchemaName
	fileFormat := fileFormatID.FileFormatName

	q := snowflake.FileFormat(fileFormat, dbName, schema).Describe()
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
