package resources

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
	if err := csvWriter.WriteAll([][]string{{ffi.DatabaseName, ffi.SchemaName, ffi.FileFormatName}}); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// FileFormat returns a pointer to the resource representing a file format.
func FileFormat() *schema.Resource {
	return &schema.Resource{
		Create: CreateFileFormat,
		Read:   ReadFileFormat,
		Update: UpdateFileFormat,
		Delete: DeleteFileFormat,

		Schema: fileFormatSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateFileFormat implements schema.CreateFunc.
func CreateFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	dbName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	fileFormatName := d.Get("name").(string)

	builder := snowflake.FileFormat(fileFormatName, dbName, schemaName)

	formatType := d.Get("format_type").(string)
	builder.WithFormatType(formatType)

	// Set optionals
	if v, ok, err := getFormatTypeOption(d, formatType, "compression"); ok && err == nil {
		builder.WithCompression(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "record_delimiter"); ok && err == nil {
		builder.WithRecordDelimiter(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "field_delimiter"); ok && err == nil {
		builder.WithFieldDelimiter(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "file_extension"); ok && err == nil {
		builder.WithFileExtension(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "skip_header"); ok && err == nil {
		builder.WithSkipHeader(v.(int))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "skip_blank_lines"); ok && err == nil {
		builder.WithSkipBlankLines(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "date_format"); ok && err == nil {
		builder.WithDateFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "time_format"); ok && err == nil {
		builder.WithTimeFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "timestamp_format"); ok && err == nil {
		builder.WithTimestampFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "binary_format"); ok && err == nil {
		builder.WithBinaryFormat(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "escape"); ok && err == nil {
		builder.WithEscape(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "escape_unenclosed_field"); ok && err == nil {
		builder.WithEscapeUnenclosedField(snowflake.EscapeString(v.(string)))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "trim_space"); ok && err == nil {
		builder.WithTrimSpace(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "field_optionally_enclosed_by"); ok && err == nil {
		builder.WithFieldOptionallyEnclosedBy(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "null_if"); ok && err == nil {
		builder.WithNullIf(expandStringListAllowEmpty(v.([]interface{})))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "error_on_column_count_mismatch"); ok && err == nil {
		builder.WithErrorOnColumnCountMismatch(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "replace_invalid_characters"); ok && err == nil {
		builder.WithReplaceInvalidCharacters(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "empty_field_as_null"); ok && err == nil {
		builder.WithEmptyFieldAsNull(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "skip_byte_order_mark"); ok && err == nil {
		builder.WithSkipByteOrderMark(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "encoding"); ok && err == nil {
		builder.WithEncoding(v.(string))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "enable_octal"); ok && err == nil {
		builder.WithEnableOctal(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "allow_duplicate"); ok && err == nil {
		builder.WithAllowDuplicate(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "strip_outer_array"); ok && err == nil {
		builder.WithStripOuterArray(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "strip_null_values"); ok && err == nil {
		builder.WithStripNullValues(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "ignore_utf8_errors"); ok && err == nil {
		builder.WithIgnoreUTF8Errors(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "binary_as_text"); ok && err == nil {
		builder.WithBinaryAsText(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "preserve_space"); ok && err == nil {
		builder.WithPreserveSpace(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "strip_outer_element"); ok && err == nil {
		builder.WithStripOuterElement(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "disable_snowflake_data"); ok && err == nil {
		builder.WithDisableSnowflakeData(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok, err := getFormatTypeOption(d, formatType, "disable_auto_convert"); ok && err == nil {
		builder.WithDisableAutoConvert(v.(bool))
	} else if err != nil {
		return err
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	q := builder.Create()
	err := snowflake.Exec(db, q)
	if err != nil {
		return fmt.Errorf("error creating file format %v err = %w", fileFormatName, err)
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
	d.SetId(dataIDInput)

	return ReadFileFormat(d, meta)
}

// ReadFileFormat implements schema.ReadFunc.
func ReadFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	ff := snowflake.FileFormat(fileFormatName, dbName, schemaName).Show()
	row := snowflake.QueryRow(db, ff)

	f, err := snowflake.ScanFileFormatShow(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] file_format (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	opts, err := snowflake.ParseFormatOptions(f.FormatOptions.String)
	if err != nil {
		return err
	}

	if err := d.Set("name", f.FileFormatName.String); err != nil {
		return err
	}

	if err := d.Set("database", f.DatabaseName.String); err != nil {
		return err
	}

	if err := d.Set("schema", f.SchemaName.String); err != nil {
		return err
	}

	if err := d.Set("format_type", opts.Type); err != nil {
		return err
	}

	if err := d.Set("compression", opts.Compression); err != nil {
		return err
	}

	if err := d.Set("record_delimiter", opts.RecordDelimiter); err != nil {
		return err
	}

	if err := d.Set("field_delimiter", opts.FieldDelimiter); err != nil {
		return err
	}

	if err := d.Set("file_extension", opts.FileExtension); err != nil {
		return err
	}

	if err := d.Set("skip_header", opts.SkipHeader); err != nil {
		return err
	}

	if err := d.Set("skip_blank_lines", opts.SkipBlankLines); err != nil {
		return err
	}

	if err := d.Set("date_format", opts.DateFormat); err != nil {
		return err
	}

	if err := d.Set("time_format", opts.TimeFormat); err != nil {
		return err
	}

	if err := d.Set("timestamp_format", opts.TimestampFormat); err != nil {
		return err
	}

	if err := d.Set("binary_format", opts.BinaryFormat); err != nil {
		return err
	}

	if err := d.Set("escape", opts.Escape); err != nil {
		return err
	}

	if err := d.Set("escape_unenclosed_field", opts.EscapeUnenclosedField); err != nil {
		return err
	}

	if err := d.Set("trim_space", opts.TrimSpace); err != nil {
		return err
	}

	if err := d.Set("field_optionally_enclosed_by", opts.FieldOptionallyEnclosedBy); err != nil {
		return err
	}

	if err := d.Set("null_if", opts.NullIf); err != nil {
		return err
	}

	if err := d.Set("error_on_column_count_mismatch", opts.ErrorOnColumnCountMismatch); err != nil {
		return err
	}

	if err := d.Set("replace_invalid_characters", opts.ReplaceInvalidCharacters); err != nil {
		return err
	}

	if err := d.Set("empty_field_as_null", opts.EmptyFieldAsNull); err != nil {
		return err
	}

	if err := d.Set("skip_byte_order_mark", opts.SkipByteOrderMark); err != nil {
		return err
	}

	if err := d.Set("encoding", opts.Encoding); err != nil {
		return err
	}

	if err := d.Set("enable_octal", opts.EnabelOctal); err != nil {
		return err
	}

	if err := d.Set("allow_duplicate", opts.AllowDuplicate); err != nil {
		return err
	}

	if err := d.Set("strip_outer_array", opts.StripOuterArray); err != nil {
		return err
	}

	if err := d.Set("strip_null_values", opts.StripNullValues); err != nil {
		return err
	}

	if err := d.Set("ignore_utf8_errors", opts.IgnoreUTF8Errors); err != nil {
		return err
	}

	if err := d.Set("binary_as_text", opts.BinaryAsText); err != nil {
		return err
	}

	if err := d.Set("preserve_space", opts.PreserveSpace); err != nil {
		return err
	}

	if err := d.Set("strip_outer_element", opts.StripOuterElement); err != nil {
		return err
	}

	if err := d.Set("disable_snowflake_data", opts.DisableSnowflakeData); err != nil {
		return err
	}

	if err := d.Set("disable_auto_convert", opts.DisableAutoConvert); err != nil {
		return err
	}

	if err := d.Set("comment", f.Comment.String); err != nil {
		return err
	}
	return nil
}

// UpdateFileFormat implements schema.UpdateFunc.
func UpdateFileFormat(d *schema.ResourceData, meta interface{}) error {
	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	builder := snowflake.FileFormat(fileFormatName, dbName, schemaName)
	fmt.Println(builder)

	db := meta.(*sql.DB)
	if d.HasChange("compression") {
		change := d.Get("compression")
		q := builder.ChangeCompression(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format compression on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("record_delimiter") {
		change := d.Get("record_delimiter")
		q := builder.ChangeRecordDelimiter(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format record delimiter on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("field_delimiter") {
		change := d.Get("field_delimiter")
		q := builder.ChangeFieldDelimiter(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format field delimiter on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("file_extension") {
		change := d.Get("file_extension")
		q := builder.ChangeFileExtension(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format file extension on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("skip_header") {
		change := d.Get("skip_header")
		q := builder.ChangeSkipHeader(change.(int))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format skip header on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("skip_blank_lines") {
		change := d.Get("skip_blank_lines")
		q := builder.ChangeSkipBlankLines(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format skip blank lines on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("date_format") {
		change := d.Get("date_format")
		q := builder.ChangeDateFormat(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format date format on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("time_format") {
		change := d.Get("time_format")
		q := builder.ChangeTimeFormat(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format time format on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("timestamp_format") {
		change := d.Get("timestamp_format")
		q := builder.ChangeTimestampFormat(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format timestamp format on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("binary_format") {
		change := d.Get("binary_format")
		q := builder.ChangeBinaryFormat(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format binary format on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("escape") {
		change := d.Get("escape")
		q := builder.ChangeEscape(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format escape on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("escape_unenclosed_field") {
		change := d.Get("escape_unenclosed_field")
		q := builder.ChangeEscapeUnenclosedField(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format escape_unenclosed_field on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("field_optionally_enclosed_by") {
		change := d.Get("field_optionally_enclosed_by")
		q := builder.ChangeFieldOptionallyEnclosedBy(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format field_optionally_enclosed_by on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("encoding") {
		change := d.Get("encoding")
		q := builder.ChangeEncoding(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format encoding on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("comment") {
		change := d.Get("comment")
		q := builder.ChangeComment(change.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format comment on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("trim_space") {
		change := d.Get("trim_space")
		q := builder.ChangeTrimSpace(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return fmt.Errorf("error updating file format trim_space on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("error_on_column_count_mismatch") {
		change := d.Get("error_on_column_count_mismatch")
		q := builder.ChangeErrorOnColumnCountMismatch(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format error_on_column_count_mismatch on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("replace_invalid_characters") {
		change := d.Get("replace_invalid_characters")
		q := builder.ChangeReplaceInvalidCharacters(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format replace_invalid_characters on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("empty_field_as_null") {
		change := d.Get("empty_field_as_null")
		q := builder.ChangeEmptyFieldAsNull(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format empty_field_as_null on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("skip_byte_order_mark") {
		change := d.Get("skip_byte_order_mark")
		q := builder.ChangeSkipByteOrderMark(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format skip_byte_order_mark on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("enable_octal") {
		change := d.Get("enable_octal")
		q := builder.ChangeEnableOctal(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format enable_octal on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("allow_duplicate") {
		change := d.Get("allow_duplicate")
		q := builder.ChangeAllowDuplicate(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format allow_duplicate on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("strip_outer_array") {
		change := d.Get("strip_outer_array")
		q := builder.ChangeStripOuterArray(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format strip_outer_array on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("strip_null_values") {
		change := d.Get("strip_null_values")
		q := builder.ChangeStripNullValues(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format strip_null_values on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("ignore_utf8_errors") {
		change := d.Get("ignore_utf8_errors")
		q := builder.ChangeIgnoreUTF8Errors(change.(bool))
		err := snowflake.Exec(db, q)
		if err != nil {
			return fmt.Errorf("error updating file format ignore_utf8_errors on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("binary_as_text") {
		change := d.Get("binary_as_text")
		q := builder.ChangeBinaryAsText(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format binary_as_text on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("preserve_space") {
		change := d.Get("preserve_space")
		q := builder.ChangePreserveSpace(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format preserve_space on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("strip_outer_element") {
		change := d.Get("strip_outer_element")
		q := builder.ChangeStripOuterElement(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format strip_outer_element on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("disable_snowflake_data") {
		change := d.Get("disable_snowflake_data")
		q := builder.ChangeDisableSnowflakeData(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format disable_snowflake_data on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("disable_auto_convert") {
		change := d.Get("disable_auto_convert")
		q := builder.ChangeDisableAutoConvert(change.(bool))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format disable_auto_convert on %v err = %w", d.Id(), err)
		}
	}

	if d.HasChange("null_if") {
		change := d.Get("null_if")
		q := builder.ChangeNullIf(expandStringListAllowEmpty(change.([]interface{})))
		if err := snowflake.Exec(db, q); err != nil {
			return fmt.Errorf("error updating file format null_if on %v err = %w", d.Id(), err)
		}
	}

	return ReadFileFormat(d, meta)
}

// DeleteFileFormat implements schema.DeleteFunc.
func DeleteFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := fileFormatID.DatabaseName
	schemaName := fileFormatID.SchemaName
	fileFormatName := fileFormatID.FileFormatName

	q := snowflake.FileFormat(fileFormatName, dbName, schemaName).Drop()
	if err := snowflake.Exec(db, q); err != nil {
		return fmt.Errorf("error deleting file format %v err = %w", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func fileFormatIDFromString(stringID string) (*fileFormatID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = fileFormatIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
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
