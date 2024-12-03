package sdk

import (
	"strings"

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

func IsStringType(_type string) bool {
	t := strings.ToUpper(_type)
	return strings.HasPrefix(t, "STRING") ||
		strings.HasPrefix(t, "VARCHAR") ||
		strings.HasPrefix(t, "CHAR") ||
		strings.HasPrefix(t, "TEXT") ||
		strings.HasPrefix(t, "NVARCHAR") ||
		strings.HasPrefix(t, "NCHAR")
}

func LegacyDataTypeFrom(newDataType datatypes.DataType) DataType {
	return DataType(newDataType.ToLegacyDataTypeSql())
}
