package sdk

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/logging"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

// DataType is based on https://docs.snowflake.com/en/sql-reference/intro-summary-data-types.
type DataType string

var allowedVectorInnerTypes = []DataType{
	DataTypeInt,
	DataTypeFloat,
}

const (
	DataTypeNumber       DataType = "NUMBER"
	DataTypeInt          DataType = "INT"
	DataTypeFloat        DataType = "FLOAT"
	DataTypeVARCHAR      DataType = "VARCHAR"
	DataTypeString       DataType = "STRING"
	DataTypeBinary       DataType = "BINARY"
	DataTypeBoolean      DataType = "BOOLEAN"
	DataTypeDate         DataType = "DATE"
	DataTypeTime         DataType = "TIME"
	DataTypeTimestamp    DataType = "TIMESTAMP"
	DataTypeTimestampLTZ DataType = "TIMESTAMP_LTZ"
	DataTypeTimestampNTZ DataType = "TIMESTAMP_NTZ"
	DataTypeTimestampTZ  DataType = "TIMESTAMP_TZ"
	DataTypeVariant      DataType = "VARIANT"
	DataTypeObject       DataType = "OBJECT"
	DataTypeArray        DataType = "ARRAY"
	DataTypeGeography    DataType = "GEOGRAPHY"
	DataTypeGeometry     DataType = "GEOMETRY"
)

var (
	DataTypeNumberSynonyms  = []string{"NUMBER", "DECIMAL", "NUMERIC", "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
	DataTypeVarcharSynonyms = []string{"VARCHAR", "CHAR", "CHARACTER", "STRING", "TEXT"}
)

func IsStringType(_type string) bool {
	t := strings.ToUpper(_type)
	return strings.HasPrefix(t, "STRING") ||
		strings.HasPrefix(t, "VARCHAR") ||
		strings.HasPrefix(t, "CHAR") ||
		strings.HasPrefix(t, "TEXT") ||
		strings.HasPrefix(t, "NVARCHAR") ||
		strings.HasPrefix(t, "NCHAR")
}

const (
	DefaultNumberPrecision = 38
	DefaultNumberScale     = 0
	DefaultVarcharLength   = 16777216
)

// ParseNumberDataTypeRaw extracts precision and scale from the raw number data type input.
// It returns defaults if it can't parse arguments, data type is different, or no arguments were provided.
func ParseNumberDataTypeRaw(rawDataType string) (int, int) {
	r := util.TrimAllPrefixes(strings.TrimSpace(strings.ToUpper(rawDataType)), DataTypeNumberSynonyms...)
	r = strings.TrimSpace(r)
	if strings.HasPrefix(r, "(") && strings.HasSuffix(r, ")") {
		parts := strings.Split(r[1:len(r)-1], ",")
		switch l := len(parts); l {
		case 1:
			precision, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err == nil {
				return precision, DefaultNumberScale
			} else {
				logging.DebugLogger.Printf(`[DEBUG] Could not parse number precision "%s", err: %v`, parts[0], err)
			}
		case 2:
			precision, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			scale, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil {
				return precision, scale
			} else {
				logging.DebugLogger.Printf(`[DEBUG] Could not parse number precision "%s" or scale "%s", errs: %v, %v`, parts[0], parts[1], err1, err2)
			}
		default:
			logging.DebugLogger.Printf("[DEBUG] Unexpected length of number arguments")
		}
	}
	logging.DebugLogger.Printf("[DEBUG] Returning default number precision and scale")
	return DefaultNumberPrecision, DefaultNumberScale
}

// ParseVarcharDataTypeRaw extracts length from the raw text data type input.
// It returns default if it can't parse arguments, data type is different, or no length argument was provided.
func ParseVarcharDataTypeRaw(rawDataType string) int {
	r := util.TrimAllPrefixes(strings.TrimSpace(strings.ToUpper(rawDataType)), DataTypeVarcharSynonyms...)
	r = strings.TrimSpace(r)
	if strings.HasPrefix(r, "(") && strings.HasSuffix(r, ")") {
		parts := strings.Split(r[1:len(r)-1], ",")
		switch l := len(parts); l {
		case 1:
			length, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err == nil {
				return length
			} else {
				logging.DebugLogger.Printf(`[DEBUG] Could not parse varchar length "%s", err: %v`, parts[0], err)
			}
		default:
			logging.DebugLogger.Printf("[DEBUG] Unexpected length of varchar arguments")
		}
	}
	logging.DebugLogger.Printf("[DEBUG] Returning default varchar length")
	return DefaultVarcharLength
}

func LegacyDataTypeFrom(newDataType datatypes.DataType) DataType {
	return DataType(newDataType.ToLegacyDataTypeSql())
}
