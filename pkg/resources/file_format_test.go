package resources_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/resources"
	. "github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestFileFormat(t *testing.T) {
	r := require.New(t)
	err := resources.FileFormat().InternalValidate(provider.Provider().Schema, true)
	r.NoError(err)
}

func TestFileFormatCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":                           "test_file_format",
		"database":                       "test_db",
		"schema":                         "test_schema",
		"format_type":                    "CSV",
		"null_if":                        []interface{}{"NULL"},
		"validate_utf8":                  true,
		"error_on_column_count_mismatch": true,
		"comment":                        "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.FileFormat().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(
			`^CREATE FILE FORMAT "test_db"."test_schema"."test_file_format" TYPE = 'CSV' NULL_IF = \('NULL'\) SKIP_BLANK_LINES = false TRIM_SPACE = false ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = false VALIDATE_UTF8 = true EMPTY_FIELD_AS_NULL = false SKIP_BYTE_ORDER_MARK = false COMMENT = 'great comment'$`,
		).WillReturnResult(sqlmock.NewResult(1, 1))
		expectReadFileFormat(mock)
		err := resources.CreateFileFormat(d, db)
		r.NoError(err)
	})
}

func TestFileFormatCreateInvalidOptions(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name":          "test_file_format",
		"database":      "test_db",
		"schema":        "test_schema",
		"format_type":   "JSON",
		"null_if":       []interface{}{"NULL"},
		"validate_utf8": true,
		"comment":       "great comment",
	}
	d := schema.TestResourceDataRaw(t, resources.FileFormat().Schema, in)
	r.NotNil(d)

	WithMockDb(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		err := resources.CreateFileFormat(d, db)
		r.EqualError(err, "validate_utf8 is an invalid format type option for format type JSON")
	})
}

func expectReadFileFormat(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{
		"created_on", "name", "database_name", "schema_name", "type", "owner", "comment", "format_options"},
	).AddRow("2019-12-23 17:20:50.088 +0000", "test_file_format", "test_db", "test_schema", "CSV", "test", "great comment", `{"TYPE":"CSV","RECORD_DELIMITER":"\n","FIELD_DELIMITER":",","FILE_EXTENSION":null,"SKIP_HEADER":0,"DATE_FORMAT":"AUTO","TIME_FORMAT":"AUTO","TIMESTAMP_FORMAT":"AUTO","BINARY_FORMAT":"HEX","ESCAPE":"NONE","ESCAPE_UNENCLOSED_FIELD":"\\","TRIM_SPACE":false,"FIELD_OPTIONALLY_ENCLOSED_BY":"NONE","NULL_IF":["\\N"],"COMPRESSION":"AUTO","ERROR_ON_COLUMN_COUNT_MISMATCH":false,"VALIDATE_UTF8":false,"SKIP_BLANK_LINES":false,"REPLACE_INVALID_CHARACTERS":false,"EMPTY_FIELD_AS_NULL":false,"SKIP_BYTE_ORDER_MARK":false,"ENCODING":"UTF8"}`)
	mock.ExpectQuery(`^SHOW FILE FORMATS LIKE 'test_file_format' IN DATABASE "test_db"$`).WillReturnRows(rows)
}
