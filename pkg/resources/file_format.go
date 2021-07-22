package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

const (
	fileFormatIDDelimiter = '|'
)

// valid format type options for each File Format Type
// https://docs.snowflake.com/en/sql-reference/sql/create-file-format.html#syntax
var formatTypeOptions = map[string][]string{
	"CSV": {
		"compression",
		"record_delimiter",
		"field_delimiter",
		"file_extension",
		"skip_header",
		"skip_blank_lines",
		"date_format",
		"time_format",
		"timestamp_format",
		"binary_format",
		"escape",
		"escape_unenclosed_field",
		"trim_space",
		"field_optionally_enclosed_by",
		"null_if",
		"error_on_column_count_mismatch",
		"replace_invalid_characters",
		"validate_utf8",
		"empty_field_as_null",
		"skip_byte_order_mark",
		"encoding",
	},
	"JSON": {
		"compression",
		"date_format",
		"time_format",
		"timestamp_format",
		"binary_format",
		"trim_space",
		"null_if",
		"file_extension",
		"enable_octal",
		"allow_duplicate",
		"strip_outer_array",
		"strip_null_values",
		"replace_invalid_characters",
		"ignore_utf8_errors",
		"skip_byte_order_mark",
	},
	"AVRO": {
		"compression",
		"trim_space",
		"null_if",
	},
	"ORC": {
		"trim_space",
		"null_if",
	},
	"PARQUET": {
		"compression",
		"binary_as_text",
		"trim_space",
		"null_if",
	},
	"XML": {
		"compression",
		"ignore_utf8_errors",
		"preserve_space",
		"strip_outer_element",
		"disable_snowflake_data",
		"disable_auto_convert",
		"skip_byte_order_mark",
	},
}

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

