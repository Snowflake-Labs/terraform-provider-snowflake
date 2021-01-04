---
page_title: "snowflake_file_format Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# Resource `snowflake_file_format`





## Schema

### Required

- **database** (String, Required) The database in which to create the file format.
- **format_type** (String, Required) Specifies the format of the input files (for data loading) or output files (for data unloading).
- **name** (String, Required) Specifies the identifier for the file format; must be unique for the database and schema in which the file format is created.
- **schema** (String, Required) The schema in which to create the file format.

### Optional

- **allow_duplicate** (Boolean, Optional) Boolean that specifies to allow duplicate object field names (only the last one will be preserved).
- **binary_as_text** (Boolean, Optional) Boolean that specifies whether to interpret columns with no defined logical data type as UTF-8 text.
- **binary_format** (String, Optional) Defines the encoding format for binary input or output.
- **comment** (String, Optional) Specifies a comment for the file format.
- **compression** (String, Optional) Specifies the current compression algorithm for the data file.
- **date_format** (String, Optional) Defines the format of date values in the data files (data loading) or table (data unloading).
- **disable_auto_convert** (Boolean, Optional) Boolean that specifies whether the XML parser disables automatic conversion of numeric and Boolean values from text to native representation.
- **disable_snowflake_data** (Boolean, Optional) Boolean that specifies whether the XML parser disables recognition of Snowflake semi-structured data tags.
- **empty_field_as_null** (Boolean, Optional) Specifies whether to insert SQL NULL for empty fields in an input file, which are represented by two successive delimiters.
- **enable_octal** (Boolean, Optional) Boolean that enables parsing of octal numbers.
- **encoding** (String, Optional) String (constant) that specifies the character set of the source data when loading data into a table.
- **error_on_column_count_mismatch** (Boolean, Optional) Boolean that specifies whether to generate a parsing error if the number of delimited columns (i.e. fields) in an input file does not match the number of columns in the corresponding table.
- **escape** (String, Optional) Single character string used as the escape character for field values.
- **escape_unenclosed_field** (String, Optional) Single character string used as the escape character for unenclosed field values only.
- **field_delimiter** (String, Optional) Specifies one or more singlebyte or multibyte characters that separate fields in an input file (data loading) or unloaded file (data unloading).
- **field_optionally_enclosed_by** (String, Optional) Character used to enclose strings.
- **file_extension** (String, Optional) Specifies the extension for files unloaded to a stage.
- **id** (String, Optional) The ID of this resource.
- **ignore_utf8_errors** (Boolean, Optional) Boolean that specifies whether UTF-8 encoding errors produce error conditions.
- **null_if** (List of String, Optional) String used to convert to and from SQL NULL.
- **preserve_space** (Boolean, Optional) Boolean that specifies whether the XML parser preserves leading and trailing spaces in element content.
- **record_delimiter** (String, Optional) Specifies one or more singlebyte or multibyte characters that separate records in an input file (data loading) or unloaded file (data unloading).
- **replace_invalid_characters** (Boolean, Optional) Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�).
- **skip_blank_lines** (Boolean, Optional) Boolean that specifies to skip any blank lines encountered in the data files.
- **skip_byte_order_mark** (Boolean, Optional) Boolean that specifies whether to skip the BOM (byte order mark), if present in a data file.
- **skip_header** (Number, Optional) Number of lines at the start of the file to skip.
- **snappy_compression** (Boolean, Optional) Boolean that specifies whether unloaded file(s) are compressed using the SNAPPY algorithm.
- **strip_null_values** (Boolean, Optional) Boolean that instructs the JSON parser to remove object fields or array elements containing null values.
- **strip_outer_array** (Boolean, Optional) Boolean that instructs the JSON parser to remove outer brackets.
- **strip_outer_element** (Boolean, Optional) Boolean that specifies whether the XML parser strips out the outer XML element, exposing 2nd level elements as separate documents.
- **time_format** (String, Optional) Defines the format of time values in the data files (data loading) or table (data unloading).
- **timestamp_format** (String, Optional) Defines the format of timestamp values in the data files (data loading) or table (data unloading).
- **trim_space** (Boolean, Optional) Boolean that specifies whether to remove white space from fields.
- **validate_utf8** (Boolean, Optional) Boolean that specifies whether to validate UTF-8 character encoding in string column data.


