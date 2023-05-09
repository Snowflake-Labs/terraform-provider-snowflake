package sdk

import (
	"strings"

	"golang.org/x/exp/slices"
)

type DataType string

const (
	DataTypeNumber       DataType = "NUMBER"
	DataTypeFloat        DataType = "FLOAT"
	DataTypeVARCHAR      DataType = "VARCHAR"
	DataTypeBinary       DataType = "BINARY"
	DataTypeBoolean      DataType = "BOOLEAN"
	DataTypeDate         DataType = "DATE"
	DataTypeTime         DataType = "TIME"
	DataTypeTimestampLTZ DataType = "TIMESTAMP_LTZ"
	DataTypeTimestampNTZ DataType = "TIMESTAMP_NTZ"
	DataTypeTimestampTZ  DataType = "TIMESTAMP_TZ"
	DataTypeVariant      DataType = "VARIANT"
	DataTypeObject       DataType = "OBJECT"
	DataTypeArray        DataType = "ARRAY"
	DataTypeGeography    DataType = "GEOGRAPHY"
	DataTypeGeometry     DataType = "GEOMETRY"

	// DataTypeUnknown is used for testing purposes only.
	DataTypeUnknown DataType = "UNKNOWN"
)

func DataTypeFromString(s string) DataType {
	dType := strings.ToUpper(s)

	switch dType {
	case "DATE":
		return DataTypeDate
	case "TIME":
		return DataTypeTime
	case "TIMESTAMP_LTZ":
		return DataTypeTimestampLTZ
	case "TIMESTAMP_TZ":
		return DataTypeTimestampTZ
	case "VARIANT":
		return DataTypeVariant
	case "OBJECT":
		return DataTypeObject
	case "ARRAY":
		return DataTypeArray
	case "GEOGRAPHY":
		return DataTypeGeography
	case "GEOMETRY":
		return DataTypeGeometry
	}

	numberSynonyms := []string{"NUMBER", "DECIMAL", "NUMERIC", "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
	if slices.ContainsFunc(numberSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeNumber
	}

	floatSynonyms := []string{"FLOAT", "FLOAT4", "FLOAT8", "DOUBLE", "DOUBLE PRECISION", "REAL"}
	if slices.ContainsFunc(floatSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeFloat
	}
	varcharSynonyms := []string{"VARCHAR", "CHAR", "CHARACTER", "STRING", "TEXT"}
	if slices.ContainsFunc(varcharSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeVARCHAR
	}
	binarySynonyms := []string{"BINARY", "VARBINARY"}
	if slices.ContainsFunc(binarySynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeBinary
	}
	booleanSynonyms := []string{"BOOLEAN", "BOOL"}
	if slices.Contains(booleanSynonyms, dType) {
		return DataTypeBoolean
	}

	timestampNTZSynonyms := []string{"DATETIME", "TIMESTAMP", "TIMESTAMP_NTZ"}
	if slices.ContainsFunc(timestampNTZSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeTimestampNTZ
	}

	return DataTypeUnknown
}
