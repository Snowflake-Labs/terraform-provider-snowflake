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

	timestampLTZSynonyms := []string{"TIMESTAMP_LTZ"}
	if slices.ContainsFunc(timestampLTZSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeTimestampLTZ
	}

	timestampTZSynonyms := []string{"TIMESTAMP_TZ"}
	if slices.ContainsFunc(timestampTZSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeTimestampTZ
	}

	timestampNTZSynonyms := []string{"DATETIME", "TIMESTAMP", "TIMESTAMP_NTZ"}
	if slices.ContainsFunc(timestampNTZSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeTimestampNTZ
	}

	timeSynonyms := []string{"TIME"}
	if slices.ContainsFunc(timeSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeTime
	}

	return DataTypeUnknown
}
