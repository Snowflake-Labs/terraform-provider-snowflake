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
)

func DataTypeFromString(s string) (DataType, error) {
	dType := strings.ToUpper(s)
	numberSynonyms := []string{"NUMBER", "DECIMAL", "NUMERIC", "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "BYTEINT"}
	if slices.ContainsFunc(numberSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeNumber, nil
	}

	floatSynonyms := []string{"FLOAT", "FLOAT4", "FLOAT8", "DOUBLE", "DOUBLE PRECISION", "REAL"}
	if slices.ContainsFunc(floatSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeFloat, nil
	}
	varcharSynonyms := []string{"VARCHAR", "CHAR", "CHARACTER", "STRING", "TEXT"}
	if slices.ContainsFunc(varcharSynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeVARCHAR, nil
	}
	binarySynonyms := []string{"BINARY", "VARBINARY"}
	if slices.ContainsFunc(binarySynonyms, func(s string) bool { return strings.HasPrefix(dType, s) }) {
		return DataTypeBinary, nil
	}
	booleanSynonyms := []string{"BOOLEAN", "BOOL"}
	if slices.Contains(booleanSynonyms, dType) {
		return DataTypeBoolean, nil
	}
	switch dType {
	case "DATE":
		return DataTypeDate, nil
	case "DATETIME":
		return DataTypeTimestampNTZ, nil
	case "TIME":
		return DataTypeTime, nil
	case "TIMESTAMP":
		return DataTypeTimestampNTZ, nil
	case "TIMESTAMP_LTZ":
		return DataTypeTimestampLTZ, nil
	case "TIMESTAMP_NTZ":
		return DataTypeTimestampNTZ, nil
	case "TIMESTAMP_TZ":
		return DataTypeTimestampTZ, nil
	case "VARIANT":
		return DataTypeVariant, nil
	case "OBJECT":
		return DataTypeObject, nil
	case "ARRAY":
		return DataTypeArray, nil
	case "GEOGRAPHY":
		return DataTypeGeography, nil
	case "GEOMETRY":
		return DataTypeGeometry, nil
	}
	return "", ErrInvalidDataType
}
