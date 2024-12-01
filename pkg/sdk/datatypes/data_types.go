package datatypes

import (
	"fmt"
	"slices"
	"strings"
)

type DataType interface {
}

type sanitizedDataTypeRaw struct {
	raw           string
	matchedByType string
}

// TODO [this PR]: test
// TODO [this PR]: support all data types
// https://docs.snowflake.com/en/sql-reference/intro-summary-data-types
func ParseDataType(raw string) (DataType, error) {
	dataTypeRaw := strings.TrimSpace(strings.ToUpper(raw))

	if idx := slices.IndexFunc(AllNumberDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseNumberDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllNumberDataTypes[idx]})
	}
	if idx := slices.Index(FloatDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseFloatDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, FloatDataTypeSynonyms[idx]})
	}
	if idx := slices.IndexFunc(AllTextDataTypes, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseTextDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, AllTextDataTypes[idx]})
	}
	if idx := slices.IndexFunc(BinaryDataTypeSynonyms, func(s string) bool { return strings.HasPrefix(dataTypeRaw, s) }); idx >= 0 {
		return parseBinaryDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, BinaryDataTypeSynonyms[idx]})
	}
	if idx := slices.Index(BooleanDataTypeSynonyms, dataTypeRaw); idx >= 0 {
		return parseBooleanDataTypeRaw(sanitizedDataTypeRaw{dataTypeRaw, BooleanDataTypeSynonyms[idx]})
	}
	return nil, fmt.Errorf("invalid data type: %s", raw)
}

// TODO [this PR]: support all data types
type TimestampLTZDataType struct{}
type TimestampTZDataType struct{}
type TimestampNTZDataType struct{}
type TimeDataType struct{}
type VectorDataType struct{}