func (ffi *fileFormatID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = fileFormatIDDelimiter
	err := csvWriter.WriteAll([][]string{{ffi.DatabaseName, ffi.SchemaName, ffi.FileFormatName}})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
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

	dbName := data.Get("database").(string)
	schemaName := data.Get("schema").(string)
	fileFormatName := data.Get("name").(string)

	builder := snowflake.FileFormat(fileFormatName, dbName, schemaName)

	formatType := data.Get("format_type").(string)
	builder.WithFormatType(formatType)

	// Set optionals
	if v, ok, err := getFormatTypeOption(data, formatType, "compression"); ok && err == nil {
		builder.WithCompression(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "record_delimiter"); ok && err == nil {
		builder.WithRecordDelimiter(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "field_delimiter"); ok && err == nil {
		builder.WithFieldDelimiter(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "file_extension"); ok && err == nil {
		builder.WithFileExtension(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "skip_header"); ok && err == nil {
		builder.WithSkipHeader(v.(int))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "skip_blank_lines"); ok && err == nil {
		builder.WithSkipBlankLines(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "date_format"); ok && err == nil {
		builder.WithDateFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "time_format"); ok && err == nil {
		builder.WithTimeFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "timestamp_format"); ok && err == nil {
		builder.WithTimestampFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "binary_format"); ok && err == nil {
		builder.WithBinaryFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "escape"); ok && err == nil {
		builder.WithEscape(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "escape_unenclosed_field"); ok && err == nil {
		builder.WithEscapeUnenclosedField(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "trim_space"); ok && err == nil {
		builder.WithTrimSpace(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "field_optionally_enclosed_by"); ok && err == nil {
		builder.WithFieldOptionallyEnclosedBy(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "null_if"); ok && err == nil {
		builder.WithNullIf(expandStringList(v.([]interface{})))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "error_on_column_count_mismatch"); ok && err == nil {
		builder.WithErrorOnColumnCountMismatch(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "replace_invalid_characters"); ok && err == nil {
		builder.WithReplaceInvalidCharacters(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "validate_utf8"); ok && err == nil {
		builder.WithValidateUTF8(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "empty_field_as_null"); ok && err == nil {
		builder.WithEmptyFieldAsNull(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "skip_byte_order_mark"); ok && err == nil {
		builder.WithSkipByteOrderMark(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "encoding"); ok && err == nil {
		builder.WithEncoding(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "enable_octal"); ok && err == nil {
		builder.WithEnableOctal(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "allow_duplicate"); ok && err == nil {
		builder.WithAllowDuplicate(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "strip_outer_array"); ok && err == nil {
		builder.WithStripOuterArray(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "strip_null_values"); ok && err == nil {
		builder.WithStripNullValues(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "ignore_utf8_errors"); ok && err == nil {
		builder.WithIgnoreUTF8Errors(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "binary_as_text"); ok && err == nil {
		builder.WithBinaryAsText(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "preserve_space"); ok && err == nil {
		builder.WithPreserveSpace(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "strip_outer_element"); ok && err == nil {
		builder.WithStripOuterElement(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "disable_snowflake_data"); ok && err == nil {
		builder.WithDisableSnowflakeData(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(data, formatType, "disable_auto_convert"); ok && err == nil {
		builder.WithDisableAutoConvert(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok := data.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()

	err := snowflake.Exec(db, q)
	if err != nil {
		return errors.Wrapf(err, "error creating file format %v", fileFormatName)
	}

	fileFormatID := &fileFormatID{
		DatabaseName:   dbName,
		SchemaName:     schemaName,
		FileFormatName: fileFormatName,
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
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	ff := snowflake.FileFormat(fileFormatName, dbName, schemaName).Show()
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

	err = data.Set("compression", opts.Compression)
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
	fileFormatID, err := fileFormatIDFromString(data.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	builder := snowflake.FileFormat(fileFormatName, dbName, schemaName)
	fmt.Println(builder)

	db := meta.(*sql.DB)
	if data.HasChange("compression") {
		change := data.Get("compression")
		q := builder.ChangeCompression(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format compression on %v", data.Id())
		}
	}

	if data.HasChange("record_delimiter") {
		change := data.Get("record_delimiter")
		q := builder.ChangeRecordDelimiter(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format record delimiter on %v", data.Id())
		}
	}

	if data.HasChange("field_delimiter") {
		change := data.Get("field_delimiter")
		q := builder.ChangeFieldDelimiter(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format field delimiter on %v", data.Id())
		}
	}

	if data.HasChange("file_extension") {
		change := data.Get("file_extension")
		q := builder.ChangeFileExtension(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format file extension on %v", data.Id())
		}

	}

	if data.HasChange("skip_header") {
		change := data.Get("skip_header")
		q := builder.ChangeSkipHeader(change.(int))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip header on %v", data.Id())
		}
	}

	if data.HasChange("skip_blank_lines") {
		change := data.Get("skip_blank_lines")
		q := builder.ChangeSkipBlankLines(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip blank lines on %v", data.Id())
		}
	}

	if data.HasChange("date_format") {
		change := data.Get("date_format")
		q := builder.ChangeDateFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format date format on %v", data.Id())
		}
	}

	if data.HasChange("time_format") {
		change := data.Get("time_format")
		q := builder.ChangeTimeFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format time format on %v", data.Id())
		}
	}

	if data.HasChange("timestamp_format") {
		change := data.Get("timestamp_format")
		q := builder.ChangeTimestampFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format timstamp format on %v", data.Id())
		}
	}

	if data.HasChange("binary_format") {
		change := data.Get("binary_format")
		q := builder.ChangeBinaryFormat(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format binary format on %v", data.Id())
		}
	}

	if data.HasChange("escape") {
		change := data.Get("escape")
		q := builder.ChangeEscape(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format escape on %v", data.Id())
		}
	}

	if data.HasChange("escape_unenclosed_field") {
		change := data.Get("escape_unenclosed_field")
		q := builder.ChangeEscapeUnenclosedField(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format escape_unenclosed_field on %v", data.Id())
		}
	}

	if data.HasChange("field_optionally_enclosed_by") {
		change := data.Get("field_optionally_enclosed_by")
		q := builder.ChangeFieldOptionallyEnclosedBy(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format field_optionally_enclosed_by on %v", data.Id())
		}
	}

	if data.HasChange("encoding") {
		change := data.Get("encoding")
		q := builder.ChangeEncoding(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format encoding on %v", data.Id())
		}
	}

	if data.HasChange("comment") {
		change := data.Get("comment")
		q := builder.ChangeComment(change.(string))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format comment on %v", data.Id())
		}
	}

	if data.HasChange("trim_space") {
		change := data.Get("trim_space")
		q := builder.ChangeTrimSpace(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format trim_space on %v", data.Id())
		}
	}

	if data.HasChange("error_on_column_count_mismatch") {
		change := data.Get("error_on_column_count_mismatch")
		q := builder.ChangeErrorOnColumnCountMismatch(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format error_on_column_count_mismatch on %v", data.Id())
		}
	}

	if data.HasChange("replace_invalid_characters") {
		change := data.Get("replace_invalid_characters")
		q := builder.ChangeReplaceInvalidCharacters(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format replace_invalid_characters on %v", data.Id())
		}
	}

	if data.HasChange("validate_utf8") {
		change := data.Get("validate_utf8")
		q := builder.ChangeValidateUTF8(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format validate_utf8 on %v", data.Id())
		}
	}

	if data.HasChange("empty_field_as_null") {
		change := data.Get("empty_field_as_null")
		q := builder.ChangeEmptyFieldAsNull(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format empty_field_as_null on %v", data.Id())
		}
	}

	if data.HasChange("skip_byte_order_mark") {
		change := data.Get("skip_byte_order_mark")
		q := builder.ChangeSkipByteOrderMark(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format skip_byte_order_mark on %v", data.Id())
		}
	}

	if data.HasChange("enable_octal") {
		change := data.Get("enable_octal")
		q := builder.ChangeEnableOctal(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format enable_octal on %v", data.Id())
		}
	}

	if data.HasChange("allow_duplicate") {
		change := data.Get("allow_duplicate")
		q := builder.ChangeAllowDuplicate(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format allow_duplicate on %v", data.Id())
		}
	}

	if data.HasChange("strip_outer_array") {
		change := data.Get("strip_outer_array")
		q := builder.ChangeStripOuterArray(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_outer_array on %v", data.Id())
		}
	}

	if data.HasChange("strip_null_values") {
		change := data.Get("strip_null_values")
		q := builder.ChangeStripNullValues(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_null_values on %v", data.Id())
		}
	}

	if data.HasChange("ignore_utf8_errors") {
		change := data.Get("ignore_utf8_errors")
		q := builder.ChangeIgnoreUTF8Errors(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format ignore_utf8_errors on %v", data.Id())
		}
	}

	if data.HasChange("binary_as_text") {
		change := data.Get("binary_as_text")
		q := builder.ChangeBinaryAsText(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format binary_as_text on %v", data.Id())
		}
	}

	if data.HasChange("preserve_space") {
		change := data.Get("preserve_space")
		q := builder.ChangePreserveSpace(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format preserve_space on %v", data.Id())
		}
	}

	if data.HasChange("strip_outer_element") {
		change := data.Get("strip_outer_element")
		q := builder.ChangeStripOuterElement(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format strip_outer_element on %v", data.Id())
		}
	}

	if data.HasChange("disable_snowflake_data") {
		change := data.Get("disable_snowflake_data")
		q := builder.ChangeDisableSnowflakeData(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format disable_snowflake_data on %v", data.Id())
		}
	}

	if data.HasChange("disable_auto_convert") {
		change := data.Get("disable_auto_convert")
		q := builder.ChangeDisableAutoConvert(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format disable_auto_convert on %v", data.Id())
		}
	}

	if data.HasChange("null_if") {
		change := data.Get("null_if")
		q := builder.ChangeNullIf(expandStringList(change.([]interface{})))
		err := snowflake.Exec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error updating file format null_if on %v", data.Id())
		}
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
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	q := snowflake.FileFormat(fileFormatName, dbName, schemaName).Drop()

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
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	q := snowflake.FileFormat(fileFormatName, dbName, schemaName).Describe()
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

func fileFormatIDFromString(stringID string) (*fileFormatID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = fileFormatIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("4 fields allowed")
	}

	return &fileFormatID{
		DatabaseName:   lines[0][0],
		SchemaName:     lines[0][1],
		FileFormatName: lines[0][2],
	}, nil
}

func getFormatTypeOption(d *schema.ResourceData, formatType, formatTypeOption string) (interface{}, bool, error) {
	validFormatTypeOptions := formatTypeOptions[formatType]
	if v, ok := d.GetOk(formatTypeOption); ok {
		if err := validateFormatTypeOptions(formatType, formatTypeOption, validFormatTypeOptions); err != nil {
			return nil, true, err
		}
		return v, true, nil
	}
	return nil, false, nil
}

func validateFormatTypeOptions(formatType, formatTypeOption string, validFormatTypeOptions []string) error {
	for _, f := range validFormatTypeOptions {
		if f == formatTypeOption {
			return nil
		}
	}
	return fmt.Errorf("%v is an invalid format type option for format type %v", formatTypeOption, formatType)
}
